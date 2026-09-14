package analyser

import (
	"errors"
	"slices"
	"sync"

	"github.com/HugoDrl/zebra/internal/parser"
)

func AnalyseLogs(
	logs []*parser.Log,
	errs []error,
	settings *AnalyserSettings,
) *CollectionMetric {
	metrics := newMetrics()
	var wg sync.WaitGroup

	for _, log := range logs {
		metrics.handleService(log)
	}

	for _, err := range errs {
		var fileErr *parser.FileError
		if errors.As(err, &fileErr) {
			metrics.FileErrors = append(metrics.FileErrors, fileErr)
		} else {
			metrics.ParsingErrorCount++
		}
	}

	wg.Wait()

	return metrics
}

func RetrieveSlowestLogs(logs []*parser.Log, numberOfSlowestLogs int) []*parser.Log {
	returningLogs := make([]*parser.Log, len(logs))
	copy(returningLogs, logs)

	slices.SortFunc(returningLogs, func(leftLog *parser.Log, rightLog *parser.Log) int {
		return int(rightLog.Duration - leftLog.Duration)
	})

	return returningLogs[:min(numberOfSlowestLogs, len(returningLogs))]
}
