package grepper

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// GrepLines returns the number of matching lines and any error.
func GrepLines(r io.Reader, w io.Writer, pattern string, opts Config) (int, error) {
	lines, err := readLines(r)
	if err != nil {
		return 0, err
	}

	matcher, err := buildMatcher(pattern, opts)
	if err != nil {
		return 0, err
	}

	matched, count := matchLines(lines, matcher, opts)

	if opts.CountOnly {
		return printCount(w, count)
	}

	toPrint := applyContext(matched, opts)

	return count, printLines(w, lines, matched, toPrint, opts)
}

// readLines reads all lines from the input reader into a slice of strings.
func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// printCount prints only the number of matched lines.
func printCount(w io.Writer, count int) (int, error) {
	_, err := fmt.Fprintln(w, count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// matchLines checks each line against the matcher and returns a boolean
// slice marking matches, along with the total count.
func matchLines(lines []string, matcher func(string) bool, opts Config) ([]bool, int) {
	matched := make([]bool, len(lines))
	count := 0
	for i, line := range lines {
		m := matcher(line)
		if opts.Invert {
			m = !m
		}
		if m {
			matched[i] = true
			count++
		}
	}
	return matched, count
}

// applyContext expands the matched lines with before/after context
// according to the grep options.
func applyContext(matched []bool, opts Config) []bool {
	n := len(matched)
	toPrint := make([]bool, n)

	after, before := opts.After, opts.Before
	if opts.Context > 0 {
		after, before = opts.Context, opts.Context
	}

	for i := 0; i < n; i++ {
		if !matched[i] {
			continue
		}
		toPrint[i] = true
		for k := max(0, i-before); k < i; k++ {
			toPrint[k] = true
		}
		for k := i + 1; k <= min(n-1, i+after); k++ {
			toPrint[k] = true
		}
	}
	return toPrint
}

// printLines writes the selected lines to the writer with optional
// line number or match prefixes.
func printLines(w io.Writer, lines []string, matched, toPrint []bool, opts Config) error {
	for i := 0; i < len(lines); i++ {
		if !toPrint[i] {
			continue
		}

		prefix := ""
		if opts.WithLineNo {
			if matched[i] {
				prefix = fmt.Sprintf("%d:", i+1)
			} else {
				prefix = fmt.Sprintf("%d-", i+1)
			}
		}

		if _, err := fmt.Fprintln(w, prefix+lines[i]); err != nil {
			return err
		}
	}
	return nil
}

// buildMatcher builds a function that tests whether a line matches
// the given pattern, considering fixed/regex and case sensitivity.
func buildMatcher(pattern string, opts Config) (func(string) bool, error) {
	if opts.Fixed {
		p := pattern
		if opts.IgnoreCase {
			p = strings.ToLower(pattern)
			return func(s string) bool { return strings.Contains(strings.ToLower(s), p) }, nil
		}
		return func(s string) bool { return strings.Contains(s, p) }, nil
	}

	pat := pattern
	if opts.IgnoreCase {
		pat = "(?i)" + pattern
	}
	re, err := regexp.Compile(pat)
	if err != nil {
		return nil, err
	}
	return func(s string) bool { return re.MatchString(s) }, nil
}
