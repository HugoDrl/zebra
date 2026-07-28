package analyser

import (
	"github.com/HugoDrl/LogParser/parser"
	"testing"
	"time"
)

func TestMetricsEquality(t *testing.T) {
	metric := &CollectionMetric{
		Lines: map[parser.Level]int{},
		ServicePerformance: map[string]ServiceMetric{
			"database": {
				Name:            "database",
				Lines:           3,
				AverageDuration: 50 * time.Millisecond,
			},
			"api": {
				Name:            "api",
				Lines:           1,
				AverageDuration: 22 * time.Millisecond,
			},
		},
	}
	tests := map[string]struct {
		input    *CollectionMetric
		expected bool
	}{
		"same pointer should be equal": {
			input:    metric,
			expected: true,
		},
		"not the same pointer but same value should be equal": {
			input: &CollectionMetric{
				Lines: map[parser.Level]int{},
				ServicePerformance: map[string]ServiceMetric{
					"database": {
						Name:            "database",
						Lines:           3,
						AverageDuration: 50 * time.Millisecond,
					},
					"api": {
						Name:            "api",
						Lines:           1,
						AverageDuration: 22 * time.Millisecond,
					},
				},
			},
			expected: true,
		},
		"service metrics with capital letters instead of minimal should detect diff": {
			input: &CollectionMetric{
				Lines: map[parser.Level]int{},
				ServicePerformance: map[string]ServiceMetric{
					"database": {
						Name:            "database",
						Lines:           3,
						AverageDuration: 50 * time.Millisecond,
					},
					// API should differ from api
					"API": {
						Name:            "API",
						Lines:           1,
						AverageDuration: 22 * time.Millisecond,
					},
				},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := metric.IsEqual(test.input)

			if output != test.expected {
				t.Fatalf("%s: expected %t - got %t", name, test.expected, output)
			}
		})
	}
}
