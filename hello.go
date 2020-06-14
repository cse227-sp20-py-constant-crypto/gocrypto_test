package main

import (
	"gocrypto_test/aes/CBC_CFB_OFB_CTR"
	"gocrypto_test/aes/GCM"
	"gocrypto_test/compare"
	"gocrypto_test/hash"
	"gocrypto_test/mac"
	"gocrypto_test/rsa"
	"gocrypto_test/signature/ECDSA"
	"gocrypto_test/streamcipher/salsa20"
)

func main() {
	cbc_cfb_ofb_ctr_test.DoTest()
	gcm_test.DoTest()
	hash_test.DoTest()
	mac_test.DoTest()
	salsa20_test.DoTest()
	rsa_test.DoTest()
	cmp_test.DoTest()
	ecdsa_test.DoTest()
}
