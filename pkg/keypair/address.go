// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package keypair

import (
	"encoding/base32"
	"errors"
	"regexp"
	"strings"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
)

const (
	// AddressBytes stores the address length in bytes
	AddressBytes = 25
	// AddressLength stores the address length as string representation
	AddressLength = 40
	// PrettyAddressLength stores the address length as pretty string representation
	PrettyAddressLength = 46
)

// ParseAddress constructs an instance of `Address` from given base32 string representation
func ParseAddress(str string) (Address, error) {
	var addr Address
	str = strings.Replace(str, "-", "", -1)
	if !core.IsChainPrefix(str) {
		return addr, errors.New("unknown chain")
	}
	b, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return addr, errors.New("can't decode address string")
	}
	copy(addr[:], b)
	return addr, nil
}

// Address is a readable string representation for a public key.
type Address [AddressBytes]byte

// PrettyString returns pretty formatted address with separators ('-').
func (addr Address) PrettyString() string {
	str := addr.String()
	els := regexp.MustCompile(`.{6}`).FindAllString(str, -1)
	return strings.Join(append(els, str[36:]), "-")
}

func (addr Address) String() string {
	return base32.StdEncoding.EncodeToString(addr[:])
}
