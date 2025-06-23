# Should - A Go Assertion Library

[![go](https://img.shields.io/badge/go-1.24-blue)](https://golang.com/)
[![codecov](https://codecov.io/gh/Kairum-Labs/should/branch/main/graph/badge.svg)](https://codecov.io/gh/Kairum-Labs/should)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kairum-Labs/should)](https://goreportcard.com/report/github.com/Kairum-Labs/should)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`Should` is a lightweight and intuitive assertion library for Go, designed to make your tests more readable and expressive. It provides **exceptionally detailed error messages** to help you debug failures faster and understand exactly what went wrong.

## Features

- **Detailed Error Messages**: Get comprehensive, contextual error information for every assertion type.
- **Smart String Handling**: Automatic multiline formatting for long strings and truncation with context.
- **Numeric Comparisons**: Detailed difference calculations with helpful hints for numeric assertions.
- **Empty/Non-Empty Checks**: Rich context about collection types, sizes, and content.
- **String Similarity**: When a string assertion fails, `Should` suggests similar strings from your collection to help you spot typos.
- **Integer Context**: When an integer assertion fails, `Should` shows the nearest values to help you understand the context.
- **Type-Safe**: Uses Go generics for type safety while maintaining a clean API.

## Installation

```bash
go get github.com/Kairum-Labs/should
```

## Quick Start

```go
package main

import (
	"testing"
	"github.com/Kairum-Labs/should"
)

func TestBasicAssertions(t *testing.T) {
	// Boolean assertions
	should.BeTrue(t, true)
	should.BeFalse(t, false)

	// Equality checks
	should.BeEqual(t, "hello", "hello")
	should.BeEqual(t, 42, 42)

	// Numeric comparisons
	should.BeGreaterThan(t, 10, 5)
	should.BeLessThan(t, 3, 7)

	// Numeric comparisons with custom messages
	should.BeGreaterThan(t, user.Age, 18, should.WithMessage("User must be adult"))
	should.BeGreaterOrEqualThan(t, score, 0, should.WithMessage("Score cannot be negative"))

	// Empty/Non-empty checks
	should.BeEmpty(t, "")
	should.BeNotEmpty(t, []int{1, 2, 3})

	// Collection operations
	users := []string{"Alice", "Bob", "Charlie"}
	should.Contain(t, users, "Alice")
	should.NotContain(t, users, "David")
	should.Contain(t, userIDs, targetID, should.WithMessage("User ID must exist in the system"))
}
```

## Detailed Error Messages

### Empty/Non-Empty Assertions

`Should` provides rich context for empty and non-empty checks:

```go
// Short string
should.BeEmpty(t, "Hello World!")
// Output:
// Expected value to be empty, but it was not:
//         Type    : string
//         Length  : 12 characters
//         Content : "Hello World!"

// Long string (automatically formatted)
longText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit..."
should.BeEmpty(t, longText)
// Output:
// Length: 516 characters, 9 lines
// 1. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// 2.  Sed do eiusmod tempor incididunt ut labore et dolore ma
// 3. gna aliqua. Ut enim ad minim veniam, quis nostrud exerci
// 4. tation ullamco laboris nisi ut aliquip ex ea commodo con
// 5. sequat. Duis aute irure dolor in reprehenderit in volupt
//
// Last lines:
// 7. xcepteur sint occaecat cupidatat non proident, sunt in c
// 8. ulpa qui officia deserunt mollit anim id est laborum. Vi
// 9. vamus sagittis lacus vel augue laoreet rutrum faucibus d

// Large slice (shows truncated content)
largeSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
should.BeEmpty(t, largeSlice)
// Output:
// Expected value to be empty, but it was not:
//         Type    : []int
//         Length  : 15 elements
//         Content : [1, 2, 3, ...] (showing first 3 of 15)

// Empty slice
should.BeNotEmpty(t, []int{})
// Output:
// Expected value to be not empty, but it was empty:
//         Type    : []int
//         Length  : 0 elements
```

### Numeric Comparisons

Get detailed information about numeric comparison failures:

```go
// Basic comparison with custom message
should.BeGreaterThan(t, 5, 10, should.WithMessage("Score validation failed"))
// Output:
// Score validation failed
// Expected value to be greater than threshold:
//         Value     : 5
//         Threshold : 10
//         Difference: -5 (value is 5 smaller)
//         Hint   : Value should be larger than threshold

// Equal values
should.BeGreaterThan(t, 42, 42)
// Output:
// Expected value to be greater than threshold:
//         Value     : 42
//         Threshold : 42
//         Difference: 0 (values are equal)
//         Hint   : Value should be larger than threshold

// Float precision
should.BeLessThan(t, 3.14, 2.71)
// Output:
// Expected value to be less than threshold:
//         Value     : 3.14
//         Threshold : 2.71
//         Difference: +0.43000000000000016 (value is 0.43000000000000016 greater)
//         Hint   : Value should be smaller than threshold

// Large numbers
should.BeLessThan(t, 1000000, 999999)
// Output:
// Expected value to be less than threshold:
//         Value     : 1000000
//         Threshold : 999999
//         Difference: +1 (value is 1 greater)
//         Hint   : Value should be smaller than threshold
```

### Struct and Object Comparisons

When comparing complex objects, `Should` shows exactly what differs:

```go
type Person struct {
    Name string
    Age  int
}

p1 := Person{Name: "John", Age: 30}
p2 := Person{Name: "Jane", Age: 25}
should.BeEqual(t, p1, p2)

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

### Length and Type Assertions

Get clear feedback on length and type mismatches.

```go
// Incorrect length
should.HaveLength(t, []string{"apple", "banana"}, 3)
// Output:
// Expected collection to have specific length:
// Type          : []string
// Expected Length: 3
// Actual Length : 2
// Difference    : -1 (1 element(s) missing)

// Incorrect type
type Dog struct{ Name string }
type Cat struct{ Name string }
var d Dog
should.BeOfType(t, Cat{Name: "Whiskers"}, d)
// Output:
// Expected value to be of specific type:
// Expected Type: should_test.Dog
// Actual Type  : should_test.Cat
// Difference   : Different concrete types
// Value        : {Name: "Whiskers"}
```

### String Similarity Detection

When checking for strings in slices, `Should` helps you find typos:

```go
users := []string{"user-one", "user_two", "UserThree", "user-3", "userThree"}
should.Contain(t, users, "user3")

// Output includes helpful suggestions:
// Expected collection to contain element:
//         Collection: [user-one, user_two, UserThree, user-3, userThree]
//         Missing   : user3
//
//         💡 Similar elements found:
//           └─ user-3 (at index 3) - 1 extra char
//           └─ userThree (at index 4) - case difference
```

### Integer Context Information

When checking for integers in slices, `Should` shows where the value would fit:

```go
numbers := []int{1, 2, 4, 5, 7, 10}
should.Contain(t, numbers, 6)

// Output includes context information:
// Collection: [..., 4, 5, 7, 10] (showing a window of 6 elements)
// Missing   : 6
```

### Set Membership Assertions

Check if a value is part of a set of allowed options.

```go
should.BeOneOf(t, "pending", []string{"active", "inactive", "suspended"})
// Output:
// Expected value to be one of the allowed options:
// Value   : "pending"
// Options : ["active", "inactive", "suspended"]
// Count   : 0 of 3 options matched
```

## API Reference

### Core Assertions

- `BeTrue(t, actual)` / `BeFalse(t, actual)` - Boolean value checks
- `BeEqual(t, actual, expected)` - Deep equality comparison with detailed diffs
- `BeNil(t, actual)` / `BeNotNil(t, actual)` - Nil pointer checks
- `BeOfType(t, actual, expected)` - Checks if a value is of a specific type
- `HaveLength(t, collection, length)` - Checks if a collection has a specific length

### Empty/Non-Empty Checks

- `BeEmpty(t, actual)` - Checks if strings, slices, arrays, maps, channels, or pointers are empty
- `BeNotEmpty(t, actual)` - Checks if values are not empty

### Numeric Comparisons

- `BeGreaterThan(t, actual, threshold)` - Numeric greater-than comparison
- `BeLessThan(t, actual, threshold)` - Numeric less-than comparison
- `BeGreaterOrEqualThan(t, actual, threshold)` - Numeric greater-than-or-equal comparison

### Collection Operations

- `BeOneOf(t, actual, options)` - Check if a value is one of a set of options
- `Contain(t, collection, element)` - Check if slice/array contains an element
- `NotContain(t, collection, element)` - Check if slice/array does not contain an element
- `ContainFunc(t, collection, predicate)` - Check if any element matches a custom predicate

### Panic Handling

- `Panic(t, func, config...)` - Assert that a function panics
- `NotPanic(t, func, config...)` - Assert that a function does not panic

Example with custom messages:

```go
// Assert function panics with custom message
should.Panic(t, func() {
    divide(1, 0)
}, should.WithMessage("Division by zero should panic"))

// Assert function doesn't panic with custom message
should.NotPanic(t, func() {
    user.Save()
}, should.WithMessage("Save operation should not panic"))
```

## Advanced Usage

### Functional Options for Assertions

`Should` uses functional options to provide a scalable way to configure assertions. This allows you to chain multiple configurations in a readable way.

#### Custom Messages with `WithMessage`

You can add custom messages to any assertion using `should.WithMessage()`:

```go
// Basic usage with a custom message
should.BeGreaterThan(t, user.Age, 18, should.WithMessage("User must be at least 18 years old"))

// Another example
should.BeGreaterOrEqualThan(t, account.Balance, 0, should.WithMessage("Account balance cannot be negative"))
```

### Custom Predicate Functions

```go
people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
    {Name: "Charlie", Age: 35},
}

// Find people over 30
should.ContainFunc(t, people, func(item any) bool {
    person, ok := item.(Person)
    if !ok {
        return false
    }
    return person.Age > 30
})

// With custom error message
should.ContainFunc(t, people, func(item any) bool {
    person, ok := item.(Person)
    if !ok {
        return false
    }
    return person.Age >= 65
}, should.WithMessage("No elderly users found"))
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

