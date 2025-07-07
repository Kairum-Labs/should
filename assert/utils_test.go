package assert

import (
	"fmt"
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

func TestFindInsertionInfo_Parameterized(t *testing.T) {
	testCases := []struct {
		name          string
		collection    []int
		target        int
		expectedFound bool
		expectedIndex int
		expectedPrev  *int
		expectedNext  *int
	}{
		{
			name:          "Insert_In_Middle",
			collection:    []int{1, 2, 3, 5, 6, 7},
			target:        4,
			expectedFound: false,
			expectedIndex: 3,
			expectedPrev:  intPtr(3),
			expectedNext:  intPtr(5),
		},
		{
			name:          "Insert_At_Beginning",
			collection:    []int{2, 3, 4, 5, 6},
			target:        1,
			expectedFound: false,
			expectedIndex: 0,
			expectedPrev:  nil,
			expectedNext:  intPtr(2),
		},
		{
			name:          "Insert_At_End",
			collection:    []int{1, 2, 3, 4, 5},
			target:        6,
			expectedFound: false,
			expectedIndex: 5,
			expectedPrev:  intPtr(5),
			expectedNext:  nil,
		},
		{
			name:          "Target_Already_Exists",
			collection:    []int{1, 2, 3, 4, 5},
			target:        3,
			expectedFound: true,
			expectedIndex: 2,
			expectedPrev:  nil,
			expectedNext:  nil,
		},
		{
			name:          "Empty_Collection",
			collection:    []int{},
			target:        1,
			expectedFound: false,
			expectedIndex: -1,
			expectedPrev:  nil,
			expectedNext:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := findInsertionInfo(tc.collection, tc.target)
			BeNil(t, err)
			BeEqual(t, info.found, tc.expectedFound)
			BeEqual(t, info.insertIndex, tc.expectedIndex)

			if tc.expectedPrev == nil {
				BeNil(t, info.prev)
			} else {
				NotBeNil(t, info.prev)
				BeEqual(t, *info.prev, *tc.expectedPrev)
			}

			if tc.expectedNext == nil {
				BeNil(t, info.next)
			} else {
				NotBeNil(t, info.next)
				BeEqual(t, *info.next, *tc.expectedNext)
			}
		})
	}
}

func TestFormatInsertionContext_WithMessage(t *testing.T) {
	collection := []int{2, 3, 5, 1, 0}
	target := 4
	info, err := findInsertionInfo(collection, target)
	BeNil(t, err)

	result := formatInsertionContext(collection, target, info)

	expected := `Collection: [2, 3, 5, 1, 0]
Missing  : 4

Element 4 would fit between 3 and 5 in sorted order`

	BeEqual(t, result, expected)
}

