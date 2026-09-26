package reader

import (
	"strings"
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/utils"
	"github.com/google/go-cmp/cmp"
)

type expectedOutput struct {
	Logs []*parser.Log
	Errs []error
}

func TestExtractLinesFromReader(t *testing.T) {
	tests := map[string]struct {
		input    extractLinesFromReaderInput
		expected expectedOutput
	}{
		"empty reader should not return any log nor error": {
			input: extractLinesFromReaderInput{
				ParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
				reader:        strings.NewReader(""),
			},
			expected: expectedOutput{
				Logs: []*parser.Log{},
				Errs: []error{},
			},
		},
		"reader with one line should return one log": {
			input: extractLinesFromReaderInput{
				ParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
				reader:        strings.NewReader("test line"),
			},
			expected: expectedOutput{
				Logs: []*parser.Log{{}},
				Errs: []error{},
			},
		},
		"reader with two lines should return two logs": {
			input: extractLinesFromReaderInput{
				ParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
				reader:        strings.NewReader("test line\ntest new line"),
			},
			expected: expectedOutput{
				Logs: []*parser.Log{{}, {}},
				Errs: []error{},
			},
		},
		"parse function that returns an error should feed the errChan with parseError": {
			input: extractLinesFromReaderInput{
				ParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, &parser.ValueError{} },
				reader:        strings.NewReader("test line"),
			},
			expected: expectedOutput{
				Logs: []*parser.Log{},
				Errs: []error{&parser.ParseError{Err: &parser.ValueError{}, Line: 1}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			logChan, errChan := extractLinesFromReader(test.input)
			logsSlice, errsSlice := utils.ExtractLogAndErrChanToSlices(logChan, errChan)
			output := expectedOutput{
				Logs: logsSlice,
				Errs: errsSlice,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
