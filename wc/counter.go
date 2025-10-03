package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Counts struct {
	Bytes int
	Runes int
	Words int
	Lines int
}

func (c *Counts) Format(spec Spec) string {
	var b strings.Builder
	b.WriteString(" ")
	if spec.Lines {
		fmt.Fprintf(&b, " %v", c.Lines)
	}
	if spec.Words {
		fmt.Fprintf(&b, " %v", c.Words)
	}
	if spec.MultiByte {
		fmt.Fprintf(&b, " %v", c.Runes)
	}
	if spec.Bytes {
		fmt.Fprintf(&b, " %v", c.Bytes)
	}

	return b.String()
}

func (c *Counts) update(s string) {
	c.Bytes += len(s)
	c.Runes += utf8.RuneCountInString(s)
	c.Words += len(strings.Fields(s))
	c.Lines++
}

func Count(_ Spec, file io.Reader) (Counts, error) {
	var counts Counts

	rd := bufio.NewReader(file)
	for {
		s, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return counts, fmt.Errorf("error reading file: %w", err)
			}
			if s != "" {
				counts.update(s)
			}
			break
		}

		counts.update(s)
	}

	return counts, nil
}
