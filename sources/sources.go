// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"errors"
	"github.com/samlecuyer/ecumene/geom"
	"github.com/samlecuyer/ecumene/query"
)

var ErrUnsupported = errors.New("Unsupported Format")

type DataSource interface {
	Query(*query.Query) chan geom.Shape
	Close()
}

func Open(ty, format, name string) (DataSource, error) {
	if ty == "file" && format == "shp" {
		return createShpSource(name)
	}
	if ty == "file" && format == "osm" {
		return createOsmSource(name)
	}
	return nil, ErrUnsupported
}
