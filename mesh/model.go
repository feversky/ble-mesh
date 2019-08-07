package mesh

import (
	"ble-mesh/mesh/def"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
	"bytes"
	"encoding/binary"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	funk "github.com/thoas/go-funk"
)

var statusCode = []string{
	"Success",
	"InvalidAddress",
	"InvalidModel",
	"InvalidAppKeyIndex",
	"InvalidNetKeyIndex",
	"InsufficientResources",
	"KeyIndexAlreadyStored",
	"InvalidPublishParameters",
	"NotaSubscribeModel",
	"StorageFailure",
	"FeatureNotSupported",
	"CannotUpdate",
	"CannotRemove",
	"CannotBind",
	"TemporarilyUnabletoChangeState",
	"CannotSet",
	"UnspecifiedError",
	"InvalidBinding",
}

var transactionIdentifier uint

var transitionStepResolution = map[uint]float32{
	0: 0.1,
	1: 1,
	2: 10,
	3: 10 * 60,
}

const (
	NODLC = 0xFFFFFFFF

	ConfigServer                       = 0x0000
	HealthServer                       = 0x0002
	GenericOnOffServer                 = 0x1000
	GenericOnOffClient                 = 0x1001
	GenericLevelServer                 = 0x1002
	GenericLevelClient                 = 0x1003
	GenericDefaultTransitionTimeServer = 0x1004
	GenericDefaultTransitionTimeClient = 0x1005
	GenericPowerOnOffServer            = 0x1006
	GenericPowerOnOffSetupServer       = 0x1007
	GenericPowerOnOffClient            = 0x1008
	GenericPowerLevelServer            = 0x1009
	GenericPowerLevelSetupServer       = 0x100A
	GenericPowerLevelClient            = 0x100B
	GenericBatteryServer               = 0x100C
	GenericBatteryClient               = 0x100D
	GenericLocationServer              = 0x100E
	GenericLocationSetupServer         = 0x100F
	GenericLocationClient              = 0x1010
	GenericAdminPropertyServer         = 0x1011
	GenericManufacturerPropertyServer  = 0x1012
	GenericUserPropertyServer          = 0x1013
	GenericClientPropertyServer        = 0x1014
	GenericPropertyClient              = 0x1015
	SensorServer                       = 0x1100
	SensorSetupServer                  = 0x1101
	SensorClient                       = 0x1102
	TimeServer                         = 0x1200
	TimeSetupServer                    = 0x1201
	TimeClient                         = 0x1202
	SceneServer                        = 0x1203
	SceneSetupServer                   = 0x1204
	SceneClient                        = 0x1205
	SchedulerServer                    = 0x1206
	SchedulerSetupServer               = 0x1207
	SchedulerClient                    = 0x1208
	LightLightnessServer               = 0x1300
	LightLightnessSetupServer          = 0x1301
	LightLightnessClient               = 0x1302
	LightCTLServer                     = 0x1303
	LightCTLSetupServer                = 0x1304
	LightCTLClient                     = 0x1305
	LightCTLTemperatureServer          = 0x1306
	LightHSLServer                     = 0x1307
	LightHSLSetupServer                = 0x1308
	LightHSLClient                     = 0x1309
	LightHSLHueServer                  = 0x130A
	LightHSLSaturationServer           = 0x130B
	LightxyLServer                     = 0x130C
	LightxyLSetupServer                = 0x130D
	LightxyLClient                     = 0x130E
	LightLCServer                      = 0x130F
	LightLCSetupServer                 = 0x1310
	LightLCClient                      = 0x1311
)

type (
	Transition struct {
		TransitionTime uint `bits:"8"`
		Delay          uint `bits:"8"`
	}

	TID struct {
		TID uint `bits:"8"`
	}

	RemainingTime struct {
		TransitionNumberOfSteps  uint `bits:"2"`
		TransitionStepResolution uint `bits:"6"`
	}

	AccessMessage struct {
		src     uint
		dst     uint
		opcode  uint
		payload []byte
	}

	AccessMessageTx struct {
		AccessMessage
		ch                 chan *AccessMessage
		expectedRespOpcode uint
		oboAck             bool
	}

	onResponseReceived func(*AccessMessage) error
	modelMsglistener   func(*AccessMessage)
)

func (r RemainingTime) parse() float32 {
	return transitionStepResolution[r.TransitionStepResolution] * float32(r.TransitionNumberOfSteps)
}

func (m *AccessMessage) Parse(out interface{}) error {
	err := utils.UnpackStructLE(m.payload, out)
	if err != nil {
		loggerModel.Error(err)
	}
	return err
}

