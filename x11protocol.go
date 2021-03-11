// Package libxgb ...
package libxgb

import (
	"bytes"
	"encoding/binary"
)

const (
	xprotoversion  card16 = 11
	xprotorevision card16 = 0
)

var (
	pack      = binary.Write
	endianess = binary.BigEndian
)

type request interface {
	Pack() request
}

type response interface{}

// connection setup
type card8 uint8
type card16 uint16
type card32 uint32

// ClientSetup ...
type ClientSetup struct {
	MajorVersion     card16
	MinorVersion     card16
	AuthProtoNameLen card16
	AuthProtoDataLen card16
	AuthProtoName    string
	AuthProtoData    string
}

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

// NewClientSetup ...
// func NewClientSetup(host, display string) (cs *ClientSetup, err error) {
// 	xa, err := xau.Xauth(host, display)
// 	if err != nil {
// 		return cs, err
// 	}
// 	cs.MajorVersion = xprotoversion
// 	cs.MinorVersion = xprotorevision
// 	cs.AuthProtoNameLen = card16(len(xa.protocol))
// 	cs.AuthProtoName = xa.protocol
// 	cs.AuthProtoDataLen = card16(len(xa.data))
// 	cs.AuthProtoData = xa.data
// 	return
// }

// Pack ...
func (cs *ClientSetup) Pack() []byte {
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
