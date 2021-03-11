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
		Err         error
	}{
		{"No Input", "", &Display{"localhost", "unix", 0, 0}, nil},
		{"With Input", "void/unix:0.10", &Display{"void", "unix", 0, 10}, nil},
		{"With only hostname as input", "localhost", &Display{"localhost", "unix", 0, 0}, nil},
		{"With bad hostname", "carbon", &Display{"localhost", "unix", 0, 0}, nil},
		{"With only colon", ":", &Display{"localhost", "unix", 0, 0}, nil},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			dp, err := parseDisplay(test.Input)
			if !reflect.DeepEqual(dp, test.Output) || err != nil {
				t.Error(err)
			}
			t.Log("Expected: ", test.Output)
			t.Log("Recieved: ", dp)
		})
	}
}
