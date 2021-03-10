package libxgb

import (
	"testing"
)

func TestClientSetup(t *testing.T) {
	if xauth, err := xau.Xauth("", ":0"); err != nil {
		t.Fatal(err)
	} else {
		println(xauth)
	}
}
