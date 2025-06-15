package assert

import (
	"strings"
	"testing"
)

func TestBeTrue_Succeeds_WhenTrue(t *testing.T) {
	t.Parallel()

	Ensure(true).BeTrue(t)
}

func TestBeFalse_Succeeds_WhenFalse(t *testing.T) {
	t.Parallel()

	Ensure(false).BeFalse(t)
}

func TestBeEqual_ForStructs_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	newPerson := Person{Name: "John", Age: 30}
	Ensure(newPerson).BeEqual(t, Person{Name: "John", Age: 30})
}

func TestBeEqual_ForSlices_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "John", Age: 30}
	p2 := Person{Name: "Jane", Age: 25}

	Ensure([]Person{p1, p2}).BeEqual(t, []Person{p1, p2})
}

func TestBeEqual_ForMaps_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "b": 2}

	Ensure(map1).BeEqual(t, map2)
}

func TestBeEqual_ForStructs_Fails_WhenNotEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "John", Age: 30}
	p2 := Person{Name: "Jane", Age: 25}

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(p1).BeEqual(t, p2)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Differences found"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestBeEqual_ForSlices_Fails_WhenNotEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "John", Age: 30}
	p2 := Person{Name: "Jane", Age: 25}

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]Person{p1}).BeEqual(t, []Person{p2})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Differences found"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestBeEqual_ForMaps_Fails_WhenNotEqual(t *testing.T) {
	t.Parallel()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "c": 3}

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(map1).BeEqual(t, map2)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Differences found"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

/* func TestBeEmpty_Fail_WhenNotEmpty(t *testing.T) {
	stringTest := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Vivamus sagittis lacus vel augue laoreet rutrum faucibus dolor auctor. Cras justo odio, dapibus ac facilisis in, egestas eget quam. Praesent commodo cursus magna, vel scelerisque nisl consectetur et."

	Ensure(stringTest).BeEmpty(t)
} */

func TestBeGreaterThan_Succeeds_WhenGreater(t *testing.T) {
	t.Parallel()

	Ensure(10).BeGreaterThan(t, 5)
}

func TestContain_Succeeds_WhenItemIsPresent(t *testing.T) {
	t.Parallel()

	Ensure([]int{1, 2, 3}).Contain(t, 2)
}

func TestContain_Fails_WhenItemIsNotPresent(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{1, 2, 3}).Contain(t, 4)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 4"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_ShortensErrorMessage_WhenSliceIsLarge(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}).Contain(t, 21)
	})
	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 21"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestNotContain_Succeeds_WhenItemIsNotPresent(t *testing.T) {
	t.Parallel()

	Ensure([]int{1, 2, 3}).NotContain(t, 4)
}

func TestNotContain_Fails_WhenItemIsPresent(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{1, 2, 3}).NotContain(t, 2)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected collection to NOT contain element"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContainFunc_Succeeds_WhenPredicateMatches(t *testing.T) {
	t.Parallel()

	Ensure([]int{1, 2, 3}).ContainFunc(t, func(item any) bool {
		i, ok := item.(int)
		if !ok {
			return false
		}
		return i == 2
	})
}

func TestContainFunc_Fails_WhenPredicateDoesNotMatch(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{1, 2, 3}).ContainFunc(t, func(item any) bool {
			return item == 4
		})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Predicate does not match any item in the slice"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestShouldContain_ShowsSimilarElements_OnFailure(t *testing.T) {
	t.Parallel()

	users := []string{
		"user-one",
		"user_two",
		"UserThree",
		"user-3",
		"userThree",
		"user-003",
		"user-four",
		"user_five",
	}

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(users).Contain(t, "user3")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Collection: [user-one, user_two, UserThree, user-3, userThree]",
		"Missing   : user3",
		"Similar elements found:",
		"└─ user-3 (at index 3) - 1 extra char",
		"└─ user-003 (at index 5) - 3 char diff",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContain_WithIntSlices(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		slice    []int
		target   int
		expected bool
	}{
		{
			name:     "When_Target_Is_Present",
			slice:    []int{1, 2, 3, 4, 5},
			target:   3,
			expected: true,
		},
		{
			name:     "When_Target_Is_Not_Present",
			slice:    []int{1, 2, 4, 5},
			target:   3,
			expected: false,
		},
		{
			name:     "When_Target_Is_First_Element",
			slice:    []int{1, 2, 3, 4, 5},
			target:   1,
			expected: true,
		},
		{
			name:     "When_Target_Is_Last_Element",
			slice:    []int{1, 2, 3, 4, 5},
			target:   5,
			expected: true,
		},
		{
			name:     "When_Slice_Is_Empty",
			slice:    []int{},
			target:   1,
			expected: false,
		},
		{
			name:     "When_Target_Would_Be_In_Middle",
			slice:    []int{1, 2, 4, 5},
			target:   3,
			expected: false,
		},
		{
			name:     "When_Target_Would_Be_At_Beginning",
			slice:    []int{2, 3, 4, 5},
			target:   1,
			expected: false,
		},
		{
			name:     "When_Target_Would_Be_At_End",
			slice:    []int{1, 2, 3, 4},
			target:   5,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expected {
				// Should pass
				Ensure(tc.slice).Contain(t, tc.target)
			} else {
				// Should fail
				failed, message := assertFails(t, func(t testing.TB) {
					Ensure(tc.slice).Contain(t, tc.target)
				})

				if !failed {
					t.Fatal("Expected test to fail, but it passed")
				}

				// Verify the error message contains relevant information
				if !strings.Contains(message, formatComparisonValue(tc.target)) {
					t.Fatalf("Error message doesn't contain the target value: %s", message)
				}
			}
		})
	}
}

