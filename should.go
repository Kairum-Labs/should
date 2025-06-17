// Package should provides a fluent assertion library for Go testing with generics support.
// It offers intuitive, readable test assertions with detailed error messages and intelligent suggestions.
//
// Example usage:
//
//	import (
//		"testing"
//		"github.com/Kairum-Labs/should"
//	)
//
//	func TestExample(t *testing.T) {
//		should.Ensure(42).BeGreaterThan(t, 10)
//		should.Ensure([]int{1, 2, 3}).Contain(t, 2)
//		should.Ensure("hello").BeEqual(t, "hello")
//	}
package should

import (
	"testing"

	"github.com/Kairum-Labs/should/assert"
)

// Ensure creates a new assertion for the given value.
// This is the main entry point for all assertions in the Should library.
//
// Example:
//
//	should.Ensure(value).BeEqual(t, expected)
//	should.Ensure(slice).Contain(t, item)
//	should.Ensure(condition).BeTrue(t)
func Ensure[T any](actual T) *assert.Assertion[T] {
	return assert.Ensure(actual)
}

// AssertionConfig provides configuration options for assertions.
// It allows for custom error message
//
// Example:
//
//	should.Ensure(value).BeEqual(t, expected, should.AssertionConfig{Message: "Custom message"})
type AssertionConfig = assert.AssertionConfig

// Panic asserts that the given function panics when executed.
// If the function does not panic, the test will fail with a descriptive error message.
//
// Example:
//
//	should.Panic(t, func() {
//		panic("expected panic")
//	})
func Panic(t testing.TB, fn func(), config ...assert.AssertionConfig) {
	assert.Panic(t, fn, config...)
}

// NotPanic asserts that the given function does not panic when executed.
// If the function panics, the test will fail with details about the panic.
//
// Example:
//
//	should.NotPanic(t, func() {
//		result := safeOperation()
//		_ = result
//	})
func NotPanic(t testing.TB, fn func(), config ...assert.AssertionConfig) {
	assert.NotPanic(t, fn, config...)
}
