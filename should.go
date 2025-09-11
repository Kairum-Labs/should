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
	"time"

	"github.com/Kairum-Labs/should/assert"
)

// Option is a functional option for configuring assertions.
type Option = assert.Option

// WithMessage creates an option for setting a custom error message.
//
// The message is treated as a plain string literal. Use this when you
// want to display a fixed message without formatting or placeholders.
//
// Example usage:
//
//	should.BeGreaterThan(t, userAge, 18, should.WithMessage("User must be adult"))
func WithMessage(message string) Option {
	return assert.WithMessage(message)
}

// WithMessagef creates an option for setting a custom error message with formatting.
//
// The message supports placeholders, similar to fmt.Sprintf, and takes
// optional arguments to replace them. Use this when you need dynamic
// content in the message.
//
// Example usage:
//
//	should.BeLessOrEqualTo(t, score, 100, should.WithMessagef("Score cannot exceed %d", 100))
func WithMessagef(message string, args ...any) Option {
	return assert.WithMessagef(message, args...)
}

// WithIgnoreCase returns an option that makes string comparisons case-insensitive.
//
// This option can be passed to assertions that perform string comparisons,
// such as StartWith and EndWith, to ensure that case differences are ignored.
//
// Example:
//
//	should.StartWith(t, "hello", "HELLO", should.WithIgnoreCase())
//	should.EndWith(t, "Hello, world", "WORLD", should.WithIgnoreCase())
func WithIgnoreCase() Option {
	return assert.WithIgnoreCase()
}

// WithStackTrace creates an option for including stack traces on NotPanic assertions.
//
// Example:
//
//	should.NotPanic(t, func() {
//		panic("expected panic")
//	}, should.WithStackTrace())
func WithStackTrace() Option {
	return assert.WithStackTrace()
}

// WithIgnoreTimezone returns an option that makes time comparisons ignore timezone/location differences.
//
// Currently, this option is only supported by BeSameTime.
//
// Example:
//
//	should.BeSameTime(t, actual, expected, should.WithIgnoreTimezone())
func WithIgnoreTimezone() Option {
	return assert.WithIgnoreTimezone()
}

