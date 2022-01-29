/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package settings

import (
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/paper"
)

const (
	pagePaperSizeKey    = "paper_size"
	pageOrientationKey  = "orientation"
	pageTopMarginKey    = "top_margin"
	pageLeftMarginKey   = "left_margin"
	pageBottomMarginKey = "bottom_margin"
	pageRightMarginKey  = "right_margin"
)

// Page holds page settings.
type Page struct {
	Size         paper.Size
	Orientation  paper.Orientation
	TopMargin    paper.Length
	LeftMargin   paper.Length
	BottomMargin paper.Length
	RightMargin  paper.Length
}

// NewPage returns new settings with factory defaults.
func NewPage() *Page {
	return &Page{
		Size:         paper.Letter,
		Orientation:  paper.Portrait,
		TopMargin:    paper.Length{Length: 0.25, Units: paper.Inch},
		LeftMargin:   paper.Length{Length: 0.25, Units: paper.Inch},
		BottomMargin: paper.Length{Length: 0.25, Units: paper.Inch},
		RightMargin:  paper.Length{Length: 0.25, Units: paper.Inch},
	}
}

// NewPageFromJSON creates new settings from a JSON object.
func NewPageFromJSON(data map[string]interface{}) *Page {
	s := NewPage()
	s.Size = paper.SizeFromString(encoding.String(data[pagePaperSizeKey]))
	s.Orientation = paper.OrientationFromString(encoding.String(data[pageOrientationKey]))
	s.TopMargin = paper.LengthFromString(encoding.String(data[pageTopMarginKey]))
	s.LeftMargin = paper.LengthFromString(encoding.String(data[pageLeftMarginKey]))
	s.BottomMargin = paper.LengthFromString(encoding.String(data[pageBottomMarginKey]))
	s.RightMargin = paper.LengthFromString(encoding.String(data[pageRightMarginKey]))
	return s
}

// ToJSON emits this object as JSON.
func (s *Page) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(pagePaperSizeKey, s.Size.Key(), false, false)
	encoder.KeyedString(pageOrientationKey, s.Orientation.Key(), false, false)
	encoder.KeyedString(pageTopMarginKey, s.TopMargin.String(), false, false)
	encoder.KeyedString(pageLeftMarginKey, s.LeftMargin.String(), false, false)
	encoder.KeyedString(pageBottomMarginKey, s.BottomMargin.String(), false, false)
	encoder.KeyedString(pageRightMarginKey, s.RightMargin.String(), false, false)
	encoder.EndObject()
}
