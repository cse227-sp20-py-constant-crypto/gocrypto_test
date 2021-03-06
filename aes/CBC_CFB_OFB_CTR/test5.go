package cbc_cfb_ofb_ctr_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit5(aesMode int, baseKey, baseIV []byte) func(uint8) func([]byte) {
	var ivArray [numberIVs][]byte
	for i := range ivArray {
		ivArray[i] = make([]byte, ivSize)
		if _, err := io.ReadFull(rand.Reader, ivArray[i]); err != nil {
			panic(err)
		}
	}
	counter := 0
	return func(class uint8) func([]byte) {
		// constant key for either class
		key := make([]byte, keySize)
		copy(key, baseKey)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		iv := make([]byte, ivSize)
		if class == 0 {
			// class-1 constant randomly picked iv
			copy(iv, baseIV)
		} else if class == 1 {
			// class-1 varying iv
			copy(iv, ivArray[counter%numberIVs])
			counter++
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

func prepareInputs5(baseMsg []byte) func() []dudect.Input {
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, numberMeasurements)
		for i := 0; i < numberMeasurements; i++ {
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
