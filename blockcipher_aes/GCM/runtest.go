package gcm_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	dudect "github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements  = 1000000                // total measurements for test 1/3/5
	smallerMeasurements = 10000                  // total measurements for test 2/4/6
	numberKeys          = numberMeasurements / 2 // number of various keys used for test 1
	numberNonce         = numberMeasurements / 2 // number of various IVs used for test 5
	trialNum            = 3                      // total trail number for this group of test
	msgTrail            = 100                    // total number of different msgs used to for test 2/6, can set to 1 for debugging
)

func DoTest() {
	fmt.Println("Start testing AES-GCM.")
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
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	msg := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0):\n key = %s\n nonce = %s\nmsg = %s\n=========="+
		"=====================\n", hex.EncodeToString(key), hex.EncodeToString(nonce), hex.EncodeToString(msg))

	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	fmt.Printf("<%s Mode Test-1>\n", "GCM")
	dudect.Dudect(spawnInit1(key, nonce), prepareInputs1(msg), true)
	fmt.Println()

	// test2
	fmt.Println("|------------------Start Test-2------------------|")
	for j := 0; j < numSpecialKeyMode; j++ {
		fmt.Printf("<%s Mode Test-2.%d>\n", "GCM", j)
		for k := 0; k < msgTrail; k++ {
			if k == 0 {
				fmt.Printf("test against base msg\n")
				dudect.Dudect(spawnInit2(j, key, nonce), prepareInputs2(msg), true)
				continue
			}
			tempMsg := make([]byte, msgSize)
			if _, err := io.ReadFull(rand.Reader, tempMsg); err != nil {
				panic(err)
			}
			fmt.Printf("test against random msg = %s\n", hex.EncodeToString(tempMsg))
			dudect.Dudect(spawnInit2(j, key, nonce), prepareInputs2(tempMsg), true)
		}
	}
	fmt.Println()

	// test3
	fmt.Println("|------------------Start Test-3------------------|")
	fmt.Printf("<%s Mode Test-3>\n", "GCM")
	dudect.Dudect(spawnInit3(key, nonce), prepareInputs3(msg), false)
	fmt.Println()

	// test4
	fmt.Println("|------------------Start Test-4------------------|")
	for j := 0; j < numSpecialMsgMode; j++ {
		fmt.Printf("<%s Mode Test-4.%d>\n", "GCM", j)
		specialMsg, f := spawnInit4(j, key, nonce)
		dudect.Dudect(f, prepareInputs4(msg, specialMsg), false)
	}
	fmt.Println()

	// test5
	fmt.Println("|------------------Start Test-5------------------|")
	fmt.Printf("<%s Mode Test-5>\n", "GCM")
	dudect.Dudect(spawnInit5(key, nonce), prepareInputs5(msg), true)
	fmt.Println()

	// test6
	fmt.Println("|------------------Start Test-6------------------|")
	for j := 0; j < numSpecialNonceMode; j++ {
		fmt.Printf("<%s Mode Test-6.%d>\n", "GCM", j)
		for k := 0; k < msgTrail; k++ {
			if k == 0 {
				fmt.Printf("test against base msg\n")
				dudect.Dudect(spawnInit6(j, key, nonce), prepareInputs6(msg), true)
				continue
			}
			tempMsg := make([]byte, msgSize)
			if _, err := io.ReadFull(rand.Reader, tempMsg); err != nil {
				panic(err)
			}
			fmt.Printf("test against random msg = %s\n", hex.EncodeToString(tempMsg))
			dudect.Dudect(spawnInit6(j, key, nonce), prepareInputs6(tempMsg), true)
		}
	}
}
