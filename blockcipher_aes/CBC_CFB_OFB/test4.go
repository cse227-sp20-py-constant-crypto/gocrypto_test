package cbc_cfb_ofb_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit4(aesMode, specialMsgMode int, baseKey, baseIV []byte) ([]byte, func(uint8) func([]byte)) {

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
		plaintext := make([]byte, msgSize)
		switch specialMsgMode {
		case 0:
			// plaintext as 0
		case 1:
			// plaintext as 1
			binary.BigEndian.PutUint32(plaintext, 1)
		case 2:
			// plaintext which encrypts to 0
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 0)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			decryptMode := cipher.NewCBCDecrypter(block, iv)
			decryptMode.CryptBlocks(plaintext, ctxt)
		case 3:
			// plaintext which encrypts to 1
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 1)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			decryptMode := cipher.NewCBCDecrypter(block, iv)
			decryptMode.CryptBlocks(plaintext, ctxt)
		default:
			panic(specialMsgMode)
		}
		return plaintext, func(_ uint8) func([]byte) {
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
		plaintext := make([]byte, msgSize)
		switch specialMsgMode {
		case 0:
			// plaintext as 0
		case 1:
			// plaintext as 1
			binary.BigEndian.PutUint32(plaintext, 1)
		case 2:
			// plaintext which encrypts to 0
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 0)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			decryptStream := cipher.NewCFBDecrypter(block, iv)
			decryptStream.XORKeyStream(plaintext, ctxt)
		case 3:
			// plaintext which encrypts to 1
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 1)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			decryptStream := cipher.NewCFBDecrypter(block, iv)
			decryptStream.XORKeyStream(plaintext, ctxt)
		default:
			panic(specialMsgMode)
		}
		return plaintext, func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		}
	case 2:
		stream := cipher.NewOFB(block, iv)
		plaintext := make([]byte, msgSize)
		switch specialMsgMode {
		case 0:
			// plaintext as 0
		case 1:
			// plaintext as 1
			binary.BigEndian.PutUint32(plaintext, 1)
		case 2:
			// plaintext which encrypts to 0
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 0)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			stream.XORKeyStream(plaintext, ctxt)
		case 3:
			// plaintext which encrypts to 1
			ctxt := make([]byte, msgSize)
			binary.BigEndian.PutUint32(ctxt, 1)
			if len(ctxt)%aes.BlockSize != 0 {
				panic("ciphertext is not a multiple of the block size")
			}
			stream.XORKeyStream(plaintext, ctxt)
		default:
			panic(specialMsgMode)
		}
		return plaintext, func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				ciphertext := make([]byte, len(plaintext))
				stream.XORKeyStream(ciphertext, plaintext)
			}
		}
	default:
		panic(aesMode)
	}
}

func prepareInputs4(baseMsg, specialMsg []byte) func() []dudect.Input {
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, smallerMeasurements)

		for i := 0; i < smallerMeasurements; i++ {
			var randByte = make([]byte, 1)
			if n, err := io.ReadFull(rand.Reader, randByte); err != nil || n != 1 {
				panic(fmt.Sprintf("Randbit failed with Err: %v, n: %v", err, n))
			}
			class := int(randByte[0]) % 2
			// class-0 constant randomly picked input
			if class == 0 {
				var data = make([]byte, msgSize)
				copy(data, baseMsg)
				inputs[i] = dudect.Input{Data: data, Class: 0}
				continue
			}
			// class-1 special input
			var data = make([]byte, msgSize)
			copy(data, specialMsg)
			inputs[i] = dudect.Input{Data: data, Class: 1}
		}
		return inputs
	}
}
