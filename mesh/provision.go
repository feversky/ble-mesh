package mesh

import (
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/aead/ecdh"

	meshCrypto "ble-mesh/mesh/crypto"
	"ble-mesh/utils"
)

type Capability struct {
	NumElements     byte
	Algorithms      uint16
	PublicKeyType   byte
	StaticOobType   byte
	OutputOobType   byte
	OutputOobAction uint16
	InputOobType    byte
	InputOobAction  uint16
}

type ProvNode struct {
	cap              *Capability
	p256             ecdh.KeyExchange
	pubKey           crypto.PublicKey
	priKey           crypto.PrivateKey
	ecdhSecret       []byte
	authValue        []byte
	random           []byte
	confirmation     []byte
	confirmationKey  []byte
	confirmationSalt []byte
	node             *Node
}

type ProvisionData struct {
	NetKey      []byte
	NetKeyIndex uint
	Flag        byte
	IvIndex     uint
	UnicastAddr uint
}

type Net struct {
	localNode  *ProvNode
	remoteNode *ProvNode
}

const useNonOob bool = true

const (
	provInvite = iota
	provCapabilities
	provStart
	provPublicKey
	provInputComplete
	provConfirmation
	provRandom
	provData
	provComplete
	provFailed
)

var (
	provChan           chan []byte
	net                Net
	confirmationInputs []byte
	loggerProv         = utils.CreateLogger("Provision")
)

var expectedLen = map[int]int{
	provCapabilities: 11,
	provPublicKey:    64,
	provConfirmation: 16,
	provRandom:       16,
	provComplete:     0,
	provFailed:       1,
}

func provisionReceive(proxyPdu []byte) {
	provChan <- proxyPdu
}

