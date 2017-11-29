package wallet

import (
	"encoding/json"
	"strconv"

	"github.com/nem-toolchain/nem-toolchain/pkg/wallet/account"
)

// UnmarshalJSON deserialize json into Wallet.
func (wlt *Wallet) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	data := f.(map[string]interface{})
	wlt.Name = data["name"].(string)
	wlt.PrivateKey = data["privateKey"].(string)

	wlt.Accounts = make(map[uint]*account.Account)
	accountsData := data["accounts"].(map[string]interface{})
	for i, values := range accountsData {
		acc, err := account.FromRaw(values)
		if err != nil {
			return err
		}

		index, err := strconv.ParseUint(i, 10, 64)
		if err != nil {
			return err
		}

		wlt.Accounts[uint(index)] = acc
	}

	return nil
}

// MarshalJSON serialize Wallet into JSON.
func (wlt Wallet) MarshalJSON() ([]byte, error) {
	auxAccounts := make(map[uint]account.SerializableAccount)
	for i, acc := range wlt.Accounts {
		auxAccounts[i] = acc.Serializable()
	}

	acc := struct {
		PrivateKey string                               `json:"privateKey"`
		Name       string                               `json:"name"`
		Accounts   map[uint]account.SerializableAccount `json:"accounts"`
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
