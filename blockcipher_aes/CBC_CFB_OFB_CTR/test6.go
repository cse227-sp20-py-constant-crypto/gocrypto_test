package cbc_cfb_ofb_ctr_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit6(aesMode, specialIVMode int, baseKey, baseIV []byte) func(uint8) func([]byte) {

	sIV := make([]byte, ivSize)
	switch specialIVMode {
	case 0:
		binary.BigEndian.PutUint32(sIV, 0)
	case 1:
		binary.BigEndian.PutUint32(sIV, 1)
	default:
		panic(specialIVMode)
	}

	return func(class uint8) func([]byte) {
		key := make([]byte, keySize)
		copy(key, baseKey)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		iv := make([]byte, ivSize)
		if class == 0 {
			// class-0 a constant randomly picked key
			copy(iv, baseIV)
		} else if class == 1 {
			// class-1 special iv
			copy(iv, sIV)
		} else {
			panic(class)
		}

		switch aesMode {
		case 0:
			mode := cipher.NewCBCEncrypter(block, iv)
			return func(plaintext []byte) {
				if len(plaintext)%aes.BlockSize != 0 {
					panic("plaintext is not a multiple of the block size")
				}
				ciphertext := make([]byte, len(plaintext))
				mode.CryptBlocks(ciphertext, plaintext)
			}
		case 1:
			stream := cipher.NewCFBEncrypter(block, iv)
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		case 2:
			stream := cipher.NewOFB(block, iv)
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		case 3:
			stream := cipher.NewCTR(block, iv)
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		default:
			panic(aesMode)
		}
	}
}

func prepareInputs6(baseMsg []byte) func() []dudect.Input {
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, smallerMeasurements)
		for i := 0; i < smallerMeasurements; i++ {
			var randByte = make([]byte, 1)
			if n, err := io.ReadFull(rand.Reader, randByte); err != nil || n != 1 {
				panic(fmt.Sprintf("Randbit failed with Err: %v, n: %v", err, n))
			}
			class := int(randByte[0]) % 2
			// class-0 and class-1 use constant randomly picked data
			data := make([]byte, msgSize)
			copy(data, baseMsg)
			if class == 0 {
				inputs[i] = dudect.Input{Data: data, Class: 0}
				continue
			}
			inputs[i] = dudect.Input{Data: data, Class: 1}
		}
		return inputs
	}
}
