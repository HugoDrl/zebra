package reader

import (
	"bufio"
	"io"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
)

type extractLinesFromReaderInput struct {
	ParseFunction parser.ParseFunction
	LogsChan      chan<- *parser.Log
	ErrsChan      chan<- error
	reader        io.Reader
}

func extractLinesFromReader(input extractLinesFromReaderInput) {
	scanner := bufio.NewScanner(input.reader)

	lineNo := 0
	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				input.ErrsChan <- err
			}
			return
		}
		lineNo++
		contentLine := scanner.Text()

		log, err := input.ParseFunction(contentLine)
		if err != nil {
			input.ErrsChan <- &parser.ParseError{
				Line: lineNo,
				Err:  err,
			}
		} else {
			input.LogsChan <- &log
		}
	}
}

type ExtractLinesFromFileInput struct {
	Root          *os.Root
	Filename      string
	ParseFunction parser.ParseFunction
}

func ExtractLogsFromFileName(input ExtractLinesFromFileInput) (<-chan *parser.Log, <-chan error) {
	logsChan := make(chan *parser.Log)
	errsChan := make(chan error)

	go func() {
		defer close(logsChan)
		defer close(errsChan)

		file, err := input.Root.Open(input.Filename)
		if err != nil {
			errsChan <- err
			return
		}

		extractLinesFromReader(extractLinesFromReaderInput{
			reader:        bufio.NewReader(file),
			ParseFunction: input.ParseFunction,
			LogsChan:      logsChan,
			ErrsChan:      errsChan,
		})
	}()
	return logsChan, errsChan
}
