package parser

import (
	"testing"
	"time"
)

func TestLogFilter(t *testing.T) {
	testLog := &Log{
		Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
		Level:    Error,
		Duration: 5 * time.Second,
		Message:  "\"response time out \"",
		Service:  "api",
		extra: map[string]string{
			"path": "/",
		},
	}
	tests := map[string]struct {
		input    ParseSettings
		expected bool
	}{
		"filter by date that includes the log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2027, 01, 01, 0, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := filterLogs(testLog, &test.input)
			if output != test.expected {
				t.Fatalf("%s: expected %t - got %t", name, test.expected, output)
			}
		})
	}
}
