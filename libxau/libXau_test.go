package libXau

import "testing"

func TestXauth(t *testing.T) {
	filename, err := XauFileName()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(filename)

	xai, err := XauReadAuth(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(xai)
}
