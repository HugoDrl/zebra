package parser_test

import (
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestGetParseFunction(t *testing.T) {
	tests := map[string]struct {
		input    parser.ParseSettings
		expected parser.ParseFunction
	}{
		"default input should return default parse function": {
			input:    parser.ParseSettings{},
			expected: parser.ParseDefaultFormatLine,
		},
		"explicit false json setting should return default parse function": {
			input: parser.ParseSettings{
				Json: false,
			},
			expected: parser.ParseDefaultFormatLine,
		},
		"explicit true on json field should return json function": {
			input: parser.ParseSettings{
				Json: true,
			},
			expected: parser.ParseJSONFormatLine,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out := parser.GetParseFunction(test.input)
			if cmp.Equal(test.expected, out) {
				t.Fatalf("expected %v - got %v", test.expected, out)
			}
		})
	}
}
