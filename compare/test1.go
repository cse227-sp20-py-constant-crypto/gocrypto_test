package cmp_test

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

func spawnInit1(cmpMode int, baseL []byte) func(uint8) func([]byte) {
	switch cmpMode {
	case 0:
		return func(_ uint8) func([]byte) {
			l := make([]byte, msgSize)
			copy(l, baseL)
			return func(r []byte) {
				subtle.ConstantTimeCompare(l, r)
			}
		}
	default:
		panic(fmt.Sprintf("cmpMode %d not within [0-%d]", cmpMode, numCMPMode-1))
	}
}

func prepareInputs1(baseR []byte) func() []dudect.Input {
	return func() []dudect.Input {
		var inputs = make([]dudect.Input, numberMeasurements)

		for i := 0; i < numberMeasurements; i++ {
			var randByte = make([]byte, 1)
			if n, err := io.ReadFull(rand.Reader, randByte); err != nil || n != 1 {
				panic(fmt.Sprintf("Randbit failed with Err: %v, n: %v", err, n))
			}
			class := int(randByte[0]) % 2
			// class-0 baseR
			if class == 0 {
				data := make([]byte, msgSize)
				copy(data, baseR)
				inputs[i] = dudect.Input{Data: data, Class: 0}
				continue
			}
			// class-1 random R of same length
			data := make([]byte, msgSize)
			if _, err := io.ReadFull(rand.Reader, data); err != nil {
				panic(err)
			}
			inputs[i] = dudect.Input{Data: data, Class: 1}
		}
		return inputs
	}
}
