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
	entity *gurps.Entity
}

// NewPage creates a new page.
func NewPage(entity *gurps.Entity) *Page {
	p := &Page{entity: entity}
	p.Self = p
	p.SetSizer(p.pageSizer)
	p.DrawCallback = p.drawSelf
	p.AdjustBorder()
	return p
}

func (p *Page) scaleSizeInsets() (scale float32, size geom32.Size, insets geom32.Insets) {
	scale = DetermineScale(p)
	sheetSettings := gurps.SheetSettingsFor(p.entity)
	w, h := sheetSettings.Page.Size.Dimensions()
	size.Width = w.Pixels() * scale
	size.Height = h.Pixels() * scale
	insets.Top = sheetSettings.Page.TopMargin.Pixels() * scale
	insets.Left = sheetSettings.Page.LeftMargin.Pixels() * scale
	insets.Bottom = sheetSettings.Page.BottomMargin.Pixels() * scale
	insets.Right = sheetSettings.Page.RightMargin.Pixels() * scale
	pH := theme.PageFooterPrimaryFont.Face().Font(theme.PageFooterPrimaryFont.Size() * scale).LineHeight()
	sH := theme.PageFooterSecondaryFont.Face().Font(theme.PageFooterSecondaryFont.Size() * scale).LineHeight()
	insets.Bottom += mathf32.Max(pH, sH) + sH
	return
}

// AdjustBorder applies the current scaling factor to the border.
func (p *Page) AdjustBorder() {
	_, _, insets := p.scaleSizeInsets()
	p.SetBorder(unison.NewEmptyBorder(insets))
}

func (p *Page) pageSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	_, pref, _ = p.scaleSizeInsets()
	return pref, pref, pref
}

func (p *Page) drawSelf(gc *unison.Canvas, _ geom32.Rect) {
	scale, size, insets := p.scaleSizeInsets()
	r := geom32.Rect{Size: size}
	gc.DrawRect(r, theme.PageColor.Paint(gc, r, unison.Fill))
	r.X += insets.Left
	r.Width -= insets.Left + insets.Right
	r.Y = r.Bottom() - insets.Bottom
	r.Height = insets.Bottom
	parent := p.Parent()
	pageNumber := parent.IndexOfChild(p) + 1

	primaryDecorations := &unison.TextDecoration{
		Font:  theme.PageFooterPrimaryFont.Face().Font(theme.PageFooterPrimaryFont.Size() * scale),
		Paint: theme.OnPageColor.Paint(gc, r, unison.Fill),
	}
	secondaryDecorations := &unison.TextDecoration{
		Font:  theme.PageFooterSecondaryFont.Face().Font(theme.PageFooterSecondaryFont.Size() * scale),
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
