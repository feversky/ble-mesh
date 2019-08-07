package mesh

import (
	"ble-mesh/mesh/crypto"
	"ble-mesh/mesh/db"
	"ble-mesh/mesh/def"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
)

type (
	Mesh struct {
		MeshName       string
		NetKeys        map[uint]*NetKey
		AppKeys        map[uint]*AppKey
		Nodes          map[uint]*Node
		Groups         map[uint]*Group
		UnicastAddress uint
		LowAddress     uint
		HighAddress    uint
		IVindex        uint
		IVupdate       uint
		SequenceNumber uint
	}

	NetKey struct {
		Index           uint
		KeyRefreshPhase uint
		Nid             uint
		Bytes           []byte
		EncryptionKey   []byte
		PrivacyKey      []byte
		BeaconKey       []byte
		IdentityKey     []byte
		NetworkId       uint64
		OldKey          *NetKey
		NewKey          *NetKey
	}

	AppKey struct {
		Index           uint
		KeyRefreshPhase uint
		Bytes           []byte
		Aid             uint
		OldKey          *AppKey
		NewKey          *AppKey
	}

	DevKey struct {
		AppKey
	}

	NodeKeyBinding struct {
		NetKeyIndex     uint
		BindedAppKeyIds []uint
	}

	Node struct {
		UUID           string
		UnicastAddress uint
		SequenceNumber uint
		Mac            string
		DeviceKey      DevKey
		BindedKeys     []NodeKeyBinding
		Friend         *Node
		LPN            bool
		// composition
		Elements []*Element
		Cid      int
		Pid      int
		Vid      int
		Crpl     int
		Features db.Features
		// status
		DefaultTTL           uint
		RelayState           def.ConfigRelayStatusMessageParameters
		AttentionTimer       uint
		SecureNetworkBeacon  bool
		GATTProxyState       uint
		NodeIdentityStates   map[uint]uint // key: net key id
		FriendState          uint
		KeyRefreshPhaseState uint
		PollTimeoutListState map[uint]uint // key: lpn address
		NetwrokTransmitState def.ConfigNetworkTransmitStatusMessageParameters
		CurrentFault         uint
	}

	Model struct {
		Element         *Element
		ModelID         uint
		PubSetting      def.ConfigModelPublicationStatusMessageParameters
		SubAddresses    []uint
		BindedAppKeyIds []uint
		State           interface{}
	}

	Element struct {
		Node           *Node
		ElementIndex   int
		Location       int
		UnicastAddress uint
		Models         []*Model
	}

	Group struct {
		Name    string
		Address uint
	}
)

var (
	confDir   string
	meshDb    *Mesh
	meshDbRaw *db.Mesh
)

func createNetKey(index uint) *NetKey {
	keyBytes := make([]byte, 16)
	rand.Read(keyBytes)
	return createNetKeyB(keyBytes, index)
}

func createNetKeyS(keyStr string, index uint) *NetKey {
	key, _ := hex.DecodeString(keyStr)
	return createNetKeyB(key, index)
}

func createNetKeyB(key []byte, index uint) *NetKey {
	nid, encryptionKey, privacyKey, _ := crypto.K2(key, []byte{0x00})
	salt, _ := crypto.S1([]byte("nkik"))
	p := append([]byte("id128"), 0x01)
	identityKey, _ := crypto.K1(key, salt, p)
	salt, _ = crypto.S1([]byte("nkbk"))
	beaconKey, _ := crypto.K1(key, salt, p)
	networkIdBytes, _ := crypto.K3(key)
	var networkId uint64
	utils.UnpackBE(networkIdBytes, "64", &networkId)
	return &NetKey{
		Bytes: key, Nid: nid, EncryptionKey: encryptionKey, PrivacyKey: privacyKey,
		IdentityKey: identityKey, BeaconKey: beaconKey, Index: index, NetworkId: networkId,
	}
}

