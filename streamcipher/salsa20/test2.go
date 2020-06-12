package salsa20_test

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/salsa20"
	"io"
)

//const (
//	smallerMeasurements = 100000
//	msgSize = 1024 * aes.BlockSize
//)

func spawnInit2(specialKeyMode int, baseKey, baseNonce []byte) func(uint8) func([]byte) {
	sKey := make([]byte, keySize)
	switch specialKeyMode {
	case 0:
		binary.BigEndian.PutUint32(sKey, 0)
	case 1:
		binary.BigEndian.PutUint32(sKey, 1)
	case 2:
		binary.BigEndian.PutUint32(sKey, 2)
	case 3:
		binary.BigEndian.PutUint32(sKey, 3)
	default:
		panic(fmt.Sprintf("specialKeyMode %d not within [0-%d]", specialKeyMode, numSpecialKeyMode-1))
	}

	return func(class uint8) func([]byte) {
		var aKey [32]byte
		key := aKey[:]

		nonce := make([]byte, nonceSize)
		copy(nonce, baseNonce)

		if class == 0 {
			// class-0 a constant randomly picked key
			copy(key, baseKey)
		} else if class == 1 {
			// class-1 special key (randomly prepicked)
			copy(key, sKey)
		} else {
			panic(class)
		}

		return func(plaintext []byte) {
			salsa20.XORKeyStream(plaintext, plaintext, nonce, &aKey)
		}
	}
}

func prepareInputs2(baseMsg []byte) func() []dudect.Input {
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
