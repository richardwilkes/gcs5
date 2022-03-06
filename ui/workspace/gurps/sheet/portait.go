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

// PortraitPanel holds the contents of the portrait block on the sheet.
type PortraitPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPortraitPanel creates a new portrait panel.
func NewPortraitPanel(entity *gurps.Entity) *PortraitPanel {
	p := &PortraitPanel{entity: entity}
	p.Self = p
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.StartAlignment,
		VAlign: unison.StartAlignment,
		VSpan:  2,
	})
	p.SetSizer(p.portraitSizer)
	p.SetBorder(&TitledBorder{Title: i18n.Text("Portrait")})
	p.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text(`Double-click to set a character portrait, or drag an image onto this block.

The dimensions of the chosen picture should be in a ratio of 3 pixels wide
for every 4 pixels tall to scale without distortion.

Recommended minimum dimensions are %dx%d.`), gurps.PortraitWidth*2, gurps.PortraitHeight*2))
	p.DrawCallback = p.drawSelf
	return p
}

func (p *PortraitPanel) portraitSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	var width, height float32
	insets := p.Border().Insets()
	parent := p.Parent()
	for parent != nil {
		if sheet, ok := parent.Self.(*Sheet); ok {
			_, idPanelPref, _ := sheet.IdentityPanel.Sizes(geom32.Size{})
			_, descPanelPref, _ := sheet.DescriptionPanel.Sizes(geom32.Size{})
			height = idPanelPref.Height + 1 + descPanelPref.Height
			break
		}
		parent = parent.Parent()
	}
	if height -= insets.Top + insets.Bottom; height > 0 {
		width = height * 0.75
	} else {
		width = gurps.PortraitWidth
		height = gurps.PortraitHeight
	}
	pref.Width = insets.Left + insets.Right + width
	pref.Height = insets.Top + insets.Bottom + height
	return pref, pref, pref
}

func (p *PortraitPanel) drawSelf(gc *unison.Canvas, r geom32.Rect) {
	r = p.ContentRect(false)
	paint := unison.ContentColor.Paint(gc, r, unison.Fill)
	gc.DrawRect(r, paint)
	if img := p.entity.Profile.Portrait(); img != nil {
		img.DrawInRect(gc, r, nil, paint)
	}
}
