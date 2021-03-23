package xproto

import (
	"encoding/binary"
)

const (
	XProtocolVersion  card16 = 11
	XProtocolRevision card16 = 0
)

var (
	pack = binary.Write
	LSB  = binary.LittleEndian
)

type card8 uint8
type card16 uint16
type card32 uint32

// pad pads data to align with 4 byte units
func pad(expr interface{}) (padding []byte) {
	switch v := expr.(type) {
	case int:
		padding = make([]byte, v)
	case card16:
	case card32:
		padding = make([]byte, (4-(uint(v)%4))%4)
	case string:
		padding = make([]byte, (4-(uint(len(v))))%4)
	}
	return
}

// ClientSetup ...
type ClientPrefix struct {
	MajorVersion     card16
	MinorVersion     card16
	_                byte
	AuthProtoNameLen card16
	AuthProtoDataLen card16
	_                card16
	AuthProtoName    *[]byte
	AuthProtoData    *[]byte
}

// NewClientSetup ...
func NewClientPrefix(authName, authData string) ClientPrefix {
	pAuthName := pad(authName)
	pAuthData := pad(authData)
	var cs ClientPrefix
	cs.MajorVersion = XProtocolVersion
	cs.MinorVersion = XProtocolRevision
	cs.AuthProtoNameLen = card16(len(pAuthName))
	cs.AuthProtoName = &pAuthName
	cs.AuthProtoDataLen = card16(len(pAuthData))
	cs.AuthProtoData = &pAuthData
	return cs
}

type SetupPrefix struct {
	Success      card8
	ReasonLen    byte
	MajorVersion card16
	MinorVersion card16
	AddBytesLen  card16
}
