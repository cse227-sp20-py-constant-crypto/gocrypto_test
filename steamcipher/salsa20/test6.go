package salsa20_test

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/salsa20"
	"io"
)

func spawnInit6(specialNonceMode int, baseKey, baseNonce []byte) func(uint8) func([]byte) {
	sNonce := make([]byte, nonceSize)
	switch specialNonceMode {
	case 0:
		binary.BigEndian.PutUint32(sNonce, 0)
	case 1:
		binary.BigEndian.PutUint32(sNonce, 1)
	default:
		panic(fmt.Sprintf("specialNonceMode %d not within [0-%d]", specialNonceMode, numSpecialNonceMode-1))
	}

	return func(class uint8) func([]byte) {
		// constant key for either class
		var aKey [32]byte
		key := aKey[:]
		copy(key, baseKey)

		nonce := make([]byte, nonceSize)
		if class == 0 {
			// class-0 a constant randomly picked key
			copy(nonce, baseNonce)
		} else if class == 1 {
			// class-1 special nonce
			copy(nonce, sNonce)
		} else {
			panic(class)
		}

		return func(plaintext []byte) {
			salsa20.XORKeyStream(plaintext, plaintext, nonce, &aKey)
		}
	}
}

func prepareInputs6(baseMsg []byte) func() []dudect.Input {
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
