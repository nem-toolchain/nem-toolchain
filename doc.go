// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package nem is a set of packages that provide many tools for Nem blockchain.
//
// nem-toolchain library contains the following packages:
//
// The core package contains core domain model.
//
// The keypair responses for account's private/public crypto keys.
//
// The vanity package implements a bundle of vanity address generators.
package nem

// blank imports help docs.
import (
	// core package
	_ "github.com/r8d8/nem-toolchain/pkg/core"
	// keypair package
	_ "github.com/r8d8/nem-toolchain/pkg/keypair"
	// vanity package
	_ "github.com/r8d8/nem-toolchain/pkg/vanity"
)
