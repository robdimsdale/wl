wundergo [![Build Status](https://travis-ci.org/robdimsdale/wundergo.svg?branch=master)](https://travis-ci.org/robdimsdale/wundergo) [![Coverage Status](https://img.shields.io/coveralls/robdimsdale/wundergo.svg)](https://coveralls.io/r/robdimsdale/wundergo?branch=master)
========

Golang API client for Wunderlist.

Copyright © 2014-2015, Robert Dimsdale. Licensed under the [MIT License](https://github.com/robdimsdale/wundergo/blob/master/LICENSE).

## Supported Golang versions

The code is tested against the latest patch versions of golang 1.2, 1.3, 1.4 and 1.5.

## Getting the code

The [**develop**](https://github.com/robdimsdale/wundergo/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/robdimsdale/wundergo/tree/master) branch points to a stable commit. All tests should pass.

## Development

### Go dependencies

There are no dependencies for the library.

The CLI binary has one dependency; install it with:

```
go get gopkg.in/yaml.v2
```

This dependency follows semantic versioning, so all versions should be safe to `go get`.

There are dependencies for the tests; they are safe to install from HEAD of
their respective repositories and hence are not vendored in.

Test dependencies are installed as follows:

```
go get -u github.com/onsi/ginkgo
go get -u github.com/onsi/gomega
go get -u github.com/nu7hatch/gouuid
```

### Running tests

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

- Roadmap: [Pivotal Tracker(https://www.pivotaltracker.com/n/projects/1235310)
