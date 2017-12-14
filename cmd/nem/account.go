// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/vanity"
	"github.com/urfave/cli"
)

func accountCommand() cli.Command {
	return cli.Command{
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
					cli.BoolFlag{
						Name:  "strip, s",
						Usage: "Strip output to private key only",
					},
				},
			},
			{
				Name:   "info",
				Usage:  "Show info for given account",
				Action: infoAction,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "address",
						Usage: "Show public address only for given private key",
					},
					cli.BoolFlag{
						Name:  "public",
						Usage: "Show public key only for given private key",
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
						Name:  "strip, s",
						Usage: "Strip output to private key only",
					},
					cli.UintFlag{
						Name:  "workers, w",
						Usage: "Number of workers for generation",
						Value: uint(runtime.NumCPU()),
					},
					cli.BoolFlag{
						Name:  "show-complexity",
						Usage: "Show additionally the specified search complexity",
					},
					cli.BoolFlag{
						Name:  "skip-estimate",
						Usage: "Skip the step to calculate estimation times to search",
					},
					cli.BoolFlag{
						Name:  "no-digits",
						Usage: "Digits in address are disallow ('234567')",
					},
					cli.StringFlag{
						Name:  "exclude",
						Usage: "Characters that must not be in the address",
					},
				},
			},
		},
	}
}

func generateAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	num := c.Uint("number")
	for i := uint(0); i < num; i++ {
		pair := keypair.Gen()
		if c.Bool("strip") {
			printlnPrivateKey(pair, true)
		} else {
			printAccountDetails(ch, pair)
		}
	}

	return nil
}

func infoAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if fi.Mode()&os.ModeCharDevice != 0 {
		return cli.NewExitError("interactive input mode isn't supported", 1)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	s := strings.Split(string(data), "\n")
	for _, s := range s[:len(s)-1] {
		pk, err := hex.DecodeString(strings.TrimSpace(s))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		pair, err := keypair.FromSeed(pk)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if c.Bool("address") {
			printlnAddress(ch, pair, true)
		} else if c.Bool("public") {
			printlnPublicKey(pair, true)
		} else {
			printAccountDetails(ch, pair)
		}
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

	workers := c.Uint("workers")
	if m := uint(runtime.NumCPU()); workers == 0 || workers > m {
		workers = m
	}

	sel, err := createSelector(c, ch)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !c.Bool("strip") && !c.Bool("skip-estimate") {
		printEstimate(workers,
			vanity.Probability(sel)/float64(num), c.Bool("show-complexity"))
	}

	rs := make(chan keypair.KeyPair)
	for i := uint(0); i < workers; i++ {
		go vanity.StartSearch(ch, sel, rs)
	}

	for i := uint(0); i < num; i++ {
		pair := <-rs
		if c.Bool("strip") {
			printlnPrivateKey(pair, true)
		} else {
			printAccountDetails(ch, pair)
		}
		go vanity.StartSearch(ch, sel, rs)
	}

	return nil
}

// createSelector create selector for vanity search,
// specified by flags.
func createSelector(c *cli.Context, ch core.Chain) (vanity.Selector, error) {
	var excludeSel vanity.Selector = vanity.TrueSelector{}
	var err error
	if c.IsSet("exclude") {
		excludeSel, err = vanity.NewExcludeSelector(c.String("exclude"))
		if err != nil {
			return nil, err
		}
	}

	var noDigitsSel vanity.Selector = vanity.TrueSelector{}
	if c.Bool("no-digits") {
		noDigitsSel, _ = vanity.NewExcludeSelector("234567")
	}

	args, err := readArgs(c)
	if err != nil {
		return nil, err
	}
	prefixes := make([]vanity.Selector, len(args))
	for i, pr := range args {
		sel, err := vanity.NewPrefixSelector(ch, strings.ToUpper(pr))
		if err != nil {
			return nil, err
		}
		prefixes[i] = sel
	}

	sel := vanity.AndSelector(
		excludeSel, noDigitsSel, vanity.OrSelector(prefixes...))

	return sel, nil
}

// printEstimate prints vanity account search time estimate.
func printEstimate(workers uint, pbty float64, cplx bool) {
	fmt.Print("Calculate accounts rate")
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			fmt.Print(".")
		}
	}()
	rate := countActualRate(workers)
	ticker.Stop()
	fmt.Printf(" %v accounts/sec\n", math.Trunc(rate))
	printEstimateDetails(pbty, rate, cplx)
}

// countActualRate counts total number of generated keypairs per second.
func countActualRate(workers uint) float64 {
	res := make(chan int, workers)
	for i := 0; i < cap(res); i++ {
		go countKeyPairs(3200, res)
	}
	rate := float64(0)
	for i := 0; i < cap(res); i++ {
		rate += float64(<-res) / 3.2
	}
	return rate
}

// countKeyPairs counts number of generated keypairs for specified interval.
func countKeyPairs(milliseconds time.Duration, res chan int) {
	timeout := time.After(time.Millisecond * milliseconds)
	for count := 0; ; count++ {
		keypair.Gen().Address(core.Mainnet)
		select {
		case <-timeout:
			res <- count
			return
		default:
			continue
		}
	}
}

// readArgs either reads positional arguments or from stdin.
func readArgs(c *cli.Context) (cli.Args, error) {
	args := c.Args()
	if args[0] == "-" {
		reader := bufio.NewReader(os.Stdin)
		raw, err := reader.ReadString('\n')
		if err != nil {
			return args, err
		}
		raw = strings.TrimSpace(raw)
		args = cli.Args(strings.Split(raw, " "))
	}

	return args, nil
}

// printEstimateDetails prints estimate search time details.
func printEstimateDetails(pbty, rate float64, cplx bool) {
	if cplx {
		fmt.Printf("Specified search complexity: %v\n", math.Trunc(1.0/pbty))
	}
	fmt.Printf("Estimate search times: %v (50%%), %v (80%%), %v (99.9%%)\n",
		timeInSeconds(vanity.NumberOfAttempts(pbty, 0.5)/rate),
		timeInSeconds(vanity.NumberOfAttempts(pbty, 0.8)/rate),
		timeInSeconds(vanity.NumberOfAttempts(pbty, 0.99)/rate))
}

// timeInSeconds formats estimated time
func timeInSeconds(val float64) string {
	val = 1e9 * math.Trunc(val)
	if val >= math.MaxInt64 || math.IsInf(val, 0) {
		return "Inf"
	}
	return time.Duration(val).String()
}

// printAccountDetails prints account details in pretty user-oriented multi-line format.
func printAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("----")
	printlnAddress(chain, pair, false)
	printlnPublicKey(pair, false)
	printlnPrivateKey(pair, false)
}

func printlnPrivateKey(pair keypair.KeyPair, strip bool) {
	printlnCustom("Private key: ", hex.EncodeToString(pair.Private), strip)
}

func printlnPublicKey(pair keypair.KeyPair, strip bool) {
	printlnCustom("Public key: ", hex.EncodeToString(pair.Public), strip)
}

func printlnAddress(chain core.Chain, pair keypair.KeyPair, strip bool) {
	printlnCustom("Address: ", pair.Address(chain).PrettyString(), strip)
}

func printlnCustom(title, value string, strip bool) {
	if !strip {
		fmt.Print(title)
	}
	fmt.Println(value)
}
