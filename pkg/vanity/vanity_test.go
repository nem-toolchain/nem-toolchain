// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"strings"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestStartSearch(t *testing.T) {
	rs := make(chan keypair.KeyPair)
	sel, _ := NewPrefixSelector(core.Testnet, "TA")
	go StartSearch(core.Testnet, sel, rs)
	pr := <-rs
	assert.True(t, strings.HasPrefix(pr.Address(core.Testnet).String(), "TA"))
}
