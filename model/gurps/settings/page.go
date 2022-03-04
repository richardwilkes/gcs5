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
	"github.com/richardwilkes/gcs/model/paper"
)

// Page holds page settings.
type Page struct {
	Size         paper.Size        `json:"paper_size"`
	Orientation  paper.Orientation `json:"orientation"`
	TopMargin    paper.Length      `json:"top_margin"`
	LeftMargin   paper.Length      `json:"left_margin"`
	BottomMargin paper.Length      `json:"bottom_margin"`
	RightMargin  paper.Length      `json:"right_margin"`
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

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (p *Page) EnsureValidity() {
	p.Size = p.Size.EnsureValid()
	p.Orientation = p.Orientation.EnsureValid()
	p.TopMargin.EnsureValidity()
	p.LeftMargin.EnsureValidity()
	p.BottomMargin.EnsureValidity()
	p.RightMargin.EnsureValidity()
}

// Clone a copy of this.
func (p *Page) Clone() *Page {
	clone := *p
	return &clone
}
