package hash_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
)

func spawnInit1(shaMode int) func(uint8) func([]byte) {
	var h hash.Hash
	switch shaMode {
	case 0:
		h = sha256.New()
	case 1:
		h = sha3.New256()
	default:
		panic(fmt.Sprintf("shaMode %d not within [0-%d]", shaMode, numSHAMode-1))
	}
	return func(_ uint8) func([]byte) {
		return func(plaintext []byte) {
			h.Write(plaintext)
			h.Sum(nil)
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