var (
	txAccessMsgs           = map[uint]*AccessMessageTx{}
	logger                 *logrus.Entry
	opcodeConfigReqRespMap = map[uint]uint{
		opConfigBeaconGet:                                opConfigBeaconStatus,
		opConfigBeaconSet:                                opConfigBeaconStatus,
		opConfigCompositionDataGet:                       opConfigCompositionDataStatus,
		opConfigDefaultTTLGet:                            opConfigDefaultTTLStatus,
		opConfigDefaultTTLSet:                            opConfigDefaultTTLStatus,
		opConfigGATTProxyGet:                             opConfigGATTProxyStatus,
		opConfigGATTProxySet:                             opConfigGATTProxyStatus,
		opConfigRelayGet:                                 opConfigRelayStatus,
		opConfigRelaySet:                                 opConfigRelayStatus,
		opConfigModelPublicationGet:                      opConfigModelPublicationStatus,
		opConfigModelPublicationSet:                      opConfigModelPublicationStatus,
		opConfigModelPublicationVirtualAddressSet:        opConfigModelPublicationStatus,
		opConfigModelSubscriptionAdd:                     opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionVirtualAddressAdd:       opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionDelete:                  opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionVirtualAddressDelete:    opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionOverwrite:               opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionVirtualAddressOverwrite: opConfigModelSubscriptionStatus,
		opConfigModelSubscriptionDeleteAll:               opConfigModelSubscriptionStatus,
		opConfigSIGModelSubscriptionGet:                  opConfigSIGModelSubscriptionList,
		opConfigVendorModelSubscriptionGet:               opConfigVendorModelSubscriptionList,
		opConfigNetKeyAdd:                                opConfigNetKeyStatus,
		opConfigNetKeyDelete:                             opConfigNetKeyStatus,
		opConfigNetKeyUpdate:                             opConfigNetKeyStatus,
		opConfigNetKeyGet:                                opConfigNetKeyList,
		opConfigAppKeyAdd:                                opConfigAppKeyStatus,
		opConfigAppKeyDelete:                             opConfigAppKeyStatus,
		opConfigAppKeyUpdate:                             opConfigAppKeyStatus,
		opConfigAppKeyGet:                                opConfigAppKeyList,
		opConfigNodeIdentityGet:                          opConfigNodeIdentityStatus,
		opConfigNodeIdentitySet:                          opConfigNodeIdentityStatus,
		opConfigModelAppBind:                             opConfigModelAppStatus,
		opConfigModelAppUnbind:                           opConfigModelAppStatus,
		opConfigSIGModelAppGet:                           opConfigSIGModelAppList,
		opConfigVendorModelAppGet:                        opConfigVendorModelAppList,
		opConfigNodeReset:                                opConfigNodeResetStatus,
		opConfigFriendGet:                                opConfigFriendStatus,
		opConfigFriendSet:                                opConfigFriendStatus,
		opConfigKeyRefreshPhaseGet:                       opConfigKeyRefreshPhaseStatus,
		opConfigKeyRefreshPhaseSet:                       opConfigKeyRefreshPhaseStatus,
		opConfigHeartbeatPublicationGet:                  opConfigHeartbeatPublicationStatus,
		opConfigHeartbeatPublicationSet:                  opConfigHeartbeatPublicationStatus,
		opConfigHeartbeatSubscriptionGet:                 opConfigHeartbeatSubscriptionStatus,
		opConfigHeartbeatSubscriptionSet:                 opConfigHeartbeatSubscriptionStatus,
		opConfigLowPowerNodePollTimeoutGet:               opConfigLowPowerNodePollTimeoutStatus,
		opConfigNetworkTransmitGet:                       opConfigNetworkTransmitStatus,
		opConfigNetworkTransmitSet:                       opConfigNetworkTransmitStatus,
	}
	opcodeReqRespMap = map[uint]uint{
		opHealthAttentionGet: opHealthAttentionStatus,
		opHealthAttentionSet: opHealthAttentionStatus,
		opHealthFaultClear:   opHealthFaultStatus,
		opHealthFaultGet:     opHealthFaultStatus,
		opHealthFaultTest:    opHealthFaultStatus,
		opHealthPeriodGet:    opHealthPeriodStatus,
		opHealthPeriodSet:    opHealthPeriodStatus,

		opGenericOnOffGet:                  opGenericOnOffStatus,
		opGenericOnOffSet:                  opGenericOnOffStatus,
		opGenericLevelGet:                  opGenericLevelStatus,
		opGenericLevelSet:                  opGenericLevelStatus,
		opGenericDeltaSet:                  opGenericLevelStatus,
		opGenericMoveSet:                   opGenericLevelStatus,
		opGenericDefaultTransitionTimeGet:  opGenericDefaultTransitionTimeStatus,
		opGenericDefaultTransitionTimeSet:  opGenericDefaultTransitionTimeStatus,
		opGenericOnPowerUpGet:              opGenericOnPowerUpStatus,
		opGenericOnPowerUpSet:              opGenericOnPowerUpStatus,
		opGenericPowerLevelGet:             opGenericPowerLevelStatus,
		opGenericPowerLevelSet:             opGenericPowerLevelStatus,
		opGenericPowerLastGet:              opGenericPowerLastStatus,
		opGenericPowerDefaultGet:           opGenericPowerDefaultStatus,
		opGenericPowerRangeGet:             opGenericPowerRangeStatus,
		opGenericBatteryGet:                opGenericBatteryStatus,
		opGenericLocationGlobalGet:         opGenericLocationGlobalStatus,
		opGenericLocationLocalGet:          opGenericLocationLocalStatus,
		opGenericLocationGlobalSet:         opGenericLocationGlobalStatus,
		opGenericLocationLocalSet:          opGenericLocationLocalStatus,
		opGenericManufacturerPropertiesGet: opGenericManufacturerPropertiesStatus,
		opGenericManufacturerPropertyGet:   opGenericManufacturerPropertyStatus,
		opGenericManufacturerPropertySet:   opGenericManufacturerPropertyStatus,
		opGenericAdminPropertiesGet:        opGenericAdminPropertiesStatus,
		opGenericAdminPropertyGet:          opGenericAdminPropertyStatus,
		opGenericAdminPropertySet:          opGenericAdminPropertyStatus,
		opGenericUserPropertiesGet:         opGenericUserPropertiesStatus,
		opGenericUserPropertyGet:           opGenericUserPropertyStatus,
		opGenericUserPropertySet:           opGenericUserPropertyStatus,
		opGenericClientPropertiesGet:       opGenericClientPropertiesStatus,

		opSensorDescriptorGet: opSensorDescriptorStatus,
		opSensorGet:           opSensorStatus,
		opSensorColumnGet:     opSensorColumnStatus,
		opSensorSeriesGet:     opSensorSeriesStatus,
		opSensorCadenceGet:    opSensorCadenceStatus,
		opSensorCadenceSet:    opSensorCadenceStatus,
		opSensorSettingsGet:   opSensorSettingsStatus,
		opSensorSettingGet:    opSensorSettingStatus,
		opSensorSettingSet:    opSensorSettingStatus,

		opTimeGet:            opTimeStatus,
		opTimeSet:            opTimeStatus,
		opTimeRoleGet:        opTimeRoleStatus,
		opTimeRoleSet:        opTimeRoleStatus,
		opTimeZoneGet:        opTimeZoneStatus,
		opTimeZoneSet:        opTimeZoneStatus,
		opTAI_UTCDeltaGet:    opTAI_UTCDeltaStatus,
		opTAI_UTCDeltaSet:    opTAI_UTCDeltaStatus,
		opSceneGet:           opSceneStatus,
		opSceneRecall:        opSceneStatus,
		opSceneRegisterGet:   opSceneRegisterStatus,
		opSceneStore:         opSceneRegisterStatus,
		opSceneDelete:        opSceneRegisterStatus,
		opSchedulerActionGet: opSchedulerActionStatus,
		opSchedulerGet:       opSchedulerStatus,
		opSchedulerActionSet: opSchedulerActionStatus,

		opLightLightnessGet:        opLightLightnessStatus,
		opLightLightnessSet:        opLightLightnessStatus,
		opLightLightnessLinearGet:  opLightLightnessLinearStatus,
		opLightLightnessLinearSet:  opLightLightnessLinearStatus,
		opLightLightnessLastGet:    opLightLightnessLastStatus,
		opLightLightnessDefaultGet: opLightLightnessDefaultStatus,
		opLightLightnessRangeGet:   opLightLightnessRangeStatus,

		opLightLightnessDefaultSet: opLightLightnessDefaultStatus,
		opLightLightnessRangeSet:   opLightLightnessRangeStatus,

		opLightCTLGet:                 opLightCTLStatus,
		opLightCTLSet:                 opLightCTLStatus,
		opLightCTLTemperatureGet:      opLightCTLTemperatureStatus,
		opLightCTLTemperatureRangeGet: opLightCTLTemperatureRangeStatus,
		opLightCTLTemperatureSet:      opLightCTLTemperatureStatus,
		opLightCTLDefaultGet:          opLightCTLDefaultStatus,

		opLightCTLDefaultSet:          opLightCTLDefaultStatus,
		opLightCTLTemperatureRangeSet: opLightCTLTemperatureRangeStatus,

		opLightHSLGet:           opLightHSLStatus,
		opLightHSLHueGet:        opLightHSLHueStatus,
		opLightHSLHueSet:        opLightHSLHueStatus,
		opLightHSLSaturationGet: opLightHSLSaturationStatus,
		opLightHSLSaturationSet: opLightHSLSaturationStatus,
		opLightHSLSet:           opLightHSLStatus,
		opLightHSLTargetGet:     opLightHSLTargetStatus,
		opLightHSLDefaultGet:    opLightHSLDefaultStatus,
		opLightHSLRangeGet:      opLightHSLRangeStatus,

		opLightHSLDefaultSet: opLightHSLDefaultStatus,
		opLightHSLRangeSet:   opLightHSLRangeStatus,

		opLightxyLGet:        opLightxyLStatus,
		opLightxyLSet:        opLightxyLStatus,
		opLightxyLTargetGet:  opLightxyLTargetStatus,
		opLightxyLDefaultGet: opLightxyLDefaultStatus,
		opLightxyLRangeGet:   opLightxyLRangeStatus,

		opLightxyLDefaultSet: opLightxyLDefaultStatus,
		opLightxyLRangeSet:   opLightxyLRangeStatus,

		opLightLCModeGet:       opLightLCModeStatus,
		opLightLCModeSet:       opLightLCModeStatus,
		opLightLCOMGet:         opLightLCOMStatus,
		opLightLCOMSet:         opLightLCOMStatus,
		opLightLCLightOnOffGet: opLightLCLightOnOffStatus,
		opLightLCLightOnOffSet: opLightLCLightOnOffStatus,
		opLightLCPropertyGet:   opLightLCPropertyStatus,
		opLightLCPropertySet:   opLightLCPropertyStatus,
	}

	modelMap = map[uint][]int{
		HealthServer:                       []int{opHealthAttentionGet, opHealthAttentionSet, opHealthFaultClear, opHealthFaultGet, opHealthFaultTest, opHealthPeriodGet, opHealthPeriodSet},
		GenericOnOffServer:                 []int{opGenericOnOffGet, opGenericOnOffSet, opGenericOnOffSetUnacknowledged},
		GenericLevelServer:                 []int{opGenericLevelGet, opGenericLevelSet, opGenericLevelSetUnacknowledged, opGenericDeltaSet, opGenericDeltaSetUnacknowledged, opGenericMoveSet, opGenericMoveSetUnacknowledged},
		GenericDefaultTransitionTimeServer: []int{opGenericDefaultTransitionTimeGet, opGenericDefaultTransitionTimeSet, opGenericDefaultTransitionTimeSetUnacknowledged},
		GenericPowerOnOffServer:            []int{opGenericOnPowerUpGet},
		GenericPowerOnOffSetupServer:       []int{opGenericOnPowerUpSet, opGenericOnPowerUpSetUnacknowledged},
		GenericPowerLevelServer:            []int{opGenericPowerLevelGet, opGenericPowerLevelSet, opGenericPowerLevelSetUnacknowledged, opGenericPowerLastGet, opGenericPowerDefaultGet, opGenericPowerRangeGet},
		GenericPowerLevelSetupServer:       []int{opGenericPowerDefaultSet, opGenericPowerDefaultSetUnacknowledged, opGenericPowerRangeSet, opGenericPowerRangeSetUnacknowledged},
		GenericBatteryServer:               []int{opGenericBatteryGet},
		GenericLocationServer:              []int{opGenericLocationGlobalGet, opGenericLocationLocalGet},
		GenericLocationSetupServer:         []int{opGenericLocationGlobalSet, opGenericLocationGlobalSetUnacknowledged, opGenericLocationLocalSet, opGenericLocationLocalSetUnacknowledged},
		GenericAdminPropertyServer:         []int{opGenericAdminPropertiesGet, opGenericAdminPropertyGet, opGenericAdminPropertySet, opGenericAdminPropertySetUnacknowledged},
		GenericManufacturerPropertyServer:  []int{opGenericManufacturerPropertiesGet, opGenericManufacturerPropertyGet, opGenericManufacturerPropertySet, opGenericManufacturerPropertySetUnacknowledged},
		GenericUserPropertyServer:          []int{opGenericUserPropertiesGet, opGenericUserPropertyGet, opGenericUserPropertySet, opGenericUserPropertySetUnacknowledged},
		GenericClientPropertyServer:        []int{opGenericClientPropertiesGet},
		SensorServer:                       []int{opSensorDescriptorGet, opSensorGet, opSensorColumnGet, opSensorSeriesGet},
		SensorSetupServer:                  []int{opSensorCadenceGet, opSensorCadenceSet, opSensorCadenceSetUnacknowledged, opSensorSettingsGet, opSensorSettingGet, opSensorSettingSet, opSensorSettingSetUnacknowledged},
		TimeServer:                         []int{opTimeGet, opTimeZoneGet, opTAI_UTCDeltaGet},
		TimeSetupServer:                    []int{opTimeSet, opTimeRoleGet, opTimeRoleSet, opTimeZoneSet, opTAI_UTCDeltaSet},
		SceneServer:                        []int{opSceneGet, opSceneRegisterGet, opSceneRecall, opSceneRecallUnacknowledged},
		SceneSetupServer:                   []int{opSceneStore, opSceneStoreUnacknowledged, opSceneDelete, opSceneDeleteUnacknowledged},
		SchedulerServer:                    []int{opSchedulerActionGet, opSchedulerGet},
		SchedulerSetupServer:               []int{opSchedulerActionSet, opSchedulerActionSetUnacknowledged},
		LightLightnessServer:               []int{opLightLightnessGet, opLightLightnessSet, opLightLightnessSetUnacknowledged, opLightLightnessLinearGet, opLightLightnessLinearSet, opLightLightnessLinearSetUnacknowledged, opLightLightnessLastGet, opLightLightnessDefaultGet, opLightLightnessRangeGet},
		LightLightnessSetupServer:          []int{opLightLightnessDefaultSet, opLightLightnessDefaultSetUnacknowledged, opLightLightnessRangeSet, opLightLightnessRangeSetUnacknowledged},
		LightCTLServer:                     []int{opLightCTLGet, opLightCTLSet, opLightCTLSetUnacknowledged, opLightCTLTemperatureRangeGet, opLightCTLDefaultGet},
		LightCTLSetupServer:                []int{opLightCTLDefaultSet, opLightCTLDefaultSetUnacknowledged, opLightCTLTemperatureRangeSet, opLightCTLTemperatureRangeSetUnacknowledged},
		LightCTLTemperatureServer:          []int{opLightCTLTemperatureGet, opLightCTLTemperatureSet, opLightCTLTemperatureSetUnacknowledged},
		LightHSLServer:                     []int{opLightHSLGet, opLightHSLSet, opLightHSLSetUnacknowledged, opLightHSLTargetGet, opLightHSLDefaultGet},
		LightHSLSetupServer:                []int{opLightHSLDefaultSet, opLightHSLDefaultSetUnacknowledged, opLightHSLRangeSet, opLightHSLRangeSetUnacknowledged},
		LightHSLHueServer:                  []int{opLightHSLHueGet, opLightHSLHueSet, opLightHSLHueSetUnacknowledged},
		LightHSLSaturationServer:           []int{opLightHSLSaturationGet, opLightHSLSaturationSet, opLightHSLSaturationSetUnacknowledged},
		LightxyLServer:                     []int{opLightxyLGet, opLightxyLSet, opLightxyLSetUnacknowledged, opLightxyLTargetGet, opLightxyLDefaultGet},
		LightxyLSetupServer:                []int{opLightxyLDefaultSet, opLightxyLDefaultSetUnacknowledged, opLightxyLRangeSet, opLightxyLRangeSetUnacknowledged},
		LightLCServer:                      []int{opLightLCModeGet, opLightLCModeSet, opLightLCModeSetUnacknowledged, opLightLCOMGet, opLightLCOMSet, opLightLCOMSetUnacknowledged, opLightLCLightOnOffGet, opLightLCLightOnOffSet, opLightLCLightOnOffSetUnacknowledged, opSensorStatus},
		LightLCSetupServer:                 []int{opLightLCPropertyGet, opLightLCPropertySet, opLightLCPropertySetUnacknowledged},
	}

	expectedRespLen = map[uint]func(d []byte) bool{
		opConfigBeaconStatus:                  func(d []byte) bool { return len(d) == 1 },
		opConfigCompositionDataStatus:         func(d []byte) bool { return len(d) > 14 },
		opConfigDefaultTTLStatus:              func(d []byte) bool { return len(d) == 1 },
		opConfigGATTProxyStatus:               func(d []byte) bool { return len(d) == 1 },
		opConfigRelayStatus:                   func(d []byte) bool { return len(d) == 2 },
		opConfigModelPublicationStatus:        func(d []byte) bool { return len(d) == 12 || len(d) == 14 },
		opConfigModelSubscriptionStatus:       func(d []byte) bool { return len(d) == 7 || len(d) == 9 },
		opConfigSIGModelSubscriptionList:      func(d []byte) bool { return len(d) >= 5 && (len(d)-5)%2 == 0 },
		opConfigVendorModelSubscriptionList:   func(d []byte) bool { return len(d) >= 7 && (len(d)-7)%2 == 0 },
		opConfigNetKeyStatus:                  func(d []byte) bool { return len(d) == 3 },
		opConfigNetKeyList:                    func(d []byte) bool { return len(d)*8%12 == 0 || len(d)*8%12 == 4 },
		opConfigAppKeyStatus:                  func(d []byte) bool { return len(d) == 4 },
		opConfigAppKeyList:                    func(d []byte) bool { return len(d) >= 3 && ((len(d)-3)*8%12 == 0 || (len(d)-3)*8%12 == 4) },
		opConfigNodeIdentityStatus:            func(d []byte) bool { return len(d) == 4 },
		opConfigModelAppStatus:                func(d []byte) bool { return len(d) == 7 || len(d) == 9 },
		opConfigSIGModelAppList:               func(d []byte) bool { return len(d) >= 5 && ((len(d)-5)*8%12 == 0 || (len(d)-5)*8%12 == 4) },
		opConfigVendorModelAppList:            func(d []byte) bool { return len(d) >= 7 && ((len(d)-7)*8%12 == 0 || (len(d)-7)*8%12 == 4) },
		opConfigNodeResetStatus:               func(d []byte) bool { return len(d) == 0 },
		opConfigFriendStatus:                  func(d []byte) bool { return len(d) == 1 },
		opConfigKeyRefreshPhaseStatus:         func(d []byte) bool { return len(d) == 4 },
		opConfigHeartbeatPublicationStatus:    func(d []byte) bool { return len(d) == 10 },
		opConfigHeartbeatSubscriptionStatus:   func(d []byte) bool { return len(d) == 9 },
		opConfigLowPowerNodePollTimeoutStatus: func(d []byte) bool { return len(d) == 5 },
		opConfigNetworkTransmitStatus:         func(d []byte) bool { return len(d) == 1 },

		//generic models
		opGenericOnOffStatus:                  func(d []byte) bool { return len(d) == 1 || len(d) == 3 },
		opGenericLevelStatus:                  func(d []byte) bool { return len(d) == 2 || len(d) == 5 },
		opGenericDefaultTransitionTimeStatus:  func(d []byte) bool { return len(d) == 1 },
		opGenericOnPowerUpStatus:              func(d []byte) bool { return len(d) == 1 },
		opGenericPowerLevelStatus:             func(d []byte) bool { return len(d) == 2 || len(d) == 5 },
		opGenericPowerLastStatus:              func(d []byte) bool { return len(d) == 2 },
		opGenericPowerDefaultStatus:           func(d []byte) bool { return len(d) == 2 },
		opGenericPowerRangeStatus:             func(d []byte) bool { return len(d) == 5 },
		opGenericBatteryStatus:                func(d []byte) bool { return len(d) == 8 },
		opGenericLocationGlobalStatus:         func(d []byte) bool { return len(d) == 10 },
		opGenericLocationLocalStatus:          func(d []byte) bool { return len(d) == 11 },
		opGenericUserPropertiesStatus:         func(d []byte) bool { return len(d)%2 == 0 },
		opGenericUserPropertyStatus:           func(d []byte) bool { return len(d) == 2 || (len(d) > 3) },
		opGenericAdminPropertiesStatus:        func(d []byte) bool { return len(d)%2 == 0 },
		opGenericAdminPropertyStatus:          func(d []byte) bool { return len(d) == 2 || len(d) > 3 },
		opGenericManufacturerPropertiesStatus: func(d []byte) bool { return len(d) == 2 },
		opGenericManufacturerPropertyStatus:   func(d []byte) bool { return len(d) == 2 || len(d) > 3 },
		opGenericClientPropertiesStatus:       func(d []byte) bool { return len(d)%2 == 0 },

		//light models
		opLightLightnessStatus:           func(d []byte) bool { return len(d) == 2 || len(d) == 5 },
		opLightLightnessLinearStatus:     func(d []byte) bool { return len(d) == 2 || len(d) == 5 },
		opLightLightnessLastStatus:       func(d []byte) bool { return len(d) == 2 },
		opLightLightnessDefaultStatus:    func(d []byte) bool { return len(d) == 2 },
		opLightLightnessRangeStatus:      func(d []byte) bool { return len(d) == 5 },
		opLightCTLStatus:                 func(d []byte) bool { return len(d) == 4 || len(d) == 9 },
		opLightCTLTemperatureStatus:      func(d []byte) bool { return len(d) == 4 || len(d) == 9 },
		opLightCTLTemperatureRangeStatus: func(d []byte) bool { return len(d) == 5 },
		opLightCTLDefaultStatus:          func(d []byte) bool { return len(d) == 6 },
		opLightHSLStatus:                 func(d []byte) bool { return len(d) == 7 },
		opLightHSLTargetStatus:           func(d []byte) bool { return len(d) == 7 },
		opLightHSLHueStatus:              func(d []byte) bool { return len(d) == 5 },
		opLightHSLSaturationStatus:       func(d []byte) bool { return len(d) == 5 },
		opLightHSLDefaultStatus:          func(d []byte) bool { return len(d) == 6 },
		opLightHSLRangeStatus:            func(d []byte) bool { return len(d) == 9 },
		opLightxyLStatus:                 func(d []byte) bool { return len(d) == 7 },
		opLightxyLTargetStatus:           func(d []byte) bool { return len(d) == 7 },
		opLightxyLDefaultStatus:          func(d []byte) bool { return len(d) == 6 },
		opLightxyLRangeStatus:            func(d []byte) bool { return len(d) == 9 },
		opLightLCModeStatus:              func(d []byte) bool { return len(d) == 1 },
		opLightLCOMStatus:                func(d []byte) bool { return len(d) == 1 },
		opLightLCLightOnOffStatus:        func(d []byte) bool { return len(d) == 3 },
		opLightLCPropertyStatus:          func(d []byte) bool { return len(d) > 2 },
	}

	respUnmarshallMap = map[uint]reflect.Type{
		opConfigBeaconStatus:                  reflect.TypeOf(def.ConfigBeaconStatusMessageParameters{}),
		opConfigCompositionDataStatus:         nil,
		opConfigDefaultTTLStatus:              reflect.TypeOf(def.ConfigDefaultTTLStatusMessageParameters{}),
		opConfigGATTProxyStatus:               reflect.TypeOf(def.ConfigGATTProxyStatusMessageParameters{}),
		opConfigRelayStatus:                   reflect.TypeOf(def.ConfigRelayStatusMessageParameters{}),
		opConfigModelPublicationStatus:        reflect.TypeOf(def.ConfigModelPublicationStatusMessageParameters{}),
		opConfigModelSubscriptionStatus:       reflect.TypeOf(def.ConfigModelSubscriptionStatusMessageParameters{}),
		opConfigNetKeyStatus:                  reflect.TypeOf(def.ConfigNetKeyStatusMessageParameters{}),
		opConfigNetKeyList:                    reflect.TypeOf(def.ConfigNetKeyListMessageParameters{}),
		opConfigAppKeyStatus:                  reflect.TypeOf(def.ConfigAppKeyStatusMessageParameters{}),
		opConfigAppKeyList:                    reflect.TypeOf(def.ConfigAppKeyListMessageParameters{}),
		opConfigNodeIdentityStatus:            reflect.TypeOf(def.ConfigNodeIdentityStatusMessageParameters{}),
		opConfigModelAppStatus:                reflect.TypeOf(def.ConfigModelAppStatusMessageParameters{}),
		opConfigSIGModelAppList:               reflect.TypeOf(def.ConfigSIGModelAppListMessageParameters{}),
		opConfigVendorModelAppList:            reflect.TypeOf(def.ConfigVendorModelAppListMessageParameters{}),
		opConfigSIGModelSubscriptionList:      reflect.TypeOf(def.ConfigSIGModelSubscriptionListMessageParameters{}),
		opConfigVendorModelSubscriptionList:   reflect.TypeOf(def.ConfigVendorModelSubscriptionListMessageParameters{}),
		opConfigFriendStatus:                  reflect.TypeOf(def.ConfigFriendStatusMessageParameters{}),
		opConfigKeyRefreshPhaseStatus:         reflect.TypeOf(def.ConfigKeyRefreshPhaseStatusMessageParameters{}),
		opConfigHeartbeatPublicationStatus:    reflect.TypeOf(def.ConfigHeartbeatPublicationStatusMessageParameters{}),
		opConfigHeartbeatSubscriptionStatus:   reflect.TypeOf(def.ConfigHeartbeatSubscriptionStatusMessageParameters{}),
		opConfigLowPowerNodePollTimeoutStatus: reflect.TypeOf(def.ConfigLowPowerNodePollTimeoutStatusMessageParameters{}),
		opConfigNetworkTransmitStatus:         reflect.TypeOf(def.ConfigNetworkTransmitStatusMessageParameters{}),
		opConfigNodeResetStatus:               nil,
		opHealthCurrentStatus:                 reflect.TypeOf(def.HealthCurrentStatusMessageParameters{}),
		opHealthFaultStatus:                   reflect.TypeOf(def.HealthFaultStatusMessageParameters{}),
		opHealthPeriodStatus:                  reflect.TypeOf(def.HealthPeriodStatusMessageParameters{}),
		opHealthAttentionStatus:               reflect.TypeOf(def.AttentionStatusMessageParameters{}),

		opGenericOnOffStatus: reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		opGenericLevelStatus: reflect.TypeOf(def.GenericLevelStatusMessageParameters{}),

		opLightLightnessStatus:           reflect.TypeOf(def.LightLightnessStatusMessageParameters{}),
		opLightLightnessLinearStatus:     reflect.TypeOf(def.LightLightnessLinearStatusMessageParameters{}),
		opLightLightnessLastStatus:       reflect.TypeOf(def.LightLightnessLastStatusMessageParameters{}),
		opLightLightnessDefaultStatus:    reflect.TypeOf(def.LightLightnessDefaultStatusMessageParameters{}),
		opLightLightnessRangeStatus:      reflect.TypeOf(def.LightLightnessRangeStatusMessageParameters{}),
		opLightCTLStatus:                 reflect.TypeOf(def.LightCTLStatusMessageParameters{}),
		opLightCTLTemperatureStatus:      reflect.TypeOf(def.LightCTLTemperatureStatusMessageParameters{}),
		opLightCTLTemperatureRangeStatus: reflect.TypeOf(def.LightCTLTemperatureRangeStatusMessageParameters{}),
		opLightCTLDefaultStatus:          reflect.TypeOf(def.LightCTLDefaultStatusMessageParameters{}),
	}

	modelStateUnmarshallMap = map[uint]reflect.Type{
		// HealthServer:                       reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		GenericOnOffServer: reflect.TypeOf(OnOffState{}),
		// GenericLevelServer:                 reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericDefaultTransitionTimeServer: reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericPowerOnOffServer:            reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericPowerOnOffSetupServer:       reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericPowerLevelServer:            reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericPowerLevelSetupServer:       reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericBatteryServer:               reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericLocationServer:              reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericLocationSetupServer:         reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericAdminPropertyServer:         reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericManufacturerPropertyServer:  reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericUserPropertyServer:          reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// GenericClientPropertyServer:        reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SensorServer:                       reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SensorSetupServer:                  reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// TimeServer:                         reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// TimeSetupServer:                    reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SceneServer:                        reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SceneSetupServer:                   reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SchedulerServer:                    reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// SchedulerSetupServer: reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		LightLightnessServer: reflect.TypeOf(LightnessState{}),
		// LightLightnessSetupServer:          reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		LightCTLServer: reflect.TypeOf(LightCtlState{}),
		// LightCTLSetupServer:                reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		LightCTLTemperatureServer: reflect.TypeOf(LightCtlState{}),
		// LightHSLServer:                     reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightHSLSetupServer:                reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightHSLHueServer:                  reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightHSLSaturationServer:           reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightxyLServer:                     reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightxyLSetupServer:                reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightLCServer:                      reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
		// LightLCSetupServer:                 reflect.TypeOf(def.GenericOnOffStatusMessageParameters{}),
	}

	loggerModel       = utils.CreateLogger("model")
	modelMsgListeners = []*modelMsglistener{}
)

