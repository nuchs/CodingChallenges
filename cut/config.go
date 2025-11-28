package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

var ErrBadSelector = errors.New("you must provide a field selector > 0")

type Config struct {
	Delimiter string
	Field     int
	In        string // '-' is for stdin
}

func NewConfig(args []string) (Config, error) {
	cfg := Config{In: "-", Delimiter: "\t"}
	var sb strings.Builder
	fs := flag.NewFlagSet("cut", flag.ContinueOnError)
	fs.SetOutput(&sb)
	fs.Usage = func() {
		sb.WriteString("wwcut - cut clone")
		fs.PrintDefaults()
	}

	fs.IntVar(
		&cfg.Field,
		"f",
		0,
		"Fields to select, can be a comma separated list or a range",
	)

	if err := fs.Parse(args); err != nil {
		return cfg, fmt.Errorf("parsing command line: %w", err)
	}
	if cfg.Field == 0 {
		return cfg, ErrBadSelector
	}

	tail := fs.Args()
	if len(tail) > 0 {
		cfg.In = tail[0]
	}

	return cfg, nil
}
