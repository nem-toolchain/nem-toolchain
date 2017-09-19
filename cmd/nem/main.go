package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/spin"
	"github.com/r8d8/nem-toolchain"
	"github.com/urfave/cli"
	"time"
)

var version = "master"

func main() {
	var chainStr string
	app := cli.NewApp()
	app.Name = "nem-cli"
	app.Version = version
	app.Author = "dubunda"
	app.Usage = "Vanity account generator for NEM"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "chain",
			Value:       "0x68",
			Usage:       "chain id",
			Destination: &chainStr,
		},
	}

	app.Action = func(c *cli.Context) error {
		spin := spin.New("\033[36m %s Working...\033[m")
		spin.Start()

		chainId, err := vanity.ToChainId(chainStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		acc, err := vanity.GenerateAccount(chainId)
		spin.Stop()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Account: ", acc)

		return nil
	}
	_ = app.Run(os.Args)
}
