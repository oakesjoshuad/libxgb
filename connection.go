package libxgb

import (
	"errors"
	"os"
)

// Display ...
type Display struct {
}

// parseDisplay
func parseDisplay(hostname string) (dp *Display, err error) {

	// checking if hostname is empty, if so, check environment variable for DISPLAY, else fail
	if hostname == "" {
		if hostname = os.Getenv("DISPLAY"); hostname == "" {
			err = errors.New("Unable to parse display name")
			return
		}
	}

	// check for launchd, <path to socket>[.<screen>]
	if dp, err = parseDisplaySocket(hostname); err == nil {
		return
	}
	return
}

// parseDisplaySocket
func parseDisplaySocket(hostname string) (dp *Display, err error) {
	err = errors.New("function not implemented")
	return
}
