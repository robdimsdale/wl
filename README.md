wundergo [![Build Status](https://travis-ci.org/robdimsdale/wundergo.svg?branch=master)](https://travis-ci.org/robdimsdale/wundergo) [![Coverage Status](https://img.shields.io/coveralls/robdimsdale/wundergo.svg)](https://coveralls.io/r/robdimsdale/wundergo?branch=master)
========

Golang API client for Wunderlist

## Supported Golang versions

The code is tested against the latest versions of golang 1.2, 1.3 and 1.4.

## Getting the code

The [**develop**](https://github.com/robdimsdale/wundergo/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/robdimsdale/wundergo/tree/master) branch points to a stable commit. All tests should pass.

## Running tests

Running the tests will require [ginkgo](http://onsi.github.io/ginkgo/).

### Unit tests

In the cloned directory run the following command:

```
ginkgo
```

### Integration tests

The integration tests require the following environment variables to be set:
`WL_CLIENT_ID` and `WL_ACCESS_TOKEN`. Values for these are obtained via the method documented at https://developer.wunderlist.com/documentation/concepts/authorization.

In the cloned directory run the following command:

```
WL_CLIENT_ID=my_client_id WL_ACCESS_TOKEN=my_access_token ginkgo integration_tests
```
## Project administration

### Tracker

Find this project on tracker at https://www.pivotaltracker.com/n/projects/1235310
