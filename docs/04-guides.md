---
title: Guides
---

## New accounts

How to create a new account for `testnet`:

```console
$ nem --chain testnet account generate
Address: TBYLAB-4RILJG-ZAUH65-YVTRDO-VH52P4-ZDPDO3-SWHS
Public key: bda2ddf2cdb65267aad0032435c5669f4391f02c681aed62b06762ddda97f1c1
Private key: 7d77192b3cc6f3770ffae4f7a85f3c40bca1256abfa1e59a3cf57cefd260af80
```

Instead of option you can specify the desired network via environment variables:

```console
$ NEM_CHAIN=testnet nem account generate
...
```

How to create multiple accounts:

```console
$ nem --chain testnet account generate -n 2
Address: TCH3WQ-DRDTOH-SAJK3C-ZJ6FT4-JT6INP-PAOUV2-K5LA
Public key: 208773e61dc7c53cc055af6ec1d5daad84912da80218db1fe2a26ad30ec459cc
Private key: 0ab24e580e14f0ac9a79c27f00c1ab6081ddc72adedd305476b4e695a4d3e3d7
----
Address: TDQJOD-DIBV7Q-6DATIR-IRLPNV-TK5XUT-2XD3Z6-HBEP
Public key: 9cb76dbcdeac1e2cc4ab8ae56ea8545cacac781dc333f334b0604cee70d59153
Private key: 3caf6c68dbda89f3c760261d76f83d75bf440509f8615395071f11721e498f3e
```

## Vanity address

How to find vanity address with predefined prefix `TCGQQK` for `testnet`:

```console
$ nem --chain testnet account vanity --skip-estimate TCGQQK
Address: TCGQQK-N5HED6-6OQ67Z-2F7GGW-Z66DWV-BFJUW6-F5WC
Public key: c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626
Private key: 4e017065d62f10223b989ff3f75a845fbe3df73d6c0e6d67cc4c59bea3213002
```

If you would like to search for multiple prefixes at the same time and without digits at all:

```console
$ nem --chain testnet account vanity --show-complexity --no-digits TABC TACB TBAC TBCA TCAB TCBA
Calculate accounts rate... 19320 accounts/sec
Specified search complexity: 1.203836e+06
Estimate search times: 0:0:43.189 (50%), 0:1:40.281 (80%), 0:4:46.940 (99.9%)
----
Address: TCBAFK-CHFHUL-FNTNTP-PMYOFI-ICEKIV-YGIXSD-IKLG
Public key: fcd12b631491585921eb8054280ebeaab894f391411ce1377f008fbfd21fb254
Private key: 7f4d2b4364ed8803ce565fb4ae2b4a97aac73f9c6629c44dcf7dde5c85e6a1af
```

As you can see from the last output, `nem-toolchain` can show specified search complexity
and calculates estimate times for three predefined accuracies: `50%`, `80%` and `99.9%`.

Important notes:

1. Mainnet addresses start with `N`, Mijin - with `M`, Testnet - with `T`.
1. Second symbols are `A`, `B`, `C`, or `D` only, so for mainnet you won't find addresses that start with `NE` or `N4`.
1. The digits `0`, `1`, `8` and `9` are not part of [Base32 encoding](https://en.wikipedia.org/wiki/Base32) and therefore will not appear in any address. 
