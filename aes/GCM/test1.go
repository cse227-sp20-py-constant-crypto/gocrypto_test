package gcm_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit1(baseKey, baseNonce []byte) func(uint8) func([]byte) {
	var keyArray [numberKeys][]byte
	for i := range keyArray {
		keyArray[i] = make([]byte, keySize)
		if _, err := io.ReadFull(rand.Reader, keyArray[i]); err != nil {
			panic(err)
		}
	}
	counter := 0
	return func(class uint8) func([]byte) {
		key := make([]byte, keySize)
		if class == 0 {
			// class-0 a constant randomly picked key
			copy(key, baseKey)
		} else if class == 1 {
			// class-1 varying keys (randomly prepicked)
			copy(key, keyArray[counter%numberKeys])
			counter++
		} else {
			panic(class)
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

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

func prepareInputs1(baseMsg []byte) func() []dudect.Input {
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
