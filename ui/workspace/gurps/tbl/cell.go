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

package tbl

import (
	"strings"

	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// CellFromCellData creates a new panel for the given cell data.
func CellFromCellData(c *node.CellData, width float32, forPage, selected bool) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  c.Alignment,
	})
	var font unison.Font
	if forPage {
		font = theme.PageFieldPrimaryFont
	} else {
		font = unison.FieldFont
	}
	switch c.Type {
	case node.Text:
		addCellLabel(c, p, width, c.Primary, font, selected)
		if c.Secondary != "" {
			if forPage {
				font = theme.PageFieldSecondaryFont
			} else {
				font = theme.FieldSecondaryFont
			}
			addCellLabel(c, p, width, c.Secondary, font, selected)
		}
	case node.Toggle:
		// TODO: Implement!
	case node.PageRef:
		addPageRefLabel(p, c.Primary, c.Secondary, font, selected)
	}
	if c.Tooltip != "" {
		p.Tooltip = unison.NewTooltipWithText(c.Tooltip)
	}
	return p
}

func addCellLabel(c *node.CellData, parent *unison.Panel, width float32, text string, f unison.Font, selected bool) {
	decoration := &unison.TextDecoration{Font: f}
	var lines []*unison.Text
	if width > 0 {
		lines = unison.NewTextWrappedLines(text, decoration, width)
	} else {
		lines = unison.NewTextLines(text, decoration)
	}
	for _, line := range lines {
		label := unison.NewLabel()
		label.Text = line.String()
		label.Font = f
		label.HAlign = c.Alignment
		if selected {
			label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
		}
		parent.AddChild(label)
	}
}

func addPageRefLabel(parent *unison.Panel, text, highlight string, f unison.Font, selected bool) {
	label := unison.NewLabel()
	label.Font = f
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	parts := strings.FieldsFunc(text, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' })
	switch len(parts) {
	case 0:
	case 1:
		label.Text = parts[0]
	default:
		label.Text = parts[0] + "+"
		parent.Tooltip = unison.NewTooltipWithText(strings.Join(parts, "\n"))
	}
	if label.Text != "" {
		const isLinkKey = "is_link"
		label.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
			c := label.OnBackgroundInk
			if _, exists := label.ClientData()[isLinkKey]; exists {
				c = theme.OnLinkColor
				gc.DrawRect(rect, theme.LinkColor.Paint(gc, rect, unison.Fill))
			}
			unison.DrawLabel(gc, label.ContentRect(false), label.HAlign, label.VAlign, label.Text, label.Font, c,
				label.Drawable, label.Side, label.Gap, !label.Enabled())
		}
		parent.MouseEnterCallback = func(where geom32.Point, mod unison.Modifiers) bool {
			label.ClientData()[isLinkKey] = true
			label.MarkForRedraw()
			return true
		}
		parent.MouseExitCallback = func() bool {
			delete(label.ClientData(), isLinkKey)
			label.MarkForRedraw()
			return true
		}
		parent.MouseDownCallback = func(where geom32.Point, button, clickCount int, mod unison.Modifiers) bool {
			var list []string
			for _, one := range strings.FieldsFunc(text, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' }) {
				if one = strings.TrimSpace(one); one != "" {
					list = append(list, one)
				}
			}
			if len(list) != 0 {
				settings.OpenPageReference(parent.Window(), list[0], highlight)
			}
			return true
		}
	}
	parent.AddChild(label)
}
