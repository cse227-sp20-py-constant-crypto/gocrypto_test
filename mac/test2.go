package mac_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/poly1305"
	"golang.org/x/crypto/sha3"
	"io"
)

func spawnInit2(macMode, specialKeyMode int, baseKey []byte) func(uint8) func([]byte) {
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
		var aKey [keySize]byte
		//key := make([]byte, keySize)
		key := aKey[:]
		if class == 0 {
			// class-0 a constant randomly picked key
			copy(key, baseKey)
		} else if class == 1 {
			// class-1 special key (randomly prepicked)
			copy(key, sKey)
		} else {
			panic(class)
		}

		switch macMode {
		case 0:
			mac := hmac.New(sha256.New, key)
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		case 1:
			mac := hmac.New(sha3.New256, key)
			return func(plaintext []byte) {
				mac.Write(plaintext)
				mac.Sum(nil)
			}
		case 2:
			mac := poly1305.New(&aKey)
			return func(plaintext []byte) {
				_, _ = mac.Write(plaintext)
				mac.Sum(nil)
			}
		default:
			panic(macMode)
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
