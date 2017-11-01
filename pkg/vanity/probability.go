// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"math"

	"strings"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/util"
)

const (
	base32FirstPosProbability  = 1.
	base32SecondPosProbability = 1. / 4
	base32OtherPosProbability  = 1. / 32
)

// NumberOfAttempts calculates number of keypairs to be generated to find account
// with pre-calculated probability `pbty` and with specified precision `prec`
func NumberOfAttempts(pbty, prec float64) float64 {
	return math.Log2(1.-prec) / math.Log2(1.-pbty)
}

// Probability determines a probability to find an address on random basis in one attempt
func Probability(sel Selector) float64 {
	res := 0.
	for _, rule := range sel.rules() {
		res += rule.probability()
	}
	if res > 1 {
		res = 1
	}
	return res
}

func (rule searchRule) probability() float64 {
	res := 1.
	if rule.exclude != nil && rule.prefix != nil {
		low, end := uint(0), uint(keypair.AddressLength-len(strings.Replace(
			rule.prefix.prefix[1:], string(PrefixPlaceholder), "", -1)))
		if len(rule.prefix.prefix) > 1 && (rule.prefix.prefix[1] != PrefixPlaceholder) {
			low, end = 2, end+1
		}
		res *= rule.exclude.probability(low, end) * rule.prefix.probability()
	} else if rule.exclude != nil {
		res *= rule.exclude.probability(0, keypair.AddressLength)
	} else if rule.prefix != nil {
		res *= rule.prefix.probability()
	}
	return res
}

func (sel excludeSelector) probability(lo, hi uint) float64 {
	if lo >= hi || hi > keypair.AddressLength {
		panic("vanity selector probability incorrect arguments")
	}
	res := 1.
	if lo == 0 {
		lo, res = 1, res*base32FirstPosProbability
		if hi == 1 {
			return res
		}
	}
	if lo == 1 {
		lo, res = 2, res*(1.-
			(float64(len(util.IntersectStrings([]string{"A", "B", "C", "D"},
				strings.Split(sel.chars, ""))))*base32SecondPosProbability))
	}
	return res *
		math.Pow(1.-float64(len(sel.chars))*base32OtherPosProbability, float64(hi-lo))
}

func (sel prefixSelector) probability() float64 {
	res := 1.
	if len(sel.prefix) > 0 && (sel.prefix[0] != PrefixPlaceholder) {
		res *= base32FirstPosProbability
	}
	if len(sel.prefix) > 1 && (sel.prefix[1] != PrefixPlaceholder) {
		res *= base32SecondPosProbability
	}
	if len(sel.prefix) > 2 {
		res *= math.Pow(base32OtherPosProbability,
			float64(len(strings.Replace(sel.prefix[2:], string(PrefixPlaceholder), "", -1))))
	}
	return res
}
