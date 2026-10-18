package reader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"time"

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
		delay := 0 * time.Second
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
			n, err := reader.Read(buffer[startIdx:])
			if err != nil {
				if errors.Is(err, io.EOF) {
					// exponential increase of delay to avoid charging CPU on read file
					if delay.Seconds() <= 15 {
						delay *= 2
					}
					if delay.Seconds() == 0 {
						delay = 1 * time.Second
					}
					continue
				}
				errsChan <- err
				return
			}
			// reset of the delay when some new lines are found
			delay = 0 * time.Second
			startIdx += n
			lineNo++
			for {
				idx := bytes.Index(buffer, []byte{'\n'})
				if idx == -1 {
					break
				}

				log, err := parseFunction(string(buffer[:idx]))
				newBuffer := make([]byte, len(buffer))
				copy(newBuffer, buffer[idx+1:])
				buffer = newBuffer
				if err != nil {
					errsChan <- &parser.ParseError{
						Line: lineNo,
						Err:  err,
					}
				} else {
					logsChan <- &log
				}
				// here, no verification is made on startIdx because we know buffer will never be full,
				// so we will always find some emtpy values. This is because this code executes only
				// if we do find some new line, so buffer will be partly emptied
				startIdx = bytes.Index(buffer, []byte{0})
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
