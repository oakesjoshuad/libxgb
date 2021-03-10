package xau

import (
	"os"
	"testing"
)

func TestXauth(t *testing.T) {
	filename, err := xauFileName()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(filename)

	xafd, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer xafd.Close()

	xai, err := xauReadAuth(xafd)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(xai)
	xafd.Close()
}
