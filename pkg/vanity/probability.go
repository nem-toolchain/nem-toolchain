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
// with pre-calculated probability `pb` and with specified precision `pr`
func NumberOfAttempts(pb, pr float64) float64 {
	return math.Log2(1.-pr) / math.Log2(1.-pb)
}

// Probability determines a probability to find an address on random basis in one attempt
func Probability(sel Selector) float64 {
	res := float64(0)
	for _, rule := range sel.rules() {
		res += rule.probability()
	}
	if res > 1 {
		res = 1
	}
	return res
}

func (rule searchRule) probability() float64 {
	res := float64(1)
	if rule.exclude != nil && rule.prefix != nil {
		res *= rule.prefix.probability() *
			rule.exclude.probability(uint(len(rule.prefix.prefix)))
	} else if rule.exclude != nil {
		res *= rule.exclude.probability(0)
	} else if rule.prefix != nil {
		res *= rule.prefix.probability()
	}
	return res
}

func (sel excludeSelector) probability(offset uint) float64 {
	if offset >= keypair.AddressLength {
		panic("wrong vanity selector probability offset")
	}
	res := float64(1)
	if offset == 0 {
		offset += 1
		res *= base32FirstPosProbability
	}
	if offset == 1 {
		offset += 1
		res *= 1. - (float64(len(util.IntersectStrings(strings.Split(sel.chars, ""),
			[]string{"A", "B", "C", "D"}))) * base32SecondPosProbability)
	}
	return res *
		math.Pow(1.-float64(len(sel.chars))*base32OtherPosProbability,
			float64(keypair.AddressLength-offset))
}

func (sel prefixSelector) probability() float64 {
	res := float64(1)
	if len(sel.prefix) > 0 {
		res *= base32FirstPosProbability
	}
	if len(sel.prefix) > 1 {
		res *= base32SecondPosProbability
	}
	if len(sel.prefix) > 2 {
		res *= math.Pow(base32OtherPosProbability, float64(len(sel.prefix)-2))
	}
	return res
}
