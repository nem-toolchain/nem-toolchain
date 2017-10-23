// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKey_EncryptDecrypt(t *testing.T) {
	pr, _ := hex.DecodeString("e3da763bde538be99e748733e24540379569d878d5678d0ff647dc4edf72cf0b")
	enc, _ := hex.DecodeString("e73e5edaac8393381aa1e5a27b71bbcd5836df93ccd60dc116c8ec0b53f44d0e4bd8472baa227297261f738c6563e43d")
	pass := "12345"

	acc := Account{}
	kp := keypair.KeyPair{}
	kp.Private = pr

	assert.NoError(t, acc.Encrypt(kp, pass))
	assert.Equal(t, acc.Encrypted, enc)

	decr, _ := acc.Decrypt(pass)
	assert.Equal(t, kp.Private, decr.Private)

}
