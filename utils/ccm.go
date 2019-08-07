// CCM Mode, defined in
// NIST Special Publication SP 800-38C.
// https://gist.github.com/hirochachacha/abb76ff71573dea2ef42
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// RunExamples run examples defined in Appendix C.
func RunExamples() {
	C4A := make([]byte, 524288/8)
	for i := range C4A {
		C4A[i] = byte(i)
	}

	examples := []struct {
		Key        []byte
		Nonce      []byte
		Data       []byte
		PlainText  []byte
		CipherText []byte
		TagLen     int
	}{
		{ // C.1
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			[]byte{0x20, 0x21, 0x22, 0x23},
			[]byte{0x71, 0x62, 0x01, 0x5b, 0x4d, 0xac, 0x25, 0x5d},
			4,
		},
		{ // C.2
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f},
			[]byte{0xd2, 0xa1, 0xf0, 0xe0, 0x51, 0xea, 0x5f, 0x62, 0x08, 0x1a, 0x77, 0x92, 0x07, 0x3d, 0x59, 0x3d, 0x1f, 0xc6, 0x4f, 0xbf, 0xac, 0xcd},
			6,
		},
		{ // C.3
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13},
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37},

			[]byte{0xe3, 0xb2, 0x01, 0xa9, 0xf5, 0xb7, 0x1a, 0x7a, 0x9b, 0x1c, 0xea, 0xec, 0xcd, 0x97, 0xe7, 0x0b, 0x61, 0x76, 0xaa, 0xd9, 0xa4, 0x42, 0x8a, 0xa5, 0x48, 0x43, 0x92, 0xfb, 0xc1, 0xb0, 0x99, 0x51},
			8,
		},
		{ // C.4
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c},
			C4A,
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f},

			[]byte{0x69, 0x91, 0x5d, 0xad, 0x1e, 0x84, 0xc6, 0x37, 0x6a, 0x68, 0xc2, 0x96, 0x7e, 0x4d, 0xab, 0x61, 0x5a, 0xe0, 0xfd, 0x1f, 0xae, 0xc4, 0x4c, 0xc4, 0x84, 0x82, 0x85, 0x29, 0x46, 0x3c, 0xcf, 0x72, 0xb4, 0xac, 0x6b, 0xec, 0x93, 0xe8, 0x59, 0x8e, 0x7f, 0x0d, 0xad, 0xbc, 0xea, 0x5b},
			14,
		},
	}

	for i, ex := range examples {
		print("Example ", i+1, ": ")

		c, err := aes.NewCipher(ex.Key)
		if err != nil {
			panic(err)
		}

		ccm, err := NewCCMWithNonceAndTagSizes(c, len(ex.Nonce), ex.TagLen)
		if err != nil {
			panic(err)
		}

		CipherText := ccm.Seal(nil, ex.Nonce, ex.PlainText, ex.Data)

		if !bytes.Equal(ex.CipherText, CipherText) {
			panic("fail")
		}

		PlainText, err := ccm.Open(nil, ex.Nonce, ex.CipherText, ex.Data)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(ex.PlainText, PlainText) {
			panic("fail")
		}

		println("ok")
	}
}

// CBC-MAC implementation
type mac struct {
	ci []byte
	p  int
	c  cipher.Block
}

func newMAC(c cipher.Block) *mac {
	return &mac{
		c:  c,
		ci: make([]byte, c.BlockSize()),
	}
}

func (m *mac) Reset() {
	for i := range m.ci {
		m.ci[i] = 0
	}
	m.p = 0
}

func (m *mac) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if m.p >= len(m.ci) {
			m.c.Encrypt(m.ci, m.ci)
			m.p = 0
		}
		m.ci[m.p] ^= c
		m.p++
	}
	return len(p), nil
}

// PadZero emulates zero byte padding.
func (m *mac) PadZero() {
	if m.p != 0 {
		m.c.Encrypt(m.ci, m.ci)
		m.p = 0
	}
}

func (m *mac) Sum(in []byte) []byte {
	if m.p != 0 {
		m.c.Encrypt(m.ci, m.ci)
		m.p = 0
	}
	return append(in, m.ci...)
}

func (m *mac) Size() int { return len(m.ci) }

func (m *mac) BlockSize() int { return 16 }

type ccm struct {
	c         cipher.Block
	mac       *mac
	nonceSize int
	tagSize   int
}

// NewCCMWithNonceAndTagSizes returns the given 128-bit, block cipher wrapped in Counter with CBC-MAC Mode, which accepts nonces of the given length.
// the formatting of this function is defined in SP800-38C, Appendix A.
// Each arguments have own valid range:
//   nonceSize should be one of the {7, 8, 9, 10, 11, 12, 13}.
//   tagSize should be one of the {4, 6, 8, 10, 12, 14, 16}.
//   Otherwise, it panics.
// The maximum payload size is defined as 1<<uint((15-nonceSize)*8)-1.
// If the given payload size exceeds the limit, it returns a error (Seal returns nil instead).
// The payload size is defined as len(plaintext) on Seal, len(ciphertext)-tagSize on Open.
func NewCCMWithNonceAndTagSizes(c cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if c.BlockSize() != 16 {
		return nil, errors.New("cipher: CCM mode requires 128-bit block cipher")
	}

	if !(7 <= nonceSize && nonceSize <= 13) {
		return nil, errors.New("cipher: invalid nonce size")
	}

	if !(4 <= tagSize && tagSize <= 16 && tagSize&1 == 0) {
		return nil, errors.New("cipher: invalid tag size")
	}

	return &ccm{
		c:         c,
		mac:       newMAC(c),
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}, nil
}

