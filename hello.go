package main

import (
	"gocrypto_test/blockcipher_aes/CBC_CFB_OFB_CTR"
	"gocrypto_test/blockcipher_aes/GCM"
	"gocrypto_test/hash"
	"gocrypto_test/mac"
	"gocrypto_test/steamcipher/salsa20"
)

func main() {
	cbc_cfb_ofb_ctr_test.DoTest()
	gcm_test.DoTest()
	hash_test.DoTest()
	mac_test.DoTest()
	salsa20_test.DoTest()
}
