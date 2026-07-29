package analyser

import (
	"testing"
	"time"

	"github.com/HugoDrl/LogParser/parser"
)

func TestHandleService(t *testing.T) {
	defaultMetric := CollectionMetric{
		Lines: map[parser.Level]int{
			parser.Warning: 1,
			parser.Info:    10,
		},
		ServicePerformance: map[string]ServiceMetric{
			"api": {
				Name:            "api",
				Lines:           11,
				AverageDuration: 44 * time.Millisecond,
			},
		},
	}
	tests := map[string]struct {
		input    *parser.Log
		expected CollectionMetric
	}{
		"empty log should not change anything": {
			input:    &parser.Log{},
			expected: defaultMetric,
		},
		"nil pointer should not change anything": {
			input:    nil,
			expected: defaultMetric,
		},
		"api log should increase the same ServiceMetric": {
			input: &parser.Log{
				Duration: 44 * time.Millisecond,
				Service:  "api",
				Level:    parser.Info,
			},
			expected: CollectionMetric{
				Lines: map[parser.Level]int{
					parser.Warning: 1,
					parser.Info:    11,
				},
				ServicePerformance: map[string]ServiceMetric{
					"api": {
						Name:            "api",
						Lines:           12,
						AverageDuration: 44 * time.Millisecond,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			metric := defaultMetric
			metric.handleService(test.input)

			if !metric.IsEqual(test.expected) {
				t.Fatalf("%s: expected %v - got %v", name, test.expected, metric)
			}
		})
	}
}
