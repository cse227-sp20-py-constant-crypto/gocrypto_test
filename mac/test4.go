package mac_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/poly1305"
	"golang.org/x/crypto/sha3"
	"io"
)

func spawnInit4(macMode, specialMsgMode int, baseKey []byte) ([]byte, func(uint8) func([]byte)) {

	// constant randomly picked key
	var aKey [keySize]byte
	//key := make([]byte, keySize)
	key := aKey[:]
	copy(key, baseKey)

	plaintext := make([]byte, msgSize)
	switch specialMsgMode {
	case 0:
	case 1:
		binary.BigEndian.PutUint32(plaintext, 1)
	default:
		panic(fmt.Sprintf("specialMsgMode %d not within [0-%d]", specialMsgMode, numSpecialMsgMode-1))
	}

	switch macMode {
	case 0:
		mac := hmac.New(sha256.New, key)
		return plaintext, func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		}
	case 1:
		mac := hmac.New(sha3.New256, key)
		return plaintext, func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		}
	case 2:
		return plaintext, func(_ uint8) func([]byte) {
			mac := poly1305.New(&aKey)
			return func(plaintext []byte) {
				_, _ = mac.Write(plaintext)
				mac.Sum(nil)
			}
		}
	default:
		panic(macMode)
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
