package libxgb_test

import (
	"reflect"
	"testing"

	"github.com/oakesjoshuad/libxgb"
)

func TestExternal(t *testing.T) {
	cases := []struct {
		Description string
		Input       string
		Output      *libxgb.Display
		Err         error
	}{
		{"No Input", "", &libxgb.Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With Input", "void/unix:0.10", &libxgb.Display{Host: "void", Protocol: "unix", Number: "0", Screen: "10"}, nil},
		{"With only hostname as input", "localhost", &libxgb.Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With bad hostname", "carbon", &libxgb.Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With only colon", ":", &libxgb.Display{Host: "localhost", Protocol: "unix", Number: "0", Screen: ""}, nil},
	}

	t.Run("Testing NewDisplay", func(t *testing.T) {
		for _, test := range cases {
			t.Run(test.Description, func(t *testing.T) {
				dp, err := libxgb.NewDisplay(test.Input)
				if !reflect.DeepEqual(dp, test.Output) || err != nil {
					t.Error(err)
				}
				t.Log("Expected: ", test.Output)
				t.Log("Recieved: ", dp)
			})
		}
	})

	t.Run("Testing open/close", func(t *testing.T) {
		dp, err := libxgb.NewDisplay("")
		if err != nil {
			t.Error(err)
		}
		if err := dp.Open(); err != nil {
			t.Error(err)
		} else {
			t.Log(dp.String())
		}
		if err := dp.Close(); err != nil {
			t.Error(err)
		}
	})
}
