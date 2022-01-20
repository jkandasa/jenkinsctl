package utils

import (
	"bufio"
	"io"
	"os"
)

type lineIterator struct {
	reader *bufio.Reader
}

func newLineIterator(rd io.Reader) *lineIterator {
	return &lineIterator{
		reader: bufio.NewReader(rd),
	}
}

func (ln *lineIterator) next() ([]byte, error) {
	var bytes []byte
	for {
		line, isPrefix, err := ln.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, line...)
		if !isPrefix {
			break
		}
	}
	return bytes, nil
}

func ReadStdIn() ([]byte, error) {
	ln := newLineIterator(os.Stdin)
	lines := make([]byte, 0)
	for {
		line, err := ln.next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		lines = append(lines, line...)
		lines = append(lines, '\n')
	}
	return []byte(lines), nil
}
