package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/ethereumproject/go-ethereum/crypto/sha3"
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

	account []Account
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

func Serialize(w Wallet) error {
	return nil
}

func Deserialize() (Wallet, error) {
	wlt := Wallet{}

	return wlt, nil
}

//{
//"privateKey": "",
//"name": "mainnet",
//"accounts": {
//"0": {
//"brain": false,
//"algo": "pass:enc",
//"encrypted": "e73e5edaac8393381aa1e5a27b71bbcd5836df93ccd60dc116c8ec0b53f44d0e4bd8472baa227297261f738c6563e43d",
//"iv": "190c85ff1e4a15262ff917b82d5e9d8c",
//"address": "NDLXS2XIAVOPOVHSUZI3N5VU4HJ6ENT24QVIGAPM",
//"label": "Primary",
//"network": 104,
//"child": "613d01ce62e43cc5bea9395e0d97942c45d661e081184245ddebae4e977f336f"
//}
//}
//}
