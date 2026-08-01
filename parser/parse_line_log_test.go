package parser

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseLineLog(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    Log
		expectedErr bool
	}{
		"standard log": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50ms",
			expected: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "test",
				Message:  "hey",
				Duration: 50 * time.Millisecond,
				Extra:    map[string]string{},
			},
		},
		"invalid: invalid service": {
			input:       "2026-01-01T00:00:00Z INVALID_SERVICE service=test message=hey duration=50ms",
			expectedErr: true,
		},
		"invalid: invalid date": {
			input:       "INVALID_DATE WARNING service=test message=hey duration=50ms",
			expectedErr: true,
		},
		"invalid: invalid duration": {
			// duration should (for now) only be in milliseconds
			input:       "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50s",
			expectedErr: true,
		},
		"invalid: missing field": {
			// duration should (for now) only be in milliseconds
			input:       "2026-01-01T00:00:00Z WARNING service=test duration=50s",
			expectedErr: true,
		},
		"log with sentence as message": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=\"this log is very important\" duration=50ms",
			expected: Log{
				Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
				Level:    Warning,
				Service:  "test",
				Message:  "\"this log is very important\"",
				Duration: 50 * time.Millisecond,
				Extra:    map[string]string{},
			},
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
			if diff := cmp.Diff(test.expected, log); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestParseLogContent(t *testing.T) {
	type ParseLogOutput struct {
		Log []*Log
		Err []error
	}

	tests := map[string]struct {
		inputContent  string
		inputSettings ParseSettings
		expected      ParseLogOutput
	}{
		"no content should return nothing": {
			inputContent:  "",
			inputSettings: ParseSettings{},
			expected:      ParseLogOutput{},
		},
		"invalid log should return an error": {
			inputContent:  "INVALID LOG",
			inputSettings: ParseSettings{},
			expected: ParseLogOutput{
				Err: []error{&ValueError{
					ErroredValue:  "INVALID",
					ExpectedValue: "time format - RFC3339",
				}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, err := parseLog(test.inputContent, &test.inputSettings)
			output := ParseLogOutput{
				Log: log,
				Err: err,
			}
			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
