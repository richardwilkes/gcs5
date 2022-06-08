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
	"fmt"

	"github.com/richardwilkes/unison"
)

type dragDrawable struct {
	label *unison.Label
}

// NewTableDragDrawable creates a new drawable for a table row drag.
func NewTableDragDrawable(data *TableDragData, svg *unison.SVG, singularName, pluralName string) unison.Drawable {
	label := unison.NewLabel()
	label.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		r := rect
		r.Inset(unison.NewUniformInsets(1))
		corner := r.Height / 2
		gc.SaveWithOpacity(0.7)
		gc.DrawRoundedRect(r, corner, corner, data.Table.SelectionInk.Paint(gc, r, unison.Fill))
		gc.DrawRoundedRect(r, corner, corner, data.Table.OnSelectionInk.Paint(gc, r, unison.Stroke))
		gc.Restore()
		label.DefaultDraw(gc, rect)
	}
	label.OnBackgroundInk = data.Table.OnSelectionInk
	label.SetBorder(unison.NewEmptyBorder(unison.Insets{
		Top:    4,
		Left:   label.Font.LineHeight(),
		Bottom: 4,
		Right:  label.Font.LineHeight(),
	}))
	if count := countRows(data.Rows); count == 1 {
		label.Text = fmt.Sprintf("1 %s", singularName)
	} else {
		label.Text = fmt.Sprintf("%d %s", count, pluralName)
	}
	if svg != nil {
		baseline := label.Font.Baseline()
		label.Drawable = &unison.DrawableSVG{
			SVG:  svg,
			Size: unison.NewSize(baseline, baseline),
		}
	}
	_, pref, _ := label.Sizes(unison.Size{})
	label.SetFrameRect(unison.Rect{Size: pref})
	return &dragDrawable{label: label}
}

func (d *dragDrawable) LogicalSize() unison.Size {
	return d.label.FrameRect().Size
}

func (d *dragDrawable) DrawInRect(canvas *unison.Canvas, rect unison.Rect, _ *unison.SamplingOptions, _ *unison.Paint) {
	d.label.Draw(canvas, rect)
}

func countRows(rows []unison.TableRowData) int {
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildRows() {
			count += countRows(row.ChildRows())
		}
	}
	return count
}
