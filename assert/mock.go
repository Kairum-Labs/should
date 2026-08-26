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

// failNowSignalType backs the failNowSignal sentinel. It carries a field
// specifically so it isn't zero-size: the Go spec allows distinct zero-size
// variables to share the same address (allocation may just return &zerobase),
// which would make pointer-identity comparison on a *struct{} unreliable.
type failNowSignalType struct{ _ byte }

// failNowSignal is the panic value mockT.FailNow uses to unwind the current
// goroutine, mirroring the real testing.T.FailNow (which calls runtime.Goexit
// and never returns to the caller). Without it, code placed after a failed
// WithFailFast assertion would keep running inside the test closure, and
// nothing would catch a library bug that forgot to actually stop.
//
// It's compared by pointer identity (not a type assertion), so
// runCatchingFailNow only ever swallows this exact panic -- an assertion
// under test that happens to panic with some other failNowSignalType-shaped
// value still surfaces as a real test failure.
var failNowSignal = &failNowSignalType{}

func (m *mockT) FailNow() {
	m.failed = true
	m.failedNow = true
	panic(failNowSignal)
}

// runCatchingFailNow runs test, catching the panic mockT.FailNow raises so
// callers get a normal return instead of a crashed test binary. Any other
// panic is re-raised so real bugs still surface.
func runCatchingFailNow(test func(t testing.TB), mock *mockT) {
	defer func() {
		if r := recover(); r != nil && r != failNowSignal {
			panic(r)
		}
	}()
	test(mock)
}

func assertFails(t *testing.T, test func(t testing.TB)) (failed bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	runCatchingFailNow(test, mock)
	return mock.failed, mock.message
}

func assertFailsFast(t *testing.T, test func(t testing.TB)) (failed bool, failedNow bool, message string) {
	t.Helper()
	mock := &mockT{T: t}
	runCatchingFailNow(test, mock)
	return mock.failed, mock.failedNow, mock.message
}
