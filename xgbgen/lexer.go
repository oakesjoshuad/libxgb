// Package xgbgen generates protocols for X directly from X header files
package xgbgen

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type tokenType int

// token holds information about a scanned token
type token struct {
	Type    tokenType
	Literal string
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
	tknLeftComment
	tknRightComment
)

var tokenMap = map[tokenType]string{
	tknError:        "Error",
	tknEOF:          "EOF",
	tknKeyword:      "Keyword",
	tknIdentifier:   "Identifier",
	tknConstant:     "Constant",
	tknLiteral:      "Literal",
	tknSymbol:       "Symbol",
	tknText:         "Text",
	tknLeftComment:  leftComment,
	tknRightComment: rightComment,
}

func (tkn *token) String() string {
	return fmt.Sprintf("{ %s : %s }", tokenMap[tkn.Type], tkn.Literal)
}

// stateFn tracks state through lexical analysis
type stateFn func(*lexer) stateFn

// lexer ...
type lexer struct {
	name    string
	scanner *bufio.Scanner
	state   stateFn
	tokens  chan token
}

// scan ...
func (l *lexer) scan() bool {
	return l.scanner.Scan()
}

func (l *lexer) text() string {
	return l.scanner.Text()
}

// emit puts a token on the l.tokens channel
func (l *lexer) emit(t tokenType, i interface{}) {
	switch i := i.(type) {
	case []string:
		if len(i) > 0 {
			l.tokens <- token{t, strings.Join(i, whitespace)}
		}
	case string:
		l.tokens <- token{t, i}
	}
}

func (l *lexer) err() error {
	return l.scanner.Err()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{Type: tknError, Literal: fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) nextToken() token {
	for {
		select {
		case tkn := <-l.tokens:
			return tkn
		default:
			l.state = l.state(l)
		}
	}
}

var stateMap map[string]stateFn

// lex returns a new lexer, given an io.Reader
func lex(name string, rdr io.Reader) *lexer {
	l := &lexer{
		name:    name,
		scanner: bufio.NewScanner(rdr),
		state:   lexText,
		tokens:  make(chan token, 2),
	}

	// initialize the state map
	stateMap = map[string]stateFn{
		leftComment:  lexLeftComment,
		rightComment: lexRightComment,
		typedef:      lexTypedef,
		structure:    lexStructure,
		pd_if:        lexDirective,
		pd_define:    lexDirective,
		pd_else:      lexDirective,
		pd_endif:     lexDirective,
	}

	l.scanner.Split(bufio.ScanWords)
	return l
}

const (
	// keywords
	leftComment  = "/*"
	rightComment = "*/"
	typedef      = "typedef"
	structure    = "struct"

	// preprocessor directives; ugly naming convention
	pd_if     = "#if"
	pd_define = "#define"
	pd_else   = "#else"
	pd_endif  = "#endif"
)

const whitespace string = " "

// lexText ...
func lexText(l *lexer) stateFn {
	var text []string
	for l.scan() {
		word := l.text()
		if fn, ok := stateMap[word]; ok {
			l.emit(tknText, text)
			return fn
		}
		text = append(text, word)
	}
	if l.err() != nil {
		return l.errorf("Error lexing text, recieved: %s", l.err())
	}
	l.emit(tknEOF, tokenMap[tknEOF])
	return nil
}

func lexLeftComment(l *lexer) stateFn {
	l.emit(tknLeftComment, l.text())
	return lexText
}

func lexRightComment(l *lexer) stateFn {
	return l.errorf("Error lexing right comment:")
}

func lexTypedef(l *lexer) stateFn {
	return l.errorf("Error lexing typedef:")
}

func lexStructure(l *lexer) stateFn {
	return l.errorf("Error lexing struct:")
}

func lexDirective(l *lexer) stateFn {
	return l.errorf("Error lexing directive:")
}

// isSpace checks for a string containing whitespace, tab space, carriage return or newline
func isSpace(input string) bool {
	switch input {
	case " ", "\t", "\r", "\n":
		return true
	}
	return false
}
