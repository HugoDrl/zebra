package parser

import (
	"testing"
	"time"
)

func TestParseLineLog(t *testing.T) {
	tests := map[string]struct{
		input string
		expected Log
		expectedErr bool
	}{
		"standard log": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50ms",
			expected: Log{
				Time: time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level: Warning,
				Service: "test",
				Message: "hey",
				Duration: 50 * time.Millisecond,
			},
		},
		"invalid service": {
			input: "2026-01-01T00:00:00Z INVALID_SERVICE service=test message=hey duration=50ms",
			expectedErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, err := parseLine(test.input)
			if (err != nil) != test.expectedErr {
				t.Fatalf(
					"%s: unexpected err behavior - got %t, expected %t (err: %s)",
					name,
					err != nil,
					test.expectedErr,
					err,
				)
			}
			if (log == nil) && test.expected.Equal(Log{}) {
				t.Skip()
			}
			if !(*log).Equal(test.expected) {
				t.Fatalf(
					"%s: expected %s - got %s",
					name,
					test.expected,
					log,
				)
			}
		})
	}
}
