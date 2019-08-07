package mesh

import (
	. "ble-mesh/mesh/def"
	"ble-mesh/utils"
	"ble-mesh/utils/errors"
)

const (
	opLightLightnessGet                         = 0x824B
	opLightLightnessSet                         = 0x824C
	opLightLightnessSetUnacknowledged           = 0x824D
	opLightLightnessStatus                      = 0x824E
	opLightLightnessLinearGet                   = 0x824F
	opLightLightnessLinearSet                   = 0x8250
	opLightLightnessLinearSetUnacknowledged     = 0x8251
	opLightLightnessLinearStatus                = 0x8252
	opLightLightnessLastGet                     = 0x8253
	opLightLightnessLastStatus                  = 0x8254
	opLightLightnessDefaultGet                  = 0x8255
	opLightLightnessDefaultStatus               = 0x8256
	opLightLightnessRangeGet                    = 0x8257
	opLightLightnessRangeStatus                 = 0x8258
	opLightLightnessDefaultSet                  = 0x8259
	opLightLightnessDefaultSetUnacknowledged    = 0x825A
	opLightLightnessRangeSet                    = 0x825B
	opLightLightnessRangeSetUnacknowledged      = 0x825C
	opLightCTLGet                               = 0x825D
	opLightCTLSet                               = 0x825E
	opLightCTLSetUnacknowledged                 = 0x825F
	opLightCTLStatus                            = 0x8260
	opLightCTLTemperatureGet                    = 0x8261
	opLightCTLTemperatureRangeGet               = 0x8262
	opLightCTLTemperatureRangeStatus            = 0x8263
	opLightCTLTemperatureSet                    = 0x8264
	opLightCTLTemperatureSetUnacknowledged      = 0x8265
	opLightCTLTemperatureStatus                 = 0x8266
	opLightCTLDefaultGet                        = 0x8267
	opLightCTLDefaultStatus                     = 0x8268
	opLightCTLDefaultSet                        = 0x8269
	opLightCTLDefaultSetUnacknowledged          = 0x826A
	opLightCTLTemperatureRangeSet               = 0x826B
	opLightCTLTemperatureRangeSetUnacknowledged = 0x826C
	opLightHSLGet                               = 0x826D
	opLightHSLHueGet                            = 0x826E
	opLightHSLHueSet                            = 0x826F
	opLightHSLHueSetUnacknowledged              = 0x8270
	opLightHSLHueStatus                         = 0x8271
	opLightHSLSaturationGet                     = 0x8272
	opLightHSLSaturationSet                     = 0x8273
	opLightHSLSaturationSetUnacknowledged       = 0x8274
	opLightHSLSaturationStatus                  = 0x8275
	opLightHSLSet                               = 0x8276
	opLightHSLSetUnacknowledged                 = 0x8277
	opLightHSLStatus                            = 0x8278
	opLightHSLTargetGet                         = 0x8279
	opLightHSLTargetStatus                      = 0x827A
	opLightHSLDefaultGet                        = 0x827B
	opLightHSLDefaultStatus                     = 0x827C
	opLightHSLRangeGet                          = 0x827D
	opLightHSLRangeStatus                       = 0x827E
	opLightHSLDefaultSet                        = 0x827F
	opLightHSLDefaultSetUnacknowledged          = 0x8280
	opLightHSLRangeSet                          = 0x8281
	opLightHSLRangeSetUnacknowledged            = 0x82
	opLightxyLGet                               = 0x8283
	opLightxyLSet                               = 0x8284
	opLightxyLSetUnacknowledged                 = 0x8285
	opLightxyLStatus                            = 0x8286
	opLightxyLTargetGet                         = 0x8287
	opLightxyLTargetStatus                      = 0x8288
	opLightxyLDefaultGet                        = 0x8289
	opLightxyLDefaultStatus                     = 0x828A
	opLightxyLRangeGet                          = 0x828B
	opLightxyLRangeStatus                       = 0x828C
	opLightxyLDefaultSet                        = 0x828D
	opLightxyLDefaultSetUnacknowledged          = 0x828E
	opLightxyLRangeSet                          = 0x828F
	opLightxyLRangeSetUnacknowledged            = 0x8290
	opLightLCModeGet                            = 0x8291
	opLightLCModeSet                            = 0x8292
	opLightLCModeSetUnacknowledged              = 0x8293
	opLightLCModeStatus                         = 0x8294
	opLightLCOMGet                              = 0x8295
	opLightLCOMSet                              = 0x8296
	opLightLCOMSetUnacknowledged                = 0x8297
	opLightLCOMStatus                           = 0x8298
	opLightLCLightOnOffGet                      = 0x8299
	opLightLCLightOnOffSet                      = 0x829A
	opLightLCLightOnOffSetUnacknowledged        = 0x829B
	opLightLCLightOnOffStatus                   = 0x829C
	opLightLCPropertyGet                        = 0x829D
	opLightLCPropertySet                        = 0x62
	opLightLCPropertySetUnacknowledged          = 0x63
	opLightLCPropertyStatus                     = 0x64
)

