wundergo [![Build Status](https://travis-ci.org/robdimsdale/wundergo.svg?branch=master)](https://travis-ci.org/robdimsdale/wundergo) [![Coverage Status](https://img.shields.io/coveralls/robdimsdale/wundergo.svg)](https://coveralls.io/r/robdimsdale/wundergo?branch=master)
========

Golang API client for Wunderlist

## Getting the code

The [**develop**](https://github.com/robdimsdale/wundergo/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/robdimsdale/wundergo/tree/master) branch points to a stable commit. All tests should pass.

## Running tests

Running the tests will require [ginkgo](http://onsi.github.io/ginkgo/).

### Unit tests

In the cloned directory run the following command:

```
ginkgo -r
```

### Integration tests

The integration tests require the following environment variables to be set:
`WL_CLIENT_ID` and `WL_ACCESS_TOKEN`. Values for these are obtained via the method documented at https://developer.wunderlist.com/documentation/concepts/authorization.

After installing the dependencies, execute the following command in the cloned directory:

```
WL_CLIENT_ID=my_client_id WL_ACCESS_TOKEN=my_access_token ginkgo -r
```
