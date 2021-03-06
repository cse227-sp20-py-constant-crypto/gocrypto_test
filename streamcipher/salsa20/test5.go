package salsa20_test

import (
	"crypto/rand"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/salsa20"
	"io"
)

func spawnInit5(baseKey, baseNonce []byte) func(uint8) func([]byte) {
	var nonceArray [numberNonce][]byte
	for i := range nonceArray {
		nonceArray[i] = make([]byte, nonceSize)
		if _, err := io.ReadFull(rand.Reader, nonceArray[i]); err != nil {
			panic(err)
		}
	}
	counter := 0

	return func(class uint8) func([]byte) {
		// constant key for either class
		var aKey [32]byte
		key := aKey[:]
		copy(key, baseKey)

		nonce := make([]byte, nonceSize)
		if class == 0 {
			// class-1 constant randomly picked nonce
			copy(nonce, baseNonce)
		} else if class == 1 {
			// class-1 varying nonce
			copy(nonce, nonceArray[counter%numberNonce])
			counter++
		} else {
			panic(class)
		}

		return func(plaintext []byte) {
			salsa20.XORKeyStream(plaintext, plaintext, nonce, &aKey)
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
