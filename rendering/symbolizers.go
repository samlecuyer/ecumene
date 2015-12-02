// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rendering

import (
	"github.com/llgcode/draw2d"
	"github.com/samlecuyer/ecumene/geom"
	"github.com/samlecuyer/ecumene/mapping"
	"github.com/samlecuyer/ecumene/query"
	"math"
)

type Symbolizer interface {
	Applies(shape geom.Shape) bool
	Draw(gc draw2d.GraphicContext, shape geom.Shape)
}

type PolygonSymbolizer struct {
	query.Filter
	r *Renderer
	s *mapping.PolygonSymbolizer
}

func (ps *PolygonSymbolizer) Draw(gc draw2d.GraphicContext, shape geom.Shape) {
	if !ps.Applies(shape) {
		return
	}
	if polygon, ok := shape.(geom.PolygonShape); ok {
		gc.SetFillColor(ps.s.Fill)
		for _, path := range polygon.Polygon() {
			l := ps.r.coordsAsPath(path)
			gc.Fill(l)
		}
	}
}

type PathSymbolizer struct {
	query.Filter
	r *Renderer
	s *mapping.PathSymbolizer
}

func (ps *PathSymbolizer) Draw(gc draw2d.GraphicContext, shape geom.Shape) {
	if !ps.Applies(shape) {
		return
	}
	gc.SetStrokeColor(ps.s.Stroke)
	gc.SetLineWidth(ps.s.Weight)
	switch specific := shape.(type) {
	case geom.LineShape:
		l := ps.r.coordsAsPath(specific.Path())
		gc.Stroke(l)
	case geom.MultiLineShape:
		for _, path := range specific.Paths() {
			l := ps.r.coordsAsPath(path)
			gc.Stroke(l)
		}
	}
}

type TextSymbolizer struct {
	query.Filter
	r *Renderer
	s *mapping.TextSymbolizer
}

func (ts *TextSymbolizer) Draw(gc draw2d.GraphicContext, shape geom.Shape) {
	if !ts.Applies(shape) {
		return
	}
	if name := shape.Attribute(ts.s.Attr); name != "" {
		gc.SetFontSize(ts.s.Size)
		l, t, r, b := gc.GetStringBounds(name)
		bb := shape.Bbox()
		dx := math.Abs(bb[2] - bb[0])
		dy := math.Abs(bb[3] - bb[1])
		ox, oy := dx/2, dy/2

		x, y := ts.r.matrix.TransformPoint(bb[0]+ox, bb[1]-oy)
		gc.SetFillColor(ts.s.Fill)
		gc.FillStringAt(name, x-(r-l)/2, y-(t-b)/2)
	}
}
