package assert

import (
	"fmt"
	"testing"
)

type mockT struct {
	*testing.T
	failed  bool
	message string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.message = fmt.Sprintf(format, args...)
}

func (m *mockT) Error(args ...interface{}) {
	m.failed = true
	m.message = fmt.Sprint(args...)
}

func (m *mockT) Helper() {
	// No-op for mock implementation
}

func (m *mockT) Failed() bool {
	return m.failed
}

func assertFails(t *testing.T, test func(t testing.TB)) (failed bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	test(mock)
	return mock.failed, mock.message
}
