package mesh

import (
	"ble-mesh/mesh/db"
	. "ble-mesh/mesh/def"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"bytes"
	"encoding/binary"

	"github.com/jinzhu/copier"
	funk "github.com/thoas/go-funk"
)

type (
	Features struct {
		Relay, Proxy, Friend, LowPower bool
	}

	CompositionElement struct {
		Location     uint16
		SigModels    []uint16
		VendorModels []uint32
	}

	Composition struct {
		Page                byte
		CID, PID, VID, CRPL uint16
		Features            Features
		Elements            []CompositionElement
	}
)

// opcodes
const (
	STATUS_SUCCESS = 0x00

	opConfigAppKeyAdd                                = 0x00
	opConfigAppKeyDelete                             = 0x8000
	opConfigAppKeyGet                                = 0x8001
	opConfigAppKeyList                               = 0x8002
	opConfigAppKeyStatus                             = 0x8003
	opConfigAppKeyUpdate                             = 0x01
	opConfigBeaconGet                                = 0x8009
	opConfigBeaconSet                                = 0x800A
	opConfigBeaconStatus                             = 0x800B
	opConfigCompositionDataGet                       = 0x8008
	opConfigCompositionDataStatus                    = 0x02
	opConfigModelPublicationSet                      = 0x03
	opConfigDefaultTTLGet                            = 0x800C
	opConfigDefaultTTLSet                            = 0x800D
	opConfigDefaultTTLStatus                         = 0x800E
	opConfigFriendGet                                = 0x800F
	opConfigFriendSet                                = 0x8010
	opConfigFriendStatus                             = 0x8011
	opConfigGATTProxyGet                             = 0x8012
	opConfigGATTProxySet                             = 0x8013
	opConfigGATTProxyStatus                          = 0x8014
	opConfigHeartbeatPublicationGet                  = 0x8038
	opConfigHeartbeatPublicationSet                  = 0x8039
	opConfigHeartbeatPublicationStatus               = 0x06
	opConfigHeartbeatSubscriptionGet                 = 0x803A
	opConfigHeartbeatSubscriptionSet                 = 0x803B
	opConfigHeartbeatSubscriptionStatus              = 0x803C
	opConfigKeyRefreshPhaseGet                       = 0x8015
	opConfigKeyRefreshPhaseSet                       = 0x8016
	opConfigKeyRefreshPhaseStatus                    = 0x8017
	opConfigLowPowerNodePollTimeoutGet               = 0x802D
	opConfigLowPowerNodePollTimeoutStatus            = 0x802E
	opConfigModelAppBind                             = 0x803D
	opConfigModelAppStatus                           = 0x803E
	opConfigModelAppUnbind                           = 0x803F
	opConfigModelPublicationGet                      = 0x8018
	opConfigModelPublicationStatus                   = 0x8019
	opConfigModelPublicationVirtualAddressSet        = 0x801A
	opConfigModelSubscriptionAdd                     = 0x801B
	opConfigModelSubscriptionDelete                  = 0x801C
	opConfigModelSubscriptionDeleteAll               = 0x801D
	opConfigModelSubscriptionOverwrite               = 0x801E
	opConfigModelSubscriptionStatus                  = 0x801F
	opConfigModelSubscriptionVirtualAddressAdd       = 0x8020
	opConfigModelSubscriptionVirtualAddressDelete    = 0x8021
	opConfigModelSubscriptionVirtualAddressOverwrite = 0x8022
	opConfigNetKeyAdd                                = 0x8040
	opConfigNetKeyDelete                             = 0x8041
	opConfigNetKeyGet                                = 0x8042
	opConfigNetKeyList                               = 0x8043
	opConfigNetKeyStatus                             = 0x8044
	opConfigNetKeyUpdate                             = 0x8045
	opConfigNetworkTransmitGet                       = 0x8023
	opConfigNetworkTransmitSet                       = 0x8024
	opConfigNetworkTransmitStatus                    = 0x8025
	opConfigNodeIdentityGet                          = 0x8046
	opConfigNodeIdentitySet                          = 0x8047
	opConfigNodeIdentityStatus                       = 0x8048
	opConfigNodeReset                                = 0x8049
	opConfigNodeResetStatus                          = 0x804A
	opConfigRelayGet                                 = 0x8026
	opConfigRelaySet                                 = 0x8027
	opConfigRelayStatus                              = 0x8028
	opConfigSIGModelAppGet                           = 0x804B
	opConfigSIGModelAppList                          = 0x804C
	opConfigSIGModelSubscriptionGet                  = 0x8029
	opConfigSIGModelSubscriptionList                 = 0x802A
	opConfigVendorModelAppGet                        = 0x804D
	opConfigVendorModelAppList                       = 0x804E
	opConfigVendorModelSubscriptionGet               = 0x802B
	opConfigVendorModelSubscriptionList              = 0x802C
	opHealthAttentionGet                             = 0x8004
	opHealthAttentionSet                             = 0x8005
	opHealthAttentionSetUnacknowledged               = 0x8006
	opHealthAttentionStatus                          = 0x8007
	opHealthCurrentStatus                            = 0x04
	opHealthFaultClear                               = 0x802F
	opHealthFaultClearUnacknowledged                 = 0x8030
	opHealthFaultGet                                 = 0x8031
	opHealthFaultStatus                              = 0x05
	opHealthFaultTest                                = 0x8032
	opHealthFaultTestUnacknowledged                  = 0x8033
	opHealthPeriodGet                                = 0x8034
	opHealthPeriodSet                                = 0x8035
	opHealthPeriodSetUnacknowledged                  = 0x8036
	opHealthPeriodStatus                             = 0x8037
)

