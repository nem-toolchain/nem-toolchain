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

// NewChain parse byte value into Chain
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

// FromString create Chain from chain name
func FromString(name string) (Chain, error) {
	var ch Chain
	var err error

	switch strings.ToLower(name) {
	case "mainnet":
		ch = Mainnet
	case "testnet":
		ch = Testnet
	case "mijin":
		ch = Mijin
	default:
		err = errors.New("invalid chain name")
	}

	return ch, err
}

// IsChainPrefix checks for existing chain prefixes
func IsChainPrefix(str string) bool {
	return regexp.MustCompile(`^[MNT]`).MatchString(str)
}

// Chain is the type of NEM chain.
type Chain struct {
	ID byte
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
