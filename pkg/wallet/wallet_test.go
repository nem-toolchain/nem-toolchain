package wallet

import (
	"testing"

	"encoding/hex"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestWalletSerialize(t *testing.T) {

}

func TestWalletDeserialize(t *testing.T) {
	wlt, _ := Deserialize("./mainnet.wlt")

	assert.Equal(t, wlt.name, "mainnet")
	assert.Equal(t, wlt.privateKey, "")

	account := wlt.accounts["0"]
	exp_address, _ := keypair.FromString("NDLXS2XIAVOPOVHSUZI3N5VU4HJ6ENT24QVIGAPM")
	assert.Equal(t, account.address, exp_address)
	assert.Equal(t, account.label, "Primary")

	exp_child, _ := hex.DecodeString("613d01ce62e43cc5bea9395e0d97942c45d661e081184245ddebae4e977f336f")
	assert.Equal(t, account.child, exp_child)

	exp_encrypted, _ := hex.DecodeString("e73e5edaac8393381aa1e5a27b71bbcd5836df93ccd60dc116c8ec0b53f44d0e4bd8472baa227297261f738c6563e43d")
	assert.Equal(t, account.encrypted, exp_encrypted)

	exp_iv, _ := hex.DecodeString("190c85ff1e4a15262ff917b82d5e9d8c")
	assert.Equal(t, account.iv, exp_iv)
	assert.Equal(t, account.network, core.Mainnet)
	assert.Equal(t, account.algo, "pass:enc")
	assert.Equal(t, account.brain, false)
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
