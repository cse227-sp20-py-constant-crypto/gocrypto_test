package salsa20_test

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/salsa20"
	"io"
)

func spawnInit4(specialMsgMode int, baseKey, baseNonce []byte) ([]byte, func(uint8) func([]byte)) {

	// constant randomly picked key
	var aKey [32]byte
	key := aKey[:]
	copy(key, baseKey)

	// constant randomly picked nonce
	nonce := make([]byte, nonceSize)
	copy(nonce, baseNonce)

	plaintext := make([]byte, msgSize)
	switch specialMsgMode {
	case 0:
		// plaintext as 0
	case 1:
		// plaintext as 1
		binary.BigEndian.PutUint32(plaintext, 1)
	case 2:
		// plaintext which encrypts to 0
		ctxt := make([]byte, msgSize)
		binary.BigEndian.PutUint32(ctxt, 0)
		salsa20.XORKeyStream(plaintext, ctxt, nonce, &aKey)
	case 3:
		// plaintext which encrypts to 1
		ctxt := make([]byte, msgSize)
		binary.BigEndian.PutUint32(ctxt, 1)
		salsa20.XORKeyStream(plaintext, ctxt, nonce, &aKey)
	default:
		panic(fmt.Sprintf("specialMsgMode %d not within [0-%d]", specialMsgMode, numSpecialMsgMode-1))
	}

	return plaintext, func(_ uint8) func([]byte) {
		key := aKey[:]
		copy(key, baseKey)

		// constant randomly picked nonce
		nonce := make([]byte, nonceSize)
		copy(nonce, baseNonce)

		return func(plaintext []byte) {
			salsa20.XORKeyStream(plaintext, plaintext, nonce, &aKey)
		}
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
