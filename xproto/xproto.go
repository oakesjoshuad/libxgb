package xproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	XProtocolVersion  Card16 = 11
	XProtocolRevision Card16 = 0
)

var (
	pack = binary.Write
	MSB  = 0x42
	LSB  = 0x6c
)

type Card8 uint8
type Card16 uint16
type Card32 uint32

// pad pads data to align with 4 byte units
func pad(data string) (int, []byte) {
	l := len(data)
	p := int((4 - uint(l)) % 4)
	padded := make([]byte, p+l)
	copy(padded, data)
	return len(data), padded
}

// ClientSetup ...
type ClientPrefix struct {
	ByteOrder        Card8
	_                byte
	MajorVersion     Card16
	MinorVersion     Card16
	AuthProtoNameLen Card16
	AuthProtoDataLen Card16
	_                Card16
}

// NewClientSetup ...
func NewClientPrefix(authName, authData string) ([]byte, error) {
	ln, pn := pad(authName)
	ld, pd := pad(authData)

	cs := new(ClientPrefix)
	cs.ByteOrder = Card8(LSB)
	cs.MajorVersion = XProtocolVersion
	cs.MinorVersion = XProtocolRevision
	cs.AuthProtoNameLen = Card16(len(authName))
	cs.AuthProtoDataLen = Card16(len(authData))
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, cs); err != nil {
		return []byte{}, fmt.Errorf("Encountered error while writing ClientPrefix to buffer: %w", err)
	}
	if n, err := buff.Write(pn); err != nil {
		return []byte{}, fmt.Errorf("Encountered error writing authname to buffer: %w", err)
	} else if n < ln {
		return []byte{}, fmt.Errorf("Encountered error writing name to buffer, only wrote %d of %d bytes", n, ln)
	}
	if n, err := buff.Write(pd); err != nil {
		return []byte{}, fmt.Errorf("Encountered error writing data to buffer: %w", err)
	} else if n < ld {
		return []byte{}, fmt.Errorf("Encountered error writing data to buffer, only wrote %d of %d bytes", n, ld)
	}
	return buff.Bytes(), nil
}

type SetupPrefix struct {
	Success      Card8
	ReasonLen    byte
	MajorVersion Card16
	MinorVersion Card16
	AddBytesLen  Card16
}