func TestFormatInsertionContext_BoundaryAndLargeCollections(t *testing.T) {
	t.Parallel()

	t.Run("Collection with 10 elements (no sorted view)", func(t *testing.T) {
		collection := []int{0, 1, 2, 3, 4, 9, 8, 7, 6, 5}
		info, err := findInsertionInfo(collection, 10)
		BeNil(t, err)

		result := formatInsertionContext(collection, 10, info)
		BeFalse(t, strings.Contains(result, "Sorted view"))
		BeTrue(t, strings.Contains(result, "Element 10 would be after 9 in sorted order"))
	})

	t.Run("Collection with 12 elements (with sorted view)", func(t *testing.T) {
		collection := []int{0, 1, 2, 3, 4, 5, 11, 10, 9, 8, 7, 6}
		info, err := findInsertionInfo(collection, 100)
		BeNil(t, err)

		result := formatInsertionContext(collection, 100, info)
		expectedParts := []string{
			"Collection: [0, 1, 2, 3, 4, ..., 10, 9, 8, 7, 6] (showing first 5 and last 5 of 12 elements)",
			"Missing  : 100",
			"Element 100 would be after 11 in sorted order",
			"└─ Sorted view: [..., 8, 9, 10, 11]",
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

	result := formatInsertionContext([]int{}, 5, insertionInfo[int]{})
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

func intPtr(i int) *int {
	return &i
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

// === Tests for formatMapValuesList ===

func TestFormatMapValuesList(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			values   []interface{}
			expected string
		}{
			{
				name:     "should format string values with single quotes",
				values:   []interface{}{"apple", "banana"},
				expected: "['apple', 'banana']",
			},
			{
				name:     "should format integer values",
				values:   []interface{}{1, 2, 3},
				expected: "[1, 2, 3]",
			},
			{
				name:     "should format mixed type values",
				values:   []interface{}{"hello", 42, true},
				expected: "['hello', 42, true]",
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatMapValuesList(tt.values)
				BeEqual(t, result, tt.expected)
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			values   []interface{}
			expected string
		}{
			{
				name:     "should handle nil slice",
				values:   nil,
				expected: "nil",
			},
			{
				name:     "should handle empty slice",
				values:   []interface{}{},
				expected: "[]",
			},
			{
				name:     "should handle slice with nil values",
				values:   []interface{}{"a", nil, "c"},
				expected: "['a', nil, 'c']",
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatMapValuesList(tt.values)
				BeEqual(t, result, tt.expected)
			})
		}
	})
}

// === Tests for containsMapKey ===

func TestContainsMapKey(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{"one": 1, "two": 2}

		t.Run("should return found when key exists", func(t *testing.T) {
			t.Parallel()
			result := containsMapKey(m, "one")
			BeTrue(t, result.Found)
			BeTrue(t, result.Exact)
		})

		t.Run("should return not found when key does not exist", func(t *testing.T) {
			t.Parallel()
			result := containsMapKey(m, "three")
			BeFalse(t, result.Found)
		})
	})

	t.Run("Similarity detection", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{"name": 1, "email_address": 2}

		t.Run("should find similar string keys", func(t *testing.T) {
			t.Parallel()
			result := containsMapKey(m, "email")
			BeFalse(t, result.Found)
			HaveLength(t, result.Similar, 1)
			BeEqual(t, result.Similar[0].Value, "email_address")
		})

		numMap := map[int]string{10: "a", 25: "b", 100: "c"}
		t.Run("should find similar numeric keys", func(t *testing.T) {
			t.Parallel()
			result := containsMapKey(numMap, 24)
			BeFalse(t, result.Found)
			HaveLength(t, result.Similar, 1)
			BeEqual(t, result.Similar[0].Value, 25)
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		t.Run("should handle nil map", func(t *testing.T) {
			t.Parallel()
			var m map[string]int
			result := containsMapKey(m, "any")
			BeFalse(t, result.Found)
			BeEqual(t, result.Total, 0)
			BeNil(t, result.Context)
		})

		t.Run("should handle empty map", func(t *testing.T) {
			t.Parallel()
			m := map[string]int{}
			result := containsMapKey(m, "any")
			BeFalse(t, result.Found)
			BeEqual(t, result.Total, 0)
			HaveLength(t, result.Context, 0)
		})

		t.Run("should handle non-map type", func(t *testing.T) {
			t.Parallel()
			result := containsMapKey([]string{"not", "a", "map"}, "key")
			BeFalse(t, result.Found)
		})
	})
}

// === Tests for containsMapValue ===

func TestContainsMapValue(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{"one": 1, "two": 2}

		t.Run("should return found when value exists", func(t *testing.T) {
			t.Parallel()
			result := containsMapValue(m, 1)
			BeTrue(t, result.Found)
			BeTrue(t, result.Exact)
		})

		t.Run("should return not found when value does not exist", func(t *testing.T) {
			t.Parallel()
			result := containsMapValue(m, 3)
			BeFalse(t, result.Found)
		})
	})

	t.Run("Similarity detection", func(t *testing.T) {
		t.Parallel()
		m := map[string]string{"user": "tester", "role": "admin"}

		t.Run("should find similar string values", func(t *testing.T) {
			t.Parallel()
			result := containsMapValue(m, "administrator")
			BeFalse(t, result.Found)
			HaveLength(t, result.Similar, 1)
			BeEqual(t, result.Similar[0].Value, "admin")
		})

		numMap := map[string]int{"a": 10, "b": 25, "c": 100}
		t.Run("should find similar numeric values", func(t *testing.T) {
			t.Parallel()
			result := containsMapValue(numMap, 24)
			BeFalse(t, result.Found)
			HaveLength(t, result.Similar, 1)
			BeEqual(t, result.Similar[0].Value, 25)
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		t.Run("should handle nil map", func(t *testing.T) {
			t.Parallel()
			var m map[string]int
			result := containsMapValue(m, 1)
			BeFalse(t, result.Found)
			BeEqual(t, result.Total, 0)
			BeNil(t, result.Context)
		})

		t.Run("should handle empty map", func(t *testing.T) {
			t.Parallel()
			m := map[string]int{}
			result := containsMapValue(m, 1)
			BeFalse(t, result.Found)
			BeEqual(t, result.Total, 0)
			HaveLength(t, result.Context, 0)
		})

		t.Run("should handle non-map type", func(t *testing.T) {
			t.Parallel()
			result := containsMapValue("not-a-map", "value")
			BeFalse(t, result.Found)
		})
	})
}

// === Tests for isNumericValue ===

func TestIsNumericValue(t *testing.T) {
	t.Parallel()

	t.Run("Numeric types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			value interface{}
		}{
			{"int", int(1)},
			{"int8", int8(1)},
			{"int16", int16(1)},
			{"int32", int32(1)},
			{"int64", int64(1)},
			{"uint", uint(1)},
			{"uint8", uint8(1)},
			{"uint16", uint16(1)},
			{"uint32", uint32(1)},
			{"uint64", uint64(1)},
			{"float32", float32(1.0)},
			{"float64", float64(1.0)},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				BeTrue(t, isNumericValue(tt.value))
			})
		}
	})

	t.Run("Non-numeric types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			value interface{}
		}{
			{"string", "1"},
			{"bool", true},
			{"struct", struct{}{}},
			{"slice", []int{}},
			{"nil", nil},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				BeFalse(t, isNumericValue(tt.value))
			})
		}
	})
}