var (
	loggerLightCli = utils.CreateLogger("LightClient")
)

type LightnessState struct {
	Lightness        uint `json:"lightness"`
	LightnessDefault uint `json:"lightnessDefault"`
	RangeMin         uint `json:"rangeMin"`
	RangeMax         uint `json:"rangeMax"`
}

type LightCtlState struct {
	CtlLightness       uint `json:"ctlLightness"`
	CtlTemperature     uint `json:"ctlTemperature"`
	LightnessDefault   uint `json:"lightnessDefault"`
	TemperatureDefault uint `json:"temperatureDefault"`
	DeltaUVDefault     int  `json:"deltaUVDefault"`
	RangeMin           uint `json:"rangeMin"`
	RangeMax           uint `json:"rangeMax"`
}

type LightCtlTemperatureState struct {
	CtlTemperature uint `json:"ctlTemperature"`
	CtlDeltaUV     int  `json:"ctlDeltaUV"`
}

func calcLightness(prc uint) uint {
	return prc*0xFFFE/100 + 1
}

func calcLightnessPrc(value uint) uint {
	return (uint(value) - 1) * 100 / 0xFFFE
}

func calcTemperature(prc uint) uint {
	return (0x4E20-0x0320)*prc/100 + 0x0320
}

func calcTemperaturePrc(value uint) uint {
	return (uint(value) - 0x320) * 100 / (0x4E20 - 0x0320)
}

func calcDeltaUV(offset float32) int {
	return int(offset * 32767)
}

func (m *Model) handleLightnessResponse(n *Node, d interface{}) error {
	resp := d.(LightLightnessStatusMessageParameters)
	state, _ := m.State.(LightnessState)
	state.Lightness = resp.PresentLightness
	m.State = state
	loggerLightCli.Debugf("present lightness: %d%%", calcLightnessPrc(resp.PresentLightness))
	return nil
}

func LightnessGet(dst uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightLightnessGet, nil, m.handleLightnessResponse)
}

func LightnessSet(dst uint, lightness uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	req := &LightLightnessSetMessageParameters{
		Lightness: lightness,
	}
	return modelSendTmplParsedWithTID(dst, opLightLightnessSet, req, m.handleLightnessResponse)
}

func (m *Model) handleLightnessLinearResponse(n *Node, d interface{}) error {
	resp := d.(LightLightnessLinearStatusMessageParameters)
	loggerLightCli.Debugf("present lightness: %d%%", calcLightnessPrc(resp.PresentLightness))
	return nil
}

func LightnessLinearGet(dst uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightLightnessLinearGet, nil, m.handleLightnessLinearResponse)
}

func LightnessLinearSet(dst uint, lightness uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	req := &LightLightnessLinearSetMessageParameters{
		Lightness: lightness,
	}
	return modelSendTmplParsedWithTID(dst, opLightLightnessLinearSet, req, m.handleLightnessLinearResponse)
}

func (m *Model) handleLightnessDefaultResponse(n *Node, d interface{}) error {
	resp := d.(LightLightnessDefaultStatusMessageParameters)
	state, _ := m.State.(LightnessState)
	state.LightnessDefault = resp.Lightness
	m.State = state
	loggerLightCli.Debugf("default lightness: %d%%", calcLightnessPrc(resp.Lightness))
	return nil
}

func LightnessDefaultGet(dst uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightLightnessDefaultGet, nil, m.handleLightnessDefaultResponse)
}

func LightnessDefaultSet(dst uint, lightness uint) error {

	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	req := &LightLightnessDefaultSetMessageParameters{
		Lightness: lightness,
	}
	return modelSendTmplParsedWithTID(dst, opLightLightnessDefaultSet, req, m.handleLightnessDefaultResponse)
}

func (m *Model) handleLightnessRangeResponse(n *Node, d interface{}) error {
	resp := d.(LightLightnessRangeStatusMessageParameters)
	if resp.StatusCode == 0 {
		state, _ := m.State.(LightnessState)
		state.RangeMin = resp.RangeMin
		state.RangeMax = resp.RangeMax
		m.State = state
		loggerLightCli.Debugf("range min: %x, max: %x", resp.RangeMin, resp.RangeMax)
	} else if resp.StatusCode == 1 {
		return errors.CannotSetRangeMin.New()
	} else if resp.StatusCode == 2 {
		return errors.CannotSetRangeMax.New()
	}
	return nil
}

func LightnessRangeGet(dst uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightLightnessRangeGet, nil, m.handleLightnessRangeResponse)
}

func LightnessRangeSet(dst uint, min uint, max uint) error {
	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	req := &LightLightnessRangeSetMessageParameters{
		RangeMin: min,
		RangeMax: max,
	}
	return modelSendTmplParsedWithTID(dst, opLightLightnessRangeSet, req, m.handleLightnessRangeResponse)
}

