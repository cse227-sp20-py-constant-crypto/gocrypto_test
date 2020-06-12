package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

//const (
//	numberMeasurements = 1000000
//	msgSize = 1024 * aes.BlockSize
//)

func spawnInit1(baseKey *rsa.PrivateKey, baseLabel []byte) func(uint8) func([]byte) {
	return func(class uint8) func([]byte) {
		var key *rsa.PrivateKey
		key = rsaKeyDeepCopy(baseKey)

		return func(plaintext []byte) {
			_, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, plaintext, baseLabel)
			if err != nil {
				panic(err)
			}
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
