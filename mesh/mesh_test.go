package mesh

import (
	"ble-mesh/mesh/crypto"
	"encoding/hex"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func Test_netPduUnpack(t *testing.T) {
	homeDir, _ := homedir.Dir()
	Init(homeDir + "/.config/ble-mesh")

	netKey := NetKey{}
	netKey.Bytes, _ = hex.DecodeString("7dd7364cd842ad18c17c2b820c84c3d6")
	netKey.Nid, netKey.EncryptionKey, netKey.PrivacyKey, _ =
		crypto.K2(netKey.Bytes, []byte{0x00})
	meshDb.IVindex = 0x12345678
	meshDb.NetKeys[0] = &netKey
	netPdu, _ := hex.DecodeString("68eca487516765b5e5bfdacbaf6cb7fb6bff871f035444ce83a670df")
	msg := networkUnpack(netPdu)
	expected, _ := hex.DecodeString("034b50057e400000010000")
	assert.Equal(t, expected, msg.plain, "they should be equal")

	netPdu, _ = hex.DecodeString("68aec467ed4901d85d806bbed248614f938067b0d983bb7b")
	msg = networkUnpack(netPdu)
	expected, _ = hex.DecodeString("00a6ac00000003")
	assert.Equal(t, expected, msg.plain, "they should be equal")
}

func Test_netPduPack(t *testing.T) {
	homeDir, _ := homedir.Dir()
	Init(homeDir + "/.config/ble-mesh")
	netKey := NetKey{}
	netKey.Bytes, _ = hex.DecodeString("7dd7364cd842ad18c17c2b820c84c3d6")
	netKey.Nid, netKey.EncryptionKey, netKey.PrivacyKey, _ =
		crypto.K2(netKey.Bytes, []byte{0x00})
	meshDb.UnicastAddress = 0x0003
	meshDb.IVindex = 0x12345678
	meshDb.NetKeys[0] = &netKey
	lowerTpPdu, _ := hex.DecodeString("8026ac01ee9dddfd2169326d23f3afdf")
	netMsg := &NetworkMessage{
		ivi:     0,
		nid:     netKey.Nid,
		ctl:     0,
		ttl:     4,
		seq:     0x3129ab,
		src:     0x03,
		dst:     0x1201,
		plain:   lowerTpPdu,
		ivIndex: meshDb.IVindex,
	}
	netPdu := networkPack(netMsg)
	expected, _ := hex.DecodeString("68cab5c5348a230afba8c63d4e686364979deaf4fd40961145939cda0e")
	assert.Equal(t, expected, netPdu, "they should be equal")

}

func Test_tpPduUnpack(t *testing.T) {
	homeDir, _ := homedir.Dir()
	Init(homeDir + "/.config/ble-mesh")
	key, _ := hex.DecodeString("63964771734fbd76e3b40519d1d94a48")
	aid, _ := crypto.K4(key)
	meshDb.AppKeys[0] = &AppKey{
		Aid: aid, Bytes: key, Index: 0,
	}
	lowerTpPdu, _ := hex.DecodeString("5a8bde6d9106ea078a")
	upperTpPdu, _ := tpDecryptMessage(lowerTpPdu, 0x1201, 0xffff, 0x07, 0x12345678, 0, 1, 0x26)
	expected, _ := hex.DecodeString("0400000000")
	assert.Equal(t, expected, upperTpPdu, "they should be equal")

}

func Test_CompositionParser(t *testing.T) {
	expected := &Composition{
		CID:      0x000c,
		PID:      0x001A,
		VID:      0x0001,
		CRPL:     0x0008,
		Features: Features{Relay: true, Proxy: true, Friend: false, LowPower: false},
		Elements: []CompositionElement{
			CompositionElement{
				Location:     0x0100,
				SigModels:    []uint16{0x0000, 0x8000, 0x0001, 0x1000, 0x1003},
				VendorModels: []uint32{0x2A003F},
			},
		},
	}
	data, _ := hex.DecodeString("000C001A0001000800030000010501000000800100001003103F002A00")
	comp := parseCompData(data)
	assert.Equal(t, *expected, *comp, "they should be equal")

}
