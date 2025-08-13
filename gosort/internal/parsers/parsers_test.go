package parsers

import (
	"testing"
)

func TestParseFloat(t *testing.T) {
	tests := []struct {
		in    string
		want  float64
		isErr bool
	}{
		{"123", 123, false},
		{"  45.6 ", 45.6, false},
		{"abc", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseFloat(tt.in)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseFloat(%q) error = %v, wantErr %v", tt.in, err, tt.isErr)
		}
		if !tt.isErr && got != tt.want {
			t.Errorf("ParseFloat(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseMonth(t *testing.T) {
	tests := []struct {
		in    string
		want  int
		isErr bool
	}{
		{"Jan", 1, false},
		{"dec", 12, false},
		{"  Feb  ", 2, false},
		{"month", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseMonth(tt.in)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseMonth(%q) error = %v, wantErr %v", tt.in, err, tt.isErr)
		}
		if !tt.isErr && got != tt.want {
			t.Errorf("ParseMonth(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseHumanNumber(t *testing.T) {
	tests := []struct {
		in    string
		want  float64
		isErr bool
	}{
		{"1K", 1024, false},
		{"2M", 2 * 1024 * 1024, false},
		{"3G", 3 * 1024 * 1024 * 1024, false},
		{"500", 500, false},
		{"bad", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseHumanNumber(tt.in)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseHumanNumber(%q) error = %v, wantErr %v", tt.in, err, tt.isErr)
		}
		if !tt.isErr && got != tt.want {
			t.Errorf("ParseHumanNumber(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
