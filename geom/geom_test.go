// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geom

import (
	"testing"
)

func TestBboxOverlap(t *testing.T) {
	bb := Bbox{-118.944862413904, 34.823301, -117.646374, 32.801462}
	b3 := Bbox{-118.30078125, 36.21093749999999, -118.125, 36.03515625}
	if bb.Overlaps(b3) {
		t.Error("these shapes don't overlap")
	}
}
