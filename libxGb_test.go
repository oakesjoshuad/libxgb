package libxgb

import (
	"reflect"
	"strconv"
	"testing"
)

func TestConnection(t *testing.T) {
	cases := []struct {
		Description string
		Input       string
		Output      *Display
		Err         error
	}{
		{"No Input", "", &Display{"localhost", "unix", 0, 0}, nil},
		{"With Input", "void/unix:0.10", &Display{"void", "unix", 0, 10}, nil},
		{"With only hostname as input", "localhost", &Display{"localhost", "unix", 0, 0}, nil},
		{"With bad hostname", "carbon", &Display{"localhost", "unix", 0, 0}, strconv.ErrSyntax},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			dp, err := parseDisplay(test.Input)
			if reflect.DeepEqual(dp, test.Output) {
				t.Log(dp)
			} else {
				t.Logf("Expected: %s\nRecieved: %s", test.Output, dp)
				t.Error(err)
			}
		})
	}
}
