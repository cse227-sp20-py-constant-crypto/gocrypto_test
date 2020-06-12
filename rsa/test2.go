package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit2(specialMsgMode int, baseKey *rsa.PrivateKey, baseLabel []byte) ([]byte, func(uint8) func([]byte)) {
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
		plaintext, _ = rsa.DecryptOAEP(sha256.New(), rand.Reader, baseKey, ctxt, baseLabel)
	case 3:
		// plaintext which encrypts to 1
		ctxt := make([]byte, msgSize)
		binary.BigEndian.PutUint32(ctxt, 1)
		plaintext, _ = rsa.DecryptOAEP(sha256.New(), rand.Reader, baseKey, ctxt, baseLabel)
	default:
		panic(fmt.Sprintf("specialMsgMode %d not within [0-%d]", specialMsgMode, numSpecialMsgMode-1))
	}

	return plaintext, func(_ uint8) func([]byte) {
		key := rsaKeyDeepCopy(baseKey)

		// constant randomly picked label
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

func prepareInputs2(baseMsg, specialMsg []byte) func() []dudect.Input {
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
