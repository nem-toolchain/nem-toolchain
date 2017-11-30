// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package account

import (
	"crypto/aes"
	"crypto/rand"
	"io"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// NewAccount create new account for selected chain.
func NewAccount(ch core.Chain) *Account {
	acc := new(Account)
	acc.Label = "Primary"
	acc.Algo = "pass:enc"
	acc.Brain = false
	acc.Network = ch

	return acc
}

// Account used to encrypt private key.
type Account struct {
	Algo      string
	Label     string
	Encrypted []byte
	Iv        []byte
	Child     []byte
	Address   keypair.Address
	Network   core.Chain
	Brain     bool
}

// Encrypt KeyPair into account.
func (acc *Account) Encrypt(key keypair.KeyPair, password string) error {
	pass := derive(password)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	acc.Encrypted = encryptData(pass, iv, key.Private)
	acc.Iv = iv
	acc.Address = key.Address(acc.Network)

	return nil
}

// Decrypt KeyPair from account.
func (acc *Account) Decrypt(password string) (keypair.KeyPair, error) {
	var key keypair.KeyPair
	pass := derive(password)
	data, err := decryptData(pass, acc.Iv, acc.Encrypted)
	if err != nil {
		return key, err
	}

	return keypair.FromSeed(data)
}
