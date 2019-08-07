package mesh

import (
	"ble-mesh/utils"
	"sync"

	"github.com/sirupsen/logrus"
)

type (
	AdvertisingBear struct {
		writeCb   func([]byte) error
		logger    *logrus.Entry
		writeLock *sync.Mutex
	}
)

/* Advertising types */
const (
	BT_LE_ADV_BEACON    = 0x2B
	BT_LE_ADV_NETWORK   = 0x2A
	BT_LE_ADV_PROVISION = 0x29
)

func (b *AdvertisingBear) Start() {
	b.logger = utils.CreateLogger("advBear")
	b.writeLock = new(sync.Mutex)
}

func (b *AdvertisingBear) Stop() {
}

func (b *AdvertisingBear) OnPduReceived(data []byte) {
	if len(data) < 2 {
		return
	}
	lenMsg := data[0]
	advType := data[1]
	payload := data[2 : lenMsg+1]
	switch advType {
	case BT_LE_ADV_PROVISION:
		// PB-ADV
	case BT_LE_ADV_NETWORK:
		networkReceive(payload)
	case BT_LE_ADV_BEACON:
		meshBeaconReceive(payload)
	}
}

func (b *AdvertisingBear) SetWriteHandle(writeFunc func([]byte) error) {
	b.writeCb = writeFunc
}

func (b *AdvertisingBear) SetMTU(mtu uint) {

}

func (b *AdvertisingBear) SendNetPdu(pdu []byte) {
	if b.writeCb != nil {
		lenAdv := len(pdu) + 1
		packet := append([]byte{byte(lenAdv), BT_LE_ADV_NETWORK}, pdu...)
		b.logger.Debugf("adv tx: % 2x", packet)
		b.writeLock.Lock()

		b.writeCb(packet)
		b.writeLock.Unlock()
	}
}

func (b *AdvertisingBear) SendProvPdu(pdu []byte) {
}