func (m *Model) handleLightCtlResponse(n *Node, d interface{}) error {
	resp := d.(LightCTLStatusMessageParameters)
	loggerLightCli.Debugf("present lightness: %d%%, temperature: %d%%",
		calcLightnessPrc(resp.PresentCTLLightness),
		calcTemperaturePrc(resp.PresentCTLTemperature))
	state, _ := m.State.(LightCtlState)
	state.CtlLightness = resp.PresentCTLLightness
	state.CtlTemperature = resp.PresentCTLTemperature
	return nil
}

func LightCtlGet(dst uint) error {
	m, err := findModelDirectly(dst, LightCTLServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightCTLGet, nil, m.handleLightCtlResponse)

}

func LightCtlSet(dst uint, lightness, temperature uint, offsetDeltaUV float32) error {
	m, err := findModelDirectly(dst, LightCTLServer)
	if err != nil {
		return err
	}
	req := &LightCTLSetMessageParameters{
		CTLLightness:   lightness,
		CTLTemperature: temperature,
		CTLDeltaUV:     calcDeltaUV(offsetDeltaUV),
	}
	return modelSendTmplParsedWithTID(dst, opLightCTLSet, req, m.handleLightCtlResponse)
}

func (m *Model) handleLightCtlTemperatureRangeResponse(n *Node, d interface{}) error {
	resp := d.(LightCTLTemperatureRangeStatusMessageParameters)
	state, _ := m.State.(LightCtlState)
	if resp.StatusCode == 0 {
		state.RangeMin = resp.RangeMin
		state.RangeMax = resp.RangeMax
		m.State = state
		loggerLightCli.Debugf("range min: %x, max: %x", resp.RangeMin, resp.RangeMax)
	} else if resp.StatusCode == 1 {
		return errors.CannotSetRangeMin.New()
	} else if resp.StatusCode == 2 {
		return errors.CannotSetRangeMax.New()
	}
	return nil
}

func LightCtlTemperatureRangeGet(dst uint) error {
	m, err := findModelDirectly(dst, LightCTLServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightCTLTemperatureRangeGet, nil, m.handleLightCtlTemperatureRangeResponse)
}

func LightCtlTemperatureRangeSet(dst uint, temp uint, offsetDeltaUV float32) error {
	m, err := findModelDirectly(dst, LightCTLServer)
	if err != nil {
		return err
	}
	req := &LightCTLTemperatureSetMessageParameters{
		CTLTemperature: temp,
		CTLDeltaUV:     calcDeltaUV(offsetDeltaUV),
	}
	return modelSendTmplParsedWithTID(dst, opLightCTLTemperatureRangeSet, req, m.handleLightCtlTemperatureRangeResponse)
}

func (m *Model) handleCtlDefaultResponse(n *Node, d interface{}) error {
	resp := d.(LightCTLDefaultStatusMessageParameters)
	state, _ := m.State.(LightCtlState)
	state.LightnessDefault = resp.Lightness
	state.TemperatureDefault = resp.Temperature
	state.DeltaUVDefault = resp.DeltaUV
	m.State = state
	// loggerLightCli.Debugf("default lightness: %d%%", calcLightnessPrc(resp.Lightness))
	return nil
}

func LightCtlDefaultGet(dst uint) error {
	m, err := findModelDirectly(dst, LightCTLServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightCTLDefaultGet, nil, m.handleCtlDefaultResponse)
}

func LightCtlDefaultSet(dst uint, lightness, temperature uint, offsetDeltaUV float32) error {

	m, err := findModelDirectly(dst, LightLightnessServer)
	if err != nil {
		return err
	}
	req := &LightCTLDefaultSetMessageParameters{
		Lightness:   lightness,
		Temperature: temperature,
		DeltaUV:     calcDeltaUV(offsetDeltaUV),
	}
	return modelSendTmplParsedWithTID(dst, opLightCTLDefaultSet, req, m.handleCtlDefaultResponse)
}

func (m *Model) handleLightCtlTemperatureResponse(n *Node, d interface{}) error {
	resp := d.(LightCTLTemperatureStatusMessageParameters)
	loggerLightCli.Debugf("present temperature: %d%%",
		calcTemperaturePrc(resp.PresentCTLTemperature))
	state, _ := m.State.(LightCtlTemperatureState)
	state.CtlTemperature = resp.PresentCTLTemperature
	state.CtlDeltaUV = resp.PresentCTLDeltaUV
	return nil
}

func LightCtlTemperatureGet(dst uint) error {
	m, err := findModelDirectly(dst, LightCTLTemperatureServer)
	if err != nil {
		return err
	}
	return modelSendTmplParsed(false, dst, opLightCTLTemperatureGet, nil, m.handleLightCtlTemperatureResponse)
}

func LightCtlTemperatureSet(dst uint, temp uint, offsetDeltaUV float32) error {
	m, err := findModelDirectly(dst, LightCTLTemperatureServer)
	if err != nil {
		return err
	}
	req := &LightCTLTemperatureSetMessageParameters{
		CTLTemperature: temp,
		CTLDeltaUV:     calcDeltaUV(offsetDeltaUV),
	}
	return modelSendTmplParsedWithTID(dst, opLightCTLTemperatureSet, req, m.handleLightCtlTemperatureResponse)
}
