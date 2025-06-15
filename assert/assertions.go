package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Ensure creates a new assertion for the given value.
//
// This is the entry point for all assertions in the Should library.
// It returns an Assertion object that provides fluent assertion methods.
//
// Example:
//
//	should.Ensure(42).BeEqual(t, 42)
//
//	should.Ensure("hello").BeNotEmpty(t)
//
//	should.Ensure([]int{1, 2, 3}).Contain(t, 2)
func Ensure[T any](value T) *Assertion[T] {
	return &Assertion[T]{value: value}
}

func fail(t testing.TB, message string, args ...any) {
	t.Helper()
	t.Errorf(message, args...)
}

// BeTrue reports a test failure if the value is not true.
//
// This assertion only works with boolean values and will fail immediately
// if the value is not a boolean type.
//
// Example:
//
//	should.Ensure(true).BeTrue(t)
//
//	should.Ensure(user.IsActive).BeTrue(t)
//
// If the input is not a boolean, the test fails immediately.
func (a *Assertion[T]) BeTrue(t testing.TB) {
	val, ok := any(a.value).(bool)
	if !ok {
		fail(t, "expected a boolean value, but got %T", a.value)
		return
	}

	if !val {
		fail(t, "Expected true, got false")
	}
}

// BeFalse reports a test failure if the value is not false.
//
// This assertion only works with boolean values and will fail immediately
// if the value is not a boolean type.
//
// Example:
//
//	should.Ensure(false).BeFalse(t)
//
//	should.Ensure(user.IsDeleted).BeFalse(t)
//
// If the input is not a boolean, the test fails immediately.
func (a *Assertion[T]) BeFalse(t testing.TB) {
	val, ok := any(a.value).(bool)
	if !ok {
		fail(t, "expected a boolean value, but got %T", a.value)
		return
	}

	if val {
		fail(t, "Expected false, got true")
	}
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
//	should.Ensure("").BeEmpty(t)
//
//	should.Ensure([]int{}).BeEmpty(t)
//
//	should.Ensure(map[string]int{}).BeEmpty(t)
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func (a *Assertion[T]) BeEmpty(t testing.TB) {
	t.Helper()
	actualValue := reflect.ValueOf(a.value)

	// Handle nil values
	if !actualValue.IsValid() {
		return // nil is considered empty
	}

	// Check if the type supports Len()
	switch actualValue.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if actualValue.Len() > 0 {
			fail(t, "%s", formatEmptyError(a.value, true))
		}
	case reflect.Ptr:
		if actualValue.IsNil() {
			return // nil pointer is considered empty
		}
		fail(t, "%s", formatEmptyError(a.value, true))
	default:
		fail(t, "BeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got %T", a.value)
	}
}

// BeNotEmpty reports a test failure if the value is empty.
//
// This assertion works with strings, slices, arrays, maps, channels, and pointers.
// For strings, non-empty means length > 0. For slices/arrays/maps/channels, non-empty means length > 0.
// For pointers, non-empty means not nil. Provides detailed error messages for empty values.
//
// Example:
//
//	should.Ensure("hello").BeNotEmpty(t)
//
//	should.Ensure([]int{1, 2, 3}).BeNotEmpty(t)
//
//	should.Ensure(&user).BeNotEmpty(t)
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func (a *Assertion[T]) BeNotEmpty(t testing.TB) {
	t.Helper()
	actualValue := reflect.ValueOf(a.value)

	// Handle nil values
	if !actualValue.IsValid() {
		fail(t, "%s", formatEmptyError(a.value, false))
		return
	}

	// Check if the type supports Len()
	switch actualValue.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if actualValue.Len() == 0 {
			fail(t, "%s", formatEmptyError(a.value, false))
		}
	case reflect.Ptr:
		if actualValue.IsNil() {
			fail(t, "%s", formatEmptyError(a.value, false))
		}
	default:
		fail(t, "BeNotEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got %T", a.value)
	}
}

