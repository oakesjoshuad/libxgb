package xgbgen

import (
	"encoding/xml"
)

type element interface{}

type xcb struct {
	XMLName  xml.Name `xml:"xcb"`
	Header   xml.Attr `xml:"header,attr"`
	Elements []element
}

func (x *xcb) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x.XMLName = start.Name
	for _, attr := range start.Attr {
		if attr.Name.Local == "header" {
			x.Header = attr
		}
	}
	for {
		tkn, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tkn.(type) {
		case xml.StartElement:
			var elt interface{}
			switch t.Name.Local {
			case "struct":
				elt = new(xcbStruct)
			case "xidtype":
				elt = new(xcbXidType)
			case "xidunion":
				elt = new(xcbXidUnion)
			case "typedef":
				elt = new(xcbTypedef)
			default:
			}
			if err := d.DecodeElement(&elt, &t); err != nil {
				return err
			}
			x.Elements = append(x.Elements, elt)
		case xml.EndElement:
			if t == start.End() {
				return nil
			}
		}
	}
}

type xcbField struct {
	XMLName xml.Name `xml:"field"`
	Type    xml.Attr `xml:"type,attr"`
	Name    xml.Attr `xml:"name,attr"`
}

type xcbPad struct {
	XMLName xml.Name `xml:"pad"`
	Bytes   xml.Attr `xml:"bytes"`
}

type xcbStruct struct {
	XMLName  xml.Name  `xml:"struct"`
	Name     xml.Attr  `xml:"name,attr"`
	Elements []element `xml:",omitempty"`
}

func (x *xcbStruct) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	x.XMLName = start.Name
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			x.Name = attr
		}
	}
	for {
		tkn, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tkn.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "field":
				elt := new(xcbField)
				if err := d.DecodeElement(elt, &t); err != nil {
					return err
				}
				x.Elements = append(x.Elements, elt)
			case "pad":
				elt := new(xcbPad)
				if err := d.DecodeElement(elt, &t); err != nil {
					return err
				}
				x.Elements = append(x.Elements, elt)
			}
		case xml.EndElement:
			if t == start.End() {
				return nil
			}
		}
	}
}

type xcbType struct {
	XMLName xml.Name `xml:"type"`
	Text    string
}

type xcbItem struct {
	XMLName xml.Name `xml:"item"`
	Name    xml.Attr `xml:"name,attr"`
	Value   xcbValue `xml:",omitempty"`
	Bit     xcbBit   `xml:",omitempty"`
}

type xcbValue struct {
	XMLName xml.Name `xml:"value"`
	Value   string
}

type xcbBit struct {
	XMLName xml.Name `xml:"bit"`
	Bit     string
}

type xcbXidType struct {
	XMLName xml.Name `xml:"xidtype"`
	Name    xml.Attr `xml:"name,attr"`
}

type xcbXidUnion struct {
	XMLName xml.Name `xml:"xidunion"`
	Name    xml.Attr `xml:"name,attr"`
	Types   []xcbType
}

type xcbTypedef struct {
	XMLName xml.Name `xml:"typedef"`
	OldName xml.Attr `xml:"oldname,attr"`
	NewName xml.Attr `xml:"newname,attr"`
}

type xcbEnum struct {
	XMLName xml.Name `xml:"enum"`
	Name    xml.Attr `xml:"name,attr"`
	Items   []xcbItem
}