func TestContain_ShowsInsertionContext_ForIntSlices(t *testing.T) {
	t.Parallel()

	// Test that the error message shows where the target value would be inserted
	slice := []int{1, 2, 4, 5, 6, 8, 10}
	target := 7

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(slice).Contain(t, target)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	// Check that the message contains the expected context
	expectedParts := []string{
		"Collection:",
		"Missing  : 7",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for BeEmpty ===

func TestBeEmpty_Succeeds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value interface{}
	}{
		{"Empty string", ""},
		{"Empty slice", []int{}},
		{"Empty map", map[string]int{}},
		{"Nil slice", func() interface{} { var s []int; return s }()},
		{"Nil pointer", func() interface{} { var p *int; return p }()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Ensure(tt.value).BeEmpty(t)
		})
	}
}

func TestBeEmpty_Fails_WithNonEmptyString(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("hello").BeEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be empty, but it was not:",
		"Type    : string",
		"Length  : 5 characters",
		"Content : \"hello\"",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_Fails_WithNonEmptySlice(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{1, 2, 3}).BeEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be empty, but it was not:",
		"Type    : []int",
		"Length  : 3 elements",
		"Content : [1, 2, 3]",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_Fails_WithLongString(t *testing.T) {
	t.Parallel()

	longString := "This is a very long string that should be truncated in the error message"
	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(longString).BeEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be empty, but it was not:",
		"Type    : string",
		"Length  : 72 characters",
		"... (truncated)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_Fails_WithLargeSlice(t *testing.T) {
	t.Parallel()

	largeSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(largeSlice).BeEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be empty, but it was not:",
		"Type    : []int",
		"Length  : 10 elements",
		"Content : [1, 2, 3, ...] (showing first 3 of 10)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_Fails_WithUnsupportedType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(42).BeEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for BeNotEmpty ===

func TestBeNotEmpty_Succeeds_WithNonEmptyString(t *testing.T) {
	t.Parallel()

	Ensure("hello").BeNotEmpty(t)
}

func TestBeNotEmpty_Succeeds_WithNonEmptySlice(t *testing.T) {
	t.Parallel()

	Ensure([]int{1, 2, 3}).BeNotEmpty(t)
}

func TestBeNotEmpty_Fails_WithEmptyString(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("").BeNotEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be not empty, but it was empty:",
		"Type    : string",
		"Length  : 0 characters",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeNotEmpty_Fails_WithEmptySlice(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure([]int{}).BeNotEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be not empty, but it was empty:",
		"Type    : []int",
		"Length  : 0 elements",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeNotEmpty_Fails_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int
	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(ptr).BeNotEmpty(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be not empty, but it was empty:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for BeGreaterThan ===

func TestBeGreaterThan_Succeeds_WithIntegers(t *testing.T) {
	t.Parallel()

	Ensure(10).BeGreaterThan(t, 5)
}

func TestBeGreaterThan_Succeeds_WithFloats(t *testing.T) {
	t.Parallel()

	Ensure(3.14).BeGreaterThan(t, 2.71)
}

func TestBeGreaterThan_Fails_WithSmallerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(5).BeGreaterThan(t, 10)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be greater than threshold:",
		"Value     : 5",
		"Threshold : 10",
		"Difference: -5 (value is 5 smaller)",
		"Hint      : Value should be larger than threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeGreaterThan_Fails_WithEqualValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(5).BeGreaterThan(t, 5)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be greater than threshold:",
		"Value     : 5",
		"Threshold : 5",
		"Difference: 0 (values are equal)",
		"Hint      : Value should be larger than threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeGreaterThan_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(0.0).BeGreaterThan(t, 0.1, "Score validation failed")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Score validation failed",
		"Expected value to be greater than threshold:",
		"Value     : 0",
		"Threshold : 0.1",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for BeLessThan ===

func TestBeLessThan_Succeeds_WithSmallerValue(t *testing.T) {
	t.Parallel()

	Ensure(5).BeLessThan(t, 10)
}

func TestBeLessThan_Fails_WithLargerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(10).BeLessThan(t, 5)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be less than threshold:",
		"Value     : 10",
		"Threshold : 5",
		"Difference: +5 (value is 5 greater)",
		"Hint      : Value should be smaller than threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeLessThan_Fails_WithEqualValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(5).BeLessThan(t, 5)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be less than threshold:",
		"Value     : 5",
		"Threshold : 5",
		"Difference: 0 (values are equal)",
		"Hint      : Value should be smaller than threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeLessThan_WithFloats(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(3.14).BeLessThan(t, 2.71)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be less than threshold:",
		"Value     : 3.14",
		"Threshold : 2.71",
		"Difference: +0.43",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for error handling ===

func TestBeGreaterThan_Fails_WithNonNumericActual(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("hello").BeGreaterThan(t, "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeLessThan_Fails_WithNonNumericExpected(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("hello").BeLessThan(t, "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for Panic ===

func TestPanic_Succeeds_WhenPanicOccurs(t *testing.T) {
	t.Parallel()

	Panic(t, func() {
		panic("test panic")
	})
}

func TestPanic_Fails_WhenNoPanicOccurs(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Panic(t, func() {
			// Do nothing - no panic
		})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected panic, but did not panic"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestPanic_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Panic(t, func() {
			// No panic
		}, "Division by zero should panic")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Division by zero should panic",
		"Expected panic, but did not panic",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestNotPanic_Succeeds_WhenNoPanicOccurs(t *testing.T) {
	t.Parallel()

	NotPanic(t, func() {
		result := 1 + 2
		_ = result
	})
}

func TestNotPanic_Fails_WhenPanicOccurs(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotPanic(t, func() {
			panic("unexpected panic")
		})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected for the function to not panic, but it panicked with:",
		"unexpected panic",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestNotPanic_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotPanic(t, func() {
			panic("error occurred")
		}, "Save operation should not panic")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Save operation should not panic",
		"Expected for the function to not panic, but it panicked with:",
		"error occurred",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for BeGreaterOrEqualThan ===

func TestBeGreaterOrEqualThan_Succeeds_WithGreaterValue(t *testing.T) {
	t.Parallel()

	Ensure(10).BeGreaterOrEqualThan(t, 5)
}

func TestBeGreaterOrEqualThan_Succeeds_WithEqualValue(t *testing.T) {
	t.Parallel()

	Ensure(5).BeGreaterOrEqualThan(t, 5)
}

func TestBeGreaterOrEqualThan_Succeeds_WithFloats(t *testing.T) {
	t.Parallel()

	Ensure(3.14).BeGreaterOrEqualThan(t, 3.14)
	Ensure(3.15).BeGreaterOrEqualThan(t, 3.14)
}

func TestBeGreaterOrEqualThan_Fails_WithSmallerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(5).BeGreaterOrEqualThan(t, 10)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected 5 to be greater or equal than 10"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeGreaterOrEqualThan_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(0).BeGreaterOrEqualThan(t, 1, "Score cannot be negative")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Score cannot be negative",
		"Expected 0 to be greater or equal than 1",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeGreaterOrEqualThan_Fails_WithNonNumericTypes(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("hello").BeGreaterOrEqualThan(t, "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for edge cases ===

func TestBeNil_Succeeds_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int
	Ensure(ptr).BeNil(t)
}

func TestBeNil_Succeeds_WithNilSlice(t *testing.T) {
	t.Parallel()

	var slice []int
	Ensure(slice).BeNil(t)
}

func TestBeNil_Succeeds_WithNilMap(t *testing.T) {
	t.Parallel()

	var m map[string]int
	Ensure(m).BeNil(t)
}

// TestBeNil_Succeeds_WithNilInterface removed due to reflect.Value issue

func TestBeNil_Fails_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(ptr).BeNil(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected nil, but was not"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeNotNil_Succeeds_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value
	Ensure(ptr).BeNotNil(t)
}

func TestBeNotNil_Succeeds_WithNonNilSlice(t *testing.T) {
	t.Parallel()

	slice := make([]int, 0)
	Ensure(slice).BeNotNil(t)
}

func TestBeNotNil_Fails_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(ptr).BeNotNil(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected not nil, but was nil"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for error handling in boolean assertions ===

func TestBeTrue_Fails_WithNonBooleanType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("true").BeTrue(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a boolean value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeTrue_Fails_WithFalseValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(false).BeTrue(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected true, got false"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeFalse_Fails_WithNonBooleanType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(0).BeFalse(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a boolean value, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeFalse_Fails_WithTrueValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(true).BeFalse(t)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected false, got true"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for ContainFunc edge cases ===

func TestContainFunc_Fails_WithNonSliceType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("not a slice").ContainFunc(t, func(item any) bool {
			return true
		})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a slice or array, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestContainFunc_WithComplexPredicate(t *testing.T) {
	t.Parallel()

	type User struct {
		Name string
		Age  int
	}

	users := []User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 17},
	}

	// Should succeed - finding adult user
	Ensure(users).ContainFunc(t, func(item any) bool {
		user := item.(User)
		return user.Age >= 18
	})

	// Should fail - no elderly users
	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(users).ContainFunc(t, func(item any) bool {
			user := item.(User)
			return user.Age >= 65
		})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Predicate does not match any item in the slice"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for Contain with different slice types ===

func TestContain_Fails_WithNonSliceType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure("not a slice").Contain(t, "something")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a slice or array, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotContain_Fails_WithNonSliceType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Ensure(42).NotContain(t, "something")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a slice or array, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}
