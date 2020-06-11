package main

import (
	"gocrypto_test/blockcipher_aes/CBC_CFB_OFB"
	"gocrypto_test/blockcipher_aes/GCM"
)

func main() {
	cbc_cfb_ofb_test.DoTest()
	gcm_test.DoTest()
}
