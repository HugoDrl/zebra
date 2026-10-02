package store

import (
	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
)

type LogsStore interface {
	RetrieveLogs(filter.Filters) ([]*parser.Log, []error, error)
	InsertLogs(<-chan *parser.Log, <-chan error) ([]*parser.Log, error)
}
