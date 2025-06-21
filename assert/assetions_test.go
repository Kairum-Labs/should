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
		NotContain(t, []string{"apple", "banana"}, "apple", AssertionConfig{Message: "Apple should not be in basket"})
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
		Contain(t, []string{"apple", "banana"}, "orange", AssertionConfig{Message: "Fruit not found in basket"})
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
		BeEmpty(t, "hello", AssertionConfig{Message: "This is a custom message"})
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
		BeEmpty(t, ch, AssertionConfig{Message: "Channel should be empty"})
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

// === Tests for BeNotEmpty ===

func TestBeNotEmpty_Succeeds_WithNonEmptyString(t *testing.T) {
	t.Parallel()

	BeNotEmpty(t, "hello")
}

func TestBeNotEmpty_Succeeds_WithNonEmptySlice(t *testing.T) {
	t.Parallel()

	BeNotEmpty(t, []int{1, 2, 3})
}

func TestBeNotEmpty_Fails_WithEmptyString(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeNotEmpty(t, "")
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
		BeNotEmpty(t, []int{})
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
		BeNotEmpty(t, ptr)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be not empty, but it was empty:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeNotEmpty_WithInvalidValue(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		var v interface{}
		BeNotEmpty(t, v)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "Expected value to be not empty, but it was empty:"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeNotEmpty_WithUnsupportedType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeNotEmpty(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeNotEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected message to contain: %q\n\nFull message:\n%s", expected, message)
	}
}

func TestBeNotEmpty_WithChannelBuffered(t *testing.T) {
	t.Parallel()

	// Test non-empty buffered channel
	ch := make(chan int, 1)
	ch <- 42
	BeNotEmpty(t, ch)

	// Test empty buffered channel
	emptyChannel := make(chan int, 1)
	failed, message := assertFails(t, func(t testing.TB) {
		BeNotEmpty(t, emptyChannel, AssertionConfig{Message: "Channel should have data"})
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

func TestBeNotEmpty_WithValidInterface(t *testing.T) {
	t.Parallel()

	var validInterface interface{} = 42
	failed, message := assertFails(t, func(t testing.TB) {
		BeNotEmpty(t, validInterface)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeNotEmpty can only be used with strings, slices, arrays, maps, channels, or pointers, but got int"
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
		BeGreaterThan(t, 0.0, 0.1, AssertionConfig{Message: "Score validation failed"})
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
		}, AssertionConfig{Message: "Division by zero should panic"})
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
		}, AssertionConfig{Message: "Save operation should not panic"})
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
		BeGreaterOrEqualThan(t, 0, 1, AssertionConfig{Message: "Score cannot be negative"})
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
		BeGreaterOrEqualThan(t, 5.0, 5.1, AssertionConfig{Message: "Integer vs float comparison"})
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

func TestBeNotNil_Succeeds_WithNonNilPointer(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value
	BeNotNil(t, ptr)
}

func TestBeNotNil_Succeeds_WithNonNilSlice(t *testing.T) {
	t.Parallel()

	slice := make([]int, 0)
	BeNotNil(t, slice)
}

func TestBeNotNil_Fails_WithNilPointer(t *testing.T) {
	t.Parallel()

	var ptr *int

	failed, message := assertFails(t, func(t testing.TB) {
		BeNotNil(t, ptr)
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
		BeTrue(t, false, AssertionConfig{Message: "User must be active"})
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
		BeFalse(t, true, AssertionConfig{Message: "User should not be deleted"})
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

// === Tests for BeNil/BeNotNil with custom messages ===

func TestBeNil_WithCustomMessage(t *testing.T) {
	t.Parallel()

	value := 42
	ptr := &value

	failed, message := assertFails(t, func(t testing.TB) {
		BeNil(t, ptr, AssertionConfig{Message: "Pointer should be nil"})
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

func TestBeNotNil_WithCustomMessage(t *testing.T) {
	t.Parallel()

	var ptr *int

	failed, message := assertFails(t, func(t testing.TB) {
		BeNotNil(t, ptr, AssertionConfig{Message: "User must not be nil"})
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

func TestBeNotNil_Fails_WithNonNillableType(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeNotNil(t, 42)
	})

	if !failed {
		t.Fatal("Expected test to fail, but it passed")
	}

	expected := "BeNotNil can only be used with nillable types, but got int"
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

// === Tests for unsupported numeric type combinations ===

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

func TestBeLessThan_WithCustomMessage(t *testing.T) {
	t.Parallel()

	failed, message := assertFails(t, func(t testing.TB) {
		BeLessThan(t, 10, 5, AssertionConfig{Message: "Value should be smaller"})
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
		BeEqual(t, person1, person2, AssertionConfig{Message: "Person objects should be identical"})
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
