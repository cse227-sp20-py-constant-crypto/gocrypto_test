package cbc_cfb_ofb_test

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	dudect "github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements  = 100000                 // total measurements for test 1/3/5
	smallerMeasurements = 10000                  // total measurements for test 2/4/6
	numberKeys          = numberMeasurements / 2 // number of various keys used for test 1
	numberIVs           = numberMeasurements / 2 // number of various IVs used for test 5
	trialNum            = 3                      // total trail number for this group of test
)

type AESMode int

const (
	CBC AESMode = iota
	CFB
	OFB
)

func (aesMode AESMode) String() string {
	names := [...]string{
		"AES-CBC",
		"AES-CFB",
		"AES-OFB",
	}
	if aesMode != CBC && aesMode != CFB && aesMode != OFB {
		return "N.A."
	}
	return names[aesMode]
}

func DoTest() {
	fmt.Println("Start testing AES CBC/CFB/OFB mode.")
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
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	msg := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0): key = %s, iv = %s,\nmsg = %s\n", hex.EncodeToString(key), hex.EncodeToString(iv), hex.EncodeToString(msg))
	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	for i := 0; i < 3; i++ {
		fmt.Printf("<%s Mode Test-1>\n", AESMode(i))
		dudect.Dudect(spawnInit1(i, key, iv), prepareInputs1(msg), true)
	}
	fmt.Println()
	// test2
	fmt.Println("|------------------Start Test-2------------------|")
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			fmt.Printf("<%s Mode Test-2.%d>\n", AESMode(i), j)
			dudect.Dudect(spawnInit2(i, j, key, iv), prepareInputs2(msg), true)
		}
	}
	fmt.Println()
	// test3
	fmt.Println("|------------------Start Test-3------------------|")
	for i := 0; i < 3; i++ {
		fmt.Printf("<%s Mode Test-3>\n", AESMode(i))
		dudect.Dudect(spawnInit3(i, key, iv), prepareInputs3(msg), false)
	}
	fmt.Println()
	// test4
	fmt.Println("|------------------Start Test-4------------------|")
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			fmt.Printf("<%s Mode Test-4.%d>\n", AESMode(i), j)
			specialMsg, f := spawnInit4(i, j, key, iv)
			dudect.Dudect(f, prepareInputs4(msg, specialMsg), false)
		}
	}
	fmt.Println()
	// test5
	fmt.Println("|------------------Start Test-5------------------|")
	for i := 0; i < 3; i++ {
		fmt.Printf("<%s Mode Test-5>\n", AESMode(i))
		dudect.Dudect(spawnInit5(i, key, iv), prepareInputs5(msg), true)
	}
	fmt.Println()
	// test6
	fmt.Println("|------------------Start Test-6------------------|")
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			fmt.Printf("<%s Mode Test-6.%d>\n", AESMode(i), j)
			dudect.Dudect(spawnInit6(i, j, key, iv), prepareInputs6(msg), true)
		}
	}
}
