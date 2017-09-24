// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package vanity implements a bundle of vanity address generators.
package vanity

import (
	"github.com/r8d8/nem-toolchain/pkg/keypair"
	"golang.org/x/crypto/ed25519"
)

// ByPrefix looking for the address in accordance with the given prefix
func ByPrefix(prefix string) keypair.KeyPair {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("assert: ed25519 generate key function internal error")
	}
	return keypair.KeyPair{priv[:32], pub}
}
