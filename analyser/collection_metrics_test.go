package analyser

import (
	"testing"
	"time"

	"github.com/HugoDrl/LogParser/parser"
)

func TestHandleService(t *testing.T) {
	getDefaultMetric := func() CollectionMetric {
		return CollectionMetric{
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
	}

	defaultMetric := getDefaultMetric()

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
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Info]++

				// Average duration should not move for this one
				perf := metric.ServicePerformance["api"]
				perf.Lines++
				metric.ServicePerformance["api"] = perf

				return metric
			}(),
		},
		"other service log should create a new ServiceMetric": {
			input: &parser.Log{
				Duration: 100 * time.Millisecond,
				Service:  "database",
				Level:    parser.Info,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Info]++
				metric.ServicePerformance["database"] = ServiceMetric{
					Name:            "database",
					Lines:           1,
					AverageDuration: 100 * time.Millisecond,
				}
				return metric
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			metric := getDefaultMetric()
			metric.handleService(test.input)

			if !metric.IsEqual(test.expected) {
				t.Fatalf("%s: expected %v - got %v", name, test.expected, metric)
			}
		})
	}
}
