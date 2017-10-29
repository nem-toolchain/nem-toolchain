// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"regexp"

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
		re *regexp.Regexp
	}{
		"":  {core.Mijin, "", regexp.MustCompile("^\\w*")},
		"M": {core.Mijin, "M", regexp.MustCompile("^M\\w*")},
		"_": {core.Mijin, "_", regexp.MustCompile("^\\w\\w*")},

		"NA":   {core.Mainnet, "NA", regexp.MustCompile("^NA\\w*")},
		"NABC": {core.Mainnet, "NABC", regexp.MustCompile("^NABC\\w*")},
		"_ABC": {core.Mainnet, "_ABC", regexp.MustCompile("^\\wABC\\w*")},

		"-T-D-2-3-4---": {
			core.Testnet,
			"TD234",
			regexp.MustCompile("^TD234\\w*"),
		},
		"TA____-BB____-C_CC___-D__D_D": {
			core.Testnet,
			"TA____BB____C_CC___D__D_D",
			regexp.MustCompile("^TA\\w\\w\\w\\wBB\\w\\w\\w\\wC\\wCC\\w\\w\\wD\\w\\wD\\wD\\w*"),
		},
	} {
		t.Run(k, func(t *testing.T) {
			sel, err := NewPrefixSelector(v.ch, k)
			assert.NoError(t, err)
			assert.Equal(t, v.pr, sel.(prefixSelector).prefix)
			assert.Equal(t, v.re, sel.(prefixSelector).re)
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
	for k, v := range map[string](*regexp.Regexp){
		"TABCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": regexp.MustCompile("^TABC\\w*"),
		"TAAAAABBBBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": regexp.MustCompile("^TA\\wAA\\wBB\\w*"),
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.True(t, prefixSelector{re: v}.Pass(addr))
		})
	}
}

func TestPrefixSelector_Pass_false(t *testing.T) {
	for k, v := range map[string](*regexp.Regexp){
		"TBACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": regexp.MustCompile("^TABC\\w*"),
		"TACBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": regexp.MustCompile("^TABC\\w*"),
		"TCABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": regexp.MustCompile("^TABC\\w*"),

		"T2AAAABBBBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": regexp.MustCompile("^TA\\wAA\\wBB\\w*"),
		"TAAA3ABBBBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": regexp.MustCompile("^TA\\wAA\\wBB\\w*"),
		"TAAAAAB4BBBBCCCCCCDDDDDDEEEEEEFFFFFF2345": regexp.MustCompile("^TA\\wAA\\wBB\\w*"),
	} {
		t.Run(k, func(t *testing.T) {
			addr, _ := keypair.ParseAddress(k)
			assert.False(t, prefixSelector{re: v}.Pass(addr))
		})
	}
}
