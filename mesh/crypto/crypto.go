package crypto

import (
	"ble-mesh/utils"
	aesCipher "crypto/aes"
	"fmt"

	"github.com/ChengjinWu/aescrypto"
	"github.com/aead/cmac/aes"
)

func AES_ECB(key, data []byte) ([]byte, error) {
	return aescrypto.AesEcbPkcs5Encrypt(data, key)
}

func AES_CCM(key, nonce, data []byte, tagSize int) (enc, mic []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	c, err := aesCipher.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	ccm, err := utils.NewCCMWithNonceAndTagSizes(c, len(nonce), tagSize)
	if err != nil {
		return nil, nil, err
	}

	cipher := ccm.Seal(nil, nonce, data, nil)
	return cipher[:len(cipher)-tagSize], cipher[len(cipher)-tagSize:], nil
}

func AES_CCM_Decrypt(key, nonce, cipher []byte, tagSize int) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	c, err := aesCipher.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ccm, err := utils.NewCCMWithNonceAndTagSizes(c, len(nonce), tagSize)
	if err != nil {
		return nil, err
	}

	return ccm.Open(nil, nonce, cipher, nil)
}

func AES_CMAC(key, data []byte) ([]byte, error) {
	return aes.Sum(data, key, aesCipher.BlockSize)
}

func S1(data []byte) ([]byte, error) {
	key := make([]byte, 16, 16)
	return AES_CMAC(key, data)
}

func K1(n, salt, p []byte) ([]byte, error) {
	t, _ := AES_CMAC(salt, n)
	return AES_CMAC(t, p)
}

func K2(n, p []byte) (nid uint, encKey, priKey []byte, err error) {
	salt, err := S1([]byte("smk2"))
	t, err := AES_CMAC(salt, n)
	t1, err := AES_CMAC(t, append(p, 0x01))
	t2, err := AES_CMAC(t, append(append(t1, p...), 0x02))
	t3, err := AES_CMAC(t, append(append(t2, p...), 0x03))
	k := append(append(t1, t2...), t3...)
	k = k[len(k)-33:]
	nid = uint(k[0] & 0x7F)
	encKey = k[1:17]
	priKey = k[17:]
	return
}

func K3(n []byte) ([]byte, error) {
	salt, _ := S1([]byte("smk3"))
	t, _ := AES_CMAC(salt, n)
	r, _ := AES_CMAC(t, append([]byte("id64"), 0x01))
	return r[len(r)-8:], nil
}

func K4(n []byte) (uint, error) {
	salt, _ := S1([]byte("smk4"))
	t, _ := AES_CMAC(salt, n)
	r, _ := AES_CMAC(t, append([]byte("id6"), 0x01))
	return uint(r[len(r)-1] & 0x3F), nil
}
