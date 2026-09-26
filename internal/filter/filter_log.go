package filter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
)

func filterLog(log *parser.Log, filters *Filters) bool {
	if !filters.StartDate.IsZero() && log.Time.Compare(filters.StartDate) < 0 {
		return false
	}
	if !filters.EndDate.IsZero() && log.Time.Compare(filters.EndDate) > 0 {
		return false
	}
	if filters.Level != "" && log.Level != filters.Level {
		return false
	}
	if filters.Service != "" && log.Service != filters.Service {
		return false
	}
	return true
}

func FilterLogs(logs []*parser.Log, filters Filters) []*parser.Log {
	filteredLogs := make([]*parser.Log, 0, len(logs))

	for _, log := range logs {
		if filterLog(log, &filters) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs
}

func ProcessRequestToFilter(r http.Request) (Filters, error) {
	filters := Filters{}
	for key := range r.URL.Query() {
		value := r.URL.Query().Get(key)
		switch key {
		case "start-date":
			startDate, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return Filters{}, err
			}
			filters.StartDate = startDate
		case "end-date":
			endDate, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return Filters{}, err
			}
			filters.EndDate = endDate
		case "level":
			level, ok := parser.ParseLevel(value)
			if !ok {
				return Filters{}, fmt.Errorf("'%s' is not a valid level", value)
			}
			filters.Level = level
		case "service":
			filters.Service = value
		}

	}
	return filters, nil
}
