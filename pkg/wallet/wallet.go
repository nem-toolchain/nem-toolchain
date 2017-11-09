// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// NewWallet creates new default wallet
func NewWallet() Wallet {
	wlt := Wallet{}
	wlt.Name = ""
	wlt.PrivateKey = ""
	wlt.Accounts = make(map[uint]Account)

	return wlt
}

// Encode encodes wallet into base64 string
func Encode(w Wallet) (string, error) {
	var encoded string

	ser, err := json.Marshal(w)
	if err != nil {
		return encoded, err
	}
	encoded = base64.StdEncoding.EncodeToString(ser)

	return encoded, nil
}

// Decode decodes wallet form a base64 string
func Decode(data string) (Wallet, error) {
	var wlt Wallet

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return wlt, err
	}
	err = json.Unmarshal(decoded, &wlt)

	return wlt, err
}

// Wallet represents wallet file
type Wallet struct {
	PrivateKey string
	Name       string
	Accounts   map[uint]Account
}

// AddAccount adds account into wallet
func (wlt *Wallet) AddAccount(ch core.Chain, pair keypair.KeyPair, password string) error {
	acc := NewAccount(ch)
	err := acc.Encrypt(pair, password)
	if err != nil {
		return err
	}

	i := len(wlt.Accounts)
	wlt.Accounts[uint(i)] = acc

	return nil
}

// UnmarshalJSON deserialize json into Wallet
func (wlt *Wallet) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	data := f.(map[string]interface{})
	wlt.Name = data["name"].(string)
	wlt.PrivateKey = data["privateKey"].(string)

	wlt.Accounts = make(map[uint]Account)
	accountsData := data["accounts"].(map[string]interface{})
	for i, values := range accountsData {
		account, err := FromRaw(values)
		if err != nil {
			return err
		}

		index, err := strconv.ParseUint(i, 10, 64)
		if err != nil {
			return err
		}

		wlt.Accounts[uint(index)] = account
	}

	return nil
}

// MarshalJSON serialize Wallet into JSON
func (wlt Wallet) MarshalJSON() ([]byte, error) {
	auxAccounts := make(map[uint]SerializableAccount)
	for i, acc := range wlt.Accounts {
		auxAccounts[i] = acc.Serializable()
	}

	acc := struct {
		PrivateKey string                       `json:"privateKey"`
		Name       string                       `json:"name"`
		Accounts   map[uint]SerializableAccount `json:"accounts"`
	}{
		PrivateKey: wlt.PrivateKey,
		Name:       wlt.Name,
		Accounts:   auxAccounts,
	}

	enc, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	return enc, nil
}