func createAppKey(index uint) *AppKey {
	keyBytes := make([]byte, 16)
	rand.Read(keyBytes)
	return createAppKeyB(keyBytes, index)
}

func createAppKeyS(keyStr string, index uint) *AppKey {
	key, _ := hex.DecodeString(keyStr)
	return createAppKeyB(key, index)
}

func createAppKeyB(key []byte, index uint) *AppKey {
	aid, _ := crypto.K4(key)
	return &AppKey{
		Aid: aid, Bytes: key, Index: index,
	}
}

func Init(configDirectory string) {
	confDir = configDirectory
	meshDbRaw = &db.Mesh{}
	err := db.ReadFromDb(path.Join(confDir, "mesh.json"), meshDbRaw)
	if err != nil {
		loggerMesh.Fatal("error happened while reading from mesh database")
	}
	// loggerMesh.Debugf("%+#v", meshDbRaw)
	meshDb = &Mesh{
		NetKeys:        make(map[uint]*NetKey),
		AppKeys:        make(map[uint]*AppKey),
		Nodes:          make(map[uint]*Node),
		Groups:         make(map[uint]*Group),
		MeshName:       meshDbRaw.MeshName,
		UnicastAddress: utils.HexStringToUint(meshDbRaw.Provisioner.UnicastAddress),
		LowAddress:     utils.HexStringToUint(meshDbRaw.Provisioner.LowAddress),
		HighAddress:    utils.HexStringToUint(meshDbRaw.Provisioner.HighAddress),
		IVindex:        meshDbRaw.IVindex,
		IVupdate:       meshDbRaw.IVupdate,
		SequenceNumber: meshDbRaw.SequenceNumber,
	}

	for _, rawKey := range meshDbRaw.NetKeys {
		key := createNetKeyS(rawKey.Key, rawKey.Index)
		if rawKey.OldKey != "" {
			oldKey := createNetKeyS(rawKey.OldKey, rawKey.Index)
			key.OldKey = oldKey
			oldKey.NewKey = key
		}
		meshDb.NetKeys[rawKey.Index] = key
	}

	for _, rawKey := range meshDbRaw.AppKeys {
		key := createAppKeyS(rawKey.Key, rawKey.Index)
		if rawKey.OldKey != "" {
			oldKey := createAppKeyS(rawKey.OldKey, rawKey.Index)
			key.OldKey = oldKey
			oldKey.NewKey = key
		}
		meshDb.AppKeys[rawKey.Index] = key
	}

	// initialize nodes
	files, err := ioutil.ReadDir(confDir)
	if err != nil {
		loggerMesh.Fatal(err)
	}

	friends := map[*Node]uint{}
	for _, f := range files {
		if f.IsDir() {
			addr, _ := strconv.ParseUint(f.Name(), 16, 16)
			nodeRaw := &db.Node{}
			err = db.ReadFromDb(path.Join(confDir, f.Name(), "node.json"), nodeRaw)
			if err != nil {
				loggerMesh.Fatal("error happened while reading from node database")
			}
			node := Node{
				UnicastAddress: utils.HexStringToUint(nodeRaw.UnicastAddress),
				Cid:            nodeRaw.Cid,
				Pid:            nodeRaw.Pid,
				Vid:            nodeRaw.Vid,
				Crpl:           nodeRaw.Crpl,
				Features:       nodeRaw.Features,
				// IVindex:        nodeRaw.IVindex,
				SequenceNumber: nodeRaw.SequenceNumber,
				Mac:            nodeRaw.Mac,
				LPN:            nodeRaw.LPN,
				DefaultTTL:     nodeRaw.TTL,
				RelayState: def.ConfigRelayStatusMessageParameters{
					Relay:                        nodeRaw.Relay,
					RelayRetransmitCount:         nodeRaw.RelayRetransmitCount,
					RelayRetransmitIntervalSteps: nodeRaw.RelayRetransmitIntervalSteps,
				},
				AttentionTimer:       nodeRaw.AttentionTimer,
				SecureNetworkBeacon:  nodeRaw.SecureNetworkBeacon,
				GATTProxyState:       nodeRaw.GATTProxyState,
				FriendState:          nodeRaw.FriendState,
				KeyRefreshPhaseState: nodeRaw.KeyRefreshPhaseState,
				NetwrokTransmitState: def.ConfigNetworkTransmitStatusMessageParameters{
					NetworkTransmitCount:         nodeRaw.NetworkTransmitCount,
					NetworkTransmitIntervalSteps: nodeRaw.NetworkTransmitIntervalSteps,
				},
				CurrentFault:       nodeRaw.CurrentFault,
				NodeIdentityStates: map[uint]uint{},
			}
			utils.InitializeStruct(reflect.ValueOf(&node).Elem(), 1)
			if nodeRaw.Friend != "" {
				friends[&node] = utils.HexStringToUint(nodeRaw.Friend)
			}
			node.DeviceKey = DevKey{}
			node.DeviceKey.Bytes, _ = hex.DecodeString(nodeRaw.DeviceKey)
			node.DeviceKey.Aid = 0
			meshDb.Nodes[uint(addr)] = &node
			for _, binding := range nodeRaw.BindedNetKeys {
				if _, ok := meshDb.NetKeys[binding.NetKeyIndex]; ok {
					b := NodeKeyBinding{
						NetKeyIndex:     binding.NetKeyIndex,
						BindedAppKeyIds: []uint{},
					}
					for _, appkeyId := range binding.BindedAppKeys {
						if _, ok := meshDb.AppKeys[appkeyId]; ok {
							b.BindedAppKeyIds = append(b.BindedAppKeyIds, appkeyId)
						}
					}
					node.NodeIdentityStates[binding.NetKeyIndex] = binding.NodeIdentityState
					node.BindedKeys = append(node.BindedKeys, b)
				}
			}
			for _, element := range nodeRaw.Elements {
				meshEle := &Element{
					Node:           &node,
					Location:       element.Location,
					ElementIndex:   element.ElementIndex,
					UnicastAddress: utils.HexStringToUint(element.UnicastAddress),
					Models:         []*Model{},
				}
				for _, model := range element.Models {
					meshModel := &Model{
						Element:         meshEle,
						ModelID:         utils.HexStringToUint(model.ModelID),
						BindedAppKeyIds: []uint{},
						SubAddresses:    []uint{},
						PubSetting: def.ConfigModelPublicationStatusMessageParameters{
							PublishPeriod: def.PublishPeriodFormat{
								NumberOfSteps:  model.PubSetting.PublishNumberOfSteps,
								StepResolution: model.PubSetting.PublishStepResolution,
							},
							PublishRetransmitCount:         model.PubSetting.PublishRetransmitCount,
							PublishRetransmitIntervalSteps: model.PubSetting.PublishRetransmitIntervalSteps,
							PublishTTL:                     model.PubSetting.PublishTTL,
							CredentialFlag:                 model.PubSetting.CredentialFlag,
							AppKeyIndex:                    model.PubSetting.AppKeyIndex,
							PublishAddress:                 model.PubSetting.PublishAddress,
						},
					}
					stateType := modelStateUnmarshallMap[meshModel.ModelID]
					if model.State != "" && model.State != "null" {
						if stateType == nil {
							loggerMesh.Fatalf("no unmarshall type of model state, model: %+v", model)
						}
						val := reflect.New(stateType)
						err := json.Unmarshal([]byte(model.State), val.Interface())
						if err != nil {
							loggerMesh.Fatalf("failed to unmarshall state of model, model: %+v", model)
						} else {
							meshModel.State = val.Interface()
						}
					} else if stateType != nil {
						meshModel.State = reflect.New(stateType).Interface()
					}
					for _, keyIdx := range model.BindedAppKeys {
						if _, ok := meshDb.AppKeys[keyIdx]; ok {
							meshModel.BindedAppKeyIds = append(meshModel.BindedAppKeyIds, keyIdx)
						}
					}
					for _, addr := range model.SubAddresses {
						meshModel.SubAddresses = append(meshModel.SubAddresses, utils.HexStringToUint(addr))
					}
					meshEle.Models = append(meshEle.Models, meshModel)
				}
				node.Elements = append(node.Elements, meshEle)
			}
		}
	}

	for _, n := range meshDb.Nodes {
		if n.LPN {
			for _, f := range meshDb.Nodes {
				if f.UnicastAddress == friends[n] {
					n.Friend = f
					break
				}
			}
		}
	}

	for _, g := range meshDbRaw.Groups {
		addr, _ := strconv.ParseInt(g.GroupAddress, 16, 32)
		group := &Group{
			Address: uint(addr),
			Name:    g.Name,
		}
		meshDb.Groups[group.Address] = group
	}
	loggerMesh.Debugf("%+#v", meshDb)
}

