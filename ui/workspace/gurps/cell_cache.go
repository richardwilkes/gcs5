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
	"strings"

	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type cellCache struct {
	width float32
	data  string
	panel *unison.Panel
}

func (c *cellCache) matches(width float32, data string) bool {
	return c != nil && c.panel != nil && c.width == width && c.data == data
}

func createAndAddCellLabel(parent *unison.Panel, width float32, text string, f unison.Font, selected bool) {
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
		if selected {
			label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
		}
		parent.AddChild(label)
	}
}

func newPageReferenceHeader() unison.TableColumnHeader {
	header := unison.NewTableColumnHeader("")
	header.Tooltip = unison.NewTooltipWithText(i18n.Text(`A reference to the book and page the item appears on
e.g. B22 would refer to "Basic Set", page 22`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  res.BookmarkSVG,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

func createAndAddPageRefCellLabel(parent *unison.Panel, text, highlight string, f unison.Font, selected bool) {
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
