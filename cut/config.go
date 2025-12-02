package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

var ErrBadSelector = errors.New("you must provide a field selector")

type Config struct {
	Delimiter string
	Spread    Spreads
	In        string // '-' is for stdin
}

func NewConfig(args []string) (Config, error) {
	cfg := Config{In: "-"}
	var sb strings.Builder
	fs := flag.NewFlagSet("cut", flag.ContinueOnError)
	fs.SetOutput(&sb)
	fs.Usage = func() {
		sb.WriteString("wwcut - cut clone")
		fs.PrintDefaults()
	}

	rStr := fs.String(
		"f",
		"",
		"Fields to select, can be a comma separated list or a range",
	)
	fs.StringVar(&cfg.Delimiter, "d", "\t", "Field delimiter")

	if err := fs.Parse(args); err != nil {
		return cfg, fmt.Errorf("parsing command line: %w", err)
	}
	if *rStr == "" {
		return cfg, ErrBadSelector
	}
	r, err := NewSpreads(*rStr)
	if err != nil {
		return cfg, err
	}
	cfg.Spread = r

	tail := fs.Args()
	if len(tail) > 0 {
		cfg.In = tail[0]
	}

	return cfg, nil
}
