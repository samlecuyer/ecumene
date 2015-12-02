// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/samlecuyer/ecumene/mapping"
	"github.com/samlecuyer/ecumene/rendering"

	"flag"
	"github.com/pkg/profile"
	"os"
)

func main() {
	flag.Parse()

	input := flag.Arg(0)
	output := flag.Arg(1)
	cwd, _ := os.Getwd()
	defer profile.Start(profile.ProfilePath(cwd)).Stop()

	m, _ := mapping.NewMap(input)

	r := rendering.NewRenderer(m, 5000, 5000)

	r.ClipToMap()
	r.DrawToFile(output)
}
