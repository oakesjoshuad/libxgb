package libxGb

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

type clientsetup struct {
	MajorVersion     card16
	MinorVersion     card16
	AuthProtoNameLen card16
	AuthProtoDataLen card16
	AuthProtoName    string
	AuthProtoData    string
}

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

func NewClientSetup(host, display string) (cs *clientsetup, err error) {
	protoname, protodata, err := xauthinfo(host, display)
	if err != nil {
		return cs, err
	}
	cs.MajorVersion = xprotoversion
	cs.MinorVersion = xprotorevision
	cs.AuthProtoNameLen = card16(len(protoname))
	cs.AuthProtoName = protoname
	cs.AuthProtoDataLen = card16(len(protodata))
	cs.AuthProtoData = protodata
	return
}

func (cs *clientsetup) Pack() []byte {
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
