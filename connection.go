package libxgb

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Display ...
type Display struct {
	Host     string
	Protocol string
	Number   uint64
	Screen   uint64
}

const (
	unixBase = "/tmp/.X11-unix/X"
)

var (
	// ErrDisplayName ...
	ErrDisplayName = errors.New("Unable to parse display")
)

func (dp *Display) String() string {
	return fmt.Sprintf("Host: %s, Protocol: %s, Number: %d, Screen: %d", dp.Host, dp.Protocol, dp.Number, dp.Screen)
}

// ParseDisplay will parse a given display string; if no string is given, it will check environment variables.
func parseDisplay(hostname string) (dp *Display, err error) {

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
		dp.Screen = 0
		if dp.Number, err = strconv.ParseUint(hostname[colon:], 10, 32); err != nil {
			return
		}
	} else {
		if dp.Number, err = strconv.ParseUint(hostname[colon:dot], 10, 32); err != nil {
			return
		}
		dot++
		if dp.Screen, err = strconv.ParseUint(hostname[dot:], 10, 32); err != nil {
			return
		}
	}
	return
}
