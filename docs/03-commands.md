---
title: Commands
---

To view details for a `nem` command at any time use `nem help` or `nem help <command>`.

```
NAME:
   nem - command-line toolchain for NEM blockchain

USAGE:
   nem [global options] command [command options] [arguments...]

VERSION:
   snapshot (9ea1a55 / 2017-10-05T01:34:13+0200)

COMMANDS:
     account  Account related bundle of actions
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --chain CHAIN  chain id from CHAIN: [mainnet|mijin|testnet] (default: "mainnet") [%NEM_CHAIN%, %CHAIN%]
   --help, -h     show help
   --version, -v  print the version
```

## Account

Account related bundle of actions.

```
NAME:
   nem account - Account related bundle of actions

USAGE:
   nem account command [command options] [arguments...]

COMMANDS:
     generate  Generate a new account
     vanity    Find vanity address by a given list of prefixes

OPTIONS:
   --help, -h  show help
```
