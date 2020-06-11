package main

import (
	"gocrypto_test/blockcipher_aes/CBC_CFB_OFB_CTR"
	"gocrypto_test/blockcipher_aes/GCM"
)

func main() {
	cbc_cfb_ofb_ctr_test.DoTest()
	gcm_test.DoTest()
}
