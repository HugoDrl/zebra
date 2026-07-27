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
		"filter by date that does not include the log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 11, 01, 0, 0, 0, 0, time	.UTC),
			},
			expected: false,
		},
		"filter by service that includes the log": {
			input: ParseSettings{
				Service: "api",
			},
			expected: true,
		},
		"filter by date that includes, but service that does not include log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				Service: "db",
			},
			expected: false,
		},
		"filter by service that includes, but level that does not include log": {
			input: ParseSettings{
				Service: "api",
				Level: Warning,
			},
			expected: false,
		},
		"filter by level that includes log": {
			input: ParseSettings{
				Level: Error,
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
