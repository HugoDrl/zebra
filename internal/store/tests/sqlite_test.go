package store_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestSQLiteStore(t *testing.T) {
	tests := map[string]struct {
		input        []*parser.Log
		inputFilters filter.Filters
		expected     []*parser.Log
	}{
		"nil as input should return empty slice": {
			input:        nil,
			inputFilters: filter.Filters{},
			expected:     []*parser.Log{},
		},
		"nil as input with filters should return empty slice": {
			input: nil,
			inputFilters: filter.Filters{
				StartDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []*parser.Log{},
		},
		"emtpy slice as input should return empty slice": {
			input:        []*parser.Log{},
			inputFilters: filter.Filters{},
			expected:     []*parser.Log{},
		},
		"emtpy slice as input with filters should return empty slice": {
			input: []*parser.Log{},
			inputFilters: filter.Filters{
				StartDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []*parser.Log{},
		},
		"single log with no filters should return the log": {
			input:        []*parser.Log{{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}},
			inputFilters: filter.Filters{},
			expected:     []*parser.Log{{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}},
		},
		"single log not concerned by start date filter should not be returned": {
			input:        []*parser.Log{{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}},
			inputFilters: filter.Filters{StartDate: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
			expected:     []*parser.Log{},
		},
		"single log not concerned by end date filter should not be returned": {
			input:        []*parser.Log{{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}},
			inputFilters: filter.Filters{EndDate: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
			expected:     []*parser.Log{},
		},
		"single log not concerned by level filter should not be returned": {
			input:        []*parser.Log{{Level: "test"}},
			inputFilters: filter.Filters{Level: "not-test"},
			expected:     []*parser.Log{},
		},
		"single log not concerned by service filter should not be returned": {
			input:        []*parser.Log{{Service: "test"}},
			inputFilters: filter.Filters{Service: "not-test"},
			expected:     []*parser.Log{},
		},
		// TODO: actually test (and implement) extras
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dbName := fmt.Sprintf("./%.6f.db", rand.ExpFloat64())
			defer os.Remove(dbName)
			store, err := store.NewSQLiteStore(dbName)
			if err != nil {
				t.Fatal(err)
				return
			}
			ctx := context.Background()

			// send all logs to a channel
			logChan := make(chan *parser.Log)
			done := make(chan error)
			go func() {
				for _, log := range test.input {
					logChan <- log
				}
				close(logChan)
			}()

			// use that channel to send
			go func() {
				_, err := store.InsertLogs(ctx, logChan, nil)
				done <- err
				close(done)
			}()

			if err := <-done; err != nil {
				t.Fatal(err)
			}
			// retrieve logs just inserted with filters
			logs, _, err := store.RetrieveLogs(ctx, test.inputFilters)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.expected, logs); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
