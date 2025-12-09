package assert

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

// customError and other error types are used for testing BeErrorAs and BeErrorIs.
// They provide distinct types and values to verify the assertion logic.
type customError struct{ msg string }

func (e customError) Error() string { return e.msg }

type anotherError struct{ msg string }

func (e anotherError) Error() string { return e.msg }

type testError struct{ msg string }

func (e testError) Error() string { return e.msg }

type edgeError struct{ msg string }

func (e edgeError) Error() string { return e.msg }

// === Tests for failWithOptions ===

func TestFailWithOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		config           *Config
		format           string
		args             []any
		expectedExact    string   // for exact match
		expectedContains []string // for partial matches
		verifyOrder      bool     // verify custom message appears before default message
	}{
		{
			name:          "without custom message",
			config:        &Config{},
			format:        "Expected condition failed",
			args:          nil,
			expectedExact: "Expected condition failed",
		},
		{
			name:   "with custom message",
			config: &Config{Message: "Custom error message"},
			format: "Expected condition failed",
			args:   nil,
			expectedContains: []string{
				"Custom error message",
				"Expected condition failed",
			},
			verifyOrder: true,
		},
		{
			name:          "with format args",
			config:        &Config{},
			format:        "Expected %d to be greater than %d",
			args:          []any{5, 10},
			expectedExact: "Expected 5 to be greater than 10",
		},
		{
			name:   "with custom message and format args",
			config: &Config{Message: "Age validation failed"},
			format: "Expected age %d to be at least %d",
			args:   []any{16, 18},
			expectedContains: []string{
				"Age validation failed",
				"Expected age 16 to be at least 18",
			},
		},
		{
			name:          "with nil config",
			config:        nil,
			format:        "Expected condition failed",
			args:          nil,
			expectedExact: "Expected condition failed",
		},
		{
			name:          "with empty custom message",
			config:        &Config{Message: ""},
			format:        "Expected condition failed",
			args:          nil,
			expectedExact: "Expected condition failed",
		},
		{
			name:   "with multiline messages",
			config: &Config{Message: "Validation failed:\n- Age is too low\n- Missing required field"},
			format: "Expected user to be valid\nDetails: Invalid age",
			args:   nil,
			expectedContains: []string{
				"Validation failed:",
				"- Age is too low",
				"- Missing required field",
				"Expected user to be valid",
				"Details: Invalid age",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			failed, message := assertFails(t, func(t testing.TB) {
				failWithOptions(t, tt.config, tt.format, tt.args...)
			})

			if !failed {
				t.Fatal("Expected test to fail, but it passed")
			}

			// Check exact match if provided
			if tt.expectedExact != "" {
				if message != tt.expectedExact {
					t.Errorf("Expected message: %q\nGot: %q", tt.expectedExact, message)
				}
			}

			// Check partial matches if provided
			for _, part := range tt.expectedContains {
				if !strings.Contains(message, part) {
					t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
				}
			}

			// Verify order if requested
			if tt.verifyOrder && len(tt.expectedContains) >= 2 {
				customIndex := strings.Index(message, tt.expectedContains[0])
				defaultIndex := strings.Index(message, tt.expectedContains[1])

				if customIndex == -1 || defaultIndex == -1 {
					t.Fatal("Both messages should be present")
				}

				if customIndex > defaultIndex {
					t.Error("Custom message should appear before default message")
				}
			}
		})
	}
}

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

func TestBeEqual_PrimitiveFormatting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		actual        interface{}
		expected      interface{}
		wantFail      bool
		expectedParts []string
	}{
		{
			name:     "different strings",
			actual:   "hello",
			expected: "world",
			wantFail: true,
			expectedParts: []string{
				"Not equal:",
				"expected: world",
				"actual  : hello",
			},
		},
		{
			name:     "different integers",
			actual:   42,
			expected: 99,
			wantFail: true,
			expectedParts: []string{
				"Not equal:",
				"expected: 99",
				"actual  : 42",
			},
		},
		{
			name:     "different booleans",
			actual:   true,
			expected: false,
			wantFail: true,
			expectedParts: []string{
				"Not equal:",
				"expected: false",
				"actual  : true",
			},
		},
		{
			name:     "floats with different values",
			actual:   3.14,
			expected: 6.28,
			wantFail: true,
			expectedParts: []string{
				"Not equal:",
				"expected: 6.28",
				"actual  : 3.14",
			},
		},
		{
			name:     "equal primitive values should pass",
			actual:   "same",
			expected: "same",
			wantFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockT := &mockT{}
			BeEqual(mockT, tt.actual, tt.expected)

			if tt.wantFail && !mockT.Failed() {
				t.Fatal("expected BeEqual to fail, but it passed")
			}
			if !tt.wantFail && mockT.Failed() {
				t.Fatalf("expected BeEqual to pass, but it failed with message: %s", mockT.message)
			}

			if tt.wantFail {
				for _, part := range tt.expectedParts {
					if !strings.Contains(mockT.message, part) {
						t.Errorf("expected error message to contain %q, got:\n%s", part, mockT.message)
					}
				}
			}
		})
	}
}

