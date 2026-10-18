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

type partialReader struct {
	Content      string
	BytesPerRead int
	currIndex    int
}

func (r *partialReader) Read(p []byte) (int, error) {
	if r.currIndex >= len(r.Content) {
		return 0, io.EOF
	}

	leaveToRead := len(r.Content) - r.currIndex
	willRead := min(r.BytesPerRead, leaveToRead)

	copy(p[:willRead], []byte(r.Content[r.currIndex:r.currIndex+willRead]))
	r.currIndex += willRead

	return willRead, nil
}

type cancellableReader struct {
	reader io.Reader
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
		inputReader        io.Reader
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
			inputReader:        strings.NewReader("test line\n"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{}},
				Errs: []error{},
			},
		},
		"reader with two lines should return two logs": {
			inputReader:        strings.NewReader("test line\ntest new line\n"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{}, {}},
				Errs: []error{},
			},
		},
		"parse function that returns an error should feed the errChan with parseError": {
			inputReader:        strings.NewReader("test line\n"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{}, &parser.ValueError{} },
			expected: expectedOutput{
				Logs: []*parser.Log{},
				Errs: []error{&parser.ParseError{Err: &parser.ValueError{}, Line: 1}},
			},
		},
		"truncated last log should not be considered": {
			inputReader:        strings.NewReader("test line\nhey i am"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{Message: s}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{Message: "test line"}},
				Errs: []error{},
			},
		},
		"any log terminated with newline should be considered done": {
			inputReader:        strings.NewReader("test line\nhey i am\n"),
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{Message: s}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{Message: "test line"}, {Message: "hey i am"}},
				Errs: []error{},
			},
		},
		"partial read of a log should be just fine": {
			inputReader: &partialReader{
				Content:      "test line\nhey i am\n",
				BytesPerRead: 3,
			},
			inputParseFunction: func(s string) (parser.Log, error) { return parser.Log{Message: s}, nil },
			expected: expectedOutput{
				Logs: []*parser.Log{{Message: "test line"}, {Message: "hey i am"}},
				Errs: []error{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			creader := cancellableReader{
				reader: test.inputReader,
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
