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

// Wrapper
type Predicate struct {
	F       func() bool
	Addr_ch chan<- keypair.Address
}

// Search for account that satisfies for all predicates
func Search(chain core.Chain, ch chan<- keypair.KeyPair, predicates []Predicate) {
	for {
		pair := keypair.Gen()
		addr := pair.Address(chain)
		for i, p := range predicates {
			p.Addr_ch <- addr
			if !p.F() {
				break
			}

			if i == len(predicates)-1 {
				ch <- pair
				return
			}
		}
	}
}

// CheckPrefix checks if address satisfies prefix
func CheckPrefix(addr keypair.Address, prefix string) bool {
	return strings.HasPrefix(addr.String(), prefix)
}

// CheckMultPrefix check is address satisfies any of prefixes
func CheckMultPrefix(addr keypair.Address, prefixes []string) bool {
	for _, p := range prefixes {
		if CheckPrefix(addr, p) {
			return true
		}
	}
	return false
}

// CheckNoDigits check address for any digit
func CheckNoDigits(addr keypair.Address) bool {
	return !strings.ContainsAny(addr.String(), "2 3 4 5 6 7")
}

func IsPrefixCorrect(prefix string) bool {
	return regexp.MustCompile("^[A-D][A-Z2-7]*$").MatchString(prefix)
}