func findNetKeyByNid(nid uint) []*NetKey {
	ret := []*NetKey{}
	for _, key := range meshDb.NetKeys {
		if key.Nid == nid || (key.OldKey != nil && key.OldKey.Nid == nid) {
			ret = append(ret, key)
		}
	}
	return ret
}

func findNetKeyByIndex(index uint) (*NetKey, error) {
	for _, key := range meshDb.NetKeys {
		if key.Index == index {
			return key, nil
		}
	}
	return nil, errors.InvalidNetKeyIndex.New().AddContextF("index:%d", index)
}

func findAppKeyByAid(aid uint) []*AppKey {
	ret := []*AppKey{}
	for _, key := range meshDb.AppKeys {
		if key.Aid == aid || (key.OldKey != nil && key.OldKey.Aid == aid) {
			ret = append(ret, key)
		}
	}
	for _, n := range meshDb.Nodes {
		if n.DeviceKey.Aid == aid {
			ret = append(ret, &n.DeviceKey.AppKey)
		}
	}
	return ret
}

func findAppKeyByIndex(index uint) (*AppKey, error) {
	for _, key := range meshDb.AppKeys {
		if key.Index == index {
			return key, nil
		}
	}
	return nil, errors.InvalidAppKeyIndex.New().AddContextF("index:%d", index)
}

