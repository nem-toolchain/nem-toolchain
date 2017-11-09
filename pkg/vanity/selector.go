// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/util"
)

const (
	// PrefixPlaceholder declares a placeholder (skipped) rune for prefix selector
	PrefixPlaceholder = '_'
	// PrefixSeparator declares a separator (ignored) rune for prefix selector
	PrefixSeparator = '-'
)

// Selector defines generic search strategy
type Selector interface {
	// Pass checks address by a given search pattern
	Pass(addr keypair.Address) bool
	// rules converts current selector into a bundle of normalized rules
	rules() []searchRule
}

// FalseSelector rejects all addresses, can be used as a default placeholder
type FalseSelector struct{}

// Pass returns always false
func (FalseSelector) Pass(keypair.Address) bool {
	return false
}

func (FalseSelector) rules() []searchRule {
	// nothing, just skip it
	return []searchRule{}
}

// TrueSelector accepts all addresses, can be used as a default placeholder
type TrueSelector struct{}

// Pass returns always true
func (TrueSelector) Pass(keypair.Address) bool {
	return true
}

func (TrueSelector) rules() []searchRule {
	// empty searchRule - always true
	return []searchRule{{}}
}

// exclude checks an address for absence of given characters
type excludeSelector struct {
	// chars determines a unique sorted list of excluded characters
	chars string
}

// NewExcludeSelector creates a new exclude selector from given string
func NewExcludeSelector(chars string) (Selector, error) {
	if !regexp.MustCompile(`^[A-Z2-7]*$`).MatchString(chars) {
		return excludeSelector{}, fmt.Errorf("incorrect exclude characters '%v'", chars)
	}
	arr := strings.Split(chars, "")
	sort.Strings(arr)
	return excludeSelector{strings.Join(util.DistinctStrings(arr), "")}, nil
}

// Pass returns true only if address doesn't contain any digits
func (sel excludeSelector) Pass(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String(), sel.chars)
}

func (sel excludeSelector) rules() []searchRule {
	return []searchRule{{exclude: &sel}}
}

// prefix checks an address by given prefix
type prefixSelector struct {
	// prefix determines a required address prefix to search
	prefix string
	// re caches a regexp.Regexp object for given prefix
	re *regexp.Regexp
}

// NewPrefixSelector creates a new prefix selector from given string
func NewPrefixSelector(ch core.Chain, prefix string) (Selector, error) {
	prefix = strings.Replace(prefix, string(PrefixSeparator), "", -1)
	str := fmt.Sprintf(`^[_%v]?([_A-D][_A-Z2-7]*)?$`, ch.ChainPrefix())
	if !regexp.MustCompile(str).MatchString(prefix) {
		return prefixSelector{}, fmt.Errorf("incorrect prefix '%v'", prefix)
	}
	regex := regexp.MustCompile(fmt.Sprintf("^%v\\w*",
		strings.Replace(prefix, string(PrefixPlaceholder), "\\w", -1)))
	return prefixSelector{prefix, regex}, nil
}

// Pass returns true only if address has a given prefix
func (sel prefixSelector) Pass(addr keypair.Address) bool {
	return sel.re.MatchString(addr.String())
}

func (sel prefixSelector) rules() []searchRule {
	return []searchRule{{prefix: &sel}}
}
