package cbc_cfb_ofb_ctr_test

import "crypto/aes"

const (
	msgSize           = 1024 * aes.BlockSize
	keySize           = 32
	ivSize            = aes.BlockSize
	numSpecialKeyMode = 4
	numSpecialIVMode  = 2
	numSpecialMsgMode = 4
)
