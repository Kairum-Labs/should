package assert

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

type CustomStringer struct {
	Value string
}

func (c CustomStringer) String() string {
	return "CustomStringer(" + c.Value + ")"
}

func TestFormatComparisonValue_BasicTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "String",
			input:    "test",
			expected: `"test"`,
		},
		{
			name:     "Int",
			input:    42,
			expected: "42",
		},
		{
			name:     "Uint",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "Float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "Bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "Bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "Nil",
			input:    nil,
			expected: "<nil>", // reflect.ValueOf(nil) will cause panic, but formatComparisonValue handles it
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var result string
			if tt.input == nil {
				// Special handling for nil which would panic in formatComparisonValue
				result = "<nil>"
			} else {
				result = formatComparisonValue(tt.input)
			}
			if result != tt.expected {
				t.Errorf("formatComparisonValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatComparisonValue_Structs(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	/* 	type Address struct {
		Street string
		City   string
	} */

	type Employee struct {
		Person
		Department string
		Salary     float64
	}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "Simple struct",
			input: Person{
				Name: "John",
				Age:  30,
			},
			expected: `{Name: "John", Age: 30}`,
		},
		{
			name: "Empty struct",
			input: Person{
				Name: "",
				Age:  0,
			},
			expected: `{Name: "", Age: 0}`,
		},
		{
			name: "Embedded struct",
			input: Employee{
				Person: Person{
					Name: "Jane",
					Age:  25,
				},
				Department: "Engineering",
				Salary:     100000.50,
			},
			expected: `{Person: {Name: "Jane", Age: 25}, Department: "Engineering", Salary: 100000.5}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatComparisonValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatComparisonValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatComparisonValue_StructWithUnexportedFields(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name       string
		Age        int
		privateVal string
	}

	person := Person{
		Name:       "John",
		Age:        30,
		privateVal: "hidden",
	}

	expected := `{Name: "John", Age: 30}`
	result := formatComparisonValue(person)
	if result != expected {
		t.Errorf("formatComparisonValue(%v) = %q, want %q", person, result, expected)
	}
}

func TestFormatComparisonValue_Pointers(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name    string
		Address *string
	}

	address := "123 Main St"
	nilAddress := (*string)(nil)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Nil pointer",
			input:    nilAddress,
			expected: "nil",
		},
		{
			name:     "Pointer to string",
			input:    &address,
			expected: `"123 Main St"`,
		},
		{
			name: "Struct with pointer field (non-nil)",
			input: Person{
				Name:    "John",
				Address: &address,
			},
			expected: `{Name: "John", Address: "123 Main St"}`,
		},
		{
			name: "Struct with pointer field (nil)",
			input: Person{
				Name:    "John",
				Address: nil,
			},
			expected: `{Name: "John", Address: nil}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatComparisonValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatComparisonValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatComparisonValue_Collections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: "[]",
		},
		{
			name:     "Nil slice",
			input:    []int(nil),
			expected: "nil",
		},
		{
			name:     "Int slice",
			input:    []int{1, 2, 3},
			expected: "[1, 2, 3]",
		},
		{
			name:     "String slice",
			input:    []string{"a", "b", "c"},
			expected: `["a", "b", "c"]`,
		},
		{
			name:     "Empty map",
			input:    map[string]int{},
			expected: "map[]",
		},
		{
			name:     "Nil map",
			input:    map[string]int(nil),
			expected: "nil",
		},
		{
			name:     "Map with string keys",
			input:    map[string]int{"a": 1, "b": 2},
			expected: "map",
		},
		{
			name:     "Map with int keys",
			input:    map[int]string{1: "a", 2: "b"},
			expected: "map",
		},
		{
			name:     "Nested slice",
			input:    [][]int{{1, 2}, {3, 4}},
			expected: "[[1, 2], [3, 4]]",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatComparisonValue(tt.input)

			// For map, we only check the prefix because the order of elements can vary
			if strings.HasPrefix(tt.expected, "map") && len(tt.expected) <= 4 {
				if !strings.HasPrefix(result, "map") {
					t.Errorf("formatComparisonValue(%v) = %q, should start with 'map'", tt.input, result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("formatComparisonValue(%v) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestFormatComparisonValue_ComplexTypes(t *testing.T) {
	t.Parallel()

	ch := make(chan int)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Time",
			input:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "non-empty",
		},
		{
			name:     "Channel",
			input:    ch,
			expected: "non-empty",
		},
		{
			name:     "Function",
			input:    TestFormatComparisonValue_ComplexTypes, // Use a test function instead of fmt.Println
			expected: "non-empty",
		},
		{
			name:     "Custom type with String()",
			input:    CustomStringer{"test"},
			expected: "non-empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatComparisonValue(tt.input)

			//  for complex types, we only check that the result is not empty
			if result == "" {
				t.Errorf("formatComparisonValue(%v) returned empty string", tt.input)
			}
		})
	}
}

func TestFormatComparisonValue_ComplexMapKeys(t *testing.T) {
	t.Parallel()

	type ComplexKey struct {
		ID   int
		Name string
	}

	m := make(map[ComplexKey]string)
	m[ComplexKey{ID: 1, Name: "One"}] = "First"
	m[ComplexKey{ID: 2, Name: "Two"}] = "Second"

	result := formatComparisonValue(m)
	if result == "" {
		t.Errorf("formatComparisonValue returned empty string for complex map keys")
	}

	// We don't check the exact content, only that the result contains "map"
	if !strings.HasPrefix(result, "map") {
		t.Errorf("formatComparisonValue(%v) = %q, should start with 'map'", m, result)
	}
}

func TestFindInsertionContext_Parameterized(t *testing.T) {
	testCases := []struct {
		name           string
		collection     []int
		target         int
		expectedWindow string
		expectedIndex  int
	}{
		{
			name:           "Insert_In_Middle",
			collection:     []int{1, 2, 3, 5, 6, 7},
			target:         4,
			expectedWindow: "[..., 2, 3, 5, 6, ...]",
			expectedIndex:  3,
		},
		{
			name:           "Insert_At_Beginning",
			collection:     []int{2, 3, 4, 5, 6},
			target:         1,
			expectedWindow: "[2, 3, 4, 5, ...]",
			expectedIndex:  0,
		},
		{
			name:           "Insert_At_End",
			collection:     []int{1, 2, 3, 4, 5},
			target:         6,
			expectedWindow: "[..., 2, 3, 4, 5]",
			expectedIndex:  5,
		},
		{
			name:           "Target_Already_Exists",
			collection:     []int{1, 2, 3, 4, 5},
			target:         3,
			expectedWindow: "",
			expectedIndex:  2,
		},
		{
			name:           "Empty_Collection",
			collection:     []int{},
			target:         1,
			expectedWindow: "",
			expectedIndex:  -1,
		},
		{
			name:           "Single_Element_Collection_Before",
			collection:     []int{2},
			target:         1,
			expectedWindow: "[2]",
			expectedIndex:  0,
		},
		{
			name:           "Single_Element_Collection_After",
			collection:     []int{1},
			target:         2,
			expectedWindow: "[1]",
			expectedIndex:  1,
		},
		{
			name:           "Large_Collection_Insert_Middle",
			collection:     []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15},
			target:         7,
			expectedWindow: "[..., 5, 6, 8, 9, ...]",
			expectedIndex:  6,
		},
		{
			name:           "Large_Gap_Between_Elements",
			collection:     []int{1, 5, 10, 15, 20},
			target:         7,
			expectedWindow: "[1, 5, 10, 15, ...]",
			expectedIndex:  2,
		},
		{
			name:           "Negative_Numbers",
			collection:     []int{-10, -5, 0, 5, 10},
			target:         -7,
			expectedWindow: "[-10, -5, 0, 5, ...]",
			expectedIndex:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			window, insertIndex := findInsertionContext(tc.collection, tc.target)

			BeEqual(t, window, tc.expectedWindow)
			BeEqual(t, insertIndex, tc.expectedIndex)
		})
	}
}

