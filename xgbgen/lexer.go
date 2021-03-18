// Package xgbgen generates protocols for X directly from X header files
package xgbgen

import (
	"bufio"
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
	text    string
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

// lex returns a new lexer, given an io.Reader
func lex(name string, rdr io.Reader) *lexer {
	l := &lexer{
		name:    name,
		scanner: bufio.NewScanner(rdr),
		state:   lexText,
		tokens:  make(chan token, 2),
	}
	l.scanner.Split(bufio.ScanLines)
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
	directive    = "#"
)

// lexText ...
func lexText(l *lexer) stateFn {
	for txt := ""; l.scan(); l.text = txt {
		if strings.HasPrefix(l.scanner.Text(), leftComment) {
			l.text = strings.TrimPrefix(l.scanner.Text(), leftComment)
			l.emit(tknText, txt)
			return lexLeftComment
		}
		if strings.HasPrefix(l.scanner.Text(), directive) {
			l.text = strings.TrimPrefix(l.scanner.Text(), directive)
			l.emit(tknText, txt)
			return lexDirective
		}
		txt += l.scanner.Text()
	}
	// check for a scan error
	if err := l.checkScanErr(); err != nil {
		l.errorf("encountered an error while scanning text: %s", err)
	}

	// we have reached the end of the line
	l.emit(tknEOF, l.text)
	return nil
}

func lexLeftComment(l *lexer) stateFn {
	l.emit(tknLeftComment, leftComment)
	return lexInsideComment
}

func lexInsideComment(l *lexer) stateFn {
	for txt := l.text; !strings.HasSuffix(txt, rightComment); l.text = txt {
		if l.scan() {
			txt = fmt.Sprintf("%s\n%s", txt, l.scanner.Text())
		} else if l.scanner.Err() != nil {
			l.errorf("encountered an error while scanning inside comment: %s", l.scanner.Err())
		} else {
			l.errorf("encountered end of file before comment termination")
		}

	}
	txt := strings.TrimSuffix(l.text, rightComment)
	l.emit(tknCommentText, txt)
	return lexRightComment
}

func lexRightComment(l *lexer) stateFn {
	l.emit(tknRightComment, rightComment)
	return lexText
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
