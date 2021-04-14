package xproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	XProtocolVersion  Card16 = 11
	XProtocolRevision Card16 = 0
	XTCPPort                 = 6000
)

var (
	pack   = binary.Write
	endian = binary.LittleEndian
	MSB    = 0x42
	LSB    = 0x6c
)

type Card8 uint8
type Card16 uint16
type Card32 uint32

type Window Card32
type Drawable Card32
type Font Card32
type Pixmap Card32
type Cursor Card32
type Colormap Card32
type GContext Card32
type Atom Card32
type VisualID Card32
type Time Card32
type KeyCode uint8
type KeySym Card32

// pad creates a byteslice containing the data string padded for 4 byte alignment
// returns the length of the byte slice and the byte slice
func pad(data string) (int, []byte) {
	l := len(data)
	p := int((4 - uint(l)) % 4)
	padded := make([]byte, p+l)
	copy(padded, data)
	return len(data), padded
}

// XConnectionClientPrefix holds the data required to initiate handshake with the Xserver
type XConnectionClientPrefix struct {
	ByteOrder        Card8
	_                byte
	MajorVersion     Card16
	MinorVersion     Card16
	AuthProtoNameLen Card16
	AuthProtoDataLen Card16
	_                Card16
}

// NewXConnectionClientPrefix returns a byteslice representation of XConnectionClientPrefix
// followed by the auth protocol and data, required for a complete handshake initiation
func NewXConnectionClientPrefix(authName, authData string) ([]byte, error) {
	ln, pn := pad(authName)
	ld, pd := pad(authData)

	cs := new(XConnectionClientPrefix)
	cs.ByteOrder = Card8(LSB)
	cs.MajorVersion = XProtocolVersion
	cs.MinorVersion = XProtocolRevision
	cs.AuthProtoNameLen = Card16(ln)
	cs.AuthProtoDataLen = Card16(ld)

	buff := new(bytes.Buffer)
	if err := binary.Write(buff, endian, cs); err != nil {
		return []byte{}, fmt.Errorf("error while writing ClientPrefix to buffer: %w", err)
	}
	if n, err := buff.Write(pn); err != nil {
		return []byte{}, fmt.Errorf("error writing authname to buffer: %w", err)
	} else if n < ln {
		return []byte{}, fmt.Errorf("error writing name to buffer, only wrote %d of %d bytes", n, ln)
	}
	if n, err := buff.Write(pd); err != nil {
		return []byte{}, fmt.Errorf("error writing data to buffer: %w", err)
	} else if n < ld {
		return []byte{}, fmt.Errorf("error writing data to buffer, only wrote %d of %d bytes", n, ld)
	}
	return buff.Bytes(), nil
}

// XConnectionSetupPrefix status codes
const (
	XConnFailed  Card8 = iota // connection refused
	XConnSuccess              // connection accepted
	XConnAuth                 // further authentication required
)

// XConnectionSetupPrefix structure to hold the information recieved from the Xserver
// following handshake initiation
type XConnectionSetupPrefix struct {
	Status       Card8
	ReasonLen    byte
	MajorVersion Card16
	MinorVersion Card16
	DataLen      Card16
}

// XConnectionSetup structure holds setup information provided by Xserver at connection
// initiation
type XConnectionSetup struct {
	ReleaseNumber      Card32
	ResourceIDBase     Card32
	ResourceIDMask     Card32
	MotionBufferSize   Card32
	VendorLength       Card16
	MaxRequestLength   Card16
	NumRoots           Card8
	NumFormats         Card8
	ImageByteOrder     Card8
	BitmapBitOrder     Card8
	BitmapScanlineUnit Card8
	BitmapScanlinePad  Card8
	MinKeyCode         KeyCode
	MaxKeyCode         KeyCode
	_                  Card32
}

type XPixmapFormat struct {
	Depth        Card8
	BitsPerPixel Card8
	ScanlinePad  Card8
	_            Card8
	_            Card32
}

type XDepth struct {
	Depth      Card8
	_          Card8
	NumVisuals Card16
	_          Card32
}

type XVisualType struct {
	BitsPerRBG      Card8
	ColormapEntries Card16
	RedMask         Card32
	GreenMask       Card32
	BlueMask        Card32
	_               Card32
}

type XWindowRoot struct {
	WindowID         Window
	DefaultColormap  Colormap
	WhitePixel       Card32
	BlackPixel       Card32
	CurrentInputMask Card32
	PixelWidth       Card16
	PixelHeight      Card16
	MMWidth          Card16
	MMHeight         Card16
	MinInstalledMaps Card16
	MaxInstalledMaps Card16
	RootVisualID     VisualID
	BackingStore     Card8
	SaveUnders       bool
	RootDepth        Card8
	NumDeptsh        Card8
}
