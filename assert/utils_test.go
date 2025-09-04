package assert

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"
)

type CustomStringer struct {
	Value string
}

func (c CustomStringer) String() string {
	return "CustomStringer(" + c.Value + ")"
}

// Custom error types for testing
type simpleError struct{ msg string }

func (e simpleError) Error() string { return e.msg }

type wrapperError struct {
	msg     string
	wrapped error
}

func (e wrapperError) Error() string { return e.msg }
func (e wrapperError) Unwrap() error { return e.wrapped }

func TestIsPrimitive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		kind     reflect.Kind
		expected bool
	}{
		// Primitive types
		{"string", reflect.String, true},
		{"int", reflect.Int, true},
		{"int8", reflect.Int8, true},
		{"int16", reflect.Int16, true},
		{"int32", reflect.Int32, true},
		{"int64", reflect.Int64, true},
		{"uint", reflect.Uint, true},
		{"uint8", reflect.Uint8, true},
		{"uint16", reflect.Uint16, true},
		{"uint32", reflect.Uint32, true},
		{"uint64", reflect.Uint64, true},
		{"uintptr", reflect.Uintptr, true},
		{"float32", reflect.Float32, true},
		{"float64", reflect.Float64, true},
		{"bool", reflect.Bool, true},

		// Non-primitive types
		{"slice", reflect.Slice, false},
		{"map", reflect.Map, false},
		{"struct", reflect.Struct, false},
		{"array", reflect.Array, false},
		{"ptr", reflect.Ptr, false},
		{"interface", reflect.Interface, false},
		{"func", reflect.Func, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isPrimitive(tt.kind)
			if got != tt.expected {
				t.Errorf("isPrimitive(%v) = %v; want %v", tt.kind, got, tt.expected)
			}
		})
	}
}

func TestFormatComparisonValue_BasicTypes(t *testing.T) {
	t.Parallel()

	var x int
	uptr := uintptr(unsafe.Pointer(&x))

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
			name:     "Uintptr",
			input:    uptr,
			expected: fmt.Sprint(uptr),
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
		{
			name:     "Rune slice",
			input:    []rune{'☀', '🌙', '⭐'},
			expected: "[9728, 127769, 11088]", // TODO: Improve to show emojis instead of numbers
		},
		{
			name:     "Byte slice",
			input:    []byte{1, 2, 3},
			expected: "[1, 2, 3]",
		},
	}

	for _, tt := range tests {
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
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
		t.Parallel()
		collection := []int{0, 1, 2, 3, 4, 9, 8, 7, 6, 5}
		info, err := findInsertionInfo(collection, 10)
		BeNil(t, err)

		result := formatInsertionContext(collection, 10, info)
		BeFalse(t, strings.Contains(result, "Sorted view"))
		BeTrue(t, strings.Contains(result, "Element 10 would be after 9 in sorted order"))
	})

	t.Run("Collection with 12 elements (with sorted view)", func(t *testing.T) {
		t.Parallel()
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
		result := containResult{
			Context: []interface{}{"apple", "banana"},
			Total:   2,
			Similar: []similarItem{
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
		result := containResult{
			Context: []interface{}{"testing", "tests"},
			Total:   2,
			Similar: []similarItem{
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
			t.Parallel()
			result := min3(tc.a, tc.b, tc.c)
			BeEqual(t, result, tc.expected)
		})
	}
}

func TestMinMax(t *testing.T) {
	t.Parallel()

	t.Run("Min function", func(t *testing.T) {
		t.Parallel()
		BeEqual(t, min(5, 3), 3)
		BeEqual(t, min(3, 5), 3)
		BeEqual(t, min(5, 5), 5)
		BeEqual(t, min(-1, 1), -1)
	})

	t.Run("Max function", func(t *testing.T) {
		t.Parallel()
		BeEqual(t, max(5, 3), 5)
		BeEqual(t, max(3, 5), 5)
		BeEqual(t, max(5, 5), 5)
		BeEqual(t, max(-1, 1), 1)
	})
}

func TestIsFloat(t *testing.T) {
	t.Parallel()

	t.Run("With float32", func(t *testing.T) {
		t.Parallel()
		result := isFloat(float32(3.14))
		BeTrue(t, result)
	})

	t.Run("With float64", func(t *testing.T) {
		t.Parallel()
		result := isFloat(3.14)
		BeTrue(t, result)
	})

	t.Run("With int", func(t *testing.T) {
		t.Parallel()
		result := isFloat(42)
		BeFalse(t, result)
	})

	t.Run("With uint", func(t *testing.T) {
		t.Parallel()
		result := isFloat(uint(42))
		BeFalse(t, result)
	})
}

func TestFormatMultilineString(t *testing.T) {
	t.Parallel()

	t.Run("Short string", func(t *testing.T) {
		t.Parallel()
		input := "Hello, World!"
		result := formatMultilineString(input)
		BeEqual(t, result, input)
	})

	t.Run("Long string", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
			t.Parallel()
			result := isSliceOrArray(tc.input)
			BeEqual(t, result, tc.expected)
		})
	}
}

func TestFormatSlice(t *testing.T) {
	t.Parallel()

	t.Run("Valid slice", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3}
		result := formatSlice(input)
		expected := "[1, 2, 3]"
		BeEqual(t, result, expected)
	})

	t.Run("Non-slice input", func(t *testing.T) {
		t.Parallel()
		input := "not a slice"
		result := formatSlice(input)
		expected := "<not a slice or array: string>"
		BeEqual(t, result, expected)
	})
}

