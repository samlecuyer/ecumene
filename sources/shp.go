// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"fmt"
	"github.com/samlecuyer/ecumene/geom"
	"github.com/samlecuyer/ecumene/query"
	"github.com/samlecuyer/ecumene/util"
	"github.com/samlecuyer/go-shp"
	"reflect"
	"strings"
)

type shpSource struct {
	r *shp.Reader
}

func (s *shpSource) Close() {
	s.r.Close()
}

func (s *shpSource) Query(q *query.Query) chan geom.Shape {
	ch := make(chan geom.Shape, 1000)
	go s.searchFor(q, ch)
	return ch
}

type shpPolygon struct {
	p     *shp.Polygon
	attrs map[string]string
}

func (p *shpPolygon) Attribute(s string) string {
	return p.attrs[s]
}

func (s *shpPolygon) Bbox() geom.Bbox {
	b := s.p.BBox()
	x0, y0 := util.Gps2webmerc(b.MinX, b.MaxY)
	x1, y1 := util.Gps2webmerc(b.MaxX, b.MinY)
	return geom.Bbox{x0, y0, x1, y1}
}

func (p *shpPolygon) Polygon() geom.Multiline {
	pgz := p.p
	lines := make(geom.Multiline, len(pgz.Parts))
	var length int32
	for i, idx := range pgz.Parts {
		ln := len(pgz.Parts)
		if i+1 >= ln {
			length = int32(len(pgz.Points)) - idx
		} else {
			length = pgz.Parts[i+1] - idx
		}
		if length == 0 {
			continue
		}

		lines[i] = make(geom.Coordinates, length)
		for j, point := range pgz.Points[idx : idx+length] {
			lng, lat := util.Gps2webmerc(point.X, point.Y)
			lines[i][j] = geom.Point{lng, lat}
		}
	}
	return lines
}

type shpPolygonZ struct {
	*shp.PolygonZ
	attrs map[string]string
}

func (p *shpPolygonZ) Attribute(s string) string {
	return p.attrs[s]
}

func (p *shpPolygonZ) Bbox() geom.Bbox {
	b := p.BBox()
	x0, y0 := util.Gps2webmerc(b.MinX, b.MaxY)
	x1, y1 := util.Gps2webmerc(b.MaxX, b.MinY)
	return geom.Bbox{x0, y0, x1, y1}
}

func (pgz *shpPolygonZ) Polygon() geom.Multiline {
	lines := make(geom.Multiline, len(pgz.Parts))
	var length int32
	for i, idx := range pgz.Parts {
		ln := len(pgz.Parts)
		if i+1 >= ln {
			length = int32(len(pgz.Points)) - idx
		} else {
			length = pgz.Parts[i+1] - idx
		}

		lines[i] = make(geom.Coordinates, length)
		for j, point := range pgz.Points[idx : idx+length] {
			lng, lat := util.Gps2webmerc(point.X, point.Y)
			lines[i][j] = geom.Point{lng, lat}
		}
	}
	return lines
}

type shpPolyLineM struct {
	*shp.PolyLineM
	attrs map[string]string
}

func (p *shpPolyLineM) Attribute(s string) string {
	return p.attrs[s]
}

func (p *shpPolyLineM) Bbox() geom.Bbox {
	b := p.BBox()
	x0, y0 := util.Gps2webmerc(b.MinX, b.MaxY)
	x1, y1 := util.Gps2webmerc(b.MaxX, b.MinY)
	return geom.Bbox{x0, y0, x1, y1}
}

func (pgz *shpPolyLineM) Paths() geom.Multiline {
	lines := make(geom.Multiline, len(pgz.Parts))
	var length int32
	for i, idx := range pgz.Parts {
		ln := len(pgz.Parts)
		if i+1 >= ln {
			length = int32(len(pgz.Points)) - idx
		} else {
			length = pgz.Parts[i+1] - idx
		}

		lines[i] = make(geom.Coordinates, length)
		for j, point := range pgz.Points[idx : idx+length] {
			lng, lat := util.Gps2webmerc(point.X, point.Y)
			lines[i][j] = geom.Point{lng, lat}
		}
	}
	return lines
}

