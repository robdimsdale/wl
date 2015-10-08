wl
==

Unofficial Wunderlist API client library and CLI, written in golang.

Copyright © 2014-2015, Robert Dimsdale. Licensed under the [MIT License](https://github.com/robdimsdale/wl/blob/master/LICENSE).

## Why?

- The library provides access to all of the endpoints documented in the [official API docs](https://developer.wunderlist.com/documentation/), plus additional useful methods like `inbox` and `delete-all-tasks`.
- The CLI is written in golang, which results in an statically-compiled CLI.

## Library

Install the library with:

```
go get gopkg.in/robdimsdale/wl.v1
```

or

```
go get github.com/robdimsdale/wl
```

### Usage

Create an instance of `wl.Client` using e.g. `oauth.NewClient()` as follows:

```
import (
  "fmt"

  "github.com/robdimsdale/wl"
  "github.com/robdimsdale/wl/logger"
  "github.com/robdimsdale/wl/oauth"
)

func main() {
  // Ignore error
  client, _ := oauth.NewClient(
    "my_access_token",
    "my_client_id",
    wl.APIURL,
    logger.NewLogger(logger.INFO),
  )

  // Ignore error
  inbox, _ := client.Inbox()
  fmt.Printf("Inbox: %v\n", inbox)
}
```

### Supported Golang versions

The code is tested against the latest patch versions of golang 1.2, 1.3, 1.4 and 1.5.

### Branches

The [**develop**](https://github.com/robdimsdale/wl/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/robdimsdale/wl/tree/master) branch points to a stable commit. All tests should pass.

## CLI Binary

A CLI is provided with support for some utility functions (e.g. list all tasks, delete all folders).

### Installation

Binaries are available on the [releases](https://github.com/robdimsdale/wl/releases) page for various operating systems and architectures.

Download the binary and place in the PATH.

#### OSX

A [homebrew tap](https://github.com/robdimsdale/homebrew-tap) is available; install the binary with:

```
brew tap robdimsdale/tap
brew install wl
```

### Usage

Access token and client id are required. Provide them via flags (`--accessToken` and `--clientID`) or with environment variables (`WL_ACCESS_TOKEN` or `WL_CLIENT_ID`)

```
$ wl inbox --accessToken my_access_token --clientID my_client_id
id: 123456789
title: inbox
created_at: 2014-08-29T19:45:34.98Z
list_type: inbox
revision: 53538
type: list
public: false

$ WL_ACCESS_TOKEN=my_access_token WL_CLIENT_ID=my_client_id wl create-task --list-id 123456789 --title "some new title"
id: 987654321
[...]
list_id: 123456789
title: some new title
completed: false
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

### CI

- CI is performed using [Concourse](http://concourse.ci): https://concourse.robdimsdale.com/pipelines/wl
