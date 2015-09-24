wundergo
========

Wunderlist API client library and CLI, written in golang.

Copyright © 2014-2015, Robert Dimsdale. Licensed under the [MIT License](https://github.com/robdimsdale/wundergo/blob/master/LICENSE).

## Why?

- The library provides access to all of the endpoints documented in the [official API docs](https://developer.wunderlist.com/documentation/), plus additional useful methods like `inbox` and `delete-all-tasks`.
- The CLI is written in golang, which results in an statically-compiled CLI.

## Library

### Supported Golang versions

The code is tested against the latest patch versions of golang 1.2, 1.3, 1.4 and 1.5.

### Getting the code

The [**develop**](https://github.com/robdimsdale/wundergo/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/robdimsdale/wundergo/tree/master) branch points to a stable commit. All tests should pass.

## CLI Binary

A CLI is provided with support for some utility functions (e.g. list all tasks, delete all folders).

### Installation

Binaries are available on the [releases](https://github.com/robdimsdale/wundergo/releases) page for various operating systems and architectures.

Download the binary and place in the PATH.

#### OSX

A [homebrew tap](https://github.com/robdimsdale/homebrew-tap) is available; install the binary with:

```
brew tap robdimsdale/tap
brew install wundergo
```

## Development

### Go dependencies

The library has no dependencies.

The CLI binary uses [godep](https://github.com/tools/godep) to manage its dependencies; install it with:

```
go get -u github.com/tools/godep
```

And restore the dependencies with:

```
godep restore
```

This will also restore the test dependencies

### Running the tests

Running the tests will require [ginkgo](http://onsi.github.io/ginkgo/).

Execute the unit tests with:

```
./scripts/unit-tests
```

The integration tests require the following environment variables to be set:
`WL_CLIENT_ID` and `WL_ACCESS_TOKEN`. Values for these are obtained via the method documented at https://developer.wunderlist.com/documentation/concepts/authorization.

In the cloned directory run the following command:

```
WL_CLIENT_ID=my_client_id WL_ACCESS_TOKEN=my_access_token ./scripts/integration_tests
```
## Project administration

### Tracker

- Roadmap: [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1235310)
