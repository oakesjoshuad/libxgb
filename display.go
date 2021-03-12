package libxgb

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// Display ...
type Display struct {
	Host, Protocol, Number, Screen string

	connection
}

type connection struct {
	net.Conn
	ctx   context.Context
	send  chan []byte
	recv  chan []byte
	errs  chan error
	close chan bool
}

const unixBase = "/tmp/.X11-unix/X"

var (
	// ErrDisplayName ...
	ErrDisplayName = errors.New("Unable to parse display")
)

func (dp *Display) String() string {
	return fmt.Sprintf("%s/%s:%s.%s", dp.Host, dp.Protocol, dp.Number, dp.Screen)
}

func (dp *Display) unixAddress() string {
	if dp.Screen != "" {
		return fmt.Sprintf("%s%s.%s", unixBase, dp.Number, dp.Screen)
	} else {
		return fmt.Sprintf("%s%s", unixBase, dp.Number)
	}
}

// Open ...
func (dp *Display) Open(pctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()
	if dp.Protocol != "unix" {
		err = dp.openTCP(ctx)
	} else {
		err = dp.openUnix(ctx)
	}
	// recieve channel
	dp.connection.recv = make(chan []byte)
	// send channel
	dp.connection.send = make(chan []byte)
	// sentinal channel to signal close
	dp.connection.close = make(chan bool)

	return
}

// TODO: comment on tx and rx functions
// TODO: write tests for this open and close
// TODO: should I binary write, or just use the given write methods
// TODO: conn type based auth, rx and tx

// Close ...
func (dp *Display) Close() error {
	dp.close <- true
	return dp.connection.Close()
}

func (dp *Display) tx() {
	select {
	case msg := <-dp.connection.send:
		if n, err := dp.connection.Write(msg); err != nil {
			dp.connection.errs <- err
		} else if n != len(msg) {
			dp.connection.errs <- errors.New("Display.tx did not transmit full message")
		}
	case <-dp.connection.close:
		close(dp.connection.send)
	}
}

func (dp *Display) openTCP(pctx context.Context) (err error) {
	return errors.New("openTCP function not implemented")
}

func (dp *Display) openUnix(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()
	if dp.Protocol != "unix" {
		return fmt.Errorf("Incorrect protocol: %s", dp.Protocol)
	}
	var dlr net.Dialer
	dlr.LocalAddr = nil
	raddr, err := net.ResolveUnixAddr(dp.Protocol, dp.unixAddress())
	if err != nil {
		return err
	}
	if dp.connection.Conn, err = dlr.DialContext(ctx, dp.Protocol, raddr.String()); err != nil {
		return err
	}
	return nil
}

// NewDisplay returns a populated Display structure with sane defaults; constructed from
// a, validated, passed in display string or the environment variable DISPLAY
func NewDisplay(hostname string) (dp *Display, err error) {

	if !strings.Contains(hostname, ":") {
		if hostname = os.Getenv("DISPLAY"); hostname == "" {
			err = ErrDisplayName
			return
		}
	}

	dp = new(Display)
	slash := strings.LastIndex(hostname, "/")
	colon := strings.LastIndex(hostname, ":")
	dot := strings.LastIndex(hostname, ".")

	if slash < 0 {
		dp.Host = "localhost"
		dp.Protocol = "unix"
	} else {
		dp.Host = hostname[:slash]
		slash++
		dp.Protocol = hostname[slash:colon]
	}
	colon++
	if dot < 0 {
		dp.Screen = ""
		if dp.Number = hostname[colon:]; dp.Number == "" {
			dp.Number = "0"
		}
	} else {
		dp.Number = hostname[colon:dot]
		dot++
		dp.Screen = hostname[dot:]
	}
	return
}
