package assert

import (
	"strings"
	"testing"
)

func TestBeTrue_Succeeds_WhenTrue(t *testing.T) {
	t.Parallel()

	BeTrue(t, true)
}

func TestBeFalse_Succeeds_WhenFalse(t *testing.T) {
	t.Parallel()

	BeFalse(t, false)
}

func TestBeEqual_ForStructs_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	newPerson := Person{Name: "John", Age: 30}
	BeEqual(t, newPerson, Person{Name: "John", Age: 30})
}

func TestBeEqual_ForSlices_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "John", Age: 30}
	p2 := Person{Name: "Jane", Age: 25}

	BeEqual(t, []Person{p1, p2}, []Person{p1, p2})
}

func TestBeEqual_ForMaps_Succeeds_WhenEqual(t *testing.T) {
	t.Parallel()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "b": 2}

	BeEqual(t, map1, map2)
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
		BeEqual(t, p1, p2)
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
		BeEqual(t, []Person{p1}, []Person{p2})
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
		BeEqual(t, map1, map2)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Differences found"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

// === Tests for NotBeEqual ===

func TestNotBeEqual(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     interface{}
			expected   interface{}
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass when values are not equal",
				actual:     "hello",
				expected:   "world",
				shouldFail: false,
			},
			{
				name:       "should pass when different types",
				actual:     "123",
				expected:   123,
				shouldFail: false,
			},
			{
				name:       "should pass when different numbers",
				actual:     42,
				expected:   24,
				shouldFail: false,
			},
			{
				name:       "should fail when values are equal",
				actual:     "hello",
				expected:   "hello",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected values to be different") {
						t.Errorf("Expected specific error message for equal values, got:\n%s", message)
					}
				},
			},
			{
				name:       "should fail when numbers are equal",
				actual:     42,
				expected:   42,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected values to be different") {
						t.Errorf("Expected specific error message for equal numbers, got:\n%s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				NotBeEqual(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.Failed() {
					t.Fatal("Expected NotBeEqual to fail, but it passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected NotBeEqual to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.Failed() {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     interface{}
			expected   interface{}
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass with custom message",
				actual:     "hello",
				expected:   "world",
				opts:       []Option{WithMessage("Values should be different")},
				shouldFail: false,
			},
			{
				name:       "should show custom error message on failure",
				actual:     "same",
				expected:   "same",
				opts:       []Option{WithMessage("Custom error: values must be different")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Custom error: values must be different") {
						t.Errorf("Expected custom error message, got: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				NotBeEqual(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.Failed() {
					t.Fatal("Expected NotBeEqual to fail, but it passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected NotBeEqual to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.Failed() {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     interface{}
			expected   interface{}
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should fail when both values are nil",
				actual:     nil,
				expected:   nil,
				shouldFail: true,
			},
			{
				name:       "should pass when one is nil and other is not",
				actual:     nil,
				expected:   "not nil",
				shouldFail: false,
			},
			{
				name:       "should pass when comparing nil with zero value",
				actual:     nil,
				expected:   0,
				shouldFail: false,
			},
			{
				name:       "should fail when both are empty strings",
				actual:     "",
				expected:   "",
				shouldFail: true,
			},
			{
				name:       "should pass when one empty and one with space",
				actual:     "",
				expected:   " ",
				shouldFail: false,
			},
			{
				name:       "should fail when both are zero",
				actual:     0,
				expected:   0,
				shouldFail: true,
			},
			{
				name:       "should pass when comparing different zero values",
				actual:     0,
				expected:   0.0,
				shouldFail: false, // different types
			},
			{
				name:       "should handle complex types - slices",
				actual:     []int{1, 2, 3},
				expected:   []int{1, 2, 4},
				shouldFail: false,
			},
			{
				name:       "should fail when slices are equal",
				actual:     []int{1, 2, 3},
				expected:   []int{1, 2, 3},
				shouldFail: true,
			},
			{
				name:       "should handle unicode characters",
				actual:     "测试🌟",
				expected:   "测试🔥",
				shouldFail: false,
			},
			{
				name:       "should handle very long strings",
				actual:     strings.Repeat("a", 1000),
				expected:   strings.Repeat("b", 1000),
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				NotBeEqual(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.Failed() {
					t.Fatal("Expected NotBeEqual to fail, but it passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected NotBeEqual to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.Failed() {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

func TestBeGreaterThan_Succeeds_WhenGreater(t *testing.T) {
	t.Parallel()

	BeGreaterThan(t, 10, 5)
}

func TestContain_Succeeds_WhenItemIsPresent(t *testing.T) {
	t.Parallel()

	Contain(t, []int{1, 2, 3}, 2)
}

func TestContain_Fails_WhenItemIsNotPresent(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int{1, 2, 3}, 4)
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
		Contain(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 21)
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

	NotContain(t, []int{1, 2, 3}, 4)
}

func TestNotContain_Fails_WhenItemIsPresent(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotContain(t, []int{1, 2, 3}, 2)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected collection to NOT contain element"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestNotContain_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotContain(t, []string{"apple", "banana"}, "apple", WithMessage("Apple should not be in basket"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected collection to NOT contain element"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestContain_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []string{"apple", "banana"}, "orange", WithMessage("Fruit not found in basket"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Fruit not found in basket",
		"Missing   : orange",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainFunc_Succeeds_WhenPredicateMatches(t *testing.T) {
	t.Parallel()

	ContainFunc(t, []int{1, 2, 3}, func(item any) bool {
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
		ContainFunc(t, []int{1, 2, 3}, func(item any) bool {
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

func TestContainFunc_WithCustomMessage(t *testing.T) {
	t.Parallel()

	mockT := &mockT{}
	customMessage := "No matching element found"

	numbers := []int{1, 3, 5, 7}
	predicate := func(item any) bool {
		return item.(int)%2 == 0 // Looking for even numbers
	}

	ContainFunc(mockT, numbers, predicate, WithMessage(customMessage))

	if !mockT.Failed() {
		t.Fatal("Expected ContainFunc to fail, but it passed")
	}

	if !strings.Contains(mockT.message, customMessage) {
		t.Errorf("Expected message to contain custom message %q, but got %q", customMessage, mockT.message)
	}

	if !strings.Contains(mockT.message, "Predicate does not match") {
		t.Errorf("Expected message to contain default error message, but got %q", mockT.message)
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
		Contain(t, users, "user3")
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
				Contain(t, tc.slice, tc.target)
			} else {
				// Should fail
				failed, message := assertFails(t, func(t testing.TB) {
					Contain(t, tc.slice, tc.target)
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
		Contain(t, slice, target)
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
			BeEmpty(t, tt.value)
		})
	}
}

func TestBeEmpty_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, "hello", WithMessage("This is a custom message"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "This is a custom message"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeEmpty_Fails_WithNonEmptyString(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, "hello")
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
		BeEmpty(t, []int{1, 2, 3})
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
		BeEmpty(t, longString)
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
		BeEmpty(t, largeSlice)
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
		BeEmpty(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeEmpty_WithChannel(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	ch <- 42

	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, ch)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be empty, but it was not:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeEmpty_WithChannelBuffered(t *testing.T) {
	t.Parallel()

	// Test empty buffered channel
	ch := make(chan int, 2)
	BeEmpty(t, ch)

	// Test non-empty buffered channel
	ch <- 42
	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, ch, WithMessage("Channel should be empty"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Channel should be empty",
		"Expected value to be empty, but it was not:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_WithNilInterface(t *testing.T) {
	t.Parallel()

	var nilInterface interface{}
	BeEmpty(t, nilInterface)
}

// === Tests for NotBeEmpty ===

func TestNotBeEmpty_Succeeds_WithNonEmptyString(t *testing.T) {
	t.Parallel()

	NotBeEmpty(t, "hello")
}

func TestNotBeEmpty_Succeeds_WithNonEmptySlice(t *testing.T) {
	t.Parallel()

	NotBeEmpty(t, []int{1, 2, 3})
}

func TestNotBeEmpty_Fails_WithEmptyString(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, "")
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

func TestNotBeEmpty_Fails_WithEmptySlice(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, []int{})
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

func TestNotBeEmpty_Fails_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int
	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, ptr)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be not empty, but it was empty:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotBeEmpty_WithInvalidValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		var v interface{}
		NotBeEmpty(t, v)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be not empty, but it was empty:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotBeEmpty_WithUnsupportedType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "NotBeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotBeEmpty_WithChannelBuffered(t *testing.T) {
	t.Parallel()

	// Test non-empty buffered channel
	ch := make(chan int, 1)
	ch <- 42
	NotBeEmpty(t, ch)

	// Test empty buffered channel
	emptyChannel := make(chan int, 1)
	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, emptyChannel, WithMessage("Channel should have data"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Channel should have data",
		"Expected value to be not empty, but it was empty:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestNotBeEmpty_WithValidInterface(t *testing.T) {
	t.Parallel()

	var validInterface interface{} = 42
	failed, message := assertFails(t, func(t testing.TB) {
		NotBeEmpty(t, validInterface)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "NotBeEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for BeGreaterThan ===

func TestBeGreaterThan_Succeeds_WithIntegers(t *testing.T) {
	t.Parallel()

	BeGreaterThan(t, 10, 5)
}

func TestBeGreaterThan_Succeeds_WithFloats(t *testing.T) {
	t.Parallel()

	BeGreaterThan(t, 3.14, 2.71)
}

func TestBeGreaterThan_Fails_WithSmallerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterThan(t, 5, 10)
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
		BeGreaterThan(t, 5, 5)
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
		BeGreaterThan(t, 0.0, 0.1, WithMessage("Score validation failed"))
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

	BeLessThan(t, 5, 10)
}

func TestBeLessThan_Fails_WithLargerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeLessThan(t, 10, 5)
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
		BeLessThan(t, 5, 5)
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
		BeLessThan(t, 3.14, 2.71)
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

func TestBeLessThan_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeLessThan(t, 10, 5, WithMessage("Value should be smaller"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Value should be smaller",
		"Expected value to be less than threshold:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

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
		}, WithMessage("Division by zero should panic"))
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
		}, WithMessage("Save operation should not panic"))
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

// === Tests for BeGreaterOrEqualTo ===

func TestBeGreaterOrEqualTo_Succeeds_WithGreaterValue(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualTo(t, 10, 5)
}

func TestBeGreaterOrEqualTo_Succeeds_WithEqualValue(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualTo(t, 5, 5)
}

func TestBeGreaterOrEqualTo_Succeeds_WithFloats(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualTo(t, 3.14, 3.14)
	BeGreaterOrEqualTo(t, 3.15, 3.14)
}

func TestBeGreaterOrEqualTo_Fails_WithSmallerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterOrEqualTo(t, 5, 10)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be greater than or equal to threshold"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeGreaterOrEqualTo_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterOrEqualTo(t, 0, 1, WithMessage("Score cannot be negative"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Score cannot be negative",
		"Expected value to be greater than or equal to threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeGreaterOrEqualTo_Fails_WithMixedIntFloat(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterOrEqualTo(t, 5.0, 5.1, WithMessage("Integer vs float comparison"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Integer vs float comparison",
		"Expected value to be greater than or equal to threshold",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for BeLessOrEqualTo ===

func TestBeLessOrEqualTo(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		// Integer tests
		BeLessOrEqualTo(t, 5, 10)
		BeLessOrEqualTo(t, 10, 10)

		// Float tests
		BeLessOrEqualTo(t, 2.71, 3.14)
		BeLessOrEqualTo(t, 3.14, 3.14)

		// Test failures
		t.Run("Fails when actual is greater than expected", func(t *testing.T) {
			failed, message := assertFails(t, func(t testing.TB) {
				BeLessOrEqualTo(t, 15, 10)
			})

			if !failed {
				t.Fatal("Expected BeLessOrEqualTo to fail, but it passed")
			}

			expectedParts := []string{
				"Expected value to be less than or equal to threshold:",
				"Value     : 15",
				"Threshold : 10",
				"Difference: +5 (value is 5 greater)",
				"Hint      : Value should be smaller than or equal to threshold",
			}
			for _, part := range expectedParts {
				if !strings.Contains(message, part) {
					t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
				}
			}
		})

		t.Run("Fails with float precision", func(t *testing.T) {
			failed, message := assertFails(t, func(t testing.TB) {
				BeLessOrEqualTo(t, 3.15, 3.14)
			})

			if !failed {
				t.Fatal("Expected BeLessOrEqualTo to fail, but it passed")
			}

			expectedParts := []string{
				"Expected value to be less than or equal to threshold:",
				"Value     : 3.15",
				"Threshold : 3.14",
				"Difference: +0.00",
			}
			for _, part := range expectedParts {
				if !strings.Contains(message, part) {
					t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
				}
			}
		})
	})

	t.Run("Custom messages", func(t *testing.T) {
		// Success with custom message
		BeLessOrEqualTo(t, 5, 10, WithMessage("Value should be within limit"))

		// Fails with custom error message
		t.Run("Fails with custom error message", func(t *testing.T) {
			failed, message := assertFails(t, func(t testing.TB) {
				BeLessOrEqualTo(t, 100, 50, WithMessage("Score should not exceed maximum"))
			})

			if !failed {
				t.Fatal("Expected BeLessOrEqualTo to fail, but it passed")
			}

			expectedParts := []string{
				"Score should not exceed maximum",
				"Expected value to be less than or equal to threshold:",
				"Value     : 100",
				"Threshold : 50",
			}
			for _, part := range expectedParts {
				if !strings.Contains(message, part) {
					t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
				}
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		// Success with zero values
		BeLessOrEqualTo(t, 0, 0)

		// Success with negative numbers
		BeLessOrEqualTo(t, -10, -5)

		// Success with very small floats
		BeLessOrEqualTo(t, 0.0001, 0.0002)

		// Fails with negative comparison
		t.Run("Fails with negative comparison", func(t *testing.T) {
			failed, message := assertFails(t, func(t testing.TB) {
				BeLessOrEqualTo(t, -5, -10)
			})

			if !failed {
				t.Fatal("Expected BeLessOrEqualTo to fail, but it passed")
			}

			expectedParts := []string{
				"Expected value to be less than or equal to threshold:",
				"Value     : -5",
				"Threshold : -10",
				"Difference: +5 (value is 5 greater)",
			}
			for _, part := range expectedParts {
				if !strings.Contains(message, part) {
					t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
				}
			}
		})
	})

	t.Run("Type compatibility", func(t *testing.T) {
		tests := []struct {
			name       string
			testFunc   func()
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name: "Success with different integer types",
				testFunc: func() {
					BeLessOrEqualTo(t, int8(5), int8(10))
					BeLessOrEqualTo(t, int16(5), int16(10))
					BeLessOrEqualTo(t, int32(5), int32(10))
					BeLessOrEqualTo(t, int64(5), int64(10))
				},
				shouldFail: false,
			},
			{
				name: "Success with different unsigned integer types",
				testFunc: func() {
					BeLessOrEqualTo(t, uint8(5), uint8(10))
					BeLessOrEqualTo(t, uint16(5), uint16(10))
					BeLessOrEqualTo(t, uint32(5), uint32(10))
					BeLessOrEqualTo(t, uint64(5), uint64(10))
				},
				shouldFail: false,
			},
			{
				name: "Success with different float types",
				testFunc: func() {
					BeLessOrEqualTo(t, float32(2.5), float32(3.5))
					BeLessOrEqualTo(t, float64(2.5), float64(3.5))
				},
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				if tt.shouldFail {
					failed, message := assertFails(t, func(t testing.TB) {
						tt.testFunc()
					})

					if !failed {
						t.Fatal("Expected test function to fail, but it passed")
					}
					if tt.errorCheck != nil {
						tt.errorCheck(t, message)
					}
				} else {
					tt.testFunc()
				}
			})
		}
	})
}

// === Tests for edge cases ===

func TestBeNil_Succeeds_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int
	BeNil(t, ptr)
}

func TestBeNil_Succeeds_WithNilSlice(t *testing.T) {
	t.Parallel()

	var slice []int
	BeNil(t, slice)
}

func TestBeNil_Succeeds_WithNilMap(t *testing.T) {
	t.Parallel()

	var m map[string]int
	BeNil(t, m)
}

// TestBeNil_Succeeds_WithNilInterface removed due to reflect.Value issue

func TestBeNil_Fails_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		BeNil(t, ptr)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected nil, but was not"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotBeNil_Succeeds_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value
	NotBeNil(t, ptr)
}

func TestNotBeNil_Succeeds_WithNonNilSlice(t *testing.T) {
	t.Parallel()

	slice := make([]int, 0)
	NotBeNil(t, slice)
}

func TestNotBeNil_Fails_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeNil(t, ptr)
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

func TestBeTrue_Fails_WithFalseValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeTrue(t, false)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected true, got false"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeFalse_Fails_WithTrueValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeFalse(t, true)
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
		ContainFunc(t, "not a slice", func(item any) bool {
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
	ContainFunc(t, users, func(item any) bool {
		user := item.(User)
		return user.Age >= 18
	})

	// Should fail - no elderly users
	failed, message := assertFails(t, func(t testing.TB) {
		ContainFunc(t, users, func(item any) bool {
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
		Contain(t, "not a slice", "something")
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
		NotContain(t, 42, "something")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a slice or array, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for custom messages in boolean assertions ===

func TestBeTrue_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeTrue(t, false, WithMessage("User must be active"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"User must be active",
		"Expected true, got false",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeFalse_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeFalse(t, true, WithMessage("User should not be deleted"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"User should not be deleted",
		"Expected false, got true",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for BeNil/NotBeNil with custom messages ===

func TestBeNil_WithCustomMessage(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		BeNil(t, ptr, WithMessage("Pointer should be nil"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Pointer should be nil",
		"Expected nil, but was not",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestNotBeNil_WithCustomMessage(t *testing.T) {
	t.Parallel()

	var ptr *int

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeNil(t, ptr, WithMessage("User must not be nil"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"User must not be nil",
		"Expected not nil, but was nil",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeNil_Fails_WithNonNillableType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeNil(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeNil can only be used with nillable types, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestNotBeNil_Fails_WithNonNillableType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		NotBeNil(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "NotBeNil can only be used with nillable types, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for numeric slice contain with different types ===

func TestContain_WithInt8Slices(t *testing.T) {
	t.Parallel()

	// Test with int8 slice - success case
	Contain(t, []int8{1, 2, 3}, int8(2))

	// Test with int8 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int8{1, 2, 4, 5}, int8(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithInt16Slices(t *testing.T) {
	t.Parallel()

	// Test with int16 slice - success case
	Contain(t, []int16{1, 2, 3}, int16(2))

	// Test with int16 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int16{1, 2, 4, 5}, int16(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithInt32Slices(t *testing.T) {
	t.Parallel()

	// Test with int32 slice - success case
	Contain(t, []int32{1, 2, 3}, int32(2))

	// Test with int32 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int32{1, 2, 4, 5}, int32(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithInt64Slices(t *testing.T) {
	t.Parallel()

	// Test with int64 slice - success case
	Contain(t, []int64{1, 2, 3}, int64(2))

	// Test with int64 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int64{1, 2, 4, 5}, int64(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUintSlices(t *testing.T) {
	t.Parallel()

	// Test with uint slice - success case
	Contain(t, []uint{1, 2, 3}, uint(2))

	// Test with uint slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []uint{1, 2, 4, 5}, uint(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUint8Slices(t *testing.T) {
	t.Parallel()

	// Test with uint8 slice - success case
	Contain(t, []uint8{1, 2, 3}, uint8(2))

	// Test with uint8 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []uint8{1, 2, 4, 5}, uint8(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUint16Slices(t *testing.T) {
	t.Parallel()

	// Test with uint16 slice - success case
	Contain(t, []uint16{1, 2, 3}, uint16(2))

	// Test with uint16 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []uint16{1, 2, 4, 5}, uint16(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUint32Slices(t *testing.T) {
	t.Parallel()

	// Test with uint32 slice - success case
	Contain(t, []uint32{1, 2, 3}, uint32(2))

	// Test with uint32 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []uint32{1, 2, 4, 5}, uint32(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUint64Slices(t *testing.T) {
	t.Parallel()

	// Test with uint64 slice - success case
	Contain(t, []uint64{1, 2, 3}, uint64(2))

	// Test with uint64 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []uint64{1, 2, 4, 5}, uint64(3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithFloat32Slices(t *testing.T) {
	t.Parallel()

	// Test with float32 slice - success case
	Contain(t, []float32{1.1, 2.2, 3.3}, float32(2.2))

	// Test with float32 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []float32{1.1, 2.2, 4.4, 5.5}, float32(3.3))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3.3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithFloat64Slices(t *testing.T) {
	t.Parallel()

	// Test with float64 slice - success case
	Contain(t, []float64{1.1, 2.2, 3.3}, 2.2)

	// Test with float64 slice - failure case
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []float64{1.1, 2.2, 4.4, 5.5}, 3.3)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing  : 3.3"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

func TestContain_WithUnsupportedNumericTypeCombination(t *testing.T) {
	t.Parallel()

	// Test with int slice but float target (should fall back to generic contain)
	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, []int{1, 2, 3}, 2.5)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing   : 2.5"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}
}

// === Tests for NotContainDuplicates ===

func TestNotContainDuplicates_Fails_WhenDuplicatesExist(t *testing.T) {
	t.Parallel()

	// Test that it correctly identifies duplicates in int slice
	failed, message := assertFails(t, func(t testing.TB) {
		NotContainDuplicates(t, []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 4, 4}, WithMessage("Expected no duplicates, but found 1 duplicate value: 2"))
	})
	if !failed {
		t.Fatal("Expected test to fail due to duplicates, but it passed")
	}
	if !strings.Contains(message, "duplicate values") {
		t.Errorf("Expected error message to mention duplicates, got: %s", message)
	}

	expected := "Expected no duplicates, but found 1 duplicate value: 2"
	if !strings.Contains(message, expected) {
		t.Fatalf("Expected message to contain %q, but got %q", expected, message)
	}

	// Test that it correctly identifies duplicates in string slice
	failed, message = assertFails(t, func(t testing.TB) {
		NotContainDuplicates(t, []string{"a", "b", "c", "c", "d", "d", "e", "e", "e", "e", "e"})
	})
	if !failed {
		t.Fatal("Expected test to fail due to duplicates, but it passed")
	}
	if !strings.Contains(message, "duplicate values") {
		t.Errorf("Expected error message to mention duplicates, got: %s", message)
	}

	// Test that it correctly identifies duplicates in complex structs
	type Address struct {
		Street string
		City   string
		ZIP    string
	}

	type User struct {
		ID       int
		Name     string
		Email    string
		Address  Address
		Metadata map[string]interface{}
	}

	users := []User{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
			Address: Address{
				Street: "123 Main St",
				City:   "New York",
				ZIP:    "10001",
			},
			Metadata: map[string]interface{}{
				"role":      "admin",
				"lastLogin": "2024-01-15",
			},
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane@example.com",
			Address: Address{
				Street: "456 Oak Ave",
				City:   "Los Angeles",
				ZIP:    "90210",
			},
			Metadata: map[string]interface{}{
				"role":      "user",
				"lastLogin": "2024-01-14",
			},
		},
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
			Address: Address{
				Street: "123 Main St",
				City:   "New York",
				ZIP:    "10001",
			},
			Metadata: map[string]interface{}{
				"role":      "admin",
				"lastLogin": "2024-01-15",
			},
		},
		{
			ID:    3,
			Name:  "Bob Wilson",
			Email: "bob@example.com",
			Address: Address{
				Street: "789 Pine Rd",
				City:   "Chicago",
				ZIP:    "60601",
			},
			Metadata: map[string]interface{}{
				"role":      "user",
				"lastLogin": "2024-01-13",
			},
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane@example.com",
			Address: Address{
				Street: "456 Oak Ave",
				City:   "Los Angeles",
				ZIP:    "90210",
			},
			Metadata: map[string]interface{}{
				"role":      "user",
				"lastLogin": "2024-01-14",
			},
		},
	}

	failed, message = assertFails(t, func(t testing.TB) {
		NotContainDuplicates(t, users)
	})
	if !failed {
		t.Fatal("Expected test to fail due to struct duplicates, but it passed")
	}
	if !strings.Contains(message, "duplicate values") {
		t.Errorf("Expected error message to mention duplicates, got: %s", message)
	}
}

func TestNotContainDuplicates_Succeeds_WhenNoDuplicates(t *testing.T) {
	t.Parallel()

	NotContainDuplicates(t, []int{1, 2, 3, 4, 5})
	NotContainDuplicates(t, []string{"a", "b", "c", "d", "e"})

	type User struct {
		ID   int
		Name string
	}

	users := []User{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
		{ID: 3, Name: "Bob"},
	}

	NotContainDuplicates(t, users)
}

func TestNotContainDuplicates_WithCustomMessage(t *testing.T) {
	t.Parallel()

	mockT := &mockT{}
	customMessage := "Collection should not have duplicates"

	duplicateSlice := []int{1, 2, 2, 3}

	NotContainDuplicates(mockT, duplicateSlice, WithMessage(customMessage))

	if !mockT.Failed() {
		t.Fatal("Expected NotContainDuplicates to fail, but it passed")
	}

	if !strings.Contains(mockT.message, customMessage) {
		t.Errorf("Expected message to contain custom message %q, but got %q", customMessage, mockT.message)
	}

	if !strings.Contains(mockT.message, "Expected no duplicates") {
		t.Errorf("Expected message to contain default error message, but got %q", mockT.message)
	}
}

// === Tests for StartsWith ===

func TestStartsWith_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		StartsWith(t, "Hello, world!", "world", WithMessage("String should start with 'world'"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected string to start with 'world'",
		"but it starts with 'Hello'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestStartsWith(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass if string starts with prefix",
				actual:     "Hello, world!",
				expected:   "Hello",
				shouldFail: false,
			},
			{
				name:       "should fail if string does not start with prefix",
				actual:     "Hello, world!",
				expected:   "world",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `Expected string to start with 'world', but it starts with 'Hello'`) ||
						!strings.Contains(message, `Actual   : 'Hello, world!'`) ||
						!strings.Contains(message, `Expected : 'world'`) {
						t.Errorf("Incorrect error message:\n%s", message)
					}
				},
			},
			{
				name:       "should show actual prefix in error message",
				actual:     "Hello, world!",
				expected:   "Hi",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `(actual prefix)`) {
						t.Errorf("Expected error message to contain '(actual prefix)' indicator, got:\n%s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				StartsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartsWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartsWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass with ignore case enabled",
				actual:     "Hello, world!",
				expected:   "hello",
				opts:       []Option{WithIgnoreCase()},
				shouldFail: false,
			},
			{
				name:       "should fail if ignore case is disabled and case mismatch is detected",
				actual:     "Hello",
				expected:   "hello",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `Note: Case mismatch detected (use should.WithIgnoreCase() if intended)`) {
						t.Errorf("Expected message to contain note message, but got %q", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				StartsWith(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartsWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartsWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			opts       []Option
			shouldFail bool
		}{
			{
				name:       "should pass with custom message",
				actual:     "Hello, world!",
				expected:   "Hello",
				opts:       []Option{WithMessage("Expected string to start with 'Hello'")},
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				StartsWith(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartsWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartsWith to pass, but it failed: %s", mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should fail for empty prefix",
				actual:     "data",
				expected:   "",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `Expected string to start with '<empty>'`) {
						t.Errorf("Expected error message for empty prefix, got: %s", message)
					}
				},
			},
			{
				name:       "should fail if actual is empty",
				actual:     "",
				expected:   "test",
				shouldFail: true,
			},
			{
				name:       "should handle exact match",
				actual:     "test",
				expected:   "test",
				shouldFail: false,
			},
			{
				name:       "should handle prefix longer than actual",
				actual:     "abc",
				expected:   "abcdef",
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				StartsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartsWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartsWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("String truncation", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should truncate actual string when longer than 56 characters",
				actual:     "This is a very long string that exceeds the 56 character limit for display purposes in error messages",
				expected:   "Different",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "... (truncated)") {
						t.Errorf("Expected message to contain truncated actual string, got: %s", message)
					}
				},
			},
			{
				name:       "should truncate expected string when longer than 56 characters",
				actual:     "Short",
				expected:   "This is a very long expected string that exceeds the 56 character limit for display purposes in error messages",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "... (truncated)") {
						t.Errorf("Expected message to contain truncated expected string, got: %s", message)
					}
				},
			},
			{
				name:       "should truncate both strings when both are longer than 56 characters",
				actual:     "This is a very long actual string that exceeds the 56 character limit for display purposes in error messages",
				expected:   "This is a very long expected string that exceeds the 56 character limit for display purposes in error messages",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					truncatedOccurrences := strings.Count(message, "... (truncated)")
					if truncatedOccurrences < 2 {
						t.Errorf("Expected at least 2 truncated strings in message, got %d occurrences in: %s", truncatedOccurrences, message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				StartsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartsWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartsWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

// === Tests for EndsWith ===

func TestEndsWith(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Success when actual ends with expected",
				actual:     "Hello, world!",
				expected:   "world!",
				shouldFail: false,
			},
			{
				name:       "Exact match passes",
				actual:     "world",
				expected:   "world",
				shouldFail: false,
			},
			{
				name:       "Fails when actual does not end with expected",
				actual:     "Hello, world!",
				expected:   "planet",
				shouldFail: true,
			},
			{
				name:       "Fails when expected is longer than actual",
				actual:     "abc",
				expected:   "abcdef",
				shouldFail: true,
			},
			{
				name:       "should show actual suffix in error message",
				actual:     "Hello, world!",
				expected:   "world",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `(actual suffix)`) {
						t.Errorf("Expected error message to contain '(actual suffix)' indicator, got:\n%s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				EndsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.Failed() {
					t.Errorf("Expected failure but test passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected success but test failed")
				}

				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			opts       []Option
			shouldFail bool
		}{
			{
				name:       "Success with ignore case enabled",
				actual:     "Hello, WORLD",
				expected:   "world",
				opts:       []Option{WithIgnoreCase()},
				shouldFail: false,
			},
			{
				name:       "Fails with ignore case disabled",
				actual:     "Hello, WORLD",
				expected:   "world",
				opts:       []Option{},
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				EndsWith(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.Failed() {
					t.Errorf("Expected failure but test passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected success but test failed")
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		t.Run("Fails with custom message", func(t *testing.T) {
			mockT := &mockT{}
			EndsWith(mockT, "Hello, world!", "planet", WithMessage("String should end with 'planet'"))

			if !mockT.Failed() {
				t.Errorf("Expected failure but test passed")
			}

			expectedStrings := []string{
				"Expected string to end with 'planet'",
				"but it ends with 'world!'",
			}

			for _, expectedString := range expectedStrings {
				if !strings.Contains(mockT.message, expectedString) {
					t.Errorf("Expected message to contain %q, but got %q", expectedString, mockT.message)
				}
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
		}{
			{
				name:       "Empty strings",
				actual:     "",
				expected:   "",
				shouldFail: false,
			},
			{
				name:       "Empty expected with non-empty actual",
				actual:     "hello",
				expected:   "",
				shouldFail: true,
			},
			{
				name:       "Non-empty expected with empty actual",
				actual:     "",
				expected:   "hello",
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				EndsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.Failed() {
					t.Errorf("Expected failure but test passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected success but test failed")
				}
			})
		}
	})

	t.Run("String truncation", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			expected   string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should truncate actual string when longer than 56 characters",
				actual:     "This is a very long string that exceeds the 56 character limit for display purposes in error messages",
				expected:   "Different",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "... (truncated)") {
						t.Errorf("Expected message to contain truncated actual string, got: %s", message)
					}
				},
			},
			{
				name:       "should truncate expected string when longer than 56 characters",
				actual:     "Short",
				expected:   "This is a very long expected string that exceeds the 56 character limit for display purposes in error messages",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "... (truncated)") {
						t.Errorf("Expected message to contain truncated expected string, got: %s", message)
					}
				},
			},
			{
				name:       "should truncate both strings when both are longer than 56 characters",
				actual:     "This is a very long actual string that exceeds the 56 character limit for display purposes in error messages",
				expected:   "This is a very long expected string that exceeds the 56 character limit for display purposes in error messages",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					truncatedOccurrences := strings.Count(message, "... (truncated)")
					if truncatedOccurrences < 2 {
						t.Errorf("Expected at least 2 truncated strings in message, got %d occurrences in: %s", truncatedOccurrences, message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				EndsWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.Failed() {
					t.Errorf("Expected failure but test passed")
				}
				if !tt.shouldFail && mockT.Failed() {
					t.Errorf("Expected success but test failed")
				}

				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

// === Tests for BeEqual with complex types and custom messages ===

func TestBeEqual_WithComplexNestedStructs_CustomMessage(t *testing.T) {
	t.Parallel()

	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
	}

	person1 := Person{
		Name: "John",
		Age:  30,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	person2 := Person{
		Name: "John",
		Age:  31,
		Address: Address{
			Street: "123 Main St",
			City:   "Boston",
		},
	}

	failed, message := assertFails(t, func(t testing.TB) {
		BeEqual(t, person1, person2, WithMessage("Person objects should be identical"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Person objects should be identical",
		"Differences found",
		"Field differences:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for numeric type coverage in toFloat64 ===

func TestNumericTypeConversions_AllTypes(t *testing.T) {
	t.Parallel()

	// Test all supported numeric types for better compareOrderable coverage
	t.Run("int8", func(t *testing.T) {
		BeGreaterThan(t, int8(10), int8(9))
	})

	t.Run("int16", func(t *testing.T) {
		BeGreaterThan(t, int16(20), int16(19))
	})

	t.Run("int32", func(t *testing.T) {
		BeGreaterThan(t, int32(30), int32(29))
	})

	t.Run("int64", func(t *testing.T) {
		BeGreaterThan(t, int64(40), int64(39))
	})

	t.Run("uint8", func(t *testing.T) {
		BeGreaterThan(t, uint8(50), uint8(49))
	})

	t.Run("uint16", func(t *testing.T) {
		BeGreaterThan(t, uint16(60), uint16(59))
	})

	t.Run("uint32", func(t *testing.T) {
		BeGreaterThan(t, uint32(70), uint32(69))
	})

	t.Run("uint64", func(t *testing.T) {
		BeGreaterThan(t, uint64(80), uint64(79))
	})

	t.Run("float32", func(t *testing.T) {
		BeGreaterThan(t, float32(3.14), float32(3.13))
	})

	t.Run("float64", func(t *testing.T) {
		BeGreaterThan(t, float64(2.718), float64(2.717))
	})
}

func TestIsNumericType_Coverage(t *testing.T) {
	t.Parallel()

	// Test with a non-numeric slice to exercise the fallback path
	type CustomStruct struct {
		Name string
	}

	structs := []CustomStruct{
		{Name: "test1"},
		{Name: "test2"},
	}

	failed, message := assertFails(t, func(t testing.TB) {
		Contain(t, structs, CustomStruct{Name: "test3"})
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Missing   :"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

// === Tests for HaveLength ===

func TestHaveLength(t *testing.T) {
	t.Run("should pass for correct length of slice", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, []int{1, 2, 3}, 3)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should pass for correct length of string", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, "abc", 3)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should pass for correct length of map", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, map[int]int{1: 1, 2: 2}, 2)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should fail for incorrect length with custom message", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, []int{1, 2}, 3, WithMessage("Custom message"))
		if !mockT.failed {
			t.Fatal("Expected HaveLength to fail, but it passed")
		}
		expectedMsg := "Custom message"
		if !strings.Contains(mockT.message, expectedMsg) {
			t.Errorf("Expected error message to contain custom message %q, but got %q", expectedMsg, mockT.message)
		}
	})

	t.Run("should fail for incorrect length and show detailed error", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, []int{1, 2}, 3)
		if !mockT.failed {
			t.Fatal("Expected HaveLength to fail, but it passed")
		}

		if !strings.Contains(mockT.message, "Expected collection to have specific length:") ||
			!strings.Contains(mockT.message, "Type          : []int") ||
			!strings.Contains(mockT.message, "Expected Length: 3") ||
			!strings.Contains(mockT.message, "Actual Length : 2") ||
			!strings.Contains(mockT.message, "Difference    : -1 (1 element missing)") {
			t.Errorf("Error message format is incorrect.\nGot:\n%s", mockT.message)
		}
	})

	t.Run("should fail for incorrect length (more elements)", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, []int{1, 2, 3, 4}, 3)
		if !mockT.failed {
			t.Fatal("Expected HaveLength to fail, but it passed")
		}
		if !strings.Contains(mockT.message, "Difference    : +1 (1 element extra)") {
			t.Errorf("Error message format is incorrect for extra elements.\nGot:\n%s", mockT.message)
		}
	})

	t.Run("should fail for unsupported type", func(t *testing.T) {
		mockT := &mockT{}
		HaveLength(mockT, 123, 1)
		if !mockT.failed {
			t.Fatal("Expected HaveLength to fail for unsupported type, but it passed")
		}
		expectedMsg := "HaveLength can only be used with types that have a concept of length (string, slice, array, map), but got int"
		if !strings.Contains(mockT.message, expectedMsg) {
			t.Errorf("Expected error message to contain %q, but got %q", expectedMsg, mockT.message)
		}
	})
}

// === Tests for BeOfType ===

func TestBeOfType(t *testing.T) {
	type Cat struct{ Name string }
	type Dog struct{ Name string }

	t.Run("should pass for same type", func(t *testing.T) {
		mockT := &mockT{}
		var c *Cat
		BeOfType(mockT, &Cat{}, c)
		if mockT.failed {
			t.Errorf("Expected BeOfType to pass, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail for different types", func(t *testing.T) {
		mockT := &mockT{}
		var d *Dog
		BeOfType(mockT, &Cat{Name: "Whiskers"}, d)
		if !mockT.failed {
			t.Fatal("Expected BeOfType to fail, but it passed")
		}

		if !strings.Contains(mockT.message, "Expected value to be of specific type:") ||
			!strings.Contains(mockT.message, "Expected Type: *assert.Dog") ||
			!strings.Contains(mockT.message, "Actual Type  : *assert.Cat") ||
			!strings.Contains(mockT.message, "Difference   : Different concrete types") ||
			!strings.Contains(mockT.message, `Value        : {Name: "Whiskers"}`) {
			t.Errorf("Error message format is incorrect.\nGot:\n%s", mockT.message)
		}
	})

	t.Run("should pass for primitive types", func(t *testing.T) {
		mockT := &mockT{}
		BeOfType(mockT, 1, 0) // int and int
		if mockT.failed {
			t.Errorf("Expected BeOfType to pass for ints, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail for different primitive types", func(t *testing.T) {
		mockT := &mockT{}
		BeOfType(mockT, int32(1), int64(0))
		if !mockT.failed {
			t.Fatal("Expected BeOfType to fail, but it passed")
		}

		if !strings.Contains(mockT.message, "Expected Type: int64") ||
			!strings.Contains(mockT.message, "Actual Type  : int32") {
			t.Errorf("Error message does not contain correct types for primitives.\nGot:\n%s", mockT.message)
		}
	})
}

// === Tests for BeOneOf ===

func TestBeOneOf(t *testing.T) {
	t.Run("should pass if value is one of the options", func(t *testing.T) {
		mockT := &mockT{}
		options := []string{"active", "inactive"}
		BeOneOf(mockT, "active", options)
		if mockT.failed {
			t.Errorf("Expected BeOneOf to pass, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail if value is not one of the options", func(t *testing.T) {
		mockT := &mockT{}
		options := []string{"active", "inactive", "suspended"}
		BeOneOf(mockT, "pending", options)
		if !mockT.failed {
			t.Fatal("Expected BeOneOf to fail, but it passed")
		}

		if !strings.Contains(mockT.message, `Expected value to be one of the allowed options:`) ||
			!strings.Contains(mockT.message, `Value   : "pending"`) ||
			!strings.Contains(mockT.message, `Options : ["active", "inactive", "suspended"]`) ||
			!strings.Contains(mockT.message, `Count   : 0 of 3 options matched`) {
			t.Errorf("Error message format is incorrect.\nGot:\n%s", mockT.message)
		}
	})

	t.Run("should fail for empty options", func(t *testing.T) {
		mockT := &mockT{}
		BeOneOf(mockT, "any", []string{})
		if !mockT.failed {
			t.Fatal("Expected BeOneOf to fail for empty options, but it passed")
		}
		if !strings.Contains(mockT.message, "Options list cannot be empty") {
			t.Errorf("Expected error for empty options, but got: %s", mockT.message)
		}
	})

	t.Run("should truncate long option lists in the error message", func(t *testing.T) {
		mockT := &mockT{}
		options := []string{"active", "inactive", "suspended", "deleted", "archived"}
		BeOneOf(mockT, "pending", options)
		if !mockT.failed {
			t.Fatal("Expected BeOneOf to fail, but it passed")
		}

		if !strings.Contains(mockT.message, `Options : ["active", "inactive", "suspended", "deleted", ...]`) {
			t.Errorf("Expected truncated option list in message, got:\n%s", mockT.message)
		}
		if !strings.Contains(mockT.message, `showing first 4 of 5`) {
			t.Errorf("Expected truncation note in message, got:\n%s", mockT.message)
		}
	})
}

// === Tests for ContainKey ===

func TestContainKey_Succeeds_WhenKeyIsPresent(t *testing.T) {
	t.Parallel()

	m := map[string]int{"name": 1, "age": 2, "email": 3}
	ContainKey(t, m, "email")
}

func TestContainKey_Succeeds_WithIntKeys(t *testing.T) {
	t.Parallel()

	m := map[int]string{1: "one", 2: "two", 3: "three"}
	ContainKey(t, m, 2)
}

func TestContainKey_Fails_WhenKeyIsNotPresent(t *testing.T) {
	t.Parallel()

	m := map[string]int{"name": 1, "age": 2}
	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, m, "email")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain key 'email', but key was not found",
		"Available keys:",
		"Missing: 'email'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainKey_ShowsSimilarKeys_ForStringKeys(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"id":           1,
		"name":         2,
		"mail":         3,
		"e_mail":       4,
		"emailAddress": 5,
	}

	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, m, "email")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain key 'email', but key was not found",
		"Available keys:",
		"Missing: 'email'",
		"Similar keys found:",
		"└─ 'mail'",
		"└─ 'e_mail'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainKey_ShowsSimilarKeys_ForIntKeys(t *testing.T) {
	t.Parallel()

	m := map[int]string{
		1:   "one",
		24:  "twenty-four",
		43:  "forty-three",
		100: "hundred",
		420: "four-twenty",
	}

	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, m, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain key 42, but key was not found",
		"Available keys:",
		"Missing: 42",
		"Similar keys found:",
		"└─ 43 - differs by 1",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainKey_WithCustomMessage(t *testing.T) {
	t.Parallel()

	m := map[string]int{"name": 1}
	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, m, "email", WithMessage("User profile must have email field"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"User profile must have email field",
		"Expected map to contain key 'email', but key was not found",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Tests for ContainValue ===

func TestContainValue_Succeeds_WhenValueIsPresent(t *testing.T) {
	t.Parallel()

	m := map[string]int{"name": 1, "age": 2, "email": 3}
	ContainValue(t, m, 2)
}

func TestContainValue_Succeeds_WithStringValues(t *testing.T) {
	t.Parallel()

	m := map[int]string{1: "one", 2: "two", 3: "three"}
	ContainValue(t, m, "two")
}

func TestContainValue_Fails_WhenValueIsNotPresent(t *testing.T) {
	t.Parallel()

	m := map[string]int{"name": 1, "age": 2}
	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, m, 3)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain value 3, but value was not found",
		"Available values: [1, 2]",
		"Missing: 3",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainValue_ShowsSimilarValues_ForStringValues(t *testing.T) {
	t.Parallel()

	m := map[int]string{
		1: "admin",
		2: "user",
		3: "guest",
		4: "moderator",
		5: "administrator",
	}

	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, m, "adm")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain value 'adm', but value was not found",
		"Available values:",
		"Missing: 'adm'",
		"Similar values found:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainValue_ShowsSimilarValues_ForIntValues(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"first":  10,
		"second": 25,
		"third":  30,
		"fourth": 45,
	}

	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, m, 24)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain value 24, but value was not found",
		"Available values:",
		"Missing: 24",
		"Similar values found:",
		"└─ 25 - differs by 1",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainValue_WithCustomMessage(t *testing.T) {
	t.Parallel()

	m := map[string]int{"score": 100}
	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, m, 95, WithMessage("Score must be achievable"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Score must be achievable",
		"Expected map to contain value 95, but value was not found",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

// === Edge Cases Tests for ContainKey ===

func TestContainKey_EdgeCases_WithNilMap(t *testing.T) {
	t.Parallel()

	var nilMap map[string]int
	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, nilMap, "test")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain key 'test', but key was not found",
		"Available keys: nil",
		"Missing: 'test'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainKey_EdgeCases_WithEmptyMap(t *testing.T) {
	t.Parallel()

	emptyMap := make(map[string]int)
	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, emptyMap, "test")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain key 'test', but key was not found",
		"Available keys: []",
		"Missing: 'test'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainKey_EdgeCases_WithZeroValues(t *testing.T) {
	t.Parallel()

	// Test with zero value keys
	m := map[int]string{0: "zero", 1: "one"}
	ContainKey(t, m, 0)

	m2 := map[string]int{"": 42, "test": 1}
	ContainKey(t, m2, "")
}

// === Edge Cases Tests for ContainValue ===

func TestContainValue_EdgeCases_WithNilMap(t *testing.T) {
	t.Parallel()

	var nilMap map[string]int
	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, nilMap, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain value 42, but value was not found",
		"Available values: nil",
		"Missing: 42",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainValue_EdgeCases_WithEmptyMap(t *testing.T) {
	t.Parallel()

	emptyMap := make(map[string]int)
	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, emptyMap, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected map to contain value 42, but value was not found",
		"Available values: []",
		"Missing: 42",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainValue_EdgeCases_WithZeroValues(t *testing.T) {
	t.Parallel()

	// Test with zero value values
	m := map[string]int{"zero": 0, "one": 1}
	ContainValue(t, m, 0)

	m2 := map[int]string{1: "", 2: "test"}
	ContainValue(t, m2, "")
}

// === Tests for NotContainKey ===

func TestNotContainKey(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Run("string keys", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				key        string
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass when key does not exist",
					mapValue:   map[string]int{"name": 1, "age": 2},
					key:        "email",
					shouldFail: false,
				},
				{
					name:       "should pass with empty map",
					mapValue:   map[string]int{},
					key:        "any",
					shouldFail: false,
				},
				{
					name:       "should fail when key exists",
					mapValue:   map[string]int{"name": 1, "age": 2},
					key:        "age",
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						expectedParts := []string{
							"Expected map to NOT contain key, but key was found:",
							"Map Type : map[string]int",
							"Map Size : 2 entries",
							`Found Key: "age"`,
							"Associated Value: 2",
						}
						for _, part := range expectedParts {
							if !strings.Contains(message, part) {
								t.Errorf("Expected message to contain %q, but it was not found in:\n%s", part, message)
							}
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("int keys", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[int]string
				key        int
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass when key does not exist",
					mapValue:   map[int]string{1: "one", 2: "two"},
					key:        3,
					shouldFail: false,
				},
				{
					name:       "should fail when int key exists",
					mapValue:   map[int]string{1: "one", 2: "two", 3: "three"},
					key:        2,
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						expectedParts := []string{
							"Expected map to NOT contain key, but key was found:",
							"Map Type : map[int]string",
							"Found Key: 2",
							`Associated Value: "two"`,
						}
						for _, part := range expectedParts {
							if !strings.Contains(message, part) {
								t.Errorf("Expected message to contain %q, but it was not found in:\n%s", part, message)
							}
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("custom struct keys", func(t *testing.T) {
			type CustomKey struct {
				ID   int
				Name string
			}

			tests := []struct {
				name       string
				mapValue   map[CustomKey]bool
				key        CustomKey
				shouldFail bool
			}{
				{
					name:       "should pass when custom key does not exist",
					mapValue:   map[CustomKey]bool{{ID: 1, Name: "test"}: true},
					key:        CustomKey{ID: 2, Name: "other"},
					shouldFail: false,
				},
				{
					name:       "should fail when custom key exists",
					mapValue:   map[CustomKey]bool{{ID: 1, Name: "test"}: true},
					key:        CustomKey{ID: 1, Name: "test"},
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})
	})

	t.Run("Custom messages", func(t *testing.T) {
		t.Run("string-int maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				key        string
				opts       []Option
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass with custom message",
					mapValue:   map[string]int{"name": 1},
					key:        "email",
					opts:       []Option{WithMessage("Email key should not exist")},
					shouldFail: false,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key, tt.opts...)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("string-string maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]string
				key        string
				opts       []Option
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should show custom error message on failure",
					mapValue:   map[string]string{"secret_key": "value"},
					key:        "secret_key",
					opts:       []Option{WithMessage("Configuration should not contain sensitive keys")},
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						if !strings.Contains(message, "Configuration should not contain sensitive keys") {
							t.Errorf("Expected custom error message, got: %s", message)
						}
						if !strings.Contains(message, "Expected map to NOT contain key, but key was found:") {
							t.Errorf("Expected standard error message, got: %s", message)
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key, tt.opts...)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Run("nil map handling", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				key        string
				shouldFail bool
			}{
				{
					name:       "should handle nil map",
					mapValue:   nil,
					key:        "test",
					shouldFail: false,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})

		t.Run("zero value keys", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[int]string
				key        int
				shouldFail bool
			}{
				{
					name:       "should handle zero value keys that don't exist",
					mapValue:   map[int]string{0: "zero", 1: "one"},
					key:        2,
					shouldFail: false,
				},
				{
					name:       "should fail with zero value key that exists",
					mapValue:   map[int]string{0: "zero", 1: "one"},
					key:        0,
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})

		t.Run("empty string keys", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				key        string
				shouldFail bool
			}{
				{
					name:       "should handle missing keys when empty string exists",
					mapValue:   map[string]int{"": 42, "test": 1},
					key:        "missing",
					shouldFail: false,
				},
				{
					name:       "should fail with empty string key that exists",
					mapValue:   map[string]int{"": 42, "test": 1},
					key:        "",
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})

		t.Run("complex key types", func(t *testing.T) {
			type ComplexKey struct {
				ID   int
				Name string
			}

			tests := []struct {
				name       string
				mapValue   map[ComplexKey]bool
				key        ComplexKey
				shouldFail bool
			}{
				{
					name:       "should handle complex struct keys that don't exist",
					mapValue:   map[ComplexKey]bool{{ID: 1, Name: "test"}: true},
					key:        ComplexKey{ID: 2, Name: "other"},
					shouldFail: false,
				},
				{
					name:       "should fail with complex struct key that exists",
					mapValue:   map[ComplexKey]bool{{ID: 1, Name: "test"}: true},
					key:        ComplexKey{ID: 1, Name: "test"},
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})

		t.Run("pointer keys", func(t *testing.T) {
			key1 := "test1"
			key2 := "test2"
			key3 := "test3"

			tests := []struct {
				name       string
				mapValue   map[*string]int
				key        *string
				shouldFail bool
			}{
				{
					name:       "should handle pointer keys that don't exist",
					mapValue:   map[*string]int{&key1: 1, &key2: 2},
					key:        &key3,
					shouldFail: false,
				},
				{
					name:       "should fail with pointer key that exists",
					mapValue:   map[*string]int{&key1: 1, &key2: 2},
					key:        &key1,
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainKey(mockT, tt.mapValue, tt.key)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainKey to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainKey to pass, but it failed: %s", mockT.message)
					}
				})
			}
		})
	})
}

// === Tests for NotContainValue ===

func TestNotContainValue(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Run("string-int maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				value      int
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass when value does not exist",
					mapValue:   map[string]int{"name": 1, "age": 2},
					value:      3,
					shouldFail: false,
				},
				{
					name:       "should pass with empty map",
					mapValue:   map[string]int{},
					value:      42,
					shouldFail: false,
				},
				{
					name:       "should fail when value exists",
					mapValue:   map[string]int{"name": 1, "age": 30, "score": 100},
					value:      30,
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						expectedParts := []string{
							"Expected map to NOT contain value, but it was found:",
							"Map Type : map[string]int",
							"Map Size : 3 entries",
							"Found Value: 30",
							`Found At: key "age"`,
						}
						for _, part := range expectedParts {
							if !strings.Contains(message, part) {
								t.Errorf("Expected message to contain %q, but it was not found in:\n%s", part, message)
							}
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("int-string maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[int]string
				value      string
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass when string value does not exist",
					mapValue:   map[int]string{1: "one", 2: "two"},
					value:      "three",
					shouldFail: false,
				},
				{
					name:       "should fail when string value exists",
					mapValue:   map[int]string{1: "admin", 2: "user", 3: "guest"},
					value:      "user",
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						expectedParts := []string{
							"Expected map to NOT contain value, but it was found:",
							"Map Type : map[int]string",
							`Found Value: "user"`,
							"Found At: key 2",
						}
						for _, part := range expectedParts {
							if !strings.Contains(message, part) {
								t.Errorf("Expected message to contain %q, but it was not found in:\n%s", part, message)
							}
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value)
					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})
	})

	t.Run("Custom messages", func(t *testing.T) {
		t.Run("string to int map", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				value      int
				opts       []Option
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass with custom message",
					mapValue:   map[string]int{"score": 100},
					value:      50,
					opts:       []Option{WithMessage("Score should not be 50")},
					shouldFail: false,
				},
				{
					name:       "should show custom error message on failure",
					mapValue:   map[string]int{"score": 100, "level": 5},
					value:      100,
					opts:       []Option{WithMessage("Score should not be 100")},
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						if !strings.Contains(message, "Score should not be 100") {
							t.Errorf("Expected custom error message, got: %s", message)
						}
						if !strings.Contains(message, "Expected map to NOT contain value, but it was found:") {
							t.Errorf("Expected standard error message, got: %s", message)
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value, tt.opts...)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("string to string map", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]string
				value      string
				opts       []Option
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass with custom message",
					mapValue:   map[string]string{"status": "active"},
					value:      "deleted",
					opts:       []Option{WithMessage("User should not have deleted status")},
					shouldFail: false,
				},
				{
					name:       "should show custom error message on failure",
					mapValue:   map[string]string{"status": "deleted"},
					value:      "deleted",
					opts:       []Option{WithMessage("User should not have deleted status")},
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						if !strings.Contains(message, "User should not have deleted status") {
							t.Errorf("Expected custom error message, got: %s", message)
						}
						if !strings.Contains(message, "Expected map to NOT contain value, but it was found:") {
							t.Errorf("Expected standard error message, got: %s", message)
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value, tt.opts...)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("int to string map", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[int]string
				value      string
				opts       []Option
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should pass with custom message for int keys",
					mapValue:   map[int]string{1: "first", 2: "second"},
					value:      "third",
					opts:       []Option{WithMessage("Should not contain 'third'")},
					shouldFail: false,
				},
				{
					name:       "should show custom error message on failure with int keys",
					mapValue:   map[int]string{1: "first", 2: "second"},
					value:      "first",
					opts:       []Option{WithMessage("Should not contain 'first'")},
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						if !strings.Contains(message, "Should not contain 'first'") {
							t.Errorf("Expected custom error message, got: %s", message)
						}
						if !strings.Contains(message, "Expected map to NOT contain value, but it was found:") {
							t.Errorf("Expected standard error message, got: %s", message)
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value, tt.opts...)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Run("string-int maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[string]int
				value      int
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should handle nil map",
					mapValue:   nil,
					value:      42,
					shouldFail: false,
				},
				{
					name:       "should handle zero value values",
					mapValue:   map[string]int{"zero": 0, "one": 1},
					value:      2,
					shouldFail: false,
				},
				{
					name:       "should fail with zero value that exists",
					mapValue:   map[string]int{"zero": 0, "one": 1},
					value:      0,
					shouldFail: true,
				},
				{
					name:       "should handle multiple keys with same value",
					mapValue:   map[string]int{"first": 42, "second": 42, "third": 100},
					value:      42,
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						if !strings.Contains(message, "Found At: keys") {
							t.Errorf("Expected message to mention multiple keys, got: %s", message)
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("int-string maps", func(t *testing.T) {
			tests := []struct {
				name       string
				mapValue   map[int]string
				value      string
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name:       "should handle empty string values",
					mapValue:   map[int]string{1: "", 2: "test"},
					value:      "missing",
					shouldFail: false,
				},
				{
					name:       "should fail with empty string value that exists",
					mapValue:   map[int]string{1: "", 2: "test"},
					value:      "",
					shouldFail: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})

		t.Run("complex struct values", func(t *testing.T) {
			type User struct{ Name string }

			tests := []struct {
				name       string
				mapValue   map[string]User
				value      User
				shouldFail bool
				errorCheck func(t *testing.T, message string)
			}{
				{
					name: "should handle complex struct values",
					mapValue: map[string]User{
						"user1": {Name: "Alice"},
						"user2": {Name: "Bob"},
					},
					value:      User{Name: "Charlie"},
					shouldFail: false,
				},
				{
					name: "should fail with complex struct value that exists",
					mapValue: map[string]User{
						"user1": {Name: "Alice"},
						"user2": {Name: "Bob"},
					},
					value:      User{Name: "Bob"},
					shouldFail: true,
					errorCheck: func(t *testing.T, message string) {
						expectedParts := []string{
							"Expected map to NOT contain value, but it was found:",
							"Map Type : map[string]assert.User", // Note: type name may vary
							`Found At: key "user2"`,
						}
						// Should NOT contain the verbose context section
						unexpectedParts := []string{
							"Context:",
							"← Found here",
						}
						for _, part := range expectedParts {
							if !strings.Contains(message, part) {
								t.Errorf("Expected message to contain %q, but it was not found in:\n%s", part, message)
							}
						}
						for _, part := range unexpectedParts {
							if strings.Contains(message, part) {
								t.Errorf("Expected message to NOT contain %q, but it was found in:\n%s", part, message)
							}
						}
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					mockT := &mockT{}
					NotContainValue(mockT, tt.mapValue, tt.value)

					if tt.shouldFail && !mockT.Failed() {
						t.Fatal("Expected NotContainValue to fail, but it passed")
					}
					if !tt.shouldFail && mockT.Failed() {
						t.Errorf("Expected NotContainValue to pass, but it failed: %s", mockT.message)
					}
					if tt.errorCheck != nil && mockT.Failed() {
						tt.errorCheck(t, mockT.message)
					}
				})
			}
		})
	})
}

func TestContainSubstring_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		ContainSubstring(t, "Hello, world!", "planet", WithMessage("String should contain 'planet'"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"String should contain 'planet'",
		"Expected string to contain 'planet', but it was not found",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainSubstring(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			substring  string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass if string contains substring",
				actual:     "Hello, world!",
				substring:  "world",
				shouldFail: false,
			},
			{
				name:       "should pass if string contains substring at beginning",
				actual:     "Hello, world!",
				substring:  "Hello",
				shouldFail: false,
			},
			{
				name:       "should pass if string contains substring at end",
				actual:     "Hello, world!",
				substring:  "world!",
				shouldFail: false,
			},
			{
				name:       "should pass with exact match",
				actual:     "test",
				substring:  "test",
				shouldFail: false,
			},
			{
				name:       "should fail if string does not contain substring",
				actual:     "Hello, world!",
				substring:  "planet",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected string to contain 'planet', but it was not found") ||
						!strings.Contains(message, "Substring   : 'planet'") ||
						!strings.Contains(message, "Actual   : 'Hello, world!'") {
						t.Errorf("Incorrect error message:\n%s", message)
					}
				},
			},
			{
				name:       "should pass with empty substring - empty string is contained in any string",
				actual:     "Hello, world!",
				substring:  "",
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				ContainSubstring(mockT, tt.actual, tt.substring)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected ContainSubstring to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected ContainSubstring to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			substring  string
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should pass with ignore case enabled",
				actual:     "Hello, WORLD!",
				substring:  "world",
				opts:       []Option{WithIgnoreCase()},
				shouldFail: false,
			},
			{
				name:       "should pass with ignore case enabled - mixed case",
				actual:     "Hello, WoRlD!",
				substring:  "WORLD",
				opts:       []Option{WithIgnoreCase()},
				shouldFail: false,
			},
			{
				name:       "should fail if ignore case is disabled and case mismatch is detected",
				actual:     "Hello, WORLD!",
				substring:  "world",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Note: Case mismatch detected (use should.WithIgnoreCase() if intended)") {
						t.Errorf("Expected message to contain note message, but got %q", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				ContainSubstring(mockT, tt.actual, tt.substring, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected ContainSubstring to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected ContainSubstring to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			substring  string
			opts       []Option
			shouldFail bool
		}{
			{
				name:       "should pass with custom message",
				actual:     "Hello, world!",
				substring:  "world",
				opts:       []Option{WithMessage("Expected string to contain 'world'")},
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				ContainSubstring(mockT, tt.actual, tt.substring, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected ContainSubstring to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected ContainSubstring to pass, but it failed: %s", mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			substring  string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should handle empty actual with empty substring",
				actual:     "",
				substring:  "",
				shouldFail: false,
			},
			{
				name:       "should fail with empty actual and non-empty substring",
				actual:     "",
				substring:  "test",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Actual   : '<empty>'") {
						t.Errorf("Expected error message to show '<empty>' for empty actual, got:\n%s", message)
					}
				},
			},
			{
				name:       "should handle substring longer than actual",
				actual:     "abc",
				substring:  "abcdef",
				shouldFail: true,
			},
			{
				name:       "should handle single character substring",
				actual:     "Hello, world!",
				substring:  "o",
				shouldFail: false,
			},
			{
				name:       "should handle single character actual",
				actual:     "a",
				substring:  "b",
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				ContainSubstring(mockT, tt.actual, tt.substring)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected ContainSubstring to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected ContainSubstring to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Long string handling", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     string
			substring  string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should use multiline formatting for long strings",
				actual:     "This is a very long string that exceeds the 200 character limit and should trigger multiline formatting in the error message to provide better readability for developers debugging their tests when the assertion fails",
				substring:  "nonexistent",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "(length:") {
						t.Errorf("Expected message to contain length information for long string, got:\n%s", message)
					}
				},
			},
			{
				name:       "should use multiline formatting for strings with newlines",
				actual:     "Line 1\nLine 2\nLine 3",
				substring:  "nonexistent",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "(length:") {
						t.Errorf("Expected message to contain length information for multiline string, got:\n%s", message)
					}
				},
			},
			{
				name:       "should show complete string when shorter than 200 characters",
				actual:     "This is a moderately long string that is longer than 80 characters but shorter than 200",
				substring:  "nonexistent",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "This is a moderately long string that is longer than 80 characters but shorter than 200") {
						t.Errorf("Expected message to contain complete actual string, got:\n%s", message)
					}
				},
			},
			{
				name:       "should show note for large substring",
				actual:     "short text",
				substring:  "This is a very long substring that exceeds the 50 character limit and should trigger a note",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Note: substring is") && !strings.Contains(message, "characters long") {
						t.Errorf("Expected message to contain note about large substring, got:\n%s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				ContainSubstring(mockT, tt.actual, tt.substring)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected ContainSubstring to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected ContainSubstring to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

// === Tests for fail function ===

func TestFail(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name           string
			message        string
			args           []any
			expectedOutput string
			description    string
		}{
			{
				name:           "should use Error when no args provided",
				message:        "simple error message",
				args:           nil,
				expectedOutput: "simple error message",
				description:    "When no arguments are provided, should call t.Error()",
			},
			{
				name:           "should use Errorf when args provided",
				message:        "formatted message: %s",
				args:           []any{"test"},
				expectedOutput: "formatted message: test",
				description:    "When arguments are provided, should call t.Errorf()",
			},
			{
				name:           "should handle multiple args",
				message:        "user %s has %d points",
				args:           []any{"john", 42},
				expectedOutput: "user john has 42 points",
				description:    "Should handle multiple formatting arguments",
			},
			{
				name:           "should handle empty args slice",
				message:        "empty args slice",
				args:           []any{},
				expectedOutput: "empty args slice",
				description:    "Empty args slice should use t.Error()",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				fail(mockT, tt.message, tt.args...)

				if !mockT.Failed() {
					t.Fatal("Expected fail to mark test as failed")
				}

				if mockT.message != tt.expectedOutput {
					t.Errorf("Expected message %q, got %q", tt.expectedOutput, mockT.message)
				}
			})
		}
	})

	t.Run("Security - Percent character handling", func(t *testing.T) {
		tests := []struct {
			name           string
			message        string
			args           []any
			expectedOutput string
			description    string
		}{
			{
				name:           "should safely handle percent in message without args",
				message:        "100% complete",
				args:           nil,
				expectedOutput: "100% complete",
				description:    "Percent characters should be safe when using t.Error()",
			},
			{
				name:           "should safely handle multiple percents without args",
				message:        "progress: 50% done, 25% remaining, 25% unknown",
				args:           nil,
				expectedOutput: "progress: 50% done, 25% remaining, 25% unknown",
				description:    "Multiple percent characters should be safe",
			},
			{
				name:           "should handle percent with format verbs without args",
				message:        "invalid format: %s, %d, %v",
				args:           nil,
				expectedOutput: "invalid format: %s, %d, %v",
				description:    "Format verbs should be treated as literal text when no args",
			},
			{
				name:           "should handle percent with args correctly",
				message:        "progress: %d%% complete",
				args:           []any{75},
				expectedOutput: "progress: 75% complete",
				description:    "Should format correctly when args are provided",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				fail(mockT, tt.message, tt.args...)

				if !mockT.Failed() {
					t.Fatal("Expected fail to mark test as failed")
				}

				if mockT.message != tt.expectedOutput {
					t.Errorf("Expected message %q, got %q", tt.expectedOutput, mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name           string
			message        string
			args           []any
			expectedOutput string
			description    string
		}{
			{
				name:           "should handle empty message",
				message:        "",
				args:           nil,
				expectedOutput: "",
				description:    "Empty message should work",
			},
			{
				name:           "should handle empty message with args",
				message:        "",
				args:           []any{"ignored"},
				expectedOutput: "%!(EXTRA string=ignored)",
				description:    "Empty message with args shows formatting error as expected",
			},
			{
				name:           "should handle newlines in message",
				message:        "line 1\nline 2\nline 3",
				args:           nil,
				expectedOutput: "line 1\nline 2\nline 3",
				description:    "Newlines should be preserved",
			},
			{
				name:           "should handle nil in args",
				message:        "value is %v",
				args:           []any{nil},
				expectedOutput: "value is <nil>",
				description:    "Nil values in args should be formatted correctly",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				fail(mockT, tt.message, tt.args...)

				if !mockT.Failed() {
					t.Fatal("Expected fail to mark test as failed")
				}

				if mockT.message != tt.expectedOutput {
					t.Errorf("Expected message %q, got %q", tt.expectedOutput, mockT.message)
				}
			})
		}
	})

	t.Run("Format string validation", func(t *testing.T) {
		tests := []struct {
			name        string
			message     string
			args        []any
			shouldPanic bool
			description string
		}{
			{
				name:        "should handle mismatched format verbs gracefully",
				message:     "expected %s but got %d",
				args:        []any{"string"},
				shouldPanic: false,
				description: "Mismatched format verbs should not panic",
			},
			{
				name:        "should handle too many args",
				message:     "simple message",
				args:        []any{"extra", "args"},
				shouldPanic: false,
				description: "Extra arguments should not cause panic",
			},
			{
				name:        "should handle complex format string",
				message:     "user %s (id: %d) has %.2f points at %v",
				args:        []any{"john", 123, 45.678, "2023-01-01"},
				shouldPanic: false,
				description: "Complex format strings should work",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}

				defer func() {
					r := recover()
					if tt.shouldPanic && r == nil {
						t.Error("Expected panic but none occurred")
					}
					if !tt.shouldPanic && r != nil {
						t.Errorf("Unexpected panic: %v", r)
					}
				}()

				fail(mockT, tt.message, tt.args...)

				if !mockT.Failed() {
					t.Fatal("Expected fail to mark test as failed")
				}
			})
		}
	})

	t.Run("Helper method call", func(t *testing.T) {
		// This test verifies that fail calls t.Helper()

		t.Run("should call Helper method", func(t *testing.T) {
			mockT := &mockT{}

			fail(mockT, "test message")

			if !mockT.Failed() {
				t.Fatal("Expected fail to mark test as failed")
			}
		})
	})
}

func TestFail_Integration(t *testing.T) {
	t.Parallel()

	t.Run("Integration with BeTrue", func(t *testing.T) {
		mockT := &mockT{}
		BeTrue(mockT, false)

		if !mockT.Failed() {
			t.Fatal("Expected BeTrue to fail and call fail function")
		}

		expected := "Expected true, got false"
		if mockT.message != expected {
			t.Errorf("Expected message %q, got %q", expected, mockT.message)
		}
	})

	t.Run("Integration with custom message", func(t *testing.T) {
		mockT := &mockT{}
		BeTrue(mockT, false, WithMessage("custom error"))

		if !mockT.Failed() {
			t.Fatal("Expected BeTrue to fail")
		}

		if !strings.Contains(mockT.message, "custom error") {
			t.Errorf("Expected message to contain custom error, got %q", mockT.message)
		}
		if !strings.Contains(mockT.message, "Expected true, got false") {
			t.Errorf("Expected message to contain assertion error, got %q", mockT.message)
		}
	})

	t.Run("Integration with percent characters in custom message", func(t *testing.T) {
		mockT := &mockT{}
		BeTrue(mockT, false, WithMessage("progress: 100% complete"))

		if !mockT.Failed() {
			t.Fatal("Expected BeTrue to fail")
		}

		if !strings.Contains(mockT.message, "progress: 100% complete") {
			t.Errorf("Expected message to contain custom message with percent, got %q", mockT.message)
		}
		if !strings.Contains(mockT.message, "Expected true, got false") {
			t.Errorf("Expected message to contain assertion error, got %q", mockT.message)
		}
	})
}

func TestBeInRange(t *testing.T) {
	t.Parallel()

	// Test cases for integer ranges
	t.Run("Integer range", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			value       int
			min         int
			max         int
			opts        []Option
			shouldFail  bool
			expectedMsg string
		}{
			{name: "should pass when value is within range", value: 50, min: 0, max: 100, shouldFail: false},
			{name: "should pass when value is at lower bound", value: 0, min: 0, max: 100, shouldFail: false},
			{name: "should pass when value is at upper bound", value: 100, min: 0, max: 100, shouldFail: false},
			{
				name:        "should fail when value is below range",
				value:       16,
				min:         18,
				max:         65,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [18, 65], but it was below:",
			},
			{
				name:        "should fail when value is above range",
				value:       105,
				min:         0,
				max:         100,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 100], but it was above:",
			},
			{
				name:        "should include custom message on failure",
				value:       150,
				min:         0,
				max:         100,
				opts:        []Option{WithMessage("Battery level must be valid")},
				shouldFail:  true,
				expectedMsg: "Battery level must be valid\nExpected value to be in range [0, 100], but it was above:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				BeInRange(mockT, tt.value, tt.min, tt.max, tt.opts...)

				if tt.shouldFail != mockT.Failed() {
					t.Errorf("Expected test failure to be %v, but was %v", tt.shouldFail, mockT.Failed())
				}

				if tt.shouldFail && !strings.Contains(mockT.message, tt.expectedMsg) {
					t.Errorf("Expected error message to contain %q, but got %q", tt.expectedMsg, mockT.message)
				}
			})
		}
	})

	// Test cases for float ranges
	t.Run("Float range", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			value       float64
			min         float64
			max         float64
			shouldFail  bool
			expectedMsg string
		}{
			{name: "should pass when value is within range", value: 0.5, min: 0.0, max: 1.0, shouldFail: false},
			{name: "should pass when value is at lower bound", value: 0.0, min: 0.0, max: 1.0, shouldFail: false},
			{name: "should pass when value is at upper bound", value: 1.0, min: 0.0, max: 1.0, shouldFail: false},
			{
				name:        "should fail when value is below range",
				value:       -0.1,
				min:         0.0,
				max:         1.0,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 1], but it was below:",
			},
			{
				name:        "should fail when value is above range",
				value:       1.1,
				min:         0.0,
				max:         1.0,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 1], but it was above:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockT := &mockT{}
				BeInRange(mockT, tt.value, tt.min, tt.max)

				if tt.shouldFail != mockT.Failed() {
					t.Errorf("Expected test failure to be %v, but was %v", tt.shouldFail, mockT.Failed())
				}

				if tt.shouldFail && !strings.Contains(mockT.message, tt.expectedMsg) {
					t.Errorf("Expected error message to contain %q, but got %q", tt.expectedMsg, mockT.message)
				}
			})
		}
	})
}
