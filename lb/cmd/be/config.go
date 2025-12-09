package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const usage = `
be - stub backend service for testing lb

Usage:
  be -n <name> -p <port>

Options:
`

var (
	ErrNoName      = errors.New("you must specify a name for the service")
	ErrInvalidPort = errors.New("you must specify a port in the range 0-65535")
)

type Config struct {
	Name string
	Port int
}

func NewConfig(args []string) (Config, error) {
	var cfg Config
	var help strings.Builder
	parser := newParser(&cfg, help)

	err := parser.Parse(args)
	if err == flag.ErrHelp {
		return cfg, errors.New(help.String())
	}
	err = validateConfig(cfg, err)

	return cfg, err
}

func newParser(cfg *Config, help strings.Builder) *flag.FlagSet {
	parser := flag.NewFlagSet("be", flag.ContinueOnError)
	parser.SetOutput(&help)
	parser.Usage = func() {
		fmt.Fprint(&help, usage)
		parser.PrintDefaults()
	}

	parser.IntVar(&cfg.Port, "p", -1, "port service is listening on")
	parser.StringVar(
		&cfg.Name,
		"n",
		"",
		"name of service, used in logging and responses",
	)

	return parser
}

func validateConfig(cfg Config, err error) error {
	var errs []error

	if err != nil {
		errs = append(errs, fmt.Errorf("parsing command line:\n\t%w", err))
	}
	if cfg.Name == "" {
		errs = append(errs, ErrNoName)
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		errs = append(errs, ErrInvalidPort)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
