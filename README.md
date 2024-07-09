# TestIt

TestIt is a very simple, but powerful, testing framework.

It allows you to avoid repetition in the following:

* mock initialization
* pre-test cleanup
* panic avoidance
* test execution
* assertions

It also helps you build small blocks that can be re-utilized between test cases without sharing states. Examples:

* Setup a group of expected calls
* Shared assertions for response fields

## Usage

### tools/tools.go

Create the tools/tools.go to declare testit as a development only dependency

```go
//go:build tools

package tools

import (
	_ "github.com/sonalys/testit"
)

```

### A very simple example for Stateful and Stateless tests:

```go
```