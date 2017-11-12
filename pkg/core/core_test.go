// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChain(t *testing.T) {
	for k, v := range map[byte]Chain{
		0x60: Mijin,
		0x68: Mainnet,
		0x98: Testnet,
	} {
		t.Run(string(k), func(t *testing.T) {
			ch, err := NewChain(k)
			assert.NoError(t, err)
			assert.Equal(t, v, ch)
		})
	}
}

func TestNewChain_fail(t *testing.T) {
	_, err := NewChain(0x12)
	assert.Error(t, err)
}

func TestFromString(t *testing.T) {
	for k, v := range map[string]Chain{
		"MIJIN":   Mijin,
		"MAINnet": Mainnet,
		"testNet": Testnet,
	} {
		t.Run(k, func(t *testing.T) {
			ch, err := FromString(k)
			assert.NoError(t, err)
			assert.Equal(t, v, ch)
		})
	}
}

func TestFromString_fail(t *testing.T) {
	for _, s := range []string{
		"",
		"123",
		"_AINnet",
	} {
		t.Run(s, func(t *testing.T) {
			_, err := FromString(s)
			assert.Error(t, err)
		})
	}
}

func TestChain_String(t *testing.T) {
	assert.Equal(t, Mijin.String(), "mijin")
	assert.Equal(t, Mainnet.String(), "mainnet")
	assert.Equal(t, Testnet.String(), "testnet")
}

func TestChain_String_panic(t *testing.T) {
	assert.Panics(t, func() { _ = Chain{0x12}.String() })
}

func TestChain_ChainPrefix(t *testing.T) {
	assert.Equal(t, "M", Mijin.Prefix())
	assert.Equal(t, "N", Mainnet.Prefix())
	assert.Equal(t, "T", Testnet.Prefix())
}

func TestIsChainPrefix_true(t *testing.T) {
	assert.True(t, IsChainPrefix("M"))
	assert.True(t, IsChainPrefix("N123"))
	assert.True(t, IsChainPrefix("TABC"))
}

func TestIsChainPrefix_false(t *testing.T) {
	assert.False(t, IsChainPrefix(""))
	assert.False(t, IsChainPrefix("123"))
	assert.False(t, IsChainPrefix("ABC"))
}

func TestChain_ChainPrefix_panic(t *testing.T) {
	assert.Panics(t, func() { Chain{byte(0x00)}.Prefix() })
}
