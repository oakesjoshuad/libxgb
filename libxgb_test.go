package libxgb

import (
	"os"
	"reflect"
	"testing"
)

func TestInternal(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Description string
		Input       string
		Output      *Display
		Err         error
	}{
		{"No Input", "", &Display{Host: hostname, Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With Input", "void/unix:0.10", &Display{Host: "void", Protocol: "unix", Number: "0", Screen: "10"}, nil},
		{"With only hostname as input", hostname, &Display{Host: hostname, Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With bad hostname", "carbon", &Display{Host: hostname, Protocol: "unix", Number: "0", Screen: ""}, nil},
		{"With only colon", ":", &Display{Host: hostname, Protocol: "unix", Number: "0", Screen: ""}, nil},
	}

	t.Run("Testing NewDisplay", func(t *testing.T) {
		for _, test := range cases {
			t.Run(test.Description, func(t *testing.T) {
				dp, err := NewDisplay(test.Input)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(dp, test.Output) {
					t.Log("Expected: ", test.Output)
					t.Log("Recieved: ", dp)
				}
			})
		}
	})

	t.Run("Testing open/close", func(t *testing.T) {
		dp, err := NewDisplay("")
		if err != nil {
			t.Error(err)
		}
		if _, err := dp.Open(); err != nil {
			t.Fatal(err)
		}
		err = dp.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
}
