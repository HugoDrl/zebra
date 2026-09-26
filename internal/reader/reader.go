package reader

import (
	"bufio"
	"io"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
)

type extractLinesFromReaderInput struct {
	ParseFunction parser.ParseFunction
	reader        io.Reader
}

func extractLinesFromReader(input extractLinesFromReaderInput) (<-chan *parser.Log, <-chan error) {
	logsChan := make(chan *parser.Log)
	errsChan := make(chan error)
	go func() {
		defer close(logsChan)
		defer close(errsChan)

		scanner := bufio.NewScanner(input.reader)

		lineNo := 0
		for {
			if ok := scanner.Scan(); !ok {
				if err := scanner.Err(); err != nil {
					errsChan <- err
				}
				return
			}
			lineNo++
			contentLine := scanner.Text()

			log, err := input.ParseFunction(contentLine)
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

	return extractLinesFromReader(extractLinesFromReaderInput{
		reader:        bufio.NewReader(file),
		ParseFunction: input.ParseFunction,
	})
}