func findNodeByAddr(addr uint) (*Node, error) {
	for _, n := range meshDb.Nodes {
		if n.UnicastAddress == addr {
			return n, nil
		}
		for _, e := range n.Elements {
			if e.UnicastAddress == addr {
				return n, nil
			}
		}
	}
	return nil, errors.NodeAddressNotFound.New().AddContextF("address:%4x", addr)
}

func (m *Model) findBindedAppKey(index uint) (*AppKey, error) {
	for _, k := range m.BindedAppKeyIds {
		if k == index {
			return findAppKeyByIndex(k)
		}
	}
	return nil, errors.ModelAppKeyBindingNotFound.New().AddContextF("model:%4x, appkeyIndex:%d", m.ModelID, index)
}

func findModelDirectly(elementAddr uint, modelId uint) (*Model, error) {
	ele, err := findElementByAddr(elementAddr)
	if err != nil {
		return nil, err
	}
	return ele.findModel(modelId)
}

func (e *Element) findModel(id uint) (*Model, error) {
	for _, m := range e.Models {
		if m.ModelID == id {
			return m, nil
		}
	}
	return nil, errors.ModelNotFound.New().AddContextF("element: %4x, model:%4x", e.UnicastAddress, id)
}

func findElementByAddr(addr uint) (*Element, error) {
	for _, n := range meshDb.Nodes {
		for _, e := range n.Elements {
			if e.UnicastAddress == addr {
				return e, nil
			}
		}
	}
	return nil, errors.ElementNotFound.New().AddContextF("element address: %4x", addr)
}

