// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"golang.org/x/crypto/sha3"
)

// FromRaw tries to create Account from provided data
func FromRaw(raw interface{}) (Account, error) {
	var account Account
	data := raw.(map[string]interface{})
	for k, v := range data {
		switch k {
		case "network":
			val := v.(float64)
			ch, err := core.NewChain(byte(val))
			if err != nil {
				return account, err
			}
			account.Network = ch
		case "label":
			account.Label = v.(string)
		case "encrypted":
			encrypted, err := hex.DecodeString(v.(string))
			if err != nil {
				return account, err
			}
			account.Encrypted = encrypted
		case "iv":
			iv, err := hex.DecodeString(v.(string))
			if err != nil {
				return account, err
			}
			account.Iv = iv
		case "address":
			addr, err := keypair.ParseAddress(v.(string))
			if err != nil {
				return account, err
			}
			account.Address = addr
		case "child":
			child, err := hex.DecodeString(v.(string))
			if err != nil {
				return account, err
			}
			account.Child = child
		case "algo":
			account.Algo = v.(string)
		case "brain":
			account.Brain = v.(bool)
		}
	}

	return account, nil
}

// Account used to encrypt private key
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

// SerializableAccount json-serializable form for Account
type SerializableAccount struct {
	Brain     bool   `json:"brain"`
	Algo      string `json:"algo"`
	Encrypted string `json:"encrypted"`
	Iv        string `json:"iv"`
	Address   string `json:"address"`
	Label     string `json:"label"`
	Network   byte   `json:"network"`
	Child     string `json:"child"`
}

// Serializable convert account into json-serializable form
func (acc Account) Serializable() SerializableAccount {
	var ser SerializableAccount
	ser.Brain = acc.Brain
	ser.Algo = acc.Algo
	ser.Encrypted = hex.EncodeToString(acc.Encrypted)
	ser.Iv = hex.EncodeToString(acc.Iv)
	ser.Address = acc.Address.String()
	ser.Label = acc.Label
	ser.Network = acc.Network.ID
	ser.Child = hex.EncodeToString(acc.Child)

	return ser
}

// Encrypt KeyPair into account
func (acc *Account) Encrypt(key keypair.KeyPair, password string) error {
	pass, err := derive(password)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(pass)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, len(key.Private))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, key.Private)

	acc.Encrypted = ciphertext
	acc.Iv = iv
	acc.Address = key.Address(acc.Network)

	return nil
}

// Decrypt KeyPair from account
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

	privKeyBytes := make([]byte, 32)
	mode := cipher.NewCBCDecrypter(block, acc.Iv)
	mode.CryptBlocks(privKeyBytes, acc.Encrypted)
	key = keypair.FromSeed(privKeyBytes)

	return key, nil
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
