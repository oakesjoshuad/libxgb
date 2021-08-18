// Package xau helps parse the .Xauthority file and returns an Xauth object for use in connecting to an Xserver
package xau

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	// FamilyLocal ...
	FamilyLocal = 256
	// FamilyWild ...
	FamilyWild = 65535
	// FamilyNetName ...
	FamilyNetName = 254
	// FamilyKrb5Principal ...
	FamilyKrb5Principal = 253
	// FamilyLocalHost ...
	FamilyLocalHost = 252

	slashDotXauthority = "/.Xauthority"
)

var (
	// ErrXauthorityLocation ...
	ErrXauthorityLocation = errors.New("Unable to locate Xauthority file")
)

// Xauth ...
type Xauth struct {
	Family  uint16
	Address string
	Number  string
	Name    string
	Data    string
}

func (xa *Xauth) String() string {
	return fmt.Sprintf("family: %d, address: %s, number: %s, name: %s, data: %s",
		xa.Family,
		xa.Address,
		xa.Number,
		xa.Name,
		xa.Data,
	)
}

// XauFileName returns the the location of the .Xauthority file
func xauFileName() (string, error) {
	if filename := os.Getenv("XAUTHORITY"); filename != "" {
		return filename, nil
	}
	if home := os.Getenv("HOME"); home != "" {
		return home + slashDotXauthority, nil
	}
	return "", ErrXauthorityLocation
}

// XauReadAuth returns Xauth
func xauReadAuth(rdr io.Reader) (xa *Xauth, err error) {
	xa = new(Xauth)
	if err != nil {
		return
	}
	if xa.Family, err = readShort(rdr); err != nil {
		return
	}
	if xa.Address, err = readString(rdr); err != nil {
		return
	}
	if xa.Number, err = readString(rdr); err != nil {
		return
	}
	if xa.Name, err = readString(rdr); err != nil {
		return
	}
	if xa.Data, err = readString(rdr); err != nil {
		return
	}
	return
}

// GetAuthByAddr ...
func GetAuthByAddr(family uint16, address, number, name string) (xa *Xauth, err error) {
	xaufilename, err := xauFileName()
	if err != nil {
		return
	}
	xaufd, err := os.Open(xaufilename)
	defer xaufd.Close()
	if err != nil {
		return
	}
	for xa, err = xauReadAuth(xaufd); err == nil; xa, err = xauReadAuth(xaufd) {
		if (family == FamilyWild || xa.Family == FamilyWild || (xa.Family == family && xa.Address == address)) && (xa.Number == number) && (xa.Name == name) {
			return
		}
	}
	return
}

// GetBestAuthByAddr ...
func GetBestAuthByAddr(family uint16, address, number string, authTypes []string) (best *Xauth, err error) {
	xaufilename, err := xauFileName()
	if err != nil {
		return
	}
	xaufd, err := os.Open(xaufilename)
	if err != nil {
		return
	}
	xa := new(Xauth)
	for xa, err = xauReadAuth(xaufd); err == nil; xa, err = xauReadAuth(xaufd) {
		if (family == FamilyWild || xa.Family == FamilyWild || (xa.Family == family && xa.Address == address)) && (xa.Number == number) {
			for _, authType := range authTypes {
				if authType == xa.Name {
					best = xa
				}
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// readShort reads 2 bytes from the buffer and returns
func readShort(rdr io.Reader) (uint16, error) {
	var buf = make([]byte, 2)
	if err := binary.Read(rdr, binary.BigEndian, buf); err != nil {
		return 0, err
	}
	return uint16(buf[0])*256 + uint16(buf[1]), nil
}

// readString
func readString(rdr io.Reader) (string, error) {
	width, err := readShort(rdr)
	if err != nil {
		return "", err
	}
	var buf = make([]byte, width)
	if err := binary.Read(rdr, binary.BigEndian, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
