package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	ctx, cancel := signal.NotifyContext(ctx, os.Kill)
	defer cancel()

	cfg, err := NewConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	Run(ctx, cfg, os.Stdout)
}
