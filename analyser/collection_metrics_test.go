package analyser

import (
	"testing"
	"time"

	"github.com/HugoDrl/LogParser/parser"
	"github.com/google/go-cmp/cmp"
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

	tests := map[string]struct {
		input    *parser.Log
		expected CollectionMetric
	}{
		"empty log should not change anything": {
			input:    &parser.Log{},
			expected: getDefaultMetric(),
		},
		"nil pointer should not change anything": {
			input:    nil,
			expected: getDefaultMetric(),
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
		"happending a log with a different duration should modify average properly": {
			input: &parser.Log{
				Duration: 100 * time.Millisecond,
				Service:  "api",
				Level:    parser.Error,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Error]++

				// Average duration should not move for this one
				perf := metric.ServicePerformance["api"]
				perf.Lines++
				perf.AverageDuration = (perf.AverageDuration*time.Duration(perf.Lines-1) + 100*time.Millisecond) / time.Duration(perf.Lines)
				metric.ServicePerformance["api"] = perf

				return metric
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			metric := getDefaultMetric()
			metric.handleService(test.input)

			if diff := cmp.Diff(test.expected, metric, cmp.AllowUnexported(CollectionMetric{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
