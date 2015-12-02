// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geom

type Shape interface {
	Geometry
	Attribute(string) string
}

type PointShape interface {
	Shape
	Point() Point
}

type LineShape interface {
	Shape
	Path() Coordinates
}

type MultiLineShape interface {
	Shape
	Paths() Multiline
}

type PolygonShape interface {
	Shape
	Polygon() Multiline
}