func TestBeEqual_TypeDifferences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		actual          interface{}
		expected        interface{}
		wantTypeInfo    bool
		expectedContent []string
	}{
		{
			name:         "int32 vs int64 with same value",
			actual:       int32(42),
			expected:     int64(42),
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: 42",
				"actual  : 42",
				"Field differences:",
				"└─ : int64 ≠ int32",
			},
		},
		{
			name:         "float32 vs float64 with same value",
			actual:       float32(3.14),
			expected:     float64(3.14),
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: 3.14",
				"actual  : 3.14",
				"Field differences:",
				"└─ : float64 ≠ float32",
			},
		},
		{
			name:         "int vs int32 with same value",
			actual:       int(100),
			expected:     int32(100),
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: 100",
				"actual  : 100",
				"Field differences:",
				"└─ : int32 ≠ int",
			},
		},
		{
			name:         "uint vs int with same value",
			actual:       uint(50),
			expected:     int(50),
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: 50",
				"actual  : 50",
				"Field differences:",
				"└─ : int ≠ uint",
			},
		},
		{
			name:         "same type and value should pass",
			actual:       int64(42),
			expected:     int64(42),
			wantTypeInfo: false,
		},
		{
			name:         "uintptr vs uint64 with same value",
			actual:       uintptr(0x1000),
			expected:     uint64(0x1000),
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: 4096",
				"actual  : 4096",
				"Field differences:",
				"└─ : uint64 ≠ uintptr",
			},
		},
		{
			name:         "rune vs int32 (should be same type)",
			actual:       rune('A'),
			expected:     int32(65),
			wantTypeInfo: false,
		},
		{
			name:         "int pointer vs int value",
			actual:       func() *int { i := 42; return &i }(),
			expected:     42,
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"Field differences:",
				"└─ : int ≠ ptr",
			},
		},
		{
			name: "struct with different slice lengths",
			actual: struct {
				Name   string
				Ages   []int
				Scores []float64
			}{
				Name:   "João",
				Ages:   []int{10, 20, 30},
				Scores: []float64{8.5, 9},
			},
			expected: struct {
				Name   string
				Ages   []int
				Scores []float64
			}{
				Name:   "João",
				Ages:   []int{10, 20},
				Scores: []float64{8.5, 9},
			},
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				`expected: {Name: "João", Ages: [10, 20], Scores: [8.5, 9]}`,
				`actual  : {Name: "João", Ages: [10, 20, 30], Scores: [8.5, 9]}`,
				"└─ Ages: length mismatch (expected: 2, actual: 3)",
			},
		},
		{
			name: "struct with slice length difference and value difference",
			actual: struct {
				Name   string
				Ages   []int
				Scores []float64
			}{
				Name:   "João",
				Ages:   []int{10, 20, 30},
				Scores: []float64{8.5, 9.5, 7.0},
			},
			expected: struct {
				Name   string
				Ages   []int
				Scores []float64
			}{
				Name:   "Maria",
				Ages:   []int{10, 20},
				Scores: []float64{8.5, 9.0, 7.0},
			},
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				`expected: {Name: "Maria", Ages: [10, 20], Scores: [8.5, 9, 7]}`,
				`actual  : {Name: "João", Ages: [10, 20, 30], Scores: [8.5, 9.5, 7]}`,
				"└─ Name: \"Maria\" ≠ \"João\"",
				"└─ Ages: length mismatch (expected: 2, actual: 3)",
				"└─ Scores.[1]: 9 ≠ 9.5",
			},
		},
		{
			name:         "slice length mismatch",
			actual:       []int{1, 2, 3},
			expected:     []int{1, 2},
			wantTypeInfo: true,
			expectedContent: []string{
				"Not equal:",
				"expected: [1, 2]",
				"actual  : [1, 2, 3]",
				"└─ : length mismatch (expected: 2, actual: 3)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockT := &mockT{}
			BeEqual(mockT, tt.actual, tt.expected)

			if tt.wantTypeInfo {
				if !mockT.Failed() {
					t.Fatal("expected BeEqual to fail due to type difference, but it passed")
				}

				for _, expectedPart := range tt.expectedContent {
					if !strings.Contains(mockT.message, expectedPart) {
						t.Errorf("expected error message to contain %q, got:\n%s", expectedPart, mockT.message)
					}
				}
			} else {
				if mockT.Failed() {
					t.Fatalf("expected BeEqual to pass with same type and value, but it failed with: %s", mockT.message)
				}
			}
		})
	}
}

// === Tests for NotBeEqual ===

func TestNotBeEqual(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()
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
				name:       "should pass with formatted custom message",
				actual:     42,
				expected:   100,
				opts:       []Option{WithMessagef("Expected value %d, got %d", 100, 42)},
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
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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

	AnyMatch(t, []int{1, 2, 3}, func(item int) bool {
		return item == 2
	})
}