var (
	loggerFoundation  = utils.CreateLogger("ConfigClient")
	foundationMethods = map[uint]string{
		opConfigAppKeyAdd:                                "ConfigAppKeyAdd",
		opConfigAppKeyDelete:                             "ConfigAppKeyDelete",
		opConfigAppKeyGet:                                "ConfigAppKeyGet",
		opConfigAppKeyUpdate:                             "ConfigAppKeyUpdate",
		opConfigBeaconGet:                                "ConfigBeaconGet",
		opConfigBeaconSet:                                "ConfigBeaconSet",
		opConfigCompositionDataGet:                       "ConfigCompositionDataGet",
		opConfigModelPublicationSet:                      "ConfigModelPublicationSet",
		opConfigDefaultTTLGet:                            "ConfigDefaultTTLGet",
		opConfigDefaultTTLSet:                            "ConfigDefaultTTLSet",
		opConfigFriendGet:                                "ConfigFriendGet",
		opConfigFriendSet:                                "ConfigFriendSet",
		opConfigGATTProxyGet:                             "ConfigGATTProxyGet",
		opConfigGATTProxySet:                             "ConfigGATTProxySet",
		opConfigHeartbeatPublicationGet:                  "ConfigHeartbeatPublicationGet",
		opConfigHeartbeatPublicationSet:                  "ConfigHeartbeatPublicationSet",
		opConfigHeartbeatSubscriptionGet:                 "ConfigHeartbeatSubscriptionGet",
		opConfigHeartbeatSubscriptionSet:                 "ConfigHeartbeatSubscriptionSet",
		opConfigKeyRefreshPhaseGet:                       "ConfigKeyRefreshPhaseGet",
		opConfigKeyRefreshPhaseSet:                       "ConfigKeyRefreshPhaseSet",
		opConfigLowPowerNodePollTimeoutGet:               "ConfigLowPowerNodePollTimeoutGet",
		opConfigModelAppBind:                             "ConfigModelAppBind",
		opConfigModelAppUnbind:                           "ConfigModelAppUnbind",
		opConfigModelPublicationGet:                      "ConfigModelPublicationGet",
		opConfigModelPublicationVirtualAddressSet:        "ConfigModelPublicationVirtualAddressSet",
		opConfigModelSubscriptionAdd:                     "ConfigModelSubscriptionAdd",
		opConfigModelSubscriptionDelete:                  "ConfigModelSubscriptionDelete",
		opConfigModelSubscriptionDeleteAll:               "ConfigModelSubscriptionDeleteAll",
		opConfigModelSubscriptionOverwrite:               "ConfigModelSubscriptionOverwrite",
		opConfigModelSubscriptionVirtualAddressAdd:       "ConfigModelSubscriptionVirtualAddressAdd",
		opConfigModelSubscriptionVirtualAddressDelete:    "ConfigModelSubscriptionVirtualAddressDelete",
		opConfigModelSubscriptionVirtualAddressOverwrite: "ConfigModelSubscriptionVirtualAddressOverwrite",
		opConfigNetKeyAdd:                                "ConfigNetKeyAdd",
		opConfigNetKeyDelete:                             "ConfigNetKeyDelete",
		opConfigNetKeyGet:                                "ConfigNetKeyGet",
		opConfigNetKeyUpdate:                             "ConfigNetKeyUpdate",
		opConfigNetworkTransmitGet:                       "ConfigNetworkTransmitGet",
		opConfigNetworkTransmitSet:                       "ConfigNetworkTransmitSet",
		opConfigNodeIdentityGet:                          "ConfigNodeIdentityGet",
		opConfigNodeIdentitySet:                          "ConfigNodeIdentitySet",
		opConfigNodeReset:                                "ConfigNodeReset",
		opConfigRelayGet:                                 "ConfigRelayGet",
		opConfigRelaySet:                                 "ConfigRelaySet",
		opConfigSIGModelAppGet:                           "ConfigSIGModelAppGet",
		opConfigSIGModelSubscriptionGet:                  "ConfigSIGModelSubscriptionGet",
		opConfigVendorModelAppGet:                        "ConfigVendorModelAppGet",
		opConfigVendorModelSubscriptionGet:               "ConfigVendorModelSubscriptionGet",
		opHealthAttentionGet:                             "HealthAttentionGet",
		opHealthAttentionSet:                             "HealthAttentionSet",
		opHealthAttentionSetUnacknowledged:               "HealthAttentionSetUnacknowledged",
		opHealthFaultClear:                               "HealthFaultClear",
		opHealthFaultClearUnacknowledged:                 "HealthFaultClearUnacknowledged",
		opHealthFaultGet:                                 "HealthFaultGet",
		opHealthFaultTest:                                "HealthFaultTest",
		opHealthFaultTestUnacknowledged:                  "HealthFaultTestUnacknowledged",
		opHealthPeriodGet:                                "HealthPeriodGet",
		opHealthPeriodSet:                                "HealthPeriodSet",
		opHealthPeriodSetUnacknowledged:                  "HealthPeriodSetUnacknowledged",
	}
)

func checkRespStatusCode(code byte) bool {
	if code == 0 {
		return true
	}
	loggerFoundation.Errorf("error: %s", statusCode[code])
	return false
}

