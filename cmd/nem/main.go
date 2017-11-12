// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Command nem responses for command line user interface
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/wallet"
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
		accountCommand(),
	}

	_ = app.Run(os.Args)
}

func encodeAction(c *cli.Context) error {
	ch, err := core.FromString(c.GlobalString("chain"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	wlt := wallet.NewWallet()
	if c.IsSet("name") {
		wlt.Name = c.String("name")
	} else {
		reader := bufio.NewReader(os.Stdin)
		nameBytes, errName := reader.ReadBytes('\n')
		if errName != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		wlt.Name = string(nameBytes)
	}

	pass, err := requestPassword()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	privKeyBytes, err := requestPrivateKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	pair, err := keypair.FromSeed(privKeyBytes)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = wlt.AddAccount(ch, pair, pass)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	walletEncoded, err := wallet.Encode(wlt)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(walletEncoded)

	return nil
}

func decodeAction(c *cli.Context) error {
	pass, err := requestPassword()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	reader := bufio.NewReader(os.Stdin)
	raw, err := reader.ReadString('\n')
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	wlt, err := wallet.Decode(raw)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, acc := range wlt.Accounts {
		pair, errDec := acc.Decrypt(pass)
		if errDec != nil {
			return cli.NewExitError(errDec.Error(), 1)
		}
		printlnPrivateKey(pair, false)
	}

	return nil
}

// requestPassword request input of password
func requestPassword() (string, error) {
	var passBytes []byte
	passBytes, err := gopass.GetPasswdPrompt("Enter password: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", cli.NewExitError(err.Error(), 1)
	}
	return string(passBytes), nil
}

// requestPassword request input of private key hex-string
func requestPrivateKey() ([]byte, error) {
	var pkBytes []byte
	val, err := gopass.GetPasswdPrompt("Enter private key: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return pkBytes, cli.NewExitError(err.Error(), 1)
	}

	return keypair.HexToPrivBytes(string(val))
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
