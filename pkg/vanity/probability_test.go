// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestNumberOfAttempts(t *testing.T) {
	assert.True(t, math.IsInf(NumberOfAttempts(0., 0.5), 0))
	assert.True(t, math.IsInf(NumberOfAttempts(0., 0.8), 0))
	assert.True(t, math.IsInf(NumberOfAttempts(0., 0.99), 0))

	assert.Equal(t, 0., NumberOfAttempts(1., 0.5))
	assert.Equal(t, 0., NumberOfAttempts(1., 0.8))
	assert.Equal(t, 0., NumberOfAttempts(1., 0.99))

	assert.Equal(t, 1., NumberOfAttempts(0.5, 0.5))
	assert.InDelta(t, 2, NumberOfAttempts(0.5, 0.8), 1)
	assert.InDelta(t, 7, NumberOfAttempts(0.5, 0.99), 1)

	assert.InDelta(t, 7., NumberOfAttempts(0.1, 0.5), 1)
	assert.InDelta(t, 15, NumberOfAttempts(0.1, 0.8), 1)
	assert.InDelta(t, 44, NumberOfAttempts(0.1, 0.99), 1)

	assert.InDelta(t, 69314, NumberOfAttempts(1e-5, 0.5), 1)
	assert.InDelta(t, 160943, NumberOfAttempts(1e-5, 0.8), 1)
	assert.InDelta(t, 460515, NumberOfAttempts(1e-5, 0.99), 1)
}

func TestProbability(t *testing.T) {
	assert.Equal(t, 0., Probability(FalseSelector{}))
	assert.Equal(t, 1., Probability(TrueSelector{}))
	assert.InDelta(t, 0.00593, Probability(excludeSelector{"ABC"}), 1e-5)
	assert.InDelta(t, 0.02374, Probability(excludeSelector{"246"}), 1e-5)
	assert.InDelta(t, 0.00024, Probability(prefixSelector{"TABC"}), 1e-5)
	assert.InDelta(t, 0.00469,
		Probability(AndSelector(
			OrSelector(excludeSelector{"2"}, excludeSelector{"4"}, excludeSelector{"6"}),
			excludeSelector{"BCD"}, prefixSelector{"TA"})), 1e-5)
}

func TestSearchRule_probability(t *testing.T) {
	assert.Equal(t, 1., searchRule{}.probability())
	assert.InDelta(t, 0.00593, searchRule{exclude: &excludeSelector{"ABC"}}.probability(), 1e-5)
	assert.InDelta(t, 0.00024, searchRule{prefix: &prefixSelector{"TABC"}}.probability(), 1e-5)
	assert.InDelta(t, 0.0000071,
		searchRule{exclude: &excludeSelector{"ABC"}, prefix: &prefixSelector{"TABC"}}.probability(),
		1e-7)
}

func TestExcludePrefix_probability(t *testing.T) {
	assert.Equal(t, 1., excludeSelector{}.probability(0))
	assert.Equal(t, 1., excludeSelector{}.probability(1))
	assert.Equal(t, 1., excludeSelector{}.probability(2))
	assert.Equal(t, 1., excludeSelector{}.probability(39))

	assert.InDelta(t, 0.22444, excludeSelector{"A"}.probability(0), 1e-5)
	assert.InDelta(t, 0.22444, excludeSelector{"A"}.probability(1), 1e-5)
	assert.InDelta(t, 0.29926, excludeSelector{"A"}.probability(2), 1e-5)
	assert.Equal(t, 0.96875, excludeSelector{"A"}.probability(39))

	assert.InDelta(t, 0.00593, excludeSelector{"ABC"}.probability(0), 1e-5)
	assert.InDelta(t, 0.00593, excludeSelector{"ABC"}.probability(1), 1e-5)
	assert.InDelta(t, 0.02374, excludeSelector{"ABC"}.probability(2), 1e-5)
	assert.Equal(t, 0.90625, excludeSelector{"ABC"}.probability(39))

	assert.InDelta(t, 0.02374, excludeSelector{"246"}.probability(0), 1e-5)
	assert.InDelta(t, 0.02374, excludeSelector{"246"}.probability(1), 1e-5)
	assert.InDelta(t, 0.02374, excludeSelector{"246"}.probability(2), 1e-5)
	assert.InDelta(t, 0.90625, excludeSelector{"246"}.probability(39), 1e-5)
}

func TestExcludePrefix_probability_panic(t *testing.T) {
	assert.Panics(t, func() { excludeSelector{}.probability(40) })
	assert.Panics(t, func() { excludeSelector{}.probability(1000) })
}

func TestPrefixPrefix_probability(t *testing.T) {
	assert.Equal(t, 1., prefixSelector{}.probability())
	assert.Equal(t, 1., prefixSelector{"T"}.probability())
	assert.Equal(t, .25, prefixSelector{"TA"}.probability())
	assert.InDelta(t, 0.00781, prefixSelector{"TAB"}.probability(), 1e-5)
	assert.InDelta(t, 0.00024, prefixSelector{"TABC"}.probability(), 1e-5)
}
