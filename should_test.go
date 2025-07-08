package should

import (
	"fmt"
	"strings"
	"testing"
)

type mockTB struct {
	testing.TB
	failed      bool
	lastMessage string
}

func (m *mockTB) Helper() {}

func (m *mockTB) Errorf(format string, args ...any) {
	m.failed = true
	m.lastMessage = fmt.Sprintf(format, args...)
}

func (m *mockTB) Error(args ...any) {
	m.failed = true
	m.lastMessage = fmt.Sprint(args...)
}

func (m *mockTB) FailNow() {
	m.failed = true
	panic("FailNow called")
}

func TestAssertions(t *testing.T) {
	t.Run("BeEqual should pass for equal values", func(t *testing.T) {
		mockT := &mockTB{}
		BeEqual(mockT, "some value", "some value")
		if mockT.failed {
			t.Errorf("Expected BeEqual to pass, but it failed with message: %q", mockT.lastMessage)
		}
	})

	t.Run("BeEqual should fail for unequal values", func(t *testing.T) {
		mockT := &mockTB{}
		BeEqual(mockT, "some value", "another value")
		if !mockT.failed {
			t.Error("Expected BeEqual to fail, but it passed")
		}
	})
}

func TestPanic(t *testing.T) {
	testCases := []struct {
		name        string
		fn          func()
		opts        []Option
		shouldFail  bool
		expectedMsg string
	}{
		{
			name:       "should pass when function panics",
			fn:         func() { panic("expected panic") },
			shouldFail: false,
		},
		{
			name:        "should fail when function does not panic",
			fn:          func() {},
			shouldFail:  true,
			expectedMsg: "Expected panic, but did not panic",
		},
		{
			name: "should fail with custom message when function does not panic",
			fn:   func() {},
			opts: []Option{
				WithMessage("custom message"),
			},
			shouldFail:  true,
			expectedMsg: "custom message\nExpected panic, but did not panic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTB{}
			Panic(mockT, tc.fn, tc.opts...)

			if tc.shouldFail != mockT.failed {
				t.Errorf("Expected test failure to be %v, but was %v", tc.shouldFail, mockT.failed)
			}

			if tc.shouldFail && !strings.Contains(mockT.lastMessage, tc.expectedMsg) {
				t.Errorf("Expected error message to contain %q, but got %q", tc.expectedMsg, mockT.lastMessage)
			}
		})
	}
}

func TestNotPanic(t *testing.T) {
	testCases := []struct {
		name        string
		fn          func()
		opts        []Option
		shouldFail  bool
		expectedMsg string
	}{
		{
			name:       "should pass when function does not panic",
			fn:         func() {},
			shouldFail: false,
		},
		{
			name:        "should fail when function panics",
			fn:          func() { panic("some panic") },
			shouldFail:  true,
			expectedMsg: "Expected for the function to not panic, but it panicked with: some panic",
		},
		{
			name: "should fail with custom message when function panics",
			fn:   func() { panic("some panic") },
			opts: []Option{
				WithMessage("custom message"),
			},
			shouldFail:  true,
			expectedMsg: "custom message\nExpected for the function to not panic, but it panicked with: some panic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTB{}
			NotPanic(mockT, tc.fn, tc.opts...)

			if tc.shouldFail != mockT.failed {
				t.Errorf("Expected test failure to be %v, but was %v", tc.shouldFail, mockT.failed)
			}

			if tc.shouldFail && !strings.Contains(mockT.lastMessage, tc.expectedMsg) {
				t.Errorf("Expected error message to contain %q, but got %q", tc.expectedMsg, mockT.lastMessage)
			}
		})
	}
}

