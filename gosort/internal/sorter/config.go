package sorter

import (
	"errors"
	"strings"
)

type Config struct {
	Key            int  // column number to sort by (starting from 1, 0 means whole line)
	Numeric        bool // -n: sort numerically
	Reverse        bool // -r: reverse order
	Unique         bool // -u: output only unique lines
	Month          bool // -M: sort by month name (Jan...Dec)
	IgnoreTrailing bool // -b: ignore trailing blanks
	CheckIfSorted  bool // -c: check if already sorted
	HumanNum       bool // -h: human-readable numeric sort (1K, 2M...)
}

// Mutually exclusive flags
const (
	FlagNumeric = "numeric"
	FlagMonth   = "month"
	FlagHuman   = "human-readable-numeric"
)

func (cfg Config) Validate() error {
	var enabled []string
	if cfg.Numeric {
		enabled = append(enabled, FlagNumeric)
	}
	if cfg.Month {
		enabled = append(enabled, FlagMonth)
	}
	if cfg.HumanNum {
		enabled = append(enabled, FlagHuman)
	}

	if len(enabled) > 1 {
		return errors.New(strings.Join(enabled, ":") + ": mutually exclusive flags")
	}

	if cfg.Key < 0 {
		return errors.New("invalid key: must be >= 0")
	}

	return nil
}
