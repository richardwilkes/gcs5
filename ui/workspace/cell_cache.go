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

package workspace

import (
	"strings"

	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const (
	savedInkKey   = "saved_ink"
	currentInkKey = "current_ink"
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
	var lines []string
	if width > 0 {
		lines = f.WrapText(text, width)
	} else {
		lines = strings.Split(text, "\n")
	}
	for _, line := range lines {
		label := unison.NewLabel()
		label.Text = line
		label.Font = f
		if selected {
			label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
		}
		parent.AddChild(label)
	}
}

func newPageReferenceHeader() unison.TableColumnHeader {
	header := unison.NewTableColumnHeader("", i18n.Text(`A reference to the book and page the item appears on (e.g. B22 would refer to "Basic Set", page 22)`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  icons.BookmarkSVG(),
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

func createAndAddPageRefCellLabel(parent *unison.Panel, text string, f unison.Font, selected bool) {
	label := unison.NewLabel()
	label.Font = f
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	parts := strings.FieldsFunc(text, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' })
	label.Text = parts[0]
	if len(parts) > 1 {
		label.Text += "+"
	}
	if label.Text != "" {
		label.MouseEnterCallback = func(where geom32.Point, mod unison.Modifiers) bool {
			clientData := label.ClientData()
			clientData[savedInkKey] = label.LabelTheme.OnBackgroundInk
			label.LabelTheme.OnBackgroundInk = theme.AccentColor
			clientData[currentInkKey] = label.LabelTheme.OnBackgroundInk
			label.MarkForRedraw()
			return false
		}
		label.MouseExitCallback = func() bool {
			clientData := label.ClientData()
			ink, ok := clientData[savedInkKey].(unison.Ink)
			if !ok {
				ink = label.LabelTheme.OnBackgroundInk
			}
			label.LabelTheme.OnBackgroundInk = ink
			delete(clientData, savedInkKey)
			delete(clientData, currentInkKey)
			label.MarkForRedraw()
			return false
		}
	}
	parent.AddChild(label)
}
