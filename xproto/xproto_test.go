package xproto

import "testing"

func TestXproto(t *testing.T) {
	cp, err := NewXConnectionClientPrefix("MIT", "asdf;lkjalkjas;ljkfd;jk38909109sd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cp)
}
