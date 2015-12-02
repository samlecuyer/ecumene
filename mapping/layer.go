// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapping

import (
	"encoding/xml"
	"github.com/samlecuyer/ecumene/sources"
)

type Layer struct {
	styles []string    `xml:"StyleName"`
	source *Datasource `xml:"Datasource"`
}

func (l *Layer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
			case "StyleName":
				var f string
				if err := d.DecodeElement(&f, &e); err != nil {
					return err
				}
				l.styles = append(l.styles, f)
			case "Datasource":
				l.source = new(Datasource)
				if err := d.DecodeElement(l.source, &e); err != nil {
					return err
				}
			}
		}
	}
	return d.Skip()
}

func (l *Layer) LoadSource() sources.DataSource {
	source := l.source
	if ds, err := sources.Open(source.Type, source.Format, source.Val); err == nil {
		return ds
	}
	return nil
}

func (l *Layer) Styles() []string {
	return l.styles
}

// TODO: remove this, this is just to keep compatibility while I'm refactoring
func (l *Layer) SourceQuery() string {
	return l.source.Query
}

type Datasource struct {
	Type   string `xml:"type,attr"`
	Format string `xml:"format,attr"`
	Val    string `xml:"name,attr"`
	Query  string `xml:"Query"`
}
