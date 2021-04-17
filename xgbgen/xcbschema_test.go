package xgbgen

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestXProtoXML(t *testing.T) {
	filename := "schemas/test.xml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	xcbproto := xcb{}
	if err := xml.Unmarshal(data, &xcbproto); err != nil {
		t.Fatal(err)
	}
	t.Log(xcbproto.Header)
	for _, elt := range xcbproto.Elements {
		t.Log(elt)
	}
}
