// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestNewPrefixSelector(t *testing.T) {
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

func TestNewPrefixSelector_fail(t *testing.T) {
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

func TestFalseSelector_Pass(t *testing.T) {
	assert.False(t, FalseSelector{}.Pass(keypair.Address{}))
}

func TestTrueSelector_Pass(t *testing.T) {
	assert.True(t, TrueSelector{}.Pass(keypair.Address{}))
}

func TestExcludeSelector_Pass_true(t *testing.T) {
	sel := excludeSelector{"BCD234"}
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	assert.True(t, sel.Pass(addr))
}

func TestExcludeSelector_Pass_false(t *testing.T) {
	sel := excludeSelector{"BCD234"}
	for _, s := range []string{
		"TBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA2",
		"TAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAA",
		"TAAAAAAAAAAAAAAAAAAA234AAAAAAAAAAAAAAAAA",
	} {
		t.Run(s, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(s)
			assert.False(t, sel.Pass(addr))
		})
	}
}

func TestPrefixSelector_Pass_true(t *testing.T) {
	sel := prefixSelector{"TABC"}
	addr, _ := keypair.ParseAddress("TABCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	assert.True(t, sel.Pass(addr))
}

func TestPrefixSelector_Pass_false(t *testing.T) {
	sel := prefixSelector{"TABC"}
	for _, s := range []string{
		"TBACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"TACBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"TCABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	} {
		t.Run(s, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(s)
			assert.False(t, sel.Pass(addr))
		})
	}
}
