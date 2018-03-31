// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/opsidian/parsley/parsley"
)

// Reader defines a text input reader
// For more efficient reading it provides methods for regexp matching.
type Reader struct {
	file        *File
	regexpCache map[string]*regexp.Regexp
}

// NewReader creates a new reader instance
// The Windows-style line endings (\r\n) are automatically replaced with Unix-style line endings (\n).
func NewReader(file *File) *Reader {
	return &Reader{
		file:        file,
		regexpCache: map[string]*regexp.Regexp{},
	}
}

// ReadRune matches any of the given runes
func (r *Reader) ReadRune(pos int, runes ...rune) (int, rune, bool) { // nolint
	if pos >= r.file.len {
		return pos, 0, false
	}

	if len(runes) == 0 {
		panic("MatchRune() should not be called with an empty rune list")
	}

	for _, runeValue := range runes {
		if runeValue < utf8.RuneSelf {
			if int8(runeValue) == int8(r.file.data[pos]) {
				return pos + 1, runeValue, true
			}
		} else {
			nextRune, width := utf8.DecodeRune(r.file.data[pos:])
			if nextRune == runeValue {
				return pos + width, runeValue, true
			}
		}
	}

	return pos, 0, false
}

// MatchString matches the given string
func (r *Reader) MatchString(pos int, str string) (int, bool) {
	if str == "" {
		panic("MatchString() should not be called with an empty string")
	}

	if len(str) > len(r.file.data)-pos {
		return pos, false
	}

	if bytes.HasPrefix(r.file.data[pos:], []byte(str)) {
		return pos + len(str), true
	}
	return pos, false
}

// MatchWord matches the given word
// It's different from MatchString() as it checks that the next character is not a word character
func (r *Reader) MatchWord(pos int, word string) (int, bool) {
	if word == "" {
		panic("MatchWord() should not be called with an empty string")
	}

	if len(word) > len(r.file.data)-pos {
		return pos, false
	}

	for i, b := range []byte(word) {
		if b >= utf8.RuneSelf {
			panic("MatchWord() should not be called with UTF8 strings")
		}
		if b != r.file.data[pos+i] {
			return pos, false
		}
	}

	if len(r.file.data)-pos-len(word) == 0 || !isWordCharacter(r.file.data[pos+len(word)]) {
		return pos + len(word), true
	}
	return pos, false
}

// ReadRegexp matches part of the input based on the given regular expression
// and returns with the full match
func (r *Reader) ReadRegexp(pos int, expr string) (int, []byte) {
	if pos >= r.file.len {
		return pos, nil
	}

	indices := r.getPattern(expr).FindIndex(r.file.data[pos:])
	if indices == nil {
		return pos, nil
	}

	return pos + indices[1], r.file.data[pos : pos+indices[1]]
}

// ReadRegexpSubmatch matches part of the input based on the given regular expression
// and returns with all capturing groups
func (r *Reader) ReadRegexpSubmatch(pos int, expr string) (int, [][]byte) {
	if pos >= r.file.len {
		return pos, nil
	}

	matches := r.getPattern(expr).FindSubmatch(r.file.data[pos:])
	if matches == nil {
		return pos, nil
	}

	return pos + len(matches[0]), matches
}

// Readf uses the given function to match the next token
func (r *Reader) Readf(pos int, f func(b []byte) ([]byte, int)) (int, []byte) {
	if pos >= r.file.len {
		return pos, nil
	}

	value, nextPos := f(r.file.data[pos:])
	if nextPos == 0 {
		if value != nil {
			panic("no value should be returned if next position is zero")
		}
		return pos, nil
	}

	if nextPos < len(value) || pos+nextPos > r.file.len {
		panic("invalid length was returned by the custom reader function")
	}

	return pos + nextPos, value
}

// Remaining returns with the remaining character count
func (r *Reader) Remaining(pos int) int {
	return r.file.len - pos
}

// IsEOF returns true if we reached the end of the buffer
func (r *Reader) IsEOF(pos int) bool {
	return pos >= r.file.len
}

// MatchWhitespaces reads all the whitespaces
func (r *Reader) MatchWhitespaces(pos int) int {
	for pos < r.file.len && isWhitespaceCharacter(r.file.data[pos]) {
		pos++
	}
	return pos
}

// Pos returns with the current position
func (r *Reader) Pos(pos int) parsley.Pos {
	return r.file.Pos(pos)
}

func (r *Reader) getPattern(expr string) *regexp.Regexp {
	rc, ok := r.regexpCache[expr]
	if !ok {
		rc = regexp.MustCompile("^(?:" + expr + ")")

		if rc.Match(nil) {
			panic(fmt.Errorf("'%s' is not allowed to match an empty input", expr))
		}

		r.regexpCache[expr] = rc
	}
	return rc
}

func isWordCharacter(b byte) bool {
	return 'a' <= b && b <= 'z' ||
		'A' <= b && b <= 'Z' ||
		'0' <= b && b <= '9' ||
		b == '_'
}

func isWhitespaceCharacter(b byte) bool {
	return b == '\t' || b == '\n' || b == '\f' || b == ' '
}
