package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"
)

const usage = `
lb - simple http load balancer

Usage:
	lb -p <port> <url urls...>

Optionns:
	`

var (
	ErrInvalidPort = errors.New("you must specify a port in the range 0-65535")
	ErrNoUrls      = errors.New("you must specify at least one backend url")
	ErrBadProtocol = errors.New("backend server urls must be over http")
	ErrMissingHost = errors.New("url does not contain a host")
)

type Config struct {
	Port int
	Urls []string
}

func NewConfig(args []string) (Config, error) {
	var help strings.Builder
	var cfg Config
	parser := newParser(&cfg, help)

	err := parser.Parse(args)
	if err == flag.ErrHelp {
		return cfg, errors.New(help.String())
	}
	cfg.Urls = parser.Args()
	err = validateConfig(cfg, err)

	return cfg, err
}

func newParser(cfg *Config, help strings.Builder) *flag.FlagSet {
	parser := flag.NewFlagSet("lb", flag.ContinueOnError)
	parser.SetOutput(&help)
	parser.Usage = func() {
		fmt.Fprint(&help, usage)
		parser.PrintDefaults()
	}

	parser.IntVar(&cfg.Port, "p", -1, "port service is listening on")
	return parser
}

func validateConfig(cfg Config, err error) error {
	var errs []error

	if err != nil {
		errs = append(errs, fmt.Errorf("parsing command line: %w", err))
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		errs = append(errs, ErrInvalidPort)
	}
	if len(cfg.Urls) == 0 {
		errs = append(errs, ErrNoUrls)
	}
	for _, uStr := range cfg.Urls {
		u, err := url.Parse(uStr)
		if err != nil {
			errs = append(errs, fmt.Errorf("bad backend url: %w", err))
			continue
		}
		if u.Scheme != "http" {
			errs = append(errs, ErrBadProtocol)
			continue
		}
		if u.Hostname() == "" {
			errs = append(errs, ErrMissingHost)
			continue
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
