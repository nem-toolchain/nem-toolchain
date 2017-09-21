// Command nem ...
package main

import (
	"fmt"
	"os"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
	"github.com/urfave/cli"
	"runtime"
)

var version = "snapshot"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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
			chainId = core.MijinId
		case "mainnet", "main", "0x68", "68":
			chainId = core.MainnetId
		case "testnet", "test", "0x98", "98":
			chainId = core.TestnetId
		default:
			panic("Unknown chain")
		}

		acc, err := keypair.GenAddress(chainId)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Account: ", acc)
		return nil
	}

	_ = app.Run(os.Args)
}
