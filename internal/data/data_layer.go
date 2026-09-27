package data

import (
	"context"

	"github.com/HugoDrl/zebra/internal/parser"
)

type DataLayer struct {
	Ctx  context.Context
	Logs []*parser.Log
	Errs []error
}

func NewDataLayer() (*DataLayer, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &DataLayer{
		Ctx:  ctx,
		Logs: make([]*parser.Log, 0),
		Errs: make([]error, 0),
	}
	return d, cancel
}

func (d *DataLayer) IngestLogChan(c <-chan *parser.Log) {
	for log := range c {
		d.Logs = append(d.Logs, log)
	}
}

func (d *DataLayer) IngestErrChan(c <-chan error) {
	for err := range c {
		d.Errs = append(d.Errs, err)
	}
}