type clientModel struct {
}

func print(v reflect.Value) {
	if v.Type().Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type().Kind() == reflect.Struct {
			print(f)
			return
		}
		tag := t.Field(i).Tag.Get("vt")
		if tag != "" {
			vt := def.AllValueTables[tag]
			fv := f.Interface().(uint)
			for k, s := range vt {
				if k.Start <= fv && k.End >= fv {
					loggerModel.Infof("%s = %x, %s", t.Field(i).Name, fv, s)
					break
				}
			}
		} else {
			loggerModel.Infof("%s = %x", t.Field(i).Name, f.Interface())
		}
	}
}

func modelSendTmplParsedWithTID(dst uint, opcode uint, data interface{}, cb func(*Node, interface{}) error) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return errors.NotAPointer.New().AddContextF("input data: %+#v", data)
	}
	transactionIdentifier++
	f := reflect.Indirect(reflect.ValueOf(data)).FieldByName("TID").Addr()
	if f.IsValid() {
		reflect.Indirect(f).Set(reflect.ValueOf(transactionIdentifier))
	}
	return modelSendTmplParsed(false, dst, opcode, data, cb)
}

func modelSendTmplParsed(async bool, dst uint, opcode uint, in interface{}, cb func(*Node, interface{}) error) error {
	var payload []byte
	var err error
	respOpcode := opcodeConfigReqRespMap[opcode]
	if respOpcode != 0 {
		loggerModel.Debugf("call method %s", foundationMethods[opcode])
	}
	if in != nil {
		payload, err = utils.PackStructLE(in)
		if err != nil {
			return err
		}
	}
	cbWrapper := func(m *AccessMessage) error {
		respOpcode := opcodeReqRespMap[opcode]
		if respOpcode == 0 {
			respOpcode = opcodeConfigReqRespMap[opcode]
		}
		retType := respUnmarshallMap[respOpcode]
		node, err := findNodeByAddr(m.src)
		if err != nil {
			return err
		}
		if retType == nil {
			loggerModel.Warn("respOpcode not registered")
			return cb(node, m)
		}
		loggerModel.Debugf("unmarshal type: %s", retType)
		val := reflect.New(retType)
		err = utils.UnpackStructLE(m.payload, val.Interface())
		if err != nil {
			return err
		}
		print(val)
		return cb(node, val.Elem().Interface())

	}
	return modelSendTmpl1(async, dst, opcode, payload, 5, cbWrapper, nil, nil)
}