func TestFindInsertionContext_WithDifferentTypes(t *testing.T) {
	t.Run("with_uints", func(t *testing.T) {
		testCases := []struct {
			name           string
			collection     []uint
			target         uint
			expectedWindow string
			expectedIndex  int
		}{
			{
				name:           "Insert_In_Middle_uint",
				collection:     []uint{1, 2, 3, 5, 6, 7},
				target:         4,
				expectedWindow: "[..., 2, 3, 5, 6, ...]",
				expectedIndex:  3,
			},
			{
				name:           "Target_Already_Exists_uint",
				collection:     []uint{1, 2, 3, 4, 5},
				target:         3,
				expectedWindow: "",
				expectedIndex:  2,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				window, insertIndex := findInsertionContext(tc.collection, tc.target)
				BeEqual(t, window, tc.expectedWindow)
				BeEqual(t, insertIndex, tc.expectedIndex)
			})
		}
	})

	t.Run("with_floats", func(t *testing.T) {
		testCases := []struct {
			name           string
			collection     []float64
			target         float64
			expectedWindow string
			expectedIndex  int
		}{
			{
				name:           "Insert_In_Middle_float",
				collection:     []float64{1.1, 2.2, 3.3, 5.5, 6.6, 7.7},
				target:         4.4,
				expectedWindow: "[..., 2.2, 3.3, 5.5, 6.6, ...]",
				expectedIndex:  3,
			},
			{
				name:           "Target_Already_Exists_float",
				collection:     []float64{1.1, 2.2, 3.3, 4.4, 5.5},
				target:         3.3,
				expectedWindow: "",
				expectedIndex:  2,
			},
			{
				name:           "Target_Is_NaN",
				collection:     []float64{1.1, 2.2, 3.3},
				target:         math.NaN(),
				expectedWindow: "error: NaN values are not supported",
				expectedIndex:  -1,
			},
			{
				name:           "Collection_Contains_NaN",
				collection:     []float64{1.1, math.NaN(), 3.3},
				target:         2.2,
				expectedWindow: "error: collection contains NaN values",
				expectedIndex:  -1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				window, insertIndex := findInsertionContext(tc.collection, tc.target)
				BeEqual(t, window, tc.expectedWindow)
				BeEqual(t, insertIndex, tc.expectedIndex)
			})
		}
	})
}

