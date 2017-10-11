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
	prv, _ := hex.DecodeString("2c52aee96f0e30f21c86b3fab7a18e927f579618818e8148e7ded1e01875ef0b")
	exp, _ := hex.DecodeString("9d1e9d01ab916dbdde0e76ba43df2246575d637db0bca090f46c1abce19a43e3")
	pub, _, _ := ed25519.GenerateKey(bytes.NewReader(prv))
	assert.Equal(t, []byte(pub), exp)
}

func TestKeyPair_Address(t *testing.T) {
	pub, _ := hex.DecodeString("c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626")
	addr, _ := FromString("TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC")
	assert.Equal(t, addr, KeyPair{Public: pub}.Address(core.Testnet))
}

func TestKeyPair_Address_ForMainnet(t *testing.T) {
	pub, _ := hex.DecodeString("9d1e9d01ab916dbdde0e76ba43df2246575d637db0bca090f46c1abce19a43e3")
	addr, _ := FromString("NAKTWAOYSE5F3J2FJWOXR56UTQLIOUXRJLBJ7CBF")
	assert.Equal(t, addr, KeyPair{Public: pub}.Address(core.Mainnet))
}

func TestKeyPair_Address_ForTestnet(t *testing.T) {
	pub, _ := hex.DecodeString("4fe5efd97360bc8a32ec105d419222eeb714e6d06fd8b895a5eedda2b0edf931")
	addr, _ := FromString("TA6XFSJYZYAIYP7FL7X2RL63647FRMB65YC6CO3G")
	assert.Equal(t, addr, KeyPair{Public: pub}.Address(core.Testnet))
}

func TestAddress_PrettyString(t *testing.T) {
	addr, _ := FromString("TCGQQKN5HED66OQ67Z2F7GGWZ66DWVBFJUW6F5WC")
	assert.Equal(t, "TCGQQK-N5HED6-6OQ67Z-2F7GGW-Z66DWV-BFJUW6-F5WC", addr.PrettyString())
}
