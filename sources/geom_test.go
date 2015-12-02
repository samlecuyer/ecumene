// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"github.com/samlecuyer/ecumene/geom"
	"testing"
)

func TestPolygonTyping(t *testing.T) {
	var p geom.Shape
	p = new(shpPolygon)
	if _, ok := p.(geom.PolygonShape); !ok {
		t.Error("shpPolygon should be a PolygonShape")
	}
	p = new(shpPolygonZ)
	if _, ok := p.(geom.PolygonShape); !ok {
		t.Error("shpPolygonZ should be a PolygonShape")
	}
}