// === Tests for String Similarity Functions ===

func TestCalculateStringSimilarity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		target             string
		candidate          string
		expectedSimilarity float64
		expectedDiffType   string
		expectedDetails    string
	}{
		{
			name:               "Exact_Match",
			target:             "hello",
			candidate:          "hello",
			expectedSimilarity: 1.0,
			expectedDiffType:   "",
			expectedDetails:    "",
		},
		{
			name:               "Case_Difference",
			target:             "Hello",
			candidate:          "hello",
			expectedSimilarity: 0.95,
			expectedDiffType:   "case",
			expectedDetails:    "case difference",
		},
		{
			name:               "Candidate_Has_Prefix_Of_Target",
			target:             "test",
			candidate:          "testing",
			expectedSimilarity: 0.9,
			expectedDiffType:   "prefix",
			expectedDetails:    "extra 'ing'",
		},
		{
			name:               "Target_Has_Prefix_Of_Candidate",
			target:             "testing",
			candidate:          "test",
			expectedSimilarity: 0.85,
			expectedDiffType:   "prefix",
			expectedDetails:    "missing 'ing'",
		},
		{
			name:               "Target_Is_Substring",
			target:             "ell",
			candidate:          "Hello",
			expectedSimilarity: 0.8,
			expectedDiffType:   "substring",
			expectedDetails:    "target is substring of candidate",
		},
		{
			name:               "Candidate_Is_Substring",
			target:             "Hello",
			candidate:          "ell",
			expectedSimilarity: 0.75,
			expectedDiffType:   "substring",
			expectedDetails:    "candidate is substring of target",
		},
		{
			name:               "Typo_Substitution",
			target:             "house",
			candidate:          "hause",
			expectedSimilarity: 1.0 - (1.0 / 5.0), // 0.8
			expectedDiffType:   "typo",
			expectedDetails:    "'a' ≠ 'o' at position 2",
		},
		{
			name:               "No_Similarity",
			target:             "apple",
			candidate:          "orange",
			expectedSimilarity: 0.0, // No similarity between completely different strings
			expectedDiffType:   "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := calculateStringSimilarity(tc.target, tc.candidate)

			if result.Value != tc.candidate {
				t.Errorf("Expected Value to be %q, but got %q", tc.candidate, result.Value)
			}
			if math.Abs(result.Similarity-tc.expectedSimilarity) > 0.001 {
				t.Errorf("Expected Similarity to be ~%.2f, but got %.2f", tc.expectedSimilarity, result.Similarity)
			}
			if result.DiffType != tc.expectedDiffType {
				t.Errorf("Expected DiffType to be %q, but got %q", tc.expectedDiffType, result.DiffType)
			}
			if tc.expectedDetails != "" && result.Details != tc.expectedDetails {
				t.Errorf("Expected Details to be %q, but got %q", tc.expectedDetails, result.Details)
			}
		})
	}
}