func TestWrappers(t *testing.T) {
	// BeTrue
	t.Run("BeTrue passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeTrue(mockT, true)
		if mockT.failed {
			t.Error("BeTrue should pass")
		}
	})
	t.Run("BeTrue fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeTrue(mockT, false)
		if !mockT.failed {
			t.Error("BeTrue should fail")
		}
	})

	// BeFalse
	t.Run("BeFalse passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeFalse(mockT, false)
		if mockT.failed {
			t.Error("BeFalse should pass")
		}
	})
	t.Run("BeFalse fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeFalse(mockT, true)
		if !mockT.failed {
			t.Error("BeFalse should fail")
		}
	})

	// BeEmpty
	t.Run("BeEmpty passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeEmpty(mockT, "")
		if mockT.failed {
			t.Error("BeEmpty should pass for empty string")
		}
	})
	t.Run("BeEmpty fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeEmpty(mockT, "not empty")
		if !mockT.failed {
			t.Error("BeEmpty should fail for non-empty string")
		}
	})

	// NotBeEmpty
	t.Run("NotBeEmpty passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotBeEmpty(mockT, "not empty")
		if mockT.failed {
			t.Error("NotBeEmpty should pass")
		}
	})
	t.Run("NotBeEmpty fails", func(t *testing.T) {
		mockT := &mockTB{}
		NotBeEmpty(mockT, "")
		if !mockT.failed {
			t.Error("NotBeEmpty should fail")
		}
	})

	// BeNil
	t.Run("BeNil passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeNil(mockT, nil)
		if mockT.failed {
			t.Error("BeNil should pass")
		}
	})
	t.Run("BeNil fails", func(t *testing.T) {
		mockT := &mockTB{}
		var x = 1
		BeNil(mockT, &x)
		if !mockT.failed {
			t.Error("BeNil should fail")
		}
	})

	// NotBeNil
	t.Run("NotBeNil passes", func(t *testing.T) {
		mockT := &mockTB{}
		var x = 1
		NotBeNil(mockT, &x)
		if mockT.failed {
			t.Error("NotBeNil should pass")
		}
	})
	t.Run("NotBeNil fails", func(t *testing.T) {
		mockT := &mockTB{}
		NotBeNil(mockT, nil)
		if !mockT.failed {
			t.Error("NotBeNil should fail")
		}
	})

	t.Run("NotBeEqual passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotBeEqual(mockT, "a", "b")
		if mockT.failed {
			t.Error("NotBeEqual should pass")
		}
	})

	// BeGreaterThan
	t.Run("BeGreaterThan passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterThan(mockT, 10, 5)
		if mockT.failed {
			t.Error("BeGreaterThan should pass")
		}
	})
	t.Run("BeGreaterThan fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterThan(mockT, 5, 10)
		if !mockT.failed {
			t.Error("BeGreaterThan should fail")
		}
	})

	// BeLessThan
	t.Run("BeLessThan passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeLessThan(mockT, 5, 10)
		if mockT.failed {
			t.Error("BeLessThan should pass")
		}
	})
	t.Run("BeLessThan fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeLessThan(mockT, 10, 5)
		if !mockT.failed {
			t.Error("BeLessThan should fail")
		}
	})

	// BeGreaterOrEqualTo
	t.Run("BeGreaterOrEqualTo passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterOrEqualTo(mockT, 10, 10)
		if mockT.failed {
			t.Error("BeGreaterOrEqualTo should pass")
		}
	})
	t.Run("BeGreaterOrEqualTo fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterOrEqualTo(mockT, 9, 10)
		if !mockT.failed {
			t.Error("BeGreaterOrEqualTo should fail")
		}
	})

	t.Run("BeLessOrEqualTo passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeLessOrEqualTo(mockT, 10, 10)
		if mockT.failed {
			t.Error("BeLessOrEqualTo should pass")
		}
	})

	// Contain
	t.Run("Contain passes", func(t *testing.T) {
		mockT := &mockTB{}
		Contain(mockT, []int{1, 2, 3}, 2)
		if mockT.failed {
			t.Error("Contain should pass")
		}
	})
	t.Run("Contain fails", func(t *testing.T) {
		mockT := &mockTB{}
		Contain(mockT, []int{1, 2, 3}, 4)
		if !mockT.failed {
			t.Error("Contain should fail")
		}
	})

	// NotContain
	t.Run("NotContain passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotContain(mockT, []int{1, 2, 3}, 4)
		if mockT.failed {
			t.Error("NotContain should pass")
		}
	})

	t.Run("NotContain fails", func(t *testing.T) {
		mockT := &mockTB{}
		NotContain(mockT, []int{1, 2, 3}, 2)
		if !mockT.failed {
			t.Error("NotContain should fail")
		}
	})

	t.Run("NotContainDuplicates passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotContainDuplicates(mockT, []int{1, 2, 3})
		if mockT.failed {
			t.Error("NotContainDuplicates should pass")
		}
	})

	t.Run("NotContainKey passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotContainKey(mockT, map[string]int{"a": 1, "b": 2}, "c")
		if mockT.failed {
			t.Error("NotContainKey should pass")
		}
	})

	t.Run("NotContainValue passes", func(t *testing.T) {
		mockT := &mockTB{}
		NotContainValue(mockT, map[string]int{"a": 1, "b": 2}, 3)
		if mockT.failed {
			t.Error("NotContainValue should pass")
		}
	})

	t.Run("StartsWith passes", func(t *testing.T) {
		mockT := &mockTB{}
		StartsWith(mockT, "Hello, world!", "Hello")
		StartsWith(mockT, "Hello, world!", "hello", IgnoreCase())
		if mockT.failed {
			t.Error("StartsWith should pass")
		}
	})

	t.Run("EndsWith passes", func(t *testing.T) {
		mockT := &mockTB{}
		EndsWith(mockT, "Hello, world", "world")
		if mockT.failed {
			t.Error("EndsWith should pass")
		}
	})

	// ContainFunc
	t.Run("ContainFunc passes", func(t *testing.T) {
		mockT := &mockTB{}
		ContainFunc(mockT, []int{1, 2, 3}, func(item any) bool { return item.(int) == 2 })
		if mockT.failed {
			t.Error("ContainFunc should pass")
		}
	})
	t.Run("ContainFunc fails", func(t *testing.T) {
		mockT := &mockTB{}
		ContainFunc(mockT, []int{1, 2, 3}, func(item any) bool { return item.(int) == 4 })
		if !mockT.failed {
			t.Error("ContainFunc should fail")
		}
	})

	t.Run("ContainSubstring passes", func(t *testing.T) {
		mockT := &mockTB{}
		ContainSubstring(mockT, "Hello, world!", "world")
		if mockT.failed {
			t.Error("ContainSubstring should pass")
		}
	})

	// HaveLength
	t.Run("HaveLength passes", func(t *testing.T) {
		mockT := &mockTB{}
		HaveLength(mockT, []int{1, 2, 3}, 3)
		if mockT.failed {
			t.Error("HaveLength should pass")
		}
	})
	t.Run("HaveLength fails", func(t *testing.T) {
		mockT := &mockTB{}
		HaveLength(mockT, []int{1, 2, 3}, 4)
		if !mockT.failed {
			t.Error("HaveLength should fail")
		}
	})

	// BeOfType
	t.Run("BeOfType passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeOfType(mockT, "hello", "world")
		if mockT.failed {
			t.Error("BeOfType should pass")
		}
	})
	t.Run("BeOfType fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeOfType(mockT, "hello", 123)
		if !mockT.failed {
			t.Error("BeOfType should fail")
		}
	})

	// BeOneOf
	t.Run("BeOneOf passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeOneOf(mockT, "a", []string{"a", "b"})
		if mockT.failed {
			t.Error("BeOneOf should pass")
		}
	})
	t.Run("BeOneOf fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeOneOf(mockT, "c", []string{"a", "b"})
		if !mockT.failed {
			t.Error("BeOneOf should fail")
		}
	})
}

func TestContainKey_Integration(t *testing.T) {
	t.Parallel()

	// Test successful cases
	userMap := map[string]int{"name": 1, "age": 2, "email": 3}
	ContainKey(t, userMap, "email")

	intMap := map[int]string{1: "one", 2: "two", 3: "three"}
	ContainKey(t, intMap, 2)
}

func TestContainValue_Integration(t *testing.T) {
	t.Parallel()

	// Test successful cases
	userMap := map[string]int{"name": 1, "age": 2, "email": 3}
	ContainValue(t, userMap, 2)

	intMap := map[int]string{1: "one", 2: "two", 3: "three"}
	ContainValue(t, intMap, "two")
}