func TestContainFunc_Fails_WhenPredicateDoesNotMatch(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		AnyMatch(t, []int{1, 2, 3}, func(item int) bool {
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
	predicate := func(item int) bool {
		return item%2 == 0 // Looking for even numbers
	}

	AnyMatch(mockT, numbers, predicate, WithMessage(customMessage))

	if !mockT.Failed() {
		t.Fatal("Expected AnyMatch to fail, but it passed")
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
		"└─ user-3 (at index 3) - 1 extra character",
		"└─ user-003 (at index 5) - 3 characters differ",
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
			t.Parallel()
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
			t.Parallel()
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

func TestBeEmpty_Fails_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, ptr)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Expected value to be empty, but it was not",
		"Type",
		"*int",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestBeEmpty_Fails_WithNonNilPointer_CustomMessage(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		BeEmpty(t, ptr, WithMessage("Pointer should be nil"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Pointer should be nil",
		"Expected value to be empty, but it was not",
		"Type",
		"*int",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
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

func TestNotPanic_Extended(t *testing.T) {
	t.Parallel()

	t.Run("With WithStackTrace option", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			testFunc      func()
			opts          []Option
			shouldFail    bool
			expectedParts []string
		}{
			{
				name: "should pass when no panic occurs",
				testFunc: func() {
					result := 1 + 2
					_ = result
				},
				opts:       []Option{WithStackTrace()},
				shouldFail: false,
			},
			{
				name: "should fail with stack trace when manual panic occurs",
				testFunc: func() {
					panic("runtime error")
				},
				opts:       []Option{WithStackTrace()},
				shouldFail: true,
				expectedParts: []string{
					"Expected for the function to not panic, but it panicked with:",
					"runtime error",
					"Stack trace:",
				},
			},
			{
				name: "should fail with stack trace when runtime panic occurs",
				testFunc: func() {
					x := 1
					y := 0
					result := x / y
					_ = result
				},
				opts:       []Option{WithStackTrace()},
				shouldFail: true,
				expectedParts: []string{
					"Expected for the function to not panic, but it panicked with:",
					"integer divide by zero",
					"Stack trace:",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.shouldFail {
					failed, message := assertFails(t, func(t testing.TB) {
						NotPanic(t, tt.testFunc, tt.opts...)
					})

					if !failed {
						t.Fatal("Expected test to fail, but it passed")
					}

					for _, part := range tt.expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
						}
					}
				} else {
					NotPanic(t, tt.testFunc, tt.opts...)
				}
			})
		}
	})

	t.Run("Combined options", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			testFunc      func()
			opts          []Option
			shouldFail    bool
			expectedParts []string
		}{
			{
				name: "should combine WithMessage and WithStackTrace",
				testFunc: func() {
					panic("database error")
				},
				opts: []Option{
					WithMessage("Database operation should not panic"),
					WithStackTrace(),
				},
				shouldFail: true,
				expectedParts: []string{
					"Database operation should not panic",
					"Expected for the function to not panic, but it panicked with:",
					"database error",
					"Stack trace:",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.shouldFail {
					failed, message := assertFails(t, func(t testing.TB) {
						NotPanic(t, tt.testFunc, tt.opts...)
					})

					if !failed {
						t.Fatal("Expected test to fail, but it passed")
					}

					for _, part := range tt.expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
						}
					}
				} else {
					NotPanic(t, tt.testFunc, tt.opts...)
				}
			})
		}
	})
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
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		// Integer tests
		BeLessOrEqualTo(t, 5, 10)
		BeLessOrEqualTo(t, 10, 10)

		// Float tests
		BeLessOrEqualTo(t, 2.71, 3.14)
		BeLessOrEqualTo(t, 3.14, 3.14)

		// Test failures
		t.Run("Fails when actual is greater than expected", func(t *testing.T) {
			t.Parallel()
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
			t.Parallel()
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
		t.Parallel()
		// Success with custom message
		BeLessOrEqualTo(t, 5, 10, WithMessage("Value should be within limit"))

		// Fails with custom error message
		t.Run("Fails with custom error message", func(t *testing.T) {
			t.Parallel()
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
		t.Parallel()
		// Success with zero values
		BeLessOrEqualTo(t, 0, 0)

		// Success with negative numbers
		BeLessOrEqualTo(t, -10, -5)

		// Success with very small floats
		BeLessOrEqualTo(t, 0.0001, 0.0002)

		// Fails with negative comparison
		t.Run("Fails with negative comparison", func(t *testing.T) {
			t.Parallel()
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
		t.Parallel()
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

func TestBeWithin(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			actual     float64
			expected   float64
			tolerance  float64
			shouldFail bool
			opts       []Option
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Pass when within tolerance",
				actual:     3.14159,
				expected:   3.1415,
				tolerance:  0.0001,
				shouldFail: false,
			},
			{
				name:       "Pass when exactly on the edge of tolerance",
				actual:     3.1416,
				expected:   3.1415,
				tolerance:  0.0001,
				shouldFail: false,
			},
			{
				name:       "Fail when outside tolerance",
				actual:     3.142,
				expected:   3.14,
				tolerance:  0.001,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Difference: 0.0020") {
						t.Errorf("Expected a specific error message about the difference, but got: %s", message)
					}
				},
			},
			{
				name:       "Fail with zero tolerance and different values",
				actual:     1.00001,
				expected:   1.0,
				tolerance:  0.0,
				shouldFail: true,
			},
			{
				name:       "Fail with zero tolerance and different values with custom message",
				actual:     1.00001,
				expected:   1.0,
				tolerance:  0.0,
				opts:       []Option{WithMessage("Custom error message")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Custom error message") {
						t.Errorf("Expected custom message, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when tolerance is negative",
				actual:     1.0,
				expected:   1.0,
				tolerance:  -0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Tolerance must be non-negative") {
						t.Errorf("Expected negative tolerance error, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when actual is NaN",
				actual:     math.NaN(),
				expected:   1.0,
				tolerance:  0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "NaN detected") {
						t.Errorf("Expected NaN detection error, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when expected is NaN",
				actual:     1.0,
				expected:   math.NaN(),
				tolerance:  0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "NaN detected") {
						t.Errorf("Expected NaN detection error, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when tolerance is NaN",
				actual:     1.0,
				expected:   1.0,
				tolerance:  math.NaN(),
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "NaN detected") {
						t.Errorf("Expected NaN detection error, got: %s", message)
					}
				},
			},
			{
				name:       "Pass when both are +Inf",
				actual:     math.Inf(1),
				expected:   math.Inf(1),
				tolerance:  0.1,
				shouldFail: false,
			},
			{
				name:       "Pass when both are -Inf",
				actual:     math.Inf(-1),
				expected:   math.Inf(-1),
				tolerance:  0.1,
				shouldFail: false,
			},
			{
				name:       "Fail when one is +Inf and other is -Inf",
				actual:     math.Inf(1),
				expected:   math.Inf(-1),
				tolerance:  0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Inf mismatch") {
						t.Errorf("Expected Inf mismatch error, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when actual is Inf and expected is finite",
				actual:     math.Inf(1),
				expected:   42.0,
				tolerance:  0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Inf mismatch") {
						t.Errorf("Expected Inf mismatch error, got: %s", message)
					}
				},
			},
			{
				name:       "Fail when expected is Inf and actual is finite",
				actual:     42.0,
				expected:   math.Inf(1),
				tolerance:  0.1,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Inf mismatch") {
						t.Errorf("Expected Inf mismatch error, got: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeWithin(mockT, tt.actual, tt.expected, tt.tolerance, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Float32 functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			actual     float32
			expected   float32
			tolerance  float32
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Pass with float32 when within tolerance",
				actual:     3.14159,
				expected:   3.1415,
				tolerance:  0.0001,
				shouldFail: false,
			},
			{
				name:       "Fail with float32 when outside tolerance",
				actual:     3.142,
				expected:   3.14,
				tolerance:  0.001,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Difference: 0.0020") {
						t.Errorf("Expected a specific error message about the difference, but got: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeWithin(mockT, tt.actual, tt.expected, tt.tolerance)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			actual     float64
			expected   float64
			tolerance  float64
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Negative values",
				actual:     -10.5,
				expected:   -10.6,
				tolerance:  0.1,
				shouldFail: false,
			},
			{
				name:       "Zero values",
				actual:     0.0,
				expected:   0.0,
				tolerance:  0.0,
				shouldFail: false,
			},
			{
				name:       "Large values",
				actual:     1000000.0,
				expected:   1000000.1,
				tolerance:  0.2,
				shouldFail: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeWithin(mockT, tt.actual, tt.expected, tt.tolerance)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
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

func TestBeError(t *testing.T) {
	t.Parallel()
	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			err        error
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Error present - should pass",
				err:        errors.New("test error"),
				shouldFail: false,
			},
			{
				name:       "Nil error - should fail",
				err:        nil,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected an error, but got nil") {
						t.Errorf("Expected nil error message, got: %s", message)
					}
				},
			},
			{
				name:       "Nil error with custom message",
				err:        nil,
				opts:       []Option{WithMessage("Custom error message")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Custom error message") {
						t.Errorf("Expected custom message, got: %s", message)
					}
					if !strings.Contains(message, "Expected an error, but got nil") {
						t.Errorf("Expected default message, got: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeError(mockT, tt.err, tt.opts...)
				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

func TestNotBeError(t *testing.T) {
	t.Parallel()
	t.Run("Basic fuctionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			err        error
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, err error, message string)
		}{
			{
				name:       "Nil error - should pass",
				err:        nil,
				shouldFail: false,
			},
			{
				name:       "Error present - should fail",
				err:        errors.New("test error"),
				shouldFail: true,
				errorCheck: func(t *testing.T, err error, message string) {
					contains := []string{
						"Expected no error, but got an error",
						"Error: \"test error\"",
						"Type: *errors.errorString",
					}
					for _, expected := range contains {
						if !strings.Contains(message, expected) {
							t.Errorf("Expected an error message, got %s", message)
						}
					}
				},
			},
			{
				name:       "An error with a custom error message",
				err:        errors.New("test an error with a custom error message"),
				shouldFail: true,
				opts:       []Option{WithMessage("Custom error message")},
				errorCheck: func(t *testing.T, err error, message string) {
					if !strings.Contains(message, "Custom error message") {
						t.Errorf("Expected custom message, got %s", message)
					}
					contains := []string{
						"Custom error message",
						"Expected no error, but got an error",
						"Error: \"test an error with a custom error message\"",
						"Type: *errors.errorString",
					}
					for _, expected := range contains {
						if !strings.Contains(message, expected) {
							t.Errorf("Expected: %s, got: %s", expected, message)
						}
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				NotBeError(mockT, tt.err, tt.opts...)
				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, tt.err, mockT.message)
				}
			})
		}
	})
}

func TestBeErrorAs(t *testing.T) {
	t.Parallel()
	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			err        error
			target     interface{}
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Error matches target type",
				err:        customError{msg: "test error"},
				target:     &customError{},
				shouldFail: false,
			},
			{
				name:       "Error doesn't match target type",
				err:        customError{msg: "test error"},
				target:     &anotherError{},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "customError") || !strings.Contains(message, "anotherError") {
						t.Errorf("Error message should contain both error types: %s", message)
					}
				},
			},
			{
				name:       "Nil error with target type",
				err:        nil,
				target:     &customError{},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected error to be") || !strings.Contains(message, "but got nil") {
						t.Errorf("Error message should indicate nil error: %s", message)
					}
				},
			},
			{
				name:       "Wrapped error matches target type",
				err:        fmt.Errorf("wrapped: %w", customError{msg: "inner error"}),
				target:     &customError{},
				shouldFail: false,
			},
			{
				name:       "Error with custom message",
				err:        anotherError{msg: "test"},
				target:     &customError{},
				opts:       []Option{WithMessage("Custom error message")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Custom error message") {
						t.Errorf("Error message should contain custom message: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeErrorAs(mockT, tt.err, tt.target, tt.opts...)
				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Options handling", func(t *testing.T) {
		t.Parallel()

		t.Run("WithMessage option", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeErrorAs(mockT, nil, &testError{}, WithMessage("Custom message for nil error"))
			if !mockT.failed {
				t.Fatal("Expected test to fail")
			}
			if !strings.Contains(mockT.message, "Custom message for nil error") {
				t.Errorf("Expected custom message in error output: %s", mockT.message)
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		t.Run("Multiple wrapped errors", func(t *testing.T) {
			t.Parallel()
			innerErr := edgeError{msg: "inner"}
			middleErr := fmt.Errorf("middle: %w", innerErr)
			outerErr := fmt.Errorf("outer: %w", middleErr)

			mockT := &mockT{}
			BeErrorAs(mockT, outerErr, &edgeError{})
			if mockT.failed {
				t.Errorf("Expected test to pass with multiple wrapped errors: %s", mockT.message)
			}
		})

		t.Run("Nil target should not panic", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeErrorAs(mockT, errors.New("test"), nil)
			if !mockT.failed {
				t.Error("Expected test to fail when target is nil")
			}
			if !strings.Contains(mockT.message, "target cannot be nil") {
				t.Errorf("Expected specific error message about nil target: %s", mockT.message)
			}
		})
	})
}

func TestBeErrorIs(t *testing.T) {
	t.Parallel()
	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()

		var (
			ErrNotFound = errors.New("not found")
			ErrInvalid  = errors.New("invalid")
		)

		tests := []struct {
			name       string
			err        error
			target     error
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "Error matches target",
				err:        ErrNotFound,
				target:     ErrNotFound,
				shouldFail: false,
			},
			{
				name:       "Error doesn't match target",
				err:        ErrNotFound,
				target:     ErrInvalid,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "not found") || !strings.Contains(message, "invalid") {
						t.Errorf("Error message should contain both error messages: %s", message)
					}
				},
			},
			{
				name:       "Nil error with target",
				err:        nil,
				target:     ErrNotFound,
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected error to be") || !strings.Contains(message, "but got nil") {
						t.Errorf("Error message should indicate nil error: %s", message)
					}
				},
			},
			{
				name:       "Wrapped error matches target",
				err:        fmt.Errorf("wrapped: %w", ErrNotFound),
				target:     ErrNotFound,
				shouldFail: false,
			},
			{
				name:       "Error with custom message",
				err:        ErrInvalid,
				target:     ErrNotFound,
				opts:       []Option{WithMessage("Custom error message")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Custom error message") {
						t.Errorf("Error message should contain custom message: %s", message)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeErrorIs(mockT, tt.err, tt.target, tt.opts...)
				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Options handling", func(t *testing.T) {
		t.Parallel()

		ErrTest := errors.New("test error")

		t.Run("WithMessage option", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeErrorIs(mockT, nil, ErrTest, WithMessage("Custom message for nil error"))
			if !mockT.failed {
				t.Fatal("Expected test to fail")
			}
			if !strings.Contains(mockT.message, "Custom message for nil error") {
				t.Errorf("Expected custom message in error output: %s", mockT.message)
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		ErrBase := errors.New("base error")

		t.Run("Multiple levels of wrapping", func(t *testing.T) {
			t.Parallel()
			level1 := fmt.Errorf("level1: %w", ErrBase)
			level2 := fmt.Errorf("level2: %w", level1)
			level3 := fmt.Errorf("level3: %w", level2)

			mockT := &mockT{}
			BeErrorIs(mockT, level3, ErrBase)
			if mockT.failed {
				t.Errorf("Expected test to pass with multiple wrapped levels: %s", mockT.message)
			}
		})

		t.Run("Nil target error", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeErrorIs(mockT, errors.New("test"), nil)
			if !mockT.failed {
				t.Error("Expected test to fail when target is nil")
			}
		})

		t.Run("Both errors nil", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeErrorIs(mockT, nil, nil)
			if !mockT.failed {
				t.Error("Expected test to fail when both errors are nil")
			}
			if !strings.Contains(mockT.message, "but got nil") {
				t.Errorf("Expected nil error message: %s", mockT.message)
			}
		})

		t.Run("String comparison in different error instances", func(t *testing.T) {
			t.Parallel()
			err1 := errors.New("same message")
			err2 := errors.New("same message")

			mockT := &mockT{}
			BeErrorIs(mockT, err1, err2)
			if !mockT.failed {
				t.Error("Expected test to fail for different error instances with same message")
			}
		})
	})
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