func handleBeaconResponse(n *Node, d interface{}) error {
	d1 := d.(ConfigBeaconStatusMessageParameters)
	n.SecureNetworkBeacon = d1.Beacon == 0x01
	return nil
}

func ConfigBeaconGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigBeaconGet, nil, handleBeaconResponse)
}

func ConfigBeaconSet(dst uint, enable bool) error {
	beacon := uint(0)
	if enable {
		beacon = 1
	}
	params := &ConfigBeaconSetMessageParameters{Beacon: beacon}
	return modelSendTmplParsed(false, dst, opConfigBeaconSet, params, handleBeaconResponse)
}

func handleCompData(src uint, data []byte) error {
	node, err := findNodeByAddr(src)
	if err != nil {
		return err
	}
	comp := parseCompData(data)
	loggerFoundation.Infof("device composition data: %+#v", comp)
	node.Cid = int(comp.CID)
	node.Pid = int(comp.PID)
	node.Vid = int(comp.VID)
	node.Crpl = int(comp.CRPL)
	node.Features = db.Features{
		Proxy:  comp.Features.Proxy,
		Friend: comp.Features.Friend,
		Relay:  comp.Features.Relay,
		Lpn:    comp.Features.LowPower,
	}
	if node.Elements == nil || len(node.Elements) == 0 {
		node.Elements = []*Element{}
		var addr uint
		if comp.Page == 0 {
			addr = src
		}
		for i, ele := range comp.Elements {
			meshEle := &Element{
				ElementIndex:   i,
				Location:       int(ele.Location),
				UnicastAddress: addr + uint(i),
				Models:         []*Model{},
			}
			for _, m := range ele.SigModels {
				meshEle.Models = append(meshEle.Models, &Model{
					ModelID:         uint(m),
					SubAddresses:    []uint{},
					BindedAppKeyIds: []uint{},
				})
			}
			for _, m := range ele.VendorModels {
				meshEle.Models = append(meshEle.Models, &Model{
					ModelID:         uint(m),
					SubAddresses:    []uint{},
					BindedAppKeyIds: []uint{},
				})
			}
			node.Elements = append(node.Elements, meshEle)
		}
	}
	return nil
	// writeNodeToDb(node)
}

func parseCompData(data []byte) *Composition {
	comp := &Composition{}
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.LittleEndian, &comp.Page)
	binary.Read(buffer, binary.LittleEndian, &comp.CID)
	binary.Read(buffer, binary.LittleEndian, &comp.PID)
	binary.Read(buffer, binary.LittleEndian, &comp.VID)
	binary.Read(buffer, binary.LittleEndian, &comp.CRPL)
	var feature uint16
	binary.Read(buffer, binary.LittleEndian, &feature)
	comp.Features = Features{
		Relay:    feature&0x01 > 0,
		Proxy:    feature&0x02 > 0,
		Friend:   feature&0x04 > 0,
		LowPower: feature&0x08 > 0,
	}
	comp.Elements = []CompositionElement{}
	nRest := len(data) - 11
	for nRest > 0 {
		var location uint16
		var numS, numV byte
		binary.Read(buffer, binary.LittleEndian, &location)
		binary.Read(buffer, binary.LittleEndian, &numS)
		binary.Read(buffer, binary.LittleEndian, &numV)
		element := CompositionElement{Location: location}
		element.SigModels = []uint16{}
		element.VendorModels = []uint32{}
		var sigId uint16
		var vendId uint32
		for i := 0; i < int(numS); i++ {
			binary.Read(buffer, binary.LittleEndian, &sigId)
			element.SigModels = append(element.SigModels, sigId)
		}
		for i := 0; i < int(numV); i++ {
			binary.Read(buffer, binary.LittleEndian, &vendId)
			element.VendorModels = append(element.VendorModels, vendId)
		}
		nRest -= 4 + int(numS)*2 + int(numV)*4
		comp.Elements = append(comp.Elements, element)
	}
	return comp
}

func ConfigCompositionDataGet(dst uint, page uint) error {
	p := &ConfigCompositionDataGetMessageParameters{Page: page}
	payload, _ := utils.PackStructLE(p)
	return modelSendTmpl(false, dst, opConfigCompositionDataGet, payload, func(msg *AccessMessage) error {
		return handleCompData(msg.src, msg.payload)
	})
}

func handleTtlResponse(n *Node, d interface{}) error {
	n.DefaultTTL = d.(ConfigDefaultTTLStatusMessageParameters).TTL
	return nil
}

func ConfigDefaultTTLGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigDefaultTTLGet, nil, handleTtlResponse)
}

func ConfigDefaultTTLSet(dst uint, ttl uint) error {
	params := &ConfigDefaultTTLSetMessageParameters{TTL: ttl}
	if ttl == 0x01 || ttl > 0x80 {
		return errors.WrongTTLSetting.New()
	}
	return modelSendTmplParsed(false, dst, opConfigDefaultTTLSet, params, handleTtlResponse)
}

func handleProxyResponse(n *Node, d interface{}) error {
	n.GATTProxyState = d.(ConfigGATTProxyStatusMessageParameters).GATTProxy
	return nil
}

func ConfigGattProxyGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigGATTProxyGet, nil, handleProxyResponse)
}

func ConfigGattProxySet(dst uint, status uint) error {
	params := &ConfigGATTProxySetMessageParameters{GATTProxy: status}
	if status > 0x03 {
		return errors.WrongGattProxySetting.New()
	}
	return modelSendTmplParsed(false, dst, opConfigGATTProxySet, params, handleProxyResponse)
}

