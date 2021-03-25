package xproto

import "testing"

func TestXproto(t *testing.T) {
	_, err := NewClientPrefix("MIT", "asdf;lkjalkjas;ljkfd;jk38909109sd")
	if err != nil {
		t.Fatal(err)
	}

}
