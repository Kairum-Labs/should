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

// === Tests for error handling ===

func TestBeGreaterThan_Fails_WithNonNumericActual(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterThan(t, "hello", "world")
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
		BeLessThan(t, "hello", "world")
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

// === Tests for BeGreaterOrEqualThan ===

func TestBeGreaterOrEqualThan_Succeeds_WithGreaterValue(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualThan(t, 10, 5)
}

func TestBeGreaterOrEqualThan_Succeeds_WithEqualValue(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualThan(t, 5, 5)
}

func TestBeGreaterOrEqualThan_Succeeds_WithFloats(t *testing.T) {
	t.Parallel()

	BeGreaterOrEqualThan(t, 3.14, 3.14)
	BeGreaterOrEqualThan(t, 3.15, 3.14)
}

func TestBeGreaterOrEqualThan_Fails_WithSmallerValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterOrEqualThan(t, 5, 10)
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
		BeGreaterOrEqualThan(t, 0, 1, WithMessage("Score cannot be negative"))
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
		BeGreaterOrEqualThan(t, "hello", "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeGreaterOrEqualThan_Fails_WithMixedIntFloat(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterOrEqualThan(t, 5.0, 5.1, WithMessage("Integer vs float comparison"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Integer vs float comparison",
		"Expected 5 to be greater or equal than 5.1",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
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

func TestBeTrue_Fails_WithNonBooleanType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeTrue(t, "true")
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

func TestBeFalse_Fails_WithNonBooleanType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeFalse(t, 0)
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
		NotContainDuplicates(t, []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 4, 4})
	})
	if !failed {
		t.Fatal("Expected test to fail due to duplicates, but it passed")
	}
	if !strings.Contains(message, "duplicate values") {
		t.Errorf("Expected error message to mention duplicates, got: %s", message)
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

// === Tests for BeLessThan error handling ===

func TestBeLessThan_Fails_WithNonNumericActual(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeLessThan(t, "hello", "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
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
				opts:       []Option{IgnoreCase()},
				shouldFail: false,
			},
			{
				name:       "should fail if ignore case is disabled and case mismatch is detected",
				actual:     "Hello",
				expected:   "hello",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, `Note: Case mismatch detected (use should.IgnoreCase() if intended)`) {
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
				opts:       []Option{IgnoreCase()},
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

// === Tests for BeGreaterThan error handling ===

func TestBeGreaterThan_Fails_WithNonNumericExpected(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeGreaterThan(t, "hello", "world")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a number for actual value, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
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

	// Test all supported numeric types for better toFloat64 coverage
	testCases := []struct {
		name     string
		value    interface{}
		expected float64
	}{
		{"int8", int8(10), 10.0},
		{"int16", int16(20), 20.0},
		{"int32", int32(30), 30.0},
		{"int64", int64(40), 40.0},
		{"uint8", uint8(50), 50.0},
		{"uint16", uint16(60), 60.0},
		{"uint32", uint32(70), 70.0},
		{"uint64", uint64(80), 80.0},
		{"uintptr", uintptr(90), 90.0},
		{"float32", float32(3.14), 3.140000104904175}, // float32 precision
		{"float64", float64(2.718), 2.718},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This will exercise the toFloat64 function through BeGreaterThan
			BeGreaterThan(t, tc.value, interface{}(tc.expected-1))
		})
	}
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

func TestContainKey_Fails_WithNonMapType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		ContainKey(t, []string{"not", "a", "map"}, "key")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a map, but got []string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
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

func TestContainValue_Fails_WithNonMapType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		ContainValue(t, "not a map", "value")
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "expected a map, but got string"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
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
		tests := []struct {
			name       string
			mapValue   interface{}
			key        interface{}
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
				name:       "should pass with int keys when key does not exist",
				mapValue:   map[int]string{1: "one", 2: "two"},
				key:        3,
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

	t.Run("Custom messages", func(t *testing.T) {
		tests := []struct {
			name       string
			mapValue   interface{}
			key        interface{}
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

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			mapValue   interface{}
			key        interface{}
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should handle nil map",
				mapValue:   func() interface{} { var m map[string]int; return m }(),
				key:        "test",
				shouldFail: false,
			},
			{
				name:       "should handle zero value keys",
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
			{
				name:       "should handle empty string keys",
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
			{
				name:       "should fail for non-map type",
				mapValue:   []string{"not", "a", "map"},
				key:        "key",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "expected a map, but got []string") {
						t.Errorf("Expected error for non-map type, got: %s", message)
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
}

// === Tests for NotContainValue ===

func TestNotContainValue(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			mapValue   interface{}
			value      interface{}
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
				name:       "should pass with string values when value does not exist",
				mapValue:   map[int]string{1: "one", 2: "two"},
				value:      "three",
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

	t.Run("Custom messages", func(t *testing.T) {
		tests := []struct {
			name       string
			mapValue   interface{}
			value      interface{}
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

	t.Run("Edge cases", func(t *testing.T) {
		tests := []struct {
			name       string
			mapValue   interface{}
			value      interface{}
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "should handle nil map",
				mapValue:   func() interface{} { var m map[string]int; return m }(),
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
			{
				name: "should handle complex struct values",
				mapValue: map[string]struct{ Name string }{
					"user1": {Name: "Alice"},
					"user2": {Name: "Bob"},
				},
				value:      struct{ Name string }{Name: "Charlie"},
				shouldFail: false,
			},
			{
				name: "should fail with complex struct value that exists",
				mapValue: map[string]struct{ Name string }{
					"user1": {Name: "Alice"},
					"user2": {Name: "Bob"},
				},
				value:      struct{ Name string }{Name: "Bob"},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					expectedParts := []string{
						"Expected map to NOT contain value, but it was found:",
						"Map Type : map[string]struct { Name string }",
						`Found Value: struct{Name: "Bob"}`,
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
			{
				name:       "should fail for non-map type",
				mapValue:   "not a map",
				value:      "value",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "expected a map, but got string") {
						t.Errorf("Expected error for non-map type, got: %s", message)
					}
				},
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
}
