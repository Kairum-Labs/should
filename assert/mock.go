package assert

import (
	"fmt"
	"testing"
)

type mockT struct {
	*testing.T
	failed    bool
	failedNow bool
	message   string
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

func (m *mockT) FailNow() {
	m.failed = true
	m.failedNow = true
}

func assertFails(t *testing.T, test func(t testing.TB)) (failed bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	test(mock)
	return mock.failed, mock.message
}

func assertFailsFast(t *testing.T, test func(t testing.TB)) (failed bool, failedNow bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	test(mock)
	return mock.failed, mock.failedNow, mock.message
}
