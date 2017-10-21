// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestNumberOfKeyPairs(t *testing.T) {
	assert.True(t, math.IsInf(NumberOfKeyPairs(0., 0.5), 0))
	assert.True(t, math.IsInf(NumberOfKeyPairs(0., 0.8), 0))
	assert.True(t, math.IsInf(NumberOfKeyPairs(0., 0.99), 0))

	assert.Equal(t, 1., NumberOfKeyPairs(0.5, 0.5))
	assert.InDelta(t, 2.32, NumberOfKeyPairs(0.5, 0.8), 0.01)
	assert.InDelta(t, 6.64, NumberOfKeyPairs(0.5, 0.99), 0.01)

	assert.Equal(t, 0., NumberOfKeyPairs(1., 0.5))
	assert.Equal(t, 0., NumberOfKeyPairs(1., 0.8))
	assert.Equal(t, 0., NumberOfKeyPairs(1., 0.99))
}
