// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/bmizerany/pat"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/samlecuyer/ecumene/mapping"
	"github.com/samlecuyer/ecumene/rendering"
	"github.com/samlecuyer/ecumene/util"
)

func main() {
	flag.Parse()

	input := flag.Arg(0)

	m, _ := mapping.NewMap(input)
	r := rendering.NewRenderer(m, 256, 256)

	mux := pat.New()
	mux.Get("/:z/:x/:y.png", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		z_str := q.Get(":z")
		x_str := q.Get(":x")
		y_str := q.Get(":y")

		x, _ := strconv.Atoi(x_str)
		y, _ := strconv.Atoi(y_str)
		z, _ := strconv.Atoi(z_str)

		lng0, lat0 := util.Num2deg(x, y, z)
		lng0, lat0 = util.Gps2webmerc(lng0, lat0)
		lng1, lat1 := util.Num2deg(x+1, y+1, z)
		lng1, lat1 = util.Gps2webmerc(lng1, lat1)

		r.Lock()
		r.ClipTo(lng0, lat0, lng1, lat1)
		tile := r.Draw()
		r.Unlock()
		png.Encode(w, tile)
	}))

	http.Handle("/", mux)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
