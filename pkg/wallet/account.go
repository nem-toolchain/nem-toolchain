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

type Account struct {
	brain     bool
	algo      string
	encrypted []byte
	iv        []byte
	address   keypair.Address
	label     string
	network   core.Chain
	child     []byte
}

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

func FromRaw(raw interface{}) (Account, error) {
	var account Account
	data := raw.(map[string]interface{})
	for k, v := range data {
		switch k {
		case "network":
			val := v.(float64)
			account.network = core.Chain{byte(val)}
		case "label":
			account.label = v.(string)
		case "encrypted":
			account.encrypted, _ = hex.DecodeString(v.(string))
		case "iv":
			account.iv, _ = hex.DecodeString(v.(string))
		case "address":
			addr, _ := keypair.ParseAddress(v.(string))
			account.address = addr
		case "child":
			account.child, _ = hex.DecodeString(v.(string))
		case "algo":
			account.algo = v.(string)
		case "brain":
			account.brain = v.(bool)
		}
	}

	return account, nil
}

func (acc Account) Serializable() SerializableAccount {
	var ser SerializableAccount
	ser.Brain = acc.brain
	ser.Algo = acc.algo
	ser.Encrypted = hex.EncodeToString(acc.encrypted)
	ser.Iv = hex.EncodeToString(acc.iv)
	ser.Address = acc.address.String()
	ser.Label = acc.label
	ser.Network = acc.network.Id
	ser.Child = hex.EncodeToString(acc.child)

	return ser
}

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

	acc.encrypted = ciphertext
	acc.iv = iv
	acc.address = key.Address(acc.network)

	return nil
}

func (acc *Account) Decrypt(password string) ([]byte, error) {
	pk := make([]byte, 32)

	pass, err := derive(password)
	if err != nil {
		return pk, err
	}

	block, err := aes.NewCipher(pass)
	if err != nil {
		return pk, err
	}

	mode := cipher.NewCBCDecrypter(block, acc.iv)
	mode.CryptBlocks(pk, acc.encrypted)

	return pk, nil
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