// BeNil reports a test failure if the value is not nil.
//
// This assertion works with pointers, interfaces, channels, functions, slices, and maps.
// It uses Go's reflection to check if the value is nil.
//
// Example:
//
//	var ptr *int
//	should.Ensure(ptr).BeNil(t)
//
//	var slice []int
//	should.Ensure(slice).BeNil(t)
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func (a *Assertion[T]) BeNil(t testing.TB) {
	t.Helper()
	if !reflect.ValueOf(a.value).IsNil() {
		fail(t, "Expected nil, but was not")
	}
}

// BeNotNil reports a test failure if the value is nil.
//
// This assertion works with pointers, interfaces, channels, functions, slices, and maps.
// It uses Go's reflection to check if the value is not nil.
//
// Example:
//
//	user := &User{Name: "John"}
//	should.Ensure(user).BeNotNil(t)
//
//	should.Ensure(make([]int, 0)).BeNotNil(t)
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func (a *Assertion[T]) BeNotNil(t testing.TB) {
	t.Helper()
	if reflect.ValueOf(a.value).IsNil() {
		fail(t, "Expected not nil, but was nil")
	}
}

// BeGreaterThan reports a test failure if the value is not greater than the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides detailed
// error messages showing the actual value, threshold, difference, and helpful hints.
// It supports optional custom error messages.
//
// Example:
//
//	should.Ensure(10).BeGreaterThan(t, 5)
//
//	should.Ensure(user.Age).BeGreaterThan(t, 18, "User must be adult")
//
//	should.Ensure(3.14).BeGreaterThan(t, 2.71)
//
// Only works with numeric types. Both values must be numeric.
func (a *Assertion[T]) BeGreaterThan(t testing.TB, expected T, msgAndArgs ...interface{}) {
	t.Helper()
	actualV := reflect.ValueOf(a.value)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", a.value)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat <= expectedAsFloat {
		customMsg := messageFromMsgAndArgs(msgAndArgs...)
		errorMsg := formatNumericComparisonError(a.value, expected, "greater")
		if customMsg != "" {
			fail(t, "%s\n%s", customMsg, errorMsg)
		} else {
			fail(t, "%s", errorMsg)
		}
	}
}

// BeLessThan reports a test failure if the value is not less than the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides detailed
// error messages showing the actual value, threshold, difference, and helpful hints.
// It supports optional custom error messages.
//
// Example:
//
//	should.Ensure(5).BeLessThan(t, 10)
//
//	should.Ensure(user.Age).BeLessThan(t, 65, "User must be under retirement age")
//
//	should.Ensure(2.71).BeLessThan(t, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func (a *Assertion[T]) BeLessThan(t testing.TB, expected T, msgAndArgs ...interface{}) {
	t.Helper()
	actualV := reflect.ValueOf(a.value)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", a.value)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat >= expectedAsFloat {
		customMsg := messageFromMsgAndArgs(msgAndArgs...)
		errorMsg := formatNumericComparisonError(a.value, expected, "less")
		if customMsg != "" {
			fail(t, "%s\n%s", customMsg, errorMsg)
		} else {
			fail(t, "%s", errorMsg)
		}
	}
}

// BeGreaterOrEqualThan reports a test failure if the value is not greater than or equal to the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides
// detailed error messages when the assertion fails. It supports optional custom error messages.
//
// Example:
//
//	should.Ensure(10).BeGreaterOrEqualThan(t, 10)
//
//	should.Ensure(user.Score).BeGreaterOrEqualThan(t, 0, "Score cannot be negative")
//
//	should.Ensure(3.14).BeGreaterOrEqualThan(t, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func (a *Assertion[T]) BeGreaterOrEqualThan(t testing.TB, expected T, msgAndArgs ...interface{}) {
	t.Helper()
	actualV := reflect.ValueOf(a.value)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", a.value)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat < expectedAsFloat {
		customMsg := messageFromMsgAndArgs(msgAndArgs...)
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected %v to be greater or equal than %v", customMsg, a.value, expected)
	}
}

