// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapping

import (
	"code.google.com/p/sadbox/color"
	"encoding/xml"
	"fmt"
	"github.com/samlecuyer/ecumene/util"
)

type Style struct {
	Name  string  `xml:"name,attr"`
	Rules []*Rule `xml:"Rule"`
}

func (s *Style) String() string {
	return fmt.Sprintf("{%s %v}", s.Name, s.Rules)
}

type Rule struct {
	Filter      string
	Symbolizers map[util.SymbolizerType]Symbolizer
}

func (r *Rule) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r.Symbolizers = make(map[util.SymbolizerType]Symbolizer)
	for {
		e, err := d.Token()
		if err != nil {
			break
		}
		switch e := e.(type) {
		case xml.EndElement:
			if e.Name.Local == start.Name.Local {
				return nil
			}
		case xml.StartElement:
			switch e.Name.Local {
			case "Filter":
				var f string
				if err := d.DecodeElement(&f, &e); err != nil {
					return err
				}
				r.Filter = f
			case "Polygon":
				s := new(PolygonSymbolizer)
				if err := d.DecodeElement(s, &e); err != nil {
					return err
				}
				r.Symbolizers[util.PolygonType] = s
			case "Path":
				s := new(PathSymbolizer)
				if err := d.DecodeElement(s, &e); err != nil {
					return err
				}
				r.Symbolizers[util.PathType] = s
			case "Text":
				s := new(TextSymbolizer)
				if err := d.DecodeElement(s, &e); err != nil {
					return err
				}
				r.Symbolizers[util.TextType] = s
			}
		}
	}
	return d.Skip()
}

type Symbolizer interface {
	Name() string
}

type PolygonSymbolizer struct {
	Fill color.Hex `xml:"fill,attr"`
}

func (s *PolygonSymbolizer) Name() string {
	return "Polygon"
}

type PathSymbolizer struct {
	Weight float64   `xml:"width,attr" default:"0.5"`
	Stroke color.Hex `xml:"stroke,attr"`
}

func (s *PathSymbolizer) Name() string {
	return "Path"
}

type TextSymbolizer struct {
	Size float64   `xml:"size,attr"`
	Fill color.Hex `xml:"fill,attr"`
	Attr string    `xml:"name,attr"`
}

func (s *TextSymbolizer) Name() string {
	return "Text"
}
