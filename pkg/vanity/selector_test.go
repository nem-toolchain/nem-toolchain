// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestNewExcludeSelector(t *testing.T) {
	for _, a := range [][2]string{
		{"", ""},
		{"A", "A"},
		{"ABC", "ABC"},
		{"ABBCCC", "ABC"},
		{"BCACBC", "ABC"},
		{"6D75234653BDACBCCDD", "234567ABCD"},
	} {
		t.Run(a[0], func(t *testing.T) {
			sel, err := NewExcludeSelector(a[0])
			assert.NoError(t, err)
			assert.Equal(t, a[1], sel.(excludeSelector).chars)
		})
	}
}

func TestNewExcludeSelector_fail(t *testing.T) {
	for _, s := range []string{
		"_",
		"123",
	} {
		t.Run(s, func(t *testing.T) {
			_, err := NewExcludeSelector(s)
			assert.Error(t, err)
		})
	}
}

func TestNewPrefixSelector(t *testing.T) {
	for k, v := range map[string]struct {
		ch core.Chain
		pr string
	}{
		"MA":                            {core.Mijin, "MA"},
		"NAB":                           {core.Mainnet, "NAB"},
		"-T-D-2-3-4---":                 {core.Testnet, "TD234"},
		"TA____-BB____-C_CC___-D__D_DD": {core.Testnet, "TA____BB____C_CC___D__D_DD"},
	} {
		t.Run(k, func(t *testing.T) {
			sel, err := NewPrefixSelector(v.ch, k)
			assert.NoError(t, err)
			assert.Equal(t, v.pr, sel.(prefixSelector).prefix)
		})
	}
}

func TestNewPrefixSelector_fail(t *testing.T) {
	for k, v := range map[string]core.Chain{
		"MA123": core.Mijin,
		"TABC":  core.Mainnet,
		"T#__":  core.Testnet,
	} {
		t.Run(k, func(t *testing.T) {
			_, err := NewPrefixSelector(v, k)
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
	for k, v := range map[string]string{
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "BCD234",
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.True(t, excludeSelector{v}.Pass(addr))
		})
	}
}

func TestExcludeSelector_Pass_false(t *testing.T) {
	for k, v := range map[string]string{
		"TBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "BCD234",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA2": "BCD234",
		"TAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAA": "BCD234",
		"TAAAAAAAAAAAAAAAAAAA234AAAAAAAAAAAAAAAAA": "BCD234",
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.False(t, excludeSelector{v}.Pass(addr))
		})
	}
}

func TestPrefixSelector_Pass_true(t *testing.T) {
	for k, v := range map[string]string{
		"TABCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "TABC",

		"TAAAAABBBBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": "TA____BB____C_CC___D__D_DD",
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.True(t, prefixSelector{v}.Pass(addr))
		})
	}
}

func TestPrefixSelector_Pass_false(t *testing.T) {
	for k, v := range map[string]string{
		"TBACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "TABC",
		"TACBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "TABC",
		"TCABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": "TABC",

		"TAAAAA2BBBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": "TA____BB____C_CC___D__D_DD",
		"TAAAAABBBBBBCC33CCDDDDDDEEEEEEFFFFFF2345": "TA____BB____C_CC___D__D_DD",
		"TAAAAABBBBBBCCCCCCDDDDD4EEEEEEFFFFFF2345": "TA____BB____C_CC___D__D_DD",
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.False(t, prefixSelector{v}.Pass(addr))
		})
	}
}
