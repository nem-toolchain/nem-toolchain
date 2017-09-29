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

	"time"

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
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:  "number, n",
							Usage: "Number of generated accounts",
							Value: 1,
						},
					},
				},
				{
					Name:   "vanity",
					Usage:  "Find vanity address by a given list of prefixes",
					Action: vanityAction,
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:  "number, n",
							Usage: "Number of generated accounts",
							Value: 1,
						},
						cli.BoolFlag{
							Name:  "no-digits",
							Usage: "Digits in address are disallow",
						},
						cli.BoolFlag{
							Name:   "skip-estimate",
							EnvVar: "NEM_SKIP_ESTIMATE",
							Usage:  "skip estimation on accounts calculation rate",
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

	num := c.Uint("number")
	for i := uint(0); i < num; i++ {
		printAccountDetails(ch, keypair.Gen())
	}

	return nil
}

func vanityAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var account_rate uint64
	if !c.Bool("skip-estimate") {
		account_rate = estimateRate()
		fmt.Printf("Estimated calculation rate - %v account/sec\n", account_rate)
	}

	var noDigitsSel vanity.Selector = vanity.TrueSelector{}
	if c.Bool("no-digits") {
		noDigitsSel = vanity.NoDigitSelector{}
	}

	var prMultiSel vanity.Selector = vanity.TrueSelector{}
	if len(c.Args()) != 0 {
		prefixes := make([]vanity.Selector, len(c.Args()))
		for i, pr := range c.Args() {
			sel, err := vanity.NewPrefixSelector(ch, strings.ToUpper(pr))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			prefixes[i] = sel
		}
		prMultiSel = vanity.OrMultiSelector(prefixes...)
	}

	sel := vanity.AndMultiSelector(noDigitsSel, prMultiSel)

	rs := make(chan keypair.KeyPair)
	for i := 0; i < runtime.NumCPU(); i++ {
		go vanity.Search(ch, sel, rs)
	}

	num := c.Uint("number")
	for i := uint(0); i < num; i++ {
		printAccountDetails(ch, <-rs)
		go vanity.Search(ch, sel, rs)
	}

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

func printAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("Address:", pair.Address(chain).PrettyString())
	fmt.Println("Public key:", hex.EncodeToString(pair.Public))
	fmt.Println("Private key:", hex.EncodeToString(pair.Private))
	fmt.Println()
}

// Estimates search rate in `account/sec`
func estimateRate() uint64 {
	test_time := 3.0
	start := time.Now()
	count := 0

	for time.Since(start).Seconds() <= test_time {
		keypair.Gen().Address(core.Testnet)
		count++
	}

	rate := count / int(test_time)
	return uint64(rate * runtime.NumCPU())
}