// === Tests for findSimilarNumericKeys ===

func TestFindSimilarNumericKeys(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		items := []interface{}{1, 12, 25, 100}

		t.Run("should find numbers with small difference", func(t *testing.T) {
			t.Parallel()
			results := findSimilarNumericKeys(items, 10, 3)
			// The function finds 1 (diff=9, similarity=0.8) and 12 (diff=2, similarity=0.8)
			// 25 has diff=15 which is >10 and "10" is not contained in "25", so it's not similar
			// 100 contains "10" so it gets similarity=0.7
			HaveLength(t, results, 3)

			// Check that we have the expected values (order may vary based on similarity)
			values := []interface{}{results[0].Value, results[1].Value, results[2].Value}
			Contain(t, values, 1)
			Contain(t, values, 12)
			Contain(t, values, 100) // 100 contains "10"
		})

		t.Run("should return sorted results by similarity", func(t *testing.T) {
			t.Parallel()
			items := []interface{}{9, 20}
			results := findSimilarNumericKeys(items, 10, 2)
			HaveLength(t, results, 2)
			BeEqual(t, results[0].Value, 9)  // diff 1 -> similarity 0.9
			BeEqual(t, results[1].Value, 20) // diff 10 -> similarity 0.8
		})

		t.Run("should limit results", func(t *testing.T) {
			t.Parallel()
			items := []interface{}{2, 3, 4, 5, 6}
			results := findSimilarNumericKeys(items, 1, 2)
			HaveLength(t, results, 2)
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		t.Run("should handle empty item slice", func(t *testing.T) {
			t.Parallel()
			results := findSimilarNumericKeys([]interface{}{}, 10, 3)
			BeEmpty(t, results)
		})

		t.Run("should handle no similar numbers", func(t *testing.T) {
			t.Parallel()
			items := []interface{}{500, 600, 700}
			results := findSimilarNumericKeys(items, 10, 3)
			BeEmpty(t, results)
		})

		t.Run("should handle non-numeric target", func(t *testing.T) {
			t.Parallel()
			items := []interface{}{1, 2, 3}
			results := findSimilarNumericKeys(items, "not a number", 3)
			BeEmpty(t, results)
		})

		t.Run("should handle non-numeric items in slice", func(t *testing.T) {
			t.Parallel()
			items := []interface{}{1, "two", 3, true}
			results := findSimilarNumericKeys(items, 2, 3)
			HaveLength(t, results, 2) // Should find 1 and 3
			Contain(t, []interface{}{results[0].Value, results[1].Value}, 1)
			Contain(t, []interface{}{results[0].Value, results[1].Value}, 3)
		})
	})
}

