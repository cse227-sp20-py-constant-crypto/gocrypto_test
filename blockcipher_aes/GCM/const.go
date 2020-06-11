package gcm_test

import "crypto/aes"

const (
	msgSize   = 1024 * aes.BlockSize
	keySize   = 32
	nonceSize = 16
)