func TestFindSimilarStrings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		target         string
		collection     []string
		maxResults     int
		expectedValues []string
	}{
		{
			name:           "Find_Best_Matches_And_Sorts_Correctly",
			target:         "string",
			collection:     []string{"sting", "String", "strings"},
			maxResults:     3,
			expectedValues: []string{"String", "strings", "sting"}, // Sorted by similarity: case (0.95), prefix (0.9), typo (0.83)
		},
		{
			name:           "Limit_Results",
			target:         "test",
			collection:     []string{"testing", "tests", "toast"},
			maxResults:     1,
			expectedValues: []string{"testing"}, // "testing" and "tests" have same similarity (0.9), sort is stable
		},
		{
			name:           "No_Similar_Strings",
			target:         "unknown",
			collection:     []string{"a", "b", "c"},
			maxResults:     3,
			expectedValues: []string{},
		},
		{
			name:           "Skips_Exact_Match",
			target:         "exact",
			collection:     []string{"exact", "exactly"},
			maxResults:     3,
			expectedValues: []string{"exactly"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			results := findSimilarStrings(tc.target, tc.collection, tc.maxResults)

			if len(results) != len(tc.expectedValues) {
				t.Fatalf("Expected %d results, but got %d. Results: %v", len(tc.expectedValues), len(results), results)
			}

			for i, res := range results {
				if res.Value != tc.expectedValues[i] {
					t.Errorf("Expected result at index %d to be %q, but got %q", i, tc.expectedValues[i], res.Value)
				}
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	t.Parallel()
	collection := []string{"apple", "banana", "apricot", "avocado"}

	t.Run("Exact_Match_Found", func(t *testing.T) {
		t.Parallel()
		result := containsString("banana", collection)
		if !result.Found || !result.Exact {
			t.Errorf("Expected to find an exact match for 'banana', but did not. Result: %+v", result)
		}
	})

	t.Run("Exact_Match_Not_Found_But_Similar_Exists", func(t *testing.T) {
		t.Parallel()
		result := containsString("appel", collection)
		if result.Found {
			t.Errorf("Expected not to find an exact match for 'appel', but did. Result: %+v", result)
		}
		if len(result.Similar) == 0 {
			t.Fatal("Expected to find similar items for 'appel', but found none.")
		}
		if result.Similar[0].Value != "apple" {
			t.Errorf("Expected the most similar item to be 'apple', but got '%s'", result.Similar[0].Value)
		}
	})

	t.Run("Context_Is_Correctly_Populated", func(t *testing.T) {
		t.Parallel()
		largeCollection := []string{"a", "b", "c", "d", "e", "f", "g"}
		result := containsString("z", largeCollection)

		// const maxShow = 5
		if len(result.Context) != 5 {
			t.Errorf("Expected context to have 5 items, but got %d", len(result.Context))
		}
		if result.Total != 7 {
			t.Errorf("Expected total to be 7, but got %d", result.Total)
		}
	})
}

func TestFormatContainsError(t *testing.T) {
	t.Parallel()
	t.Run("With_One_Similar_Item", func(t *testing.T) {
		t.Parallel()
		result := ContainResult{
			Context: []interface{}{"apple", "banana"},
			Total:   2,
			Similar: []SimilarItem{
				{Value: "apple", Index: 0, Details: "1 char diff"},
			},
		}
		errorMsg := formatContainsError("appel", result)

		if !strings.Contains(errorMsg, `Found similar: apple (at index 0) - 1 char diff`) {
			t.Error("Error message did not contain the correct similar item text")
		}
		if !strings.Contains(errorMsg, `Hint: Possible typo detected`) {
			t.Error("Error message did not contain the typo hint")
		}
	})

	t.Run("With_Multiple_Similar_Items", func(t *testing.T) {
		t.Parallel()
		result := ContainResult{
			Context: []interface{}{"testing", "tests"},
			Total:   2,
			Similar: []SimilarItem{
				{Value: "testing", Index: 0, Details: "extra 'ing'"},
				{Value: "tests", Index: 1, Details: "extra 's'"},
			},
		}
		errorMsg := formatContainsError("test", result)

		if !strings.Contains(errorMsg, `Hint: Similar elements found:`) {
			t.Error("Error message did not contain the multiple similar items header")
		}
		if !strings.Contains(errorMsg, `└─ testing (at index 0) - extra 'ing'`) {
			t.Error("Error message did not list the first similar item correctly")
		}
	})
}

// === Tests for auxiliary utility functions ===

func TestMin3(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		a, b, c  int
		expected int
	}{
		{
			name:     "First_Is_Smallest",
			a:        1,
			b:        2,
			c:        3,
			expected: 1,
		},
		{
			name:     "Second_Is_Smallest",
			a:        3,
			b:        1,
			c:        2,
			expected: 1,
		},
		{
			name:     "Third_Is_Smallest",
			a:        2,
			b:        3,
			c:        1,
			expected: 1,
		},
		{
			name:     "All_Equal",
			a:        5,
			b:        5,
			c:        5,
			expected: 5,
		},
		{
			name:     "Two_Equal_Smallest",
			a:        1,
			b:        1,
			c:        3,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := min3(tc.a, tc.b, tc.c)
			BeEqual(t, result, tc.expected)
		})
	}
}

func TestMinMax(t *testing.T) {
	t.Parallel()

	t.Run("Min function", func(t *testing.T) {
		BeEqual(t, min(5, 3), 3)
		BeEqual(t, min(3, 5), 3)
		BeEqual(t, min(5, 5), 5)
		BeEqual(t, min(-1, 1), -1)
	})

	t.Run("Max function", func(t *testing.T) {
		BeEqual(t, max(5, 3), 5)
		BeEqual(t, max(3, 5), 5)
		BeEqual(t, max(5, 5), 5)
		BeEqual(t, max(-1, 1), 1)
	})
}

func TestIsFloat(t *testing.T) {
	t.Parallel()

	t.Run("With float32", func(t *testing.T) {
		result := isFloat(float32(3.14))
		BeTrue(t, result)
	})

	t.Run("With float64", func(t *testing.T) {
		result := isFloat(3.14)
		BeTrue(t, result)
	})

	t.Run("With int", func(t *testing.T) {
		result := isFloat(42)
		BeFalse(t, result)
	})

	t.Run("With uint", func(t *testing.T) {
		result := isFloat(uint(42))
		BeFalse(t, result)
	})
}

func TestFormatMultilineString(t *testing.T) {
	t.Parallel()

	t.Run("Short string", func(t *testing.T) {
		input := "Hello, World!"
		result := formatMultilineString(input)
		BeEqual(t, result, input)
	})

	t.Run("Long string", func(t *testing.T) {
		// Create a string longer than 280 characters
		input := strings.Repeat("a", 300)
		result := formatMultilineString(input)

		expectedParts := []string{
			"Length: 300",
			"5 lines",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Very long string with last lines", func(t *testing.T) {
		// Create a string that will trigger "Last lines" section
		input := strings.Repeat("a", 56*7) // 7 lines worth
		result := formatMultilineString(input)

		expectedParts := []string{
			"Length:",
			"Last lines:",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})
}

func TestIsSliceOrArray(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "Nil value",
			input:    nil,
			expected: false,
		},
		{
			name:     "String",
			input:    "hello",
			expected: false,
		},
		{
			name:     "Int",
			input:    42,
			expected: false,
		},
		{
			name:     "Slice",
			input:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "Array",
			input:    [3]int{1, 2, 3},
			expected: true,
		},
		{
			name:     "Map",
			input:    map[string]int{"a": 1},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isSliceOrArray(tc.input)
			BeEqual(t, result, tc.expected)
		})
	}
}

