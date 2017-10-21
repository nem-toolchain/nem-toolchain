// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
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
