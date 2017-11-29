// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package account

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"golang.org/x/crypto/sha3"
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
	pass, err := derive(password)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(pass)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, len(key.Private)+aes.BlockSize)
	padding := make([]byte, aes.BlockSize)
	paddedData := append(padding, key.Private...)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	acc.Encrypted = ciphertext
	acc.Iv = iv
	acc.Address = key.Address(acc.Network)

	return nil
}

// Decrypt KeyPair from account.
func (acc *Account) Decrypt(password string) (keypair.KeyPair, error) {
	var key keypair.KeyPair
	pass, err := derive(password)
	if err != nil {
		return key, err
	}

	block, err := aes.NewCipher(pass)
	if err != nil {
		return key, err
	}

	privKeyBytes := make([]byte, 48)
	mode := cipher.NewCBCDecrypter(block, acc.Iv)
	mode.CryptBlocks(privKeyBytes, acc.Encrypted)

	return keypair.FromSeed(privKeyBytes[16:])
}

func derive(password string) ([]byte, error) {
	pass := []byte(password)
	hash := sha3.NewKeccak256()
	h := pass

	for i := 0; i < 20; i++ {
		_, err := hash.Write(h)
		if err != nil {
			return h, err
		}

		h = hash.Sum(nil)
		hash.Reset()
	}

	return h, nil
}
