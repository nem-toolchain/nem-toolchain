// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package core contains core domain model.
package core

import "regexp"

// Supported predefined chains.
var (
	Mijin   = Chain{byte(0x60)}
	Mainnet = Chain{byte(0x68)}
	Testnet = Chain{byte(0x98)}
)

// IsChainPrefix checks for existing chain prefixes
func IsChainPrefix(str string) bool {
	return regexp.MustCompile(`^[MNT]`).MatchString(str)
}

// Chain is the type of NEM chain.
type Chain struct {
	Id byte
}

// ChainPrefix returns unique chain prefix
func (ch Chain) ChainPrefix() string {
	switch ch {
	case Mijin:
		return "M"
	case Mainnet:
		return "N"
	case Testnet:
		return "T"
	}
	panic("unknown chain")
}
