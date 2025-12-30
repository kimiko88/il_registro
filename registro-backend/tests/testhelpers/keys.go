package testhelpers

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return priv, &priv.PublicKey
}
