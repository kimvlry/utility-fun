package cutter_test

import (
	"bytes"
	"cut/internal/cutter"
	"cut/internal/parser"
	"os/exec"
	"strings"
	"testing"
)

// runUnixCut runs unix cut and returns it's output to compare result with
func runUnixCut(input, fieldSpec, delimiter string, separated bool) (string, error) {
	args := []string{"-f", fieldSpec}
	if delimiter != "\t" {
		args = append(args, "-d", delimiter)
	}
	if separated {
		args = append(args, "-s")
	}

	cmd := exec.Command("cut", args...)
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}

// runGoCut returns simple cut implementation result
func runGoCut(input, fieldSpec, delimiter string, separated bool) (string, error) {
	fields, err := parser.ParseFieldSpec(fieldSpec)
	if err != nil {
		return "", err
	}
	ex := cutter.NewExtractor(delimiter, fields, separated)

	var buf bytes.Buffer
	for _, line := range strings.Split(strings.TrimRight(input, "\n"), "\n") {
		if out := ex.Extract(line); out != "" {
			buf.WriteString(out)
			buf.WriteByte('\n')
		}
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// compare unix cut and simple cut results
func TestExtractor_AgainstUnixCut(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldSpec string
		delimiter string
		separated bool
	}{
		{
			"single_field",
			"a,b,c\nd,e,f",
			"2",
			",",
			false,
		},
		{
			"range_fields",
			"1,2,3,4,5\n6,7,8,9,10",
			"2-4",
			",",
			false,
		},
		{
			"mixed_fields",
			"alpha,beta,gamma,delta",
			"1,3-4",
			",",
			false,
		},
		{
			"ignore_out_of_bounds",
			"x,y",
			"1,3-5",
			",",
			false,
		},
		{
			"separated_flag",
			"no_delimiter_here\nfoo,bar",
			"1",
			",",
			true,
		},
		{
			"default_delimiter_tab",
			"a\tb\tc\nd\te\tf",
			"1,3",
			"\t",
			false,
		},
		{
			"single_char_delimiter",
			"a|b|c|d",
			"2,4",
			"|",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sysOut, err := runUnixCut(tt.input, tt.fieldSpec, tt.delimiter, tt.separated)
			if err != nil {
				t.Fatalf("unix cut failed: %v", err)
			}

			goOut, err := runGoCut(tt.input, tt.fieldSpec, tt.delimiter, tt.separated)
			if err != nil {
				t.Fatalf("go cut failed: %v", err)
			}

			if sysOut != goOut {
				t.Errorf("mismatch:\n--- unix cut ---\n%q\n--- go cut ---\n%q", sysOut, goOut)
			}
		})
	}
}
