# magic
[![Tests](https://github.com/onrik/magic/workflows/Tests/badge.svg)](https://github.com/onrik/magic/actions)
[![Coverage Status](https://coveralls.io/repos/github/onrik/magic/badge.svg?branch=master)](https://coveralls.io/github/onrik/magic?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/magic)](https://goreportcard.com/report/github.com/onrik/magic)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/onrik/magic)](https://pkg.go.dev/github.com/onrik/magic)

## Magic converter for different structs
By default:
1. Map same types with same names
1. Map slices with same types 
1. Map types to pointers and backwards (for example: int to *int) 
1. Returns error on types mismatch
1. Returns lint of unconverted fields

By options:
1. Custom converters for different types
1. Custom mapping for different fields names

## Examples
### Simple
```go
package main

import (
	"fmt"

	"github.com/onrik/magic/v2"
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

	_, err := magic.Map(user1, &user2)
	fmt.Println(err)
	fmt.Printf("%+v\n", user2)
}
```

### Custom converter
```go
package main

import (
	"fmt"
	"reflect"

	"time"

	"github.com/onrik/magic/v2"
)

func timeToUnix(from, to reflect.Value) (bool, error) {
	t, ok := from.Interface().(time.Time)
	if !ok {
		return false, nil
	}

	_, ok = to.Interface().(int64)
	if !ok {
		return false, nil
	}

	to.SetInt(t.Unix())

	return true, nil
}

type User1 struct {
	ID       int
	Name     string
	Password string
	Age      int
	Created  time.Time
}

type User2 struct {
	ID      int
	Name    string
	Created int64
}

func main() {
	user1 := User1{
		ID:       1,
		Name:     "John",
		Password: "111",
		Age:      21,
		Created:  time.Now(),
	}
	user2 := User2{}

	_, err := magic.Map(user1, &user2, magic.WithConverters(timeToUnix))
	fmt.Println(err)
	fmt.Printf("%+v\n", user2)
}
```


### Fields mapping
```go
package main

import (
	"fmt"
	"time"

	"github.com/onrik/magic/v2"
)

type User1 struct {
	ID       int
	Name     string
	Password string
	Age      int
	Created  time.Time
}

type User2 struct {
	ID         int
	Name       string
	Registered time.Time
}

func main() {
	user1 := User1{
		ID:       1,
		Name:     "John",
		Password: "111",
		Age:      21,
		Created:  time.Now(),
	}
	user2 := User2{}

	_, err := magic.Map(user1, &user2, magic.WithMapping(map[string]string{
		"Created": "Registered",
	}))
	fmt.Println(err)
	fmt.Printf("%+v\n", user2)
}
```

