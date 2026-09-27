package reader

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/utils"
	"github.com/google/go-cmp/cmp"
)

type cancellableReader struct {
	reader strings.Reader
	ctx    context.Context
	cancel func()
}

func (r *cancellableReader) Read(p []byte) (int, error) {
	i, err := r.reader.Read(p)
	if err != nil && errors.Is(err, io.EOF) {
		r.cancel()
	}
	return i, err
}

type expectedOutput struct {
	Logs []*parser.Log
	Errs []error
}

func TestExtractLinesFromReader(t *testing.T) {
	tests := map[string]struct {
		inputReader        *strings.Reader
		inputParseFunction parser.ParseFunction
		expected           expectedOutput
	}{
		"empty reader should not return any log nor error": {
			inputReader:        strings.NewReader(""),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{},
				Errs: []error{},
			},
		},
		"reader with one line should return one log": {
			inputReader:        strings.NewReader("test line"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{}},
				Errs: []error{},
			},
		},
		"reader with two lines should return two logs": {
			inputReader:        strings.NewReader("test line\ntest new line"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{}, {}},
				Errs: []error{},
			},
		},
		"parse function that returns an error should feed the errChan with parseError": {
			inputReader:        strings.NewReader("test line"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, &parser.ValueError{} },
			expected: expectedOutput{
				Logs: []*parser.Log{},
				Errs: []error{&parser.ParseError{Err: &parser.ValueError{}, Line: 1}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			creader := cancellableReader{
				reader: *test.inputReader,
				ctx:    ctx,
				cancel: cancel,
			}
			logChan, errChan := extractLinesFromReader(creader.ctx, &creader, test.inputParseFunction)
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
