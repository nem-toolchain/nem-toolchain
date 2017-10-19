package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/hex"
)

func TestWalletSerialize(t *testing.T) {

}

func TestWalletDeserialize(t *testing.T) {

}

func TestPrivateKey_EncryptDecrypt(t *testing.T) {
	pr, _ := hex.DecodeString("2a91e1d5c110a8d0105aad4683f962c2a56663a3cad46666b16d243174673d90")
	exp, _ := hex.DecodeString("8cd87bc513857a7079d182a6e19b370e907107d97bd3f81a85bcebcc4b5bd3b5")

	pass := derive([]byte("TestTest"))
	assert.Equal(t, pass, exp)

	enc, _ := encrypt(pr, pass)
	dec, _ := decrypt(enc, pass)
	assert.Equal(t, dec, pr)
}
