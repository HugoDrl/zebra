package store

import (
	"context"

	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
)

type LogsStore interface {
	RetrieveLogs(context.Context, filter.Filters) ([]*parser.Log, []error, error)
	InsertLogs(context.Context, <-chan *parser.Log, <-chan error) ([]*parser.Log, error)
}
