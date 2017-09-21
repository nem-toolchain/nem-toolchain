// Command nem ...
package main

import (
	"fmt"
	"os"

	"runtime"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
	"github.com/urfave/cli"
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
		var chain core.Chain

		switch chainStr {
		case "mijin", "0x60", "60":
			chain = core.Mijin
		case "mainnet", "main", "0x68", "68":
			chain = core.Mainnet
		case "testnet", "test", "0x98", "98":
			chain = core.Testnet
		default:
			panic("Unknown chain")
		}

		acc, err := keypair.GenAddress(chain)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Account: ", acc)
		return nil
	}

	_ = app.Run(os.Args)
}
