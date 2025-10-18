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

func (c *Counts) Format(spec Spec, src string) string {
	var b strings.Builder
	if spec.Lines {
		fmt.Fprintf(&b, "%6d", c.Lines)
	}
	if spec.Words {
		fmt.Fprintf(&b, "%6d", c.Words)
	}
	if spec.MultiByte {
		fmt.Fprintf(&b, "%6d", c.Runes)
	}
	if spec.Bytes {
		fmt.Fprintf(&b, "%6d", c.Bytes)
	}
	fmt.Fprintf(&b, " %s", src)

	return b.String()
}

func (c *Counts) update(s string) {
	length := len(s)
	c.Bytes += length
	c.Runes += utf8.RuneCountInString(s)
	c.Words += len(strings.Fields(s))
	if length > 0 && s[length-1] == '\n' {
		c.Lines++
	}
}

func Count(file io.Reader) (Counts, error) {
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
