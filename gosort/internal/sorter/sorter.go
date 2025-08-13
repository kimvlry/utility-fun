package sorter

import (
	"errors"
	"gosort/internal/parsers"
	"sort"
	"strings"
)

// sortableLine represents a single input line along with precomputed
// values (sorting keys) that are used during the sort comparison.
//
// The idea is to avoid repeated parsing of the same string during sorting.
// Instead, we do all parsing once in advance, store the results here,
// and then compare only these prepared values in O(1) time.

// Example of prepared for input lines ["Jan\t100", "Feb\t50"] with -k 2 -n:
// [
//
//	{ original: "Jan\t100", keyStr: "100", keyNum: 100, keyMonth: 0 },
//	{ original: "Feb\t50",  keyStr: "50",  keyNum: 50,  keyMonth: 0 },
//
// ]
// Sorting will directly compare keyNum = 100 vs 50, no reparsing needed.
type sortableLine struct {
	original string  // original — the full line as it appeared in the input.
	keyStr   string  // keyStr   — the string representation of the key to sort by (if sorting lexicographically).
	keyNum   float64 // keyNum   — the numeric representation of the key (used for numeric or human-readable numeric sort).
	keyMonth int     // keyMonth — the numeric index of a month (1 = Jan, 2 = Feb, ..., 12 = Dec) for month-name sort.
}

// uniquePrepared removes duplicate lines from an already sorted prepared slice.
//
// It assumes the slice is sorted according to the active sort criteria.
// The comparison for uniqueness (-u) is done on the 'original' field to ensure
// that duplicates are detected exactly as they appear in the input (full-line match).
func uniquePrepared(prepared []sortableLine) []sortableLine {
	if len(prepared) == 0 {
		return prepared
	}
	out := []sortableLine{prepared[0]}
	for i := 1; i < len(prepared); i++ {
		if prepared[i].original != prepared[i-1].original {
			out = append(out, prepared[i])
		}
	}
	return out
}

// SortLines sorts lines according to the provided Config.
func SortLines(lines []string, cfg Config) ([]string, error) {
	// Precompute sort keys for all lines
	prepared := make([]sortableLine, len(lines))
	for i, line := range lines {
		if cfg.IgnoreTrailing {
			line = strings.TrimRight(line, " ")
		}

		key := extractField(line, cfg)
		sl := sortableLine{
			original: line,
			keyStr:   key,
		}

		// Parse numeric values if needed
		if cfg.Numeric {
			if num, err := parsers.ParseFloat(key); err == nil {
				sl.keyNum = num
			}
		}

		// Parse human-readable sizes if needed
		if cfg.HumanNum {
			if num, err := parsers.ParseHumanNumber(key); err == nil {
				sl.keyNum = num
			}
		}

		// Parse month names if needed
		if cfg.Month {
			if m, err := parsers.ParseMonth(key); err == nil {
				sl.keyMonth = m
			}
		}

		prepared[i] = sl
	}

	// Define comparison function using precomputed keys
	less := func(i, j int) bool {
		a := prepared[i]
		b := prepared[j]

		// Month sort has the highest priority if enabled
		if cfg.Month && a.keyMonth != 0 && b.keyMonth != 0 {
			if cfg.Reverse {
				return a.keyMonth > b.keyMonth
			}
			return a.keyMonth < b.keyMonth
		}

		// Human-readable or numeric sort
		if (cfg.HumanNum || cfg.Numeric) && a.keyNum != 0 && b.keyNum != 0 {
			if cfg.Reverse {
				return a.keyNum > b.keyNum
			}
			return a.keyNum < b.keyNum
		}

		// Fallback to lexicographical comparison
		if cfg.Reverse {
			return b.keyStr < a.keyStr
		}
		return a.keyStr < b.keyStr
	}

	// -c flag: only check if sorted
	if cfg.CheckIfSorted {
		if sort.SliceIsSorted(prepared, less) {
			return nil, nil
		}
		return nil, errors.New("input is not sorted")
	}

	// Sort using precomputed keys
	sort.Slice(prepared, less)

	// Remove duplicates if -u flag is set
	if cfg.Unique {
		prepared = uniquePrepared(prepared)
	}

	// Extract sorted original lines
	result := make([]string, len(prepared))
	for i, sl := range prepared {
		result[i] = sl.original
	}

	return result, nil
}

// extractField extracts the key field from a line according to -k flag.
// If Key <= 0, the entire line is returned.
// Fields are split by tab character.
func extractField(line string, cfg Config) string {
	if cfg.Key <= 0 {
		return line
	}
	cols := strings.Split(line, "\t")
	if cfg.Key-1 < len(cols) {
		return cols[cfg.Key-1]
	}
	return ""
}
