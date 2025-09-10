package assert

import (
	"fmt"
	"time"
)

// Option is a functional option for configuring assertions.
type Option interface {
	Apply(config *Config)
}

// Config provides configuration options for assertions.
// It allows for custom error messages and future extensibility.
type Config struct {
	Message    string
	IgnoreCase bool
	StackTrace bool
	Time       TimeOptions
	/*
		 	Description    string
			DeepComparison bool
	*/
}

type TimeOptions struct {
	IgnoreTimezone bool
	TruncateUnit   time.Duration
}

// message implements the Option interface for custom messages.
type message string

// ignoreCase is a boolean flag for ignoring case in comparisons.
type ignoreCase bool

// stackTrace is a boolean flag for including stack traces on NotPanic assertions.
type stackTrace bool

// ignoreTimezone configures time comparisons to ignore timezone/location differences
type ignoreTimezone bool

// truncateDuration configures time comparisons to truncate both values before comparing
type truncateDuration time.Duration

// Apply sets the custom message in the config.
func (m message) Apply(c *Config) {
	c.Message = string(m)
}

func (i ignoreCase) Apply(c *Config) {
	c.IgnoreCase = bool(i)
}

func (s stackTrace) Apply(c *Config) {
	c.StackTrace = bool(s)
}

// Apply implements Option for ignoreTimezone
func (i ignoreTimezone) Apply(c *Config) {
	c.Time.IgnoreTimezone = bool(i)
}

// Apply implements Option for truncateDuration
func (u truncateDuration) Apply(c *Config) {
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

// WithStackTrace creates an option for including stack traces on NotPanic assertions.
func WithStackTrace() Option {
	return stackTrace(true)
}

// WithIgnoreTimezone creates an option for ignoring timezone when comparing times.
// When enabled, comparisons use calendar components (year, month, day, hour, minute, second[, ns])
// and do not consider the Location/offset.
func WithIgnoreTimezone() Option {
	return ignoreTimezone(true)
}

// WithTruncate truncates the actual and expected times to the specified unit before comparing.
//
// This is useful for asserting that two times are the same up to a certain level of precision,
// ignoring differences in smaller units.
func WithTruncate(unit time.Duration) Option {
	return truncateDuration(unit)
}
