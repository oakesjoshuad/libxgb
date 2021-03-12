package libxgb

import (
	"context"
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
		{"No Input", "", &Display{"localhost", "unix", "0", "", connection{}}, nil},
		{"With Input", "void/unix:0.10", &Display{"void", "unix", "0", "10", connection{}}, nil},
		{"With only hostname as input", "localhost", &Display{"localhost", "unix", "0", "", connection{}}, nil},
		{"With bad hostname", "carbon", &Display{"localhost", "unix", "0", "", connection{}}, nil},
		{"With only colon", ":", &Display{"localhost", "unix", "0", "", connection{}}, nil},
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
		ctx := context.Background()
		defer ctx.Done()
		dp, err := NewDisplay("")
		if err != nil {
			t.Error(err)
		}
		if err := dp.Open(ctx); err != nil {
			t.Error(err)
		} else {
			t.Log(dp.connection.RemoteAddr())
		}
	})
}
