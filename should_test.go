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