func handleRelayResponse(n *Node, d interface{}) error {
	resp := d.(ConfigRelayStatusMessageParameters)
	n.RelayState = resp
	return nil
}

func ConfigRelayGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigRelayGet, nil, handleRelayResponse)
}

func ConfigRelaySet(dst uint, relay, cnt, step uint) error {
	params := &ConfigRelaySetMessageParameters{
		Relay:                        relay,
		RelayRetransmitCount:         cnt,
		RelayRetransmitIntervalSteps: step,
	}
	return modelSendTmplParsed(false, dst, opConfigRelaySet, params, handleRelayResponse)
}

// func (m *AccessMessage) publicationStatus() {
// 	statusLiteral := map[uint]string{
// 		0x00: "Success",
// 		0x01: "Invalid Address",
// 		0x02: "Invalid Model",
// 		0x03: "Invalid AppKey Index",
// 		0x04: "Invalid NetKey Index",
// 		0x05: "Insufficient Resources",
// 		0x06: "Key Index Already Stored",
// 		0x07: "Invalid Publish Parameters",
// 		0x08: "Not a Subscribe Model",
// 		0x09: "Storage Failure",
// 		0x0A: "Feature Not Supported",
// 		0x0B: "Cannot Update",
// 		0x0C: "Cannot Remove",
// 		0x0D: "Cannot Bind",
// 		0x0E: "Temporarily Unable to Change State",
// 		0x0F: "Cannot Set",
// 		0x10: "Unspecified Error",
// 		0x11: "Invalid Binding",
// 	}
// 	credFlagLiteral := map[uint]string{
// 		0: "Master security material is used for Publishing",
// 		1: "Friendship security material is used for Publishin",
// 	}
// 	var status, elementAddress, publishAddress, appKeyIndex, credentialFlag, publishTTL,
// 		publishPeriod, publishRetransmitCount, publishRetransmitIntervalStep, modelIdentifier uint
// 	bitString := "8,16,16,12,1,03,8,8,3,5"
// 	utils.UnpackLE(m.payload, bitString, &status, &elementAddress, &publishAddress, &appKeyIndex, &credentialFlag,
// 		&publishTTL, &publishPeriod, &publishRetransmitCount, &publishRetransmitIntervalStep)
// 	if len(m.payload) == 12 {
// 		bitString += "16"
// 	} else if len(m.payload) == 14 {
// 		bitString += "32"
// 	} else {
// 		loggerFoundation.Error("received wrong size of publication status")
// 	}
// 	utils.UnpackLE(m.payload[10:], bitString, &modelIdentifier)
// 	loggerFoundation.Infof("status: %d, %s", status, statusLiteral[status])
// 	loggerFoundation.Infof("elementAddress: %4x", elementAddress)
// 	loggerFoundation.Infof("publishAddress: %4x", publishAddress)
// 	loggerFoundation.Infof("appKeyIndex: %d", appKeyIndex)
// 	loggerFoundation.Infof("credential flag: %d, %s", credentialFlag, credFlagLiteral[credentialFlag])
// 	loggerFoundation.Infof("publish ttl: %d", publishTTL)
// 	loggerFoundation.Infof("publish period: %d", publishPeriod)
// 	loggerFoundation.Infof("publish restransmit count: %d", publishRetransmitCount)
// 	loggerFoundation.Infof("publish restransmit interval step: %d", publishRetransmitIntervalStep)
// 	loggerFoundation.Infof("model identifier: %8x", modelIdentifier)
// }

func handlePublicationResponse(n *Node, d interface{}) error {
	resp := d.(ConfigModelPublicationStatusMessageParameters)
	ele, err := findElementByAddr(resp.ElementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(resp.ModelIdentifier)
	if err != nil {
		return err
	}
	if resp.Status == STATUS_SUCCESS {
		m.PubSetting = resp
		return nil
	}
	return errors.InvalidResponse.New()
}

func ConfigModelPublicationGet(elementAddress uint, modelId uint) error {
	params := &ConfigModelPublicationGetMessageParameters{
		ElementAddress:  elementAddress,
		ModelIdentifier: modelId,
	}
	node, _ := findNodeByAddr(elementAddress)
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelPublicationGet, params, handlePublicationResponse)
}

func ConfigModelPublicationSet(elementAddress, publishAddress, appKeyIndex, credentialFlag, publishTTL,
	numSteps, stepResolution, publishRetransmitCount, publishRetransmitIntervalStep, modelIdentifier uint) error {
	params := &ConfigModelPublicationSetMessageParameters{
		ElementAddress: elementAddress,
		PublishAddress: publishAddress,
		AppKeyIndex:    appKeyIndex,
		CredentialFlag: credentialFlag,
		PublishTTL:     publishTTL,
		PublishPeriod: PublishPeriodFormat{
			NumberOfSteps:  numSteps,
			StepResolution: stepResolution,
		},
		PublishRetransmitCount:         publishRetransmitCount,
		PublishRetransmitIntervalSteps: publishRetransmitIntervalStep,
		ModelIdentifier:                modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelPublicationSet, params, handlePublicationResponse)
}

