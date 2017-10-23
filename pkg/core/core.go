// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package core contains core domain model.
package core

import (
	"errors"
	"regexp"
)

// Chain is the type of NEM chain.
type Chain struct {
	Id byte
}

// Supported predefined chains.
var (
	Mijin   = Chain{byte(0x60)}
	Mainnet = Chain{byte(0x68)}
	Testnet = Chain{byte(0x98)}
)

// Parse byte value into Chain
func NewChain(val byte) (Chain, error) {
	var ch Chain
	var err error

	switch val {
	case byte(0x68):
		ch = Mainnet
	case byte(0x98):
		ch = Testnet
	case byte(0x60):
		ch = Mijin
	default:
		err = errors.New("invalid chain id")
	}

	return ch, err
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

// IsChainPrefix checks for existing chain prefixes
func IsChainPrefix(str string) bool {
	return regexp.MustCompile(`^[MNT]`).MatchString(str)
}

func (ch Chain) String() string {
	switch ch {
	case Mijin:
		return "mijin"
	case Mainnet:
		return "mainnet"
	case Testnet:
		return "testnet"
	}
	panic("unknown chain")
}
