package mesh

import (
	"ble-mesh/utils"
	"sync"

	"github.com/sirupsen/logrus"
)

type GattProxyBear struct {
	proxyChan chan []byte
	proxyPdu  []byte
	writeCb   func([]byte) error
	mtu       int
	logger    *logrus.Entry
	writeLock *sync.Mutex
}

const (
	COMPLETE = iota
	FIRST
	CONTINUATION
	LAST
)

const (
	NETWORK = iota
	BEACON
	PROXY_CONFIG
	PROVISION
)

func (b *GattProxyBear) Start() {
	b.logger = utils.CreateLogger("GattBear")
	b.writeLock = new(sync.Mutex)
	if b.proxyChan == nil {
		b.proxyChan = make(chan []byte)
	}
	go b.assemblyProc()
}

func (b *GattProxyBear) Stop() {
	close(b.proxyChan)
}

func (b *GattProxyBear) OnPduReceived(data []byte) {
	b.proxyChan <- data
}

func (b *GattProxyBear) SetWriteHandle(writeFunc func([]byte) error) {
	b.writeCb = writeFunc
}

func (b *GattProxyBear) SetMTU(mtu uint) {
	b.mtu = int(mtu)
}

func (b *GattProxyBear) SendNetPdu(pdu []byte) {
	b.proxySend(byte(NETWORK), pdu)
}

func (b *GattProxyBear) SendProvPdu(pdu []byte) {
	b.proxySend(byte(PROVISION), pdu)
}

func (b *GattProxyBear) segmentation(proxyType byte, pdu []byte) [][]byte {
	list := [][]byte{}
	genSeg := func(header byte, body []byte) []byte { return append([]byte{byte(header<<6 | proxyType)}, body...) }

	pduMtu := b.mtu - 1
	if len(pdu) <= pduMtu {
		return append(list, genSeg(COMPLETE, pdu))
	}
	var chunk []byte
	chunk, pdu = pdu[:pduMtu], pdu[pduMtu:]
	list = append(list, genSeg(FIRST, chunk))
	for {
		if len(pdu) <= pduMtu {
			return append(list, genSeg(LAST, pdu))
		}
		chunk, pdu = pdu[:pduMtu], pdu[pduMtu:]
		list = append(list, genSeg(CONTINUATION, chunk))
	}
}

func (b *GattProxyBear) proxySend(proxyType byte, data []byte) {
	if b.writeCb == nil {
		return
	}
	b.logger.Debugf("proxy tx: % 2x", data)
	for _, seg := range b.segmentation(proxyType, data) {
		b.writeLock.Lock()
		b.writeCb(seg)
		b.writeLock.Unlock()
	}
}

func (b *GattProxyBear) assemblyProc() {
	for {
		packet, more := <-b.proxyChan
		if !more {
			return
		}
		b.logger.Debugf("PROXY RX: % 2x", packet)
		sar := packet[0] >> 6
		msgType := packet[0] & 0x3F

		switch sar {
		case COMPLETE:
			b.sendMessageToUpperLayer(msgType, packet[1:])
		case FIRST:
			b.proxyPdu = []byte{}
			b.proxyPdu = append(b.proxyPdu, packet[1:]...)
		case CONTINUATION:
			b.proxyPdu = append(b.proxyPdu, packet[1:]...)
		case LAST:
			b.proxyPdu = append(b.proxyPdu, packet[1:]...)
			b.sendMessageToUpperLayer(msgType, b.proxyPdu)
		}
	}
}

func (b *GattProxyBear) sendMessageToUpperLayer(msgType byte, data []byte) {
	switch msgType {
	case PROVISION:
		provisionReceive(data)
	case NETWORK:
		networkReceive(data)
	case BEACON:
		meshBeaconReceive(data)
	}
}