// func ConfigModelPublicationVirtualAddressSet(elementAddress uint, publishAddress uuid.UUID, appKeyIndex, credentialFlag, publishTTL,
// 	publishPeriod, publishRetransmitCount, publishRetransmitIntervalStep, modelIdentifier uint) {
// 	node, _ := findNodeByAddr(elementAddress)
// 	bitString := "16,B,12,1,03,8,8,3,5," + strconv.Itoa(int(getModelIdBitLength(modelIdentifier)))
// 	virtualAddr, _ := publishAddress.MarshalBinary()
// 	payload, _ := utils.PackLE(bitString, elementAddress, virtualAddr, appKeyIndex, credentialFlag,
// 		publishTTL, publishPeriod, publishRetransmitCount, publishRetransmitIntervalStep, modelIdentifier)
// 	modelSendTmpl(node.UnicastAddress, opConfigModelPublicationVirtualAddressSet, payload, func(msg *AccessMessage) {
// 		msg.publicationStatus()
// 	})
// }

// func (m *AccessMessage) subscriptionStatus() {
// 	var status, elementAddress, address, modelId uint
// 	bitString := "8,16,16,"

// 	if len(m.payload) == 7 {
// 		bitString += "16"
// 	} else if len(m.payload) == 9 {
// 		bitString += "32"
// 	} else {
// 		loggerFoundation.Error("received wrong size of subscription status")
// 	}
// 	err := utils.UnpackLE(m.payload, bitString, &status, &elementAddress, &address, &modelId)
// 	if err != nil {
// 		loggerFoundation.Error(err)
// 	}
// 	loggerFoundation.Infof("got subscription status: %s, element address: %4x, address: %4x, model id: %8x",
// 		statusCode[status], elementAddress, address, modelId)
// }

