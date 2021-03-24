package xproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	XProtocolVersion  Card16 = 11
	XProtocolRevision Card16 = 0
)

var (
	pack = binary.Write
	LSB  = binary.LittleEndian
)

type Card8 uint8
type Card16 uint16
type Card32 uint32

// pad pads data to align with 4 byte units
func pad(data string) []byte {
	padding := make([]byte, (4-(uint(len(data))))%4)
	copy(padding, []byte(data))
	return padding
}

// ClientSetup ...
type ClientPrefix struct {
	MajorVersion     Card16
	MinorVersion     Card16
	_                byte
	AuthProtoNameLen Card16
	AuthProtoDataLen Card16
	_                Card16
}

// NewClientSetup ...
func NewClientPrefix(authName, authData string) ([]byte, error) {
	pAuthName := pad(authName)
	pAuthData := pad(authData)
	cs := new(ClientPrefix)
	cs.MajorVersion = XProtocolVersion
	cs.MinorVersion = XProtocolRevision
	cs.AuthProtoNameLen = Card16(len(pAuthName))
	cs.AuthProtoDataLen = Card16(len(pAuthData))

	buff := new(bytes.Buffer)
	if err := binary.Write(buff, LSB, cs); err != nil {
		return []byte{}, fmt.Errorf("Error writing client prefix to buffer: %w", err)
	}
	fmt.Println(buff.Bytes())
	if n, err := buff.Write(pAuthName); err != nil {
		return []byte{}, fmt.Errorf("Error writing AuthName to buffer: %w", err)
	} else if n < len(pAuthName) {
		return []byte{}, errors.New("Error writing AuthName to buffer")
	}
	if n, err := buff.Write(pAuthData); err != nil {
		return buff.Bytes(), fmt.Errorf("Error writing AuthData to buffer: %w", err)
	} else if n < len(pAuthData) {
		return []byte{}, errors.New("Error writing AuthData to buffer")
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
