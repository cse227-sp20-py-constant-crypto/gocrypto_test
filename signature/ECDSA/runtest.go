package ecdsa_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	dudect "github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements  = 1000000                  // total measurements for test 1/3
	smallerMeasurements = 100000                   // total measurements for test 2/4
	numberKeys          = numberMeasurements / 100 // number of various keys used for test 1
	trialNum            = 3                        // total trail number for this group of test
)

func DoTest() {
	fmt.Println("Start testing ECDSA.")
	for i := 0; i < trialNum; i++ {
		fmt.Printf("<<<<<<<<<<<<<<<Trail %d>>>>>>>>>>>>>>>>\n", i+1)
		test()
	}
}

func test() {
	for curveMode := 0; curveMode < 1; curveMode++ {
		var cur elliptic.Curve
		switch curveMode {
		case 0:
			cur = elliptic.P256()
			fmt.Println("<NIST-P256>")
		case 1:
			cur = elliptic.P384()
			fmt.Println("<NIST-P384>")
		default:
			panic(curveMode)
		}

		var key *ecdsa.PrivateKey
		var err error
		key, err = ecdsa.GenerateKey(cur, rand.Reader)
		if err != nil {
			panic(err)
		}
		msg := make([]byte, msgSize)
		if _, err := io.ReadFull(rand.Reader, msg); err != nil {
			panic(err)
		}
		fmt.Printf("Randomly chosen baseline parameters (class-0):\n key.X = %d\n key.Y = %d\n msg = %s\n=========="+
			"=====================\n", &key.PublicKey.X, &key.PublicKey.Y, hex.EncodeToString(msg))

		// test1
		fmt.Println("|------------------Start Test-1------------------|")
		dudect.Dudect(spawnInit1(key), prepareInputs1(curveMode, msg), false)
		fmt.Println()

		// test3
		fmt.Println("|------------------Start Test-3------------------|")
		dudect.Dudect(spawnInit3(curveMode, key), prepareInputs3(curveMode, msg), true)
		fmt.Println()
	}
}
