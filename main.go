package main

import (
	"fmt"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/server"
)

func main() {
	var logs []*parser.Log
	var errs []error
	if err := server.NewServer(
		&server.DataLayer{
			Logs: logs,
			Errs: errs,
		},
	).StartServer(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
