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
$ nem --chain testnet account vanity --skip-estimate --no-digits TABC TACB TBAC TBCA TCAB TCBA
Address: TBACRS-TXXWHM-LZYKPI-ULCOZU-WFMVIX-UMVLYT-LMKM
Public key: bb326c920bb5b42d5d99df602bb82fdcdd922911ef3b46e73af654babba43698
Private key: 2d9a61cee0a3b210c2cde438b7f620931049320fd8dffbda04e28e9dd0fbfdef
```

Or you can go further and customize excluded characters:

```console
> nem-snapshot --chain testnet account vanity --show-complexity --exclude 246 TABC TACB TBAC TBCA TCAB TCBA
Calculate accounts rate... 18475 accounts/sec
Specified search complexity: 23620
Estimate search times: 0s (50%), 2s (80%), 5s (99.9%)
----
Address: TACBF3-5NBHOJ-75H3EL-QUJ3VZ-W3X73L-3BGENR-7YE5
Public key: b2b32e526d9105d937c731d1a0f470acfae8ce1bc22f100d1210f962557feda0
Private key: 6720d6a3d1601684c648d3f1fc2aa1109cb655c3fe0f86e458e55f8d7319b1e7
```

As you can see from the last output, `nem-toolchain` can show specified search complexity
and calculates estimate times for three predefined accuracies: `50%`, `80%` and `99.9%`.

Important notes:

1. Mainnet addresses start with `N`, Mijin - with `M`, Testnet - with `T`.
1. Second symbols are `A`, `B`, `C`, or `D` only, so for mainnet you won't find addresses that start with `NE` or `N4`.
1. The digits `0`, `1`, `8` and `9` are not part of [Base32 encoding](https://en.wikipedia.org/wiki/Base32) and therefore will not appear in any address. 
