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
		config      []AssertionConfig
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
			config: []AssertionConfig{
				{Message: "custom message"},
			},
			shouldFail:  true,
			expectedMsg: "custom message\nExpected panic, but did not panic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTB{}
			Panic(mockT, tc.fn, tc.config...)

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
		config      []AssertionConfig
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
			config: []AssertionConfig{
				{Message: "custom message"},
			},
			shouldFail:  true,
			expectedMsg: "custom message\nExpected for the function to not panic, but it panicked with: some panic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTB{}
			NotPanic(mockT, tc.fn, tc.config...)

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

	// BeNotEmpty
	t.Run("BeNotEmpty passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeNotEmpty(mockT, "not empty")
		if mockT.failed {
			t.Error("BeNotEmpty should pass")
		}
	})
	t.Run("BeNotEmpty fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeNotEmpty(mockT, "")
		if !mockT.failed {
			t.Error("BeNotEmpty should fail")
		}
	})

	// BeNil
	t.Run("BeNil passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeNil[any](mockT, nil)
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

	// BeNotNil
	t.Run("BeNotNil passes", func(t *testing.T) {
		mockT := &mockTB{}
		var x = 1
		BeNotNil(mockT, &x)
		if mockT.failed {
			t.Error("BeNotNil should pass")
		}
	})
	t.Run("BeNotNil fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeNotNil[any](mockT, nil)
		if !mockT.failed {
			t.Error("BeNotNil should fail")
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

	// BeGreaterOrEqualThan
	t.Run("BeGreaterOrEqualThan passes", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterOrEqualThan(mockT, 10, 10)
		if mockT.failed {
			t.Error("BeGreaterOrEqualThan should pass")
		}
	})
	t.Run("BeGreaterOrEqualThan fails", func(t *testing.T) {
		mockT := &mockTB{}
		BeGreaterOrEqualThan(mockT, 9, 10)
		if !mockT.failed {
			t.Error("BeGreaterOrEqualThan should fail")
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
}
