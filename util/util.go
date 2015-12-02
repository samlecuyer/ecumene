// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"math"
)

func Gps2webmerc(lng, lat float64) (float64, float64) {
	lat_rad := lat * math.Pi / 180
	return lng, math.Asinh(math.Tan(lat_rad)) * 180 / math.Pi
}

func Num2deg(x, y, z int) (lng, lat float64) {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lng = float64(x)/math.Exp2(float64(z))*360.0 - 180.0
	return
}

func Deg2num(lng, lat float64, zoom int) (x, y int) {
	n := math.Exp2(float64(zoom))
	x = int(math.Floor((lng + 180.0) / 360.0 * n))
	lat_rad := lat * math.Pi / 180
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat_rad)+1.0/math.Cos(lat_rad))/math.Pi) / 2.0 * n))
	return
}

type SymbolizerType uint

const (
	TextType SymbolizerType = 1 << iota
	PointType
	PathType
	PolygonType
)
