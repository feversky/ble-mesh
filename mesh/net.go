package mesh

import (
	"ble-mesh/mesh/crypto"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"bytes"
	"container/list"
)

type NetworkMessage struct {
	ivi     uint
	nid     uint
	ctl     uint
	ttl     uint
	seq     uint
	src     uint
	dst     uint
	plain   []byte
	ivIndex uint
	netKey  *NetKey
	ackFunc onAckReceived
}

const cacheSize = 50

var (
	netRxChan chan []byte
	loggerNet = utils.CreateLogger("Net")
	cache     = list.New()
)

func networkReceive(proxyPdu []byte) {
	// loggerNet.Debugf("bear Rx: %+#v", proxyPdu)
	netRxChan <- proxyPdu
}

func netRxProc() {
	for {
		netPdu, more := <-netRxChan
		if !more {
			return
		}
		netMsg, err := networkUnpack(netPdu)
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		if netMsg == nil || netMsg.src == meshDb.UnicastAddress {
			continue
		}
		if netMsg.dst == meshDb.UnicastAddress || isVirtualAddr(netMsg.dst) || isGroupAddr(netMsg.dst) {
			transportReceive(netMsg)
		}
	}
}

func startNet() {
	netRxChan = make(chan []byte)
	go netRxProc()
}

func stopNet() {
	close(netRxChan)
}

func networkUnpack(netPdu []byte) (out *NetworkMessage, err error) {
	defer func() {
		if out != nil && out.netKey.NewKey != nil {
			//this is an old key
		}
	}()
	// cache
	for e := cache.Front(); e != nil; e = e.Next() {
		if bytes.Equal(e.Value.([]byte), netPdu) {
			// duplicated message, already handled
			return nil, nil
		}
	}
	if cache.Len() > cacheSize {
		cache.Remove(cache.Front())
	}
	cache.PushBack(netPdu)

	msg := NetworkMessage{}
	err = utils.UnpackBE(netPdu[:1], "1, 7", &msg.ivi, &msg.nid)
	if err != nil {
		return nil, err
	}
	netKeys := findNetKeyByNid(msg.nid)
	if len(netKeys) == 0 {
		return nil, errors.NetKeyNotFoundByNid.New().AddContext(msg.nid)
	}
	for _, netKey := range netKeys {
		obfuscated := netPdu[1:7]
		// When this state is active, a node shall transmit using the current IV Index
		// and shall process messages from the current IV Index and also the current IV Index - 1.
		msg.ivIndex = meshDb.IVindex
		if (msg.ivIndex & 0x01) != msg.ivi {
			msg.ivIndex--
		}
		privacyPlain, err := utils.PackBE("040, 32, B8", msg.ivIndex, netPdu[7:14])
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		pecb, err := crypto.AES_ECB(netKey.PrivacyKey, privacyPlain)
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		deobfuscated := make([]byte, len(obfuscated), len(obfuscated))
		for i, d := range obfuscated {
			deobfuscated[i] = d ^ pecb[i]
		}
		err = utils.UnpackBE(deobfuscated, "1, 7, 24, 16", &msg.ctl, &msg.ttl, &msg.seq, &msg.src)
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		netMicLen := 4
		if msg.ctl == 1 {
			netMicLen = 8
		}
		// netMic, encrypted := netPdu[-netMicLen:], netPdu[7:-netMicLen]
		nonce, err := genNetworkNonce(msg.src, msg.seq, msg.ivIndex, msg.ctl, msg.ttl)
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		plainNet, err := crypto.AES_CCM_Decrypt(
			netKey.EncryptionKey,
			nonce,
			netPdu[7:],
			netMicLen)
		if plainNet == nil || err != nil {
			loggerNet.Error(err)
			continue
		}
		err = utils.UnpackBE(plainNet[:2], "16", &msg.dst)
		if err != nil {
			loggerNet.Error(err)
			continue
		}
		msg.plain = plainNet[2:]
		msg.netKey = netKey
		if node, _ := findNodeByAddr(msg.src); node != nil {
			node.SequenceNumber = uint(msg.seq)
		}
		loggerNet.Debugf("NET Rx: %+#v", msg)
		return &msg, nil
	}

	return nil, errors.NoValidNetKeyForDecryption.New().AddContext(netPdu)
}

func networkSend(msg *NetworkMessage) error {
	netPdu, err := networkPack(msg)
	if err != nil {
		return err
	}
	netBear.SendNetPdu(netPdu)
	return nil
	// privacy_random = bitstring.pack('pad:40, uintbe:32, bytes:7',
	// iv_index, network_pdu[:7]).bytes
}

func networkPack(msg *NetworkMessage) ([]byte, error) {
	micSize := 4
	if msg.ctl == 1 {
		micSize = 8
	}
	key := msg.netKey
	nonce, err := genNetworkNonce(msg.src, msg.seq, msg.ivIndex, msg.ctl, msg.ttl)
	if err != nil {
		return nil, err
	}
	plain, err := utils.PackBE("16,B8", msg.dst, msg.plain)
	if err != nil {
		return nil, err
	}
	cipher, netMic, err := crypto.AES_CCM(key.EncryptionKey, nonce, plain, micSize)
	if err != nil {
		return nil, err
	}
	preObfuscated, err := utils.PackBE("1,7,24,16", msg.ctl, msg.ttl, msg.seq, msg.src)
	if err != nil {
		return nil, err
	}
	privacyRandom, err := utils.PackBE("040,32,B8", meshDb.IVindex, append(cipher, netMic...)[:7])
	if err != nil {
		return nil, err
	}

	pecb, err := crypto.AES_ECB(key.PrivacyKey, privacyRandom)
	if err != nil {
		return nil, err
	}
	pecb = pecb[:6]
	obfuscated := make([]byte, len(pecb), len(pecb))
	for i, d := range preObfuscated {
		obfuscated[i] = d ^ pecb[i]
	}
	netPdu, err := utils.PackBE("1,7,B8,B8,B8", meshDb.IVindex&0x01, key.Nid, obfuscated, cipher, netMic)
	loggerNet.Debugf("net message to sent:%+#v", msg)
	return netPdu, err
}
