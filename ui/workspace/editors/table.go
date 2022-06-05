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

package editors

import (
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

func newTable(parent *unison.Panel, provider TableProvider) {
	table := unison.NewTable()
	table.DividerInk = theme.HeaderColor
	table.Padding.Top = 0
	table.Padding.Bottom = 0
	table.HierarchyColumnIndex = provider.HierarchyColumnIndex()
	table.HierarchyIndent = unison.FieldFont.LineHeight()
	table.MinimumRowHeight = unison.FieldFont.LineHeight()
	headers := provider.Headers()
	widget.TableSetupColumnSizes(table, headers)
	table.SetLayoutData(&unison.FlexLayoutData{
		MinSize: geom.Size[float32]{Height: unison.FieldFont.LineHeight()},
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
		HGrab:   true,
		VGrab:   true,
	})
	widget.TableInstallStdCallbacks(table)
	table.FrameChangeCallback = func() {
		table.SizeColumnsToFitWithExcessIn(provider.ExcessWidthColumnIndex())
	}
	tableHeader := widget.TableCreateHeader(table, headers)
	tableHeader.BackgroundInk = theme.HeaderColor
	tableHeader.DividerInk = theme.HeaderColor
	tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, unison.Insets{Bottom: 1}, false)
	tableHeader.SetBorder(tableHeader.HeaderBorder)
	tableHeader.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	table.SetTopLevelRows(provider.RowData(table))
	table.InstallCmdHandlers(constants.OpenEditorItemID, func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.OpenEditor(widget.FindRebuildable(table), table) })
	table.InstallCmdHandlers(constants.OpenOnePageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenPageRef(table) })
	table.InstallCmdHandlers(constants.OpenEachPageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenEachPageRef(table) })
	parent.AddChild(tableHeader)
	parent.AddChild(table)
}
