package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit4(baseKey, anotherKey *rsa.PrivateKey, baseLabel []byte) func(uint8) func([]byte) {
	return func(class uint8) func([]byte) {
		var key *rsa.PrivateKey
		if class == 0 {
			// class-0 a constant randomly picked key
			key = rsaKeyDeepCopy(baseKey)
		} else if class == 1 {
			// class-1 another constant randomly picked key
			key = rsaKeyDeepCopy(anotherKey)
		} else {
			panic(class)
		}

		label := make([]byte, labelSize)
		copy(label, baseLabel)

		return func(plaintext []byte) {
			_, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, plaintext, baseLabel)
			if err != nil {
				panic(err)
			}
		}
	}
}

func prepareInputs4(baseMsg []byte) func() []dudect.Input {
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
