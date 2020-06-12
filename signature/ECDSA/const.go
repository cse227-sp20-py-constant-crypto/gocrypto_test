package ecdsa_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	msgSize = 32 // a typical 32byte symmetric key
)

func ecdsaKeyDeepCopy(src *ecdsa.PrivateKey) *ecdsa.PrivateKey {
	x509Encoded, _ := x509.MarshalECPrivateKey(src)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	block, _ := pem.Decode(pemEncoded)
	dst, _ := x509.ParseECPrivateKey(block.Bytes)
	return dst
}
