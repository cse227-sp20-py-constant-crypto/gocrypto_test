package ecdsa_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit1(baseKey *ecdsa.PrivateKey) func(uint8) func([]byte) {
	return func(class uint8) func([]byte) {
		return func(plaintext []byte) {
			hash := sha256.Sum256(plaintext)
			_, _, err := ecdsa.Sign(rand.Reader, baseKey, hash[:])
			if err != nil {
				panic(err)
			}
		}
	}
}

func prepareInputs1(curveMode int, baseMsg []byte) func() []dudect.Input {
	var inputSize int
	switch curveMode {
	case 0:
		inputSize = numberMeasurements
	case 1:
		inputSize = smallerMeasurements
	default:
		panic(curveMode)
	}
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, inputSize)

		for i := 0; i < inputSize; i++ {
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
