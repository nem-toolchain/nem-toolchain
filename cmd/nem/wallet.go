// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/nem-toolchain/nem-toolchain/pkg/wallet"
	"github.com/urfave/cli"
)

func walletCommand() cli.Command {
	return cli.Command{
		Name:  "wallet",
		Usage: "Wallet related bundle of actions",
		Subcommands: []cli.Command{
			{
				Name:   "generate",
				Usage:  "Generate a new wallet",
				Action: encodeAction,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name",
						Usage: "Name of wallet",
					},
				},
			},
			{
				Name:   "read",
				Usage:  "Read wallet",
				Action: decodeAction,
			},
		},
	}
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
		fmt.Print("Enter name: ")
		reader := bufio.NewReader(os.Stdin)
		nameBytes, errName := reader.ReadBytes('\n')
		if errName != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		wlt.Name = string(nameBytes)
	}

	pass := requestPassword()
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
	pass := requestPassword()
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

// requestPrivateKey request input of private key
func requestPrivateKey() ([]byte, error) {
	pk := requestHiddenString("Enter private key: ")
	return keypair.HexToPrivBytes(pk)
}

// requestPassword request input of password
func requestPassword() string {
	return requestHiddenString("Enter password: ")
}

// requestHiddenString hides requested input
func requestHiddenString(prompt string) string {
	var password string
	fmt.Print(prompt)
	fmt.Print("\033[8m") // Hide input
	fmt.Scan(&password)
	fmt.Print("\033[28m") // Show input
	return password
}