type shpPolyLine struct {
	*shp.PolyLine
	attrs map[string]string
}

func (p *shpPolyLine) Attribute(s string) string {
	return p.attrs[s]
}

func (p *shpPolyLine) Bbox() geom.Bbox {
	b := p.BBox()
	x0, y0 := util.Gps2webmerc(b.MinX, b.MaxY)
	x1, y1 := util.Gps2webmerc(b.MaxX, b.MinY)
	return geom.Bbox{x0, y0, x1, y1}
}

func (pgz *shpPolyLine) Paths() geom.Multiline {
	lines := make(geom.Multiline, len(pgz.Parts))
	var length int32
	for i, idx := range pgz.Parts {
		ln := len(pgz.Parts)
		if i+1 >= ln {
			length = int32(len(pgz.Points)) - idx
		} else {
			length = pgz.Parts[i+1] - idx
		}

		lines[i] = make(geom.Coordinates, length)
		for j, point := range pgz.Points[idx : idx+length] {
			lng, lat := util.Gps2webmerc(point.X, point.Y)
			lines[i][j] = geom.Point{lng, lat}
		}
	}
	return lines
}

type shpPoint struct {
	x, y  float64
	attrs map[string]string
}

func (p *shpPoint) Attribute(s string) string {
	return p.attrs[s]
}

func (p *shpPoint) Bbox() geom.Bbox {
	x, y := util.Gps2webmerc(p.x, p.y)
	return geom.Bbox{x, y, x, y}
}

func (p *shpPoint) Point() geom.Point {
	lng, lat := util.Gps2webmerc(p.x, p.y)
	return geom.Point{lng, lat}
}

func (s *shpSource) searchFor(q *query.Query, ch chan geom.Shape) {
	defer close(ch)
	defer s.r.Close()

	b := s.r.BBox()
	x0, y0 := util.Gps2webmerc(b.MinX, b.MaxY)
	x1, y1 := util.Gps2webmerc(b.MaxX, b.MinY)
	b3 := geom.Bbox{x0, y0, x1, y1}

	if !b3.Overlaps(q.Bounds) {
		return
	}

	fields := make([]string, len(s.r.Fields()))

	var fieldsToGrab []int
	if q.Sel != nil {
		var f string
		for i, field := range s.r.Fields() {
			f = strings.Trim(field.String(), " \x00")
			fields[i] = f
			for _, name := range q.Sel.Fields {
				if name == f {
					fieldsToGrab = append(fieldsToGrab, i)
					fmt.Println(f)
				}
			}
		}
	}

	for s.r.Next() {
		n, p := s.r.Shape()
		b := p.BBox()
		x0, y0 = util.Gps2webmerc(b.MinX, b.MaxY)
		x1, y1 = util.Gps2webmerc(b.MaxX, b.MinY)
		b3 = geom.Bbox{x0, y0, x1, y1}

		if b3.Overlaps(q.Bounds) {
			attrs := make(map[string]string)
			if q.Sel != nil {
				for _, i := range fieldsToGrab {
					if val := s.r.ReadAttribute(n, i); val != "" {
						attrs[fields[i]] = val
					}
				}
			}
			switch underlying := p.(type) {
			case *shp.Polygon:
				ch <- &shpPolygon{underlying, attrs}
			case *shp.PolygonZ:
				ch <- &shpPolygonZ{underlying, attrs}
			case *shp.PolyLine:
				ch <- &shpPolyLine{underlying, attrs}
			case *shp.PolyLineM:
				ch <- &shpPolyLineM{underlying, attrs}
			case *shp.Point:
				ch <- &shpPoint{underlying.X, underlying.Y, attrs}
			default:
				fmt.Println(reflect.TypeOf(p).Elem())
			}
		}
	}
}

func createShpSource(name string) (DataSource, error) {
	f, err := shp.Open(name)
	if err != nil {
		return nil, err
	}
	return &shpSource{f}, nil
}
