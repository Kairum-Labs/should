package assert

import (
	"fmt"
	"time"
)

// Option configures an assertion.
//
// Options are created by the With* helpers.
type Option interface {
	apply(config *config)
}

// config stores option state while an assertion is evaluated.
type config struct {
	Message    string
	IgnoreCase bool
	// Deprecated: NotPanic includes stack traces by default.
	StackTrace bool
	FailFast   bool
	Time       timeOptions
	/*
		 	Description    string
			DeepComparison bool
	*/
}

// Config provides configuration options for assertions.
// It allows for custom error messages and future extensibility.
//
// Deprecated: Config is an implementation detail. Use the With* helpers instead.
// Config will be removed in a future major release.
type Config = config

// timeOptions stores configuration for time comparisons.
type timeOptions struct {
	IgnoreTimezone bool
	TruncateUnit   time.Duration
}

// TimeOptions provides configuration for time comparisons.
//
// Deprecated: TimeOptions is an implementation detail. Use the With* helpers instead.
// TimeOptions will be removed in a future major release.
type TimeOptions = timeOptions

// message implements the Option interface for custom messages.
type message string

// ignoreCase is a boolean flag for ignoring case in comparisons.
type ignoreCase bool

// deprecatedStackTrace preserves compatibility with WithStackTrace.
type deprecatedStackTrace struct{}

// failFast is a boolean flag for stopping the test on assertion failure.
type failFast bool

// ignoreTimezone configures time comparisons to ignore timezone/location differences
type ignoreTimezone bool

// truncateDuration configures time comparisons to truncate both values before comparing
type truncateDuration time.Duration

// apply sets the custom message in the config.
func (m message) apply(c *config) {
	c.Message = string(m)
}

func (i ignoreCase) apply(c *config) {
	c.IgnoreCase = bool(i)
}

func (deprecatedStackTrace) apply(*config) {}

func (f failFast) apply(c *config) {
	c.FailFast = bool(f)
}

// apply implements Option for ignoreTimezone.
func (i ignoreTimezone) apply(c *config) {
	c.Time.IgnoreTimezone = bool(i)
}

// apply implements Option for truncateDuration.
func (u truncateDuration) apply(c *config) {
	c.Time.TruncateUnit = time.Duration(u)
}

// WithMessage creates an option for setting a custom error message.
//
// The message is treated as a plain string literal. Use this when you
// want to display a fixed message without formatting or placeholders.
//
// Example usage:
//
//	should.BeGreaterThan(t, userAge, 18, should.WithMessage("User must be adult"))
//
// See also: [WithMessagef] for messages that include formatting placeholders.
func WithMessage(msg string) Option {
	return message(msg)
}

// WithMessagef creates an option for setting a custom error message with formatting.
//
// The message supports placeholders, similar to fmt.Sprintf, and takes
// optional arguments to replace them. Use this when you need dynamic
// content in the message.
//
// Example usage:
//
//	should.BeLessOrEqualTo(t, score, 100, should.WithMessagef("Score cannot exceed %d", 100))
func WithMessagef(msg string, args ...any) Option {
	return message(fmt.Sprintf(msg, args...))
}

// WithIgnoreCase creates an option for ignoring case in comparisons.
func WithIgnoreCase() Option {
	return ignoreCase(true)
}

// WithStackTrace is retained for compatibility and has no effect because NotPanic
// includes stack traces by default.
//
// Deprecated: NotPanic includes stack traces by default.
func WithStackTrace() Option {
	return deprecatedStackTrace{}
}

// WithFailFast stops the test immediately when an assertion fails, instead of
// marking it failed and allowing execution to continue.
//
// It has the same goroutine restrictions as [testing.T.FailNow]. Use it for
// preconditions where continuing after a failed assertion would be misleading
// or could cause a panic.
func WithFailFast() Option {
	return failFast(true)
}

// WithIgnoreTimezone creates an option for ignoring timezone when comparing times in [BeSameTime].
// When enabled, comparisons use calendar components (year, month, day, hour, minute, second[, ns])
// and do not consider the Location/offset.
func WithIgnoreTimezone() Option {
	return ignoreTimezone(true)
}

// WithTruncate truncates the actual and expected times to the specified unit before comparing them in [BeSameTime].
//
// This is useful for asserting that two times are the same up to a certain level of precision,
// ignoring differences in smaller units.
func WithTruncate(unit time.Duration) Option {
	return truncateDuration(unit)
}
