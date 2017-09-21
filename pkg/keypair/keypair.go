// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package keypair responses for private, public and address account logic
package keypair

import (
	"encoding/base32"
	"os"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Account keypair
type KeyPair struct {
	public  []byte
	private []byte
}

// GenAddress generates a new address for required chain on crypto random basis.
// It’s a run-time error for unknown chain.
func GenAddress(chainId byte) (string, error) {
	pair, err := NewKeyPair()
	if err != nil {
		os.Exit(-1)
	}
	return ToAddress(pair.public, chainId)
}

// NewKeyPair generates a public/private key pair using entropy from crypto rand.
func NewKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	return KeyPair{
		pub, priv[:32],
	}, err
}

// ToAddress converts public key to public account address.
// It’s a run-time error for unknown chain.
func ToAddress(pubKey []byte, chainId byte) (string, error) {
	if !IsValidChainId(chainId) {
		return "", core.ErrInvalidChain
	}
	h := sha3.SumKeccak256(pubKey)
	r := ripemd160.New()
	_, err := r.Write(h[:])
	if err != nil {
		return "", err
	}
	b := append([]byte{chainId}, r.Sum(nil)...)
	h = sha3.SumKeccak256(b)
	a := append(b, h[:4]...)
	return base32.StdEncoding.EncodeToString(a), nil
}

// IsValidChainId checks chain id for existence
func IsValidChainId(id byte) bool {
	for _, i := range []byte{core.MijinId, core.MainnetId, core.TestnetId} {
		if i == id {
			return true
		}
	}
	return false
}
