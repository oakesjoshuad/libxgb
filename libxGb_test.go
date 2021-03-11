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
		{"No Input", "", &Display{"localhost", "unix", "0", ""}, nil},
		{"With Input", "void/unix:0.10", &Display{"void", "unix", "0", "10"}, nil},
		{"With only hostname as input", "localhost", &Display{"localhost", "unix", "0", ""}, nil},
		{"With bad hostname", "carbon", &Display{"localhost", "unix", "0", ""}, nil},
		{"With only colon", ":", &Display{"localhost", "unix", "0", ""}, nil},
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

	t.Run("Testing open", func(t *testing.T) {
		dp, err := NewDisplay("")
		if err != nil {
			t.Error(err)
		}
		if cp, err := dp.Open(); err != nil {
			t.Error(err)
		} else {
			t.Log(cp.LocalAddr())
			t.Log(cp.RemoteAddr())

		}
	})
}