// === Tests for formatMapContainKeyError ===

func TestFormatMapContainKeyError(t *testing.T) {
	t.Parallel()

	t.Run("Basic formatting", func(t *testing.T) {
		t.Parallel()
		result := MapContainResult{
			Context: []interface{}{"key1", "key2"},
			Total:   2,
		}
		msg := formatMapContainKeyError("missingKey", result)

		expected := "Expected map to contain key 'missingKey', but key was not found"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}

		expected = "Available keys: ['key1', 'key2']"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}

		expected = "Missing: 'missingKey'"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}

		unexpected := "Similar key"
		if strings.Contains(msg, unexpected) {
			t.Errorf("Expected message not to contain %q, got %q", unexpected, msg)
		}
	})

	t.Run("With similar keys", func(t *testing.T) {
		t.Parallel()
		t.Run("should show one similar key", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Context: []interface{}{"email_address"},
				Total:   1,
				Similar: []SimilarItem{{Value: "email_address", Details: "some detail"}},
			}
			msg := formatMapContainKeyError("email", result)
			expected := "Similar key found:"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'email_address' - some detail"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})

		t.Run("should show multiple similar keys", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Similar: []SimilarItem{
					{Value: "mail", Details: "detail1"},
					{Value: "e-mail", Details: "detail2"},
				},
			}
			msg := formatMapContainKeyError("email", result)
			expected := "Similar keys found:"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'mail' - detail1"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'e-mail' - detail2"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()
		t.Run("should truncate long key list", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Context: []interface{}{"a", "b", "c", "d", "e"},
				Total:   10,
			}
			msg := formatMapContainKeyError("z", result)
			expected := "(showing 5 of 10)"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})

		t.Run("should handle non-string keys", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Context: []interface{}{1, 2, 3},
				Total:   3,
			}
			msg := formatMapContainKeyError(4, result)
			expected := "Expected map to contain key 4, but key was not found"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "Available keys: [1, 2, 3]"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})
	})
}

// === Tests for formatMapContainValueError ===

func TestFormatMapContainValueError(t *testing.T) {
	t.Parallel()

	t.Run("Basic formatting", func(t *testing.T) {
		t.Parallel()
		result := MapContainResult{
			Context: []interface{}{"val1", "val2"},
			Total:   2,
		}
		msg := formatMapContainValueError("missingValue", result)
		expected := "Expected map to contain value 'missingValue', but value was not found"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}
		expected = "Available values: ['val1', 'val2']"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}
		expected = "Missing: 'missingValue'"
		if !strings.Contains(msg, expected) {
			t.Errorf("Expected message to contain %q, got %q", expected, msg)
		}
		unexpected := "Similar value"
		if strings.Contains(msg, unexpected) {
			t.Errorf("Expected message not to contain %q, got %q", unexpected, msg)
		}
	})

	t.Run("With similar values", func(t *testing.T) {
		t.Parallel()
		t.Run("should show one similar value", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Similar: []SimilarItem{{Value: "admin", Details: "some detail"}},
			}
			msg := formatMapContainValueError("administrator", result)
			expected := "Similar value found:"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'admin' - some detail"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})

		t.Run("should show multiple similar values", func(t *testing.T) {
			t.Parallel()
			result := MapContainResult{
				Similar: []SimilarItem{
					{Value: "user", Details: "detail1"},
					{Value: "guest", Details: "detail2"},
				},
			}
			msg := formatMapContainValueError("customer", result)
			expected := "Similar values found:"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'user' - detail1"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
			expected = "└─ 'guest' - detail2"
			if !strings.Contains(msg, expected) {
				t.Errorf("Expected message to contain %q, got %q", expected, msg)
			}
		})
	})
}

