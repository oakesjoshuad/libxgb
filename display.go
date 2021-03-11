package libxgb

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// Display ...
type Display struct {
	Host, Protocol, Number, Screen string
}

const unixBase = "/tmp/.X11-unix/X"

var (
	// ErrDisplayName ...
	ErrDisplayName = errors.New("Unable to parse display")
)

func (dp *Display) String() string {
	return fmt.Sprintf("%s/%s:%s.%s", dp.Host, dp.Protocol, dp.Number, dp.Screen)
}

// Open ...
func (dp *Display) Open() (ncp net.Conn, err error) {
	if dp.Protocol != "unix" {
		return dp.openTCP()
	}
	ncp = new(net.UnixConn)
	ncp, err = dp.openUnix()
	return
}

func (dp *Display) openTCP() (*net.TCPConn, error) {
	return nil, errors.New("openTCP function not implemented")
}

func (dp *Display) openUnix() (*net.UnixConn, error) {
	if dp.Protocol != "unix" {
		return nil, errors.New("Incorrect Protocol")
	}
	address := strings.Join([]string{unixBase, dp.Number}, "")
	if dp.Screen != "" {
		address = strings.Join([]string{address, dp.Screen}, ".")
	}
	uap, err := net.ResolveUnixAddr(dp.Protocol, address)
	if err != nil {
		return nil, err
	}
	return net.DialUnix(dp.Protocol, nil, uap)
}

// NewDisplay will parse a given display string; if no string is given, it will check environment variables.
func NewDisplay(hostname string) (dp *Display, err error) {

	// checking if hostname is empty, if so, check environment variable for DISPLAY, else fail
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
