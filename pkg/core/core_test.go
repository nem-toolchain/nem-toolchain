// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsChainPrefix(t *testing.T) {
	assert.True(t, IsChainPrefix("M"))
	assert.True(t, IsChainPrefix("N123"))
	assert.True(t, IsChainPrefix("TABC"))

	assert.False(t, IsChainPrefix(""))
	assert.False(t, IsChainPrefix("123"))
	assert.False(t, IsChainPrefix("ABC"))
}

func TestChain_ChainPrefix(t *testing.T) {
	assert.Equal(t, "M", Mijin.ChainPrefix())
	assert.Equal(t, "N", Mainnet.ChainPrefix())
	assert.Equal(t, "T", Testnet.ChainPrefix())
}

func TestChain_ChainPrefix_panic(t *testing.T) {
	assert.Panics(t, func() { Chain{byte(0x00)}.ChainPrefix() })
}