func provProc() {
	defer onProvisionFinished()
	state := provInvite
	confirmationInputs := make([]byte, 0)
	pduOut := genProvInvite()
	confirmationInputs = append(confirmationInputs, pduOut[1:]...)
	provBear.SendProvPdu(pduOut)
	state = provCapabilities
	for {
		provPdu, more := <-provChan
		if !more {
			return
		}
		loggerProv.Debugf("provision pdu RX: % 2x", provPdu)
		pduType := int(provPdu[0])
		data := provPdu[1:]
		if len(data) != expectedLen[pduType] {
			loggerProv.Errorf("unexpected pdu length, pdu type:%d", pduType)
			goto failed
		}
		if pduType == provFailed {
			loggerProv.Errorf("provision failed, error code:%d", data[0])
			goto failed
		}

		switch state {
		case provInvite:
			state = provCapabilities
		case provCapabilities:
			if state != pduType {
				loggerProv.Errorf("expect a provCapabilities pdu, but recevied a %d", pduType)
				goto failed
			}
			cap := &Capability{}
			err := utils.ReadStructFromBuffer(data, cap)
			if err != nil {
				loggerProv.Errorf("decoding error while processing capability, error:%s", err)
				goto failed
			}

			confirmationInputs = append(confirmationInputs, provPdu[1:]...)

			pduOut = genProvStart()
			confirmationInputs = append(confirmationInputs, pduOut[1:]...)
			provBear.SendProvPdu(pduOut)
			time.Sleep(time.Second)

			state = provPublicKey

			p256 := ecdh.Generic(elliptic.P256())

			priKey, pubKey, err := p256.GenerateKey(rand.Reader)
			if err != nil {
				loggerProv.Errorf("Failed to generate private/public key pair: %s\n", err)
			}
			provPubKey, _ := pubKey.(ecdh.Point)
			confirmationInputs = append(
				append(confirmationInputs, provPubKey.X.Bytes()...),
				provPubKey.Y.Bytes()...)

			net.remoteNode.cap = cap
			net.localNode = &ProvNode{}
			net.localNode.pubKey = pubKey
			net.localNode.priKey = priKey
			net.localNode.p256 = p256
			provBear.SendProvPdu(genProvPublicKey())
		case provPublicKey:
			x := new(big.Int).SetBytes(data[:32])
			y := new(big.Int).SetBytes(data[32:])
			net.remoteNode.pubKey = &ecdh.Point{X: x, Y: y}
			ecdhSecret := net.localNode.p256.ComputeSecret(net.localNode.priKey, net.remoteNode.pubKey)
			confirmationInputs = append(confirmationInputs, data...)
			confirmationSalt, _ := meshCrypto.S1(confirmationInputs)
			confirmationKey, _ := meshCrypto.K1(ecdhSecret, confirmationSalt, []byte("prck"))
			authValue := make([]byte, 16, 16)
			randoms := make([]byte, 16, 16)
			rand.Read(randoms)
			net.localNode.ecdhSecret = ecdhSecret
			net.localNode.confirmationSalt = confirmationSalt
			net.localNode.confirmationKey = confirmationKey
			net.localNode.authValue = authValue
			net.localNode.random = randoms
			net.localNode.confirmation, _ = meshCrypto.AES_CMAC(confirmationKey, append(randoms, authValue...))

			loggerProv.Debugf("ecdhSecret: %x", ecdhSecret)
			loggerProv.Debugf("confirmationInputs: %x", confirmationInputs)
			loggerProv.Debugf("confirmationSalt: %x", confirmationSalt)
			loggerProv.Debugf("confirmationKey: %x", confirmationKey)
			loggerProv.Debugf("randoms: %x", randoms)
			loggerProv.Debugf("confirmation: %x", net.localNode.confirmation)
			loggerProv.Debugf("node public key, x:%x, y:%x", x, y)

			provBear.SendProvPdu(genProvConfirmation())
			state = provConfirmation
		case provConfirmation:
			loggerProv.Debugf("node confirmation: %x", data)
			net.remoteNode.confirmation = data
			provBear.SendProvPdu(genProvRandom())
			state = provRandom
		case provRandom:
			loggerProv.Debugf("node random: %x", data)
			net.remoteNode.random = data
			calcConfirmation, _ := meshCrypto.AES_CMAC(net.localNode.confirmationKey, append(data, net.localNode.authValue...))
			if bytes.Equal(calcConfirmation, net.remoteNode.confirmation) {
				loggerProv.Infof("provision confirmation successful")
			} else {
				loggerProv.Infof("provision confirmation failed")
				goto failed
			}
			provisioningSalt, _ := meshCrypto.S1(
				append(
					append(net.localNode.confirmationSalt, net.localNode.random...),
					net.remoteNode.random...))
			sessionKey, _ := meshCrypto.K1(net.localNode.ecdhSecret, provisioningSalt, []byte("prsk"))
			sessionNonce, _ := meshCrypto.K1(net.localNode.ecdhSecret, provisioningSalt, []byte("prsn"))
			devKey, _ := meshCrypto.K1(net.localNode.ecdhSecret, provisioningSalt, []byte("prdk"))
			//Provisioning Data = Network Key || Key Index || Flags || IV Index || Unicast Address

			netKey, _ := findNetKeyByIndex(0)
			if netKey == nil {
				loggerProv.Error("netkey 000 does not exist")
				return
			}
			unicastAddr := calcUnicastAddr(net.remoteNode.node.UUID)
			b := make([]byte, 2)
			binary.LittleEndian.PutUint16(b, uint16(netKey.Index))
			//Key Refresh Flag    0: Key Refresh Phase 0      1: Key Refresh Phase 2
			keyRefresh := uint(0)
			if netKey.KeyRefreshPhase == 2 {
				keyRefresh = 1
			}
			provData, _ := utils.PackBE("B8,B8,8,32,16", netKey.Bytes, b,
				byte((meshDb.IVupdate<<1)+keyRefresh), meshDb.IVindex, unicastAddr)
			loggerProv.Debugf("provision data: % 2x", provData)

			enc, tag, _ := meshCrypto.AES_CCM(sessionKey, sessionNonce[len(sessionNonce)-13:], provData, 8)
			loggerProv.Infof("sending provision data...")
			provBear.SendProvPdu(genProvData(enc, tag))

			net.remoteNode.node.UnicastAddress = unicastAddr
			net.remoteNode.node.DeviceKey = DevKey{}
			net.remoteNode.node.DeviceKey.Bytes = devKey
			net.remoteNode.node.DeviceKey.Aid = 0
			net.remoteNode.node.NodeIdentityStates = map[uint]uint{}
			net.remoteNode.node.BindedKeys = []NodeKeyBinding{
				NodeKeyBinding{NetKeyIndex: netKey.Index, BindedAppKeyIds: []uint{}},
			}
			state = provComplete
		case provComplete:
			loggerProv.Infof("provision successful")
			meshDb.Nodes[net.remoteNode.node.UnicastAddress] = net.remoteNode.node
			writeNodeToDb(net.remoteNode.node)
			return
		}

	}

failed:
	loggerProv.Errorf("error happened during provision")
}

func startProvision(uuid string) {
	net.remoteNode = &ProvNode{}
	net.remoteNode.node = &Node{}
	net.remoteNode.node.UUID = uuid
	provChan = make(chan []byte)
	go provProc()
}

func stopProvision() {
	close(provChan)
}

func genProvInvite() []byte {
	return []byte{provInvite, 10}
}

func genProvStart() []byte {
	var algorithm byte
	var publicKey byte
	var authMethod byte
	var authAction byte
	var authSize byte
	return []byte{provStart, algorithm, publicKey, authMethod, authAction, authSize}
}

func genProvPublicKey() []byte {
	pubKey := net.localNode.pubKey.(ecdh.Point)
	return append(
		append([]byte{provPublicKey}, pubKey.X.Bytes()...),
		pubKey.Y.Bytes()...)
}

func genProvConfirmation() []byte {
	return append([]byte{provConfirmation}, net.localNode.confirmation...)
}

func genProvRandom() []byte {
	return append([]byte{provRandom}, net.localNode.random...)
}

func genProvData(data, mic []byte) []byte {
	return append(append([]byte{provData}, data...), mic...)
}
