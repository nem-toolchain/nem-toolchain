package example

import (
	"testing"
	"fmt"
	"encoding/hex"
)

func TestAccountGenerator(t *testing.T) {
	acc := "TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC"

	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")
	account := toAccount(pub, TestnetId)
	fmt.Println(account)
	if account != acc {
		t.Fatalf("Mismatched account: expected %v but received %v", acc, account)
	}
}
