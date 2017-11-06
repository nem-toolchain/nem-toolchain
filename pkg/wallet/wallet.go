// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/base64"
	"encoding/json"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// New creates new default wallet
func New(chain core.Chain) Wallet {
	wlt := Wallet{}
	wlt.Chain = chain
	wlt.PrivateKey = ""
	wlt.Accounts = make(map[string]Account)

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
	Chain      core.Chain
	Accounts   map[string]Account
}

// AddAccount adds account into wallet
func (wlt *Wallet) AddAccount(pair keypair.KeyPair, password string) error {
	acc := Account{}
	err := acc.Encrypt(pair, password)
	if err != nil {
		return err
	}

	i := string(len(wlt.Accounts))
	wlt.Accounts[i] = acc

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
	wlt.Chain, err = core.FromString(data["name"].(string))
	if err != nil {
		return err
	}
	wlt.PrivateKey = data["privateKey"].(string)

	wlt.Accounts = make(map[string]Account)
	val := data["accounts"].(map[string]interface{})
	for i, v := range val {
		account, err := FromRaw(v)
		if err != nil {
			return err
		}

		wlt.Accounts[i] = account
	}

	return nil
}

// MarshalJSON serialize Wallet into JSON
func (wlt Wallet) MarshalJSON() ([]byte, error) {
	auxAccounts := make(map[string]SerializableAccount)
	for i, acc := range wlt.Accounts {
		auxAccounts[i] = acc.Serializable()
	}

	acc := struct {
		PrivateKey string                         `json:"privateKey"`
		Name       string                         `json:"name"`
		Accounts   map[string]SerializableAccount `json:"accounts"`
	}{
		PrivateKey: wlt.PrivateKey,
		Name:       wlt.Chain.String(),
		Accounts:   auxAccounts,
	}

	enc, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	return enc, nil
}
