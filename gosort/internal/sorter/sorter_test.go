package sorter

import (
	"reflect"
	"testing"
)

func TestSortLines_IndividualFlags(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		cfg     Config
		want    []string
		wantErr bool
	}{
		{
			"numeric",
			[]string{"10", "2"},
			Config{Numeric: true},
			[]string{"2", "10"},
			false,
		},

		{
			"reverse",
			[]string{"a", "b"},
			Config{Reverse: true},
			[]string{"b", "a"},
			false,
		},

		{
			"unique",
			[]string{"a", "a", "b"},
			Config{Unique: true},
			[]string{"a", "b"},
			false,
		},

		{
			"column2",
			[]string{"a\t2", "b\t1"},
			Config{Key: 2},
			[]string{"b\t1", "a\t2"},
			false,
		},

		{
			"month",
			[]string{"Feb", "Jan"},
			Config{Month: true},
			[]string{"Jan", "Feb"},
			false,
		},

		{
			"ignoreBlanks",
			[]string{"a ", "a"},
			Config{IgnoreTrailing: true},
			[]string{"a", "a"},
			false,
		},

		{
			"human",
			[]string{"1K", "500"},
			Config{HumanNum: true},
			[]string{"500", "1K"},
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

			got, err := SortLines(tt.lines, tt.cfg)
			if err != nil {
				t.Fatalf("SortLines() error = %v, wantErr false", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestSortLines_Combinations(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		cfg     Config
		want    []string
		wantErr bool
	}{
		{
			"numeric_reverse",
			[]string{"10", "2"},
			Config{Numeric: true, Reverse: true},
			[]string{"10", "2"},
			false,
		},

		{
			"column2_numeric_unique",
			[]string{"a\t2", "b\t2", "c\t1"},
			Config{Key: 2, Numeric: true, Unique: true},
			[]string{"c\t1", "a\t2", "b\t2"},
			false,
		},

		{
			"month_ignoreBlanks",
			[]string{" Feb", "Jan"},
			Config{Month: true, IgnoreTrailing: true},
			[]string{"Jan", " Feb"},
			false,
		},

		{
			"human_reverse",
			[]string{"1K", "500"},
			Config{HumanNum: true, Reverse: true},
			[]string{"1K", "500"},
			false,
		},

		{
			"conflict_n_M",
			[]string{"Jan", "10"},
			Config{Numeric: true, Month: true},
			nil,
			true,
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

			got, err := SortLines(tt.lines, tt.cfg)
			if err != nil {
				t.Fatalf("SortLines() error = %v, wantErr false", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
