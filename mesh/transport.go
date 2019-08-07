package mesh

import (
	"ble-mesh/mesh/crypto"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"time"

	funk "github.com/thoas/go-funk"
)

var (
	tpRxChan       chan *NetworkMessage
	tpTxChan       chan *NetworkMessage
	segAckChan     chan *SegmentAckMessage
	segMsgTimeouts map[*NetworkMessage]*AckTimeout
	loggerTp       = utils.CreateLogger("tp")
	tpSars         map[uint64]*TpSar
)

type (
	AckTimeout struct {
		remainingTries int
		timeOut        time.Time
		duration       time.Duration
		segO           int
		fullmask       uint
	}

	TpSar struct {
		cipher    []byte
		recvdSegs []int
		segN      int
		nid       uint
		id        uint64
	}

	SegmentAckMessage struct {
		src, obo, seqZero, blockAck uint
	}

	onAckReceived func(ack *SegmentAckMessage)
)

func (s *TpSar) getRecvdFlags() uint {
	var r uint
	for _, s := range s.recvdSegs {
		r |= 0x00000001 << uint(s)
	}
	return r
}

func calcSarId(src, dst, seqAuth uint) uint64 {
	return uint64(src)<<48 | uint64(dst) | uint64(seqAuth)&0xFFFFFFFF
}

const (
	SEGMENT_SIZE      = 12
	MAX_TRANSPORT_PDU = 15
)

func transportReceive(netMsg *NetworkMessage) {
	tpRxChan <- netMsg
}

func transportRxProc() {
	for {
		netMsg, more := <-tpRxChan
		if !more {
			return
		}
		var err error
		if netMsg == nil {
			// duplicated message
			continue
		}
		if netMsg.ctl == 0 {
			err = tpHandleAccessMessageRx(netMsg)
		} else {
			err = tpHandleControlMessageRx(netMsg)
		}
		if err != nil {
			loggerTp.Error(err)
		}
	}
}

func tpHandleAccessMessageRx(netMsg *NetworkMessage) error {
	var accessPdu []byte
	var seg, akf, aid, szmic uint
	err := utils.UnpackBE(netMsg.plain[:1], "1,1,6", &seg, &akf, &aid)
	if err != nil {
		return err
	}
	if seg == 1 {
		var seqZero uint
		var segO, segN int
		err := utils.UnpackBE(netMsg.plain[1:4], "1,13,5,5", &szmic, &seqZero, &segO, &segN)
		if err != nil {
			return err
		}
		seqAuth := seqAuth(netMsg.seq, seqZero)
		// calculate a unique id to store the sar
		sarId := calcSarId(netMsg.src, netMsg.dst, seqAuth)
		tpSar := tpSars[sarId]
		if tpSar == nil {
			// it's a new fresh message
			tpSar = &TpSar{recvdSegs: []int{}, id: sarId, nid: netMsg.nid, segN: segN}
			go func() {
				time.Sleep(10 * time.Second)
				if tpSar != nil {
					delete(tpSars, tpSar.id)
				}
			}()
			tpSars[sarId] = tpSar
		}
		/* Sanity Check--> certain things must match */
		if tpSar.segN != segN || tpSar.nid != netMsg.nid {
			return errors.TransportSarFailed.New()
		}
		if funk.ContainsInt(tpSar.recvdSegs, segO) {
			// if segO already recevied, just send the ack
			return tpSendSegAck(netMsg.netKey, netMsg.dst, netMsg.src, seqZero, tpSar.getRecvdFlags())
		}
		// update recevied segments
		tpSar.recvdSegs = append(tpSar.recvdSegs, segO)
		// send ack
		err = tpSendSegAck(netMsg.netKey, netMsg.dst, netMsg.src, seqZero, tpSar.getRecvdFlags())
		segment := netMsg.plain[4:]
		if tpSar.cipher == nil {
			tpSar.cipher = make([]byte, SEGMENT_SIZE*(segN+1), SEGMENT_SIZE*(segN+1))
		}
		for i, d := range segment {
			tpSar.cipher[segO*SEGMENT_SIZE+i] = d
		}
		if segO == segN {
			// final segment
			tpSar.cipher = tpSar.cipher[:SEGMENT_SIZE*int(segN)+len(segment)]
		}
		if len(tpSar.recvdSegs) == int(segN+1) {
			// all received
			loggerTp.Debugf("TP assembly: %v", tpSar.cipher)
			accessPdu, err = tpDecryptMessage(tpSar.cipher, netMsg.src, netMsg.dst, seqAuth, netMsg.ivIndex, szmic, akf, aid)
			if err != nil {
				return err
			}
		} else {
			// wait for next segment
			return nil
		}
		// delete the entry for some delay to ignore duplicated messages
		go func() {
			time.Sleep(3 * time.Second)
			if tpSar != nil {
				delete(tpSars, tpSar.id)
			}
		}()
	} else {
		szmic = 0
		accessPdu, err = tpDecryptMessage(netMsg.plain[1:], netMsg.src, netMsg.dst,
			netMsg.seq, netMsg.ivIndex, szmic, akf, aid)
		if err != nil {
			return err
		}
	}
	loggerTp.Debugf("TP RX: % 2x", accessPdu)
	modelMessageReceive(netMsg.src, netMsg.dst, accessPdu)
	return nil
}

