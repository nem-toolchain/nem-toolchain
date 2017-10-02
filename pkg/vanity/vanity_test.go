// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestIsPrefixCorrect(t *testing.T) {
	assert.True(t, isPrefixCorrect(core.Mijin, "MA"))
	assert.True(t, isPrefixCorrect(core.Mainnet, "NAB"))
	assert.True(t, isPrefixCorrect(core.Testnet, "TD234"))

	assert.False(t, isPrefixCorrect(core.Mijin, "MA123"))
	assert.False(t, isPrefixCorrect(core.Mijin, "NABC"))
}
