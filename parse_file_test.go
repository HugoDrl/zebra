package main_test

import (
	"os"
	"sort"
	"testing"
	"time"

	zebra "github.com/HugoDrl/zebra"
	"github.com/HugoDrl/zebra/parser"
	"github.com/google/go-cmp/cmp"
)

func emptyChan[T any](channel chan T) []T {
	array := make([]T, 0)
	for value := range channel {
		array = append(array, value)
	}
	if len(array) == 0 {
		return nil
	}
	return array
}

func TestParseLogsFromFile(t *testing.T) {
	type expectedStruct struct {
		Logs []*parser.Log
		Errs []error
	}
	tests := map[string]struct {
		inputFileContent   []string
		inputParseSettings parser.ParseSettings
		expected           expectedStruct
	}{
		"empty file should return nothing": {
			inputFileContent: []string{""},
			inputParseSettings: parser.ParseSettings{
				Files: []string{"temp_parse_logs_from_file.txt"},
			},
		},
		"file with a valid log should parse it": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{"temp_parse_logs_from_file.txt"},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
			},
		},
		"several files with several valid logs should parse them all": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
				"2026-01-01T00:00:10Z INFO service=database message=\"hello again from database\" duration=12ms",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file1.txt",
					"temp_parse_logs_from_file2.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 10, 0, time.UTC),
						Level:    parser.Info,
						Service:  "database",
						Message:  `"hello again from database"`,
						Duration: parser.Duration(12 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
			},
		},
		"several files, but one has invalid content should return both logs and errors": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
				"invalid content for logs",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file1.txt",
					"temp_parse_logs_from_file2.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{&parser.ValueError{ErroredValue: "invalid", ExpectedValue: "time format - RFC3339"}},
			},
		},
		"one file with several logs should parse them all": {
			inputFileContent: []string{
				`2026-01-01T00:00:00Z WARNING service=database message="hello from database" duration=10ms
2026-01-01T00:00:10Z INFO service=database message="hello again from database" duration=12ms`,
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file_parse_all.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 10, 0, time.UTC),
						Level:    parser.Info,
						Service:  "database",
						Message:  `"hello again from database"`,
						Duration: parser.Duration(12 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Preparation of file (with content written) and channels
			for i, fileName := range test.inputParseSettings.Files {
				err := os.WriteFile(fileName, []byte(test.inputFileContent[i]), 0o777)
				if err != nil {
					t.Fatal(err)
				}
			}
			// Channels are 100 values max.
			// Because current tests architecture implies to fill channels before closing and reading,
			// This means current tests should not be more than 100 values for each
			logChan := make(chan *parser.Log, 100)
			errsChan := make(chan error, 100)

			// Actual tested function
			zebra.ProcessFiles(&test.inputParseSettings, logChan, errsChan)

			// Clean files
			for _, file := range test.inputParseSettings.Files {
				os.Remove(file)
			}
			// Retrieve channels informations and test
			logs := emptyChan(logChan)
			errs := emptyChan(errsChan)

			// Sort logs by time to be predictible for comparison
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].Time.Compare(logs[j].Time) < 0
			})

			output := expectedStruct{
				Logs: logs,
				Errs: errs,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