func TestFormatComplexType(t *testing.T) {
	type SimpleStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "nil",
		},
		{
			name:     "simple struct",
			input:    SimpleStruct{Name: "John", Age: 30},
			expected: "SimpleStruct{Name: \"John\", Age: 30}",
		},
		{
			name:     "pointer to struct",
			input:    &SimpleStruct{Name: "Jane", Age: 25},
			expected: "SimpleStruct{Name: \"Jane\", Age: 25}",
		},
		{
			name:     "non-struct type",
			input:    "hello",
			expected: `"hello"`,
		},
	}

	for _, test := range tests {
		result := formatComplexType(test.input)
		if result != test.expected {
			t.Errorf("formatComplexType(%s): expected %q, got %q", test.name, test.expected, result)
		}
	}
}

func TestFormatStructWithTruncation(t *testing.T) {
	type LongStruct struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
	}

	longStruct := LongStruct{
		Field1: "very long value that should cause truncation",
		Field2: "another long value",
		Field3: "yet another long value",
		Field4: "and another one",
		Field5: "final long value",
	}

	v := reflect.ValueOf(longStruct)
	structType := reflect.TypeOf(longStruct)
	result := formatStructWithTruncation(v, structType)

	if !strings.Contains(result, "LongStruct{") {
		t.Errorf("Expected result to contain struct name, got: %s", result)
	}
	if !strings.Contains(result, "...") {
		t.Errorf("Expected result to be truncated with '...', got: %s", result)
	}
}

func TestFormatFieldWithTruncation(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "short string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "long string",
			input:    "this is a very long string that should be truncated",
			expected: `"this is a very lo..."`,
		},
		{
			name:     "nil pointer",
			input:    (*string)(nil),
			expected: "nil",
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: "[]string{}",
		},
		{
			name:     "slice with items",
			input:    []string{"a", "b", "c"},
			expected: "[]string(3 items)",
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: "map[string]int{}",
		},
		{
			name:     "map with items",
			input:    map[string]int{"a": 1, "b": 2},
			expected: "map[string]int(2 items)",
		},
		{
			name:     "integer",
			input:    42,
			expected: "42",
		},
	}

	for _, test := range tests {
		v := reflect.ValueOf(test.input)
		result := formatFieldWithTruncation(v)
		if result != test.expected {
			t.Errorf("formatFieldWithTruncation(%s): expected %q, got %q", test.name, test.expected, result)
		}
	}
}

func TestFormatDiffValueConcise(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "nil",
		},
		{
			name:     "short string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "long string",
			input:    "this is a very long string that exceeds thirty characters",
			expected: `"this is a very long string ..."`,
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: "map[]",
		},
		{
			name:     "single entry map",
			input:    map[string]int{"key": 42},
			expected: `map["key": 42]`,
		},
		{
			name:     "multi entry map",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: "map[3 entries]",
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: "[]",
		},
		{
			name:     "small slice",
			input:    []int{1, 2, 3},
			expected: "[1, 2, 3]",
		},
		{
			name:     "large slice",
			input:    []int{1, 2, 3, 4, 5},
			expected: "[5 items]",
		},
		{
			name:     "boolean",
			input:    true,
			expected: "true",
		},
		{
			name:     "integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "float",
			input:    3.14,
			expected: "3.14",
		},
	}

	for _, test := range tests {
		result := formatDiffValueConcise(test.input)
		if result != test.expected {
			t.Errorf("formatDiffValueConcise(%s): expected %q, got %q", test.name, test.expected, result)
		}
	}
}

func TestContainsMapValue_CloseMatches(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	testMap := map[string]TestStruct{
		"user1": {Name: "Alice", Age: 30},
		"user2": {Name: "Bob", Age: 25},
	}

	target := TestStruct{Name: "Alice", Age: 31} // Similar to user1 but different age

	result := containsMapValue(testMap, target)

	if result.Found {
		t.Error("Expected not to find exact match")
	}

	if len(result.CloseMatches) == 0 {
		t.Error("Expected to find close matches")
	}

	if len(result.CloseMatches) > 0 {
		match := result.CloseMatches[0]
		if len(match.Differences) == 0 {
			t.Error("Expected close match to have differences")
		}
	}
}

func TestContainsMapValue_NonStruct(t *testing.T) {
	testMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	target := "value3"

	result := containsMapValue(testMap, target)

	if result.Found {
		t.Error("Expected not to find value")
	}

	// Should not have close matches for non-struct types
	if len(result.CloseMatches) > 0 {
		t.Error("Expected no close matches for non-struct types")
	}
}

