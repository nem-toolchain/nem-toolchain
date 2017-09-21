// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package keypair responses for private, public and address account subjects
package keypair

import (
	"encoding/base32"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Keypair is a bundle of a private/public key pair
type KeyPair struct {
	private []byte
	public  []byte
}

// NewKeyPair generates a private/public key pair using entropy from crypto rand.
func NewKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	return KeyPair{
		priv[:32], pub,
	}, err
}

// GenAddress generates a new address for required chain on crypto random basis.
func GenAddress(chain core.Chain) (string, error) {
	pair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	return ToAddress(pair.public, chain)
}

// ToAddress converts public key to public account address string representation.
func ToAddress(pubKey []byte, chain core.Chain) (string, error) {
	h := sha3.SumKeccak256(pubKey)
	r := ripemd160.New()
	_, err := r.Write(h[:])
	if err != nil {
		return "", err
	}
	b := append([]byte{chain.Id}, r.Sum(nil)...)
	h = sha3.SumKeccak256(b)
	a := append(b, h[:4]...)
	return base32.StdEncoding.EncodeToString(a), nil
}
