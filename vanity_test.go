package example

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAccount(t *testing.T) {
	var assert = assert.New(t)
	acc := "TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC"
	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")

	account, err := toAccount(pub, TestnetId)
	assert.NoError(err)

	if account != acc {
		t.Fatalf("Mismatched account: expected %v but received %v", acc, account)
	}
}
