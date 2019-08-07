package mesh

import (
	. "ble-mesh/mesh/def"
	"ble-mesh/utils"
)

const (
	opGenericOnOffGet               = 0x8201
	opGenericOnOffSet               = 0x8202
	opGenericOnOffSetUnacknowledged = 0x8203
	opGenericOnOffStatus            = 0x8204

	opGenericLevelGet               = 0x8205
	opGenericLevelSet               = 0x8206
	opGenericLevelSetUnacknowledged = 0x8207
	opGenericLevelStatus            = 0x8208
	opGenericDeltaSet               = 0x8209
	opGenericDeltaSetUnacknowledged = 0x820A
	opGenericMoveSet                = 0x820B
	opGenericMoveSetUnacknowledged  = 0x820C

	opGenericDefaultTransitionTimeGet               = 0x820D
	opGenericDefaultTransitionTimeSet               = 0x820E
	opGenericDefaultTransitionTimeSetUnacknowledged = 0x820F
	opGenericDefaultTransitionTimeStatus            = 0x8210

	opGenericOnPowerUpGet               = 0x8211
	opGenericOnPowerUpStatus            = 0x8212
	opGenericOnPowerUpSet               = 0x8213
	opGenericOnPowerUpSetUnacknowledged = 0x8214

	opGenericPowerLevelGet               = 0x8215
	opGenericPowerLevelSet               = 0x8216
	opGenericPowerLevelSetUnacknowledged = 0x8217
	opGenericPowerLevelStatus            = 0x8218
	opGenericPowerLastGet                = 0x8219
	opGenericPowerLastStatus             = 0x821A
	opGenericPowerDefaultGet             = 0x821B
	opGenericPowerDefaultStatus          = 0x821C
	opGenericPowerRangeGet               = 0x821D
	opGenericPowerRangeStatus            = 0x821E

	opGenericPowerDefaultSet               = 0x821F
	opGenericPowerDefaultSetUnacknowledged = 0x8220
	opGenericPowerRangeSet                 = 0x8221

	opGenericPowerRangeSetUnacknowledged           = 0x8222
	opGenericBatteryGet                            = 0x8223
	opGenericBatteryStatus                         = 0x8224
	opGenericLocationGlobalGet                     = 0x8225
	opGenericLocationGlobalStatus                  = 0x40
	opGenericLocationLocalGet                      = 0x8226
	opGenericLocationLocalStatus                   = 0x8227
	opGenericLocationGlobalSet                     = 0x41
	opGenericLocationGlobalSetUnacknowledged       = 0x42
	opGenericLocationLocalSet                      = 0x8228
	opGenericLocationLocalSetUnacknowledged        = 0x8229
	opGenericManufacturerPropertiesGet             = 0x822A
	opGenericManufacturerPropertiesStatus          = 0x43
	opGenericManufacturerPropertyGet               = 0x822B
	opGenericManufacturerPropertySet               = 0x44
	opGenericManufacturerPropertySetUnacknowledged = 0x45
	opGenericManufacturerPropertyStatus            = 0x46
	opGenericAdminPropertiesGet                    = 0x822C
	opGenericAdminPropertiesStatus                 = 0x47
	opGenericAdminPropertyGet                      = 0x822D
	opGenericAdminPropertySet                      = 0x48
	opGenericAdminPropertySetUnacknowledged        = 0x49
	opGenericAdminPropertyStatus                   = 0x4A
	opGenericUserPropertiesGet                     = 0x822E
	opGenericUserPropertiesStatus                  = 0x4B
	opGenericUserPropertyGet                       = 0x822F
	opGenericUserPropertySet                       = 0x4C
	opGenericUserPropertySetUnacknowledged         = 0x4D
	opGenericUserPropertyStatus                    = 0x4E
	opGenericClientPropertiesGet                   = 0x4F
	opGenericClientPropertiesStatus                = 0x50
)

var (
	loggerOnOffCli = utils.CreateLogger("OnOffClient")
)

type OnOffState struct {
	OnOff uint `json:"onoff"`
}

func (m *Model) handleOnOffResponse(n *Node, d interface{}) error {
	resp := d.(GenericOnOffStatusMessageParameters)
	m.State = OnOffState{
		OnOff: resp.PresentOnOff,
	}
	return nil
}

func GenericOnOffGet(dst uint) error {
	m, err := findModelDirectly(dst, GenericOnOffServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opGenericOnOffGet, nil, m.handleOnOffResponse)
}

func genericOnOffSet(ack bool, dst uint, onoff uint) error {
	m, err := findModelDirectly(dst, GenericOnOffServer)
	if err != nil {
		return err
	}
	params := &GenericOnOffSetMessageParameters{
		OnOff: onoff,
	}
	op := uint(opGenericOnOffSetUnacknowledged)
	if ack {
		op = opGenericOnOffSet
	}
	return modelSendTmplParsedWithTID(dst, op, params, m.handleOnOffResponse)
}

func GenericOnOffSet(dst, onoff uint) error {
	return genericOnOffSet(true, dst, onoff)
}

func GenericOnOffSetUnacknowledged(dst, onoff uint) error {
	return genericOnOffSet(false, dst, onoff)
}

func handleGenericLevelResponse(n *Node, d interface{}) error {
	return nil
}

func GenericLevelGet(dst uint) error {
	return modelSendTmplParsed(false, dst, opGenericLevelGet, nil, handleGenericLevelResponse)
}

func genericLevelSet(ack bool, dst uint, level uint) error {
	params := &GenericLevelSetMessageParameters{
		Level: level,
	}
	op := uint(opGenericLevelSetUnacknowledged)
	if ack {
		op = opGenericLevelSet
	}
	return modelSendTmplParsedWithTID(dst, op, params, handleGenericLevelResponse)
}

func GenericLevelSet(dst, level uint) error {
	return genericLevelSet(true, dst, level)
}

func GenericLevelSetUnacknowledged(dst, level uint) error {
	return genericLevelSet(false, dst, level)
}