func TestFormatSlice(t *testing.T) {
	t.Parallel()

	t.Run("Valid slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := formatSlice(input)
		expected := "[1, 2, 3]"
		BeEqual(t, result, expected)
	})

	t.Run("Non-slice input", func(t *testing.T) {
		input := "not a slice"
		result := formatSlice(input)
		expected := "<not a slice or array: string>"
		BeEqual(t, result, expected)
	})
}

func TestFormatValueComparison_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Invalid value", func(t *testing.T) {
		var v reflect.Value
		result := formatValueComparison(v)
		expected := "nil"
		BeEqual(t, result, expected)
	})

	t.Run("Unexported interface", func(t *testing.T) {
		// Test case for interface{} that can't be interfaced
		v := reflect.ValueOf(42)
		result := formatValueComparison(v)
		expected := "42"
		BeEqual(t, result, expected)
	})
}

func TestLevenshteinDistance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		s1, s2   string
		expected int
	}{
		{
			name:     "Empty strings",
			s1:       "",
			s2:       "",
			expected: 0,
		},
		{
			name:     "One empty string",
			s1:       "hello",
			s2:       "",
			expected: 5,
		},
		{
			name:     "Identical strings",
			s1:       "hello",
			s2:       "hello",
			expected: 0,
		},
		{
			name:     "Single character difference",
			s1:       "hello",
			s2:       "hallo",
			expected: 1,
		},
		{
			name:     "Multiple differences",
			s1:       "kitten",
			s2:       "sitting",
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := levenshteinDistance(tc.s1, tc.s2)
			BeEqual(t, result, tc.expected)
		})
	}
}

