// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/base64"
	"encoding/json"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

type Wallet struct {
	PrivateKey string
	Name       string
	Accounts   map[string]Account
}

func New(chain core.Chain) Wallet {
	wlt := Wallet{}
	wlt.Name = chain.String()
	wlt.PrivateKey = ""
	wlt.Accounts = make(map[string]Account)

	return wlt
}

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

func (wlt *Wallet) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	data := f.(map[string]interface{})
	wlt.Name = data["name"].(string)
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

func (wlt Wallet) MarshalJSON() ([]byte, error) {
	aux_accounts := make(map[string]SerializableAccount)
	for i, acc := range wlt.Accounts {
		aux_accounts[i] = acc.Serializable()
	}

	acc := struct {
		PrivateKey string                         `json:"privateKey"`
		Name       string                         `json:"name"`
		Accounts   map[string]SerializableAccount `json:"accounts"`
	}{
		PrivateKey: wlt.PrivateKey,
		Name:       wlt.Name,
		Accounts:   aux_accounts,
	}

	enc, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func Serialize(w Wallet) (string, error) {
	var encoded string

	ser, err := json.Marshal(w)
	if err != nil {
		return encoded, err
	}
	encoded = base64.StdEncoding.EncodeToString(ser)

	return encoded, nil
}

func Deserialize(data string) (Wallet, error) {
	var wlt Wallet

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return wlt, err
	}

	err = json.Unmarshal(decoded, &wlt)

	return wlt, err
}
