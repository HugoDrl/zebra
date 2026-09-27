package reader

import (
	"bufio"
	"context"
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

		scanner := bufio.NewScanner(reader)

		lineNo := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if ok := scanner.Scan(); !ok {
				if err := scanner.Err(); err != nil {
					errsChan <- err
					return
				}
				continue
			}
			lineNo++
			contentLine := scanner.Text()

			log, err := parseFunction(contentLine)
			if err != nil {
				errsChan <- &parser.ParseError{
					Line: lineNo,
					Err:  err,
				}
			} else {
				logsChan <- &log
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
