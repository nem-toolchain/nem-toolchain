// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"strings"

	"regexp"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
)

type Predicate interface {
	call(addr keypair.Address) bool
}

// NoDigitPredicate checks account for absence of digits
type NoDigitPredicate struct{}

// PrefixPredicate checks account for selected prefix
type PrefixPredicate struct {
	Prefix string
}

// MultPrefixPredicate checks account for any of specified prefixes
type MultPrefixPredicate struct {
	Prefixes []string
}

func (nd NoDigitPredicate) call(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String(), "2 3 4 5 6 7")
}

func (pr PrefixPredicate) call(addr keypair.Address) bool {
	return checkPrefix(addr, pr.Prefix)
}

func (mpr MultPrefixPredicate) call(addr keypair.Address) bool {
	for _, p := range mpr.Prefixes {
		if checkPrefix(addr, p) {
			return true
		}
	}
	return false
}

// Search for account that satisfies for all predicates
func Search(chain core.Chain, ch chan<- keypair.KeyPair, predicates []Predicate) {
	for {
		pair := keypair.Gen()
		addr := pair.Address(chain)
		for i, p := range predicates {
			if !p.call(addr) {
				break
			}

			if i == len(predicates)-1 {
				ch <- pair
				return
			}
		}
	}
}


func checkPrefix(addr keypair.Address, prefix string) bool {
	return strings.HasPrefix(addr.String(), prefix)
}

// IsPrefixCorrect verify that prefix can be used
func IsPrefixCorrect(prefix string) bool {
	return regexp.MustCompile("^[A-D][A-Z2-7]*$").MatchString(prefix)
}
