---
title: Guides
---

## New accounts

How to create a new account for testnet:

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
Address: NAU6YK-24PNZX-ADXUPG-56BXUV-II4XYO-52JXAN-47UV
Public key: aef0136c2b99a08475ad8ec3b4bb56410d034f02053674531514501d3f6975fc
Private key: ddd749e79a22d6a79fd551cbedcf99a5160f8dd5984445daccb1a1fbf9b77ea6

Address: NAQZWC-ZUDCZ5-4IRKJN-WXJ76C-HLFSP5-RWYCXN-XJCZ
Public key: 77717697183da6fa770400c8b7187cc7bdeb098dcd6c539f0d7c08d822cb44c8
Private key: fdd38df995c7fbb7b3de12613a5578f26bbe625333b538401dac4c24bab8b6f5
```

## Vanity addresses

How to find vanity address with predefined prefix `TCGQQK` for testnet:

```console
$ nem --chain testnet account vanity TCGQQK
Address: TCGQQK-N5HED6-6OQ67Z-2F7GGW-Z66DWV-BFJUW6-F5WC
Public key: c342dbf7cdd3096c4c3910c511a57049e62847dd5030c7e644bc855acc1fd626
Private key: 4e017065d62f10223b989ff3f75a845fbe3df73d6c0e6d67cc4c59bea3213002
```

If you would like to search for multiple prefixes at the same time and without digits at all:

```console
$ nem --chain testnet account vanity --no-digits TABC TACB TBAC TBCA TCAB TCBA
Address: TACBLF-CJBFVE-TPTUIP-HBIRVI-PHQKKD-OSJMOF-KGNU
Public key: 8a68fdf463b4531f409369ffe368f9d78eb5e0b713459b767fbb4c4bfd148667
Private key: ae4e943300554508d52c863329a53e40787e994e7c2733d54c378fb88421d387
```

Important notes:

1. Mainnet addresses start with `N`, Mijin - with `M`, Testnet - with `T`.
1. Second symbols are `A`, `B`, `C`, or `D` only, so for mainnet you won't find addresses that start with `NE` or `N4`.
1. The digits `0`, `1`, `8` and `9` are not part of [Base32 encoding](https://en.wikipedia.org/wiki/Base32) and therefore will not appear in any address. 
