package main

import (
	"fmt"
	"io"
)

func Encode(data io.Reader) ([]byte, error) {
	_, err := generateEncodingTable(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encoding table: %w", err)
	}

	return nil, nil
}