// BeEqual reports a test failure if the two values are not deeply equal.
//
// This assertion uses Go's reflect.DeepEqual for comparison and provides detailed
// error messages showing exactly what differs between the values. For complex objects,
// it shows field-by-field differences to help identify the specific mismatches.
//
// Example:
//
//	should.Ensure("hello").BeEqual(t, "hello")
//
//	should.Ensure(42).BeEqual(t, 42)
//
//	should.Ensure(user).BeEqual(t, expectedUser)
//
// Works with any comparable types. Uses deep comparison for complex objects.
func (a *Assertion[T]) BeEqual(t testing.TB, expected T, msgAndArgs ...interface{}) {
	t.Helper()

	if reflect.DeepEqual(a.value, expected) {
		return
	}

	diffs := findDifferences(expected, a.value)

	var differences []string
	differencesOutput := "Field differences:\n"
	for _, diff := range diffs {
		differencesOutput += fmt.Sprintf("  └─ %s: %s ≠ %s\n", diff.Path, formatDiffValue(diff.Expected), formatDiffValue(diff.Actual))
	}

	customMsg := messageFromMsgAndArgs(msgAndArgs...)
	if customMsg != "" {
		customMsg += "\n"
	}

	message := fmt.Sprintf(
		"%sNot equal:\nexpected: %v\nactual  : %v",
		customMsg,
		formatComparisonValue(expected),
		formatComparisonValue(a.value),
	)

	differences = append(differences, message, differencesOutput)

	diffMessage := strings.Join(differences, "\n")
	fail(t, "Differences found:\n%s", diffMessage)
}

// Contain reports a test failure if the slice or array does not contain the expected value.
//
// If the value is a []string and the expected is a string, it provides more detailed output
// to help identify similar or near-matching elements. For []int, it shows insertion context
// to help understand where the missing value would fit.
//
// Example:
//
//	should.Ensure(users).Contain(t, "user3")
//
//	should.Ensure([]int{1, 2, 3}).Contain(t, 2)
//
//	should.Ensure([]string{"apple", "banana"}).Contain(t, "apple")
//
// If the input is not a slice or array, the test fails immediately.
func (a *Assertion[T]) Contain(t testing.TB, expected any, msgAndArgs ...interface{}) {
	if !isSliceOrArray(a.value) {
		fail(t, "expected a slice or array, but got %T", a.value)
		return
	}

	// Try to convert to slice of strings for special treatment
	if collection, ok := any(a.value).([]string); ok {
		if target, ok := expected.(string); ok {
			result := containsString(target, collection)
			if result.Found {
				return
			}
			output := formatContainsError(target, result)
			fail(t, "%s", output)
			return
		}
	}

	if collection, ok := any(a.value).([]int); ok {
		if target, ok := expected.(int); ok {
			window, insertIndex := findInsertionContext(collection, target)
			output := formatInsertionContext(collection, target, window)
			if window == "" && insertIndex != -1 {
				return
			}
			fail(t, "%s", output)
			return
		}
	}

	// Fallback for other types of slice
	actualValue := reflect.ValueOf(a.value)
	for i := range actualValue.Len() {
		item := actualValue.Index(i).Interface()
		if reflect.DeepEqual(item, expected) {
			return
		}
	}

	// If not found, fail with a generic message
	fail(t, "\nExpected %v to contain %v", formatSlice(a.value), formatComparisonValue(expected))
}