func (node *Node) findNodeNetKeyByAppKeyIndex(appKeyId uint) (*NetKey, error) {
	for _, b := range node.BindedKeys {
		for _, a := range b.BindedAppKeyIds {
			if a == appKeyId {
				return findNetKeyByIndex(b.NetKeyIndex)
			}
		}
	}
	return nil, errors.NoNetKeyBindedToAppKey.New().AddContextF("appkeyIndex:%d", appKeyId)
}

func (node *Node) findNodeNetKeyByIndex(index uint) (*NetKey, error) {
	for _, key := range node.BindedKeys {
		if key.NetKeyIndex == index {
			return findNetKeyByIndex(index)
		}
	}
	return nil, errors.NetKeyNotBindedToNode.New().AddContextF("node:%4x, netkeyIndex:%d", node.UnicastAddress, index)
}

func (node *Node) findNodeAppKeyByIndex(index uint) (*AppKey, error) {
	for _, k := range node.BindedKeys {
		for _, ak := range k.BindedAppKeyIds {
			if ak == index {
				return findAppKeyByIndex(index)
			}
		}
	}
	return nil, errors.AppKeyNotBindedToNode.New().AddContextF("node:%4x, appkeyIndex:%d", node.UnicastAddress, index)
}

func (node *Node) findNodeKeyBindingByNetKeyIndex(index uint) (*NodeKeyBinding, error) {
	for _, bd := range node.BindedKeys {
		if bd.NetKeyIndex == index {
			return &bd, nil
		}
	}
	return nil, errors.NetKeyNotBindedToNode.New().AddContextF("node:%4x, appkeyIndex:%d", node.UnicastAddress, index)
}

func (node *Node) findLpnFriends() []*Node {
	list := []*Node{}
	for _, n := range meshDb.Nodes {
		if n.Friend == node {
			list = append(list, n)
		}
	}
	return list
}

func writeMeshToDb() {
	meshDbRaw.IVindex = meshDb.IVindex
	meshDbRaw.IVupdate = meshDb.IVupdate
	meshDbRaw.SequenceNumber = meshDb.SequenceNumber
	for _, netKey := range meshDb.NetKeys {
		for i := 0; i < len(meshDbRaw.NetKeys); i++ {
			if meshDbRaw.NetKeys[i].Index == netKey.Index {
				meshDbRaw.NetKeys[i].KeyRefreshPhase = netKey.KeyRefreshPhase
				meshDbRaw.NetKeys[i].Key = hex.EncodeToString(netKey.Bytes)
			}
		}
	}
	for _, appKey := range meshDb.AppKeys {
		for i := 0; i < len(meshDbRaw.AppKeys); i++ {
			if meshDbRaw.AppKeys[i].Index == appKey.Index {
				meshDbRaw.AppKeys[i].Key = hex.EncodeToString(appKey.Bytes)
			}
		}
	}
	pathDir := path.Join(confDir, "mesh.json")
	db.WriteToDb(pathDir, meshDbRaw)
}

