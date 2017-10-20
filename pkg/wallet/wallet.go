package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/ethereumproject/go-ethereum/crypto/sha3"
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

type Wallet struct {
	privateKey string
	name       string
	accounts   map[string]Account
}

type EncryptedKey struct {
	ciphertext []byte
	iv         []byte
}

func encrypt(key []byte, password []byte) (EncryptedKey, error) {
	var encrypted EncryptedKey
	block, err := aes.NewCipher(password)

	if err != nil {
		return encrypted, err
	}

	ciphertext := make([]byte, len(key))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return encrypted, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, key)

	encrypted = EncryptedKey{
		ciphertext: ciphertext,
		iv:         iv,
	}

	return encrypted, nil
}

func decrypt(encrypted EncryptedKey, password []byte) ([]byte, error) {
	pk := make([]byte, 32)
	block, err := aes.NewCipher(password)

	if err != nil {
		return pk, err
	}

	mode := cipher.NewCBCDecrypter(block, encrypted.iv)
	mode.CryptBlocks(pk, encrypted.ciphertext)

	return pk, nil
}

func derive(pass []byte) []byte {
	hash := sha3.NewKeccak256()
	h := pass

	for i := 0; i < 20; i++ {
		hash.Write(h)
		h = hash.Sum(nil)
		hash.Reset()
	}

	return h
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
		var account Account
		data := v.(map[string]interface{})
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
				addr, _ := keypair.FromString(v.(string))
				account.address = addr
			case "child":
				account.child, _ = hex.DecodeString(v.(string))
			case "algo":
				account.algo = v.(string)
			case "brain":
				account.brain = v.(bool)
			}
		}

		wlt.accounts[i] = account
	}

	return nil
}

func Serialize(w Wallet) error {
	return nil
}

func Deserialize(path string) (Wallet, error) {
	var wlt Wallet

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return wlt, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return wlt, err
	}

	err = json.Unmarshal(decoded, &wlt)
	if err != nil {
		return wlt, err
	}

	return wlt, nil
}