func (ccm *ccm) NonceSize() int {
	return ccm.nonceSize
}

func (ccm *ccm) Overhead() int {
	return ccm.tagSize
}

func (ccm *ccm) Seal(dst, nonce, plaintext, data []byte) []byte {
	if len(nonce) != ccm.nonceSize {
		panic("cipher: incorrect nonce length given to CCM")
	}

	// AEAD interface doesn't provide a way to return errors.
	// So it returns nil instead.
	if maxUvarint(15-ccm.nonceSize) < uint64(len(plaintext)) {
		return nil
	}

	ret, ciphertext := sliceForAppend(dst, len(plaintext)+ccm.mac.Size())

	// Formatting of the Counter Blocks are defined in A.3.
	Ctr := make([]byte, 16)               // Ctr0
	Ctr[0] = byte(15 - ccm.nonceSize - 1) // [q-1]3
	copy(Ctr[1:], nonce)                  // N

	S0 := ciphertext[len(plaintext):] // S0
	ccm.c.Encrypt(S0, Ctr)

	Ctr[15] = 1 // Ctr1

	ctr := cipher.NewCTR(ccm.c, Ctr)

	ctr.XORKeyStream(ciphertext, plaintext)

	T := ccm.getTag(Ctr, data, plaintext)

	xorBytes(S0, S0, T) // T^S0

	return ret[:len(plaintext)+ccm.tagSize]
}

func (ccm *ccm) Open(dst, nonce, ciphertext, data []byte) ([]byte, error) {
	if len(nonce) != ccm.nonceSize {
		return nil, errors.New("cipher: incorrect nonce length given to CCM")
	}

	if len(ciphertext) <= ccm.tagSize {
		return nil, errors.New("cipher: incorrect ciphertext length given to CCM")
	}

	if maxUvarint(15-ccm.nonceSize) < uint64(len(ciphertext)-ccm.tagSize) {
		return nil, errors.New("cipher: len(ciphertext)-tagSize exceeds the maximum payload size")
	}

	ret, plaintext := sliceForAppend(dst, len(ciphertext)-ccm.tagSize)

	// Formatting of the Counter Blocks are defined in A.3.
	Ctr := make([]byte, 16)               // Ctr0
	Ctr[0] = byte(15 - ccm.nonceSize - 1) // [q-1]3
	copy(Ctr[1:], nonce)                  // N

	S0 := make([]byte, 16) // S0
	ccm.c.Encrypt(S0, Ctr)

	Ctr[15] = 1 // Ctr1

	ctr := cipher.NewCTR(ccm.c, Ctr)

	ctr.XORKeyStream(plaintext, ciphertext[:len(plaintext)])

	T := ccm.getTag(Ctr, data, plaintext)

	xorBytes(T, T, S0)

	if !bytes.Equal(T[:ccm.tagSize], ciphertext[len(plaintext):]) {
		return nil, errors.New("cipher: message authentication failed")
	}

	return ret, nil
}

// getTag reuses a Ctr block for making the B0 block because of some parts are the same.
// For more details, see A.2 and A.3.
func (ccm *ccm) getTag(Ctr, data, plaintext []byte) []byte {
	ccm.mac.Reset()

	B := Ctr                                                // B0
	B[0] |= byte(((ccm.tagSize - 2) / 2) << 3)              // [(t-2)/2]3
	putUvarint(B[1+ccm.nonceSize:], uint64(len(plaintext))) // Q

	if len(data) > 0 {
		B[0] |= 1 << 6 // Adata

		ccm.mac.Write(B)

		if len(data) < (1<<15 - 1<<7) {
			putUvarint(B[:2], uint64(len(data)))

			ccm.mac.Write(B[:2])
		} else if len(data) <= 1<<31-1 {
			B[0] = 0xff
			B[1] = 0xfe
			putUvarint(B[2:6], uint64(len(data)))

			ccm.mac.Write(B[:6])
		} else {
			B[0] = 0xff
			B[1] = 0xff
			putUvarint(B[2:10], uint64(len(data)))

			ccm.mac.Write(B[:10])
		}
		ccm.mac.Write(data)
		ccm.mac.PadZero()
	} else {
		ccm.mac.Write(B)
	}

	ccm.mac.Write(plaintext)
	ccm.mac.PadZero()

	return ccm.mac.Sum(nil)
}

func maxUvarint(n int) uint64 {
	return 1<<uint(n*8) - 1
}

// put uint64 as big endian.
func putUvarint(bs []byte, u uint64) {
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(u >> uint(8*(len(bs)-1-i)))
	}
}

// defined in crypto/cipher/gcm.go
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

// defined in crypto/cipher/xor.go
func xorBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}
