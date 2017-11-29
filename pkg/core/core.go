// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package core contains core domain model.
package core

import (
	"errors"
	"regexp"
	"strings"
)

// Supported predefined chains.
var (
	Mijin   = Chain{byte(0x60)}
	Mainnet = Chain{byte(0x68)}
	Testnet = Chain{byte(0x98)}
)

// Chain is the type of NEM chain.
type Chain struct {
	ID byte
}

// NewChain parses byte value into a chain.
func NewChain(v byte) (ch Chain, err error) {
	switch v {
	case byte(0x68):
		ch = Mainnet
	case byte(0x98):
		ch = Testnet
	case byte(0x60):
		ch = Mijin
	default:
		err = errors.New("core: invalid chain id")
	}
	return
}

// FromString creates a chain from a chain name.
func FromString(s string) (ch Chain, err error) {
	switch strings.ToLower(s) {
	case "mijin":
		ch = Mijin
	case "mainnet":
		ch = Mainnet
	case "testnet":
		ch = Testnet
	default:
		err = errors.New("core: invalid chain name")
	}
	return
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

// Prefix returns unique chain prefix.
func (ch Chain) Prefix() string {
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

// IsChainPrefix checks for existing chain prefixes.
func IsChainPrefix(s string) bool {
	return regexp.MustCompile(`^[MNT]`).MatchString(s)
}
