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

package gurps

import (
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type templateContent struct {
	unison.Panel
	flex *unison.FlexLayout
}

// newTemplateContent creates a new page.
func newTemplateContent() *templateContent {
	p := &templateContent{
		flex: &unison.FlexLayout{
			Columns:  1,
			VSpacing: 1,
		},
	}
	p.Self = p
	p.SetBorder(unison.NewEmptyBorder(geom32.NewUniformInsets(4)))
	p.SetLayout(p)
	return p
}

func (p *templateContent) LayoutSizes(_ *unison.Panel, _ geom32.Size) (min, pref, max geom32.Size) {
	w, _ := settings.Global().Sheet.Page.Size.Dimensions()
	_, size, _ := p.flex.LayoutSizes(p.AsPanel(), geom32.Size{Width: w.Pixels()})
	pref.Width = w.Pixels()
	pref.Height = size.Height
	return pref, pref, pref
}

func (p *templateContent) PerformLayout(_ *unison.Panel) {
	p.flex.PerformLayout(p.AsPanel())
}

// ApplyPreferredSize to this panel.
func (p *templateContent) ApplyPreferredSize() {
	r := p.FrameRect()
	_, pref, _ := p.Sizes(geom32.Size{})
	r.Size = pref
	p.SetFrameRect(r)
	p.ValidateLayout()
}