// NotContain reports a test failure if the slice or array contains the expected value.
//
// This assertion works with slices and arrays of any type and provides detailed
// error messages showing where the unexpected element was found.
//
// Example:
//
//	should.Ensure(users).NotContain(t, "bannedUser")
//
//	should.Ensure([]int{1, 2, 3}).NotContain(t, 4)
//
//	should.Ensure([]string{"apple", "banana"}).NotContain(t, "orange")
//
// If the input is not a slice or array, the test fails immediately.
func (a *Assertion[T]) NotContain(t testing.TB, expected any, msgAndArgs ...interface{}) {
	if !isSliceOrArray(a.value) {
		fail(t, "expected a slice or array, but got %T", a.value)
		return
	}

	actualValue := reflect.ValueOf(a.value)

	foundOutput := []string{}
	for i := range actualValue.Len() {
		item := actualValue.Index(i).Interface()
		if reflect.DeepEqual(item, expected) {
			foundOutput = append(foundOutput, fmt.Sprintf("\nCollection: %s", formatSlice(a.value)))
			foundOutput = append(foundOutput, fmt.Sprintf("Found: %s at index %d", formatComparisonValue(item), i))
			output := strings.Join(foundOutput, "\n")
			fail(t, "\nExpected collection to NOT contain element: %s", output)
			return
		}
	}
}

// ContainFunc reports a test failure if no element in the slice or array matches the predicate function.
//
// This assertion allows for custom matching logic by providing a predicate function
// that will be called for each element in the collection. The test passes if any element
// makes the predicate return true.
//
// Example:
//
//	should.Ensure(users).ContainFunc(t, func(item any) bool {
//		user := item.(User)
//		return user.Age > 18
//	})
//
//	should.Ensure(numbers).ContainFunc(t, func(item any) bool {
//		return item.(int) % 2 == 0
//	})
//
// If the input is not a slice or array, the test fails immediately.
func (a *Assertion[T]) ContainFunc(t testing.TB, expected func(TItem any) bool, msgAndArgs ...interface{}) {
	if !isSliceOrArray(a.value) {
		fail(t, "expected a slice or array, but got %T", a.value)
		return
	}

	actualValue := reflect.ValueOf(a.value)

	for i := range actualValue.Len() {
		item := actualValue.Index(i).Interface()
		if expected(item) {
			return
		}
	}

	fail(t, "\nPredicate does not match any item in the slice")
}

// Panic reports a test failure if the given function does not panic.
//
// This assertion executes the provided function and expects it to panic.
// It captures and recovers from the panic to prevent the test from crashing.
// Supports optional custom error messages.
//
// Example:
//
//	should.Panic(t, func() {
//		panic("expected panic")
//	})
//
//	should.Panic(t, func() {
//		divide(1, 0)
//	}, "Division by zero should panic")
//
// The function parameter must not be nil.
func Panic(t testing.TB, fn func(), msgAndArgs ...interface{}) {
	t.Helper()
	panicked, _ := didPanic(fn)
	if !panicked {
		customMsg := messageFromMsgAndArgs(msgAndArgs...)
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected panic, but did not panic", customMsg)
	}
}

// NotPanic reports a test failure if the given function panics.
//
// This assertion executes the provided function and expects it to complete normally
// without panicking. If a panic occurs, it captures the panic value and includes it
// in the error message. Supports optional custom error messages.
//
// Example:
//
//	should.NotPanic(t, func() {
//		result := add(1, 2)
//		_ = result
//	})
//
//	should.NotPanic(t, func() {
//		user.Save()
//	}, "Save operation should not panic")
//
// The function parameter must not be nil.
func NotPanic(t testing.TB, fn func(), msgAndArgs ...interface{}) {
	t.Helper()
	panicked, r := didPanic(fn)
	if panicked {
		customMsg := messageFromMsgAndArgs(msgAndArgs...)
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected for the function to not panic, but it panicked with: %v", customMsg, r)
	}
}

// didPanic executes a function and reports whether it panicked, returning the recovered value.
func didPanic(fn func()) (panicked bool, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			recovered = r
		}
	}()
	fn()
	return
}

func toFloat64(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}
