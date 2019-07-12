# magic
[![Build Status](https://travis-ci.org/onrik/magic.svg?branch=master)](https://travis-ci.org/onrik/magic)
[![Coverage Status](https://coveralls.io/repos/github/onrik/magic/badge.svg?branch=master)](https://coveralls.io/github/onrik/magic?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/magic)](https://goreportcard.com/report/github.com/onrik/magic)
[![GoDoc](https://godoc.org/github.com/onrik/magic?status.svg)](https://godoc.org/github.com/onrik/magic)

## Magic converter for differect structs
By default:
1. Map same types with same names
1. Map slices with same types 
1. Map types to pointers and backwards (for example: int to *int) 
1. Returns error on types mismatch

By options:
1. Custom converters for different types
1. Custom mapping for different fields names

Examples:
```go
package main

import (
	"fmt"

	"github.com/onrik/magic"
)

type User1 struct {
	ID       int
	Name     string
	Password string
	Age      int
}

type User2 struct {
	ID   int
	Name string
	Age  *int
}

func main() {
	user1 := User1{
		ID:       1,
		Name:     "John",
		Password: "111",
		Age:      21,
	}
	user2 := User2{}

	err := magic.Map(user1, &user2)
	fmt.Println(err)
	fmt.Printf("%+v\n", user2)
}
```