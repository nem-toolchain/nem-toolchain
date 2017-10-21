// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
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

func (a Account) Serializable() SerializableAccount {
	var ser SerializableAccount
	ser.Brain = a.brain
	ser.Algo = a.algo
	ser.Encrypted = hex.EncodeToString(a.encrypted)
	ser.Iv = hex.EncodeToString(a.iv)
	ser.Address = a.address.String()
	ser.Label = a.label
	ser.Network = a.network.Id
	ser.Child = hex.EncodeToString(a.child)

	return ser
}
