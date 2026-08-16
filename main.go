package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HugoDrl/zebra/analyser"
	"github.com/HugoDrl/zebra/parser"
)

func extractLinesFromFile(reader *bufio.Reader, json bool, outChan chan<- *parser.Log, errsChan chan<- error) {
	scanner := bufio.NewScanner(reader)

	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				errsChan <- err
			}
			return
		}
		contentLine := scanner.Text()

		var log parser.Log
		var err error

		if json {
			log, err = parser.ParseJSONFormatLine(string(contentLine))
		} else {
			log, err = parser.ParseDefaultFormatLine(string(contentLine))
		}
		if err != nil {
			errsChan <- err
		} else {
			outChan <- &log
		}
	}
}

func ProcessFiles(
	settings *parser.ParseSettings,
	outChan chan<- *parser.Log,
	errsChan chan<- error,
	wg *sync.WaitGroup,
) {
	for _, filepath := range settings.Files {
		wg.Go(func() {
			reader, err := os.OpenFile(filepath, os.O_RDONLY, 0)
			if err != nil {
				errsChan <- err
				return
			}
			r := bufio.NewReader(reader)
			extractLinesFromFile(r, settings.Json, outChan, errsChan)
		})
	}
}

func initSettings() (*parser.ParseSettings, *analyser.AnalyserSettings, error) {
	files := flag.String("files", "", "log files to analyse")
	json := flag.Bool("json", false, "wether or not format to parse is json")
	startDate := flag.String("start", "", "logs date to start from")
	endDate := flag.String("end", "", "logs date to end to")
	service := flag.String("service", "", "filter logs by service")
	level := flag.String("level", "", "filter logs by level")
	slowestLogs := flag.Int("top", 0, "number of slowest logs to show")
	flag.Parse()

	if *files == "" {
		return nil, nil, errors.New("Please specify file(s) separated by a comma using --files flag")
	}

	var processedStartDate time.Time
	var processedEndDate time.Time
	var processErr error
	if *startDate != "" {
		processedStartDate, processErr = time.Parse(time.RFC3339, *startDate)
		if processErr != nil {
			return nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}
	if *endDate != "" {
		processedEndDate, processErr = time.Parse(time.RFC3339, *endDate)
		if processErr != nil {
			return nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}

	parsingSettings := parser.ParseSettings{
		Files:     strings.Split(*files, ","),
		Json:      *json,
		StartDate: processedStartDate,
		EndDate:   processedEndDate,
		Level:     parser.Level(*level),
		Service:   *service,
	}
	analyserSettings := analyser.AnalyserSettings{
		SlowestLogsToRetrieve: *slowestLogs,
	}
	return &parsingSettings, &analyserSettings, nil
}

func main() {
	parsingSettings, analyserSettings, err := initSettings()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
		os.Exit(1)
	}

	out := make(chan *parser.Log)
	errs := make(chan error)
	var wgFiles sync.WaitGroup
	ProcessFiles(parsingSettings, out, errs, &wgFiles)

	go func() {
		wgFiles.Wait()
		close(out)
		close(errs)
	}()

	metrics := analyser.AnalyseLogs(out, errs, analyserSettings)
	metrics.Display()
}
