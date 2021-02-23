package libxGb

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const (
	FamilyLocalHost = 252
	FamilyWild      = 65535
	FamilyNetName   = 254
	FamilyLocal     = 256
	slashdotxauth   = "/.XAuthority"
)

var (
	Endianess = binary.BigEndian

	ErrAuthFileLocation = errors.New("Unable to parse XAuthority file location")
)

type xauth struct {
	family   uint16
	address  string
	number   string
	protocol string
	data     string
}

func readBytes(r io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if err := binary.Read(r, Endianess, buf); err != nil {
		return 0, err
	}
	output := uint16(buf[0])*256 + uint16(buf[1])
	return output, nil
}

func readString(r io.Reader) (string, error) {
	if l, err := readBytes(r); err != nil {
		return "", err
	} else {
		buf := make([]byte, l)
		if err := binary.Read(r, Endianess, buf); err != nil {
			return "", err
		}
		return string(buf), nil
	}
}

func readAuth(r io.Reader) (xa xauth, err error) {

	if xa.family, err = readBytes(r); err != nil {
		return
	}
	if xa.address, err = readString(r); err != nil {
		return
	}
	if xa.number, err = readString(r); err != nil {
		return
	}
	if xa.protocol, err = readString(r); err != nil {
		return
	}
	if xa.data, err = readString(r); err != nil {
		return
	}
	return
}

func xauthinfo(host, display string) (protocol, data string, err error) {

	if host == "" || host == "localhost" {
		if host, err = os.Hostname(); err != nil {
			return
		}
	}

	var filename string
	if filename = os.Getenv("XAUTHORITY"); filename == "" {
		if filename = os.Getenv("HOME"); filename == "" {
			err = ErrAuthFileLocation
			return
		}
		filename += slashdotxauth
	}

	authfile := new(os.File)
	if authfile, err = os.OpenFile(filename, os.O_RDONLY, 0644); err != nil {
		return
	}

	var xa xauth
	for xa, err = readAuth(authfile); err == nil; xa, err = readAuth(authfile) {
		family := xa.family == FamilyWild || (xa.family == FamilyLocal && xa.address == host)
		disp := xa.number == "" || xa.number == display
		if family && disp {
			protocol = xa.protocol
			data = xa.data
			return
		}
	}
	return
}
