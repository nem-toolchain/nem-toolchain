// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"strings"

	"fmt"
	"regexp"
	"sort"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/util"
)

// NewExcludeSelector creates a new exclude selector from given string
func NewExcludeSelector(chars string) (Selector, error) {
	if !regexp.MustCompile(`^[A-Z2-7]*$`).MatchString(chars) {
		return excludeSelector{}, fmt.Errorf("incorrect exclude characters '%v'", chars)
	}
	arr := strings.Split(chars, "")
	sort.Strings(arr)
	return excludeSelector{strings.Join(util.DistinctStrings(arr), "")}, nil
}

// NewPrefixSelector creates a new prefix selector from given string
func NewPrefixSelector(ch core.Chain, prefix string) (Selector, error) {
	str := fmt.Sprintf(`^%v[A-D][A-Z2-7]*$`, ch.ChainPrefix())
	if !regexp.MustCompile(str).MatchString(prefix) {
		return prefixSelector{}, fmt.Errorf("incorrect prefix '%v'", prefix)
	}
	return prefixSelector{prefix}, nil
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

// exclude checks an address for absence of given characters
type excludeSelector struct {
	// chars determines a unique sorted list of excluded characters
	chars string
}

// prefix checks an address by given prefix
type prefixSelector struct {
	// prefix determines a required address prefix to search
	prefix string
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
func (sel excludeSelector) Pass(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String()[1:], sel.chars)
}

func (sel excludeSelector) rules() []searchRule {
	return []searchRule{{exclude: &sel}}
}

// Pass returns true only if address has a given prefix
func (sel prefixSelector) Pass(addr keypair.Address) bool {
	return strings.HasPrefix(addr.String(), sel.prefix)
}

func (sel prefixSelector) rules() []searchRule {
	return []searchRule{{prefix: &sel}}
}
