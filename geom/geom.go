// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geom

import (
	"math"
)

type Geometry interface {
	Bbox() Bbox
}

type Bbox [4]float64

func (bb Bbox) ExpandToFit(other Bbox) Bbox {
	return Bbox{
		math.Min(bb[0], other[0]),
		math.Max(bb[1], other[1]),
		math.Max(bb[2], other[2]),
		math.Min(bb[3], other[3]),
	}
}

func (r Bbox) Overlaps(s Bbox) bool {
	// r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
	// r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
	return r[0] < s[2] && s[0] < r[2] &&
		r[3] < s[1] && s[3] < r[1]
}

type Point [2]float64

func (p *Point) Bbox() Bbox {
	return Bbox{p[0], p[1], p[0], p[1]}
}

type Coordinates []Point
type Multiline []Coordinates

func (points Multiline) Bbox() Bbox {
	p := points[0]
	bb := p.Bbox()
	for i := 1; i < len(points); i++ {
		bb = bb.ExpandToFit(points[i].Bbox())
	}
	return bb
}

func (points Coordinates) Bbox() Bbox {
	p := points[0]
	bb := p.Bbox()
	for i := 1; i < len(points); i++ {
		bb = bb.ExpandToFit(points[i].Bbox())
	}
	return bb
}
