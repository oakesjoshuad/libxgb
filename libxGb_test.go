package libxgb

import (
	"reflect"
	"testing"
)

func TestConnection(t *testing.T) {
	cases := []struct {
		Description string
		Input       string
		Output      *Display
	}{
		{"No Input", "", &Display{Host: "localhost", Protocol: "", Number: 0, Screen: 0}},
		{"With Input", "void/unix:0.10", &Display{Host: "void", Protocol: "unix", Number: 0, Screen: 10}},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			dp, err := parseDisplay(test.Input)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			if reflect.DeepEqual(dp, test.Output) {
				t.Log(dp)
			} else {
				t.Logf("Expected: %s\nRecieved: %s", test.Output, dp)
				t.Fail()
			}
		})
	}
}
