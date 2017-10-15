// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestNewPrefixSelector_normal(t *testing.T) {
	for k, v := range map[core.Chain]string{
		core.Mijin:   "MA",
		core.Mainnet: "NAB",
		core.Testnet: "TD234",
	} {
		t.Run(k.ChainPrefix(), func(t *testing.T) {
			_, err := NewPrefixSelector(k, v)
			assert.NoError(t, err)
		})
	}
}

func TestNewPrefixSelector_error(t *testing.T) {
	for _, p := range []string{
		"MA123",
		"NABC",
	} {
		t.Run(p, func(t *testing.T) {
			_, err := NewPrefixSelector(core.Mijin, p)
			assert.Error(t, err)
		})
	}
}

func TestTrueSelector_pass(t *testing.T) {
	sel := ExcludeSelector{"BCD234"}

	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	assert.True(t, sel.Pass(addr))
}
