package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Arr []byte
	Bin uint16
	Str string
}

func Test_struct_pack(t *testing.T) {
	data := &testStruct{[]byte{0x02, 0x03}, 0x57ba, "ss"}
	expected := []byte{0x02, 0x03, 0x57, 0xba, 115, 115}

	gen, _ := WriteStructToBuffer(data)
	assert.Equal(t, expected, gen, "they should be equal")

}
