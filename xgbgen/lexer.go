// Package xgbgen generates protocols for X directly from X header files
package xgbgen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"text/scanner"
)

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenError
	tokenLeftComment
	tokenComment
	tokenRightComment
	tokenText
)

var tokenMap = map[tokenType]string{
	tokenEOF:          "end of file",
	tokenError:        "error",
	tokenLeftComment:  "/*",
	tokenRightComment: "*/",
	tokenText:         "text",
}

type token struct {
	tokenType
	tokenLiteral string
}

func (t *token) String() string {
	return fmt.Sprintf("{ tokenType: %s, tokenLiteral: %s", tokenMap[t.tokenType], t.tokenLiteral)
}

type lexer struct {
	filename string
	scanner  *bufio.Scanner
	state    stateFn
	tokens   chan token
}

type stateFn func(*lexer) stateFn

func lex(filename string, input io.Reader) *lexer {
	l := &lexer{
		filename: filename,
		scanner:  bufio.NewScanner(input),
		state:    lexText,
		tokens:   make(chan token, 2),
	}
	return l
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) scan(split bufio.SplitFunc) rune {
	l.scanner.Split(split)
	return
}

func (l *lexer) peek() rune {
	return l.scanner.Peek()
}

func (l *lexer) parse() token {
	for {
		select {
		case tok := <-l.tokens:
			return tok
		default:
			l.state = l.state(l)
		}
	}
}

const (
	leftComment  = "/*"
	rightComment = "*/"
)

func lexText(l *lexer) stateFn {
	var buff bytes.Buffer
	for tok := l.scan(); tok != scanner.EOF; tok = l.scan() {
		switch tok {
		case '/':
			if l.peek() == '*' {
				if buff.Len() > 0 {
					l.tokens <- token{tokenText, buff.String()}
				}
				return lexLeftComment
			}
		default:
			if _, err := buff.WriteString(l.scanner.TokenText()); err != nil {
				return l.errorf("line: %d; %s", l.scanner.Line, err)
			}
		}
	}
	if buff.Len() > 0 {
		l.tokens <- token{tokenType: tokenText, tokenLiteral: buff.String()}
	}
	l.tokens <- token{tokenType: tokenEOF}
	return nil
}

func lexLeftComment(l *lexer) stateFn {
	if tok := l.scan(); tok == '*' {
		l.tokens <- token{tokenLeftComment, leftComment}
		return lexComment
	}
	return l.errorf("parsing error at line: %d", l.scanner.Line)
}

func lexComment(l *lexer) stateFn {
	var buff bytes.Buffer
	for tok := l.scan(); tok != scanner.EOF; tok = l.scan() {
		switch tok {
		case '*':
			if l.peek() == '/' {
				if buff.Len() > 0 {
					l.tokens <- token{tokenComment, buff.String()}
				}
				return lexRightComment
			}
		default:
			if _, err := buff.WriteString(l.scanner.TokenText()); err != nil {
				return l.errorf("line: %d; %s", l.scanner.Line, err)
			}
		}
	}
	return l.errorf("non terminated comment")
}

func lexRightComment(l *lexer) stateFn {
	if tok := l.scan(); tok == '/' {
		l.tokens <- token{tokenRightComment, rightComment}
		return lexText
	}
	return l.errorf("malformed comment, line: %d", l.scanner.Line)
}
