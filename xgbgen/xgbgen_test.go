package xgbgen

import (
	"os"
	"testing"
)

const filename = "./testdata/Xproto.h"

func TestLexer(t *testing.T) {
	filehandle, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	lxr := lex(filename, filehandle)

	t.Run("test scan", func(t *testing.T) {
		for tok := lxr.parse(); tok.tokenType != tokenEOF; tok = lxr.parse() {
			t.Logf("%s", tok.String())
		}

	})
}
