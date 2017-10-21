// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package keypair responses for account's private/public crypto keys.
package keypair

import (
	"encoding/base32"

	"strings"

	"regexp"

	"errors"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const (
	// AddressBytes stores the address length
	AddressBytes = 25
	// AddressLength stores the address string representation length
	AddressLength = 40
)

// Address is a readable string representation for a public key.
type Address [AddressBytes]byte

// KeyPair is a private/public crypto key pair.
type KeyPair struct {
	Private []byte
	Public  []byte
}

// Gen generates a new private/public key pair using entropy from crypto rand.
func Gen() KeyPair {
	pub, pr, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("assert: ed25519 GenerateKey function internal error")
	}
	return KeyPair{pr[:32], pub}
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

// ParseAddress constructs an instance of `Address` from given base32 string representation
func ParseAddress(str string) (Address, error) {
	var addr Address
	str = strings.Replace(str, "-", "", -1)
	if !core.IsChainPrefix(str) {
		return addr, errors.New("unknown chain")
	}
	b, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return addr, errors.New("can't decode address string")
	}
	copy(addr[:], b)
	return addr, nil
}

// PrettyString returns pretty formatted address with separators ('-').
func (addr Address) PrettyString() string {
	str := addr.String()
	els := regexp.MustCompile(`.{6}`).FindAllString(str, -1)
	return strings.Join(append(els, str[36:]), "-")
}

func (addr Address) String() string {
	return base32.StdEncoding.EncodeToString(addr[:])
}
