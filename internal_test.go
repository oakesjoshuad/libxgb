package libxgb

import (
	"reflect"
	"testing"
)

func TestInternal(t *testing.T) {
	cases := []struct {
		Description string
		Input       string
		Output      *Display
		Err         error
	}{
		{"No Input", "", &Display{Host: "void", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With Input", "void/unix:0.10", &Display{Host: "void", Protocol: "unix", Number: "0", Screen: "10"}, nil},
		{"With only hostname as input", "localhost", &Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With bad hostname", "carbon", &Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With only colon", ":", &Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
	}

	t.Run("Testing NewDisplay", func(t *testing.T) {
		for _, test := range cases {
			t.Run(test.Description, func(t *testing.T) {
				dp, err := NewDisplay(test.Input)
				if !reflect.DeepEqual(dp, test.Output) || err != nil {
					t.Error(err)
				}
				t.Log("Expected: ", test.Output)
				t.Log("Recieved: ", dp)
			})
		}
	})

	t.Run("Testing open/close", func(t *testing.T) {
		dp, err := NewDisplay("")
		if err != nil {
			t.Error(err)
		}
		if err := dp.Open(); err != nil {
			t.Error(err)
		}
		if err := dp.Close(); err != nil {
			t.Error(err)
		}
	})
}
