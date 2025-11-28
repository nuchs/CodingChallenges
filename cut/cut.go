package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Cut(in io.Reader, out io.Writer, cfg Config) error {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer flushBuffer(writer)

	for scanner.Scan() {
		if scanner.Err() != nil {
			return fmt.Errorf("reading from input: %w", scanner.Err())
		}
		line := scanner.Text()
		trimmed := fmt.Sprintln(applyCuts(cfg, line))
		if _, err := writer.WriteString(trimmed); err != nil {
			return fmt.Errorf("writing to output: %w", err)
		}
	}

	return nil
}

func applyCuts(cfg Config, line string) string {
	parts := strings.Split(line, cfg.Delimiter)

	if len(parts) == 1 {
		return parts[0]
	}
	if cfg.Field > len(parts) {
		return ""
	}
	return parts[cfg.Field-1]
}

func flushBuffer(buf *bufio.Writer) {
	if err := buf.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush output: %s\n", err)
	}
}
