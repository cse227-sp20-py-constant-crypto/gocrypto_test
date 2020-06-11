package gcm_test

import (
	"crypto/aes"
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
	nonce := make([]byte, aes.BlockSize)
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
	for j := 0; j < 4; j++ {
		fmt.Printf("<%s Mode Test-2.%d>\n", "GCM", j)
		dudect.Dudect(spawnInit2(j, key, nonce), prepareInputs2(msg), true)
	}
	fmt.Println()

	// test3
	fmt.Println("|------------------Start Test-3------------------|")
	fmt.Printf("<%s Mode Test-3>\n", "GCM")
	dudect.Dudect(spawnInit3(key, nonce), prepareInputs3(msg), false)
	fmt.Println()

	// test4
	fmt.Println("|------------------Start Test-4------------------|")
	for j := 0; j < 2; j++ {
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
	for j := 0; j < 2; j++ {
		fmt.Printf("<%s Mode Test-6.%d>\n", "GCM", j)
		dudect.Dudect(spawnInit6(j, key, nonce), prepareInputs6(msg), true)
	}
}