func TestBeTrue_WithMessagef(t *testing.T) {
	t.Parallel()
	failed, message := assertFails(t, func(t testing.TB) {
		BeTrue(t, false, WithMessagef("custom: %d %s", 42, "error"))
	})
	if !failed {
		t.Error("Expected assertion to fail")
	}
	expected := "custom: 42 error"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain %q, got %q", expected, message)
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

// === Tests for AnyMatch edge cases ===

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
	AnyMatch(t, users, func(user User) bool {
		return user.Age >= 18
	})

	// Should fail - no elderly users
	failed, message := assertFails(t, func(t testing.TB) {
		AnyMatch(t, users, func(user User) bool {
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
		NotContainDuplicates(
			t,
			[]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 4, 4},
			WithMessage("Expected no duplicates, but found 1 duplicate value: 2"),
		)
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

// === Tests for StartWith ===

func TestStartsWith_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		StartWith(t, "Hello, world!", "world", WithMessage("String should start with 'world'"))
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
		t.Parallel()
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
				t.Parallel()
				mockT := &mockT{}
				StartWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		t.Parallel()
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
				t.Parallel()
				mockT := &mockT{}
				StartWith(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		t.Parallel()
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
				t.Parallel()
				mockT := &mockT{}
				StartWith(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartWith to pass, but it failed: %s", mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
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
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				StartWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("String truncation", func(t *testing.T) {
		t.Parallel()
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
				StartWith(mockT, tt.actual, tt.expected)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected StartWith to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected StartWith to pass, but it failed: %s", mockT.message)
				}
				if tt.errorCheck != nil && mockT.failed {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})
}

// === Tests for EndWith ===

func TestEndsWith(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
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
				EndWith(mockT, tt.actual, tt.expected)

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
		t.Parallel()
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
				EndWith(mockT, tt.actual, tt.expected, tt.opts...)

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
		t.Parallel()
		t.Run("Fails with custom message", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			EndWith(mockT, "Hello, world!", "planet", WithMessage("String should end with 'planet'"))

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
		t.Parallel()
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
				EndWith(mockT, tt.actual, tt.expected)

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
		t.Parallel()
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
				EndWith(mockT, tt.actual, tt.expected)

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
		t.Parallel()
		BeGreaterThan(t, int8(10), int8(9))
	})

	t.Run("int16", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, int16(20), int16(19))
	})

	t.Run("int32", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, int32(30), int32(29))
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, int64(40), int64(39))
	})

	t.Run("uint8", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, uint8(50), uint8(49))
	})

	t.Run("uint16", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, uint16(60), uint16(59))
	})

	t.Run("uint32", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, uint32(70), uint32(69))
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, uint64(80), uint64(79))
	})

	t.Run("float32", func(t *testing.T) {
		t.Parallel()
		BeGreaterThan(t, float32(3.14), float32(3.13))
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()
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
	t.Parallel()

	t.Run("should pass for correct length of slice", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		HaveLength(mockT, []int{1, 2, 3}, 3)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should pass for correct length of string", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		HaveLength(mockT, "abc", 3)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should pass for correct length of map", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		HaveLength(mockT, map[int]int{1: 1, 2: 2}, 2)
		if mockT.failed {
			t.Errorf("Expected HaveLength to pass, but it failed with message: %q", mockT.message)
		}
	})

	t.Run("should fail for incorrect length with custom message", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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

func TestBeSameTime(t *testing.T) {
	t.Parallel()

	// Helper times for testing
	baseTime := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	sameTimeUTC := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	sameTimeEST := time.Date(2023, 12, 25, 10, 30, 45, 123456789, time.FixedZone("EST", -5*3600))
	differentTime := time.Date(2023, 12, 25, 15, 30, 46, 123456789, time.UTC)
	differentNanos := time.Date(2023, 12, 25, 15, 30, 45, 987654321, time.UTC)

	t.Run("Basic functionality", func(t *testing.T) {
		tests := []struct {
			name       string
			actual     time.Time
			expected   time.Time
			opts       []Option
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name:       "exact same time should pass",
				actual:     baseTime,
				expected:   sameTimeUTC,
				shouldFail: false,
			},
			{
				name:       "different time should fail with custom message",
				actual:     baseTime,
				expected:   differentTime,
				opts:       []Option{WithMessage("Expected times to match but they differ")},
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					if !strings.Contains(message, "Expected times to match but they differ") {
						t.Errorf("Expected error message to contain custom message, got: %s", message)
					}
				},
			},
			{
				name:       "different times should fail",
				actual:     baseTime,
				expected:   differentTime,
				shouldFail: true,
			},
			{
				name:       "same time different timezone should pass without options",
				actual:     baseTime,
				expected:   sameTimeEST,
				shouldFail: false, // Same instant, different representation
			},
			{
				name:       "same time different timezone should pass with IgnoreTimezone",
				actual:     baseTime,
				expected:   sameTimeEST,
				opts:       []Option{WithIgnoreTimezone()},
				shouldFail: false,
			},
			{
				name:       "different nanoseconds should fail without options",
				actual:     baseTime,
				expected:   differentNanos,
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeSameTime(mockT, tt.actual, tt.expected, tt.opts...)

				if tt.shouldFail && !mockT.failed {
					t.Fatal("Expected test to fail, but it passed")
				}
				if !tt.shouldFail && mockT.failed {
					t.Errorf("Expected test to pass, but it failed: %s", mockT.message)
				}

				if tt.shouldFail && tt.errorCheck != nil {
					tt.errorCheck(t, mockT.message)
				}
			})
		}
	})

	t.Run("Options handling", func(t *testing.T) {
		t.Parallel()

		t.Run("Multiple options work together", func(t *testing.T) {
			t.Parallel()
			t1 := time.Date(2023, 6, 15, 14, 30, 25, 123456789, time.UTC)
			t2 := time.Date(2023, 6, 15, 9, 30, 25, 987654321, time.FixedZone("EST", -5*3600))

			mockT := &mockT{}
			BeSameTime(mockT, t1, t2, WithIgnoreTimezone(), WithTruncate(time.Second))

			if mockT.failed {
				t.Errorf("Times should be equal with both options: %s", mockT.message)
			}
		})

		t.Run("WithMessage formats correctly", func(t *testing.T) {
			t.Parallel()
			t1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
			t2 := time.Date(2023, 1, 1, 12, 0, 1, 0, time.UTC)

			mockT := &mockT{}
			BeSameTime(mockT, t1, t2, WithMessage("Time validation failed"))

			if !mockT.failed {
				t.Fatal("Expected test to fail")
			}

			if !strings.Contains(mockT.message, "Time validation failed") {
				t.Errorf("Error message should contain custom message, got: %s", mockT.message)
			}
		})
	})

	t.Run("WithIgnoreTimezone works correctly", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			actual     time.Time
			expected   time.Time
			shouldFail bool
		}{
			{
				name:       "pass when same instant but different positive fixed zones",
				actual:     time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC),
				expected:   time.Date(2023, 1, 1, 18, 30, 0, 0, time.FixedZone("UTC+3", 3*3600)),
				shouldFail: false,
			},
			{
				name:       "pass when same instant but with negative fixed zones",
				actual:     time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC),
				expected:   time.Date(2023, 1, 1, 10, 30, 0, 0, time.FixedZone("UTC-5", -5*3600)),
				shouldFail: false,
			},
			{
				name:       "pass when same instant but with named timezone (e.g., America/Sao_Paulo)",
				actual:     time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC),
				expected:   time.Date(2023, 1, 1, 12, 30, 0, 0, time.FixedZone("America/Sao_Paulo", -3*3600)),
				shouldFail: false,
			},
			{
				name:       "fail when instants are different with different timezones",
				actual:     time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC),
				expected:   time.Date(2023, 1, 1, 15, 30, 1, 0, time.FixedZone("UTC+3", 3*3600)),
				shouldFail: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				BeSameTime(mockT, tt.actual, tt.expected, WithIgnoreTimezone())

				if mockT.failed != tt.shouldFail {
					t.Errorf("Test failed for scenario '%s'. Expected failure: %t, but got failure: %t. Error: %s",
						tt.name, tt.shouldFail, mockT.failed, mockT.message)
				}
			})
		}
	})

	t.Run("WithTruncate works correctly", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name         string
			actual       time.Time
			expected     time.Time
			truncateUnit time.Duration
			shouldFail   bool
		}{
			{
				name:         "pass when times are truncated to seconds",
				actual:       time.Date(2023, 1, 1, 12, 0, 0, 1, time.UTC),
				expected:     time.Date(2023, 1, 1, 12, 0, 0, 999999999, time.UTC),
				truncateUnit: time.Second,
				shouldFail:   false,
			},
			{
				name:         "fail when times are different even after truncating to seconds",
				actual:       time.Date(2023, 1, 1, 12, 0, 1, 0, time.UTC),
				expected:     time.Date(2023, 1, 1, 12, 0, 2, 0, time.UTC),
				truncateUnit: time.Second,
				shouldFail:   true,
			},
			{
				name:         "pass when times are truncated to minutes",
				actual:       time.Date(2023, 1, 1, 12, 1, 10, 0, time.UTC),
				expected:     time.Date(2023, 1, 1, 12, 1, 50, 0, time.UTC),
				truncateUnit: time.Minute,
				shouldFail:   false,
			},
			{
				name:         "fail when times are different even after truncating to minutes",
				actual:       time.Date(2023, 1, 1, 12, 1, 0, 0, time.UTC),
				expected:     time.Date(2023, 1, 1, 12, 2, 0, 0, time.UTC),
				truncateUnit: time.Minute,
				shouldFail:   true,
			},
			{
				name:         "edge case: truncating just below a minute boundary",
				actual:       time.Date(2023, 1, 1, 12, 0, 59, 999999999, time.UTC),
				expected:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				truncateUnit: time.Minute,
				shouldFail:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockT := &mockT{}
				BeSameTime(mockT, tt.actual, tt.expected, WithTruncate(tt.truncateUnit))

				if mockT.failed != tt.shouldFail {
					t.Errorf("Test failed for scenario '%s'. Expected failure: %t, but got failure: %t. Error: %s",
						tt.name, tt.shouldFail, mockT.failed, mockT.message)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		t.Run("zero times", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeSameTime(mockT, time.Time{}, time.Time{})

			if mockT.failed {
				t.Errorf("Zero times should be equal: %s", mockT.message)
			}
		})

		t.Run("zero vs non-zero time", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			BeSameTime(mockT, time.Time{}, time.Now())

			if !mockT.failed {
				t.Error("Zero time should not equal non-zero time")
			}
		})

		t.Run("daylight saving time transitions", func(t *testing.T) {
			t.Parallel()
			// Test DST transitions - times that are NOT the same instant
			loc, err := time.LoadLocation("America/New_York")
			if err != nil {
				t.Skip("Could not load timezone data")
			}

			// These are different actual times
			beforeDST := time.Date(2023, 3, 12, 1, 30, 0, 0, loc)
			afterDST := time.Date(2023, 3, 12, 3, 30, 0, 0, loc)

			mockT := &mockT{}
			BeSameTime(mockT, beforeDST, afterDST)

			if !mockT.failed {
				t.Error("Different DST times should not be equal")
			}
		})

		t.Run("maximum time difference", func(t *testing.T) {
			t.Parallel()
			t1 := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			t2 := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

			mockT := &mockT{}
			BeSameTime(mockT, t1, t2)

			if !mockT.failed {
				t.Error("Maximum time difference should fail")
			}

			// Just verify we got an error message
			if len(mockT.message) == 0 {
				t.Error("Error message should not be empty for large time differences")
			}
		})

		t.Run("nanosecond precision edge cases", func(t *testing.T) {
			t.Parallel()
			// Times that differ by exactly 1 nanosecond
			t1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
			t2 := time.Date(2023, 1, 1, 12, 0, 0, 1, time.UTC)

			// Without ignoring nanoseconds - should fail
			mockT1 := &mockT{}
			BeSameTime(mockT1, t1, t2)
			if !mockT1.failed {
				t.Error("Times differing by 1ns should fail without IgnoreNanoseconds")
			}

			// With ignoring nanoseconds - should pass
			mockT2 := &mockT{}
			BeSameTime(mockT2, t1, t2, WithTruncate(time.Second))
			if mockT2.failed {
				t.Errorf("Times should be equal when ignoring nanoseconds: %s", mockT2.message)
			}
		})

		t.Run("timezone offset edge cases", func(t *testing.T) {
			t.Parallel()
			// Test extreme timezone offsets - SAME INSTANT
			utc := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
			// Same instant in UTC+14
			plus14 := utc.In(time.FixedZone("UTC+14", 14*3600))
			// Same instant in UTC-12
			minus12 := utc.In(time.FixedZone("UTC-12", -12*3600))

			// These should pass even without IgnoreTimezone because they're the same instant
			mockT1 := &mockT{}
			BeSameTime(mockT1, utc, plus14)
			if mockT1.failed {
				t.Errorf("Same instant should pass regardless of timezone representation: %s", mockT1.message)
			}

			mockT2 := &mockT{}
			BeSameTime(mockT2, utc, minus12)
			if mockT2.failed {
				t.Errorf("Same instant should pass regardless of timezone representation: %s", mockT2.message)
			}
		})

		t.Run("different calendar dates same instant", func(t *testing.T) {
			t.Parallel()
			// Test case where calendar dates differ but it's the same instant
			utc := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			// Same instant, but previous day in a western timezone
			west := time.Date(2022, 12, 31, 14, 0, 0, 0, time.FixedZone("UTC-10", -10*3600))

			mockT := &mockT{}
			BeSameTime(mockT, utc, west)

			if mockT.failed {
				t.Errorf("Same instant should pass even with different calendar dates: %s", mockT.message)
			}
		})
	})
}

