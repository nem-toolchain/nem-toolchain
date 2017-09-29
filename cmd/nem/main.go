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
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "n",
					Usage: "Number of generated accounts",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "generate",
					Usage:  "Generate a new account",
					Action: generateAction,
				},
				{
					Name:   "vanity",
					Usage:  "Find vanity address by a given list of prefixes",
					Action: vanityAction,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "no-digits",
							Usage: "Digits in address are disallow",
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

	pairs := make([]keypair.KeyPair, 0)
	count := c.GlobalUint64("n")
	if count > 0 {
		for i := uint64(0); i < count; i++ {
			pairs = append(pairs, keypair.Gen())
		}
	} else {
		pairs = append(pairs, keypair.Gen())
	}

	printAccountDetails(ch, pairs...)
	return nil
}

func vanityAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var noDigitsSel vanity.Selector = vanity.TrueSelector{}
	if c.Bool("no-digits") {
		noDigitsSel = vanity.NoDigitSelector{}
	}

	var prMultiSel vanity.Selector = vanity.TrueSelector{}
	if len(c.Args()) != 0 {
		prefixes := make([]vanity.Selector, len(c.Args()))
		for i, pr := range c.Args() {
			sel, err := vanity.PrefixSelectorFrom(ch, strings.ToUpper(pr))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			prefixes[i] = sel
		}
		prMultiSel = vanity.OrMultiSelector(prefixes...)
	}

	rs := make(chan keypair.KeyPair)
	pairs := make([]keypair.KeyPair, 0)
	count := c.GlobalUint64("n")
	run := func() {
		for i := 0; i < runtime.NumCPU(); i++ {
			go vanity.Search(ch, vanity.AndMultiSelector(noDigitsSel, prMultiSel), rs)
		}
	}

	run()
	if count > 0 {
		for i := uint64(0); i < count; {
			pairs = append(pairs, <-rs)
			i++

			if i != 0 && i%uint64(runtime.NumCPU()) == 0 {
				run()
			}
		}
	} else {
		pairs = append(pairs, <-rs)
	}

	printAccountDetails(ch, pairs...)
	return nil
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

func printAccountDetails(chain core.Chain, pairs ...keypair.KeyPair) {
	for _, pair := range pairs {
		fmt.Println("Address:", pair.Address(chain).PrettyString())
		fmt.Println("Public key:", hex.EncodeToString(pair.Public))
		fmt.Println("Private key:", hex.EncodeToString(pair.Private))
		fmt.Printf("\n")
	}
}
