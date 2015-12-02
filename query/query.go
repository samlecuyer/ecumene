// Copyright 2015 Sam L'ecuyer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package query

import (
	"fmt"
	"github.com/samlecuyer/ecumene/geom"
	"regexp"
	"strings"
)

type Assertion func(string) bool

// TODO: I need to find a way to query ahead of time
type Query struct {
	Bounds  geom.Bbox
	Filters []Filter
	Sel     *Select
}

func NewQuery(bb geom.Bbox) *Query {
	return &Query{Bounds: bb}
}

func (q *Query) Bounded(bb geom.Bbox) *Query {
	q.Bounds = bb
	return q
}

type Select struct {
	Fields []string `xml:"fields,attr"`
	// Table string `xml:"table,attr"`
}

// This is possibly the worst sql parser ever.
func (q *Query) Select(selstmt string) *Query {
	r, err := regexp.Compile(`^([a-zA-Z\,]+)`)
	if err != nil {
		fmt.Println(err)
	} else {
		w := r.FindAllStringSubmatch(selstmt, -1)
		if len(w) == 1 && len(w[0]) == 2 {
			q.Sel = &Select{strings.Split(w[0][1], ",")}
		}
	}
	return q
}

func doesExpr(part string, shape geom.Shape) bool {
	var name string
	var value string

	if strings.Contains(part, " = ") {
		parts := strings.Split(part, " = ")
		fmt.Sscanf(parts[0], "%s", &name)
		value = strings.TrimSpace(parts[1])
		return shape.Attribute(name) == value
	} else if strings.Contains(part, " contains ") {
		parts := strings.Split(part, " contains ")
		fmt.Sscanf(parts[0], "%s", &name)
		value = strings.TrimSpace(parts[1])
		attr := shape.Attribute(name)
		return strings.Contains(attr, value)
	} else {
		fmt.Sscanf(part, "%s", &name)
		return shape.Attribute(name) != ""
	}
	return false
}

func doesOr(filter string, shape geom.Shape) bool {
	for _, part := range strings.Split(string(filter), " or ") {
		if doesExpr(part, shape) {
			return true
		}
	}
	return false
}

func doesAnd(filter Filter, shape geom.Shape) bool {
	matches := false
	for i, part := range strings.Split(string(filter), " and ") {
		does := doesOr(part, shape)
		if i == 0 {
			matches = does
		} else {
			matches = does && matches
		}
	}
	return matches
}

type Filter string

func (f Filter) Applies(shape geom.Shape) bool {
	return f == "" || doesAnd(f, shape)
}
