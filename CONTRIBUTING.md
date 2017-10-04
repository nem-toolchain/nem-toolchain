# How to contribute

If you would like to contribute code you can do so through GitHub by forking
the repository and open a pull request against the master branch.

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md).
By participating, you are expected to uphold this code.

## How to build

Prerequisites are:

* [`git`](https://git-scm.com)
* `make`
* [`Go 1.8+`](http://golang.org/doc/install)

Install the build and lint dependencies:

``` sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

``` sh
$ make test
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

``` sh
$ make build
```

Make sure your `PATH` includes the `bin` directory so your local builds can be easily used:

```bash
export PATH=$PATH:./bin
```

Simplest first test after build is:

```bash
nem -h
```

When you are satisfied with the changes, we suggest you run:

``` sh
$ make fmt && make ci
```

Which runs all the formatters, linters and tests.

## Categorized commits

All commits started with `[` will be included to release notes automatically, good candidates are:

* `[FIX]` - bug or other issue fix (non-breaking change which fixes an issue)
* `[IMPROVEMENT]` - existing feature improvement (non-breaking change which expend some functionality)
* `[FEATURE]` - new feature (non-breaking change which adds functionality)
* `[BREAKING]` - breaking change (fix or feature that would cause existing functionality to change)
* `[MISC]` - not fall under any of the above categories (non-breaking non-functional change)

## License

By contributing your code, you agree to license your contribution under the terms of the
[MIT License](LICENSE).

If you are adding a new file it should have a header like this:

```
// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
 ```
