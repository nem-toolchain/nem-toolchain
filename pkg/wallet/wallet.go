// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Wallet struct {
	privateKey string
	name       string
	accounts   map[string]Account
}

func (wlt *Wallet) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	data := f.(map[string]interface{})
	wlt.name = data["name"].(string)
	wlt.privateKey = data["privateKey"].(string)

	wlt.accounts = make(map[string]Account)
	val := data["accounts"].(map[string]interface{})
	for i, v := range val {
		account, err := FromRaw(v)
		if err != nil {
			return err
		}

		wlt.accounts[i] = account
	}

	return nil
}

func (wlt Wallet) MarshalJSON() ([]byte, error) {
	aux_accounts := make(map[string]SerializableAccount)
	for i, acc := range wlt.accounts {
		aux_accounts[i] = acc.Serializable()
	}

	acc := struct {
		PrivateKey string                         `json:"privateKey"`
		Name       string                         `json:"name"`
		Accounts   map[string]SerializableAccount `json:"accounts"`
	}{
		PrivateKey: wlt.privateKey,
		Name:       wlt.name,
		Accounts:   aux_accounts,
	}

	enc, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func ReadWallet(path string) (Wallet, error) {
	var wlt Wallet

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return wlt, err
	}

	wlt, err = Deserialize(string(data))

	return wlt, err
}

func WriteWallet(path string, wlt Wallet) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatal(e)
		}
	}()

	ser, err := Serialize(wlt)
	if err != nil {
		return err
	}

	_, err = f.WriteString(ser)
	if err != nil {
		return err
	}
	err = f.Sync()

	return err
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
