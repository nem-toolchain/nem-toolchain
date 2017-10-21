// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKey_EncryptDecrypt(t *testing.T) {
	pr, _ := hex.DecodeString("2a91e1d5c110a8d0105aad4683f962c2a56663a3cad46666b16d243174673d90")
	exp, _ := hex.DecodeString("8cd87bc513857a7079d182a6e19b370e907107d97bd3f81a85bcebcc4b5bd3b5")

	pass, _ := derive([]byte("TestTest"))
	assert.Equal(t, pass, exp)

	enc, _ := encrypt(pr, pass)
	dec, _ := decrypt(enc, pass)
	assert.Equal(t, dec, pr)
}
