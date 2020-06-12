package cmp_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	dudect "github.com/Reapor-Yurnero/godudect"
	"io"
)

const (
	numberMeasurements = 1000000 // total measurements for test 1/3/5
	trialNum           = 3       // total trail number for this group of test
	numCMPMode         = int(NUM)
)

type HASHMode int

const (
	subtleCMP HASHMode = iota
	NUM
)

func (aesMode HASHMode) String() string {
	names := [...]string{
		"subtle.ConstantTimeCompare",
	}
	if aesMode != subtleCMP {
		return "N.A."
	}
	return names[aesMode]
}

func DoTest() {
	fmt.Println("Start testing hash functions.")
	for i := 0; i < trialNum; i++ {
		fmt.Printf("<<<<<<<<<<<<<<<Trail %d>>>>>>>>>>>>>>>>\n", i+1)
		test()
	}
}

func test() {
	baseL := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, baseL); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0):\nmsg = %s\n", hex.EncodeToString(baseL))

	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	for i := 0; i < numCMPMode; i++ {
		fmt.Printf("<%s Test-1>\n", HASHMode(i))
		dudect.Dudect(spawnInit1(i, baseL), prepareInputs1(baseL), true) // need to be true to make sure for each measurement, the msg is a fresh object (not already in the mem)
	}
	fmt.Println()

	// test2
	fmt.Println("|------------------Start Test-2------------------|")
	for i := 0; i < numCMPMode; i++ {
		fmt.Printf("<%s Test-2>\n", HASHMode(i))
		dudect.Dudect(spawnInit2(i, baseL), prepareInputs2(), true) // need to be true to make sure for each measurement, the msg is a fresh object (not already in the mem)
	}
	fmt.Println()
}
