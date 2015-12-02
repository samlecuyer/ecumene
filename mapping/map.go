// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapping

import (
	"code.google.com/p/sadbox/color"
	"encoding/xml"
	"fmt"
	"github.com/samlecuyer/ecumene/geom"
	"os"
)

type Map struct {
	XMLName xml.Name `xml:"Map"`
	Styles  []*Style `xml:"Style"`
	Layers  []*Layer `xml:"Layer"`

	Stroke  color.Hex `xml:"stroke,attr"`
	BgColor color.Hex `xml:"bgcolor,attr"`

	bounds geom.Bbox
}

type Include struct {
	XMLName xml.Name `xml:"Include"`
	Styles  []*Style `xml:"Style"`
	Layers  []*Layer `xml:"Layer"`
}

func (m *Map) Bounds() geom.Bbox {
	return m.bounds
}

func (m *Map) setBounds(extent string) (err error) {
	_, err = fmt.Sscanf(extent, "%f %f %f %f", &m.bounds[0], &m.bounds[3], &m.bounds[2], &m.bounds[1])
	return
}

func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "extent" {
			m.setBounds(attr.Value)
		}
		if attr.Name.Local == "bgcolor" {
			m.BgColor = color.Hex(attr.Value)
		}
	}
	for {
		e, err := d.Token()
		if err != nil {
			return err
		}
		switch e := e.(type) {
		case xml.EndElement:
			if e.Name.Local == start.Name.Local {
				return nil
			}
		case xml.StartElement:
			switch e.Name.Local {
			case "Style":
				style := new(Style)
				if err := d.DecodeElement(style, &e); err != nil {
					return err
				}
				m.Styles = append(m.Styles, style)
			case "Layer":
				layer := new(Layer)
				if err := d.DecodeElement(layer, &e); err != nil {
					return err
				}
				m.Layers = append(m.Layers, layer)
			case "Include":
				var name string
				if err := d.DecodeElement(&name, &e); err != nil {
					return err
				}
				include, err := loadInclude(name)
				if err != nil {
					return err
				}
				m.Styles = append(m.Styles, include.Styles...)
				m.Layers = append(m.Layers, include.Layers...)
			}
		}
	}
	return d.Skip()
}

func NewMap(path string) (*Map, error) {
	return loadFile(path)
}

func loadFile(path string) (*Map, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := new(Map)
	decoder := xml.NewDecoder(file)
	decoder.Decode(m)
	return m, nil
}

func loadInclude(path string) (*Include, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := new(Include)
	decoder := xml.NewDecoder(file)
	decoder.Decode(m)
	return m, nil
}

func (m *Map) FindStyle(name string) *Style {
	for _, style := range m.Styles {
		if style.Name == name {
			return style
		}
	}
	return nil
}
