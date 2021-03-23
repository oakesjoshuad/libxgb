package libxgb

import (
	"reflect"
	"testing"

	"github.com/oakesjoshuad/libxgb/xau"
	"github.com/oakesjoshuad/libxgb/xproto"
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
		} else {
		}
		if xa, err := xau.GetAuthByAddr(xau.FamilyLocal, dp.Host, dp.Number, MIT); err != nil {
			t.Log(err)
		} else {
			cs := xproto.NewClientSetup(xa.AuthName, xa.AuthData)
			dp.Send(Message{len(cs), cs})
		}
		if msg := dp.CheckMessage(); msg.Length != 0 {
			t.Log(string(msg.Payload))
		}
		if err := dp.Close(); err != nil {
			t.Error(err)
		}
	})
}