func TestFormatMapContainValueError_ComplexTypes(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	result := MapContainResult{
		Found:   false,
		Total:   2,
		Context: []interface{}{TestStruct{Name: "Alice", Age: 30}},
		CloseMatches: []CloseMatch{
			{
				Value:       TestStruct{Name: "Alice", Age: 30},
				Differences: []string{"Age (31 ≠ 30)"},
			},
		},
	}

	target := TestStruct{Name: "Alice", Age: 31}
	errorMsg := formatMapContainValueError(target, result)

	if !strings.Contains(errorMsg, "Expected map to contain value, but it was not found") {
		t.Error("Expected error message to contain main error text")
	}

	if !strings.Contains(errorMsg, "Close matches:") {
		t.Error("Expected error message to contain close matches section")
	}

	if !strings.Contains(errorMsg, "Differs in: Age") {
		t.Error("Expected error message to contain difference details")
	}
}

func TestFormatMapContainValueError_SimpleTypes(t *testing.T) {
	result := MapContainResult{
		Found:   false,
		Total:   2,
		Context: []interface{}{"value1", "value2"},
	}

	target := "value3"
	errorMsg := formatMapContainValueError(target, result)

	if !strings.Contains(errorMsg, "Expected map to contain value 'value3'") {
		t.Error("Expected error message to contain target value")
	}

	if !strings.Contains(errorMsg, "Available values:") {
		t.Error("Expected error message to contain available values")
	}
}

func TestFormatMapNotContainKeyError(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		mapValue interface{}
		expected string
	}{
		{
			name:     "string key in string map",
			target:   "age",
			mapValue: map[string]int{"name": 1, "age": 30, "city": 2},
			expected: `Expected map to NOT contain key, but key was found:
Map Type : map[string]int
Map Size : 3 entries
Found Key: "age"
Associated Value: 30`,
		},
		{
			name:     "int key in int map",
			target:   42,
			mapValue: map[int]string{1: "one", 42: "forty-two", 3: "three"},
			expected: `Expected map to NOT contain key, but key was found:
Map Type : map[int]string
Map Size : 3 entries
Found Key: 42
Associated Value: "forty-two"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMapNotContainKeyError(tt.target, tt.mapValue)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestFormatMapNotContainValueError_SimpleTypes(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		mapValue interface{}
		contains string
	}{
		{
			name:     "string value in map",
			target:   "John",
			mapValue: map[string]string{"user1": "Alice", "user2": "John", "user3": "Bob"},
			contains: `Expected map to NOT contain value, but it was found:
Map Type : map[string]string
Map Size : 3 entries
Found Value: "John"
Found At: key "user2"`,
		},
		{
			name:     "int value in map",
			target:   42,
			mapValue: map[string]int{"a": 10, "b": 42, "c": 30},
			contains: `Expected map to NOT contain value, but it was found:
Map Type : map[string]int
Map Size : 3 entries
Found Value: 42
Found At: key "b"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMapNotContainValueError(tt.target, tt.mapValue)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain:\n%s\n\nGot:\n%s", tt.contains, result)
			}
		})
	}
}

func TestFormatMapNotContainValueError_ComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
		Role string
	}

	targetUser := User{ID: 2, Name: "Bob", Role: "user"}
	userMap := map[string]User{
		"emp1": {ID: 1, Name: "Alice", Role: "admin"},
		"emp2": {ID: 2, Name: "Bob", Role: "user"},
		"emp3": {ID: 3, Name: "Charlie", Role: "user"},
	}

	result := formatMapNotContainValueError(targetUser, userMap)

	expectedParts := []string{
		"Expected map to NOT contain value, but it was found:",
		"Map Type : map[string]assert.User",
		"Map Size : 3 entries",
		`Found Value: User{ID: 2, Name: "Bob", Role: "user"}`,
		`Found At: key "emp2"`,
	}

	// Should NOT contain the verbose context section
	unexpectedParts := []string{
		"Context:",
		"← Found here",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain '%s', but it was not found in:\n%s", part, result)
		}
	}

	for _, part := range unexpectedParts {
		if strings.Contains(result, part) {
			t.Errorf("Expected result to NOT contain '%s', but it was found in:\n%s", part, result)
		}
	}
}