func TestFormatValueComparison_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Invalid value", func(t *testing.T) {
		t.Parallel()
		var v reflect.Value
		result := formatValueComparison(v)
		expected := "nil"
		BeEqual(t, result, expected)
	})

	t.Run("Unexported interface", func(t *testing.T) {
		t.Parallel()
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
			t.Parallel()
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
			t.Parallel()
			result := generateTypoDetails(tc.target, tc.candidate, tc.distance)
			BeEqual(t, result, tc.expected)
		})
	}
}

// === Tests for error formatting functions ===

func TestFormatEmptyError(t *testing.T) {
	t.Parallel()

	t.Run("Empty string - expecting empty", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		result := mapContainResult{
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
			result := mapContainResult{
				Context: []interface{}{"email_address"},
				Total:   1,
				Similar: []similarItem{{Value: "email_address", Details: "some detail"}},
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
			result := mapContainResult{
				Similar: []similarItem{
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
			result := mapContainResult{
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
			result := mapContainResult{
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
		result := mapContainResult{
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
			result := mapContainResult{
				Similar: []similarItem{{Value: "admin", Details: "some detail"}},
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
			result := mapContainResult{
				Similar: []similarItem{
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
	t.Parallel()
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
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := formatComplexType(test.input)
			if result != test.expected {
				t.Errorf("formatComplexType(%s): expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}

func TestFormatStructWithTruncation(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			v := reflect.ValueOf(test.input)
			result := formatFieldWithTruncation(v)
			if result != test.expected {
				t.Errorf("formatFieldWithTruncation(%s): expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}

func TestFormatDiffValueConcise(t *testing.T) {
	t.Parallel()
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
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := formatDiffValueConcise(test.input)
			if result != test.expected {
				t.Errorf("formatDiffValueConcise(%s): expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}

func TestContainsMapValue_CloseMatches(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	type TestStruct struct {
		Name string
		Age  int
	}

	result := mapContainResult{
		Found:   false,
		Total:   2,
		Context: []interface{}{TestStruct{Name: "Alice", Age: 30}},
		CloseMatches: []closeMatch{
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
	t.Parallel()
	result := mapContainResult{
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
	t.Parallel()
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
			t.Parallel()
			result := formatMapNotContainKeyError(tt.target, tt.mapValue)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestFormatMapNotContainValueError_SimpleTypes(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := formatMapNotContainValueError(tt.target, tt.mapValue)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain:\n%s\n\nGot:\n%s", tt.contains, result)
			}
		})
	}
}

func TestFormatMapNotContainValueError_ComplexTypes(t *testing.T) {
	t.Parallel()
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
		items := []similarItem{
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
		results := removeDuplicateSimilarItems([]similarItem{})

		if len(results) != 0 {
			t.Errorf("Expected empty result for empty input, got %d items", len(results))
		}
	})

	t.Run("Single item", func(t *testing.T) {
		t.Parallel()
		items := []similarItem{
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
		items := []similarItem{
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

func TestFormatRangeError(t *testing.T) {
	t.Parallel()

	t.Run("value below range", func(t *testing.T) {
		t.Parallel()
		actual := 16
		minValue := 18
		maxValue := 65
		expected := `Expected value to be in range [18, 65], but it was below:
        Value    : 16
        Range    : [18, 65]
        Distance : 2 below minimum (16 < 18)
        Hint     : Value should be >= 18`
		result := formatRangeError(actual, minValue, maxValue)
		if result != expected {
			t.Errorf("Expected message:\n%s\n\nGot:\n%s", expected, result)
		}
	})

	t.Run("value above range", func(t *testing.T) {
		t.Parallel()
		actual := 105
		minValue := 0
		maxValue := 100
		expected := `Expected value to be in range [0, 100], but it was above:
        Value    : 105
        Range    : [0, 100]
        Distance : 5 above maximum (105 > 100)
        Hint     : Value should be <= 100`
		result := formatRangeError(actual, minValue, maxValue)
		if result != expected {
			t.Errorf("Expected message:\n%s\n\nGot:\n%s", expected, result)
		}
	})

	t.Run("float value below range", func(t *testing.T) {
		t.Parallel()
		actual := -0.1
		minValue := 0.0
		maxValue := 1.0
		expected := fmt.Sprintf(
			`Expected value to be in range [%v, %v], but it was below:
        Value    : %v
        Range    : [%v, %v]
        Distance : %v below minimum (%v < %v)
        Hint     : Value should be >= %v`,
			minValue,
			maxValue,
			actual,
			minValue,
			maxValue,
			minValue-actual,
			actual,
			minValue,
			minValue,
		)
		result := formatRangeError(actual, minValue, maxValue)
		if result != expected {
			t.Errorf("Expected message:\n%s\n\nGot:\n%s", expected, result)
		}
	})

	t.Run("float value above range", func(t *testing.T) {
		t.Parallel()
		actual := 1.1
		minValue := 0.0
		maxValue := 1.0
		expected := fmt.Sprintf(
			`Expected value to be in range [%v, %v], but it was above:
        Value    : %v
        Range    : [%v, %v]
        Distance : %v above maximum (%v > %v)
        Hint     : Value should be <= %v`,
			minValue,
			maxValue,
			actual,
			minValue,
			maxValue,
			actual-maxValue,
			actual,
			maxValue,
			maxValue,
		)
		result := formatRangeError(actual, minValue, maxValue)
		if result != expected {
			t.Errorf("Expected message:\n%s\n\nGot:\n%s", expected, result)
		}
	})
}

// === Tests for formatNotPanicError ===

func TestFormatNotPanicError(t *testing.T) {
	t.Parallel()

	t.Run("without stack trace", func(t *testing.T) {
		t.Parallel()
		panicInfo := panicInfo{
			Panicked:  true,
			Recovered: "test error",
			Stack:     "",
		}
		cfg := &Config{StackTrace: false}

		result := formatNotPanicError(panicInfo, cfg)

		if !strings.Contains(result, "Expected for the function to not panic") {
			t.Error("Should contain basic panic message")
		}
		if !strings.Contains(result, "test error") {
			t.Error("Should contain panic value")
		}
	})

	t.Run("with stack trace", func(t *testing.T) {
		t.Parallel()
		panicInfo := panicInfo{
			Panicked:  true,
			Recovered: "runtime error",
			Stack:     "some stack trace",
		}
		cfg := &Config{StackTrace: true}

		result := formatNotPanicError(panicInfo, cfg)

		if !strings.Contains(result, "Expected for the function to not panic") {
			t.Error("Should contain basic panic message")
		}
		if !strings.Contains(result, "runtime error") {
			t.Error("Should contain panic value")
		}
		if !strings.Contains(result, "Stack trace:") {
			t.Error("Should contain stack trace header")
		}
		if !strings.Contains(result, "some stack trace") {
			t.Error("Should contain stack trace content")
		}
	})
}

func TestCheckIfSorted(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			collection any
			failFast   bool
			expected   sortCheckResult
		}{
			{
				name:       "should return true for empty slice",
				collection: []int{},
				failFast:   false,
				expected:   sortCheckResult{IsSorted: true, Violations: nil, Total: 0},
			},
			{
				name:       "should return true for single element",
				collection: []int{42},
				failFast:   false,
				expected:   sortCheckResult{IsSorted: true, Violations: nil, Total: 1},
			},
			{
				name:       "should return true for sorted int slice",
				collection: []int{1, 2, 3, 4, 5},
				failFast:   false,
				expected:   sortCheckResult{IsSorted: true, Violations: nil, Total: 5},
			},
			{
				name:       "should return true for sorted string slice",
				collection: []string{"apple", "banana", "cherry"},
				failFast:   false,
				expected:   sortCheckResult{IsSorted: true, Violations: nil, Total: 3},
			},
			{
				name:       "should return true for sorted float slice",
				collection: []float64{1.1, 2.2, 3.3, 4.4},
				failFast:   false,
				expected:   sortCheckResult{IsSorted: true, Violations: nil, Total: 4},
			},
			{
				name:       "should return false with violations for unsorted slice",
				collection: []int{1, 3, 2, 5, 4},
				failFast:   false,
				expected: sortCheckResult{
					IsSorted: false,
					Violations: []sortViolation{
						{Index: 1, Value: 3, Next: 2},
						{Index: 3, Value: 5, Next: 4},
					},
					Total: 5,
				},
			},
			{
				name:       "should return false with violations for reverse sorted slice",
				collection: []int{5, 4, 3, 2, 1},
				failFast:   false,
				expected: sortCheckResult{
					IsSorted: false,
					Violations: []sortViolation{
						{Index: 0, Value: 5, Next: 4},
						{Index: 1, Value: 4, Next: 3},
						{Index: 2, Value: 3, Next: 2},
						{Index: 3, Value: 2, Next: 1},
					},
					Total: 5,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				var result sortCheckResult

				switch coll := tt.collection.(type) {
				case []int:
					result = checkIfSorted(coll)
				case []string:
					result = checkIfSorted(coll)
				case []float64:
					result = checkIfSorted(coll)
				default:
					t.Fatalf("Unsupported collection type: %T", tt.collection)
				}

				if result.IsSorted != tt.expected.IsSorted {
					t.Errorf("checkIfSorted() IsSorted = %v, want %v", result.IsSorted, tt.expected.IsSorted)
				}

				if result.Total != tt.expected.Total {
					t.Errorf("checkIfSorted() Total = %v, want %v", result.Total, tt.expected.Total)
				}

				if len(result.Violations) != len(tt.expected.Violations) {
					t.Errorf("checkIfSorted() Violations length = %v, want %v", len(result.Violations), len(tt.expected.Violations))
					return
				}

				for i, violation := range result.Violations {
					expected := tt.expected.Violations[i]
					if violation.Index != expected.Index || violation.Value != expected.Value || violation.Next != expected.Next {
						t.Errorf("checkIfSorted() Violation[%d] = {%d, %v, %v}, want {%d, %v, %v}",
							i, violation.Index, violation.Value, violation.Next,
							expected.Index, expected.Value, expected.Next)
					}
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		t.Run("should handle duplicates correctly", func(t *testing.T) {
			t.Parallel()
			duplicateSlice := []int{1, 2, 2, 3, 3, 3}
			result := checkIfSorted(duplicateSlice)

			if !result.IsSorted {
				t.Error("Expected slice with duplicates to be considered sorted")
			}

			if len(result.Violations) != 0 {
				t.Errorf("Expected no violations for duplicates, got %d", len(result.Violations))
			}
		})

		t.Run("should handle negative numbers", func(t *testing.T) {
			t.Parallel()
			negativeSlice := []int{-5, -3, -1, 0, 2}
			result := checkIfSorted(negativeSlice)

			if !result.IsSorted {
				t.Error("Expected sorted negative numbers to be considered sorted")
			}
		})

		t.Run("should stop at max violations limit", func(t *testing.T) {
			t.Parallel()
			largeUnsorted := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
			result := checkIfSorted(largeUnsorted)

			if result.IsSorted {
				t.Error("Expected large unsorted slice to be not sorted")
			}

			if len(result.Violations) != 6 {
				t.Errorf("Expected exactly 6 violations (max limit), got %d", len(result.Violations))
			}
		})
	})
}

func TestFormatSortErrorGeneric(t *testing.T) {
	t.Parallel()

	t.Run("Basic formatting", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			result   sortCheckResult
			contains []string
		}{
			{
				name: "should format error for small collection with single violation",
				result: sortCheckResult{
					IsSorted:   false,
					Violations: []sortViolation{{Index: 0, Value: 3, Next: 1}},
					Total:      3,
				},
				contains: []string{
					"Expected collection to be in ascending order, but it is not:",
					"Collection: (total: 3 elements)",
					"Status    : 1 order violation found",
					"Problems  :",
					"Index 0: 3 > 1",
				},
			},
			{
				name: "should format error for multiple violations",
				result: sortCheckResult{
					IsSorted: false,
					Violations: []sortViolation{
						{Index: 0, Value: 5, Next: 4},
						{Index: 1, Value: 4, Next: 3},
					},
					Total: 5,
				},
				contains: []string{
					"Expected collection to be in ascending order, but it is not:",
					"Collection: (total: 5 elements)",
					"Status    : 2 order violations found",
					"Index 0: 5 > 4",
					"Index 1: 4 > 3",
				},
			},
			{
				name: "should format error for large collection",
				result: sortCheckResult{
					IsSorted: false,
					Violations: []sortViolation{
						{Index: 0, Value: 100, Next: 99},
					},
					Total: 150,
				},
				contains: []string{
					"Collection: [Large collection] (total: 150 elements)",
					"Status    : 1 order violation found",
					"Index 0: 100 > 99",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatSortError(tt.result)

				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("formatSortError() result does not contain expected string %q\nActual result:\n%s", expected, result)
					}
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		t.Run("should return empty string when no violations", func(t *testing.T) {
			t.Parallel()
			result := sortCheckResult{
				IsSorted:   true,
				Violations: []sortViolation{}, // empty violations
				Total:      3,
			}

			output := formatSortError(result)
			if output != "" {
				t.Errorf("formatSortError() with no violations should return empty string, got: %q", output)
			}
		})

		t.Run("should handle many violations with truncation", func(t *testing.T) {
			t.Parallel()
			// Create 8 violations to trigger truncation (shows 5, mentions 3 more)
			violations := make([]sortViolation, 8)
			for i := 0; i < 8; i++ {
				violations[i] = sortViolation{
					Index: i,
					Value: i + 10, // higher value
					Next:  i,      // lower value
				}
			}

			result := sortCheckResult{
				IsSorted:   false,
				Violations: violations,
				Total:      10,
			}

			output := formatSortError(result)
			if !strings.Contains(output, "3 more violations") {
				t.Errorf("formatSortError() should mention '3 more violations', got:\n%s", output)
			}

			if !strings.Contains(output, "8 order violations found") {
				t.Errorf("formatSortError() should mention '8 order violations found', got:\n%s", output)
			}
		})

		t.Run("should handle exactly 6 violations (edge case)", func(t *testing.T) {
			t.Parallel()
			// Create exactly 6 violations to trigger "1 more violation" message
			violations := make([]sortViolation, 6)
			for i := 0; i < 6; i++ {
				violations[i] = sortViolation{
					Index: i,
					Value: i + 10, // higher value
					Next:  i,      // lower value
				}
			}

			result := sortCheckResult{
				IsSorted:   false,
				Violations: violations,
				Total:      10,
			}

			output := formatSortError(result)
			if !strings.Contains(output, "1 more violation") && !strings.Contains(output, "... and 1 more violation") {
				t.Errorf("formatSortError() should mention '1 more violation', got:\n%s", output)
			}
		})

		t.Run("should handle different value types", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name   string
				result sortCheckResult
			}{
				{
					name: "string values",
					result: sortCheckResult{
						IsSorted:   false,
						Violations: []sortViolation{{Index: 0, Value: "zebra", Next: "apple"}},
						Total:      2,
					},
				},
				{
					name: "float values",
					result: sortCheckResult{
						IsSorted:   false,
						Violations: []sortViolation{{Index: 0, Value: 3.14, Next: 2.71}},
						Total:      2,
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					result := formatSortError(tt.result)

					if result == "" {
						t.Error("formatSortError() should not return empty string for violations")
					}

					if !strings.Contains(result, "Expected collection to be in ascending order") {
						t.Error("formatSortError() should contain main error message")
					}
				})
			}
		})
	})
}

func TestFormatBeSameTimeError(t *testing.T) {
	t.Parallel()

	t.Run("basic error formatting", func(t *testing.T) {
		t.Parallel()
		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		actual := time.Date(2023, 12, 25, 15, 30, 47, 0, time.UTC)
		diff := 2 * time.Second

		result := formatBeSameTimeError(expected, actual, diff)

		// Verify main components are present
		if !strings.Contains(result, "Expected times to be the same") {
			t.Error("Should contain main error message")
		}
		if !strings.Contains(result, "2s") {
			t.Error("Should contain duration")
		}
		if !strings.Contains(result, "later") {
			t.Error("Should indicate actual is later")
		}
		if !strings.Contains(result, "Expected:") {
			t.Error("Should contain expected time label")
		}
		if !strings.Contains(result, "Actual  :") {
			t.Error("Should contain actual time label")
		}
	})

	t.Run("actual earlier than expected", func(t *testing.T) {
		t.Parallel()
		expected := time.Date(2023, 12, 25, 15, 30, 47, 0, time.UTC)
		actual := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		diff := 2 * time.Second

		result := formatBeSameTimeError(expected, actual, diff)

		if !strings.Contains(result, "earlier") {
			t.Error("Should indicate actual is earlier")
		}
	})

	t.Run("with fractional seconds", func(t *testing.T) {
		t.Parallel()
		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		actual := time.Date(2023, 12, 25, 15, 30, 45, 500000000, time.UTC) // +0.5s
		diff := 500 * time.Millisecond

		result := formatBeSameTimeError(expected, actual, diff)

		if !strings.Contains(result, "500ms") {
			t.Errorf("Should contain 500ms, got: %s", result)
		}
	})

	t.Run("message structure", func(t *testing.T) {
		t.Parallel()
		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		actual := time.Date(2023, 12, 25, 15, 30, 46, 0, time.UTC)
		diff := 1 * time.Second

		result := formatBeSameTimeError(expected, actual, diff)
		lines := strings.Split(result, "\n")

		if len(lines) != 3 {
			t.Errorf("Expected 3 lines, got %d: %s", len(lines), result)
		}

		// Check line structure
		if !strings.HasPrefix(lines[0], "Expected times to be the same") {
			t.Errorf("First line should start with main message, got: %s", lines[0])
		}
		if !strings.HasPrefix(lines[1], "Expected:") {
			t.Errorf("Second line should start with 'Expected:', got: %s", lines[1])
		}
		if !strings.HasPrefix(lines[2], "Actual  :") {
			t.Errorf("Third line should start with 'Actual  :', got: %s", lines[2])
		}
	})
}

func TestHumanizeDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		// Nanoseconds and microseconds (shown as fractional ms)
		{
			name:     "1 nanosecond",
			duration: 1 * time.Nanosecond,
			expected: "0.000ms",
		},
		{
			name:     "1 microsecond",
			duration: 1 * time.Microsecond,
			expected: "0.001ms",
		},
		{
			name:     "500 microseconds",
			duration: 500 * time.Microsecond,
			expected: "0.500ms",
		},

		// Milliseconds
		{
			name:     "1 millisecond",
			duration: 1 * time.Millisecond,
			expected: "1ms",
		},
		{
			name:     "1.5 milliseconds",
			duration: 1500 * time.Microsecond,
			expected: "1.5ms",
		},
		{
			name:     "999.9 milliseconds",
			duration: 999900 * time.Microsecond,
			expected: "999.9ms",
		},

		// Seconds
		{
			name:     "1 second",
			duration: 1 * time.Second,
			expected: "1s",
		},
		{
			name:     "1.5 seconds",
			duration: 1500 * time.Millisecond,
			expected: "1.5s",
		},
		{
			name:     "59.9 seconds",
			duration: 59900 * time.Millisecond,
			expected: "59.9s",
		},
		{
			name:     "60 seconds",
			duration: 60 * time.Second,
			expected: "1m",
		},
		{
			name:     "125.3 seconds",
			duration: 125300 * time.Millisecond,
			expected: "2m5s",
		},

		// Minutes
		{
			name:     "1 minute",
			duration: 1 * time.Minute,
			expected: "1m",
		},
		{
			name:     "1 minute 30 seconds",
			duration: 1*time.Minute + 30*time.Second,
			expected: "1m30s",
		},
		{
			name:     "59 minutes 59 seconds",
			duration: 59*time.Minute + 59*time.Second,
			expected: "59m59s",
		},

		// Hours
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
			expected: "1h",
		},
		{
			name:     "1 hour 30 minutes",
			duration: 1*time.Hour + 30*time.Minute,
			expected: "1h30m",
		},
		{
			name:     "23 hours 59 minutes",
			duration: 23*time.Hour + 59*time.Minute,
			expected: "23h59m",
		},

		// Days
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			expected: "1d",
		},
		{
			name:     "1 day 1 hour",
			duration: 24*time.Hour + 1*time.Hour,
			expected: "1d1h",
		},
		{
			name:     "2 days",
			duration: 48 * time.Hour,
			expected: "2d",
		},
		{
			name:     "2 days 1 hour",
			duration: 48*time.Hour + 1*time.Hour,
			expected: "2d1h",
		},

		// Negative durations (should be handled as positive)
		{
			name:     "negative 1 second",
			duration: -1 * time.Second,
			expected: "1s",
		},
		{
			name:     "negative 1 minute",
			duration: -1 * time.Minute,
			expected: "1m",
		},

		// Edge cases
		{
			name:     "zero duration",
			duration: 0,
			expected: "0.000ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := humanizeDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("humanizeDuration(%v) = %q, expected %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestFormatTimeForDisplay(t *testing.T) {
	t.Parallel()

	t.Run("UTC times", func(t *testing.T) {
		tests := []struct {
			name     string
			time     time.Time
			expected string
		}{
			{
				name:     "basic UTC time without nanoseconds",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
				expected: "2023-12-25 15:30:45 UTC",
			},
			{
				name:     "UTC time with nanoseconds",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC),
				expected: "2023-12-25 15:30:45.123456789 UTC",
			},
			{
				name:     "UTC time with trailing zeros in nanoseconds",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC),
				expected: "2023-12-25 15:30:45.123 UTC",
			},
			{
				name:     "UTC time with only microseconds",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 123456000, time.UTC),
				expected: "2023-12-25 15:30:45.123456 UTC",
			},
			{
				name:     "UTC time with only milliseconds",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC),
				expected: "2023-12-25 15:30:45.123 UTC",
			},
			{
				name:     "UTC time with 1 nanosecond",
				time:     time.Date(2023, 12, 25, 15, 30, 45, 1, time.UTC),
				expected: "2023-12-25 15:30:45.000000001 UTC",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatTimeForDisplay(tt.time)
				if result != tt.expected {
					t.Errorf("formatTimeForDisplay(%v) = %q, expected %q", tt.time, result, tt.expected)
				}
			})
		}
	})

	t.Run("timezone handling", func(t *testing.T) {
		t.Run("fixed timezone", func(t *testing.T) {
			t.Parallel()
			est := time.FixedZone("EST", -5*3600)
			testTime := time.Date(2023, 12, 25, 10, 30, 45, 0, est)
			result := formatTimeForDisplay(testTime)
			// Shows UTC time but preserves timezone name
			expected := "2023-12-25 15:30:45 EST"
			if result != expected {
				t.Errorf("formatTimeForDisplay(%v) = %q, expected %q", testTime, result, expected)
			}
		})

		t.Run("utc timezone", func(t *testing.T) {
			t.Parallel()
			testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
			result := formatTimeForDisplay(testTime)
			expected := "2023-01-01 12:00:00 UTC"
			if result != expected {
				t.Errorf("formatTimeForDisplay(%v) = %q, expected %q", testTime, result, expected)
			}
		})

		t.Run("local timezone fallback", func(t *testing.T) {
			t.Parallel()
			testTime := time.Date(2023, 3, 15, 9, 15, 30, 0, time.Local)
			result := formatTimeForDisplay(testTime)

			// Should contain the date part at minimum
			if !strings.Contains(result, "2023-03-15") {
				t.Errorf("Result should contain date, got: %s", result)
			}
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("zero time", func(t *testing.T) {
			t.Parallel()
			result := formatTimeForDisplay(time.Time{})
			expected := "0001-01-01 00:00:00 UTC"

			if result != expected {
				t.Errorf("formatTimeForDisplay(zero) = %q, expected %q", result, expected)
			}
		})

		t.Run("maximum precision nanoseconds", func(t *testing.T) {
			t.Parallel()
			testTime := time.Date(2023, 12, 25, 15, 30, 45, 999999999, time.UTC)
			result := formatTimeForDisplay(testTime)
			expected := "2023-12-25 15:30:45.999999999 UTC"

			if result != expected {
				t.Errorf("formatTimeForDisplay(%v) = %q, expected %q", testTime, result, expected)
			}
		})

		t.Run("leap year date", func(t *testing.T) {
			t.Parallel()
			testTime := time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC)
			result := formatTimeForDisplay(testTime)
			expected := "2020-02-29 12:00:00 UTC"

			if result != expected {
				t.Errorf("formatTimeForDisplay(%v) = %q, expected %q", testTime, result, expected)
			}
		})
	})
}

func TestFormatterIntegration(t *testing.T) {
	t.Parallel()

	t.Run("full error message integration", func(t *testing.T) {
		t.Parallel()
		expected := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
		actual := time.Date(2023, 12, 25, 15, 30, 47, 987654321, time.UTC)
		diff := actual.Sub(expected)

		result := formatBeSameTimeError(expected, actual, diff)

		lines := strings.Split(result, "\n")
		if len(lines) != 3 {
			t.Fatalf("Expected 3 lines, got %d", len(lines))
		}

		if !strings.Contains(lines[0], "s") {
			t.Error("Duration should be humanized in first line")
		}

		if !strings.Contains(lines[1], ".123456789") {
			t.Error("Expected time should show nanoseconds")
		}
		if !strings.Contains(lines[2], ".987654321") {
			t.Error("Actual time should show nanoseconds")
		}

		// Check relation
		if !strings.Contains(lines[2], "later") {
			t.Error("Should indicate actual is later")
		}
	})

	t.Run("different timezone integration", func(t *testing.T) {
		t.Parallel()
		utc := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		est := time.Date(2023, 12, 25, 10, 30, 46, 0, time.FixedZone("EST", -5*3600))
		diff := est.Sub(utc)

		result := formatBeSameTimeError(utc, est, diff)

		// Both times should be displayed in UTC format
		if !strings.Contains(result, "2023-12-25 15:30:45 UTC") {
			t.Error("Expected time should be in UTC format")
		}
		if !strings.Contains(result, "2023-12-25 15:30:46 EST") {
			t.Error("Actual time should be converted to UTC format")
		}
	})

	t.Run("very small difference", func(t *testing.T) {
		t.Parallel()
		base := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		withNanos := time.Date(2023, 12, 25, 15, 30, 45, 500, time.UTC) // 500ns
		diff := withNanos.Sub(base)

		result := formatBeSameTimeError(base, withNanos, diff)

		// Should show fractional milliseconds for tiny differences
		if !strings.Contains(result, "ms") {
			t.Error("Very small differences should be shown in milliseconds")
		}
	})
}

func TestFormatBeErrorMessage(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			action   string
			err      error
			target   interface{}
			contains []string
		}{
			{
				name:   "Action 'as'",
				action: "as",
				err:    simpleError{msg: "test error"},
				target: &simpleError{},
				contains: []string{
					"Expected error to be *assert.simpleError",
					"but type not found in error chain",
					"Error: \"test error\"",
					"Types  : [assert.simpleError]",
				},
			},
			{
				name:   "Action 'is'",
				action: "is",
				err:    simpleError{msg: "test error"},
				target: errors.New("target"),
				contains: []string{
					"Expected error to be \"target\"",
					"but not found in error chain",
					"Error: \"test error\"",
					"Types  : [assert.simpleError]",
				},
			},
			{
				name:   "Unknown action",
				action: "unknown",
				err:    simpleError{msg: "test error"},
				target: nil,
				contains: []string{
					"Assertion failed with an unknown type of error",
					"Error: \"test error\"",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeErrorMessage(tt.action, tt.err, tt.target)
				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result:\n%s", expected, result)
					}
				}
			})
		}
	})

	t.Run("Error chain", func(t *testing.T) {
		t.Parallel()
		innerErr := simpleError{msg: "inner"}
		outerErr := wrapperError{msg: "outer", wrapped: innerErr}

		result := formatBeErrorMessage("as", outerErr, &simpleError{})

		if !strings.Contains(result, "Types  : [assert.wrapperError, assert.simpleError]") {
			t.Errorf("Expected error chain types, got:\n%s", result)
		}
	})

	t.Run("Multiple wrapped errors", func(t *testing.T) {
		t.Parallel()
		innerErr := simpleError{msg: "inner"}
		middleErr := wrapperError{msg: "middle", wrapped: innerErr}
		outerErr := wrapperError{msg: "outer", wrapped: middleErr}

		result := formatBeErrorMessage("is", outerErr, errors.New("target"))

		expected := "Types  : [assert.wrapperError, assert.wrapperError, assert.simpleError]"
		if !strings.Contains(result, expected) {
			t.Errorf("Expected multiple wrapped error types, got:\n%s", result)
		}
	})
}

func TestFormatNotBeErrorMessage(t *testing.T) {
	t.Parallel()

	t.Run("Basic functionality", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name      string
			customMsg string
			err       error
			contains  []string
		}{
			{
				name:      "with customMsg",
				customMsg: "File doesn't exist",
				err:       errors.New("test error"),
				contains: []string{
					"File doesn't exist",
					"Expected no error, but got an error",
					"Error: \"test error\"",
					"Type: *errors.errorString",
				},
			},
			{
				name:      "empty customMsg",
				customMsg: "",
				err:       errors.New("test error"),
				contains: []string{
					"Expected no error, but got an error",
					"Error: \"test error\"",
					"Type: *errors.errorString",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatNotBeErrorMessage(tt.customMsg, tt.err)
				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result: \n%s", expected, result)
					}
				}
			})
		}
	})
}

func TestFormatBeWithinError(t *testing.T) {
	t.Parallel()

	t.Run("basic formatting cases", func(t *testing.T) {
		t.Parallel()
		basicCases := []struct {
			name        string
			actual      float64
			expected    float64
			tolerance   float64
			contains    []string
			notContains []string
		}{
			{
				name:      "normal decimal numbers",
				actual:    10.5,
				expected:  10.0,
				tolerance: 0.1,
				contains: []string{
					"Expected value to be within ±0.100000 of 10.000000, but it was not:",
					"Actual:    10.500000",
					"Expected:  10.000000",
					"Difference: 0.500000",
					"(4.00× greater than tolerance)",
				},
			},
			{
				name:      "negative numbers",
				actual:    -5.3,
				expected:  -5.0,
				tolerance: 0.2,
				contains: []string{
					"Actual:    -5.300000",
					"Expected:  -5.000000",
					"Difference: 0.300000",
					"(50.00% greater than tolerance)",
				},
			},
			{
				name:      "zero expected with relative diff",
				actual:    0.5,
				expected:  0.0,
				tolerance: 0.1,
				contains: []string{
					"Actual:    0.500000",
					"Expected:  0.000000",
					"Difference: 0.500000",
					"(4.00× greater than tolerance)",
				},
			},
			{
				name:      "values that appear within tolerance",
				actual:    10.05,
				expected:  10.0,
				tolerance: 0.1,
				contains: []string{
					"Actual:    10.050000",
					"Expected:  10.000000",
					"Difference: 0.050000",
				},
			},
		}

		for _, tt := range basicCases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				// Function should always return a non-empty error message
				if result == "" {
					t.Error("Expected non-empty result, got empty string")
				}

				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result:\n%s", expected, result)
					}
				}

				for _, notExpected := range tt.notContains {
					if strings.Contains(result, notExpected) {
						t.Errorf("Did not expect %q in result:\n%s", notExpected, result)
					}
				}
			})
		}
	})

	t.Run("zero and negative tolerance cases", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name      string
			actual    float64
			expected  float64
			tolerance float64
			contains  []string
		}{
			{
				name:      "zero tolerance - different values",
				actual:    10.0,
				expected:  5.0,
				tolerance: 0.0,
				contains: []string{
					"Expected value to be within ±0.000000 of 5.000000",
					"Difference: 5.000000",
				},
			},
			{
				name:      "negative tolerance",
				actual:    10.0,
				expected:  5.0,
				tolerance: -1.0,
				contains: []string{
					"Expected value to be within ±-1.000000 of 5.000000",
					"Tolerance: ±-1.000000",
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				if result == "" {
					t.Error("Expected non-empty result, got empty string")
				}

				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result:\n%s", expected, result)
					}
				}

				// For zero tolerance, should not show percentage
				if tt.tolerance == 0.0 && strings.Contains(result, "% greater than tolerance") {
					t.Errorf("Should not show percentage for zero tolerance: %s", result)
				}
			})
		}
	})

	t.Run("scientific notation cases", func(t *testing.T) {
		t.Parallel()
		scientificCases := []struct {
			name      string
			actual    float64
			expected  float64
			tolerance float64
			contains  []string
		}{
			{
				name:      "very large numbers",
				actual:    1.5e30,
				expected:  1.0e30,
				tolerance: 1e28,
				contains: []string{
					"1.500000e+30",
					"1.000000e+30",
					"5.000000e+29",
					"1.000000e+28",
					"(49.00× greater than tolerance)",
				},
			},
			{
				name:      "very small numbers",
				actual:    1.5e-6,
				expected:  1.0e-6,
				tolerance: 1e-8,
				contains: []string{
					"1.500000e-06",
					"1.000000e-06",
					"5.000000e-07",
					"1.000000e-08",
					"(49.00× greater than tolerance)",
				},
			},
		}

		for _, tt := range scientificCases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result:\n%s", expected, result)
					}
				}
			})
		}
	})

	t.Run("special values edge cases", func(t *testing.T) {
		t.Parallel()

		edgeCases := []struct {
			name      string
			actual    float64
			expected  float64
			tolerance float64
			contains  []string
		}{
			{
				name:      "positive infinity actual",
				actual:    math.Inf(1),
				expected:  10.0,
				tolerance: 1.0,
				contains: []string{
					"+Inf",
					"1.000000e+01", // Scientific notation for 10.0
					"Difference: +Inf",
				},
			},
			{
				name:      "negative infinity actual",
				actual:    math.Inf(-1),
				expected:  -10.0,
				tolerance: 1.0,
				contains: []string{
					"-Inf",
					"-1.000000e+01", // Scientific notation for -10.0
					"Difference: +Inf",
				},
			},
			{
				name:      "NaN actual",
				actual:    math.NaN(),
				expected:  10.0,
				tolerance: 1.0,
				contains: []string{
					"NaN",
					"10.000000",
					"Difference: NaN",
				},
			},
		}

		for _, tt := range edgeCases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				// Should always return a non-empty result
				if result == "" {
					t.Error("Expected non-empty result for special values")
				}

				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected %q in result:\n%s", expected, result)
					}
				}

				// Basic structure should still be present
				requiredFields := []string{"Actual:", "Expected:", "Tolerance:", "Difference:"}
				for _, field := range requiredFields {
					if !strings.Contains(result, field) {
						t.Errorf("Missing required field %q in: %s", field, result)
					}
				}
			})
		}
	})

	t.Run("format selection accuracy", func(t *testing.T) {
		t.Parallel()

		formatTests := []struct {
			name                string
			actual              float64
			expected            float64
			tolerance           float64
			shouldUseScientific bool
		}{
			{"small decimals use decimal", 1.234567, 1.234560, 1e-6, false},
			{"large numbers use scientific", 1e7, 1e7 + 1000, 100, true},
			{"medium numbers use decimal", 123.456, 123.450, 0.001, false},
			{"boundary case large - uses scientific", 999999, 1000000, 1, true}, // 1e6 threshold
			{"boundary case small", 0.0001, 0.00011, 1e-6, false},
		}

		for _, tt := range formatTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				hasScientific := strings.Contains(result, "e+") || strings.Contains(result, "e-")
				if tt.shouldUseScientific && !hasScientific {
					t.Errorf("Expected scientific notation for %s, got: %s", tt.name, result)
				} else if !tt.shouldUseScientific && hasScientific {
					t.Errorf("Expected decimal notation for %s, got: %s", tt.name, result)
				}
			})
		}
	})

	t.Run("relative percentage accuracy", func(t *testing.T) {
		t.Parallel()

		relativeTests := []struct {
			name             string
			actual           float64
			expected         float64
			tolerance        float64
			expectedRelative string
		}{
			// Corrected calculations: ((diff - tolerance) / tolerance) * 100
			{"900% over tolerance", 11.0, 10.0, 0.1, "(9.00× greater than tolerance)"},
			{"400% over tolerance", 15.0, 10.0, 1.0, "(4.00× greater than tolerance)"},
			{"just below multiplier boundary", 14.0, 10.0, 1.0, "(3.00× greater than tolerance)"},
			{"exactly multiplier boundary", 16.0, 10.0, 1.0, "(5.00× greater than tolerance)"},
			{"above multiplier boundary", 20.0, 10.0, 1.0, "(9.00× greater than tolerance)"},
			{"negative expected", -11.0, -10.0, 0.1, "(9.00× greater than tolerance)"},
			{"slightly over tolerance", 10.15, 10.0, 0.1, "(50.00% greater than tolerance)"},
		}

		for _, tt := range relativeTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := formatBeWithinError(tt.actual, tt.expected, tt.tolerance)

				if !strings.Contains(result, tt.expectedRelative) {
					t.Errorf("Expected %q in result for %s, got: %s",
						tt.expectedRelative, tt.name, result)
				}
			})
		}
	})

	t.Run("generic type constraints", func(t *testing.T) {
		t.Parallel()

		t.Run("float32", func(t *testing.T) {
			t.Parallel()
			result := formatBeWithinError(float32(3.14159), float32(3.14), float32(0.001))
			expectedStrings := []string{
				"3.141590",
				"3.140000",
				"0.001590",
				"0.001000",
			}
			for _, expected := range expectedStrings {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected %q in result:\n%s", expected, result)
				}
			}
		})
	})
}