func modelSendTmpl(async bool, dst uint, opcode uint, payload []byte, cb onResponseReceived) error {
	return modelSendTmpl1(async, dst, opcode, payload, 5, cb, nil, nil)
}

func modelSendTmpl1(async bool, dst uint, opcode uint, payload []byte, timeout uint, cb onResponseReceived, toFunc func(), ackFunc onAckReceived) error {
	// if txAccessMsgs[dst] != nil {
	// 	loggerConfCli.Warnf("previous config message not finished")
	// 	return
	// }
	var node *Node
	var err error
	if node, err = findNodeByAddr(dst); err != nil {
		return err
	}
	msg := &AccessMessageTx{
		AccessMessage: AccessMessage{
			src:     meshDb.UnicastAddress,
			dst:     dst,
			opcode:  opcode,
			payload: payload,
		},
		ch:                 make(chan *AccessMessage),
		expectedRespOpcode: opcodeReqRespMap[opcode] + opcodeConfigReqRespMap[opcode],
	}

	waitForResp := func() error {
		loggerModel.Debugf("msg Tx: %+#v", msg.AccessMessage)
		// only allow one request per node
		txAccessMsgs[dst] = msg
		f := func() error {
			to := time.NewTimer(time.Second * time.Duration(timeout))
			select {
			case msgRx := <-txAccessMsgs[dst].ch:
				validator := expectedRespLen[msgRx.opcode]
				if validator == nil {
					loggerModel.Error("missing DLC")
				}
				if msg.expectedRespOpcode == msgRx.opcode && validator != nil && validator(msgRx.payload) {
					if cb != nil {
						err := cb(msgRx)
						if err != nil {
							return err
						}
						writeNodeToDb(node)
					}
				} else {
					return errors.DataLengthCheckFailed.New()
				}
				delete(txAccessMsgs, dst)
				return nil
			case <-to.C:
				if toFunc != nil {
					toFunc()
				}
				delete(txAccessMsgs, dst)
				return errors.Timeout.New()
			}
		}
		if async {
			go f()
		} else {
			return f()
		}
		return nil
	}

	ackFuncWrapper := func(ack *SegmentAckMessage) {
		if ack.obo == 1 && ack.src != dst {
			lpnNode, _ := findNodeByAddr(dst)
			friendNode, _ := findNodeByAddr(ack.src)
			if lpnNode != nil {
				lpnNode.Friend = friendNode
			}
		}
	}

	pdu := generateRequest(opcode, payload)
	if _, ok := opcodeConfigReqRespMap[opcode]; ok {
		// it's a reponse of configuration request
		tpSendAccessMsgWithDevKey(pdu, dst, 5, ackFuncWrapper)
		return waitForResp()
	} else if _, ok := opcodeReqRespMap[opcode]; ok {
		//find the binded appkeys of the model

		var targetModel *Model
		for _, e := range node.Elements {
			for _, m := range e.Models {
				if msgs, ok := modelMap[m.ModelID]; ok {
					if funk.ContainsInt(msgs, int(opcode)) {
						targetModel = m
						break
					}
				}
			}
		}
		if targetModel == nil {
			return errors.OpcodeNotSupportedByNode.New().AddContextF("opcode:%4x, node:%4x", opcode, node.UnicastAddress)
		}
		if len(targetModel.BindedAppKeyIds) == 0 {
			return errors.NoAppKeyBindedToModel.New().AddContextF("node:%4x, model:%4x", node.UnicastAddress, targetModel.ModelID)
		}
		for _, appKeyId := range targetModel.BindedAppKeyIds {
			netKey, err := node.findNodeNetKeyByAppKeyIndex(appKeyId)
			if err != nil {
				return err
			}
			// onekey or all keys???
			var appKey *AppKey
			appKey, err = node.findNodeAppKeyByIndex(appKeyId)
			if err != nil {
				return err
			}
			tpSendAccessMsgWithAppKey(appKey, netKey, pdu, dst, 5)
		}
		// wait for each out message or just one transaction
		return waitForResp()
	} else {
		// vendor models
	}
	return nil
}