func writeNodeToDb(node *Node) {
	nodeRaw := &db.Node{
		UUID:                         node.UUID,
		Cid:                          node.Cid,
		Pid:                          node.Pid,
		Vid:                          node.Vid,
		Crpl:                         node.Crpl,
		Features:                     node.Features,
		SequenceNumber:               node.SequenceNumber,
		LPN:                          node.LPN,
		TTL:                          node.DefaultTTL,
		Relay:                        node.RelayState.Relay,
		RelayRetransmitCount:         node.RelayState.RelayRetransmitCount,
		RelayRetransmitIntervalSteps: node.RelayState.RelayRetransmitIntervalSteps,
		AttentionTimer:               node.AttentionTimer,
		SecureNetworkBeacon:          node.SecureNetworkBeacon,
		GATTProxyState:               node.GATTProxyState,
		FriendState:                  node.FriendState,
		KeyRefreshPhaseState:         node.KeyRefreshPhaseState,
		NetworkTransmitCount:         node.NetwrokTransmitState.NetworkTransmitCount,
		NetworkTransmitIntervalSteps: node.NetwrokTransmitState.NetworkTransmitIntervalSteps,
		CurrentFault:                 node.CurrentFault,
	}
	if node.Friend != nil {
		nodeRaw.Friend = strconv.FormatUint(uint64(node.Friend.UnicastAddress), 16)
	}
	nodeRaw.DeviceKey = hex.EncodeToString(node.DeviceKey.Bytes)
	nodeRaw.UnicastAddress = strconv.FormatUint(uint64(node.UnicastAddress), 16)
	nodeRaw.BindedNetKeys = []db.BindedNetKey{}
	for _, key := range node.BindedKeys {
		bindedNetKey := db.BindedNetKey{
			NetKeyIndex:       key.NetKeyIndex,
			BindedAppKeys:     []uint{},
			NodeIdentityState: node.NodeIdentityStates[key.NetKeyIndex],
		}
		for _, k := range key.BindedAppKeyIds {
			bindedNetKey.BindedAppKeys = append(bindedNetKey.BindedAppKeys, k)
		}

		nodeRaw.BindedNetKeys = append(nodeRaw.BindedNetKeys, bindedNetKey)
	}

	if nodeRaw.Elements == nil || len(nodeRaw.Elements) == 0 {
		nodeRaw.Elements = []db.Element{}
		for _, ele := range node.Elements {
			eleRaw := db.Element{
				ElementIndex:   ele.ElementIndex,
				Location:       ele.Location,
				UnicastAddress: strconv.FormatUint(uint64(ele.UnicastAddress), 16),
				Models:         []db.Model{},
			}
			for _, m := range ele.Models {
				state, _ := json.Marshal(m.State)
				mRaw := db.Model{
					ModelID: utils.UintToHexString(m.ModelID),
					PubSetting: db.PubSetting{
						PublishAddress:                 m.PubSetting.PublishAddress,
						AppKeyIndex:                    m.PubSetting.AppKeyIndex,
						CredentialFlag:                 m.PubSetting.CredentialFlag,
						PublishTTL:                     m.PubSetting.PublishTTL,
						PublishNumberOfSteps:           m.PubSetting.PublishPeriod.NumberOfSteps,
						PublishStepResolution:          m.PubSetting.PublishPeriod.StepResolution,
						PublishRetransmitCount:         m.PubSetting.PublishRetransmitCount,
						PublishRetransmitIntervalSteps: m.PubSetting.PublishRetransmitIntervalSteps,
					},
					SubAddresses: []string{},
					State:        string(state),
				}
				mRaw.BindedAppKeys = []uint{}
				for _, a := range m.BindedAppKeyIds {
					mRaw.BindedAppKeys = append(mRaw.BindedAppKeys, a)
				}
				for _, a := range m.SubAddresses {
					mRaw.SubAddresses = append(mRaw.SubAddresses, utils.UintToHexString(a))
				}

				eleRaw.Models = append(eleRaw.Models, mRaw)
			}
			nodeRaw.Elements = append(nodeRaw.Elements, eleRaw)
		}
	}

	pathNodeDir := path.Join(confDir, nodeRaw.UnicastAddress)
	if _, err := os.Stat(pathNodeDir); os.IsNotExist(err) {
		os.Mkdir(pathNodeDir, 0755)
	}

	pathNode := path.Join(pathNodeDir, "node.json")
	db.WriteToDb(pathNode, nodeRaw)
}

func deleteNode(addr uint) {
	pathNodeDir := path.Join(confDir, strconv.FormatUint(uint64(addr), 16))
	os.RemoveAll(pathNodeDir)
}

func GetDb() *Mesh {
	return meshDb
}

func GetNode(addr uint) *Node {
	n, _ := findNodeByAddr(addr)
	return n
}
