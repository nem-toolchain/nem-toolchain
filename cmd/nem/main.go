// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Command nem responses for command line user interface
package main

import (
	"fmt"
	"os"

	"runtime"

	"encoding/hex"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
	"github.com/urfave/cli"
)

const version = "snapshot"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cli.NewApp()
	app.Name = "nem"
	app.Usage = "command-line toolchain for Nem blockchain"
	app.Version = version
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
			Name:    "account",
			Aliases: []string{"a"},
			Usage:   "Account related bundle of actions",
			Subcommands: []cli.Command{
				{
					Name:    "generate",
					Aliases: []string{"g"},
					Usage:   "Generate a new account",
					Action: func(c *cli.Context) error {
						var chain core.Chain
						switch c.GlobalString("chain") {
						case "mijin":
							chain = core.Mijin
						case "mainnet":
							chain = core.Mainnet
						case "testnet":
							chain = core.Testnet
						default:
							return cli.NewExitError("unknown chain", 1)
						}

						pair := keypair.Gen()
						fmt.Println("Address:", pair.Address(chain).PrettyString())
						fmt.Println("Public key:", hex.EncodeToString(pair.Public))
						fmt.Println("Private key:", hex.EncodeToString(pair.Private))

						return nil
					},
				},
			},
		},
	}

	_ = app.Run(os.Args)
}
