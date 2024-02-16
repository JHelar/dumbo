package lex

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

type runeReader struct {
	currentRune rune
	reader      *bufio.Reader
}

func newRuneReaderFromFile(file *os.File) *runeReader {
	return &runeReader{
		reader: bufio.NewReader(file),
	}
}

func newRuneReader(content []byte) *runeReader {
	return &runeReader{
		reader: bufio.NewReader(bytes.NewReader(content)),
	}
}

func (rr *runeReader) Next() (rune, error) {
	if rr.currentRune != 0 {
		r := rr.currentRune
		rr.currentRune = 0

		return r, nil
	}

	r, _, err := rr.reader.ReadRune()

	if err != nil {
		return 0, err
	}

	return r, nil
}

func (rr *runeReader) MustNext() rune {
	r, err := rr.Next()

	if err != nil {
		panic(err)
	}

	return r
}

func (rr *runeReader) Peak() (rune, error) {
	if rr.currentRune != 0 {
		return rr.currentRune, nil
	}

	r, _, err := rr.reader.ReadRune()

	if err != nil {
		log.Printf("Error peaking rune: %s", err.Error())
		return 0, err
	}

	rr.currentRune = r
	return r, nil
}
