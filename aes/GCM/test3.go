package gcm_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

//const (
//	numberMeasurements = 1000000
//	msgSize = 1024 * aes.BlockSize
//)

func spawnInit3(baseKey, baseNonce []byte) func(uint8) func([]byte) {

	return func(_ uint8) func([]byte) {
		// constant randomly picked key
		key := make([]byte, keySize)
		copy(key, baseKey)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		// constant randomly picked nonce
		nonce := make([]byte, nonceSize)
		copy(nonce, baseNonce)

		aesgcm, err := cipher.NewGCMWithNonceSize(block, nonceSize)
		if err != nil {
			panic(err)
		}

		return func(plaintext []byte) {
			aesgcm.Seal(nil, nonce, plaintext, nil)
		}
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
