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
		input    CollectionMetric
		expected bool
	}{
		"same value should be equal": {
			input: CollectionMetric{
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
			input: CollectionMetric{
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
		"different query should result in neq": {
			input: CollectionMetric{
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
				Query: "hey this is different query",
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

func TestEmptyMetricsEquality(t *testing.T) {
	metric := CollectionMetric{}
	tests := map[string]struct {
		input    CollectionMetric
		expected bool
	}{
		"same value should be equal": {
			input:    CollectionMetric{},
			expected: true,
		},
		"nil values should be equal": {
			input: CollectionMetric{
				Lines:              nil,
				ServicePerformance: nil,
				FileErrors:         nil,
				ParsingErrorCount:  0,
				SlowestInput:       nil,
				Query:              "",
			},
			expected: true,
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
