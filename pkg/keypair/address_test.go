// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package keypair

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress_length(t *testing.T) {
	for _, s := range []string{
		"",
		"TABC",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",   // not enough (-1)
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABC", // to much (+1)
	} {
		t.Run(s, func(t *testing.T) {
			_, err := ParseAddress(s)
			assert.Error(t, err)
		})
	}
}

func TestParseAddress_encoding(t *testing.T) {
	for _, s := range []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA_",
	} {
		t.Run(s, func(t *testing.T) {
			_, err := ParseAddress(s)
			assert.Error(t, err)
		})
	}
}

func TestParseAddress_pretty(t *testing.T) {
	addr, err := ParseAddress("TAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAA")
	assert.NoError(t, err)
	assert.Equal(t, "TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", addr.String())
}

func TestAddress_PrettyString(t *testing.T) {
	addr, _ := ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	assert.Equal(t, "TAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAAAA-AAAA", addr.PrettyString())
}

func TestAddress_String(t *testing.T) {
	addr, _ := ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	assert.Equal(t, "TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", addr.String())
}
