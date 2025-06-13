# Should - A Go Assertion Library

[![go](https://img.shields.io/badge/go-1.24-blue)](https://golang.com/)
[![codecov](https://codecov.io/gh/Kairum-Labs/should/branch/main/graph/badge.svg)](https://codecov.io/gh/Kairum-Labs/should)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kairum-Labs/should)](https://goreportcard.com/report/github.com/Kairum-Labs/should)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`Should` is a lightweight, intuitive, and fluent assertion library for Go, designed to make your tests more readable and expressive. It provides clear, detailed error messages to help you debug failures faster.

## Features

- **Fluent API**: Chain assertions in a natural, readable way.
- **Detailed Diffs**: Get clear output for struct and slice comparisons.
- **Extensible**: Easily add your own custom assertions.
- **String Similarity**: When a string assertion fails, `Should` suggests similar strings from your collection to help you spot typos.
- **Integer Context**: When an integer assertion fails, `Should` shows the nearest values to help you understand the context.
- **Type-Safe**: Uses Go generics for type safety while maintaining a clean API.
- **High-Performance**: Optimized implementations for common types and operations.

## Installation

```bash
go get github.com/Andrei-hub11/should
```

## Usage

Here's how you can use `Should` in your tests:

```go
package main

import (
	"testing"
	"github.com/Andrei-hub11/should"
)

func TestUser(t *testing.T) {
	// Simple true/false assertions
	should.Ensure(true).BeTrue(t)
	should.Ensure(false).BeFalse(t)

	// Equality checks for structs, slices, and maps
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "John", Age: 30}
	p2 := Person{Name: "John", Age: 30}
	should.Ensure(p1).BeEqual(t, p2)

	// Slice contains/not contains
	users := []string{"Alice", "Bob", "Charlie"}
	should.Ensure(users).Contain(t, "Alice")
	should.Ensure(users).NotContain(t, "David")

	// Integer slice contains
	numbers := []int{1, 2, 4, 5}
	should.Ensure(numbers).Contain(t, 2) // Passes
	// should.Ensure(numbers).Contain(t, 3) // Would fail with context about where 3 would fit
}
```

## Advanced Examples

### Working with Slices and Maps

```go
// Advanced slice operations
people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
    {Name: "Charlie", Age: 35},
}

// Custom predicate functions for more complex assertions
should.Ensure(people).ContainFunc(t, func(item any) bool {
    person, ok := item.(Person)
    if !ok {
        return false
    }
    return person.Age > 30
})

// Map assertions
userMap := map[string]int{"Alice": 25, "Bob": 30, "Charlie": 35}
should.Ensure(userMap).BeEqual(t, map[string]int{"Alice": 25, "Bob": 30, "Charlie": 35})
```

### Detailed Error Messages

When an assertion fails, `Should` provides detailed error messages:

```go
// When this fails, you get a detailed diff showing exactly what's different
p1 := Person{Name: "John", Age: 30}
p2 := Person{Name: "Jane", Age: 25}
should.Ensure(p1).BeEqual(t, p2)

// Output:
// Differences found:
// Not equal:
// expected: {Name: "Jane", Age: 25}
// actual  : {Name: "John", Age: 30}
//
// Field differences:
//   └─ Name: "Jane" ≠ "John"
//   └─ Age: 25 ≠ 30
```

### String Similarity Detection

When checking for strings in slices, `Should` helps you find typos:

```go
users := []string{"user-one", "user_two", "UserThree", "user-3", "userThree"}
should.Ensure(users).Contain(t, "user3")

// Output includes helpful suggestions:
// Collection: [user-one, user_two, UserThree, user-3, userThree]
// Missing   : user3
// Similar elements found:
// └─ user-3 (at index 3) - 1 extra char
// └─ userThree (at index 4) - 5 char diff
```

### Integer Context Information

When checking for integers in slices, `Should` shows where the value would fit:

```go
numbers := []int{1, 2, 4, 5, 7, 10}
should.Ensure(numbers).Contain(t, 6)

// Output includes context information:
// Collection: [..., 4, 5, 7, 10]
// Missing   : 6
```