func tpHandleControlMessageRx(netMsg *NetworkMessage) error {
	var seg, opcode uint
	err := utils.UnpackBE(netMsg.plain[:1], "1,7", &seg, &opcode)
	if err != nil {
		return err
	}
	if seg == 0 {
		if opcode == 0 {
			var obo, seqZero, blockAck uint
			err := utils.UnpackBE(netMsg.plain[1:7], "1,13,02,32", &obo, &seqZero, &blockAck)
			if err != nil {
				return err
			}
			ackMsg := &SegmentAckMessage{
				src:      netMsg.src,
				obo:      obo,
				seqZero:  seqZero,
				blockAck: blockAck,
			}
			loggerTp.Debugf("received ACK: %+#v", ackMsg)
			segAckChan <- ackMsg
		}
	}
	return nil
}

func tpDecryptMessage(cipher []byte, src, dst, seq, ivIndex, szmic, akf, aid uint) ([]byte, error) {
	for _, key := range findAppKeyByAid(aid) {
		function := genApplicationNonce
		if akf == 0 {
			if node, _ := findNodeByAddr(src); node != nil {
				key = &node.DeviceKey.AppKey
				function = genDeviceNonce
			}
		}
		if key != nil && key.Bytes != nil {
			tagSize := 4
			if szmic == 1 {
				tagSize = 8
			}
			nonce, err := function(src, seq, ivIndex, szmic, dst)
			accessPlain, err := crypto.AES_CCM_Decrypt(
				key.Bytes,
				nonce,
				cipher,
				tagSize)
			if err == nil {
				return accessPlain, err
			}
		}
	}
	return nil, errors.NoValidAppKeyForDecryption.New()
}

func seqAuth(seq, seqZero uint) uint {
	mask := uint(0x1FFF)
	seqAuth := seqZero & mask

	seqAuth |= seq &^ mask
	if seqAuth > seq {
		seqAuth -= (mask + 1)
	}

	return seqAuth
}

func transportTxProc() {
	for {
		netMsg, more := <-tpTxChan
		if !more {
			return
		}
		time.Sleep(time.Millisecond * 10)
		netMsg.seq = meshDb.SequenceNumber
		meshDb.SequenceNumber++
		networkSend(netMsg)
	}

}

func transportSegmentTimeoutProc() {
	for {
		to := time.NewTimer(time.Millisecond * 10)
		select {
		case ack, more := <-segAckChan:
			if !more {
				return
			}
			for netMsg, meta := range segMsgTimeouts {
				mask := uint(1) << uint(meta.segO)
				if netMsg.dst == ack.src {
					// todo: check seqZero
					if mask&ack.blockAck != 0 {
						loggerTp.Debugf("ack for segment %d received", meta.segO)
						delete(segMsgTimeouts, netMsg)
					} else {
						tpTxChan <- netMsg
						meta.remainingTries--
						meta.timeOut = time.Now().Add(meta.duration)
					}
					if ack.blockAck == meta.fullmask {
						netMsg.ackFunc(ack)
					}
				}
			}
		case <-to.C:
		}
		for netMsg, meta := range segMsgTimeouts {
			if time.Now().After(meta.timeOut) {
				if meta.remainingTries == 0 {
					delete(segMsgTimeouts, netMsg)
					continue
				}
				loggerTp.Debugf("timeout occurs, resend segment %d", meta.segO)
				tpTxChan <- netMsg
				meta.remainingTries--
				meta.timeOut = time.Now().Add(meta.duration)
			}
		}
	}
}

func tpSendAccessMsgWithAppKey(appKey *AppKey, netKey *NetKey, payload []byte, dst uint, ttl uint) error {
	_, err := findNodeByAddr(dst)
	if err != nil {
		return err
	}
	return tpSendAccessMsg(1, appKey, netKey, payload, dst, ttl, nil)
}

func tpSendAccessMsgWithDevKey(payload []byte, dst uint, ttl uint, ackFunc onAckReceived) error {
	node, err := findNodeByAddr(dst)
	if err != nil {
		return err
	}
	netKey, err := node.findNodeNetKeyByIndex(node.BindedKeys[0].NetKeyIndex)
	if err != nil {
		return err
	}
	return tpSendAccessMsg(0, &node.DeviceKey.AppKey, netKey, payload, dst, ttl, ackFunc)
}

