// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Command nem responses for command line user interface
package main

import (
	"fmt"
	"os"
	"time"

	"runtime"

	"strings"

	"math"

	"encoding/hex"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/vanity"
	"github.com/urfave/cli"
)

var (
	// date stores build timestamp
	date string
	// commit stores actual commit hash
	commit string
	// version indicates actual version
	version string
)

func main() {
	app := cli.NewApp()
	app.Name = "nem"
	app.Usage = "command-line toolchain for NEM blockchain"

	if version == "" {
		app.Version = "git"
	} else {
		app.Version = fmt.Sprintf("%v (%v / %v)", version, commit, date)
	}

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
							Name:  "skip-estimate",
							Usage: "Skip the step to calculate estimation times to search",
						},
						cli.BoolFlag{
							Name:  "show-complexity",
							Usage: "Show additionally the specified search complexity",
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
		fmt.Println("----")
	}

	return nil
}

func vanityAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	num := c.Uint("number")
	if num == 0 {
		return nil
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

	if !c.Bool("skip-estimate") {
		fmt.Print("Calculate accounts rate")
		ticker := time.NewTicker(time.Second)
		go func() {
			for range ticker.C {
				fmt.Print(".")
			}
		}()
		rate := float64(0)
		res := make(chan int, runtime.NumCPU())
		for i := 0; i < cap(res); i++ {
			go func(res chan<- int) {
				res <- keypair.CountKeyPairs(3200)
			}(res)
		}
		for i := 0; i < cap(res); i++ {
			rate += float64(<-res) / 3.2
		}
		ticker.Stop()
		fmt.Printf(" %v accounts/sec\n", math.Trunc(rate))

		pbty := vanity.Probability(sel) / float64(num)
		if c.Bool("show-complexity") {
			fmt.Printf("Specified search complexity: %v\n", math.Trunc(1.0/pbty))
		}
		fmt.Printf("Estimate search times: %v (50%%), %v (80%%), %v (99.9%%)\n",
			timeInSeconds(vanity.NumberOfKeyPairs(pbty, 0.5)/rate),
			timeInSeconds(vanity.NumberOfKeyPairs(pbty, 0.8)/rate),
			timeInSeconds(vanity.NumberOfKeyPairs(pbty, 0.99)/rate))
		fmt.Println("----")
	}

	rs := make(chan keypair.KeyPair)
	for i := 0; i < runtime.NumCPU(); i++ {
		go vanity.StartSearch(ch, sel, rs)
	}

	for i := uint(0); i < num; i++ {
		if i != 0 {
			fmt.Println("----")
		}
		printAccountDetails(ch, <-rs)
		go vanity.StartSearch(ch, sel, rs)
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

// Format estimated time
func timeInSeconds(val float64) string {
	val = 1e9 * math.Trunc(val)
	if val >= math.MaxInt64 || math.IsInf(val, 0) {
		return "Inf"
	}
	return time.Duration(val).String()
}

// Pretty print account details
func printAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("Address:", pair.Address(chain).PrettyString())
	fmt.Println("Public key:", hex.EncodeToString(pair.Public))
	fmt.Println("Private key:", hex.EncodeToString(pair.Private))
}