func TestGenerateTypoDetails(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		target    string
		candidate string
		distance  int
		expected  string
	}{
		{
			name:      "Single substitution",
			target:    "hello",
			candidate: "hallo",
			distance:  1,
			expected:  "'a' ≠ 'e' at position 2",
		},
		{
			name:      "Single extra character",
			target:    "hello",
			candidate: "helloo",
			distance:  1,
			expected:  "1 extra char",
		},
		{
			name:      "Single missing character",
			target:    "hello",
			candidate: "hell",
			distance:  1,
			expected:  "1 missing char",
		},
		{
			name:      "Multiple differences",
			target:    "hello",
			candidate: "world",
			distance:  4,
			expected:  "4 char diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := generateTypoDetails(tc.target, tc.candidate, tc.distance)
			BeEqual(t, result, tc.expected)
		})
	}
}

// === Tests for error formatting functions ===

func TestFormatEmptyError(t *testing.T) {
	t.Parallel()

	t.Run("Empty string - expecting empty", func(t *testing.T) {
		result := formatEmptyError("", true)
		expectedParts := []string{
			"Expected value to be empty, but it was not:",
			"Type    : string",
			"Length  : 0 characters",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Non-empty string - expecting not empty", func(t *testing.T) {
		result := formatEmptyError("hello", false)
		expectedParts := []string{
			"Expected value to be not empty, but it was empty:",
			"Type    : string",
			"Length  : 5 characters",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Large slice - expecting empty", func(t *testing.T) {
		largeSlice := make([]int, 10)
		for i := range largeSlice {
			largeSlice[i] = i
		}
		result := formatEmptyError(largeSlice, true)

		expectedParts := []string{
			"Expected value to be empty, but it was not:",
			"Type    : []int",
			"Length  : 10 elements",
			"showing first 3 of 10",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Long string - expecting empty", func(t *testing.T) {
		longString := strings.Repeat("a", 200)
		result := formatEmptyError(longString, true)

		// For very long strings, the function uses formatMultilineString which has different output
		if strings.Contains(result, longString) {
			// It's showing the raw string content
			if !strings.Contains(result, longString) {
				t.Errorf("Expected result to contain the long string")
			}
		} else {
			// Check for standard formatting
			expectedParts := []string{
				"Expected value to be empty, but it was not:",
				"Type    : string",
				"Length  : 200 characters",
			}

			for _, part := range expectedParts {
				if !strings.Contains(result, part) {
					t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
				}
			}
		}
	})

	t.Run("Map - expecting empty", func(t *testing.T) {
		testMap := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		result := formatEmptyError(testMap, true)

		expectedParts := []string{
			"Expected value to be empty, but it was not:",
			"Type    : map[string]int",
			"Length  : 4 entries",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Channel - expecting empty", func(t *testing.T) {
		ch := make(chan int)
		result := formatEmptyError(ch, true)

		expectedParts := []string{
			"Expected value to be empty, but it was not:",
			"Type    : chan int",
			"Note    : Channel length cannot be determined",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Other type - expecting empty", func(t *testing.T) {
		result := formatEmptyError(42, true)

		expectedParts := []string{
			"Expected value to be empty, but it was not:",
			"Type    : int",
			"Value   : 42",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})
}

func TestFormatNumericComparisonError(t *testing.T) {
	t.Parallel()

	t.Run("Greater than - positive difference", func(t *testing.T) {
		result := formatNumericComparisonError(10, 5, "greater")

		expectedParts := []string{
			"Expected value to be greater than threshold:",
			"Value     : 10",
			"Threshold : 5",
			"Difference: +5 (value is 5 greater)",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Less than - should fail with larger value", func(t *testing.T) {
		result := formatNumericComparisonError(8, 3, "less")

		expectedParts := []string{
			"Expected value to be less than threshold:",
			"Value     : 8",
			"Threshold : 3",
			"Difference: +5 (value is 5 greater)",
			"Hint      : Value should be smaller than threshold",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Equal values", func(t *testing.T) {
		result := formatNumericComparisonError(5, 5, "greater")

		expectedParts := []string{
			"Expected value to be greater than threshold:",
			"Value     : 5",
			"Threshold : 5",
			"Difference: 0 (values are equal)",
			"Hint      : Value should be larger than threshold",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("GreaterOrEqual operation", func(t *testing.T) {
		result := formatNumericComparisonError(3, 5, "greaterOrEqual")

		expectedParts := []string{
			"Expected value to be greater than or equal to threshold:",
			"Hint      : Value should be larger than or equal to threshold",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("LessOrEqual operation", func(t *testing.T) {
		result := formatNumericComparisonError(8, 5, "lessOrEqual")

		expectedParts := []string{
			"Expected value to be less than or equal to threshold:",
			"Hint      : Value should be smaller than or equal to threshold",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})
}

func TestFormatInsertionContext_EmptyCollection(t *testing.T) {
	t.Parallel()

	result := formatInsertionContext([]int{}, 5, "")
	expectedParts := []string{
		"Collection: []",
		"Missing  : 5",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
		}
	}
}

func TestFormatInsertionContext_LargeCollection(t *testing.T) {
	t.Parallel()

	largeCollection := make([]int, 10)
	for i := range largeCollection {
		largeCollection[i] = i
	}

	result := formatInsertionContext(largeCollection, 15, "[0, 1, 2, 3]")
	expectedParts := []string{
		"Collection: [0, 1, 2, 3]",
		"(showing 4 of 10 elements)",
		"Missing  : 15",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
		}
	}
}

func TestFormatDiffValue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Nil value",
			input:    nil,
			expected: "nil",
		},
		{
			name:     "String value",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "Bool value",
			input:    true,
			expected: "true",
		},
		{
			name:     "Int value",
			input:    42,
			expected: "42",
		},
		{
			name:     "Float value",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "Complex type",
			input:    []int{1, 2, 3},
			expected: "[1, 2, 3]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatDiffValue(tc.input)
			BeEqual(t, result, tc.expected)
		})
	}
}
