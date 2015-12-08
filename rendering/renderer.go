// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rendering

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/samlecuyer/ecumene/geom"
	"github.com/samlecuyer/ecumene/mapping"
	"github.com/samlecuyer/ecumene/query"
	"github.com/samlecuyer/ecumene/util"
	"code.google.com/p/sadbox/color"
	"image"
	"image/draw"
	"log"
	"math"
	"reflect"
	"sync"
)

type Renderer struct {
	m             *mapping.Map
	width, height float64
	bbox          geom.Bbox
	layers        [][]geom.Shape
	matrix        draw2d.Matrix
	sync.Mutex
}

type SubImage interface {
	SubImage(r image.Rectangle) image.Image
}

func NewRenderer(m *mapping.Map, width, height int) *Renderer {
	r := &Renderer{
		m: m, width: float64(width), height: float64(height),
	}
	return r
}

func (r *Renderer) ClipTo(lng0, lat0, lng1, lat1 float64) {
	r.bbox[0] = lng0
	r.bbox[1] = lat0
	r.bbox[2] = lng1
	r.bbox[3] = lat1
	log.Println("clipped to: ", r.bbox)
}

func (r *Renderer) ClipToMap() error {
	b := r.m.Bounds()
	x0, y0, _ := r.m.Srs.Forward(b[0], b[1])
	x1, y1, _ := r.m.Srs.Forward(b[2], b[3])
	r.bbox = geom.Bbox{x0, y0, x1, y1}

	x0, y0, _ = r.m.Srs.Forward(b[0], b[3])
	x1, y1, _ = r.m.Srs.Forward(b[2], b[1])
	r.bbox = r.bbox.ExpandToFit(geom.Bbox{x0, y0, x1, y1})

	x0, y0, _ = r.m.Srs.Forward(0, b[3])
	x1, y1, _ = r.m.Srs.Forward(0, b[1])
	r.bbox = r.bbox.ExpandToFit(geom.Bbox{x0, y0, x1, y1})

	log.Println("clipped to: ", r.bbox)
	return nil
}

func (r *Renderer) DrawToFile(filename string) error {
	dest := r.Draw()
	return draw2dimg.SaveToPngFile(filename, dest)
}

func (r *Renderer) DrawAdjustedToFile(filename string) error {
	dest := r.Draw()
	if subimage, ok := dest.(SubImage); ok {
		bb := r.bbox
		x0, y0, x1, y1 := int(bb[0]), int(bb[1]), int(bb[2]), int(bb[3])
		dest = subimage.SubImage(image.Rect(x0, y0, x1, y1))
	}
	return draw2dimg.SaveToPngFile(filename, dest)
}

func (r *Renderer) Draw() image.Image {
	pixelsX, pixelsY := int(r.width), int(r.height)

	dest := image.NewRGBA(image.Rect(0, 0, pixelsX, pixelsY))
	draw.Draw(dest, dest.Bounds(), &image.Uniform{r.m.BgColor}, image.ZP, draw.Src)

	draw2d.SetFontFolder("/Library/Fonts/")
	draw2d.SetFontNamer(func(fontData draw2d.FontData) string {
		return fontData.Name + ".ttf"
	})
	gc := draw2dimg.NewGraphicContext(dest)
	// gc.DPI = 300

	gc.SetLineCap(draw2d.RoundCap)
	gc.SetFillColor(r.m.BgColor)
	gc.SetStrokeColor(r.m.Stroke)
	gc.SetFontData(draw2d.FontData{Name: "Georgia", Family: draw2d.FontFamilySerif, Style: draw2d.FontStyleNormal})

	dx := math.Abs(r.bbox[2] - r.bbox[0])
	dy := math.Abs(r.bbox[3] - r.bbox[1])

	pxf, pyf := float64(pixelsX), float64(pixelsY)
	r1, r2 := (pxf / dx), (pyf / dy)
	r0 := math.Min(r1, r2)
	w, h := dx*r0, dy*r0
	ox, oy := (pxf-w)/2, (pyf-h)/2
	img_box := [4]float64{ox, oy, ox + w, oy + h}

	r.matrix = draw2d.NewMatrixFromRects(r.bbox, img_box)

	for _, layer := range r.m.Layers {
		q := query.NewQuery(r.m.Bounds()).Select(layer.SourceQuery())
		if ds := layer.LoadSource(); ds != nil {
			defer ds.Close()
			for shp := range ds.Query(q) {
				var symbolizerType util.SymbolizerType
				switch shp.(type) {
				case geom.LineShape, geom.MultiLineShape:
					symbolizerType = util.PathType
				case geom.PolygonShape:
					symbolizerType = util.PolygonType
				}
				for _, symbolizer := range r.findSymbolizers(layer, symbolizerType) {
					symbolizer.Draw(gc, shp)
				}
				for _, symbolizer := range r.findSymbolizers(layer, util.TextType) {
					symbolizer.Draw(gc, shp)
				}
			}
		}
	}

	return dest
}

func (r *Renderer) graticule(gc draw2d.GraphicContext) {
	b := r.m.Bounds()
	d2r := math.Pi/180.
	gc.SetFillColor(AlphaHex("#ce4251"))
	// iterate over all the latitudes
	padding := 20 * d2r
	dxy := 0.001
	for phi := b[1]; phi > b[3]; phi -= padding {
		x, y, _ := r.m.Srs.Forward(b[0], phi)
		x, y = r.matrix.TransformPoint(x, y)
		gc.MoveTo(phi, b[0])
		for lam := b[0] + dxy; lam < b[2]; lam += dxy {
			x, y, _ = r.m.Srs.Forward(lam, phi)
			x, y = r.matrix.TransformPoint(x, y)
			gc.LineTo(x, y)
		}
		gc.Stroke()
	}
	for lam := b[0]; lam <= b[2]; lam += padding {
		x, y, _ := r.m.Srs.Forward(lam, b[1])
		x, y = r.matrix.TransformPoint(x, y)
		gc.MoveTo(lam, b[1])
		for phi := b[1] + dxy; phi >= b[3]; phi -= dxy {
			x, y, _ = r.m.Srs.Forward(lam, phi)
			x, y = r.matrix.TransformPoint(x, y)
			gc.LineTo(x, y)
		}
		gc.Stroke()
	}
}

func (r *Renderer) coordsAsPath(coords geom.Coordinates) *draw2d.Path {
	path := new(draw2d.Path)
	for i, point := range coords {
		x, y, _ := r.m.Srs.Forward(point[0], point[1])
		x, y = r.matrix.TransformPoint(x, y)
		if math.IsNaN(x) || math.IsInf(x, 1) {
			continue
		}
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	return path
}

func (r *Renderer) findSymbolizers(layer *mapping.Layer, filter util.SymbolizerType) []Symbolizer {
	var symbolizers []Symbolizer
	for _, styleName := range layer.Styles() {
		if style := r.m.FindStyle(styleName); style != nil {
			for _, rule := range style.Rules {
				if ps, ok := rule.Symbolizers[filter]; ok {
					switch specific := ps.(type) {
					case *mapping.PolygonSymbolizer:
						symbolizers = append(symbolizers, &PolygonSymbolizer{query.Filter(rule.Filter), r, specific})
					case *mapping.PathSymbolizer:
						symbolizers = append(symbolizers, &PathSymbolizer{query.Filter(rule.Filter), r, specific})
					case *mapping.TextSymbolizer:
						symbolizers = append(symbolizers, &TextSymbolizer{query.Filter(rule.Filter), r, specific})
					default:
						log.Println(reflect.TypeOf(ps).Elem())
					}
				}
			}
		}
	}
	return symbolizers
}
