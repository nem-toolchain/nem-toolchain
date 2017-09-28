// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Command nem responses for command line user interface
package main

import (
	"fmt"
	"os"

	"encoding/hex"

	"runtime"

	"strings"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
	"github.com/r8d8/nem-toolchain/pkg/vanity"
	"github.com/urfave/cli"
)

var (
	// BuildTime stores build timestamp
	BuildTime = "undefined"
	// CommitHash stores actual commit hash
	CommitHash = "undefined"
	// Version indicates actual version
	Version = "undefined"
)

func main() {
	app := cli.NewApp()
	app.Name = "nem"
	app.Usage = "command-line toolchain for Nem blockchain"
	app.Version = fmt.Sprintf("%v (%v / %v)", Version, CommitHash, BuildTime)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "chain",
			Value:  "mainnet",
			EnvVar: "NEM_CHAIN,CHAIN",
			Usage:  "chain id from `CHAIN`: [mainnet|mijin|testnet]",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "account",
			Usage: "Account related bundle of actions",
			Subcommands: []cli.Command{
				{
					Name:   "generate",
					Usage:  "Generate a new account",
					Action: generateAction,
				},
				{
					Name:   "vanity",
					Usage:  "Find vanity address by given prefix",
					Action: vanityAction,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:   "no-digits",
							EnvVar: "NEM_NO_DIGIT",
							Usage:  "disallow digits in account",
						},
						cli.StringSliceFlag{
							Name:   "any-pos",
							Value:  nil,
							EnvVar: "NEM_ANY_POS",
							Usage:  "list of prefixes for vanity search",
						},
					},
				},
			},
		},
	}

	_ = app.Run(os.Args)
}

func generateAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	printAccountDetails(ch, keypair.Gen())
	return nil
}

func vanityAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var pr string
	var prefixes []string
	switch len(c.Args()) {
	case 0:
		return cli.NewExitError("wrong args - prefix is not specified", 1)
	case 1:
		pr := strings.ToUpper(c.Args().First())
		if !vanity.IsPrefixCorrect(pr) {
			return cli.NewExitError("wrong args - invalid prefix format", 1)
		}
		pr = prependPrefix(ch, pr)

	default:
		for _, pr := range c.GlobalStringSlice("any-pos") {
			if !vanity.IsPrefixCorrect(pr) {
				return cli.NewExitError("wrong args - invalid prefix format", 1)
			}
			pr = prependPrefix(ch, pr)
			prefixes = append(prefixes, pr)
		}
	}

	rs := make(chan keypair.KeyPair)
	for i := 0; i < runtime.NumCPU(); i++ {
		predicates := make([]vanity.Predicate, 0)
		addr_ch := make(chan keypair.Address, 1)

		if len(prefixes) != 0 {
			predicates = append(predicates, vanity.Predicate{
				F: func() bool {
					addr := <-addr_ch
					vanity.CheckMultPrefix(addr, prefixes)
					return true
				},
				Addr_ch: addr_ch,
			})
		} else {
			predicates = append(predicates, vanity.Predicate{
				F: func() bool {
					addr := <-addr_ch
					vanity.CheckPrefix(addr, pr)
					return true
				},
				Addr_ch: addr_ch,
			})
		}

		if c.GlobalBool("no-digits") {
			predicates = append(predicates, vanity.Predicate{
				F: func() bool {
					addr := <-addr_ch
					return vanity.CheckNoDigits(addr)
				},
				Addr_ch: addr_ch,
			})
		}

		go vanity.Search(ch, rs, predicates)
	}
	printAccountDetails(ch, <-rs)
	return nil
}

func prependPrefix(ch core.Chain, pr string) string {
	switch ch {
	case core.Mijin:
		pr = "M" + pr
	case core.Mainnet:
		pr = "N" + pr
	case core.Testnet:
		pr = "T" + pr
	}
	return pr
}

func chainGlobalOption(c *cli.Context) (core.Chain, error) {
	var ch core.Chain
	switch c.GlobalString("chain") {
	case "mijin":
		ch = core.Mijin
	case "mainnet":
		ch = core.Mainnet
	case "testnet":
		ch = core.Testnet
	default:
		return ch, fmt.Errorf("unknown chain '%v'", c.GlobalString("chain"))
	}
	return ch, nil
}

func printAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("Address:", pair.Address(chain).PrettyString())
	fmt.Println("Public key:", hex.EncodeToString(pair.Public))
	fmt.Println("Private key:", hex.EncodeToString(pair.Private))
}