// === Tests for BeOfType ===

func TestBeOfType(t *testing.T) {
	t.Parallel()
	type Cat struct{ Name string }
	type Dog struct{ Name string }

	t.Run("should pass for same type", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		var c *Cat
		BeOfType(mockT, &Cat{}, c)
		if mockT.failed {
			t.Errorf("Expected BeOfType to pass, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail for different types", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		var d *Dog
		BeOfType(mockT, &Cat{Name: "Whiskers"}, d)
		if !mockT.failed {
			t.Fatal("Expected BeOfType to fail, but it passed")
		}

		if !strings.Contains(mockT.message, "Expected value to be of specific type:") ||
			!strings.Contains(mockT.message, "Expected Type: *assert.Dog") ||
			!strings.Contains(mockT.message, "Actual Type  : *assert.Cat") ||
			!strings.Contains(mockT.message, "Difference   : Different concrete types") {
			t.Errorf("Error message format is incorrect.\nGot:\n%s", mockT.message)
		}
	})

	t.Run("should pass for primitive types", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		BeOfType(mockT, 1, 0) // int and int
		if mockT.failed {
			t.Errorf("Expected BeOfType to pass for ints, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail for different primitive types", func(t *testing.T) {
		t.Parallel()
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
	t.Parallel()

	t.Run("should pass if value is one of the options", func(t *testing.T) {
		t.Parallel()
		mockT := &mockT{}
		options := []string{"active", "inactive"}
		BeOneOf(mockT, "active", options)
		if mockT.failed {
			t.Errorf("Expected BeOneOf to pass, but it failed: %s", mockT.message)
		}
	})

	t.Run("should fail if value is not one of the options", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		t.Run("string keys", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		t.Parallel()
		t.Run("string-int maps", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		t.Parallel()
		t.Run("nil map handling", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		t.Parallel()
		t.Run("string-int maps", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		t.Parallel()
		t.Run("string to int map", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		t.Parallel()
		t.Run("string-int maps", func(t *testing.T) {
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
			t.Parallel()
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
					t.Parallel()
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
		ContainSubstring(t, "Hello, world!", "planet", WithMessage(`String should contain "planet"`))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		`String should contain "planet"`,
		`Expected string to contain "planet", but it was not found`,
	}

	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", part, message)
		}
	}
}

func TestContainSubstring_CaseMismatchWithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		ContainSubstring(t, "Hello, WORLD!", "world", WithMessage("Custom case mismatch error"))
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expectedParts := []string{
		"Custom case mismatch error",
		`Expected string to contain "world", but found case difference`,
		`Substring: "world"`,
		`Found    : "WORLD" at position 7`,
		"Note: Case mismatch detected (use should.WithIgnoreCase() if intended)",
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
		t.Parallel()
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
					if !strings.Contains(message, `Expected string to contain "planet", but it was not found`) ||
						!strings.Contains(message, `Substring   : "planet"`) ||
						!strings.Contains(message, `Actual   : "Hello, world!"`) {
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
			{
				name:       "should use fall-back error with case mismatch note when exact pattern not found",
				actual:     "Connection failed: TİMEOUT", // Turkish İ creates case folding edge case
				substring:  "timeout",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					expectedParts := []string{
						`Expected string to contain "timeout", but it was not found`,
						"Note: Case mismatch detected (use should.WithIgnoreCase() if intended)",
					}
					for _, part := range expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain %q, but got:\n%s", part, message)
						}
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
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
		t.Parallel()
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
					expectedParts := []string{
						`Expected string to contain "world", but found case difference`,
						`Substring: "world"`,
						`Found    : "WORLD" at position 7`,
						`Note: Case mismatch detected (use should.WithIgnoreCase() if intended)`,
					}
					for _, part := range expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain %q, but got:\n%s", part, message)
						}
					}
				},
			},
			{
				name:       "should show simplified case mismatch error for mixed case",
				actual:     `Get "http://127.0.0.1:56748": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`,
				substring:  "timeout",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					expectedParts := []string{
						`Expected string to contain "timeout", but found case difference`,
						`Substring: "timeout"`,
						`Found    : "Timeout" at position 64`,
						"Note: Case mismatch detected (use should.WithIgnoreCase() if intended)",
					}
					for _, part := range expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain %q, but got:\n%s", part, message)
						}
					}
				},
			},
			{
				name:       "should show simplified case mismatch error at beginning",
				actual:     "HELLO world",
				substring:  "hello",
				shouldFail: true,
				errorCheck: func(t *testing.T, message string) {
					expectedParts := []string{
						`Expected string to contain "hello", but found case difference`,
						`Substring: "hello"`,
						`Found    : "HELLO" at position 0`,
						"Note: Case mismatch detected (use should.WithIgnoreCase() if intended)",
					}
					for _, part := range expectedParts {
						if !strings.Contains(message, part) {
							t.Errorf("Expected message to contain %q, but got:\n%s", part, message)
						}
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()
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
					if !strings.Contains(message, `Actual   : "<empty>"`) {
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
				t.Parallel()
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
		t.Parallel()
		tests := []struct {
			name       string
			actual     string
			substring  string
			shouldFail bool
			errorCheck func(t *testing.T, message string)
		}{
			{
				name: "should use multiline formatting for long strings",
				actual: "This is a very long string that exceeds the 200 character limit and should trigger " +
					"multiline formatting in the error message to provide better readability for developers " +
					"debugging their tests when the assertion fails",
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
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()
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
				t.Parallel()
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
		t.Parallel()

		t.Run("should call Helper method", func(t *testing.T) {
			t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
			minValue    int
			maxValue    int
			opts        []Option
			shouldFail  bool
			expectedMsg string
		}{
			{name: "should pass when value is within range", value: 50, minValue: 0, maxValue: 100, shouldFail: false},
			{name: "should pass when value is at lower bound", value: 0, minValue: 0, maxValue: 100, shouldFail: false},
			{name: "should pass when value is at upper bound", value: 100, minValue: 0, maxValue: 100, shouldFail: false},
			{
				name:        "should fail when value is below range",
				value:       16,
				minValue:    18,
				maxValue:    65,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [18, 65], but it was below:",
			},
			{
				name:        "should fail when value is above range",
				value:       105,
				minValue:    0,
				maxValue:    100,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 100], but it was above:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeInRange(mockT, tt.value, tt.minValue, tt.maxValue, tt.opts...)

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
			minValue    float64
			maxValue    float64
			shouldFail  bool
			expectedMsg string
		}{
			{name: "should pass when value is within range", value: 0.5, minValue: 0.0, maxValue: 1.0, shouldFail: false},
			{name: "should pass when value is at lower bound", value: 0.0, minValue: 0.0, maxValue: 1.0, shouldFail: false},
			{name: "should pass when value is at upper bound", value: 1.0, minValue: 0.0, maxValue: 1.0, shouldFail: false},
			{
				name:        "should fail when value is below range",
				value:       -0.1,
				minValue:    0.0,
				maxValue:    1.0,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 1], but it was below:",
			},
			{
				name:        "should fail when value is above range",
				value:       1.1,
				minValue:    0.0,
				maxValue:    1.0,
				shouldFail:  true,
				expectedMsg: "Expected value to be in range [0, 1], but it was above:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				mockT := &mockT{}
				BeInRange(mockT, tt.value, tt.minValue, tt.maxValue)

				if tt.shouldFail != mockT.Failed() {
					t.Errorf("Expected test failure to be %v, but was %v", tt.shouldFail, mockT.Failed())
				}

				if tt.shouldFail && !strings.Contains(mockT.message, tt.expectedMsg) {
					t.Errorf("Expected error message to contain %q, but got %q", tt.expectedMsg, mockT.message)
				}
			})
		}
	})

	t.Run("Custom messages", func(t *testing.T) {
		t.Parallel()
		t.Run("should include custom message on failure when below range", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			opts := []Option{WithMessage("Value is out of bounds")}
			BeInRange(mockT, 10, 20, 30, opts...)

			if !mockT.Failed() {
				t.Fatal("Expected test to fail, but it passed")
			}

			expectedMsg := "Value is out of bounds\nExpected value to be in range [20, 30], but it was below:"
			if !strings.Contains(mockT.message, expectedMsg) {
				t.Errorf("Expected error message to contain %q, but got %q", expectedMsg, mockT.message)
			}
		})

		t.Run("should include custom message on failure when above range", func(t *testing.T) {
			t.Parallel()
			mockT := &mockT{}
			opts := []Option{WithMessage("Battery level must be valid")}
			BeInRange(mockT, 150, 0, 100, opts...)

			if !mockT.Failed() {
				t.Fatal("Expected test to fail, but it passed")
			}

			expectedMsg := "Battery level must be valid\nExpected value to be in range [0, 100], but it was above:"
			if !strings.Contains(mockT.message, expectedMsg) {
				t.Errorf("Expected error message to contain %q, but got %q", expectedMsg, mockT.message)
			}
		})
	})
}

func TestBeSorted(t *testing.T) {
	t.Parallel()
	// Helper function to reduce repetition
	runSortTest := func(collection any, shouldFail bool, msgCheck string) {
		mockTest := &mockT{}

		switch elements := collection.(type) {
		case []int:
			BeSorted(mockTest, elements)
		case [5]int:
			BeSorted(mockTest, elements[:]) // Convert array to slice
		case [3]int:
			BeSorted(mockTest, elements[:]) // Convert array to slice
		case []float64:
			BeSorted(mockTest, elements)
		case []string:
			BeSorted(mockTest, elements)
		case [3]string:
			BeSorted(mockTest, elements[:]) // Convert array to slice
		default:
			t.Fatalf("Unsupported type in test: %T", collection)
		}

		if shouldFail && !mockTest.Failed() {
			t.Error("Expected failure but test passed")
		}
		if !shouldFail && mockTest.Failed() {
			t.Errorf("Expected success but test failed: %s", mockTest.message)
		}
		if msgCheck != "" && mockTest.failed && !strings.Contains(mockTest.message, msgCheck) {
			t.Errorf("Expected message to contain %q, got:\n%s", msgCheck, mockTest.message)
		}
	}

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()

		// Successful cases
		runSortTest([]int{}, false, "")               // empty
		runSortTest([]int{42}, false, "")             // single element
		runSortTest([]int{1, 2, 3, 4, 5}, false, "")  // sorted
		runSortTest([5]int{1, 2, 3, 4, 5}, false, "") // sorted array
		runSortTest([]int{5, 5, 5}, false, "")        // duplicates

		// Failure cases
		runSortTest([]int{3, 1, 2}, true, "Index 0: 3 > 1")
		runSortTest([3]int{3, 1, 2}, true, "Index 0: 3 > 1")
		runSortTest([]int{5, 4, 3, 2, 1}, true, "4 order violations found")
	})

	t.Run("Type variations", func(t *testing.T) {
		t.Parallel()

		// Different types - sorted
		runSortTest([]float64{1.1, 2.2, 3.3}, false, "")
		runSortTest([]string{"apple", "banana", "cherry"}, false, "")
		runSortTest([3]string{"apple", "banana", "cherry"}, false, "")

		// Different types - unsorted
		runSortTest([]float64{1.1, 3.3, 2.2}, true, "Index 1: 3.3 > 2.2")
		runSortTest([]string{"banana", "apple"}, true, "Index 0: banana > apple")
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		runSortTest([]int{-5, -3, -1, 0, 2}, false, "")         // negatives sorted
		runSortTest([]int{-1, -5, 0}, true, "Index 0: -1 > -5") // negatives unsorted
		runSortTest([]float64{1.0, 1.1, 1.2}, false, "")        // float precision
	})

	t.Run("Large collection", func(t *testing.T) {
		t.Parallel()

		largeSlice := generateUnsortedLargeSlice(1000)
		runSortTest(largeSlice, true, "[Large collection]")
	})

	t.Run("Custom message", func(t *testing.T) {
		t.Parallel()

		customMessageMock := &mockT{}
		customMsg := "Values must be in order"
		BeSorted(customMessageMock, []int{3, 1, 2}, WithMessage(customMsg))

		if !customMessageMock.Failed() {
			t.Error("Expected failure but test passed")
		}
		if !strings.Contains(customMessageMock.message, customMsg) {
			t.Errorf("Expected custom message in error")
		}
	})

}

// generateUnsortedLargeSlice creates a large slice with some violations for testing
func generateUnsortedLargeSlice(size int) []int {
	slice := make([]int, size)
	for i := range size {
		slice[i] = i
	}
	// Add some violations
	if size > 10 {
		slice[5] = slice[5] + 10           // Make element 5 larger
		slice[size/2] = slice[size/2] - 10 // Make middle element smaller
		slice[size-5] = slice[size-5] - 20 // Make near-end element smaller
	}
	return slice
}
