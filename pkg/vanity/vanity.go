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

// StartSearch starts search process for vanity account address that satisfies a given selector
func StartSearch(chain core.Chain, selector Selector, ch chan<- keypair.KeyPair) {
	for {
		pair := keypair.Gen()
		if !selector.Pass(pair.Address(chain)) {
			continue
		}
		ch <- pair
		return
	}
}

// AndMultiSelector combines several selectors into a sequential multi selector chain (`AND` logic)
func AndMultiSelector(selectors ...Selector) Selector {
	return seqMultiSelector{selectors}
}

// OrMultiSelector combines several selectors into a parallel multi selector chain (`OR` logic)
func OrMultiSelector(selectors ...Selector) Selector {
	return parMultiSelector{selectors}
}

// Selector defines generic search strategy
type Selector interface {
	// Pass checks address by a given search pattern
	Pass(addr keypair.Address) bool
	// rules converts current selector into a bundle of normalized rules
	rules() []searchRule
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

// NewPrefixSelector creates a new prefix selector from given string
func NewPrefixSelector(ch core.Chain, prefix string) (PrefixSelector, error) {
	if !isPrefixCorrect(ch, prefix) {
		return PrefixSelector{}, fmt.Errorf("incorrect prefix '%v'", prefix)
	}
	return PrefixSelector{prefix}, nil
}

// Pass returns always false
func (FalseSelector) Pass(keypair.Address) bool {
	return false
}

func (FalseSelector) rules() []searchRule {
	// nothing, just skip it
	return []searchRule{}
}

// Pass returns always true
func (TrueSelector) Pass(keypair.Address) bool {
	return true
}

func (TrueSelector) rules() []searchRule {
	// empty searchRule - always true
	return []searchRule{{}}
}

// Pass returns true only if address doesn't contain any digits
func (NoDigitSelector) Pass(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String()[2:], "234567")
}

func (sel NoDigitSelector) rules() []searchRule {
	return []searchRule{{noDigitSelector: &sel}}
}

// Pass returns true only if address has a given prefix
func (sel PrefixSelector) Pass(addr keypair.Address) bool {
	return strings.HasPrefix(addr.String(), sel.Prefix)
}

func (sel PrefixSelector) rules() []searchRule {
	return []searchRule{{prefixSelector: &sel}}
}

// seqMultiSelector allows nested selectors to be combined into a sequential multi selector chain (`AND` logic)
type seqMultiSelector struct {
	items []Selector
}

// parMultiSelector allows nested selectors to be combined into a parallel multi selector chain (`OR` logic)
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

func (sel seqMultiSelector) rules() []searchRule {
	res := []searchRule{}
	for _, it := range sel.items {
		res = combineSearchRules(res, it.rules())
	}
	return res
}

func (sel parMultiSelector) Pass(addr keypair.Address) bool {
	for _, it := range sel.items {
		if it.Pass(addr) {
			return true
		}
	}
	return false
}

func (sel parMultiSelector) rules() []searchRule {
	res := []searchRule{}
	for _, it := range sel.items {
		res = append(res, it.rules()...)
	}
	return res
}

// combineSearchRules joins two separate bundle of search rules into a new one (AND logic)
func combineSearchRules(rls1 []searchRule, rls2 []searchRule) []searchRule {
	res := make([]searchRule, len(rls1)*len(rls2))
	for i, r1 := range rls1 {
		for j, r2 := range rls1 {
			res[i*j] = r1.merge(r2)
		}
	}
	return res
}

// searchRule is one row declarative search rule used for example to calculate probability
type searchRule struct {
	noDigitSelector *NoDigitSelector
	prefixSelector  *PrefixSelector
}

func (rule searchRule) merge(other searchRule) searchRule {
	res := rule
	if res.noDigitSelector == nil {
		res.noDigitSelector = other.noDigitSelector
	}
	if res.prefixSelector == nil {
		res.prefixSelector = other.prefixSelector
	}
	return res
}

// isPrefixCorrect verify that prefix can be used
func isPrefixCorrect(ch core.Chain, prefix string) bool {
	str := fmt.Sprintf("^%v[A-D][A-Z2-7]*$", ch.ChainPrefix())
	return regexp.MustCompile(str).MatchString(prefix)
}
