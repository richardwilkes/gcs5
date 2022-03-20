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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// PageList holds a list for a sheet page.
type PageList struct {
	unison.Panel
	tableHeader          *unison.TableHeader
	table                *unison.Table
	topLevelRowsCallback func(table *unison.Table) []unison.TableRowData
}

// NewAdvantagesPageList creates the advantages page list.
func NewAdvantagesPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewAdvantageTableHeaders(true), 0, 0,
		tbl.NewAdvantageRowData(func() []*gurps.Advantage { return entity.Advantages }, true))
}

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewEquipmentTableHeaders(entity, true, true), 2, 2,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return entity.CarriedEquipment }, true, true))
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewEquipmentTableHeaders(entity, true, false), 1, 1,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return entity.OtherEquipment }, true, false))
}

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewSkillTableHeaders(true), 0, 0,
		tbl.NewSkillRowData(func() []*gurps.Skill { return entity.Skills }, true))
}

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewSpellTableHeaders(true), 0, 0,
		tbl.NewSpellRowData(func() []*gurps.Spell { return entity.Spells }, true))
}

// NewNotesPageList creates the notes page list.
func NewNotesPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewNoteTableHeaders(true), 0, 0,
		tbl.NewNoteRowData(func() []*gurps.Note { return entity.Notes }, true))
}

// NewConditionalModifiersPageList creates the conditional modifiers page list.
func NewConditionalModifiersPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewConditionalModifierTableHeaders(i18n.Text("Condition")), -1, 1,
		tbl.NewConditionalModifierRowData(entity, false))
}

// NewReactionsPageList creates the reaction modifiers page list.
func NewReactionsPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewConditionalModifierTableHeaders(i18n.Text("Reaction")), -1, 1,
		tbl.NewConditionalModifierRowData(entity, true))
}

// NewMeleeWeaponsPageList creates the melee weapons page list.
func NewMeleeWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewWeaponTableHeaders(true), -1, 0,
		tbl.NewWeaponRowData(func() []*gurps.Weapon { return entity.EquippedWeapons(weapon.Melee) }, true))
}

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(tbl.NewWeaponTableHeaders(false), -1, 0,
		tbl.NewWeaponRowData(func() []*gurps.Weapon { return entity.EquippedWeapons(weapon.Ranged) }, false))
}

func newPageList(columnHeaders []unison.TableColumnHeader, hierarchyColumnIndex, excessWidthIndex int, topLevelRows func(table *unison.Table) []unison.TableRowData) *PageList {
	p := &PageList{
		table:                unison.NewTable(),
		topLevelRowsCallback: topLevelRows,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	p.SetBorder(unison.NewLineBorder(theme.HeaderColor, 0, geom32.NewUniformInsets(1), false))
	p.table.DividerInk = theme.HeaderColor
	p.table.MinimumRowHeight = theme.PageFieldPrimaryFont.LineHeight()
	p.table.Padding.Top = 0
	p.table.Padding.Bottom = 0
	p.table.HierarchyColumnIndex = hierarchyColumnIndex
	p.table.HierarchyIndent = theme.PageFieldPrimaryFont.LineHeight()
	p.table.PreventUserColumnResize = true
	p.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range p.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom32.Size{})
		p.table.ColumnSizes[i].AutoMinimum = pref.Width
		p.table.ColumnSizes[i].AutoMaximum = 800
		p.table.ColumnSizes[i].Minimum = pref.Width
		p.table.ColumnSizes[i].Maximum = 10000
	}
	p.table.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	p.table.FrameChangeCallback = func() {
		p.table.SizeColumnsToFitWithExcessIn(excessWidthIndex)
	}
	p.tableHeader = unison.NewTableHeader(p.table, columnHeaders...)
	p.tableHeader.BackgroundInk = theme.HeaderColor
	p.tableHeader.DividerInk = theme.HeaderColor
	p.tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, geom32.Insets{Bottom: 1}, false)
	p.tableHeader.SetBorder(p.tableHeader.HeaderBorder)
	p.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fixed.F64d4FromString(s1); err == nil {
			var n2 fixed.F64d4
			if n2, err = fixed.F64d4FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	p.tableHeader.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.tableHeader.DrawCallback = func(gc *unison.Canvas, dirty geom32.Rect) {
		sortedOn := -1
		for i, hdr := range p.tableHeader.ColumnHeaders {
			if hdr.SortState().Order == 0 {
				sortedOn = i
				break
			}
		}
		if sortedOn != -1 {
			gc.DrawRect(dirty, p.tableHeader.BackgroundInk.Paint(gc, dirty, unison.Fill))
			r := p.tableHeader.ColumnFrame(sortedOn)
			r.X -= p.table.Padding.Left
			r.Width += p.table.Padding.Left + p.table.Padding.Right
			gc.DrawRect(r, theme.MarkerColor.Paint(gc, r, unison.Fill))
			save := p.tableHeader.BackgroundInk
			p.tableHeader.BackgroundInk = unison.Transparent
			p.tableHeader.DefaultDraw(gc, dirty)
			p.tableHeader.BackgroundInk = save
		} else {
			p.tableHeader.DefaultDraw(gc, dirty)
		}
	}
	p.table.SetTopLevelRows(p.topLevelRowsCallback(p.table))
	p.AddChild(p.tableHeader)
	p.AddChild(p.table)
	return p
}

// Sync the underlying data.
func (p *PageList) Sync() {
	p.table.SetTopLevelRows(p.topLevelRowsCallback(p.table))
}
