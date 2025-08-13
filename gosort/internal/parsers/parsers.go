package parsers

import (
	"errors"
	"strconv"
	"strings"
)

// ParseFloat tries to parse a string into float64.
// It trims spaces before parsing.
// Returns an error if the string cannot be parsed as a valid number.
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// ParseMonth converts a month name (Jan, Feb, ... Dec) into its numeric index (1-12).
// It is case-insensitive and ignores surrounding spaces.
func ParseMonth(s string) (int, error) {
	months := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}
	normalized := strings.ToLower(strings.TrimSpace(s))
	if val, ok := months[normalized]; ok {
		return val, nil
	}
	return 0, errors.New("invalid month: " + s)
}

// ParseHumanNumber parses strings like "1K", "2M", "3G" into numeric values
// using binary multiples (K=1024, M=1024^2, G=1024^3).
// If no suffix is present, it falls back to a regular float parsing.
// It ignores case and surrounding spaces.
func ParseHumanNumber(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty string")
	}

	unitMultipliers := map[string]float64{
		"k": 1024,
		"m": 1024 * 1024,
		"g": 1024 * 1024 * 1024,
	}

	lower := strings.ToLower(s)
	for unit, multiplier := range unitMultipliers {
		if strings.HasSuffix(lower, unit) {
			numPart := strings.TrimSuffix(lower, unit)
			val, err := strconv.ParseFloat(numPart, 64)
			if err != nil {
				return 0, err
			}
			return val * multiplier, nil
		}
	}

	// No recognized suffix — parse as plain float
	return strconv.ParseFloat(lower, 64)
}
