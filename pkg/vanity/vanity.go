// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"strings"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
)

// FindByPrefix looking for the address in accordance with the given prefix
func FindByPrefix(chain core.Chain, prefix string, numGoroutine int) keypair.KeyPair {
	ch := make(chan keypair.KeyPair)
	for i := 0; i < numGoroutine; i++ {
		go lookByPrefix(chain, prefix, ch)
	}
	return <-ch
}

func lookByPrefix(chain core.Chain, prefix string, ch chan<- keypair.KeyPair) {
	for {
		pair := keypair.Gen()
		if checkByPrefix(chain, pair, prefix) {
			ch <- pair
			break
		}
	}
}

func checkByPrefix(chain core.Chain, pair keypair.KeyPair, prefix string) bool {
	return strings.HasPrefix(pair.Address(chain).String(), prefix)
}
