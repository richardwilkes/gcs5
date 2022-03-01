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
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var _ unison.Border = &TitledBorder{}

// TitledBorder provides a titled line border that scales.
type TitledBorder struct {
	Owner unison.Paneler
	Title string
}

func (t *TitledBorder) scaleInsetsText(paint *unison.Paint) (scale float32, insets geom32.Insets, text *unison.Text) {
	scale = DetermineScale(t.Owner)
	text = unison.NewText(t.Title, &unison.TextDecoration{
		Font:  theme.PageLabelPrimaryFont.Face().Font(theme.PageLabelPrimaryFont.Size() * scale),
		Paint: paint,
	})
	insets = geom32.Insets{
		Top:    text.Height() + scale*2,
		Left:   scale,
		Bottom: scale,
		Right:  scale,
	}
	return
}

// Insets implements unison.Border
func (t *TitledBorder) Insets() geom32.Insets {
	_, insets, _ := t.scaleInsetsText(nil)
	return insets
}

// Draw implements unison.Border
func (t *TitledBorder) Draw(gc *unison.Canvas, rect geom32.Rect) {
	scale, insets, text := t.scaleInsetsText(theme.OnHeaderColor.Paint(gc, rect, unison.Fill))
	clip := rect
	clip.Inset(insets)
	path := unison.NewPath()
	path.SetFillType(unison.EvenOdd)
	path.Rect(rect)
	path.Rect(clip)
	gc.DrawPath(path, theme.HeaderColor.Paint(gc, rect, unison.Fill))
	text.Draw(gc, rect.X+(rect.Width-text.Width())/2, rect.Y+scale+text.Baseline())
}
