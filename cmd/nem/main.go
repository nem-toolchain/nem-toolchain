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

	"bufio"

	"io/ioutil"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/vanity"
	"github.com/nem-toolchain/nem-toolchain/pkg/wallet"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
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
					Name:   "import",
					Usage:  "Import account from wallet file",
					Action: importAction,
				},
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
							Name:  "save-wallet, w",
							Usage: "Save to wallet file",
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
						cli.UintFlag{
							Name:  "workers, w",
							Usage: "Number of workers for generation",
							Value: uint(runtime.NumCPU()),
						},
						cli.StringFlag{
							Name:  "exclude",
							Usage: "Characters that must not be in the address",
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
						cli.BoolFlag{
							Name:  "save-wallet, w",
							Usage: "Save to wallet file",
						},
					},
				},
				{
					Name:   "info",
					Usage:  "extract info from account",
					Action: info,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "address",
							Usage: "Show address for supplied private key",
						},
						cli.BoolFlag{
							Name:  "public",
							Usage: "Show public key for supplied private key",
						},
					},
				},
			},
		},
	}

	_ = app.Run(os.Args)
}

func importAction(c *cli.Context) error {
	if len(c.Args()) > 1 {
		return cli.NewExitError(errors.New("invalid path to wallet file"), 1)
	}

	raw, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Enter password:")
	pass, err := terminal.ReadPassword(0)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	wlt, err := wallet.Deserialize(string(raw))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, acc := range wlt.Accounts {
		kp, err_dec := acc.Decrypt(string(pass))
		if err_dec != nil {
			return cli.NewExitError(err_dec.Error(), 1)
		}
		fmt.Println("Address:", acc.Address)
		fmt.Println("Private key:", hex.EncodeToString(kp.Private))
	}

	return err
}

func generateAction(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var pass string
	if c.Bool("save-wallet") {
		pass, err = requestPassword()
		if err != nil {
			return err
		}
	}

	num := c.Uint("number")
	for i := uint(0); i < num; i++ {
		kp := keypair.Gen()
		if c.Bool("save-wallet") {
			wlt, err := createWallet(ch, kp, pass)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			fmt.Printf("Generated wallet: %+v\n", wlt)
		} else {
			printAccountDetails(ch, kp)
			fmt.Println("----")
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

	var excludeSel vanity.Selector = vanity.TrueSelector{}
	if c.IsSet("exclude") {
		excludeSel, err = vanity.NewExcludeSelector(c.String("exclude"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	var noDigitsSel vanity.Selector = vanity.TrueSelector{}
	if c.Bool("no-digits") {
		noDigitsSel, _ = vanity.NewExcludeSelector("234567")
	}

	var prMultiSel vanity.Selector = vanity.TrueSelector{}
	if len(c.Args()) != 0 {
		prefixes := make([]vanity.Selector, len(c.Args()))
		for i, pr := range c.Args() {
			sel, err_pr := vanity.NewPrefixSelector(ch, strings.ToUpper(pr))
			if err_pr != nil {
				return cli.NewExitError(err_pr.Error(), 1)
			}
			prefixes[i] = sel
		}
		prMultiSel = vanity.OrSelector(prefixes...)
	}

	sel := vanity.AndSelector(excludeSel, noDigitsSel, prMultiSel)

	workers := c.Uint("workers")
	if m := uint(runtime.NumCPU()); workers == 0 || workers > m {
		workers = m
	}

	if !c.Bool("skip-estimate") {
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
		printEstimateDetails(
			vanity.Probability(sel)/float64(num), rate, c.Bool("show-complexity"))
		fmt.Println("----")
	}

	rs := make(chan keypair.KeyPair)
	for i := uint(0); i < workers; i++ {
		go vanity.StartSearch(ch, sel, rs)
	}

	var pass string
	if c.Bool("save-wallet") {
		pass, err = requestPassword()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	for i := uint(0); i < num; i++ {
		if i != 0 {
			fmt.Println("----")
		}
		if c.Bool("save-wallet") {
			wlt, err := createWallet(ch, <-rs, pass)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			fmt.Printf("Generated wallet: %+v\n", wlt)
		} else {
			printAccountDetails(ch, <-rs)
		}
		go vanity.StartSearch(ch, sel, rs)
	}

	return nil
}

func info(c *cli.Context) error {
	ch, err := chainGlobalOption(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Enter private key: ")
	reader := bufio.NewReader(os.Stdin)
	privateKeyStr, err := reader.ReadString('\n')
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	pkBytes, err := hex.DecodeString(strings.TrimSpace(privateKeyStr))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	pair := keypair.FromSeed(pkBytes)
	if c.Bool("address") {
		printAddress(ch, pair)
	}

	if c.Bool("public") {
		printPublicKey(pair)
	}

	if !c.Bool("address") && !c.Bool("public") {
		printAccountDetails(ch, pair)
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

// countActualRate counts total number of generated keypairs per second
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

// countKeyPairs counts number of generated keypairs for specified interval
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

// printEstimateDetails prints estimate search time details
func printEstimateDetails(pbty, rate float64, compl bool) {
	if compl {
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

// printAccountDetails prints account details in pretty user-oriented multi-line format
func printAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	printAddress(chain, pair)
	printPublicKey(pair)
	printPrivateKey(pair)
}

func printAddress(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("Address:", pair.Address(chain).PrettyString())
}

func printPublicKey(pair keypair.KeyPair) {
	fmt.Println("Public key:", hex.EncodeToString(pair.Public))
}

func printPrivateKey(pair keypair.KeyPair) {
	fmt.Println("Private key:", hex.EncodeToString(pair.Private))
}

func createWallet(chain core.Chain, pair keypair.KeyPair, pass string) (wallet.Wallet, error) {
	wlt := wallet.New(chain)
	err := wlt.AddAccount(pair, pass)

	return wlt, err
}

func requestPassword() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	return reader.ReadString('\n')
}
