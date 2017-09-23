// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package keypair

import (
	"encoding/hex"
	"testing"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestKeyPair_Address(t *testing.T) {
	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")
	assert.Equal(t,
		"TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC", KeyPair{Public: pub}.Address(core.Testnet))
}
