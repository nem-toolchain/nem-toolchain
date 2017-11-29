// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"math"
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
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
	assert.InDelta(t, 0.00024, Probability(prefixSelector{prefix: "TABC"}), 1e-5)

	assert.Equal(t, 1.,
		Probability(OrSelector(excludeSelector{""}, prefixSelector{prefix: ""})))

	assert.InDelta(t, 0.00469,
		Probability(AndSelector(
			OrSelector(excludeSelector{"2"}, excludeSelector{"4"}, excludeSelector{"6"}),
			excludeSelector{"BCD"}, prefixSelector{prefix: "TA"})), 1e-5)
}

func TestSearchRule_probability(t *testing.T) {
	assert.Equal(t, 1., searchRule{}.probability())

	assert.InDelta(t, 0.00593,
		searchRule{exclude: &excludeSelector{"ABC"}}.probability(), 1e-5)
	assert.InDelta(t, 0.00024,
		searchRule{prefix: &prefixSelector{prefix: "TA_B__C___"}}.probability(), 1e-5)
	assert.InDelta(t, 0.000031,
		searchRule{prefix: &prefixSelector{prefix: "T_A__B___C"}}.probability(), 1e-6)
	assert.InDelta(t, 0.0000071, searchRule{
		exclude: &excludeSelector{"ABC"},
		prefix:  &prefixSelector{prefix: "TA_B__C___"}}.probability(), 1e-7)
	assert.InDelta(t, 0.00000024, searchRule{
		exclude: &excludeSelector{"ABC"},
		prefix:  &prefixSelector{prefix: "T_A__B___C"}}.probability(), 1e-8)
}

func TestExcludeSelector_probability(t *testing.T) {
	assert.Equal(t, 1., excludeSelector{}.probability(0, keypair.AddressLength))
	assert.Equal(t, 1., excludeSelector{}.probability(1, keypair.AddressLength))
	assert.Equal(t, 1., excludeSelector{}.probability(2, keypair.AddressLength))
	assert.Equal(t, 1., excludeSelector{}.probability(39, keypair.AddressLength))

	assert.Equal(t, 1., excludeSelector{"246A"}.probability(0, 1))
	assert.Equal(t, 0.75, excludeSelector{"246A"}.probability(0, 2))
	assert.Equal(t, 0.75, excludeSelector{"246A"}.probability(1, 2))
	assert.Equal(t, 0.65625, excludeSelector{"246A"}.probability(1, 3))
	assert.Equal(t, 0.875, excludeSelector{"246A"}.probability(2, 3))
	assert.Equal(t, 0.875, excludeSelector{"246A"}.probability(39, keypair.AddressLength))

	assert.InDelta(t, 0.22444,
		excludeSelector{"A"}.probability(0, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.22444,
		excludeSelector{"A"}.probability(1, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.29926,
		excludeSelector{"A"}.probability(2, keypair.AddressLength), 1e-5)
	assert.Equal(t, 0.96875,
		excludeSelector{"A"}.probability(39, keypair.AddressLength))

	assert.InDelta(t, 0.00593,
		excludeSelector{"ABC"}.probability(0, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.00593,
		excludeSelector{"ABC"}.probability(1, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.02374,
		excludeSelector{"ABC"}.probability(2, keypair.AddressLength), 1e-5)
	assert.Equal(t, 0.90625,
		excludeSelector{"ABC"}.probability(39, keypair.AddressLength))

	assert.InDelta(t, 0.02374,
		excludeSelector{"246"}.probability(0, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.02374,
		excludeSelector{"246"}.probability(1, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.02374,
		excludeSelector{"246"}.probability(2, keypair.AddressLength), 1e-5)
	assert.InDelta(t, 0.90625,
		excludeSelector{"246"}.probability(39, keypair.AddressLength), 1e-5)
}

func TestExcludeSelector_probability_panic(t *testing.T) {
	assert.Panics(t, func() { excludeSelector{}.probability(0, keypair.AddressLength+1) })
	assert.Panics(t, func() { excludeSelector{}.probability(keypair.AddressLength, 0) })

	assert.Panics(t, func() { excludeSelector{}.probability(1000, 123) })
}

func TestPrefixSelector_probability(t *testing.T) {
	assert.Equal(t, 1., prefixSelector{}.probability())
	assert.Equal(t, 1., prefixSelector{prefix: "T"}.probability())
	assert.Equal(t, .25, prefixSelector{prefix: "TA"}.probability())

	assert.InDelta(t, 0.00781, prefixSelector{prefix: "TAB"}.probability(), 1e-5)
	assert.InDelta(t, 0.00024, prefixSelector{prefix: "TABC"}.probability(), 1e-5)

	assert.InDelta(t, 0.00098, prefixSelector{prefix: "___B__C___"}.probability(), 1e-5)
	assert.InDelta(t, 0.00024, prefixSelector{prefix: "_A_B__C___"}.probability(), 1e-5)
	assert.InDelta(t, 0.00024, prefixSelector{prefix: "TA_B__C___"}.probability(), 1e-5)
	assert.InDelta(t, 0.000031, prefixSelector{prefix: "T_A__B___C"}.probability(), 1e-6)
}
