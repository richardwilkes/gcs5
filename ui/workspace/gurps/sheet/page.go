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
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/menus"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// Page holds a logical page worth of content.
type Page struct {
	unison.Panel
	flex   *unison.FlexLayout
	entity *gurps.Entity
}

// NewPage creates a new page.
func NewPage(entity *gurps.Entity) *Page {
	p := &Page{
		entity: entity,
		flex: &unison.FlexLayout{
			Columns:  1,
			HSpacing: 1,
			VSpacing: 1,
		},
	}
	p.Self = p
	p.SetBorder(unison.NewEmptyBorder(p.insets()))
	p.SetSizer(p.pageSizer)
	p.SetLayout(p)
	p.DrawCallback = p.drawSelf
	return p
}

// LayoutSizes implements unison.Layout
func (p *Page) LayoutSizes(_ *unison.Panel, _ geom32.Size) (min, pref, max geom32.Size) {
	pref = p.prefSize()
	return pref, pref, pref
}

// PerformLayout implements unison.Layout
func (p *Page) PerformLayout(_ *unison.Panel) {
	p.flex.PerformLayout(p.AsPanel())
}

func (p *Page) prefSize() geom32.Size {
	w, h := gurps.SheetSettingsFor(p.entity).Page.Size.Dimensions()
	return geom32.Size{
		Width:  w.Pixels(),
		Height: h.Pixels(),
	}
}

func (p *Page) insets() geom32.Insets {
	sheetSettings := gurps.SheetSettingsFor(p.entity)
	insets := geom32.Insets{
		Top:    sheetSettings.Page.TopMargin.Pixels(),
		Left:   sheetSettings.Page.LeftMargin.Pixels(),
		Bottom: sheetSettings.Page.BottomMargin.Pixels(),
		Right:  sheetSettings.Page.RightMargin.Pixels(),
	}
	height := theme.PageFooterSecondaryFont.LineHeight()
	insets.Bottom += mathf32.Max(theme.PageFooterPrimaryFont.LineHeight(), height) + height
	return insets
}

func (p *Page) pageSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	pref = p.prefSize()
	return pref, pref, pref
}

func (p *Page) drawSelf(gc *unison.Canvas, _ geom32.Rect) {
	insets := p.insets()
	r := geom32.Rect{Size: p.prefSize()}
	gc.DrawRect(r, theme.PageColor.Paint(gc, r, unison.Fill))
	r.X += insets.Left
	r.Width -= insets.Left + insets.Right
	r.Y = r.Bottom() - insets.Bottom
	r.Height = insets.Bottom
	parent := p.Parent()
	pageNumber := parent.IndexOfChild(p) + 1

	primaryDecorations := &unison.TextDecoration{
		Font:  theme.PageFooterPrimaryFont,
		Paint: theme.OnPageColor.Paint(gc, r, unison.Fill),
	}
	secondaryDecorations := &unison.TextDecoration{
		Font:  theme.PageFooterSecondaryFont,
		Paint: primaryDecorations.Paint,
	}

	var title string
	if gurps.SheetSettingsFor(p.entity).UseTitleInFooter {
		title = p.entity.Profile.Title
	} else {
		title = p.entity.Profile.Name
	}
	center := unison.NewText(title, primaryDecorations)
	left := unison.NewText(fmt.Sprintf(i18n.Text("%s is copyrighted ©%s by %s"), cmdline.AppName,
		cmdline.CopyrightYears, cmdline.CopyrightHolder), secondaryDecorations)
	right := unison.NewText(fmt.Sprintf(i18n.Text("Modified %s"), p.entity.ModifiedOn), secondaryDecorations)
	if pageNumber&1 == 0 {
		left, right = right, left
	}
	y := r.Y + mathf32.Max(mathf32.Max(left.Baseline(), right.Baseline()), center.Baseline())
	left.Draw(gc, r.X, y)
	center.Draw(gc, r.X+(r.Width-center.Width())/2, y)
	right.Draw(gc, r.Right()-right.Width(), y)
	y = r.Y + mathf32.Max(mathf32.Max(left.Height(), right.Height()), center.Height())

	center = unison.NewText(menus.WebSiteDomain, secondaryDecorations)
	left = unison.NewText(i18n.Text("All rights reserved"), secondaryDecorations)
	right = unison.NewText(fmt.Sprintf(i18n.Text("Page %d of %d"), pageNumber, len(parent.Children())), secondaryDecorations)
	if pageNumber&1 == 0 {
		left, right = right, left
	}
	y += mathf32.Max(mathf32.Max(left.Baseline(), right.Baseline()), center.Baseline())
	left.Draw(gc, r.X, y)
	center.Draw(gc, r.X+(r.Width-center.Width())/2, y)
	right.Draw(gc, r.Right()-right.Width(), y)
}