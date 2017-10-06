// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"math"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

const (
	base32FirstPosProbability  = 1.
	base32SecondPosProbability = 1. / 4
	base32OtherPosProbability  = 1. / 32

	base32OtherCharProbability = 26. * base32OtherPosProbability
	//base32OtherDigitProbability = 6. * base32OtherPosProbability
)

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

// Calculate amount of keypairs to be generated
// to find account with probability `pbty`
func NumberOfKeyPairs(pbty, prec float64) float64 {
	return math.Log(1-prec) / math.Log(1-pbty)
}

func (rule searchRule) probability() float64 {
	res := float64(1)
	if rule.noDigitSelector != nil && rule.prefixSelector != nil {
		res *= rule.noDigitSelector.probability(uint(len(rule.prefixSelector.Prefix)))
	} else if rule.noDigitSelector != nil {
		res *= rule.noDigitSelector.probability(0)
	}
	if rule.prefixSelector != nil {
		res *= rule.prefixSelector.probability()
	}
	return res
}

func (NoDigitSelector) probability(offset uint) float64 {
	// first two chars are always not digits
	if offset < 2 {
		offset = 2
	}
	return math.Pow(base32OtherCharProbability, float64(keypair.AddressLength-offset))
}

func (sel PrefixSelector) probability() float64 {
	switch len(sel.Prefix) {
	case 1:
		return base32FirstPosProbability
	case 2:
		return base32SecondPosProbability
	default:
		return base32SecondPosProbability *
			math.Pow(base32OtherPosProbability, float64(len(sel.Prefix)-2))
	}
}
