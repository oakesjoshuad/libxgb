// libxau locates, parses the Xauthority file
package libXau

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const (
	slashDotXauthority  = "/.Xauthority"
	FamilyLocal         = 256
	FamilyWild          = 65535
	FamilyNetname       = 254
	FamilyKrb5Principal = 253
	FamilyLocalHost     = 252
)

var (
	errXauthorityLocation = errors.New("Unable to locate Xauthority file")
)

// XAuthInfo
type XAuthInfo struct {
	Family  uint16
	Address string
	Number  string
	Name    string
	Data    string
}

// XauFileName returns the the location of the .Xauthority file
func XauFileName() (string, error) {
	if filename := os.Getenv("XAUTHORITY"); filename != "" {
		return filename, nil
	}
	if home := os.Getenv("HOME"); home != "" {
		return home + slashDotXauthority, nil
	}
	return "", errXauthorityLocation
}

// XauReadAuth returns XAuthInfo
func XauReadAuth(filename string) (xa *XAuthInfo, err error) {
	xa = new(XAuthInfo)
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	if xa.Family, err = readShort(fd); err != nil {
		return
	}
	if xa.Address, err = readString(fd); err != nil {
		return
	}
	if xa.Number, err = readString(fd); err != nil {
		return
	}
	if xa.Name, err = readString(fd); err != nil {
		return
	}
	if xa.Data, err = readString(fd); err != nil {
		return
	}

	return
}

// read_short reads 2 bytes from the buffer and returns
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
