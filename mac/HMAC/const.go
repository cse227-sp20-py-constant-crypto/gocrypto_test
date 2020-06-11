package hmac_test

import "crypto/aes"

const (
	msgSize           = 1024 * aes.BlockSize
	keySize           = 32
	numSpecialKeyMode = 4
	numSpecialMsgMode = 2
)
