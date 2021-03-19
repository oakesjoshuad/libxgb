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
	tknDirective
	tknIdentifier
	tknConstant
	tknLiteral
	tknSymbol
	tknText
	tknLeftComment
	tknCommentText
	tknRightComment
	tknInlineComment
)

var tokenMap = map[tokenType]string{
	tknError:         "Error",
	tknEOF:           "EOF",
	tknKeyword:       "Keyword",
	tknDirective:     "Directive",
	tknIdentifier:    "Identifier",
	tknConstant:      "Constant",
	tknLiteral:       "Literal",
	tknSymbol:        "Symbol",
	tknText:          "Text",
	tknLeftComment:   "Begin Comment",
	tknCommentText:   "Text",
	tknRightComment:  "End Comment",
	tknInlineComment: "Inline Comment",
}

// String is a stringer method returning a string representation of a token
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
	linepos int
	tokens  chan token
}

// scan ...
func (l *lexer) scan() bool {
	return l.scanner.Scan()
}

func (l *lexer) checkScanErr() error {
	return l.scanner.Err()
}

// emit checks that a token is not empty and puts it on the token channel
func (l *lexer) emit(t tokenType, s string) {
	if !isEmpty(s) {
		l.tokens <- token{t, s}
	}
}

// errorf ...
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tknError, fmt.Sprintf(format, args...)}
	return nil
}

// nextToken returns a token from the token channel if it exists or moves to the next state
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
		linepos: 0,
		tokens:  make(chan token, 2),
	}
	l.scanner.Split(bufio.ScanLines)

	key = map[string]stateFn{
		cStruct:  lexCStruct,
		cTypedef: lexCTypedef,
		cDefine:  lexCDirective,
	}

	return l
}

const (
	// symbols
	leftBrace     = "{"
	rightBrace    = "}"
	leftBracket   = "["
	rightBracket  = "]"
	leftComment   = "/*"
	rightComment  = "*/"
	inlineComment = "//"
	semicolon     = ";"

	// keywords our lexer should be interested in
	cStruct  = "struct"
	cTypedef = "typedef"
	cDefine  = "#define"
)

var key map[string]stateFn

// lexText ...
func lexText(l *lexer) stateFn {
	var txt string
	for txt := ""; l.scan(); txt += l.scanner.Text() {
		// comment bodies are non standardized, checking of our line is prefixed with
		// an opening comment identifier, '/*' or '//'
		if strings.HasPrefix(l.scanner.Text(), leftComment) {
			l.emit(tknText, txt)
			return lexLeftComment
		} else if strings.HasPrefix(l.scanner.Text(), inlineComment) {
			l.emit(tknText, txt)
			return lexInlineComment
		}
		// otherwise, I'm looking for a keyword
		if words := strings.Fields(l.scanner.Text()); len(words) > 0 {
			word := words[l.linepos]
			if fn, ok := key[word]; ok {
				l.emit(tknText, txt)
				return fn
			}
		}
	}

	// check for a scan error
	if err := l.checkScanErr(); err != nil {
		l.errorf("encountered an error while scanning text: %s", err)
	}

	// we have reached the end of the line
	l.emit(tknEOF, txt)
	return nil
}

// lexLeftCommenet  ...
func lexLeftComment(l *lexer) stateFn {
	if strings.HasSuffix(l.scanner.Text(), rightComment) {
		l.emit(tknLeftComment, leftComment)
		return lexInlineComment
	}
	l.emit(tknLeftComment, leftComment)
	return lexInsideComment
}

// lexInsideComment ...
func lexInsideComment(l *lexer) stateFn {
	var txt string
	for txt = strings.TrimPrefix(l.scanner.Text(), leftComment); l.scan(); txt += l.scanner.Text() {
		if strings.HasSuffix(l.scanner.Text(), rightComment) {
			txt += strings.TrimSuffix(l.scanner.Text(), rightComment)
			l.emit(tknCommentText, txt)
			return lexRightComment
		}
	}
	if err := l.checkScanErr(); err != nil {
		return l.errorf("encountered an error scanning a comment: %s", err)
	}
	return l.errorf("encounted EOF while lexing inside comment: %s", txt)
}

// lexRightComment ...
func lexRightComment(l *lexer) stateFn {
	l.emit(tknRightComment, rightComment)
	return lexText
}

// lexInlineComment ...
func lexInlineComment(l *lexer) stateFn {
	txt := strings.Trim(l.scanner.Text(), leftComment+rightComment+inlineComment)
	l.emit(tknInlineComment, txt)
	return lexText
}

// lexCDirective ...
func lexCDirective(l *lexer) stateFn {
	return l.errorf("lexCDirective Not Implemented")
}

// lexCTypedef ...
func lexCTypedef(l *lexer) stateFn {
	return l.errorf("Error parsing Ctypedef")
}

// lexCStruct ...
func lexCStruct(l *lexer) stateFn {
	return l.errorf("Error lexing struct:")
}

const (
	// directives
	define = "#define"
)

func lexConstant(l *lexer) stateFn {
	return l.errorf("lexConstant Not Implemented")
}
func lexIdentifier(l *lexer) stateFn {
	return l.errorf("lexIdentifier Not Implemented")
}

func lexLiteral(l *lexer) stateFn {
	return l.errorf("lexLiteral Not Implemented")
}

const (
	whitespace     = ' '
	newline        = '\n'
	tab            = '\t'
	carriagereturn = '\r'
)

func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isWord(s string) bool {
	for _, r := range s {
		if !isAlphaNumeric(r) {
			return false
		}
	}
	return true
}

func isEmpty(s string) bool {
	for _, r := range s {
		if !isSpace(r) {
			return false
		}
	}
	return true
}

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
