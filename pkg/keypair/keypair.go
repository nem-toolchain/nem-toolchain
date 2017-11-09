// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package keypair responses for account's private/public crypto keys.
package keypair

import (
	"bytes"
	"io"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const (
	// PrivateBytes stores the private key length in bytes
	PrivateBytes = 32
	// PublicBytes stores the public key length in bytes
	PublicBytes = 32
)

// KeyPair is a private/public crypto key pair.
type KeyPair struct {
	Private []byte
	Public  []byte
}

// Gen generates a new private/public key pair using entropy from crypto rand.
func Gen() KeyPair {
	return FromSeed(nil)
}

// FromSeed generates a new private/public key pair using specified seed data
func FromSeed(seed []byte) KeyPair {
	var r io.Reader
	if seed != nil {
		r = bytes.NewReader(
			append(make([]byte, PrivateBytes-len(seed)), seed...))
	}

	pub, pr, err := ed25519.GenerateKey(r)
	if err != nil {
		panic("assert: ed25519 GenerateKey function internal error")
	}
	return KeyPair{pr[:PrivateBytes], pub}
}

// Address converts a key pair into corresponding address string representation.
func (pair KeyPair) Address(chain core.Chain) Address {
	h := sha3.SumKeccak256(pair.Public)

	r := ripemd160.New()
	_, err := r.Write(h[:])
	if err != nil {
		panic("assert: Ripemd160 hash function internal error")
	}

	b := append([]byte{chain.Id}, r.Sum(nil)...)

	h = sha3.SumKeccak256(b)
	a := append(b, h[:4]...)

	addr := Address{}
	copy(addr[:], a[:])
	return addr
}
