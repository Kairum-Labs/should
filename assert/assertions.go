package should

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Ensure[T any](value T) *Assertion[T] {
	return &Assertion[T]{value: value}
}

func fail(t testing.TB, message string, args ...any) {
	t.Helper()
	t.Errorf(message, args...)
}

// BeTrue verifies if the value is true.
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

// BeFalse verifies if the value is false.
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

// BeEmpty verifies if the value is empty.
func (a *Assertion[T]) BeEmpty(t testing.TB) {
	t.Helper()
	actualValue := reflect.ValueOf(a.value)

	if actualValue.Kind() == reflect.String {
		fail(t, "\nExpected empty, but was not\n%s", formatMultilineString(actualValue.String()))
		return
	}

	if actualValue.Len() > 0 {
		fail(t, "Expected empty, but was not")
	}
}

// BeNotEmpty verifies if the value is not empty.
func (a *Assertion[T]) BeNotEmpty(t testing.TB) {
	t.Helper()
	actualValue := reflect.ValueOf(a.value)
	if actualValue.Len() == 0 {
		fail(t, "Expected not empty, but was empty")
	}
}

// BeNil verifies if the value is nil.
func (a *Assertion[T]) BeNil(t testing.TB) {
	t.Helper()
	if !reflect.ValueOf(a.value).IsNil() {
		fail(t, "Expected nil, but was not")
	}
}

// BeNotNil verifies if the value is not nil.
func (a *Assertion[T]) BeNotNil(t testing.TB) {
	t.Helper()
	if reflect.ValueOf(a.value).IsNil() {
		fail(t, "Expected not nil, but was nil")
	}
}

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
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected %v to be greater than %v", customMsg, a.value, expected)
	}
}

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
		if customMsg != "" {
			customMsg += "\n"
		}
		fail(t, "%sExpected %v to be less than %v", customMsg, a.value, expected)
	}
}

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

// ShouldBeEqual asserts that two objects are deeply equal.
// It provides detailed output on what the differences are.
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
// to help identify similar or near-matching elements.
//
// Example:
//
//	Ensure(users).Contain(t, "user3")
//
//	Ensure([]int{1, 2, 3}).Contain(t, 2)
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

// Panic asserts that the given function panics.
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

// NotPanic asserts that the given function does not panic.
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
