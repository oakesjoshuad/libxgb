// Package xgbgen generates protocols for X directly from X header files
package xgbgen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
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
	tknCommentText
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
	tknLeftComment:  "Begin Comment",
	tknCommentText:  "Text",
	tknRightComment: "End Comment",
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
	words   []string
	tokens  chan token
}

// scan ...
func (l *lexer) scan() bool {
	return l.scanner.Scan()
}

func (l *lexer) checkScanErr() error {
	return l.scanner.Err()
}

func (l *lexer) emit(t tokenType, s string) {
	l.tokens <- token{t, s}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tknError, fmt.Sprintf(format, args...)}
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
	l.scanner.Split(bufio.ScanLines)

	// initialize the state map
	stateMap = map[string]stateFn{
		leftComment:  lexLeftComment,
		rightComment: lexRightComment,
		typedef:      lexTypedef,
		structure:    lexStructure,
		pd_if:        lexDirective,
		pd_define:    lexConstant,
		pd_else:      lexDirective,
		pd_endif:     lexDirective,
	}

	return l
}

const (
	// symbols
	leftBrace    = "{"
	rightBrace   = "}"
	leftBracket  = "["
	rightBracket = "]"
	leftComment  = "/*"
	rightComment = "*/"
	// keywords
	typedef   = "typedef"
	structure = "struct"

	// preprocessor directives; ugly naming convention
	pd_if     = "#if"
	pd_define = "#define"
	pd_else   = "#else"
	pd_endif  = "#endif"
)

// lexText ...
func lexText(l *lexer) stateFn {
	buf := new(bytes.Buffer)
	for l.scan() {
		// check for a scan error
		if err := l.checkScanErr(); err != nil {
			l.errorf("encountered an error while scanning text: %s", err)
		}
		// split the scanned text into a slice of words
		words := strings.Split(l.scanner.Text(), " ")
		// if the first word is a keyword, emit a non-empty buffer as text and return the appropriate stateFn
		if fn, ok := stateMap[words[0]]; ok {
			if buf.Len() > 0 {
				l.emit(tknText, buf.String())
			}
			// retain split input for the following states, sans keyword
			l.words = words[1:]
			return fn
		}
		// write the original text into the local buffer preserving formatting
		buf.WriteString(l.scanner.Text())
	}
	// we have reached the end of the line
	l.emit(tknEOF, tokenMap[tknEOF])
	return nil
}

func lexLeftComment(l *lexer) stateFn {
	l.emit(tknLeftComment, leftComment)
	return lexComment
}

func lexComment(l *lexer) stateFn {
	return l.errorf("Not Implemented")
}

func lexRightComment(l *lexer) stateFn {
	return l.errorf("Not Implemented")
}

func lexTypedef(l *lexer) stateFn {
	return l.errorf("Error lexing typedef:")
}

func lexStructure(l *lexer) stateFn {
	return l.errorf("Error lexing struct:")
}

func lexDirective(l *lexer) stateFn {
	return l.errorf("Error lexing directive, function not implemented")
}

func lexConstant(l *lexer) stateFn {
	return l.errorf("Not Implemented")
}

func lexIdentifier(l *lexer) stateFn {
	return l.errorf("Not Implemented")
}

func lexValue(l *lexer) stateFn {
	return l.errorf("Not Implemented")
}

const (
	whitespace     = ' '
	newline        = '\n'
	tab            = '\t'
	carriagereturn = '\r'
)

func isSpace(r rune) bool {
	switch r {
	case whitespace, newline, tab, carriagereturn:
		return true
	}
	return false
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
