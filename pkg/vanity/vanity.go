// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"strings"

	"regexp"

	"fmt"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
)

// Search for vanity account address that satisfies a given selector
func Search(chain core.Chain, selector Selector, ch chan<- keypair.KeyPair) {
	for {
		pair := keypair.Gen()
		if !selector.Pass(pair.Address(chain)) {
			continue
		}
		ch <- pair
		return
	}
}

// Selector defines generic search strategy
type Selector interface {
	// Pass checks address by a given search pattern
	Pass(addr keypair.Address) bool
}

// FalseSelector rejects all addresses, can be used as a default placeholder
type FalseSelector struct{}

// TrueSelector accepts all addresses, can be used as a default placeholder
type TrueSelector struct{}

// NoDigitSelector checks an address for absence of digits
type NoDigitSelector struct{}

// PrefixSelector checks an address by given prefix
type PrefixSelector struct {
	// Prefix determines a required address prefix to search
	Prefix string
}

// PrefixSelectorFrom creates a new prefix selector from given string
func PrefixSelectorFrom(ch core.Chain, prefix string) (PrefixSelector, error) {
	if !isPrefixCorrect(ch, prefix) {
		return PrefixSelector{}, fmt.Errorf("incorrect prefix '%v'", prefix)
	}
	return PrefixSelector{prefix}, nil
}

// Pass always returns false
func (FalseSelector) Pass(addr keypair.Address) bool {
	return false
}

// Pass always returns true
func (TrueSelector) Pass(addr keypair.Address) bool {
	return true
}

// Pass returns true only if address doesn't contain any digits
func (NoDigitSelector) Pass(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String(), "234567")
}

// Pass returns true only if address has a given prefix
func (pr PrefixSelector) Pass(addr keypair.Address) bool {
	return strings.HasPrefix(addr.String(), pr.Prefix)
}

// AndMultiSelector combines several selectors into a sequential multi selector chain (`AND` logic)
func AndMultiSelector(selectors ...Selector) Selector {
	return seqMultiSelector{selectors}
}

// OrMultiSelector combines several selectors into a parallel multi selector chain (`OR` logic)
func OrMultiSelector(selectors ...Selector) Selector {
	return parMultiSelector{selectors}
}

// seqMultiSelector allows multiple selectors to be combined into a sequential multi selector chain (`AND` logic)
type seqMultiSelector struct {
	items []Selector
}

// parMultiSelector allows multiple selectors to be combined into a parallel multi selector chain (`OR` logic)
type parMultiSelector struct {
	items []Selector
}

func (sel seqMultiSelector) Pass(addr keypair.Address) bool {
	for _, it := range sel.items {
		if !it.Pass(addr) {
			return false
		}
	}
	return true
}

func (sel parMultiSelector) Pass(addr keypair.Address) bool {
	for _, it := range sel.items {
		if it.Pass(addr) {
			return true
		}
	}
	return false
}

// IsPrefixCorrect verify that prefix can be used
func isPrefixCorrect(ch core.Chain, prefix string) bool {
	str := fmt.Sprintf("^%v[A-D][A-Z2-7]*$", ch.ChainPrefix())
	return regexp.MustCompile(str).MatchString(prefix)
}
