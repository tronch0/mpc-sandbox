package main

import (
	"crypto/rand"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func generateIdentity() crypto.PrivKey {
	sk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	return sk
}
