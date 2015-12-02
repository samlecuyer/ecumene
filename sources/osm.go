// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"encoding/xml"
	"fmt"
	"github.com/samlecuyer/ecumene/geom"
	"github.com/samlecuyer/ecumene/query"
	"github.com/samlecuyer/ecumene/util"
	"image"
	"math"
	"os"
)

type osmSource struct {
	osm *Osm
}

func (s *osmSource) Close() {}

func (s *osmSource) Query(q *query.Query) chan geom.Shape {
	ch := make(chan geom.Shape)
	go s.searchFor(q, ch)
	return ch
}

func createOsmSource(name string) (DataSource, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := new(Osm)
	decoder := xml.NewDecoder(file)
	decoder.Decode(m)
	return &osmSource{m}, nil
}

func (s *osmSource) searchFor(q *query.Query, ch chan geom.Shape) {
	// TODO: we need to actually work
	defer close(ch)
}

type RefId int

type Node struct {
	Id   RefId   `xml:"id,attr"`
	Lat  float64 `xml:"lat,attr"`
	Lng  float64 `xml:"lon,attr"`
	Tags []*Tag  `xml:"tag"`
}

func (n *Node) String() string {
	return fmt.Sprintf("(%v, %v)", n.Lng, n.Lat)
}

type Relation struct {
	Id      RefId     `xml:"id,attr"`
	Members []*Member `xml:"member"`
	Tags    []*Tag    `xml:"tag"`
}

func (r *Relation) Name() string {
	for _, t := range r.Tags {
		if t.K == "name" {
			return t.V
		}
	}
	return ""
}

func (r *Relation) String() string {
	return fmt.Sprintf("%s (%d)", r.Name(), len(r.Members))
}

type Way struct {
	Id    RefId     `xml:"id,attr"`
	Nodes []NodeRef `xml:"nd"`
	Tags  []*Tag    `xml:"tag"`
}

func (w *Way) Name() string {
	for _, t := range w.Tags {
		if t.K == "name" {
			return t.V
		}
	}
	return "unknown"
}

func (w *Way) String() string {
	return fmt.Sprintf("%s (%d)", w.Name(), len(w.Nodes))
}

type NodeRef struct {
	Id RefId `xml:"ref,attr"`
}

type Member struct {
	Ref  RefId  `xml:"ref,attr"`
	Type string `xml:"type,attr"`
	Role string `xml:"role,attr"`
}

type Tag struct {
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

type Bounds struct {
	Minlat float64 `xml:"minlat,attr"`
	MinLng float64 `xml:"minlon,attr"`
	Maxlat float64 `xml:"maxlat,attr"`
	MaxLng float64 `xml:"maxlon,attr"`
}

func (b *Bounds) ComputeTiles(zoom int) Grid {
	sx, sy := util.Deg2num(math.Min(b.MinLng, b.MaxLng), math.Min(b.Minlat, b.Maxlat), zoom)
	ex, ey := util.Deg2num(math.Max(b.MinLng, b.MaxLng), math.Max(b.Minlat, b.Maxlat), zoom)
	return Grid{sx, sy, ex, ey, zoom}
}

func (b *Bounds) TileRect(zoom int) image.Rectangle {
	sx, sy := util.Deg2num(b.MinLng, b.Minlat, zoom)
	ex, ey := util.Deg2num(b.MaxLng, b.Maxlat, zoom)
	return image.Rect(sx, sy, ex, ey)
}

type Grid struct {
	sx, sy, ex, ey int
	zoom           int
}

type Osm struct {
	XMLName   xml.Name    `xml:"osm"`
	Bounds    Bounds      `xml:"bounds"`
	Nodes     []*Node     `xml:"node"`
	Ways      []*Way      `xml:"way"`
	Relations []*Relation `xml:"relation"`
}

func (osm *Osm) NodesMaps() map[RefId]*Node {
	m := make(map[RefId]*Node)
	for _, t := range osm.Nodes {
		m[t.Id] = t
	}
	return m
}

type coord struct {
	lng float64
	lat float64
}
