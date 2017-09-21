// Package nem is a set of packages that provide many tools for Nem blockchain.
//
// nem-toolchain library contains the following packages:
//
// The core package contains core domain model.
//
// The keypair package responses for private, public and address account subjects.
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
