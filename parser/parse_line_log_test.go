package parser

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseLineLog(t *testing.T) {
	type ParseLineOutput struct {
		Log Log
		Err error
	}
	tests := map[string]struct {
		input    string
		expected ParseLineOutput
	}{
		"standard log": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Log: Log{
					Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
					Level:    Warning,
					Service:  "test",
					Message:  "hey",
					Duration: 50 * time.Millisecond,
					Extra:    map[string]string{},
				},
			},
		},
		"invalid: invalid service": {
			input: "2026-01-01T00:00:00Z INVALID_SERVICE service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Err: &ValueError{
					ExpectedValue: "",
					ErroredValue:  "level",
				},
			},
		},
		"invalid: not enough mandatory arguments": {
			input: "2026-01-01T00:00:00Z",
			expected: ParseLineOutput{
				Err: &ParseError{Reason: "Not enough arguments - expected 2 - found 1"},
			},
		},
		"invalid: invalid date": {
			input: "INVALID_DATE WARNING service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Err: &ValueError{
					ExpectedValue: "time format - RFC3339",
					ErroredValue:  "INVALID_DATE",
				},
			},
		},
		"invalid: invalid duration in seconds instead of ms": {
			// duration should (for now) only be in milliseconds
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50s",
			expected: ParseLineOutput{
				Err: &ValueError{
					ExpectedValue: "duration value in milliseconds on format xxxms",
					ErroredValue:  "50s",
				},
			},
		},
		"invalid: invalid duration is not a number": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=NaNms",
			expected: ParseLineOutput{
				Err: &ValueError{
					ExpectedValue: "duration value in milliseconds on format xxxms",
					ErroredValue:  "NaNms",
				},
			},
		},
		"invalid: missing field": {
			input: "2026-01-01T00:00:00Z WARNING service=test duration=50ms",
			expected: ParseLineOutput{
				Err: &ValueError{
					ExpectedValue: "message field, service field, and duration field, should not be empty.",
					ErroredValue:  "message: , service: test, duration: 50000000",
				},
			},
		},
		"wrong format key=value pairs": {
			input: "2026-01-01T00:00:00Z WARNING service=test=api duration=50ms",
			expected: ParseLineOutput{
				Err: &ParseError{Reason: "Wrong format in key-values design (found [service test api])"},
			},
		},
		"log with sentence as message": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=\"this log is very important\" duration=50ms",
			expected: ParseLineOutput{
				Log: Log{
					Time:     time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
					Level:    Warning,
					Service:  "test",
					Message:  "\"this log is very important\"",
					Duration: 50 * time.Millisecond,
					Extra:    map[string]string{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, err := parseLine(test.input)
			output := ParseLineOutput{
				Log: log,
				Err: err,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
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
		"valid log should parse it": {
			inputContent:  "2026-01-01T00:00:00Z WARNING service=database duration=50ms message=\"query executed\"",
			inputSettings: ParseSettings{},
			expected: ParseLogOutput{
				Log: []*Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    Warning,
						Service:  "database",
						Duration: 50 * time.Millisecond,
						Message:  `"query executed"`,
						Extra:    map[string]string{},
					},
				},
			},
		},
		"log with extra fields should add it": {
			inputContent:  "2026-01-01T00:00:00Z WARNING service=database duration=50ms message=\"query executed\" hey=0 extra_field=\"an extra field\"",
			inputSettings: ParseSettings{},
			expected: ParseLogOutput{
				Log: []*Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    Warning,
						Service:  "database",
						Duration: 50 * time.Millisecond,
						Message:  `"query executed"`,
						Extra: map[string]string{
							"hey":         "0",
							"extra_field": `"an extra field"`,
						},
					},
				},
			},
		},
		"multiple logs should add them": {
			inputContent:  "2026-01-01T00:00:00Z WARNING service=database duration=50ms message=\"query executed\"\n2026-01-02T01:00:00Z INFO service=api duration=10ms message=ok",
			inputSettings: ParseSettings{},
			expected: ParseLogOutput{
				Log: []*Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    Warning,
						Service:  "database",
						Duration: 50 * time.Millisecond,
						Message:  `"query executed"`,
						Extra:    map[string]string{},
					},
					{
						Time:     time.Date(2026, 1, 2, 1, 0, 0, 0, time.UTC),
						Level:    Info,
						Service:  "api",
						Message:  "ok",
						Duration: 10 * time.Millisecond,
						Extra:    map[string]string{},
					},
				},
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