func ConfigModelSubscriptionAdd(elementAddress, address, modelIdentifier uint) error {
	params := &ConfigModelSubscriptionAddMessageParameters{
		ElementAddress:  elementAddress,
		Address:         address,
		ModelIdentifier: modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)
	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	if utils.Contains(m.SubAddresses, address) {
		return errors.AddressAlreadyInSubscriptionList.New()
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelSubscriptionAdd, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelSubscriptionStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			m.SubAddresses = append(m.SubAddresses, resp.Address)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigModelSubscriptionDelete(elementAddress, address, modelIdentifier uint) error {
	params := &ConfigModelSubscriptionDeleteMessageParameters{
		ElementAddress:  elementAddress,
		Address:         address,
		ModelIdentifier: modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)

	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	if !utils.Contains(m.SubAddresses, address) {
		return errors.AddressNotInSubscriptionList.New()
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelSubscriptionDelete, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelSubscriptionStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			idx := 0
			for i, a := range m.SubAddresses {
				if a == address {
					idx = i
					m.SubAddresses = append(m.SubAddresses[:idx], m.SubAddresses[idx+1:]...)
					return nil
				}
			}
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigModelSubscriptionOverwrite(elementAddress, address, modelIdentifier uint) error {
	params := &ConfigModelSubscriptionOverwriteMessageParameters{
		ElementAddress:  elementAddress,
		Address:         address,
		ModelIdentifier: modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)

	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelSubscriptionOverwrite, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelSubscriptionStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			m.SubAddresses = append([]uint{}, resp.Address)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

// func ConfigModelSubscriptionVirtualAddressAdd(elementAddress uint, address uuid.UUID, modelIdentifier uint) {
// 	node, _ := findNodeByAddr(elementAddress)
// 	bitString := "16,B," + strconv.Itoa(int(getModelIdBitLength(modelIdentifier)))
// 	virtualAddr, _ := address.MarshalBinary()
// 	payload, _ := utils.PackLE(bitString, elementAddress, virtualAddr, modelIdentifier)
// 	modelSendTmpl(node.UnicastAddress, opConfigModelSubscriptionVirtualAddressDelete, payload, func(msg *AccessMessage) {
// 		msg.subscriptionStatus()
// 	})
// }

// func ConfigModelSubscriptionVirtualAddressDelete(elementAddress uint, address uuid.UUID, modelIdentifier uint) {
// 	node, _ := findNodeByAddr(elementAddress)
// 	bitString := "16,B," + strconv.Itoa(int(getModelIdBitLength(modelIdentifier)))
// 	virtualAddr, _ := address.MarshalBinary()
// 	payload, _ := utils.PackLE(bitString, elementAddress, virtualAddr, modelIdentifier)
// 	modelSendTmpl(node.UnicastAddress, opConfigModelSubscriptionVirtualAddressAdd, payload, func(msg *AccessMessage) {
// 		msg.subscriptionStatus()
// 	})
// }

// func ConfigModelSubscriptionVirtualAddressOverwrite(elementAddress uint, address uuid.UUID, modelIdentifier uint) {
// 	node, _ := findNodeByAddr(elementAddress)
// 	bitString := "16,B," + strconv.Itoa(int(getModelIdBitLength(modelIdentifier)))
// 	virtualAddr, _ := address.MarshalBinary()
// 	payload, _ := utils.PackLE(bitString, elementAddress, virtualAddr, modelIdentifier)
// 	modelSendTmpl(node.UnicastAddress, opConfigModelSubscriptionVirtualAddressOverwrite, payload, func(msg *AccessMessage) {
// 		msg.subscriptionStatus()
// 	})
// }

func ConfigModelSubscriptionDeleteAll(elementAddress, modelIdentifier uint) error {
	params := &ConfigModelSubscriptionDeleteAllMessageParameters{
		ElementAddress:  elementAddress,
		ModelIdentifier: modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)

	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelSubscriptionDeleteAll, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelSubscriptionStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			m.SubAddresses = []uint{}
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigSigModelSubscriptionGet(elementAddress, modelIdentifier uint) error {
	params := &ConfigSIGModelAppGetMessageParameters{
		ElementAddress:  elementAddress,
		ModelIdentifier: modelIdentifier,
	}

	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	node, _ := findNodeByAddr(elementAddress)
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigSIGModelSubscriptionGet, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigSIGModelSubscriptionListMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			vendorAddrs := []uint{}
			for _, addr := range m.SubAddresses {
				if addr > 0xFFFF {
					vendorAddrs = append(vendorAddrs, addr)
				}
			}
			m.SubAddresses = append(resp.Addresses, vendorAddrs...)
			return nil
		}
		return errors.InvalidResponse.New()
		// virtual address is not considered
		// var status, elementAddr, modelId uint
		// utils.UnpackLE(msg.payload, "8,16,16", &status, &elementAddr, &modelId)
		// p := msg.payload[5:]
		// ids := make([]uint16, uint(len(msg.payload)-5)/2)
		// for i := 0; i < len(ids); i++ {
		// 	ids[i] = binary.LittleEndian.Uint16(p[i : i+2])
		// }
		// loggerFoundation.Infof("subscriptions: status: %s, element address: %4x, model: %4x, id:%+#v",
		// 	statusCode[status], elementAddr, modelId, ids)
	})
}

func ConfigVendorModelSubscriptionGet(elementAddress, modelIdentifier uint) error {
	params := &ConfigVendorModelSubscriptionGetMessageParameters{
		ElementAddress:  elementAddress,
		ModelIdentifier: modelIdentifier,
	}
	node, _ := findNodeByAddr(elementAddress)

	ele, err := findElementByAddr(elementAddress)
	if err != nil {
		return err
	}
	m, err := ele.findModel(modelIdentifier)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigVendorModelSubscriptionGet, params, func(n *Node, d interface{}) error {
		// virtual address is not considered
		resp := d.(ConfigVendorModelSubscriptionListMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			sigAddrs := []uint{}
			for _, addr := range m.SubAddresses {
				if addr <= 0xFFFF {
					sigAddrs = append(sigAddrs, addr)
				}
			}
			m.SubAddresses = append(resp.Addresses, sigAddrs...)

			return nil
		}
		return errors.InvalidResponse.New()
		// var status, elementAddr, modelId uint
		// utils.UnpackLE(msg.payload, "8,16,32", &status, &elementAddr, &modelId)
		// p := msg.payload[7:]
		// ids := make([]uint16, uint(len(msg.payload)-7)/2)
		// for i := 0; i < len(ids); i++ {
		// 	ids[i] = binary.LittleEndian.Uint16(p[i : i+2])
		// }
		// loggerFoundation.Infof("subscriptions: status: %s, element address: %4x, model: %4x, id:%+#v",
		// 	statusCode[status], elementAddr, modelId, ids)
	})
}

func ConfigNetKeyAdd(dst uint, netkeyIndex uint) error {
	node, _ := findNodeByAddr(dst)
	if k, _ := node.findNodeNetKeyByIndex(netkeyIndex); k != nil {
		return errors.NetKeyAlreadyBindedToNode.New().AddContextF("node:%4x, netkeyIndex:%d", dst, netkeyIndex)
	}
	key, err := findNetKeyByIndex(netkeyIndex)
	if err != nil {
		return err
	}
	params := &ConfigNetKeyAddMessageParameters{
		NetKeyIndex: netkeyIndex,
		NetKey:      key.Bytes,
	}
	return modelSendTmplParsed(false, dst, opConfigNetKeyAdd, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigNetKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			node.BindedKeys = append(node.BindedKeys, NodeKeyBinding{NetKeyIndex: netkeyIndex})
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigNetKeyUpdate(dst uint, netkeyIndex uint) error {
	newKey := createNetKey(netkeyIndex)
	oldkey, err := findNetKeyByIndex(netkeyIndex)
	if err != nil {
		return err
	}
	params := &ConfigNetKeyUpdateMessageParameters{
		NetKeyIndex: netkeyIndex,
		NetKey:      newKey.Bytes,
	}
	return modelSendTmplParsed(false, dst, opConfigNetKeyUpdate, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigNetKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			copier.Copy(oldkey, newKey)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigNetKeyDelete(dst uint, netkeyIndex uint) error {
	node, _ := findNodeByAddr(dst)
	slIdx := -1
	for i, k := range node.BindedKeys {
		if k.NetKeyIndex == netkeyIndex {
			slIdx = i
			break
		}
	}
	_, err := findNetKeyByIndex(netkeyIndex)
	if err != nil {
		return err
	}
	params := &ConfigNetKeyDeleteMessageParameters{
		NetKeyIndex: netkeyIndex,
	}
	return modelSendTmplParsed(false, dst, opConfigNetKeyDelete, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigNetKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			node.BindedKeys = append(node.BindedKeys[:slIdx],
				node.BindedKeys[slIdx+1:]...)
			return nil
		}
		return errors.InvalidResponse.New()
		//todo: remove appkey bindings
	})
}

func ConfigNetKeyGet(dst uint) error {
	node, _ := findNodeByAddr(dst)
	return modelSendTmplParsed(false, dst, opConfigNetKeyGet, nil, func(n *Node, d interface{}) error {
		resp := d.(ConfigNetKeyListMessageParameters)
		newBinding := []NodeKeyBinding{}
		for _, keyId := range resp.NetKeyIndexes {
			for _, b := range node.BindedKeys {
				if b.NetKeyIndex == keyId {
					newBinding = append(newBinding, b)
					break
				} else {
					newBinding = append(newBinding, NodeKeyBinding{NetKeyIndex: keyId})
				}
			}
		}
		node.BindedKeys = newBinding
		return nil
	})
}

func ConfigAppKeyAdd(dst, netKeyIndex, appKeyIndex uint) error {
	node, _ := findNodeByAddr(dst)
	if key, _ := node.findNodeAppKeyByIndex(appKeyIndex); key != nil {
		return errors.AppKeyAlreadyBindedToNode.New().AddContextF("node:%4x, appkeyIndex:%d", dst, appKeyIndex)
	}
	_, err := node.findNodeNetKeyByIndex(netKeyIndex)
	if err != nil {
		return err
	}
	appkey, err := findAppKeyByIndex(appKeyIndex)
	if err != nil {
		return err
	}
	params := &ConfigAppKeyAddMessageParameters{
		NetKeyIndex: netKeyIndex,
		AppKeyIndex: appKeyIndex,
		AppKey:      appkey.Bytes,
	}
	return modelSendTmplParsed(false, dst, opConfigAppKeyAdd, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigAppKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			for i := range n.BindedKeys {
				if n.BindedKeys[i].NetKeyIndex == netKeyIndex {
					if !funk.Contains(n.BindedKeys[i].BindedAppKeyIds, appKeyIndex) {
						n.BindedKeys[i].BindedAppKeyIds = append(n.BindedKeys[i].BindedAppKeyIds, appKeyIndex)
					}
				}
			}
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigAppKeyUpdate(dst, netKeyIndex, appKeyIndex uint) error {
	node, _ := findNodeByAddr(dst)
	_, err := node.findNodeAppKeyByIndex(appKeyIndex)
	if err != nil {
		return err
	}
	newKey := createAppKey(appKeyIndex)
	params := &ConfigAppKeyUpdateMessageParameters{
		NetKeyIndex: netKeyIndex,
		AppKeyIndex: appKeyIndex,
		AppKey:      newKey.Bytes,
	}
	return modelSendTmplParsed(false, dst, opConfigAppKeyUpdate, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigAppKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			loggerFoundation.Infoln("app key updated successfully")
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigAppKeyDelete(dst, netKeyIndex, appKeyIndex uint) error {
	var binding *NodeKeyBinding
	var j int
	node, _ := findNodeByAddr(dst)
	for ni, k := range node.BindedKeys {
		if k.NetKeyIndex == netKeyIndex && k.BindedAppKeyIds != nil {
			for ai, ak := range k.BindedAppKeyIds {
				if ak == appKeyIndex {
					binding = &node.BindedKeys[ni]
					j = ai
				}
			}
		}
	}
	if binding == nil {
		return errors.AppKeyNotBindedToNode.New().AddContextF("appkeyIndex:%d, node:%4x", appKeyIndex, dst)
	}
	params := &ConfigAppKeyDeleteMessageParameters{
		NetKeyIndex: netKeyIndex,
		AppKeyIndex: appKeyIndex,
	}

	return modelSendTmplParsed(false, dst, opConfigAppKeyDelete, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigAppKeyStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			binding.BindedAppKeyIds = append(binding.BindedAppKeyIds[:j],
				binding.BindedAppKeyIds[j+1:]...)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigAppKeyGet(dst, netKeyIndex uint) error {
	params := &ConfigAppKeyGetMessageParameters{
		NetKeyIndex: netKeyIndex,
	}
	node, _ := findNodeByAddr(dst)
	_, err := node.findNodeNetKeyByIndex(netKeyIndex)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opConfigAppKeyGet, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigAppKeyListMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			for i := range node.BindedKeys {
				if node.BindedKeys[i].NetKeyIndex == netKeyIndex {
					node.BindedKeys[i].BindedAppKeyIds = append([]uint{}, resp.AppKeyIndexes...)
					break
				}
			}
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func handleNodeIdentityResponse(n *Node, d interface{}) error {
	resp := d.(ConfigNodeIdentityStatusMessageParameters)
	if resp.Status == STATUS_SUCCESS {
		n.NodeIdentityStates[resp.NetKeyIndex] = resp.Identity
		return nil
	}
	return errors.InvalidResponse.New()
}

func ConfigNodeIdentityGet(dst, netKeyIndex uint) error {
	params := &ConfigNodeIdentityGetMessageParameters{
		NetKeyIndex: netKeyIndex,
	}
	return modelSendTmplParsed(false, dst, opConfigNodeIdentityGet, params, handleNodeIdentityResponse)
}

func ConfigNodeIdentitySet(dst, netKeyIndex, identity uint) error {
	params := &ConfigNodeIdentitySetMessageParameters{
		NetKeyIndex: netKeyIndex,
		Identity:    identity,
	}
	return modelSendTmplParsed(false, dst, opConfigNodeIdentitySet, params, handleNodeIdentityResponse)
}

func ConfigModelAppBind(elementAddr, appKeyIndex, modelId uint) error {
	node, err := findNodeByAddr(elementAddr)
	if err != nil {
		return err
	}
	e, err := findElementByAddr(elementAddr)
	if err != nil {
		return err
	}
	m, err := e.findModel(modelId)
	if err != nil {
		return err
	}
	if k, _ := m.findBindedAppKey(appKeyIndex); k != nil {
		return errors.AppKeyAlreadyBindedToModel.New()
	}
	_, err = node.findNodeAppKeyByIndex(appKeyIndex)
	if err != nil {
		return err
	}

	params := &ConfigModelAppBindMessageParameters{
		ElementAddress:  elementAddr,
		AppKeyIndex:     appKeyIndex,
		ModelIdentifier: modelId,
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelAppBind, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelAppStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			m.BindedAppKeyIds = append(m.BindedAppKeyIds, resp.AppKeyIndex)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func ConfigModelAppUnbind(elementAddr, appKeyIndex, modelId uint) error {
	var index uint
	var model *Model
	node, err := findNodeByAddr(elementAddr)
	if err != nil {
		return err
	}
	e, err := findElementByAddr(elementAddr)
	if err != nil {
		return err
	}
	m, err := e.findModel(modelId)
	if err != nil {
		return err
	}
	for i, k := range m.BindedAppKeyIds {
		if k == appKeyIndex {
			model = m
			index = uint(i)
		}
	}

	if model == nil {
		return errors.ModelAppKeyBindingNotFound.New().AddContextF("appkeyIndex:%d, model:%4x", appKeyIndex, modelId)
	}
	params := &ConfigModelAppUnbindMessageParameters{
		ElementAddress:  elementAddr,
		AppKeyIndex:     appKeyIndex,
		ModelIdentifier: modelId,
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opConfigModelAppUnbind, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigModelAppStatusMessageParameters)
		if resp.Status == STATUS_SUCCESS {
			model.BindedAppKeyIds = append(model.BindedAppKeyIds[:index], model.BindedAppKeyIds[index+1:]...)
			// writeNodeToDb(node)
			return nil
		}
		return errors.InvalidResponse.New()
	})
}

func configModelAppGet(elementAddr, modelId uint) error {
	var model *Model
	node, err := findNodeByAddr(elementAddr)
	if err != nil {
		return err
	}
	e, err := findElementByAddr(elementAddr)
	if err != nil {
		return err
	}
	model, err = e.findModel(modelId)
	if err != nil {
		return err
	}

	var params interface{}
	params = &ConfigSIGModelAppGetMessageParameters{
		ElementAddress:  elementAddr,
		ModelIdentifier: modelId,
	}
	opcode := uint(opConfigSIGModelAppGet)
	if !isSigModel(modelId) {
		opcode = opConfigVendorModelAppGet
		params = &ConfigVendorModelAppGetMessageParameters{
			ElementAddress:  elementAddr,
			ModelIdentifier: modelId,
		}
	}
	return modelSendTmplParsed(false, node.UnicastAddress, opcode, params, func(n *Node, d interface{}) error {
		sigResp, ok := d.(ConfigSIGModelAppListMessageParameters)
		vndResp, ok1 := d.(ConfigVendorModelAppListMessageParameters)
		if ok && sigResp.Status == STATUS_SUCCESS {
			model.BindedAppKeyIds = append([]uint{}, sigResp.AppKeyIndexes...)
		} else if ok1 && vndResp.Status == STATUS_SUCCESS {
			model.BindedAppKeyIds = append([]uint{}, vndResp.AppKeyIndexes...)
		}
		return nil
	})
}

func ConfigSigModelAppGet(elementAddr, modelId uint) error {
	return configModelAppGet(elementAddr, modelId)
}

func ConfigVendorModelAppGet(elementAddr, modelId uint) error {
	return configModelAppGet(elementAddr, modelId)
}

func ConfigNodeReset(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigNodeReset, nil, func(n *Node, d interface{}) error {
		deleteNode(dst)
		loggerFoundation.Info("node reset finished")
		return nil
	})
}

func handleFriendResponse(n *Node, d interface{}) error {
	resp := d.(ConfigFriendStatusMessageParameters)
	n.FriendState = resp.Friend
	return nil
}

func ConfigFriendGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigFriendGet, nil, handleFriendResponse)
}

func ConfigFriendSet(dst, friendState uint) error {
	params := &ConfigFriendSetMessageParameters{
		Friend: friendState,
	}
	return modelSendTmplParsed(false, dst, opConfigFriendSet, params, handleFriendResponse)
}

// todo: key refresh, heartbeat

func ConfigLowPowerNodePollTimeoutGet(friendAddr, lpnAddr uint) error {
	params := &ConfigLowPowerNodePollTimeoutGetMessageParameters{
		LPNAddress: lpnAddr,
	}
	return modelSendTmplParsed(false, friendAddr, opConfigLowPowerNodePollTimeoutGet, params, func(n *Node, d interface{}) error {
		resp := d.(ConfigLowPowerNodePollTimeoutStatusMessageParameters)
		n.PollTimeoutListState[resp.LPNAddress] = resp.PollTimeout
		return nil
	})
}

func handleNetworkTransmitResponse(n *Node, d interface{}) error {
	resp := d.(ConfigNetworkTransmitStatusMessageParameters)
	n.NetwrokTransmitState = resp
	return nil
}

func ConfigNetworkTransmitGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opConfigNetworkTransmitGet, nil, handleNetworkTransmitResponse)
}

func ConfigNetworkTransmitSet(dst uint, count, step uint) error {
	params := &ConfigNetworkTransmitSetMessageParameters{
		NetworkTransmitCount:         count,
		NetworkTransmitIntervalSteps: step,
	}
	return modelSendTmplParsed(false, dst, opConfigNetworkTransmitSet, params, handleNetworkTransmitResponse)
}

// todo: health messages
