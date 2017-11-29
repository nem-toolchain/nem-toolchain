package account

import (
	"encoding/hex"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// FromRaw tries to create Account from provided data.
func FromRaw(raw interface{}) (*Account, error) {
	var account = new(Account)
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

// SerializableAccount json-serializable form for Account.
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

// Serializable convert account into json-serializable form.
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
