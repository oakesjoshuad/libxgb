package libxgb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/oakesjoshuad/libxgb/xau"
	"github.com/oakesjoshuad/libxgb/xproto"
)

// Display ...
type Display struct {
	Host, Protocol, Number, Screen string
	// encapsulating the connection to expose only needed functionality.
	connection connection
	ctx        context.Context
	// channels to buffer communication
	send  chan Message
	recv  chan Message
	errs  chan error
	close chan bool
}

type connection struct {
	net.Conn
}

// Message is the primary method of interacting with the Xserver through the display connection. It consists of the message payload in byte string form and the message length.
type Message struct {
	Length  int
	Payload []byte
}

// unixBase contains the file path of the unix "socket"
const unixBase = "/tmp/.X11-unix/X"

// Authorization types
const (
	// MIT authorization
	MIT = "MIT-MAGIC-COOKIE-1"
)

var authNames = []string{MIT}

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
	}
	return fmt.Sprintf("%s%s", unixBase, dp.Number)
}

// Open ...
func (dp *Display) Open() ([]byte, error) {
	ctx := context.TODO()
	defer ctx.Done()

	if dp.Protocol != "unix" {
		if err := dp.openTCP(ctx); err != nil {
			return nil, err
		}
	} else {
		if err := dp.openUnix(ctx); err != nil {
			return nil, err
		}
	}

	xauth, err := xau.GetBestAuthByAddr(xau.FamilyLocal, dp.Host, dp.Number, []string{MIT})
	if err != nil {
		return nil, fmt.Errorf("error getting xauth: %w", err)
	}
	cpfx, err := xproto.NewXConnectionClientPrefix(xauth.Name, xauth.Data)
	if err != nil {
		return nil, fmt.Errorf("error generating client prefix: %w", err)
	}
	if n, err := dp.connection.Write(cpfx); err != nil {
		return nil, fmt.Errorf("error writing client prefix to connection: %w", err)
	} else if n < len(cpfx) {
		return nil, fmt.Errorf("error writing cpfx to connection, %d of %d bytes written", n, len(cpfx))
	}

	spfx, err := xproto.NewXConnectionSetupPrefix(dp.connection)
	if err != nil {
		return nil, err
	}

	// if we failed, parse the reason for failure.
	if spfx.Status == xproto.XConnectionFailed {
		buf := new(bytes.Buffer)
		if ln, err := buf.ReadFrom(dp.connection); err != nil {
			return nil, fmt.Errorf("error reading reason for refused connection: %w", err)
		} else if ln < int64(spfx.ReasonLen) {
			return nil, fmt.Errorf("error parsing reason for refused connection")
		}
		return nil, fmt.Errorf("Xserver refused connection: %s", buf.String())
	} else if spfx.Status == xproto.XConnectionAuth {
		buf := new(bytes.Buffer)
		if ln, err := buf.ReadFrom(dp.connection); err != nil {
			return nil, fmt.Errorf("error reading extra authentication information reason: %w", err)
		} else if ln < int64(spfx.ReasonLen) {
			return nil, fmt.Errorf("error parsing reason for extra authentication, %d bytes of %d bytes read", ln, spfx.ReasonLen)
		}
		return nil, fmt.Errorf("Xserver requires extra authentication: %s", buf.String())
	}

	// successfully initiated connection with Xserver, parse the returned setup information

	// recieve channel
	//dp.recv = make(chan Message)
	// send channel
	//dp.send = make(chan Message)
	// error channel
	//dp.errs = make(chan error)
	// sentinal channel to signal close
	//dp.close = make(chan bool)

	//go dp.rxtx()
	//go dp.err()

	return nil, nil
}

// Send puts a Message on the send channel
func (dp *Display) Send(msg Message) {
	dp.send <- msg
}

// Close sends close signal to all channels and closes the connection
func (dp *Display) Close() error {
	//dp.close <- true
	//close(dp.close)
	return dp.connection.Close()
}

// tx transmit Message to Xserver
func (dp *Display) rxtx() {
	select {
	case msg := <-dp.send:
		if n, err := dp.connection.Write(msg.Payload); err != nil {
			dp.errs <- err
		} else if n != int(msg.Length) {
			dp.errs <- errors.New("Display.tx did not transmit full Message")
		}
	case <-dp.close:
		close(dp.send)
		close(dp.recv)
		close(dp.errs)
	default:
		var buff bytes.Buffer
		if n, err := buff.ReadFrom(dp.connection); err != nil {
			dp.errs <- fmt.Errorf("Error reading response from Xserver: %w", err)
		} else if n > 0 {
			dp.recv <- Message{Length: int(n), Payload: buff.Bytes()}
		}
	}
}

func (dp *Display) err() {
	select {
	case <-dp.close:
		close(dp.errs)
	case err := <-dp.errs:
		fmt.Println(err)
	}
}

// openTCP ...
func (dp *Display) openTCP(pctx context.Context) (err error) {
	return errors.New("openTCP function not implemented")
}

// openUnix ...
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
// a validated input string or the environment variable DISPLAY
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
		if dp.Host, err = os.Hostname(); err != nil {
			return
		}
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
