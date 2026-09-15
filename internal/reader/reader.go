package reader

import (
	"bufio"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
)

type ParseLinesInput struct {
	ParseFunction parser.ParseFunction
	LogsChan      chan<- *parser.Log
	ErrsChan      chan<- error
}

type extractLinesFromReaderInput struct {
	ParseLinesInput
	reader *bufio.Reader
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
	ParseLinesInput
	root     *os.Root
	Filename string
}

func ExtractLogsFromFileName(input ExtractLinesFromFileInput) {
	file, err := input.root.Open(input.Filename)
	if err != nil {
		input.ErrsChan <- err
		return
	}

	reader := bufio.NewReader(file)
	extractLinesFromReader(extractLinesFromReaderInput{
		reader: reader,
		ParseLinesInput: ParseLinesInput{
			ParseFunction: input.ParseLinesInput.ParseFunction,
			LogsChan:      input.LogsChan,
			ErrsChan:      input.ErrsChan,
		},
	})
}
