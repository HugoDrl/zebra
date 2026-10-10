package reader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
)

func extractLinesFromReader(ctx context.Context, reader io.Reader, parseFunction parser.ParseFunction) (<-chan *parser.Log, <-chan error) {
	logsChan := make(chan *parser.Log)
	errsChan := make(chan error)
	go func() {
		defer close(logsChan)
		defer close(errsChan)

		buffer := make([]byte, 1024)
		startIdx := 0
		lineNo := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, err := reader.Read(buffer[startIdx:])
			if err != nil {
				if errors.Is(err, io.EOF) {
					continue
				}
				errsChan <- err
				return
			}
			startIdx += n
			lineNo++
			for {
				idx := bytes.Index(buffer, []byte{'\n'})
				if idx == -1 {
					break
				}

				log, err := parseFunction(string(buffer[:idx]))
				buffer = buffer[idx+1:]
				if err != nil {
					errsChan <- &parser.ParseError{
						Line: lineNo,
						Err:  err,
					}
				} else {
					logsChan <- &log
				}
				startIdx = 0
			}
		}
	}()
	return logsChan, errsChan
}

type ExtractLinesFromFileInput struct {
	Ctx           context.Context
	Root          *os.Root
	Filename      string
	ParseFunction parser.ParseFunction
}

func ExtractLogsFromFileName(input ExtractLinesFromFileInput) (<-chan *parser.Log, <-chan error) {
	file, err := input.Root.Open(input.Filename)
	if err != nil {
		echan := make(chan error, 1)
		echan <- err
		close(echan)
		return nil, echan
	}

	return extractLinesFromReader(input.Ctx, file, input.ParseFunction)
}
