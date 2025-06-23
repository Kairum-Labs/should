// Package should provides a simple, elegant, and readable assertion library for Go testing.
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
//		should.BeGreaterThan(t, 42, 10)
//		should.Contain(t, []int{1, 2, 3}, 2)
//		should.BeEqual(t, "hello", "hello")
//	}
package should

import (
	"testing"

	"github.com/Kairum-Labs/should/assert"
)

// Option is a functional option for configuring assertions.
type Option = assert.Option

// WithMessage creates an option for setting a custom error message.
//
// Example:
//
//	should.BeEqual(t, value, expected, should.WithMessage("Custom message"))
func WithMessage(message string) Option {
	return assert.WithMessage(message)
}

// BeTrue reports a test failure if the value is not true.
//
// This assertion only works with boolean values and will fail immediately
// if the value is not a boolean type.
//
// Example:
//
//	should.BeTrue(t, true)
//
//	should.BeTrue(t, user.IsActive, should.WithMessage("User must be active"))
//
// If the input is not a boolean, the test fails immediately.
func BeTrue[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeTrue(t, actual, opts...)
}

// BeFalse reports a test failure if the value is not false.
//
// This assertion only works with boolean values and will fail immediately
// if the value is not a boolean type.
//
// Example:
//
//	should.BeFalse(t, false)
//
//	should.BeFalse(t, user.IsDeleted, should.WithMessage("User should not be deleted"))
//
// If the input is not a boolean, the test fails immediately.
func BeFalse[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeFalse(t, actual, opts...)
}

// BeEmpty reports a test failure if the value is not empty.
//
// This assertion works with strings, slices, arrays, maps, channels, and pointers.
// For strings, empty means zero length. For slices/arrays/maps/channels, empty means zero length.
// For pointers, empty means nil. Provides detailed error messages showing the type,
// length, and content of non-empty values.
//
// Example:
//
//	should.BeEmpty(t, "")
//
//	should.BeEmpty(t, []int{}, should.WithMessage("List should be empty"))
//
//	should.BeEmpty(t, map[string]int{})
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func BeEmpty[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeEmpty(t, actual, opts...)
}

// BeNotEmpty reports a test failure if the value is empty.
//
// This assertion works with strings, slices, arrays, maps, channels, and pointers.
// For strings, non-empty means length > 0. For slices/arrays/maps/channels, non-empty means length > 0.
// For pointers, non-empty means not nil. Provides detailed error messages for empty values.
//
// Example:
//
//	should.BeNotEmpty(t, "hello")
//
//	should.BeNotEmpty(t, []int{1, 2, 3}, should.WithMessage("List must have items"))
//
//	should.BeNotEmpty(t, &user)
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func BeNotEmpty[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeNotEmpty(t, actual, opts...)
}

// BeNil reports a test failure if the value is not nil.
//
// This assertion works with pointers, interfaces, channels, functions, slices, and maps.
// It uses Go's reflection to check if the value is nil.
//
// Example:
//
//	var ptr *int
//	should.BeNil(t, ptr)
//
//	var slice []int
//	should.BeNil(t, slice, should.WithMessage("Slice should be nil"))
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func BeNil[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeNil(t, actual, opts...)
}

// BeNotNil reports a test failure if the value is nil.
//
// This assertion works with pointers, interfaces, channels, functions, slices, and maps.
// It uses Go's reflection to check if the value is not nil.
//
// Example:
//
//	user := &User{Name: "John"}
//	should.BeNotNil(t, user, should.WithMessage("User must not be nil"))
//
//	should.BeNotNil(t, make([]int, 0))
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func BeNotNil[T any](t testing.TB, actual T, opts ...Option) {
	t.Helper()
	assert.BeNotNil(t, actual, opts...)
}

// BeGreaterThan reports a test failure if the value is not greater than the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides detailed
// error messages showing the actual value, threshold, difference, and helpful hints.
// It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeGreaterThan(t, 10, 5)
//
//	should.BeGreaterThan(t, user.Age, 18, should.WithMessage("User must be adult"))
//
//	should.BeGreaterThan(t, 3.14, 2.71)
//
// Only works with numeric types. Both values must be numeric.
func BeGreaterThan[T any](t testing.TB, actual, expected T, opts ...Option) {
	t.Helper()
	assert.BeGreaterThan(t, actual, expected, opts...)
}

// BeLessThan reports a test failure if the value is not less than the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides detailed
// error messages showing the actual value, threshold, difference, and helpful hints.
// It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeLessThan(t, 5, 10)
//
//	should.BeLessThan(t, user.Age, 65, should.WithMessage("User must be under retirement age"))
//
//	should.BeLessThan(t, 2.71, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func BeLessThan[T any](t testing.TB, actual, expected T, opts ...Option) {
	t.Helper()
	assert.BeLessThan(t, actual, expected, opts...)
}

// BeGreaterOrEqualThan reports a test failure if the value is not greater than or equal to the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides
// detailed error messages when the assertion fails. It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeGreaterOrEqualThan(t, 10, 10)
//
//	should.BeGreaterOrEqualThan(t, user.Score, 0, should.WithMessage("Score cannot be negative"))
//
//	should.BeGreaterOrEqualThan(t, 3.14, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func BeGreaterOrEqualThan[T any](t testing.TB, actual, expected T, opts ...Option) {
	t.Helper()
	assert.BeGreaterOrEqualThan(t, actual, expected, opts...)
}

