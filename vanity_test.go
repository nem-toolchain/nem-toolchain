package vanity

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAccount(t *testing.T) {
	var assert = assert.New(t)
	acc := "TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC"
	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")

	account, err := ToAccount(pub, TestnetId)
	assert.NoError(err)

	if account != acc {
		t.Fatalf("Mismatched account: expected %v but received %v", acc, account)
	}
}

func TestToChainId(t *testing.T) {
	var assert = assert.New(t)

	_, err := ToChainId("mainnet")
	assert.NoError(err)
	_, err = ToChainId("main")
	assert.NoError(err)
	_, err = ToChainId("0x68")
	assert.NoError(err)
	_, err = ToChainId("68")
	assert.NoError(err)
	_, err = ToChainId("testnet")
	assert.NoError(err)
	_, err = ToChainId("test")
	assert.NoError(err)
	_, err = ToChainId("0x98")
	assert.NoError(err)
	_, err = ToChainId("98")
	assert.NoError(err)
	_, err = ToChainId("mijin")
	assert.NoError(err)
	_, err = ToChainId("0x60")
	assert.NoError(err)
	_, err = ToChainId("60")
	assert.NoError(err)

	_, err = ToChainId("mainNet")
	assert.Error(err)
	_, err = ToChainId("Main")
	assert.Error(err)
	_, err = ToChainId("TE")
	assert.Error(err)
	_, err = ToChainId("testN")
	assert.Error(err)
	_, err = ToChainId("0x100")
	assert.Error(err)
	_, err = ToChainId("0x")
	assert.Error(err)
	_, err = ToChainId("")
	assert.Error(err)
	_, err = ToChainId("miiijni")
	assert.Error(err)
	_, err = ToChainId("101010")
	assert.Error(err)
}

func TestGenerateAccount(t *testing.T) {
	var assert = assert.New(t)
	_, err := GenerateAccount(TestnetId)
	assert.NoError(err)
	//strings.HasPrefix(acc, "TC")

	_, err = GenerateAccount(MainnetId)
	assert.NoError(err)
	//strings.HasPrefix(acc, "NA")

	_, err = GenerateAccount(MijinId)
	assert.NoError(err)
	//strings.HasPrefix(acc, "TA")

	_, err = GenerateAccount(0x00)
	assert.Error(err)
	_, err = GenerateAccount(0xFF)
	assert.Error(err)
}

func TestIsValidChainId(t *testing.T) {
	var assert = assert.New(t)
	assert.True(IsValidChainId(TestnetId))
	assert.True(IsValidChainId(MainnetId))
	assert.True(IsValidChainId(MijinId))
	assert.False(IsValidChainId(0x10))
	assert.False(IsValidChainId(0x00))
}
