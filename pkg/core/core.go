// Package core contains core domain model
package core

import (
	"errors"
)

// Supported chains
const (
	MijinId   = byte(0x60)
	MainnetId = byte(0x68)
	TestnetId = byte(0x98)
)

// ErrInvalidChain indicates invalid chain id.
var ErrInvalidChain = errors.New("invalid chain id")

// IsValidChainId checks chain id for existence
func IsValidChainId(id byte) bool {
	for _, i := range []byte{MijinId, MainnetId, TestnetId} {
		if i == id {
			return true
		}
	}
	return false
}
