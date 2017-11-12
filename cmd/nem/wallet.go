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

func walletCommand() cli.Command {
	return cli.Command{
		Name:  "wallet",
		Usage: "Wallet related bundle of actions",
		Subcommands: []cli.Command{
			{
				Name:   "generate",
				Usage:  "Generate a new wallet",
				Action: encodeAction,
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
