package libxGb

import "testing"

func TestClientSetup(t *testing.T) {
	if xauth, err := xauthinfo("", ":0"); err != nil {
		t.Fatal(err)
	} else {
		println(xauth)
	}
}
