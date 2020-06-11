package cbc_cfb_ofb_ctr_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit3(aesMode int, baseKey, baseIV []byte) func(uint8) func([]byte) {
	// constant randomly picked key
	key := make([]byte, keySize)
	copy(key, baseKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// constant randomly picked iv
	iv := make([]byte, ivSize)
	copy(iv, baseIV)

	switch aesMode {
	case 0:
		mode := cipher.NewCBCEncrypter(block, iv)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				if len(plaintext)%aes.BlockSize != 0 {
					panic("plaintext is not a multiple of the block size")
				}
				ciphertext := make([]byte, len(plaintext))
				mode.CryptBlocks(ciphertext, plaintext)
			}
		}
	case 1:
		stream := cipher.NewCFBEncrypter(block, iv)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		}
	case 2:
		stream := cipher.NewOFB(block, iv)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		}
	case 3:
		stream := cipher.NewCTR(block, iv)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		}
	default:
		panic(aesMode)
	}
}

func prepareInputs3(baseMsg []byte) func() []dudect.Input {
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, numberMeasurements)

		for i := 0; i < numberMeasurements; i++ {
			var randByte = make([]byte, 1)
			if n, err := io.ReadFull(rand.Reader, randByte); err != nil || n != 1 {
				panic(fmt.Sprintf("Randbit failed with Err: %v, n: %v", err, n))
			}
			class := int(randByte[0]) % 2
			// class-0 constant randomly picked input
			if class == 0 {
				temp := make([]byte, msgSize)
				copy(temp, baseMsg)
				inputs[i] = dudect.Input{Data: temp, Class: 0}
				continue
			}
			// class-1 varying random inputs
			var data = make([]byte, msgSize)
			if _, err := io.ReadFull(rand.Reader, data); err != nil {
				panic(err)
			}
			inputs[i] = dudect.Input{Data: data, Class: 1}
		}
		return inputs
	}
}
