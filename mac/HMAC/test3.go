package hmac_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/sha3"
	"io"
)

func spawnInit3(hmacMode int, baseKey []byte) func(uint8) func([]byte) {
	// constant randomly picked key
	key := make([]byte, keySize)
	copy(key, baseKey)

	switch hmacMode {
	case 0:
		mac := hmac.New(sha256.New, key)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		}
	case 1:
		mac := hmac.New(sha3.New256, key)
		return func(_ uint8) func([]byte) {
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		}
	default:
		panic(hmacMode)
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