func tpSendAccessMsg(akf uint, appKey *AppKey, netKey *NetKey, payload []byte, dst uint, ttl uint, ackFunc onAckReceived) error {
	if appKey == nil || netKey == nil {
		return errors.NilPointer.New()
	}
	var szmic, seg, ctl uint
	seq := meshDb.SequenceNumber
	function := genDeviceNonce
	if akf == 1 {
		function = genApplicationNonce
	}
	transMic := 4
	if szmic == 1 {
		transMic = 8
	}
	nonce, err := function(
		meshDb.UnicastAddress,
		seq,
		meshDb.IVindex,
		szmic,
		dst,
	)
	if err != nil {
		return err
	}
	cipher, mic, err := crypto.AES_CCM(appKey.Bytes, nonce, payload, transMic)
	if err != nil {
		return err
	}
	upperTpPdu := append(cipher, mic...)
	ivIndex := meshDb.IVindex
	if meshDb.IVupdate == 1 {
		ivIndex--
	}
	if len(upperTpPdu) > MAX_TRANSPORT_PDU {
		msgs := []*NetworkMessage{}
		seg = 1
		segments := funk.Chunk(upperTpPdu, SEGMENT_SIZE).([][]byte)
		seqZero := seq & 0x1fff
		segN := len(segments) - 1
		for segO, segment := range segments {
			pdu, err := utils.PackBE("1,1,6,1,13,5,5,B8", seg, akf, appKey.Aid, szmic, seqZero, segO, segN, segment)
			if err != nil {
				return err
			}
			netMsg := &NetworkMessage{
				ivi:     ivIndex & 0x01,
				nid:     netKey.Nid,
				ctl:     ctl,
				ttl:     ttl,
				src:     meshDb.UnicastAddress,
				dst:     dst,
				plain:   pdu,
				ivIndex: ivIndex,
				netKey:  netKey,
				ackFunc: ackFunc,
			}
			msgs = append(msgs, netMsg)
			go func(index int) {
				time.Sleep(time.Duration(200*index) * time.Millisecond)
				tpTxChan <- netMsg
				duration := time.Duration(550+50*ttl) * time.Millisecond
				segMsgTimeouts[netMsg] = &AckTimeout{
					remainingTries: 3,
					segO:           index,
					duration:       duration,
					timeOut:        time.Now().Add(duration),
					fullmask:       (1 << uint(segN)) - 1,
				}
			}(segO)
		}
	} else {
		pdu, err := utils.PackBE("1,1,6,B8", seg, akf, appKey.Aid, upperTpPdu)
		if err != nil {
			return err
		}
		netMsg := &NetworkMessage{
			ivi:     ivIndex & 0x01,
			nid:     netKey.Nid,
			ctl:     ctl,
			ttl:     ttl,
			seq:     seq,
			src:     meshDb.UnicastAddress,
			dst:     dst,
			plain:   pdu,
			ivIndex: ivIndex,
			netKey:  netKey,
		}
		tpTxChan <- netMsg
	}
	return nil
}

// SEG		1	0 = Unsegmented Message
// Opcode	7	0x00 = Segment Acknowledgment Message
// OBO		1	Friend on behalf of a Low Power node
// SeqZero	13	SeqZero of the Upper Transport PDU
// RFU		2	Reserved for Future Use
// BlockAck	32	Block acknowledgment for segments
func tpSendSegAck(key *NetKey, src, dst, seqZero, recvdFlags uint) error {

	/* We don't ACK from multicast destinations */
	if src&0x8000 > 0 {
		return nil
	}
	var seg, opCode, obo int
	pdu, err := utils.PackBE("1,7,1,13,02,32", seg, opCode, obo, seqZero, recvdFlags)
	if err != nil {
		return err
	}
	ivIndex := meshDb.IVindex
	if meshDb.IVupdate == 1 {
		ivIndex--
	}
	/* We don't ACK segments as a Low Power Node */

	/* If we are acking our LPN Friend, queue, don't send */
	loggerTp.Debugf("TP ACK TX: % 2x", pdu)
	netMsg := &NetworkMessage{
		ivi:     ivIndex & 0x01,
		nid:     key.Nid,
		ctl:     1,
		ttl:     5,
		src:     meshDb.UnicastAddress,
		dst:     dst,
		plain:   pdu,
		ivIndex: ivIndex,
		netKey:  key,
	}
	tpTxChan <- netMsg
	return nil
}

func startTransport() {
	tpSars = map[uint64]*TpSar{}
	tpRxChan = make(chan *NetworkMessage)
	tpTxChan = make(chan *NetworkMessage, 10)
	segAckChan = make(chan *SegmentAckMessage)
	segMsgTimeouts = map[*NetworkMessage]*AckTimeout{}
	go transportRxProc()
	go transportTxProc()
	go transportSegmentTimeoutProc()
}

func stopTransport() {
	close(tpRxChan)
	close(tpTxChan)
	close(segAckChan)
}