func registerModelMessageRxListener(cb *modelMsglistener) {
	modelMsgListeners = append(modelMsgListeners, cb)
}

func unregisterModelMessageRxListener(cb *modelMsglistener) {
	var index int
	for i, l := range modelMsgListeners {
		if l == cb {
			index = i
			break
		}
	}
	modelMsgListeners = append(modelMsgListeners[:index], modelMsgListeners[index+1:]...)
}

func modelMessageReceive(src, dst uint, data []byte) {
	opcode, payload := extractPayload(data)
	msg := &AccessMessage{src: src, dst: dst, opcode: opcode, payload: payload}
	loggerModel.Debugf("model Rx: %+#v", msg)
	if txMsg, ok := txAccessMsgs[src]; ok {
		if txMsg.expectedRespOpcode == msg.opcode {
			txMsg.ch <- msg
		} else {
			loggerModel.Errorf("unexpected response, expected opcode: %4x, actual: %4x", txMsg.expectedRespOpcode, msg.opcode)
		}
	}
	// todo: message publication
}

func extractPayload(data []byte) (uint, []byte) {
	opType := data[0] >> 6
	opSize := 1
	opCode := uint(data[0])
	switch opType {
	case 0x02:
		opSize = 2
		opCode = uint(binary.BigEndian.Uint16(data[:2]))
	case 0x03:
		opSize = 3
		opCode = uint(binary.BigEndian.Uint32(append([]byte{0x00}, data[:3]...)))
	}
	return opCode, data[opSize:]
}

func generateRequest(op uint, args ...interface{}) []byte {
	buffer := bytes.NewBuffer([]byte{})
	if op <= 0xFF {
		binary.Write(buffer, binary.BigEndian, byte(op))
	} else if op <= 0xFFFF {
		binary.Write(buffer, binary.BigEndian, uint16(op))
	} else {
		buf3 := []byte{}
		binary.BigEndian.PutUint32(buf3, uint32(op))
		binary.Write(buffer, binary.BigEndian, buf3[1:])
	}
	if args != nil {
		for _, arg := range args {
			switch reflect.TypeOf(arg).Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				binary.Write(buffer, binary.BigEndian, reflect.ValueOf(arg).Interface())
			case reflect.Slice:
				s := reflect.ValueOf(arg).Bytes()
				binary.Write(buffer, binary.BigEndian, s)
			}
		}
	}
	return buffer.Bytes()
}
