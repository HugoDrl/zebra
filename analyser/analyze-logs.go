package analyser

import (
	"errors"
	"github.com/HugoDrl/LogParser/parser"
	"sync"
)

func AnalyseLogs(
	logChan <-chan *parser.Log,
	errChan <-chan error,
	settings *AnalyserSettings,
) *CollectionMetric {
	metrics := newMetrics()
	var wg sync.WaitGroup

	wg.Go(func() {
		for log := range logChan {
			metrics.Lines[log.Level]++
			metrics.handleService(log)
			metrics.handleSlowestLogs(settings.SlowestLogsToRetrieve, log)
		}
	})

	wg.Go(func() {
		for err := range errChan {
			var fileErr *parser.FileError
			if errors.As(err, &fileErr) {
				metrics.FileErrors = append(metrics.FileErrors, fileErr)
			} else {
				metrics.ParsingErrorCount++
			}
		}
	})

	wg.Wait()

	return metrics
}
