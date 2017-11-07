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
	privBytes, _ := hex.DecodeString("e3da763bde538be99e748733e24540379569d878d5678d0ff647dc4edf72cf0b")
	pass := "12345"

	acc := Account{}
	kp := keypair.FromSeed(privBytes)

	assert.NoError(t, acc.Encrypt(kp, pass))
	//assert.Equal(t, len(acc.Encrypted), 48)

	decr, _ := acc.Decrypt(pass)
	assert.Equal(t, kp.Private, decr.Private)

}
