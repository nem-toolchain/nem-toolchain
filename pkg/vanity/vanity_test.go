// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrefixCorrect(t *testing.T) {
	assert.True(t, isPrefixCorrect("-A"))
	assert.True(t, isPrefixCorrect("-AB"))
	assert.True(t, isPrefixCorrect("_D234"))

	assert.False(t, isPrefixCorrect("_A123"))
	assert.False(t, isPrefixCorrect("_ABC_"))
	assert.False(t, isPrefixCorrect("_XABC"))
}
