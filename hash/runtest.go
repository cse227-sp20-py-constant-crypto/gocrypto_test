package hash_test

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
)

type HASHMode int

const (
	SHA256 HASHMode = iota
	SHA3256
)

func (aesMode HASHMode) String() string {
	names := [...]string{
		"SHA256",
		"SHA3-256",
	}
	if aesMode != SHA256 && aesMode != SHA3256 {
		return "N.A."
	}
	return names[aesMode]
}

func DoTest() {
	fmt.Println("Start testing AES-CBC/CFB/OFB.")
	for i := 0; i < trialNum; i++ {
		fmt.Printf("<<<<<<<<<<<<<<<Trail %d>>>>>>>>>>>>>>>>\n", i+1)
		test()
	}
}

func test() {
	msg := make([]byte, msgSize)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	fmt.Printf("Randomly chosen baseline parameters (class-0):\nmsg = %s\n", hex.EncodeToString(msg))

	// test1
	fmt.Println("|------------------Start Test-1------------------|")
	for i := 0; i < 2; i++ {
		fmt.Printf("<%s Mode Test-1>\n", HASHMode(i))
		dudect.Dudect(spawnInit1(i), prepareInputs1(msg), true)
	}
	fmt.Println()

}
