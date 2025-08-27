package cutter

import (
	"strings"
)

type Extractor struct {
	delimiter string
	fields    []int
	separated bool
}

func NewExtractor(delimiter string, fields []int, separated bool) *Extractor {
	return &Extractor{
		delimiter: delimiter,
		fields:    fields,
		separated: separated,
	}
}

func (e *Extractor) Extract(line string) string {
	if e.separated && !strings.Contains(line, e.delimiter) {
		return ""
	}

	parts := strings.Split(line, e.delimiter)
	selected := make([]string, 0, len(e.fields))

	for _, f := range e.fields {
		idx := f - 1
		if idx >= 0 && idx < len(parts) {
			selected = append(selected, parts[idx])
		}
	}

	if len(selected) == 0 {
		return ""
	}
	return strings.Join(selected, e.delimiter)
}
