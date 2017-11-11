---
title: Quick Start
---

`nem-toolchain` is distributed in a binary form and can be installed manually via the
[tarball releases](https://github.com/nem-toolchain/nem-toolchain/releases/latest)
or one of the options below.

The quickest way to get `nem` is to run the following command, which installs to the
local `bin` by default.

```console
$ curl -sL https://git.io/getnem | bash
```

Or, if you want to lock at specific version:

```console
$ curl -sL https://git.io/getnem | VERSION=vX.Y.Z bash
```

Verify installation with:

```
$ bin/nem -v
```

You can include the `bin` directory in your `PATH` to simplify further usage:

```console
export PATH=$PATH:./bin
```
