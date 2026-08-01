package analyser_test

import (
	"testing"

	"github.com/HugoDrl/LogParser/analyser"
	"github.com/HugoDrl/LogParser/parser"
	"github.com/google/go-cmp/cmp"
)

func feedThenCloseChan[T any](values ...T) chan T {
	channel := make(chan T, len(values))
	for _, value := range values {
		channel <- value
	}
	close(channel)
	return channel
}

func TestLogAnalyze(t *testing.T) {
	type LogAnalyzeInput struct {
		logChan  chan *parser.Log
		errChan  chan error
		settings *analyser.AnalyserSettings
	}

	tests := map[string]struct {
		input    LogAnalyzeInput
		expected *analyser.CollectionMetric
	}{
		"empty everything should result in empty metrics": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{}...),
				errChan: feedThenCloseChan([]error{}...),
			},
			expected: &analyser.CollectionMetric{
				Lines:              map[parser.Level]int{},
				ServicePerformance: map[string]analyser.ServiceMetric{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := analyser.AnalyseLogs(test.input.logChan, test.input.errChan, test.input.settings)
			if diff := cmp.Diff(
				test.expected,
				output,
			); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
