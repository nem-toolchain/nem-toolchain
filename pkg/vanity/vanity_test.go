package vanity

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenAddress(t *testing.T) {
	_, err := GenAddress(MijinId)
	assert.NoError(t, err)

	_, err = GenAddress(TestnetId)
	assert.NoError(t, err)

	_, err = GenAddress(MainnetId)
	assert.NoError(t, err)

	_, err = GenAddress(0x00)
	assert.Error(t, err)

	_, err = GenAddress(0xFF)
	assert.Error(t, err)
}

func TestToAddress(t *testing.T) {
	acc := "TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC"
	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")

	account, err := ToAddress(pub, TestnetId)
	assert.NoError(t, err)

	if account != acc {
		t.Fatalf("Mismatched account: expected %v but received %v", acc, account)
	}
}

func TestIsValidChainId(t *testing.T) {
	assert.True(t, IsValidChainId(TestnetId))
	assert.True(t, IsValidChainId(MainnetId))
	assert.True(t, IsValidChainId(MijinId))

	assert.False(t, IsValidChainId(0x10))
	assert.False(t, IsValidChainId(0x00))
}
