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
		res *= rule.prefix.probability()
		if len(rule.prefix.prefix) > 1 && (rule.prefix.prefix[1] != PrefixPlaceholder) {
			res *= rule.exclude.probability(2,
				uint(keypair.AddressLength-2-len(strings.Replace(
					rule.prefix.prefix[2:], string(PrefixPlaceholder), "", -1))))
		} else {
			res *= rule.exclude.probability(0,
				uint(keypair.AddressLength-len(strings.Replace(
					rule.prefix.prefix[1:], string(PrefixPlaceholder), "", -1))))
		}
	} else if rule.exclude != nil {
		res *= rule.exclude.probability(0, keypair.AddressLength)
	} else if rule.prefix != nil {
		res *= rule.prefix.probability()
	}
	return res
}

func (sel excludeSelector) probability(offset, length uint) float64 {
	if offset > keypair.AddressLength {
		panic("wrong vanity selector probability offset")
	}
	if offset+length > keypair.AddressLength {
		length = keypair.AddressLength - offset
	}
	res := 1.
	if length == 0 {
		return res
	}
	if offset == 0 {
		offset, length, res = 1, length-1, res*base32FirstPosProbability
	}
	if length == 0 {
		return res
	}
	if offset == 1 {
		_, length, res = 2, length-1, res*(1.-
			(float64(len(util.IntersectStrings([]string{"A", "B", "C", "D"},
				strings.Split(sel.chars, ""))))*base32SecondPosProbability))
	}
	if length == 0 {
		return res
	}
	return res *
		math.Pow(1.-float64(len(sel.chars))*base32OtherPosProbability, float64(length))
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
