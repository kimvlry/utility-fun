package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFieldSpec parses string formatted as "1,3-5,7" to a slice: []int{1,3,4,5,7}
func ParseFieldSpec(spec string) ([]int, error) {
	result := make(map[int]struct{})

	parts := strings.Split(spec, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start <= 0 || end <= 0 || start > end {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := start; i <= end; i++ {
				result[i] = struct{}{}
			}
		} else {
			val, err := strconv.Atoi(part)
			if err != nil || val <= 0 {
				return nil, fmt.Errorf("invalid filed number: %s", part)
			}
			result[val] = struct{}{}
		}
	}

	out := make([]int, 0, len(result))
	for k := range result {
		out = append(out, k)
	}

	for i := 0; i < len(out)-1; i++ {
		for j := i + 1; j < len(out); j++ {
			if out[i] > out[j] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}
