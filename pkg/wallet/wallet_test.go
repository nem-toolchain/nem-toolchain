// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"testing"

	"encoding/hex"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestWalletEncode(t *testing.T) {
	exp := "eyJwcml2YXRlS2V5IjoiIiwibmFtZSI6Im1haW5uZXQiLCJhY2N" +
		"vdW50cyI6eyIwIjp7ImJyYWluIjpmYWxzZSwiYWxnbyI6InBhc3M6ZW5j" +
		"IiwiZW5jcnlwdGVkIjoiZTczZTVlZGFhYzgzOTMzODFhYTFlNWEyN2I3MWJi" +
		"Y2Q1ODM2ZGY5M2NjZDYwZGMxMTZjOGVjMGI1M2Y0NGQwZTRiZDg0NzJiYWEyM" +
		"jcyOTcyNjFmNzM4YzY1NjNlNDNkIiwiaXYiOiIxOTBjODVmZjFlNGExNTI2MmZ" +
		"mOTE3YjgyZDVlOWQ4YyIsImFkZHJlc3MiOiJORExYUzJYSUFWT1BPVkhTVVpJ" +
		"M041VlU0SEo2RU5UMjRRVklHQVBNIiwibGFiZWwiOiJQcmltYXJ5IiwibmV0d" +
		"29yayI6MTA0LCJjaGlsZCI6IjYxM2QwMWNlNjJlNDNjYzViZWE5Mzk1ZTBkOT" +
		"c5NDJjNDVkNjYxZTA4MTE4NDI0NWRkZWJhZTRlOTc3ZjMzNmYifX19"
	wlt, _ := Decode(exp)
	ser, _ := Encode(wlt)

	assert.Equal(t, ser, exp)
}

func TestWalletDecode(t *testing.T) {
	exp := "eyJwcml2YXRlS2V5IjoiIiwibmFtZSI6Im1haW5uZXQiLCJhY2N" +
		"vdW50cyI6eyIwIjp7ImJyYWluIjpmYWxzZSwiYWxnbyI6InBhc3M6ZW5j" +
		"IiwiZW5jcnlwdGVkIjoiZTczZTVlZGFhYzgzOTMzODFhYTFlNWEyN2I3MWJi" +
		"Y2Q1ODM2ZGY5M2NjZDYwZGMxMTZjOGVjMGI1M2Y0NGQwZTRiZDg0NzJiYWEyM" +
		"jcyOTcyNjFmNzM4YzY1NjNlNDNkIiwiaXYiOiIxOTBjODVmZjFlNGExNTI2MmZ" +
		"mOTE3YjgyZDVlOWQ4YyIsImFkZHJlc3MiOiJORExYUzJYSUFWT1BPVkhTVVpJ" +
		"M041VlU0SEo2RU5UMjRRVklHQVBNIiwibGFiZWwiOiJQcmltYXJ5IiwibmV0d" +
		"29yayI6MTA0LCJjaGlsZCI6IjYxM2QwMWNlNjJlNDNjYzViZWE5Mzk1ZTBkOT" +
		"c5NDJjNDVkNjYxZTA4MTE4NDI0NWRkZWJhZTRlOTc3ZjMzNmYifX19"
	wlt, _ := Decode(exp)

	assert.Equal(t, wlt.Chain, core.Mainnet)
	assert.Equal(t, wlt.PrivateKey, "")

	account := wlt.Accounts["0"]
	expAddress, _ := keypair.ParseAddress("NDLXS2XIAVOPOVHSUZI3N5VU4HJ6ENT24QVIGAPM")
	assert.Equal(t, account.Address, expAddress)
	assert.Equal(t, account.Label, "Primary")

	expChild, _ := hex.DecodeString("613d01ce62e43cc5bea9395e0d97942c45d661e081184245ddebae4e977f336f")
	assert.Equal(t, account.Child, expChild)

	expEncrypted, _ := hex.DecodeString("e73e5edaac8393381aa1e5a27b71bbcd5836df93ccd60dc116c8ec0b53f44d0e4bd8472baa227297261f738c6563e43d")
	assert.Equal(t, account.Encrypted, expEncrypted)

	expIv, _ := hex.DecodeString("190c85ff1e4a15262ff917b82d5e9d8c")
	assert.Equal(t, account.Iv, expIv)
	assert.Equal(t, account.Network, core.Mainnet)
	assert.Equal(t, account.Algo, "pass:enc")
	assert.Equal(t, account.Brain, false)
}
