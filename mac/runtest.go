package mac_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements  = 1000000                // total measurements for test 1/3/5
	smallerMeasurements = 10000                  // total measurements for test 2/4/6
	numberKeys          = numberMeasurements / 2 // number of various keys used for test 1
	trialNum            = 3                      // total trail number for this group of test
	msgTrail            = 100                    // total number of different msgs used to for test 2/6, can set to 1 for debugging
	numHMACMode         = int(NUM)
)

type HMACMode int

const (
	hmacSHA256 HMACMode = iota
	hmacSHA3256
	Poly1305
	NUM
)

func (hmacMode HMACMode) String() string {
	names := [...]string{
		"HMAC-SHA256",
		"HMAC-SHA3256",
		"Poly1305",
	}
	if hmacMode != hmacSHA256 && hmacMode != hmacSHA3256 && hmacMode != Poly1305 {
		return "N.A."
	}
	return names[hmacMode]
}

func DoTest() {
	fmt.Println("Start testing HMAC.")
	for i := 0; i < trialNum; i++ {
		fmt.Printf("<<<<<<<<<<<<<<<Trail %d>>>>>>>>>>>>>>>>\n", i+1)
		test()
	}
}

func test() {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	msg := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0):\n key = %s\nmsg = %s\n", hex.EncodeToString(key), hex.EncodeToString(msg))

	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	for i := 2; i < numHMACMode; i++ {
		fmt.Printf("<%s Test-1>\n", HMACMode(i))
		dudect.Dudect(spawnInit1(i, key), prepareInputs1(msg), true)
	}
	fmt.Println()

	// test2
	fmt.Println("|------------------Start Test-2------------------|")
	for i := 0; i < numHMACMode; i++ {
		for j := 0; j < numSpecialKeyMode; j++ {
			fmt.Printf("<%s Mode Test-2.%d>\n", HMACMode(i), j)
			for k := 0; k < msgTrail; k++ {
				if k == 0 {
					fmt.Printf("test against msg\n")
					dudect.Dudect(spawnInit2(i, j, key), prepareInputs2(msg), true)
					fmt.Println()
					continue
				}
				tempMsg := make([]byte, msgSize)
				if _, err := io.ReadFull(rand.Reader, tempMsg); err != nil {
					panic(err)
				}
				fmt.Printf("test against random msg = %s\n", hex.EncodeToString(tempMsg))
				dudect.Dudect(spawnInit2(i, j, key), prepareInputs2(tempMsg), true)
				fmt.Println()
			}
		}
	}

	// test3
	fmt.Println("|------------------Start Test-3------------------|")
	for i := 0; i < numHMACMode; i++ {
		fmt.Printf("<%s Test-3>\n", HMACMode(i))
		if HMACMode(i) != Poly1305 {
			dudect.Dudect(spawnInit3(i, key), prepareInputs3(msg), false)
			continue
		}
		fmt.Println("?")
		dudect.Dudect(spawnInit3(i, key), prepareInputs3(msg), true)
	}
	fmt.Println()

	// test4
	fmt.Println("|------------------Start Test-4------------------|")
	for i := 0; i < numHMACMode; i++ {
		for j := 0; j < numSpecialMsgMode; j++ {
			fmt.Printf("<%s Mode Test-4.%d>\n", HMACMode(i), j)
			specialMsg, f := spawnInit4(i, j, key)
			if HMACMode(i) != Poly1305 {
				dudect.Dudect(f, prepareInputs4(msg, specialMsg), false)
				continue
			}
			dudect.Dudect(f, prepareInputs4(msg, specialMsg), true)
		}
	}
	fmt.Println()
}
