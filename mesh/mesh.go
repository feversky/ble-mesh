package mesh

import (
	"ble-mesh/mesh/crypto"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

const (
	UNASSIGNED_ADDRESS   = 0x0000
	PROXIES_ADDRESS      = 0xfffc
	FRIENDS_ADDRESS      = 0xfffd
	RELAYS_ADDRESS       = 0xfffe
	ALL_NODES_ADDRESS    = 0xffff
	VIRTUAL_ADDRESS_LOW  = 0x8000
	VIRTUAL_ADDRESS_HIGH = 0xbfff
	GROUP_ADDRESS_LOW    = 0xc000
	GROUP_ADDRESS_HIGH   = 0xff00
)

const (
	GATT_BEAR = iota
	ADV_BEAR
)

var (
	loggerMesh    = utils.CreateLogger("Mesh")
	netBear       Bear
	provBear      Bear
	currentBear   uint
	provStoppedCb func()
)

func genNetworkNonce(src, seq, ivIndex, ctl, ttl uint) ([]byte, error) {
	data, err := utils.PackBE("8, 1, 7, 24, 16, 016, 32", 0x00, ctl, ttl, seq, src, ivIndex)
	if err != nil {
		err = errors.FailedToGenerateNonce.New().AddContext(err)
	}
	return data, err
}

func genApplicationNonce(src, seq, ivIndex, szmic, dst uint) ([]byte, error) {
	data, err := utils.PackBE("8, 1, 07, 24, 16, 16, 32", 0x01, szmic, seq, src, dst, ivIndex)
	if err != nil {
		err = errors.FailedToGenerateNonce.New().AddContext(err)
	}
	return data, err
}

func genDeviceNonce(src, seq, ivIndex, szmic, dst uint) ([]byte, error) {
	data, err := utils.PackBE("8, 1, 07, 24, 16, 16, 32", 0x02, szmic, seq, src, dst, ivIndex)
	if err != nil {
		err = errors.FailedToGenerateNonce.New().AddContext(err)
	}
	return data, err
}

func StartMeshNetwork() {
	netBear.Start()
	startNet()
	startTransport()
}

func StopMeshNetwork() {
	netBear.Stop()
	stopNet()
	stopTransport()
}

func StartMeshProvision(uuid string, provStopped func()) {
	provBear.Start()
	startProvision(uuid)
	provStoppedCb = provStopped
}

func onProvisionFinished() {
	if provStoppedCb != nil {
		provStoppedCb()
	}
	StopMeshProvision()
}

func StopMeshProvision() {
	provBear.Stop()
	stopProvision()
}

func SetNetworkBear(b Bear) {
	netBear = b
}

func SetProvisionBear(b Bear) {
	provBear = b
}

func OnClose() {
	writeMeshToDb()
	for _, n := range meshDb.Nodes {
		writeNodeToDb(n)
	}
}

func RefreshNetKey(index uint) error {
	//generate a new key and save
	keyBytes := make([]byte, 16)
	rand.Read(keyBytes)
	newKey := createNetKeyB(keyBytes, index)
	oldKey, _ := findNetKeyByIndex(index)
	oldKey.NewKey = newKey
	newKey.OldKey = oldKey
	meshDb.NetKeys[index] = newKey
	failedNodes := []*Node{}
	timeoutNodes := []*Node{}
	friends := map[*Node]*Node{}
	writeMeshToDb()

	finalResult := make(chan bool)
	maxTimeout := 1000
	resps := map[uint]*AccessMessage{}
	var responseHandler modelMsglistener
	responseHandler = func(msg *AccessMessage) {
		if msg.opcode == opConfigNetKeyStatus {
			resps[msg.src] = msg
		}
	}
	go func(to *int) {
		defer unregisterModelMessageRxListener(&responseHandler)
		registerModelMessageRxListener(&responseHandler)
		startTime := time.Now()
		for {
			time.Sleep(time.Millisecond * 100)
			if time.Now().After(startTime.Add(time.Duration(*to) * time.Second)) {
				finalResult <- false
				break
			}
			if len(resps) == len(meshDb.Nodes) {
				for src, resp := range resps {
					netKey := binary.LittleEndian.Uint16(resp.payload[1:3]) & 0x0FFF
					if netKey == uint16(index) {
						status := statusCode[resp.payload[0]]
						// loggerConfCli.Infof("got response for net key %d, status: %s", netKey, status)
						if status == "Success" {
							//copy netkey from db
							node, _ := findNodeByAddr(src)
							oldkey, _ := node.findNodeNetKeyByIndex(index)
							copier.Copy(oldkey, newKey)
							continue
						}
					}
					finalResult <- false
					break
				}
				finalResult <- true
				break
			}
		}
	}(&maxTimeout)

	// opConfigNetKeyUpdate
	var wg sync.WaitGroup
	for _, node := range meshDb.Nodes {
		wg.Add(1)
		payload, _ := utils.PackLE("12,04,B", newKey.Index, newKey.Bytes)
		modelSendTmpl1(false, node.UnicastAddress, opConfigNetKeyUpdate, payload, 10, func(msg *AccessMessage) error {
			netKey := binary.LittleEndian.Uint16(msg.payload[1:3]) & 0x0FFF
			status := statusCode[msg.payload[0]]
			loggerMesh.Infof("got response for net key %d, status: %s", netKey, status)
			if status == "Success" {
				//copy netkey from db
				if friends[node] != nil {
					delete(friends, node)
				}
			} else {
				loggerMesh.Errorf("cannot update netkey for node %x, response: %s", node.UnicastAddress, status)
				failedNodes = append(failedNodes, node)
			}
			wg.Done()
			//todo: implement error
			return nil
		}, func() {
			//timeout
			if friends[node] == nil {
				timeoutNodes = append(timeoutNodes, node)
			}
			wg.Done()
		}, func(ack *SegmentAckMessage) {
			// obo ack received
			if ack.obo == 1 {
				friend, _ := findNodeByAddr(ack.src)
				// update mesh topology
				node.Friend = friend
				node.LPN = true
				writeNodeToDb(node)

				friends[node] = friend
			}
		})
	}
	wg.Wait()
	errStr := ""
	for _, n := range failedNodes {
		errStr += fmt.Sprintf("fail to update netkey for node %x\n", n.UnicastAddress)
	}
	for _, n := range timeoutNodes {
		errStr += fmt.Sprintf("timeout to update netkey for node %x\n", n.UnicastAddress)
	}
	if errStr != "" {
		// return errors.New(errStr)
	}

	// On receiving Segment Acknowledgments with the OBO field set to 1 to key update messages
	// sent to a Low Power node, a Configuration Client may perform a PollTimeout List procedure
	// to the Low Power node's Friend node (identifying the Friend node using the value of SRC field
	// of the Segment Acknowledgment) in order to obtain the current value of the PollTimeout timer,
	// and schedule retries of NetKey or AppKey updates based on this value.
	timeoutFriendNodes := []*Node{}
	// lpnPollTimeout := map[*Node]uint{}
	var maxPollTimeout uint
	for lpn, friend := range friends {
		wg.Add(1)
		payload, _ := utils.PackLE("16", lpn.UnicastAddress)
		modelSendTmpl1(false, friend.UnicastAddress, opConfigLowPowerNodePollTimeoutGet, payload, 5,
			func(msg *AccessMessage) error {
				var addr, timeout uint
				utils.UnpackLE(msg.payload, "16,32", &addr, &timeout)
				// lpnPollTimeout[lpn] = timeout / 10
				if maxPollTimeout < timeout {
					maxPollTimeout = timeout
				}
				loggerFoundation.Infof("current timeout of LPN %4x: %dms", addr, timeout*100)
				wg.Done()
				//todo: implement error
				return nil
			},
			func() {
				//timeout should not happen
				timeoutFriendNodes = append(timeoutFriendNodes, friend)
				wg.Done()
			}, nil)
	}
	wg.Wait()
	for _, n := range timeoutFriendNodes {
		errStr += fmt.Sprintf("timeout occurs during PollTimeout to node %x\n", n.UnicastAddress)
	}
	if errStr != "" {
		// return errors.New(errStr)
	}
	maxTimeout = int(maxPollTimeout)

	res := <-finalResult
	if res {
		loggerMesh.Info("netkey refresh phase 1 end")
		newKey.KeyRefreshPhase = 1
		return nil
	}
	// return errors.New(string.Sprintf("timeout occurs during PollTimeout to address %4x", friend.UnicastAddress))

	for _, n := range meshDb.Nodes {
		var found bool
		for _, rn := range resps {
			if rn.src == n.UnicastAddress {
				found = true
				break
			}
		}
		if !found {
			errStr += fmt.Sprintf("no response from node %x\n", n.UnicastAddress)
		}
	}
	// return errors.New(errStr)
	return nil
}

// handles mesh beacon
func meshBeaconReceive(proxyPdu []byte) {
	// netChan <- proxyPdu
	beaconType, beaconData := proxyPdu[0], proxyPdu[1:]
	switch beaconType {
	case 0x00:
		//  unprovisioned beacon, ignore
	case 0x01:
		// all proxy and relay node shall send beacon, other nodes might send beacon
		BEACON_AUTH_SIZE := 8
		beacon := beaconData[:len(beaconData)-BEACON_AUTH_SIZE]
		auth := beaconData[len(beaconData)-BEACON_AUTH_SIZE:]
		// BEACON_FORMAT = 'pad:6, uint:1, uint:1, bytes:8, uintbe:32'
		var ivUpdateFlag, keyRefresh, ivIndex uint
		var networkId uint64
		utils.UnpackBE(beacon, "06, 1, 1, 64, 32", &ivUpdateFlag, &keyRefresh, &networkId, &ivIndex)
		// todo: key refresh, use both new key and old key
		for _, netKey := range meshDb.NetKeys {
			if netKey.NetworkId == networkId {
				authVerify, _ := crypto.AES_CMAC(netKey.BeaconKey, beacon)
				if bytes.Equal(auth, authVerify[:8]) {
					loggerMesh.Infof("received mesh beacon: keyReresh %d, IV index %d, IV update flag %d", keyRefresh, ivIndex, ivUpdateFlag)
					// todo: KeyRefresh procedure
					// netKey.KeyRefresh = keyRefresh
					if ivIndex == meshDb.IVindex {
						meshDb.IVupdate = ivUpdateFlag
					}
					if ivIndex == meshDb.IVindex+1 && ivUpdateFlag == 1 {
						// key refresh procedure
						meshDb.IVindex = ivIndex
						meshDb.IVupdate = ivUpdateFlag
					}
				}
			}
		}
	}
}

func isUnicastAddr(addr uint) bool {
	return addr > UNASSIGNED_ADDRESS && addr <= VIRTUAL_ADDRESS_LOW
}

func isUnassigned(addr uint) bool {
	return addr == UNASSIGNED_ADDRESS
}
func isGroupAddr(addr uint) bool {
	return addr > GROUP_ADDRESS_LOW && addr <= GROUP_ADDRESS_HIGH
}
func isVirtualAddr(addr uint) bool {
	return addr > VIRTUAL_ADDRESS_LOW && addr <= VIRTUAL_ADDRESS_HIGH
}
func isAllAddr(addr uint) bool {
	return addr == ALL_NODES_ADDRESS
}

// Sig Model: 16, Vendor Model: 32
func getModelIdBitLength(id uint) uint {
	if isSigModel(id) {
		return 16
	}
	return 32
}

func isSigModel(id uint) bool {
	return id <= 0xFFFF
}

func calcUnicastAddr(uuid string) uint {
	var maxAddr uint
	for k, n := range meshDb.Nodes {
		if k > maxAddr {
			maxAddr = k
		}
		if n.UUID == uuid {
			return n.UnicastAddress
		}
	}
	if meshDb.LowAddress > maxAddr {
		return meshDb.LowAddress
	}
	return maxAddr + uint(len(meshDb.Nodes[maxAddr].Elements))

}

func SetNode(nodeNew *Node) error {
	node, _ := findNodeByAddr(nodeNew.UnicastAddress)

	if !reflect.DeepEqual(node.BindedKeys, nodeNew.BindedKeys) {
		// refresh netkey list
		err := ConfigNetKeyGet(node.UnicastAddress)
		if err != nil {
			return err
		}
		// delete netkeys
		for _, bd := range node.BindedKeys {
			if keyBd, _ := nodeNew.findNodeKeyBindingByNetKeyIndex(bd.NetKeyIndex); keyBd == nil {
				// cannot find the key in new node setting, delete it
				err := ConfigNetKeyDelete(node.UnicastAddress, bd.NetKeyIndex)
				if err != nil {
					return err
				}
			}
		}
		// add netkey
		for _, bd := range nodeNew.BindedKeys {
			if keyBd, _ := node.findNodeKeyBindingByNetKeyIndex(bd.NetKeyIndex); keyBd == nil {
				// cannot find the key in old node setting, add it
				err := ConfigNetKeyAdd(node.UnicastAddress, bd.NetKeyIndex)
				if err != nil {
					return err
				}
			}
		}

		// refresh appkey list
		for _, bd := range nodeNew.BindedKeys {
			err := ConfigAppKeyGet(node.UnicastAddress, bd.NetKeyIndex)
			if err != nil {
				return err
			}
			// delete appkey
			appkeysNew := bd.BindedAppKeyIds
			bindingOld, _ := node.findNodeKeyBindingByNetKeyIndex(bd.NetKeyIndex)
			appkeysOld := bindingOld.BindedAppKeyIds
			for _, keyOld := range appkeysOld {
				if !utils.Contains(appkeysNew, keyOld) {
					err := ConfigAppKeyDelete(node.UnicastAddress, bd.NetKeyIndex, keyOld)
					if err != nil {
						return err
					}
				}
			}
			for _, keyNew := range appkeysNew {
				if !utils.Contains(appkeysOld, keyNew) {
					err := ConfigAppKeyAdd(node.UnicastAddress, bd.NetKeyIndex, keyNew)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// model appkey binding
	for _, elementNew := range nodeNew.Elements {
		for _, modelNew := range elementNew.Models {
			elementOld, _ := findElementByAddr(elementNew.UnicastAddress)
			modelOld, _ := elementOld.findModel(modelNew.ModelID)
			if modelOld.BindedAppKeyIds == nil {
				modelOld.BindedAppKeyIds = []uint{}
			}
			if modelNew.BindedAppKeyIds == nil {
				modelNew.BindedAppKeyIds = []uint{}
			}
			if !reflect.DeepEqual(modelOld.BindedAppKeyIds, modelNew.BindedAppKeyIds) {
				err := configModelAppGet(elementOld.UnicastAddress, modelOld.ModelID)
				if err != nil {
					return err
				}
				appkeysNew := modelNew.BindedAppKeyIds
				appkeysOld := modelOld.BindedAppKeyIds
				// delete appkey
				for _, keyOld := range appkeysOld {
					if !utils.Contains(appkeysNew, keyOld) {
						err := ConfigModelAppUnbind(elementOld.UnicastAddress, keyOld, modelOld.ModelID)
						if err != nil {
							return err
						}
					}
				}
				// bind new appkey
				for _, keyNew := range appkeysNew {
					if !utils.Contains(appkeysOld, keyNew) {
						err := ConfigModelAppBind(node.UnicastAddress, keyNew, modelOld.ModelID)
						if err != nil {
							return err
						}
					}
				}
			}

			// publication
			// todo: support virtual address
			if !reflect.DeepEqual(modelOld.PubSetting, modelNew.PubSetting) {
				pubSetting := modelNew.PubSetting
				err := ConfigModelPublicationSet(
					elementOld.UnicastAddress,
					pubSetting.PublishAddress,
					pubSetting.AppKeyIndex,
					pubSetting.CredentialFlag,
					pubSetting.PublishTTL,
					pubSetting.PublishPeriod.NumberOfSteps,
					pubSetting.PublishPeriod.StepResolution,
					pubSetting.PublishRetransmitCount,
					pubSetting.PublishRetransmitIntervalSteps,
					modelOld.ModelID,
				)
				if err != nil {
					return err
				}
			}

			// subscription
			// todo: support virtual address
			if modelOld.SubAddresses == nil {
				modelOld.SubAddresses = []uint{}
			}
			if modelNew.SubAddresses == nil {
				modelNew.SubAddresses = []uint{}
			}
			if !reflect.DeepEqual(modelOld.SubAddresses, modelNew.SubAddresses) {
				// logger.Debug(modelOld.SubAddresses, modelNew.SubAddresses)
				// if modelOld.SubAddresses == nil || len(modelOld.SubAddresses) == 0 {
				err := configModelAppGet(elementOld.UnicastAddress, modelOld.ModelID)
				if err != nil {
					return err
				}
				// }
				if len(modelOld.SubAddresses) > 0 {
					err := ConfigModelSubscriptionDeleteAll(elementOld.UnicastAddress, modelOld.ModelID)
					if err != nil {
						return err
					}
				}
				for _, subAddr := range modelNew.SubAddresses {
					err := ConfigModelSubscriptionAdd(elementOld.UnicastAddress, subAddr, modelOld.ModelID)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// beacon
	if node.SecureNetworkBeacon != nodeNew.SecureNetworkBeacon {
		err := ConfigBeaconSet(node.UnicastAddress, nodeNew.SecureNetworkBeacon)
		if err != nil {
			return err
		}
	}

	// ttl
	if node.DefaultTTL != nodeNew.DefaultTTL {
		err := ConfigDefaultTTLSet(node.UnicastAddress, nodeNew.DefaultTTL)
		if err != nil {
			return err
		}
	}

	// gatt proxy
	if node.GATTProxyState != nodeNew.GATTProxyState {
		err := ConfigGattProxySet(node.UnicastAddress, nodeNew.GATTProxyState)
		if err != nil {
			return err
		}
	}

	// friend
	if node.Features.Friend {
		if node.FriendState != nodeNew.FriendState {
			err := ConfigFriendSet(node.UnicastAddress, nodeNew.FriendState)
			if err != nil {
				return err
			}
		}
		if node.FriendState == 1 {
			for _, lpn := range node.findLpnFriends() {
				err := ConfigLowPowerNodePollTimeoutGet(node.UnicastAddress, lpn.UnicastAddress)
				if err != nil {
					return err
				}
			}
		}

	}

	// node identity
	for netkeyId, identity := range nodeNew.NodeIdentityStates {
		if node.NodeIdentityStates[netkeyId] != identity {
			err := ConfigNodeIdentitySet(node.UnicastAddress, netkeyId, identity)
			if err != nil {
				return err
			}
		}
	}

	// network transmit
	if node.NetwrokTransmitState.NetworkTransmitCount != nodeNew.NetwrokTransmitState.NetworkTransmitCount ||
		node.NetwrokTransmitState.NetworkTransmitIntervalSteps != nodeNew.NetwrokTransmitState.NetworkTransmitIntervalSteps {
		err := ConfigNetworkTransmitSet(
			node.UnicastAddress,
			nodeNew.NetwrokTransmitState.NetworkTransmitCount,
			nodeNew.NetwrokTransmitState.NetworkTransmitIntervalSteps,
		)
		if err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(node.RelayState, nodeNew.RelayState) {
		err := ConfigRelaySet(
			node.UnicastAddress,
			nodeNew.RelayState.Relay,
			nodeNew.RelayState.RelayRetransmitCount,
			nodeNew.RelayState.RelayRetransmitIntervalSteps,
		)
		if err != nil {
			return err
		}
	}

	if node.AttentionTimer != nodeNew.AttentionTimer {
	}

	writeNodeToDb(node)
	// todo: heartbeat publication & subscription
	return nil
}

func ResetNode(addr uint) error {
	_, err := findNodeByAddr(addr)
	if err != nil {
		return err
	}
	return ConfigNodeReset(addr)
}
