package utils

import (
	"ble-mesh/mesh/def"
	"testing"

	"github.com/stretchr/testify/assert"
)

type structA struct {
	Id uint `bits:"8"`
	A  uint `bits:"5"`
	B  uint `bits:"3"`
	C  uint `bits:"8/16"`
}

type structB struct {
	Id uint   `bits:"8"`
	A  uint   `bits:"5"`
	B  uint   `bits:"3"`
	C  []uint `bits:"8"`
}

type structC struct {
	A uint   `bits:"5"`
	B uint   `bits:"3"`
	C []byte `bits:"8"`
}

func Test_Unpack(t *testing.T) {
	data := []byte{0x1F, 0xEF, 0xD2, 0x40, 0x59}
	var a, b, c, d, e int
	var z []byte
	UnpackBE(data, "03,1,8,4,20,4", &a, &b, &c, &d, &e)
	assert.Equal(t, 1, a, "they should be equal")
	assert.Equal(t, 254, b, "they should be equal")
	assert.Equal(t, 15, c, "they should be equal")
	assert.Equal(t, 0xD2405, d, "they should be equal")
	assert.Equal(t, 9, e, "they should be equal")
	z, _ = PackBE("03,1,8,4,20,4", a, b, c, d, e)
	assert.Equal(t, data, z, "they should be equal")

	data = []byte{0xF8, 0xEF, 0xD2, 0x40, 0x59}
	UnpackLE(data, "03,1,8,4,20,4", &a, &b, &c, &d, &e)
	assert.Equal(t, 1, a, "they should be equal")
	assert.Equal(t, 0xFF, b, "they should be equal")
	assert.Equal(t, 0x0E, c, "they should be equal")
	assert.Equal(t, 0x940D2, d, "they should be equal")
	assert.Equal(t, 5, e, "they should be equal")
	z, _ = PackLE("03,1,8,4,20,4", a, b, c, d, e)
	assert.Equal(t, data, z, "they should be equal")
}

func Test_StructUnpack(t *testing.T) {
	s := new(structA)
	e := structA{0x51, 0x0E, 0x03, 0x1234}
	data := []byte{0x51, 0x73, 0x12, 0x34}
	UnpackStructBE(data, s)
	assert.Equal(t, e, *s, "should equal")

	e = structA{0x51, 0x0E, 0x03, 0x12}
	data = []byte{0x51, 0x73, 0x12}
	UnpackStructBE(data, s)
	assert.Equal(t, e, *s, "should equal")

	r := new(structB)
	z := structB{0x51, 0x0E, 0x03, []uint{0x12, 0x34}}
	data = []byte{0x51, 0x73, 0x12, 0x34}
	UnpackStructBE(data, r)
	assert.Equal(t, z, *r, "should equal")

	p := structC{0x0F, 0x03, []byte{0x20, 0x77}}
	data = []byte{0x7b, 0x20, 0x77}
	x, _ := PackStructBE(p)
	assert.Equal(t, data, x, "should equal")

	params := &def.ConfigModelPublicationSetMessageParameters{
		ElementAddress: 0x1000,
		PublishAddress: 0xc000,
		AppKeyIndex:    0,
		CredentialFlag: 1,
		PublishTTL:     5,
		PublishPeriod: def.PublishPeriodFormat{
			NumberOfSteps:  3,
			StepResolution: 2,
		},
		PublishRetransmitCount:         5,
		PublishRetransmitIntervalSteps: 2,
		ModelIdentifier:                0x1000,
	}
	data, _ = PackStructBE(params)

	param1 := &def.ConfigModelSubscriptionAddMessageParameters{
		ElementAddress:  0x1000,
		Address:         0x3333,
		ModelIdentifier: 0xFFFFF,
	}
	data, _ = PackStructBE(param1)
}
