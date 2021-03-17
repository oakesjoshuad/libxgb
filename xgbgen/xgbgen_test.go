package xgbgen

import (
	"strings"
	"testing"
)

var testinput = `
/* Definitions for the X window system used by server and c bindings */
`

func TestLexer(t *testing.T) {

	filehandle := strings.NewReader(testinput)
	lxr := lex("testinput", filehandle)

	t.Run("test scan", func(t *testing.T) {
		for tkn := lxr.nextToken(); tkn.Type != tknEOF; tkn = lxr.nextToken() {
			if tkn.Type == tknError {
				t.Log(tkn.String())
				break
			}
			t.Log(tkn.String())
		}
	})
}
