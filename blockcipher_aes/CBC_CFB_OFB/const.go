package cbc_cfb_ofb_test

import "crypto/aes"

const (
	msgSize = 1024 * aes.BlockSize
	keySize = 32
	ivSize  = aes.BlockSize
)
