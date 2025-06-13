package should

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

func assertFails(t *testing.T, test func(t testing.TB)) (failed bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	test(mock)
	return mock.failed, mock.message
}
