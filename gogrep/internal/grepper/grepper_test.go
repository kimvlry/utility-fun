package grepper

import (
	"bytes"
	"errors"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// runUnixGrep executes the system grep command with given input, pattern, and flags.
// Returns output lines or error.
func runUnixGrep(input, pattern string, flags []string) ([]string, error) {
	cmdArgs := append(flags, pattern)
	cmd := exec.Command("grep", cmdArgs...)
	cmd.Stdin = strings.NewReader(input)

	out, err := cmd.Output()
	if err != nil {
		// If grep returns exit code 1, it means no match
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == 1 {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	return lines, nil
}

// runGoGrep executes this Go implementation with the given input, pattern, and config.
// Returns output lines or error.
func runGoGrep(input, pattern string, cfg Config) ([]string, error) {
	var out bytes.Buffer
	_, err := GrepLines(strings.NewReader(input), &out, pattern, cfg)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	return lines, nil
}

// helper to build flag slice from Config for system grep
func buildGrepFlags(cfg Config) []string {
	flags := []string{}

	numFlags := map[string]int{
		"-A": cfg.After,
		"-B": cfg.Before,
		"-C": cfg.Context,
	}
	for f, val := range numFlags {
		if val > 0 {
			flags = append(flags, f, strconv.Itoa(val))
		}
	}

	boolFlags := map[string]bool{
		"-c": cfg.CountOnly,
		"-i": cfg.IgnoreCase,
		"-v": cfg.Invert,
		"-F": cfg.Fixed,
		"-n": cfg.WithLineNo,
	}
	for f, set := range boolFlags {
		if set {
			flags = append(flags, f)
		}
	}

	return flags
}

func TestGrepLines_IndividualFlags(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		cfg     Config
		pattern string
		wantErr bool
	}{
		{
			"after_context",
			"alpha\nbeta\ngamma\n",
			Config{After: 1},
			"alpha",
			false,
		},
		{
			"before_context",
			"alpha\nbeta\ngamma\n",
			Config{Before: 1},
			"beta",
			false,
		},
		{"context",
			"alpha\nbeta\ngamma\n",
			Config{Context: 1},
			"beta",
			false,
		},
		{
			"count_only",
			"alpha\nbeta\nalpha\n",
			Config{CountOnly: true},
			"alpha",
			false,
		},
		{
			"ignore_case",
			"Alpha\nbeta\n",
			Config{IgnoreCase: true},
			"alpha",
			false,
		},
		{
			"invert_match",
			"alpha\nbeta\n",
			Config{Invert: true},
			"alpha",
			false,
		},
		{
			"fixed_strings",
			"alpha[123]\n",
			Config{Fixed: true},
			"alpha[123]",
			false,
		},
		{
			"line_number",
			"alpha\nbeta\n",
			Config{WithLineNo: true},
			"beta",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			goLines, err := runGoGrep(tt.input, tt.pattern, tt.cfg)
			if err != nil {
				t.Fatalf("GrepLines() error = %v", err)
			}

			flags := buildGrepFlags(tt.cfg)
			sysLines, err := runUnixGrep(tt.input, tt.pattern, flags)
			if err != nil {
				t.Fatalf("system grep error = %v", err)
			}

			if !reflect.DeepEqual(goLines, sysLines) {
				t.Errorf("got %#v\nwant %#v", goLines, sysLines)
			}
		})
	}
}

func TestGrepLines_Combinations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		cfg     Config
		pattern string
		wantErr bool
	}{
		{
			"context_line_number_ignorecase",
			"alpha\nBeta\ngamma\n",
			Config{Context: 1, WithLineNo: true, IgnoreCase: true},
			"beta",
			false,
		},

		{
			"after_invert",
			"alpha\nbeta\ngamma\n",
			Config{After: 1, Invert: true},
			"beta",
			false,
		},

		{
			"count_only_conflict_with_after",
			"alpha\nbeta\n",
			Config{CountOnly: true, After: 1},
			"alpha",
			false,
		},

		{
			"negative_after",
			"alpha\nbeta\n",
			Config{After: -1},
			"alpha",
			true,
		},

		{
			"count_with_line_number",
			"alpha\nbeta\ngamma\nbeta\n",
			Config{CountOnly: true, WithLineNo: true},
			"beta",
			false,
		},

		{
			"count_with_after",
			"alpha\nbeta\ngamma\ndelta\n",
			Config{CountOnly: true, After: 2},
			"beta",
			false,
		},

		{
			"count_with_after_multiple_matches",
			"foo\nbar\nbaz\nfoo\nbar\nbaz\n",
			Config{CountOnly: true, After: 2},
			"foo",
			false,
		},

		{
			"count_with_line_number_single_match",
			"line1\nline2\nmatch\nline4\n",
			Config{CountOnly: true, WithLineNo: true},
			"match",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			goLines, err := runGoGrep(tt.input, tt.pattern, tt.cfg)
			if err != nil {
				t.Fatalf("GrepLines() error = %v", err)
			}

			flags := buildGrepFlags(tt.cfg)
			sysLines, err := runUnixGrep(tt.input, tt.pattern, flags)
			if err != nil {
				t.Fatalf("system grep error = %v", err)
			}

			if !reflect.DeepEqual(goLines, sysLines) {
				t.Errorf("got %#v\nwant %#v", goLines, sysLines)
			}
		})
	}
}
