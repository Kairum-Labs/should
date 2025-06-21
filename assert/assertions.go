// Package assert provides the underlying implementation for the Should assertion library.
// It contains the core assertion logic, which is then exposed through the top-level
// `should` package. This package handles value comparisons, error formatting,
// and detailed difference reporting.
package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

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
//	should.BeTrue(t, true)
//
//	should.BeTrue(t, user.IsActive, should.AssertionConfig{Message: "User must be active"})
//
// If the input is not a boolean, the test fails immediately.
func BeTrue[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	val, ok := any(actual).(bool)
	if !ok {
		fail(t, "expected a boolean value, but got %T", actual)
		return
	}

	if !val {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			fail(t, "%s\nExpected true, got false", customMsg)
		} else {
			fail(t, "Expected true, got false")
		}
	}
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
//	should.BeFalse(t, user.IsDeleted, should.AssertionConfig{Message: "User should not be deleted"})
//
// If the input is not a boolean, the test fails immediately.
func BeFalse[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	val, ok := any(actual).(bool)
	if !ok {
		fail(t, "expected a boolean value, but got %T", actual)
		return
	}

	if val {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			fail(t, "%s\nExpected false, got true", customMsg)
		} else {
			fail(t, "Expected false, got true")
		}
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
//	should.BeEmpty(t, "")
//
//	should.BeEmpty(t, []int{}, should.AssertionConfig{Message: "List should be empty"})
//
//	should.BeEmpty(t, map[string]int{})
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func BeEmpty[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	actualValue := reflect.ValueOf(actual)

	// Handle nil values
	if !actualValue.IsValid() {
		return // nil is considered empty
	}

	// Check if the type supports Len()
	switch actualValue.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if actualValue.Len() > 0 {
			var customMsg string
			if len(config) > 0 {
				customMsg = config[0].Message
			}
			errorMsg := formatEmptyError(actual, true)
			if customMsg != "" {
				fail(t, "%s\n%s", customMsg, errorMsg)
			} else {
				fail(t, "%s", errorMsg)
			}
		}
	case reflect.Ptr:
		if actualValue.IsNil() {
			return // nil pointer is considered empty
		}
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		errorMsg := formatEmptyError(actual, true)
		if customMsg != "" {
			fail(t, "%s\n%s", customMsg, errorMsg)
		} else {
			fail(t, "%s", errorMsg)
		}
	default:
		fail(t, "BeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got %T", actual)
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
//	should.BeNotEmpty(t, "hello")
//
//	should.BeNotEmpty(t, []int{1, 2, 3}, should.AssertionConfig{Message: "List must have items"})
//
//	should.BeNotEmpty(t, &user)
//
// Only works with strings, slices, arrays, maps, channels, or pointers.
func BeNotEmpty[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	actualValue := reflect.ValueOf(actual)

	// Handle nil values
	if !actualValue.IsValid() {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		errorMsg := formatEmptyError(actual, false)
		if customMsg != "" {
			fail(t, "%s\n%s", customMsg, errorMsg)
		} else {
			fail(t, "%s", errorMsg)
		}
		return
	}

	// Check if the type supports Len()
	switch actualValue.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if actualValue.Len() == 0 {
			var customMsg string
			if len(config) > 0 {
				customMsg = config[0].Message
			}
			errorMsg := formatEmptyError(actual, false)
			if customMsg != "" {
				fail(t, "%s\n%s", customMsg, errorMsg)
			} else {
				fail(t, "%s", errorMsg)
			}
		}
	case reflect.Ptr:
		if actualValue.IsNil() {
			var customMsg string
			if len(config) > 0 {
				customMsg = config[0].Message
			}
			errorMsg := formatEmptyError(actual, false)
			if customMsg != "" {
				fail(t, "%s\n%s", customMsg, errorMsg)
			} else {
				fail(t, "%s", errorMsg)
			}
		}
	default:
		fail(t, "BeNotEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got %T", actual)
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
//	should.BeNil(t, ptr)
//
//	var slice []int
//	should.BeNil(t, slice, should.AssertionConfig{Message: "Slice should be nil"})
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func BeNil[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	v := reflect.ValueOf(actual)

	if !v.IsValid() {
		return // A nil interface is considered nil.
	}

	kind := v.Kind()
	nillable := kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Interface ||
		kind == reflect.Map ||
		kind == reflect.Ptr ||
		kind == reflect.Slice

	if !nillable {
		fail(t, "BeNil can only be used with nillable types, but got %T", actual)
		return
	}

	if !v.IsNil() {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			fail(t, "%s\nExpected nil, but was not", customMsg)
		} else {
			fail(t, "Expected nil, but was not")
		}
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
//	should.BeNotNil(t, user, should.AssertionConfig{Message: "User must not be nil"})
//
//	should.BeNotNil(t, make([]int, 0))
//
// Only works with nillable types (pointers, interfaces, channels, functions, slices, maps).
func BeNotNil[T any](t testing.TB, actual T, config ...AssertionConfig) {
	t.Helper()
	v := reflect.ValueOf(actual)

	isNil := !v.IsValid()
	if v.IsValid() {
		kind := v.Kind()
		nillable := kind == reflect.Chan ||
			kind == reflect.Func ||
			kind == reflect.Interface ||
			kind == reflect.Map ||
			kind == reflect.Ptr ||
			kind == reflect.Slice

		if !nillable {
			fail(t, "BeNotNil can only be used with nillable types, but got %T", actual)
			return
		}
		isNil = v.IsNil()
	}

	if isNil {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			fail(t, "%s\nExpected not nil, but was nil", customMsg)
		} else {
			fail(t, "Expected not nil, but was nil")
		}
	}
}

// BeGreaterThan reports a test failure if the value is not greater than the expected threshold.
//
// This assertion works with all numeric types (int, float, etc.) and provides detailed
// error messages showing the actual value, threshold, difference, and helpful hints.
// It supports optional custom error messages through AssertionConfig.
//
// Example:
//
//	should.BeGreaterThan(t, 10, 5)
//
//	should.BeGreaterThan(t, user.Age, 18, should.AssertionConfig{Message: "User must be adult"})
//
//	should.BeGreaterThan(t, 3.14, 2.71)
//
// Only works with numeric types. Both values must be numeric.
func BeGreaterThan[T any](t testing.TB, actual T, expected T, config ...AssertionConfig) {
	t.Helper()
	actualV := reflect.ValueOf(actual)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", actual)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat <= expectedAsFloat {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		errorMsg := formatNumericComparisonError(actual, expected, "greater")
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
// It supports optional custom error messages through AssertionConfig.
//
// Example:
//
//	should.BeLessThan(t, 5, 10)
//
//	should.BeLessThan(t, user.Age, 65, should.AssertionConfig{Message: "User must be under retirement age"})
//
//	should.BeLessThan(t, 2.71, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func BeLessThan[T any](t testing.TB, actual T, expected T, config ...AssertionConfig) {
	t.Helper()
	actualV := reflect.ValueOf(actual)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", actual)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat >= expectedAsFloat {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		errorMsg := formatNumericComparisonError(actual, expected, "less")
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
// detailed error messages when the assertion fails. It supports optional custom error messages through AssertionConfig.
//
// Example:
//
//	should.BeGreaterOrEqualThan(t, 10, 10)
//
//	should.BeGreaterOrEqualThan(t, user.Score, 0, should.AssertionConfig{Message: "Score cannot be negative"})
//
//	should.BeGreaterOrEqualThan(t, 3.14, 3.14)
//
// Only works with numeric types. Both values must be numeric.
func BeGreaterOrEqualThan[T any](t testing.TB, actual T, expected T, config ...AssertionConfig) {
	t.Helper()
	actualV := reflect.ValueOf(actual)
	expectedV := reflect.ValueOf(expected)

	actualAsFloat, actualOk := toFloat64(actualV)
	expectedAsFloat, expectedOk := toFloat64(expectedV)

	if !actualOk {
		fail(t, "expected a number for actual value, but got %T", actual)
		return
	}

	if !expectedOk {
		fail(t, "expected a number for expected value, but got %T", expected)
		return
	}

	if actualAsFloat < expectedAsFloat {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected %v to be greater or equal than %v", customMsg, actual, expected)
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
//	should.BeEqual(t, "hello", "hello")
//
//	should.BeEqual(t, 42, 42)
//
//	should.BeEqual(t, user, expectedUser, should.AssertionConfig{Message: "User objects should match"})
//
// Works with any comparable types. Uses deep comparison for complex objects.
func BeEqual[T any](t testing.TB, actual T, expected T, config ...AssertionConfig) {
	t.Helper()

	if reflect.DeepEqual(actual, expected) {
		return
	}

	diffs := findDifferences(expected, actual)

	var differences []string
	differencesOutput := "Field differences:\n"
	for _, diff := range diffs {
		differencesOutput += fmt.Sprintf("  └─ %s: %s ≠ %s\n", diff.Path, formatDiffValue(diff.Expected), formatDiffValue(diff.Actual))
	}

	var customMsg string
	if len(config) > 0 {
		customMsg = config[0].Message
	}
	if customMsg != "" {
		customMsg += "\n"
	}

	message := fmt.Sprintf(
		"%sNot equal:\nexpected: %v\nactual  : %v",
		customMsg,
		formatComparisonValue(expected),
		formatComparisonValue(actual),
	)

	differences = append(differences, message, differencesOutput)

	diffMessage := strings.Join(differences, "\n")
	fail(t, "Differences found:\n%s", diffMessage)
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
//	should.Contain(t, []float64{1.1, 2.2}, 1.5, should.AssertionConfig{Message: "Expected value missing"})
//
//	should.Contain(t, []string{"apple", "banana"}, "apple")
//
// If the input is not a slice or array, the test fails immediately.
func Contain[T any](t testing.TB, actual T, expected any, config ...AssertionConfig) {
	t.Helper()
	if !isSliceOrArray(actual) {
		fail(t, "expected a slice or array, but got %T", actual)
		return
	}

	// Handle string slices with intelligent similarity detection
	if collection, ok := any(actual).([]string); ok {
		if target, ok := expected.(string); ok {
			result := containsString(target, collection)
			if result.Found {
				return
			}
			var customMsg string
			if len(config) > 0 {
				customMsg = config[0].Message
			}
			output := formatContainsError(target, result)
			if customMsg != "" {
				fail(t, "%s\n%s", customMsg, output)
			} else {
				fail(t, "%s", output)
			}
			return
		}
	}

	// Handle numeric slices with insertion context
	actualValue := reflect.ValueOf(actual)
	elemType := actualValue.Type().Elem()

	// Check if it's a numeric type and provide insertion context
	if isNumericType(elemType) {
		found, output := handleNumericSliceContain(actual, expected)
		if found {
			return
		}
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
		if customMsg != "" {
			fail(t, "%s\n%s", customMsg, output)
		} else {
			fail(t, "%s", output)
		}
		return
	}

	// Generic fallback for other types
	for i := range actualValue.Len() {
		item := actualValue.Index(i).Interface()
		if reflect.DeepEqual(item, expected) {
			return
		}
	}

	// If not found, fail with a detailed message
	var customMsg string
	if len(config) > 0 {
		customMsg = config[0].Message
	}
	baseMsg := fmt.Sprintf("Expected collection to contain element:\n  Collection: %s\n  Missing   : %s",
		formatSlice(actual), formatComparisonValue(expected))

	if customMsg != "" {
		fail(t, "%s\n%s", customMsg, baseMsg)
	} else {
		fail(t, "%s", baseMsg)
	}
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
//	should.NotContain(t, []string{"apple", "banana"}, "orange", should.AssertionConfig{Message: "Should not have orange"})
//
// If the input is not a slice or array, the test fails immediately.
func NotContain[T any](t testing.TB, actual T, expected any, config ...AssertionConfig) {
	t.Helper()
	if !isSliceOrArray(actual) {
		fail(t, "expected a slice or array, but got %T", actual)
		return
	}

	actualValue := reflect.ValueOf(actual)

	foundOutput := []string{}
	for i := range actualValue.Len() {
		item := actualValue.Index(i).Interface()
		if reflect.DeepEqual(item, expected) {
			foundOutput = append(foundOutput, fmt.Sprintf("\nCollection: %s", formatSlice(actual)))
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
//	should.ContainFunc(t, users, func(item any) bool {
//		user := item.(User)
//		return user.Age > 18
//	})
//
//	should.ContainFunc(t, numbers, func(item any) bool {
//		return item.(int) % 2 == 0
//	}, should.AssertionConfig{Message: "No even numbers found"})
//
// If the input is not a slice or array, the test fails immediately.
func ContainFunc[T any](t testing.TB, actual T, expected func(TItem any) bool, config ...AssertionConfig) {
	t.Helper()
	if !isSliceOrArray(actual) {
		fail(t, "expected a slice or array, but got %T", actual)
		return
	}

	actualValue := reflect.ValueOf(actual)

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
// Supports optional custom error messages through AssertionConfig.
//
// Example:
//
//	should.Panic(t, func() {
//		panic("expected panic")
//	})
//
//	should.Panic(t, func() {
//		divide(1, 0)
//	}, should.AssertionConfig{Message: "Division by zero should panic"})
//
// The function parameter must not be nil.
func Panic(t testing.TB, fn func(), config ...AssertionConfig) {
	t.Helper()
	panicked, _ := didPanic(fn)
	if !panicked {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
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
// in the error message. Supports optional custom error messages through AssertionConfig.
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
//	}, should.AssertionConfig{Message: "Save operation should not panic"})
//
// The function parameter must not be nil.
func NotPanic(t testing.TB, fn func(), config ...AssertionConfig) {
	t.Helper()
	panicked, r := didPanic(fn)
	if panicked {
		var customMsg string
		if len(config) > 0 {
			customMsg = config[0].Message
		}
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

// isNumericType checks if a reflect.Type represents a numeric type
func isNumericType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// handleNumericSliceContain handles contain operations for numeric slices with insertion context
func handleNumericSliceContain(collection any, target any) (found bool, output string) {
	// Handle different numeric slice types
	switch coll := collection.(type) {
	case []int:
		if t, ok := target.(int); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []int8:
		if t, ok := target.(int8); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []int16:
		if t, ok := target.(int16); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []int32:
		if t, ok := target.(int32); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []int64:
		if t, ok := target.(int64); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []uint:
		if t, ok := target.(uint); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []uint8:
		if t, ok := target.(uint8); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []uint16:
		if t, ok := target.(uint16); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []uint32:
		if t, ok := target.(uint32); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []uint64:
		if t, ok := target.(uint64); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []float32:
		if t, ok := target.(float32); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	case []float64:
		if t, ok := target.(float64); ok {
			window, insertIndex := findInsertionContext(coll, t)
			if window == "" && insertIndex != -1 {
				return true, ""
			}
			return false, formatInsertionContext(coll, t, window)
		}
	}

	// Fallback for unsupported numeric types or type mismatches
	return false, fmt.Sprintf("Expected collection to contain element:\n  Collection: %s\n  Missing   : %s",
		formatSlice(collection), formatComparisonValue(target))
}
