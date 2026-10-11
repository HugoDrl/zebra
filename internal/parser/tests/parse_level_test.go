package parser_test

import (
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
)

func TestLevelParse(t *testing.T) {
	tests := map[string]struct {
		input      string
		expected   parser.Level
		expectedOk bool
	}{
		"WARNING level should exist": {
			input:      "WARNING",
			expected:   parser.Warning,
			expectedOk: true,
		},
		"warning level should exist": {
			input:      "warning",
			expected:   parser.Warning,
			expectedOk: true,
		},
		"waRNiNg level minimal should exist": {
			input:      "waRNiNg",
			expected:   parser.Warning,
			expectedOk: true,
		},
		"error level should exist": {
			input:      "error",
			expected:   parser.Error,
			expectedOk: true,
		},
		"info level should exist": {
			input:      "info",
			expected:   parser.Info,
			expectedOk: true,
		},
		"debug level should exist": {
			input:      "DEBUG",
			expected:   parser.Debug,
			expectedOk: true,
		},
		"fatal level should exist": {
			input:      "Fatal",
			expected:   parser.Fatal,
			expectedOk: true,
		},
		"incorrect level should not exist": {
			input:      "incorrect",
			expected:   "",
			expectedOk: false,
		},
		"emtpy level should not exist": {
			input:      "",
			expected:   "",
			expectedOk: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			level, ok := parser.ParseLevel(test.input)

			if ok != test.expectedOk {
				t.Fatalf("%s: error in OK behavior: expected %t - got %t", name, test.expectedOk, ok)
			}
			if level != test.expected {
				t.Fatalf("%s: expected %s - got %s", name, test.expected, level)
			}
		})
	}
}
