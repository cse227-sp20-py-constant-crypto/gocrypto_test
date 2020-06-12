package rsa_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	msgSize           = 32 // a typical 32byte symmetric key
	keySize           = 1024
	labelSize         = 16
	numSpecialMsgMode = 4
)

func rsaKeyDeepCopy(src *rsa.PrivateKey) *rsa.PrivateKey {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(src)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	block, _ := pem.Decode(privkeyPem)
	dst, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return dst
}
