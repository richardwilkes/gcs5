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

package widget

import (
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/unison"
)

var _ unison.Border = &TitledBorder{}

// TitledBorder provides a titled line border.
type TitledBorder struct {
	Title string
	Font  unison.Font
}

func (t *TitledBorder) font() unison.Font {
	if t.Font == nil {
		return theme.PageLabelPrimaryFont
	}
	return t.Font
}

// Insets implements unison.Border
func (t *TitledBorder) Insets() unison.Insets {
	return unison.Insets{
		Top:    t.font().LineHeight() + 2,
		Left:   1,
		Bottom: 1,
		Right:  1,
	}
}

// Draw implements unison.Border
func (t *TitledBorder) Draw(gc *unison.Canvas, rect unison.Rect) {
	clip := rect
	clip.Inset(t.Insets())
	path := unison.NewPath()
	path.SetFillType(unison.EvenOdd)
	path.Rect(rect)
	path.Rect(clip)
	gc.DrawPath(path, theme.HeaderColor.Paint(gc, rect, unison.Fill))
	text := unison.NewText(t.Title, &unison.TextDecoration{
		Font:  t.font(),
		Paint: theme.OnHeaderColor.Paint(gc, rect, unison.Fill),
	})
	text.Draw(gc, rect.X+(rect.Width-text.Width())/2, rect.Y+1+text.Baseline())
}