// BeEqual reports a test failure if the two values are not deeply equal.
//
// This assertion uses Go's reflect.DeepEqual for comparison and provides detailed
// error messages showing exactly what differs between the values. For complex objects,
// it shows field-by-field differences to help identify the specific mismatches.
//
// Example:
//
//	should.BeEqual(t, "hello", "hello")
//
//	should.BeEqual(t, 42, 42)
//
//	should.BeEqual(t, user, expectedUser, should.WithMessage("User objects should match"))
//
// Works with any comparable types. Uses deep comparison for complex objects.
func BeEqual[T any](t testing.TB, actual, expected T, opts ...Option) {
	t.Helper()
	assert.BeEqual(t, actual, expected, opts...)
}

// Contain reports a test failure if the slice or array does not contain the expected value.
//
// This assertion provides intelligent error messages based on the type of collection:
// - For []string: Shows similar elements and typo detection
// - For numeric slices ([]int, []float64, etc.): Shows insertion context and sorted position
// - For other types: Shows formatted collection with clear error messages
// Supports all slice and array types.
//
// Example:
//
//	should.Contain(t, users, "user3")
//
//	should.Contain(t, []int{1, 2, 3}, 2)
//
//	should.Contain(t, []float64{1.1, 2.2}, 1.5, should.WithMessage("Expected value missing"))
//
//	should.Contain(t, []string{"apple", "banana"}, "apple")
//
// If the input is not a slice or array, the test fails immediately.
func Contain[T any](t testing.TB, actual T, expected any, opts ...Option) {
	t.Helper()
	assert.Contain(t, actual, expected, opts...)
}

// NotContain reports a test failure if the slice or array contains the expected value.
//
// This assertion works with slices and arrays of any type and provides detailed
// error messages showing where the unexpected element was found.
//
// Example:
//
//	should.NotContain(t, users, "bannedUser")
//
//	should.NotContain(t, []int{1, 2, 3}, 4)
//
//	should.NotContain(t, []string{"apple", "banana"}, "orange", should.WithMessage("Should not have orange"))
//
// If the input is not a slice or array, the test fails immediately.
func NotContain[T any](t testing.TB, actual T, expected any, opts ...Option) {
	t.Helper()
	assert.NotContain(t, actual, expected, opts...)
}

// ContainFunc reports a test failure if no element in the slice or array matches the predicate function.
//
// This assertion allows for custom matching logic by providing a predicate function
// that will be called for each element in the collection. The test passes if any element
// makes the predicate return true.
//
// Example:
//
//	should.ContainFunc(t, users, func(item any) bool {
//		user := item.(User)
//		return user.Age > 18
//	})
//
//	should.ContainFunc(t, numbers, func(item any) bool {
//		return item.(int) % 2 == 0
//	}, should.WithMessage("No even numbers found"))
//
// If the input is not a slice or array, the test fails immediately.
func ContainFunc[T any](t testing.TB, actual T, predicate func(item any) bool, opts ...Option) {
	t.Helper()
	assert.ContainFunc(t, actual, predicate, opts...)
}

// Panic asserts that the given function panics when executed.
// If the function does not panic, the test will fail with a descriptive error message.
//
// Example:
//
//	should.Panic(t, func() {
//		panic("expected panic")
//	})
func Panic(t testing.TB, fn func(), opts ...Option) {
	t.Helper()
	assert.Panic(t, fn, opts...)
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
func NotPanic(t testing.TB, fn func(), opts ...Option) {
	t.Helper()
	assert.NotPanic(t, fn, opts...)
}

// HaveLength reports a test failure if the collection does not have the expected length.
//
// This assertion works with strings, slices, arrays, and maps.
// It provides a detailed error message showing the expected and actual lengths,
// along with the difference.
//
// Example:
//
//	should.HaveLength(t, []int{1, 2, 3}, 3)
//	should.HaveLength(t, "hello", 5)
func HaveLength(t testing.TB, actual any, expected int, opts ...Option) {
	t.Helper()
	assert.HaveLength(t, actual, expected, opts...)
}

// BeOfType reports a test failure if the value is not of the expected type.
//
// This assertion checks if the type of the actual value matches the type
// of the expected value (using an instance of the expected type).
//
// Example:
//
//	type MyType struct{}
//	var v MyType
//	should.BeOfType(t, MyType{}, v)
func BeOfType(t testing.TB, actual, expected any, opts ...Option) {
	t.Helper()
	assert.BeOfType(t, actual, expected, opts...)
}

// BeOneOf reports a test failure if the value is not one of the provided options.
//
// This assertion checks if the actual value is present in the slice of allowed options.
// It uses deep comparison to check for equality.
//
// Example:
//
//	status := "pending"
//	allowedStatus := []string{"active", "inactive"}
//	should.BeOneOf(t, status, allowedStatus)
func BeOneOf[T any](t testing.TB, actual T, options []T, opts ...Option) {
	t.Helper()
	assert.BeOneOf(t, actual, options, opts...)
}
