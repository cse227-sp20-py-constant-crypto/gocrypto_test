package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	dudect "github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements  = 1000000                   // total measurements for test 1/3
	smallerMeasurements = 10000                     // total measurements for test 2/4
	numberKeys          = numberMeasurements / 1000 // number of various keys used for test 1
	trialNum            = 3                         // total trail number for this group of test
	keyCompareTrail     = 100                       // total number of different msgs used to for test 2/6, can set to 1 for debugging
)

func DoTest() {
	fmt.Println("Start testing RSA OAEP.")
	for i := 0; i < trialNum; i++ {
		fmt.Printf("<<<<<<<<<<<<<<<Trail %d>>>>>>>>>>>>>>>>\n", i+1)
		test()
	}
}

func test() {
	var key *rsa.PrivateKey
	var err error
	key, err = rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		panic(err)
	}
	label := make([]byte, labelSize)
	if _, err := io.ReadFull(rand.Reader, label); err != nil {
		panic(err)
	}
	msg := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0):\n key.N = %d\n key.e = %d\n label = %s\n msg = %s\n=========="+
		"=====================\n", &key.N, &key.E, hex.EncodeToString(label), hex.EncodeToString(msg))

	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	fmt.Printf("<%s Test-1>\n", "RSA OAEP")
	dudect.Dudect(spawnInit1(key, label), prepareInputs1(msg), false)
	fmt.Println()

	// test2
	fmt.Println("|------------------Start Test-2------------------|")
	for j := 0; j < numSpecialMsgMode; j++ {
		fmt.Printf("<%s Test-2.%d>\n", "RSA OAEP", j)
		specialMsg, f := spawnInit2(j, key, label)
		dudect.Dudect(f, prepareInputs2(msg, specialMsg), false)
		fmt.Println()
	}

	// test3
	fmt.Println("|------------------Start Test-3------------------|")
	fmt.Printf("<%s Mode Test-3>\n", "RSA OAEP")
	dudect.Dudect(spawnInit3(key, label), prepareInputs3(msg), true)
	fmt.Println()

	// test4
	fmt.Println("|------------------Start Test-4------------------|")
	for i := 0; i < keyCompareTrail; i++ {
		var key1, key2 *rsa.PrivateKey
		var err error
		key1, err = rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			panic(err)
		}
		key2, err = rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Randomly chosen two key:\n key1: N = %d, e = %d\n key2: N = %d, e = %d\n", &key1.N, &key1.E, &key2.N, &key2.E)
		dudect.Dudect(spawnInit4(key1, key2, label), prepareInputs4(msg), true)
		fmt.Println()
	}
}
