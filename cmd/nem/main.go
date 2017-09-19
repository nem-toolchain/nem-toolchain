package main

import (
	"fmt"
	"os"

	"github.com/r8d8/nem-toolchain"
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
		chainId, err := vanity.ToChainId(chainStr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		acc, err := vanity.GenerateAccount(chainId)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Account: ", acc)
		return nil
	}

	_ = app.Run(os.Args)
}
