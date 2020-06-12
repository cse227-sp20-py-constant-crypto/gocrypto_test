package ecdsa_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit3(curveMode int, baseKey *ecdsa.PrivateKey) func(uint8) func([]byte) {
	var keyArray [numberKeys]*ecdsa.PrivateKey
	var err error
	var cur elliptic.Curve
	switch curveMode {
	case 0:
		cur = elliptic.P256()
	case 1:
		cur = elliptic.P384()
	default:
		panic(curveMode)
	}
	for i := range keyArray {
		keyArray[i], err = ecdsa.GenerateKey(cur, rand.Reader)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("random key pool ready")
	counter := 0
	return func(class uint8) func([]byte) {
		var key *ecdsa.PrivateKey
		if class == 0 {
			// class-0 a constant randomly picked key
			key = ecdsaKeyDeepCopy(baseKey)
		} else if class == 1 {
			// class-1 varying keys (randomly prepicked)
			key = ecdsaKeyDeepCopy(keyArray[counter%numberKeys])
			counter++
		} else {
			panic(class)
		}

		return func(plaintext []byte) {
			hash := sha256.Sum256(plaintext)
			_, _, err := ecdsa.Sign(rand.Reader, key, hash[:])
			if err != nil {
				panic(err)
			}
		}
	}
}

func prepareInputs3(curveMode int, baseMsg []byte) func() []dudect.Input {
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
