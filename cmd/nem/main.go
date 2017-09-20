package main

import (
	"fmt"
	"os"

	"github.com/r8d8/nem-toolchain/pkg/vanity"
	"github.com/urfave/cli"
)

var version = "snapshot"

func main() {
	var chainStr string
	app := cli.NewApp()
	app.Name = "nem"
	app.Version = version
	app.Author = "dubunda"
	app.Usage = "Vanity account address generator for Nem blockchain"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "chain",
			Value:       "0x68",
			Usage:       "chain id",
			Destination: &chainStr,
		},
	}

	app.Action = func(c *cli.Context) error {
		var chainId byte

		switch chainStr {
		case "mijin", "0x60", "60":
			chainId = vanity.MijinId
		case "mainnet", "main", "0x68", "68":
			chainId = vanity.MainnetId
		case "testnet", "test", "0x98", "98":
			chainId = vanity.TestnetId
		default:
			panic("Unknown chain")
		}

		acc, err := vanity.GenAddress(chainId)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Account: ", acc)
		return nil
	}

	_ = app.Run(os.Args)
}