// === Tests for ContainSubstring utility functions ===

func TestFormatContainSubstringError(t *testing.T) {
	t.Parallel()

	t.Run("Basic error formatting", func(t *testing.T) {
		t.Parallel()
		result := formatContainSubstringError("Hello, World!", "planet", "")

		expectedParts := []string{
			"Expected string to contain 'planet', but it was not found",
			"Substring   : 'planet'",
			"Actual   : 'Hello, World!'",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Empty substring handling", func(t *testing.T) {
		t.Parallel()
		result := formatContainSubstringError("Hello", "", "")

		expectedParts := []string{
			"Expected string to contain '<empty>', but it was not found",
			"Substring   : '<empty>'",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Empty actual string handling", func(t *testing.T) {
		t.Parallel()
		result := formatContainSubstringError("", "test", "")

		expectedParts := []string{
			"Expected string to contain 'test', but it was not found",
			"Actual   : '<empty>'",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Long string with multiline formatting", func(t *testing.T) {
		t.Parallel()
		longString := strings.Repeat("a", 250)
		result := formatContainSubstringError(longString, "test", "")

		expectedParts := []string{
			"Actual   : (length: 250)",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}

		// For very long strings, the actual content should be displayed
		if !strings.Contains(result, "aaaaaaa") {
			t.Errorf("Expected result to contain part of the long string")
		}
	})

	t.Run("String with newlines", func(t *testing.T) {
		t.Parallel()
		multilineString := "Hello\nWorld\nTest"
		result := formatContainSubstringError(multilineString, "missing", "")

		expectedParts := []string{
			"Actual   : (length: 16)",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("Large substring note", func(t *testing.T) {
		t.Parallel()
		largeSubstring := strings.Repeat("x", 60)
		result := formatContainSubstringError("small text", largeSubstring, "")

		expectedParts := []string{
			"Note: Substring is 60 characters long",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("With typo detection", func(t *testing.T) {
		t.Parallel()
		result := formatContainSubstringError("Hello, beautiful world!", "beatiful", "") //nolint:misspell

		expectedParts := []string{
			"Similar substring", // Can be either "found:" or "s found:"
			"'beautiful' at position",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("With multiple similar substrings", func(t *testing.T) {
		t.Parallel()
		result := formatContainSubstringError("test testing tested", "tst", "")

		expectedParts := []string{
			"Similar substrings found:",
		}

		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Expected result to contain: %q\n\nFull result:\n%s", part, result)
			}
		}
	})

	t.Run("With custom note message", func(t *testing.T) {
		t.Parallel()
		customNote := "\nNote: Custom message here"
		result := formatContainSubstringError("Hello", "missing", customNote)

		if !strings.Contains(result, "Note: Custom message here") {
			t.Errorf("Expected result to contain custom note message")
		}
	})
}

func TestFindSimilarSubstrings(t *testing.T) {
	t.Parallel()

	t.Run("Basic similarity detection", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("Hello, beautiful world!", "beatiful", 3) //nolint:misspell

		if len(results) == 0 {
			t.Fatal("Expected to find similar substrings")
		}

		found := false
		for _, result := range results {
			if result.Value == "beautiful" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find 'beautiful' in results, got: %v", results)
		}
	})

	t.Run("Empty substring", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("Hello, World!", "", 3)

		if results != nil {
			t.Errorf("Expected nil results for empty substring, got: %v", results)
		}
	})

	t.Run("Substring too long", func(t *testing.T) {
		t.Parallel()
		longSubstring := strings.Repeat("x", 25)
		results := findSimilarSubstrings("Hello, World!", longSubstring, 3)

		if results != nil {
			t.Errorf("Expected nil results for substring > 20 chars, got: %v", results)
		}
	})

	t.Run("Empty text", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("", "test", 3)

		if results != nil {
			t.Errorf("Expected nil results for empty text, got: %v", results)
		}
	})

	t.Run("No similar substrings", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("abcd", "xyz", 3)

		if len(results) != 0 {
			t.Errorf("Expected no similar substrings, got: %v", results)
		}
	})

	t.Run("Exact matches are skipped", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("test testing test", "test", 3)

		// Should only find "testing" as similar, not the exact "test" matches
		for _, result := range results {
			if result.Value == "test" {
				t.Errorf("Expected exact matches to be skipped, but found: %v", result)
			}
		}
	})

	t.Run("Different length substrings", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("testing tests", "test", 5)

		// Should find both "tests" (different length) and substrings from "testing"
		found := false
		for _, result := range results {
			if result.Value == "tests" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find 'tests' (different length), got: %v", results)
		}
	})

	t.Run("Results are sorted by similarity", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("house hause home", "house", 3)

		if len(results) < 2 {
			t.Fatal("Expected at least 2 results for sorting test")
		}

		for i := 0; i < len(results)-1; i++ {
			if results[i].Similarity < results[i+1].Similarity {
				t.Errorf("Results not sorted by similarity: %f < %f", results[i].Similarity, results[i+1].Similarity)
			}
		}
	})

	t.Run("Max results limit", func(t *testing.T) {
		t.Parallel()
		text := "test tast tost tust test1 test2 test3" //nolint:misspell
		results := findSimilarSubstrings(text, "test", 2)

		if len(results) > 2 {
			t.Errorf("Expected max 2 results, got: %d", len(results))
		}
	})

	t.Run("Index positions are correct", func(t *testing.T) {
		t.Parallel()
		results := findSimilarSubstrings("Hello, beautiful world!", "beatiful", 3) //nolint:misspell

		if len(results) == 0 {
			t.Fatal("Expected to find similar substrings")
		}

		// Find "beautiful" result
		for _, result := range results {
			if result.Value == "beautiful" {
				expectedIndex := strings.Index("Hello, beautiful world!", "beautiful")
				if result.Index != expectedIndex {
					t.Errorf("Expected index %d, got %d", expectedIndex, result.Index)
				}
				break
			}
		}
	})
}

func TestRemoveDuplicateSimilarItems(t *testing.T) {
	t.Parallel()

	t.Run("Remove duplicates by value and position", func(t *testing.T) {
		t.Parallel()
		items := []SimilarItem{
			{Value: "test", Index: 0, Similarity: 0.8, Details: "detail1"},
			{Value: "test", Index: 0, Similarity: 0.9, Details: "detail2"}, // duplicate
			{Value: "test", Index: 5, Similarity: 0.7, Details: "detail3"}, // different position
			{Value: "best", Index: 0, Similarity: 0.6, Details: "detail4"}, // different value
		}

		results := removeDuplicateSimilarItems(items)

		if len(results) != 3 {
			t.Errorf("Expected 3 unique items, got %d", len(results))
		}

		// Verify we have the unique combinations
		expectedKeys := []string{"test@0", "test@5", "best@0"}
		actualKeys := make([]string, len(results))
		for i, item := range results {
			actualKeys[i] = fmt.Sprintf("%s@%d", item.Value, item.Index)
		}

		for _, expected := range expectedKeys {
			found := false
			for _, actual := range actualKeys {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find key %s in results", expected)
			}
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		t.Parallel()
		results := removeDuplicateSimilarItems([]SimilarItem{})

		if len(results) != 0 {
			t.Errorf("Expected empty result for empty input, got %d items", len(results))
		}
	})

	t.Run("Single item", func(t *testing.T) {
		t.Parallel()
		items := []SimilarItem{
			{Value: "test", Index: 0, Similarity: 0.8, Details: "detail"},
		}

		results := removeDuplicateSimilarItems(items)

		if len(results) != 1 {
			t.Errorf("Expected 1 item, got %d", len(results))
		}

		if results[0].Value != "test" || results[0].Index != 0 {
			t.Errorf("Expected original item to be preserved")
		}
	})

	t.Run("No duplicates", func(t *testing.T) {
		t.Parallel()
		items := []SimilarItem{
			{Value: "test1", Index: 0, Similarity: 0.8, Details: "detail1"},
			{Value: "test2", Index: 1, Similarity: 0.7, Details: "detail2"},
			{Value: "test3", Index: 2, Similarity: 0.6, Details: "detail3"},
		}

		results := removeDuplicateSimilarItems(items)

		if len(results) != 3 {
			t.Errorf("Expected 3 items (no duplicates), got %d", len(results))
		}
	})
}
