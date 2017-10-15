// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"strings"

	"regexp"

	"fmt"

	"sort"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/util"
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

// ExcludeSelector checks an address for absence of given characters
type ExcludeSelector struct {
	// chars determines a unique sorted list of excluded characters
	chars string
}

// PrefixSelector checks an address by given prefix
type PrefixSelector struct {
	// prefix determines a required address prefix to search
	prefix string
}

// NewExcludeSelector creates a new exclude selector from given string
func NewExcludeSelector(chars string) (ExcludeSelector, error) {
	if !regexp.MustCompile(`^[A-Z2-7]*$`).MatchString(chars) {
		return ExcludeSelector{}, fmt.Errorf("incorrect exclude characters '%v'", chars)
	}
	arr := strings.Split(chars, "")
	sort.Strings(arr)
	return ExcludeSelector{strings.Join(arr, "")}, nil
}

// NewPrefixSelector creates a new prefix selector from given string
func NewPrefixSelector(ch core.Chain, prefix string) (PrefixSelector, error) {
	str := fmt.Sprintf(`^%v[A-D][A-Z2-7]*$`, ch.ChainPrefix())
	if !regexp.MustCompile(str).MatchString(prefix) {
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
func (sel ExcludeSelector) Pass(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String()[1:], sel.chars)
}

func (sel ExcludeSelector) merge(other ExcludeSelector) ExcludeSelector {
	arr := append(
		strings.Split(sel.chars, ""),
		strings.Split(other.chars, "")...)
	sort.Strings(arr)
	return ExcludeSelector{strings.Join(util.DistinctStrings(arr), "")}
}

func (sel ExcludeSelector) rules() []searchRule {
	return []searchRule{{excludeSelector: &sel}}
}

// Pass returns true only if address has a given prefix
func (sel PrefixSelector) Pass(addr keypair.Address) bool {
	return strings.HasPrefix(addr.String(), sel.prefix)
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
	res := []searchRule{{}}
	for _, it := range sel.items {
		n := []searchRule{}
		for _, r1 := range res {
			for _, r2 := range it.rules() {
				n = append(n, r1.merge(r2))
			}
		}
		res = n
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

// searchRule is one row declarative search rule used for example to calculate probability
type searchRule struct {
	excludeSelector *ExcludeSelector
	prefixSelector  *PrefixSelector
}

func (rule searchRule) merge(other searchRule) searchRule {
	if rule.excludeSelector == nil {
		rule.excludeSelector = other.excludeSelector
	} else if other.excludeSelector != nil {
		s := rule.excludeSelector.merge(*other.excludeSelector)
		rule.excludeSelector = &s
	}
	if rule.prefixSelector == nil {
		rule.prefixSelector = other.prefixSelector
	}
	return rule
}
