package parser

import (
	"testing"
	"time"
)

func TestLogEqual(t *testing.T) {
	tests := map[string]struct {
		log1     Log
		log2     Log
		expected bool
	}{
		"empty logs should be equal": {
			log1:     Log{},
			log2:     Log{},
			expected: true,
		},
		"regular logs that are equal": {
			log1: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
			},
			log2: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
			},
			expected: true,
		},
		"logs that don't have the same message": {
			log1: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "test",
				Duration: 0,
			},
			log2: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
			},
			expected: false,
		},
		"logs that don't have the same time": {
			log1: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 1, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "test",
				Duration: 0,
			},
			log2: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
			},
			expected: false,
		},
		"log that have metadata": {
			log1: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 1, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
			},
			log2: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "db",
				Message:  "",
				Duration: 0,
				extra: map[string]string{
					"extra": "true",
				},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := test.log1.Equal(test.log2)
			if output != test.expected {
				t.Fatalf("%s: expected %t - got %t", name, test.expected, output)
			}
		})
	}
}
