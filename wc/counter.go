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
	Src   string
}

func (c *Counts) Format(spec Spec, fw FieldWidths) string {
	var b strings.Builder
	if spec.Lines {
		fmt.Fprintf(&b, "%*d", fw.Line, c.Lines)
	}
	if spec.Words {
		fmt.Fprintf(&b, "%*d", fw.Word, c.Words)
	}
	if spec.MultiByte {
		fmt.Fprintf(&b, "%*d", fw.Rune, c.Runes)
	}
	if spec.Bytes {
		fmt.Fprintf(&b, "%*d", fw.Byte, c.Bytes)
	}
	fmt.Fprintf(&b, " %s", c.Src)

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

func Count(file io.Reader, src string) (Counts, error) {
	counts := Counts{Src: src}

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
