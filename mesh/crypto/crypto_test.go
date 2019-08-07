package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_s1(t *testing.T) {

	salt, _ := S1([]byte("test"))
	expected, _ := hex.DecodeString("b73cefbd641ef2ea598c2b6efb62f79c")
	assert.Equal(t, expected, salt, "they should be equal")

}

func Test_k1(t *testing.T) {
	n, _ := hex.DecodeString("3216d1509884b533248541792b877f98")
	salt, _ := hex.DecodeString("2ba14ffa0df84a2831938d57d276cab4")
	p, _ := hex.DecodeString("5a09d60797eeb4478aada59db3352a0d")
	expected, _ := hex.DecodeString("f6ed15a8934afbe7d83e8dcb57fcf5d7")
	actual, _ := K1(n, salt, p)
	assert.Equal(t, expected, actual, "they should be equal")
}

func Test_k2(t *testing.T) {
	n, _ := hex.DecodeString("f7a2a44f8e8a8029064f173ddc1e2b00")
	p := []byte{0}
	expectedEncKey, _ := hex.DecodeString("9f589181a0f50de73c8070c7a6d27f46")
	expectedPriKey, _ := hex.DecodeString("4c715bd4a64b938f99b453351653124f")
	nid, encKey, priKey, _ := K2(n, p)
	assert.Equal(t, uint(0x7f), nid, "they should be equal")
	assert.Equal(t, expectedEncKey, encKey, "they should be equal")
	assert.Equal(t, expectedPriKey, priKey, "they should be equal")
}

func Test_k3(t *testing.T) {
	n, _ := hex.DecodeString("f7a2a44f8e8a8029064f173ddc1e2b00")
	actual, _ := K3(n)

	expected, _ := hex.DecodeString("ff046958233db014")
	assert.Equal(t, expected, actual, "they should be equal")

}

func Test_k4(t *testing.T) {
	n, _ := hex.DecodeString("3216d1509884b533248541792b877f98")
	actual, _ := K4(n)

	assert.Equal(t, uint(0x38), actual, "they should be equal")

}

func Test_ccm(t *testing.T) {
	data, _ := hex.DecodeString("efb2255e6422d330088e09bb015ed707056700010203040b0c")
	sessionKey, _ := hex.DecodeString("c80253af86b33dfa450bbdb2a191fea3")
	sessionNonce, _ := hex.DecodeString("da7ddbe78b5f62b81d6847487e")
	actual, tag, _ := AES_CCM(sessionKey, sessionNonce, data, 8)
	expectedEnc, _ := hex.DecodeString("d0bd7f4a89a2ff6222af59a90a60ad58acfe3123356f5cec29")
	expectedTag, _ := hex.DecodeString("73e0ec50783b10c7")

	assert.Equal(t, expectedEnc, actual, "they should be equal")
	assert.Equal(t, expectedTag, tag, "they should be equal")

	plain, _ := AES_CCM_Decrypt(sessionKey, sessionNonce, append(expectedEnc, expectedTag...), 8)
	assert.Equal(t, data, plain, "they should be equal")
}

func Test_ecb(t *testing.T) {
	expected, _ := hex.DecodeString("6ca487507564")
	key, _ := hex.DecodeString("8b84eedec100067d670971dd2aa700cf")
	pdu, _ := hex.DecodeString("68eca487516765b5e5bfdacbaf6cb7fb6bff871f035444ce83a670df")
	data, _ := hex.DecodeString("000000000012345678")
	actual, _ := AES_ECB(key, append(data, pdu[7:14]...))
	assert.Equal(t, expected, actual[:6], "they should be equal")
}
