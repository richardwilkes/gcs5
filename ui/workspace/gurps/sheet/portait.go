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

package sheet

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Portrait displays and allows editing of an Entity's portrait.
type Portrait struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPortrait creates a new portrait panel.
func NewPortrait(entity *gurps.Entity) *Portrait {
	p := &Portrait{entity: entity}
	p.Self = p
	p.SetSizer(p.portraitSizer)
	p.SetBorder(&TitledBorder{Title: i18n.Text("Portrait")})
	p.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text(`Double-click to set a character portrait, or drag an image onto this block.

The dimensions of the chosen picture should be in a ratio of 3 pixels wide
for every 4 pixels tall to scale without distortion.

Dimensions of %dx%d are ideal.`), gurps.PortraitWidth*2, gurps.PortraitHeight*2))
	p.DrawCallback = p.drawSelf
	return p
}

func (p *Portrait) portraitSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	insets := p.Border().Insets()
	pref.Width = insets.Left + insets.Right + gurps.PortraitWidth
	pref.Height = insets.Top + insets.Bottom + gurps.PortraitHeight
	return pref, pref, pref
}

func (p *Portrait) drawSelf(gc *unison.Canvas, r geom32.Rect) {
	r = p.ContentRect(false)
	paint := unison.ContentColor.Paint(gc, r, unison.Fill)
	gc.DrawRect(r, paint)
	if img := p.entity.Profile.Portrait(); img != nil {
		img.DrawInRect(gc, r, nil, paint)
	}
}
