// Package xgbgen generates protocols for X directly from X header files
package xgbgen

import (
	"bufio"
	"fmt"
	"io"
)

type tokenType int

// token holds information about a scanned token
type token struct {
	tokenType
	literal string
}

// tokenTypes
const (
	tknError tokenType = iota
	tknEOF
	tknKeyword
	tknIdentifier
	tknConstant
	tknLiteral
	tknSymbol
	tknText
)

var tokenMap = map[tokenType]string{
	tknError:      "Error",
	tknEOF:        "EOF",
	tknKeyword:    "Keyword",
	tknIdentifier: "Identifier",
	tknConstant:   "Constant",
	tknLiteral:    "Literal",
	tknSymbol:     "Symbol",
}

// stateFn tracks state through lexical analysis
type stateFn func(*lexer) stateFn

// lexer ...
type lexer struct {
	name   string
	reader *bufio.Reader
	state  stateFn
	tokens chan token
}

const (
	// read delimeters
	whitespace byte = ' '
)

func (l *lexer) read(delim byte) (string, error) {
	return l.reader.ReadString(delim)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tknError, fmt.Sprintf(format, args...)}
	return nil
}

// lex returns a new lexer, given an io.Reader
func lex(name string, rdr io.Reader) (l *lexer) {
	l = &lexer{
		name:   name,
		reader: bufio.NewReader(rdr),
		state:  lexText,
		tokens: make(chan token, 2),
	}
	return l
}

const (
// keywords
)

// lexText ...
func lexText(l *lexer) stateFn {
	text, err := l.read(whitespace)
	if err != nil {
		l.errorf("Error reading while lexing text, %s", err)
	}

	return nil
}
