// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package keypair

import (
	"encoding/hex"
	"testing"

	"bytes"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

func TestKeyPair_Gen(t *testing.T) {
	pr, _ := hex.DecodeString("2c52aee96f0e30f21c86b3fab7a18e927f579618818e8148e7ded1e01875ef0b")
	exp, _ := hex.DecodeString("9d1e9d01ab916dbdde0e76ba43df2246575d637db0bca090f46c1abce19a43e3")
	act, _, _ := ed25519.GenerateKey(bytes.NewReader(pr))
	assert.Equal(t, exp, []byte(act))
}

func TestKeyPair_Address_mainnet(t *testing.T) {
	pub, _ := hex.DecodeString("9d1e9d01ab916dbdde0e76ba43df2246575d637db0bca090f46c1abce19a43e3")
	addr, _ := ParseAddress("NAKTWAOYSE5F3J2FJWOXR56UTQLIOUXRJLBJ7CBF")
	assert.Equal(t, addr, KeyPair{Public: pub}.Address(core.Mainnet))
}

func TestKeyPair_Address_testnet(t *testing.T) {
	pub, _ := hex.DecodeString("4fe5efd97360bc8a32ec105d419222eeb714e6d06fd8b895a5eedda2b0edf931")
	addr, _ := ParseAddress("TA6XFSJYZYAIYP7FL7X2RL63647FRMB65YC6CO3G")
	assert.Equal(t, addr, KeyPair{Public: pub}.Address(core.Testnet))
}

func TestParseAddress_wrongLength(t *testing.T) {
	for _, str := range []string{
		"",
		"TABC",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",   // not enough (-1)
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABC", // to much (+1)
	} {
		t.Run(str, func(t *testing.T) {
			_, err := ParseAddress(str)
			assert.Error(t, err)
		})
	}
}

func TestParseAddress_wrongEncoding(t *testing.T) {
	for _, str := range []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1",
		"TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA_",
	} {
		t.Run(str, func(t *testing.T) {
			_, err := ParseAddress(str)
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
