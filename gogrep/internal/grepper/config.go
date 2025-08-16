package grepper

import (
	"errors"
)

type Config struct {
	After   int // -A N: print N lines of trailing context after each matching line
	Before  int // -B N: print N lines of leading context before each matching line
	Context int // -C N: print N lines of context around each matching line (sets both -A and -B unless explicitly overridden)

	CountOnly  bool // -c: print only the count of matching lines instead of the lines themselves
	IgnoreCase bool // -i: ignore case distinctions when matching
	Invert     bool // -v: invert the match, selecting non-matching lines
	Fixed      bool // -F: interpret the pattern as a fixed string instead of a regular expression
	WithLineNo bool // -n: prefix each output line with its line number
}

// Validate checks the configuration for invalid or conflicting options.
// Returns an error if any constraints are violated.
func (c *Config) Validate() error {
	if c.Context < 0 || c.After < 0 || c.Before < 0 {
		return errors.New("invalid arguments: context, before, and after must be non-negative")
	}
	return nil
}
