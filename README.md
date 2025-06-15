# Should - A Go Assertion Library

[![go](https://img.shields.io/badge/go-1.24-blue)](https://golang.com/)
[![codecov](https://codecov.io/gh/Kairum-Labs/should/branch/main/graph/badge.svg)](https://codecov.io/gh/Kairum-Labs/should)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kairum-Labs/should)](https://goreportcard.com/report/github.com/Kairum-Labs/should)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`Should` is a lightweight, intuitive, and fluent assertion library for Go, designed to make your tests more readable and expressive. It provides **exceptionally detailed error messages** to help you debug failures faster and understand exactly what went wrong.

## Features

- **Fluent API**: Chain assertions in a natural, readable way.
- **Detailed Error Messages**: Get comprehensive, contextual error information for every assertion type.
- **Smart String Handling**: Automatic multiline formatting for long strings and truncation with context.
- **Numeric Comparisons**: Detailed difference calculations with helpful hints for numeric assertions.
- **Empty/Non-Empty Checks**: Rich context about collection types, sizes, and content.
- **String Similarity**: When a string assertion fails, `Should` suggests similar strings from your collection to help you spot typos.
- **Integer Context**: When an integer assertion fails, `Should` shows the nearest values to help you understand the context.
- **Type-Safe**: Uses Go generics for type safety while maintaining a clean API.
- **High-Performance**: Optimized implementations for common types and operations.

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
	should.Ensure(true).BeTrue(t)
	should.Ensure(false).BeFalse(t)

	// Equality checks
	should.Ensure("hello").BeEqual(t, "hello")
	should.Ensure(42).BeEqual(t, 42)

	// Numeric comparisons
	should.Ensure(10).BeGreaterThan(t, 5)
	should.Ensure(3).BeLessThan(t, 7)

	// Empty/Non-empty checks
	should.Ensure("").BeEmpty(t)
	should.Ensure([]int{1, 2, 3}).BeNotEmpty(t)

	// Collection operations
	users := []string{"Alice", "Bob", "Charlie"}
	should.Ensure(users).Contain(t, "Alice")
	should.Ensure(users).NotContain(t, "David")
}
```

## Detailed Error Messages

### Empty/Non-Empty Assertions

`Should` provides rich context for empty and non-empty checks:

```go
// Short string
should.Ensure("Hello World!").BeEmpty(t)
// Output:
// Expected value to be empty, but it was not:
//         Type    : string
//         Length  : 12 characters
//         Content : "Hello World!"

// Long string (automatically formatted)
longText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit..."
should.Ensure(longText).BeEmpty(t)
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
should.Ensure(largeSlice).BeEmpty(t)
// Output:
// Expected value to be empty, but it was not:
//         Type    : []int
//         Length  : 15 elements
//         Content : [1, 2, 3, ...] (showing first 3 of 15)

// Empty slice
should.Ensure([]int{}).BeNotEmpty(t)
// Output:
// Expected value to be not empty, but it was empty:
//         Type    : []int
//         Length  : 0 elements
```

### Numeric Comparisons

Get detailed information about numeric comparison failures:

```go
// Basic comparison with custom message
should.Ensure(5).BeGreaterThan(t, 10, "Score validation failed")
// Output:
// Score validation failed
// Expected value to be greater than threshold:
//         Value     : 5
//         Threshold : 10
//         Difference: -5 (value is 5 smaller)
//         💡 Hint   : Value should be larger than threshold

// Equal values
should.Ensure(42).BeGreaterThan(t, 42)
// Output:
// Expected value to be greater than threshold:
//         Value     : 42
//         Threshold : 42
//         Difference: 0 (values are equal)
//         💡 Hint   : Value should be larger than threshold

// Float precision
should.Ensure(3.14).BeLessThan(t, 2.71)
// Output:
// Expected value to be less than threshold:
//         Value     : 3.14
//         Threshold : 2.71
//         Difference: +0.43000000000000016 (value is 0.43000000000000016 greater)
//         💡 Hint   : Value should be smaller than threshold

// Large numbers
should.Ensure(1000000).BeLessThan(t, 999999)
// Output:
// Expected value to be less than threshold:
//         Value     : 1000000
//         Threshold : 999999
//         Difference: +1 (value is 1 greater)
//         💡 Hint   : Value should be smaller than threshold
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
should.Ensure(numbers).Contain(t, 6)

// Output includes context information:
// Collection: [..., 4, 5, 7, 10] (showing a window of 6 elements)
// Missing   : 6
```

## API Reference

### Core Assertions

- `BeTrue(t)` / `BeFalse(t)` - Boolean value checks
- `BeEqual(t, expected)` - Deep equality comparison with detailed diffs
- `BeNil(t)` / `BeNotNil(t)` - Nil pointer checks

### Empty/Non-Empty Checks

- `BeEmpty(t)` - Checks if strings, slices, arrays, maps, channels, or pointers are empty
- `BeNotEmpty(t)` - Checks if values are not empty

### Numeric Comparisons

- `BeGreaterThan(t, threshold)` - Numeric greater-than comparison
- `BeLessThan(t, threshold)` - Numeric less-than comparison
- `BeGreaterOrEqualThan(t, threshold)` - Numeric greater-than-or-equal comparison

### Collection Operations

- `Contain(t, element)` - Check if slice/array contains an element
- `NotContain(t, element)` - Check if slice/array does not contain an element
- `ContainFunc(t, predicate)` - Check if any element matches a custom predicate

### Panic Handling

- `Panic(t, func)` - Assert that a function panics
- `NotPanic(t, func)` - Assert that a function does not panic

## Advanced Usage

### Custom Predicate Functions

```go
people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
    {Name: "Charlie", Age: 35},
}

// Find people over 30
should.Ensure(people).ContainFunc(t, func(item any) bool {
    person, ok := item.(Person)
    if !ok {
        return false
    }
    return person.Age > 30
})
```

### Type Safety

```go
// This won't compile - type mismatch
// should.Ensure("hello").BeGreaterThan(t, 42)

// This works - same types
should.Ensure(42).BeGreaterThan(t, 30)
should.Ensure(3.14).BeLessThan(t, 4.0)
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

