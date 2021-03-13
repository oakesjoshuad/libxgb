package libxgb

import (
	"bytes"
	"encoding/binary"
)

const (
	XProtocolVersion  card16 = 11
	XProtocolRevision card16 = 0
)

var (
	pack      = binary.Write
	endianess = binary.BigEndian
)

type request interface {
	pack() []byte
}

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
type ClientSetup struct {
	MajorVersion     card16
	MinorVersion     card16
	AuthProtoNameLen card16
	AuthProtoDataLen card16
	AuthProtoName    string
	AuthProtoData    string
}

// NewClientSetup ...
func NewClientSetup(authName, authData string) []byte {
	var cs ClientSetup
	cs.MajorVersion = XProtocolVersion
	cs.MinorVersion = XProtocolRevision
	cs.AuthProtoNameLen = card16(len(authName))
	cs.AuthProtoName = authName
	cs.AuthProtoDataLen = card16(len(authData))
	cs.AuthProtoData = authData
	return cs.pack()
}

// Pack ...
func (cs *ClientSetup) pack() []byte {
	var buf bytes.Buffer
	pack(&buf, endianess, binary.BigEndian)
	pack(&buf, endianess, pad(1))
	pack(&buf, endianess, cs.MajorVersion)
	pack(&buf, endianess, cs.MinorVersion)
	pack(&buf, endianess, cs.AuthProtoNameLen)
	pack(&buf, endianess, cs.AuthProtoDataLen)
	pack(&buf, endianess, pad(2))
	pack(&buf, endianess, cs.AuthProtoName)
	pack(&buf, endianess, pad(cs.AuthProtoName))
	pack(&buf, endianess, cs.AuthProtoName)
	pack(&buf, endianess, pad(cs.AuthProtoData))
	return buf.Bytes()
}