// WithTruncate truncates the actual and expected times to the specified unit before comparing them for equality.
//
// This is useful for asserting that two times are the same up to a certain level of precision,
// ignoring differences in smaller units.
//
// Example:
//
//	time1 := time.Date(2024, 8, 10, 15, 30, 0, 1_000_000, time.UTC)
//	time2 := time.Date(2024, 8, 10, 15, 30, 0, 999_999_999, time.UTC)
//
//	// This assertion will pass because both times truncate to 15:30:00.
//	should.BeSameTime(t, time1, time2, should.WithTruncate(time.Second))
func WithTruncate(unit time.Duration) Option {
	return assert.WithTruncate(unit)
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
func BeTrue(t testing.TB, actual bool, opts ...Option) {
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
func BeFalse(t testing.TB, actual bool, opts ...Option) {
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
func BeEmpty(t testing.TB, actual any, opts ...Option) {
	t.Helper()
	assert.BeEmpty(t, actual, opts...)
}

// NotBeEmpty reports a test failure if the value is empty.
//
// This assertion works with strings, slices, arrays, maps, channels, and pointers.
// For strings, non-empty means length > 0. For slices/arrays/maps/channels, non-empty means length > 0.
// For pointers, non-empty means not nil. Provides detailed error messages for empty values.
//
// Example:
//
//	should.NotBeEmpty(t, "hello")
//
//	should.NotBeEmpty(t, []int{1, 2, 3}, should.WithMessage("List must have items"))
//
//	should.NotBeEmpty(t, &user)
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func NotBeEmpty(t testing.TB, actual any, opts ...Option) {
	t.Helper()
	assert.NotBeEmpty(t, actual, opts...)
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
func BeNil(t testing.TB, actual any, opts ...Option) {
	t.Helper()
	assert.BeNil(t, actual, opts...)
}

// NotBeNil reports a test failure if the value is nil.
//
// This assertion works with pointers, interfaces, channels, functions, slices, and maps.
// It uses Go's reflection to check if the value is not nil.
//
// Example:
//
//	user := &User{Name: "John"}
//	should.NotBeNil(t, user, should.WithMessage("User must not be nil"))
//
//	should.NotBeNil(t, make([]int, 0))
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func NotBeNil(t testing.TB, actual any, opts ...Option) {
	t.Helper()
	assert.NotBeNil(t, actual, opts...)
}

// BeError reports a test failure if the provided error is nil.
//
// This assertion is useful to ensure that a function call actually
// produced an error when one is expected. It provides clear failure
// messages showing when an error was expected but not returned.
// It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeError(t, err)
//	should.BeError(t, err, should.WithMessage("Expected a validation error"))
func BeError(t testing.TB, err error, opts ...Option) {
	t.Helper()
	assert.BeError(t, err)
}

// BeErrorAs reports a test failure if the provided error does not match
// the target type using errors.As.
//
// This assertion is useful when you need to verify that an error
// can be unwrapped into a specific type, such as a custom error struct.
// It supports optional custom error messages through Option.
//
// Example:
//
//	var pathErr *os.PathError
//	should.BeErrorAs(t, err, &pathErr)
//	should.BeErrorAs(t, err, &MyCustomError{}, should.WithMessage("Expected custom error type"))
func BeErrorAs(t *testing.T, err error, target interface{}, opts ...Option) {
	t.Helper()
	assert.BeErrorAs(t, err, target)
}

// BeErrorIs reports a test failure if the provided error is not equal to
// the target error using errors.Is.
//
// This assertion is useful to check if an error matches a specific sentinel
// value, such as io.EOF or custom exported error variables.
// It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeErrorIs(t, err, io.EOF)
//	should.BeErrorIs(t, err, ErrUnauthorized, should.WithMessage("Expected unauthorized error"))
func BeErrorIs(t *testing.T, err error, target error, opts ...Option) {
	t.Helper()
	assert.BeErrorIs(t, err, target)
}

// BeGreaterThan reports a test failure if the value is not greater than the expected threshold.
//
// This assertion works with all numeric types and provides detailed
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
// Only works with numeric types. Both values must be of the same type.
func BeGreaterThan[T assert.Ordered](t testing.TB, actual T, expected T, opts ...Option) {
	t.Helper()
	assert.BeGreaterThan(t, actual, expected, opts...)
}

// BeLessThan reports a test failure if the value is not less than the expected threshold.
//
// This assertion works with all numeric types and provides detailed
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
// Only works with numeric types. Both values must be of the same type.
func BeLessThan[T assert.Ordered](t testing.TB, actual T, expected T, opts ...Option) {
	t.Helper()
	assert.BeLessThan(t, actual, expected, opts...)
}

// BeGreaterOrEqualTo reports a test failure if the value is not greater than or equal to the expected threshold.
//
// This assertion works with all numeric types and provides
// detailed error messages when the assertion fails. It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeGreaterOrEqualTo(t, 10, 10)
//
//	should.BeGreaterOrEqualTo(t, user.Score, 0, should.WithMessage("Score cannot be negative"))
//
//	should.BeGreaterOrEqualTo(t, 3.14, 3.14)
//
// Only works with numeric types. Both values must be of the same type.
func BeGreaterOrEqualTo[T assert.Ordered](t testing.TB, actual T, expected T, opts ...Option) {
	t.Helper()
	assert.BeGreaterOrEqualTo(t, actual, expected, opts...)
}

// BeLessOrEqualTo reports a test failure if the value is not less than or equal to the expected threshold.
//
// This assertion works with all numeric types and provides
// detailed error messages when the assertion fails. It supports optional custom error messages through Option.
//
// Example:
//
//	should.BeLessOrEqualTo(t, 5, 10)
//
//	should.BeLessOrEqualTo(t, user.Age, 65, should.WithMessage("User must be under retirement age"))
//
//	should.BeLessOrEqualTo(t, 3.14, 3.14)
//
// Only works with numeric types. Both values must be of the same type.
func BeLessOrEqualTo[T assert.Ordered](t testing.TB, actual T, expected T, opts ...Option) {
	t.Helper()
	assert.BeLessOrEqualTo(t, actual, expected, opts...)
}

// BeInRange reports a test failure if the value is not within the specified range (inclusive).
//
// This assertion works with all numeric types and provides detailed
// error messages when the assertion fails, indicating whether the value is
// above or below the range and by how much.
//
// Example:
//
//	should.BeInRange(t, 25, 18, 65)
//
//	should.BeInRange(t, 99.5, 0.0, 100.0)
//
//	should.BeInRange(t, 200, 200, 299, should.WithMessage("HTTP status should be 2xx"))
//
// Only works with numeric types. All values must be of the same type.
func BeInRange[T assert.Ordered](t testing.TB, actual T, minValue T, maxValue T, opts ...Option) {
	t.Helper()
	assert.BeInRange(t, actual, minValue, maxValue, opts...)
}

// BeSorted reports a test failure if the slice is not sorted in ascending order.
//
// This assertion works with slices of any ordered type (integers, floats, strings).
// For arrays, convert to slice using slice syntax: myArray[:].
// Provides detailed error messages showing order violations with indices and values.
//
// Example:
//
//	should.BeSorted(t, []int{1, 2, 3, 4, 5})
//
//	should.BeSorted(t, []string{"a", "b", "c"})
//
//	should.BeSorted(t, myArray[:]) // for arrays
//
// Only works with slices of ordered types (cmp.Ordered constraint).
func BeSorted[T assert.Sortable](t testing.TB, actual []T, opts ...Option) {
	t.Helper()
	assert.BeSorted(t, actual, opts...)
}

// BeEqual reports a test failure if the two values are not deeply equal.
//
// Uses reflect.DeepEqual for comparison. For primitive types (string, int, float, bool, etc.),
// it shows a simple message. For complex objects (structs, slices, maps), it shows
// field-by-field differences.
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

// NotBeEqual reports a test failure if the two values are deeply equal.
//
// This assertion uses Go's reflect.DeepEqual for comparison and provides detailed
// error messages showing exactly what differs between the values. For complex objects,
// it shows field-by-field differences to help identify the specific mismatches.
//
// Example:
//
//	should.NotBeEqual(t, "hello", "world")
//
//	should.NotBeEqual(t, 42, 43)
//
//	should.NotBeEqual(t, user, expectedUser, should.WithMessage("User objects should not match"))
func NotBeEqual[T any](t testing.TB, actual, expected T, opts ...Option) {
	t.Helper()
	assert.NotBeEqual(t, actual, expected, opts...)
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
func Contain(t testing.TB, actual any, expected any, opts ...Option) {
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
func NotContain(t testing.TB, actual any, expected any, opts ...Option) {
	t.Helper()
	assert.NotContain(t, actual, expected, opts...)
}

// AnyMatch reports a test failure if no element in the slice matches the predicate function.
//
// This assertion allows custom matching logic by providing a predicate function
// that will be called for each element in the slice. The test passes if any element
// makes the predicate return true.
//
// Example:
//
//	type User struct { Age int }
//
//	users := []User{{Age: 16}, {Age: 21}}
//	should.AnyMatch(t, users, func(user User) bool {
//		return user.Age > 18
//	})
//
//	numbers := []int{1, 3, 5, 8}
//	should.AnyMatch(t, numbers, func(n int) bool {
//		return n%2 == 0
//	}, should.WithMessage("No even numbers found"))
func AnyMatch[T any](t testing.TB, actual []T, predicate func(T) bool, opts ...Option) {
	t.Helper()
	assert.AnyMatch(t, actual, predicate, opts...)
}

// StartWith reports a test failure if the string does not start with the expected substring.
//
// This assertion checks if the actual string starts with the expected substring.
// It provides a detailed error message showing the expected and actual strings,
// along with a note if the case mismatch is detected.
//
// Example:
//
//	should.StartWith(t, "Hello, world!", "hello")
//
//	should.StartWith(t, "Hello, world!", "hello", should.WithIgnoreCase())
//
//	should.StartWith(t, "Hello, world!", "world", should.WithMessage("Expected string to start with 'world'"))
//
// Note: The assertion is case-sensitive by default. Use should.WithIgnoreCase() to ignore case.
func StartWith(t testing.TB, actual string, expected string, opts ...Option) {
	t.Helper()
	assert.StartWith(t, actual, expected, opts...)
}

// EndWith reports a test failure if the string does not end with the expected substring.
//
// This assertion checks if the actual string ends with the expected substring.
// It provides a detailed error message showing the expected and actual strings,
// along with a note if the case mismatch is detected.
//
// Example:
//
//	should.EndWith(t, "Hello, world!", "world")
//
//	should.EndWith(t, "Hello, world", "WORLD", should.WithIgnoreCase())
//
//	should.EndWith(t, "Hello, world!", "world", should.WithMessage("Expected string to end with 'world'"))
//
// Note: The assertion is case-sensitive by default. Use should.WithIgnoreCase() to ignore case.
func EndWith(t testing.TB, actual string, expected string, opts ...Option) {
	t.Helper()
	assert.EndWith(t, actual, expected, opts...)
}

// ContainSubstring reports a test failure if the string does not contain the expected substring.
//
// This assertion checks if the actual string contains the expected substring.
// It provides a detailed error message showing the expected and actual strings,
// with intelligent formatting for very long strings, and includes a note if
// case mismatch is detected. For needles up to 20 characters, it also provides
// typo detection using Levenshtein distance to suggest similar substrings.
//
// Example:
//
//	should.ContainSubstring(t, "Hello, world!", "world")
//
//	should.ContainSubstring(t, "Hello, World", "WORLD", should.WithIgnoreCase())
//
//	should.ContainSubstring(t, longText, "keyword", should.WithMessage("Expected keyword to be present"))
//
// Note: The assertion is case-sensitive by default. Use should.WithIgnoreCase() to ignore case.
// Typo detection is automatically enabled for needles up to 20 characters for performance.
func ContainSubstring(t testing.TB, actual string, substring string, opts ...Option) {
	t.Helper()
	assert.ContainSubstring(t, actual, substring, opts...)
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

// BeSameTime reports a test failure if two `time.Time` values do not represent the same time.
//
// By default, the comparison is timezone-sensitive and nanosecond-precise. You can customize
// the behavior with functional options:
//
// - should.WithIgnoreTimezone(): compares the instants regardless of the timezone/location
//
// - should.WithTruncate(unit): truncates both times to the specified precision before comparison
//
// Example:
//
//	should.BeSameTime(t, time1, time2)
//
//	should.BeSameTime(
//	    t,
//	    actual,
//	    expected,
//	    should.WithIgnoreTimezone(),
//	    should.WithTruncate(time.Second),
//	)
func BeSameTime(t testing.TB, actual time.Time, expected time.Time, opts ...Option) {
	t.Helper()
	assert.BeSameTime(t, actual, expected, opts...)
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

// ContainKey reports a test failure if the map does not contain the expected key.
//
// This assertion works with maps of any key type and provides intelligent error messages:
// - For string keys: Shows similar keys and typo detection
// - For numeric keys: Shows similar keys with numeric differences
// - For other types: Shows formatted keys with clear error messages
// Supports all map types.
//
// Example:
//
//	userMap := map[string]int{"name": 1, "age": 2}
//	should.ContainKey(t, userMap, "email")
//
//	should.ContainKey(t, map[int]string{1: "one", 2: "two"}, 3, should.WithMessage("Key must exist"))
func ContainKey[K comparable, V any](t testing.TB, actual map[K]V, expectedKey K, opts ...Option) {
	t.Helper()
	assert.ContainKey(t, actual, expectedKey, opts...)
}

// ContainValue reports a test failure if the map does not contain the expected value.
//
// This assertion works with maps of any value type and provides intelligent error messages:
// - For string values: Shows similar values and typo detection
// - For numeric values: Shows similar values with numeric differences
// - For other types: Shows formatted values with clear error messages
// Supports all map types.
//
// Example:
//
//	userMap := map[string]int{"name": 1, "age": 2}
//	should.ContainValue(t, userMap, 3)
//
//	should.ContainValue(t, map[int]string{1: "one", 2: "two"}, "three", should.WithMessage("Value must exist"))
func ContainValue[K comparable, V any](t testing.TB, actual map[K]V, expectedValue V, opts ...Option) {
	t.Helper()
	assert.ContainValue(t, actual, expectedValue, opts...)
}

// NotContainDuplicates reports a test failure if the slice or array contains duplicate values.
//
// This assertion works with slices and arrays of any type and provides detailed
// error messages showing where the duplicate values were found.
//
// Example:
//
//	should.NotContainDuplicates(t, []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 4, 4})
//
//	should.NotContainDuplicates(t, []string{"John", "John"})
//
// If the input is not a slice or array, the test fails immediately.
func NotContainDuplicates(t testing.TB, actual any, opts ...Option) {
	t.Helper()
	assert.NotContainDuplicates(t, actual, opts...)
}

// NotContainKey reports a test failure if the map contains the expected key.
//
// This assertion works with maps of any key type and provides detailed error messages
// showing where the key was found, including the map type, size, and context around
// the found key. Supports all map types.
//
// Example:
//
//	userMap := map[string]int{"name": 1, "age": 2}
//	should.NotContainKey(t, userMap, "age") // This will fail
//
//	should.NotContainKey(t, map[int]string{1: "one", 2: "two"}, 3, should.WithMessage("Key should not exist"))
func NotContainKey[K comparable, V any](t testing.TB, actual map[K]V, expectedKey K, opts ...Option) {
	t.Helper()
	assert.NotContainKey(t, actual, expectedKey, opts...)
}

// NotContainValue reports a test failure if the map contains the expected value.
//
// This assertion works with maps of any value type and provides detailed error messages
// showing where the value was found, including the map type, size, and context around
// the found value. Supports all map types.
//
// Example:
//
//	userMap := map[string]int{"name": 1, "age": 2}
//	should.NotContainValue(t, userMap, 2) // This will fail
//
//	should.NotContainValue(t, map[int]string{1: "one", 2: "two"}, "three", should.WithMessage("Value should not exist"))
func NotContainValue[K comparable, V any](t testing.TB, actual map[K]V, expectedValue V, opts ...Option) {
	t.Helper()
	assert.NotContainValue(t, actual, expectedValue, opts...)
}
