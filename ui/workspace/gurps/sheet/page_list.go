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
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/menus"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath"
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
	canPerformMap        map[int]func() bool
	performMap           map[int]func()
}

// NewAdvantagesPageList creates the advantages page list.
func NewAdvantagesPageList(entity *gurps.Entity) *PageList {
	p := newPageList(tbl.NewAdvantageTableHeaders(true), 0, 0,
		tbl.NewAdvantageRowData(func() []*gurps.Advantage { return entity.Advantages }, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	return p
}

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(entity *gurps.Entity) *PageList {
	p := newPageList(tbl.NewEquipmentTableHeaders(entity, true, true), 2, 2,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return entity.CarriedEquipment }, true, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	return p
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(entity *gurps.Entity) *PageList {
	p := newPageList(tbl.NewEquipmentTableHeaders(entity, true, false), 1, 1,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return entity.OtherEquipment }, true, false))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	return p
}

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(entity *gurps.Entity) *PageList {
	p := newPageList(tbl.NewSkillTableHeaders(true), 0, 0,
		tbl.NewSkillRowData(func() []*gurps.Skill { return entity.Skills }, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	return p
}

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(entity *gurps.Entity) *PageList {
	p := newPageList(tbl.NewSpellTableHeaders(true), 0, 0,
		tbl.NewSpellRowData(func() []*gurps.Spell { return entity.Spells }, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	return p
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
		canPerformMap:        make(map[int]func() bool),
		performMap:           make(map[int]func()),
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
	p.table.MouseDownCallback = func(where geom32.Point, button, clickCount int, mod unison.Modifiers) bool {
		p.table.RequestFocus()
		return p.table.DefaultMouseDown(where, button, clickCount, mod)
	}
	p.table.FrameChangeCallback = func() {
		p.table.SizeColumnsToFitWithExcessIn(excessWidthIndex)
	}
	p.table.CanPerformCmdCallback = func(_ interface{}, id int) bool {
		if f, ok := p.canPerformMap[id]; ok {
			return f()
		}
		return false
	}
	p.table.PerformCmdCallback = func(_ interface{}, id int) {
		if f, ok := p.performMap[id]; ok {
			f()
		}
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

func (p *PageList) installPerformHandlers(id int, can func() bool, do func()) {
	p.canPerformMap[id] = can
	p.performMap[id] = do
}

func (p *PageList) installIncrementHandler() {
	p.installPerformHandlers(menus.IncrementItemID, func() bool {
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Advantage:
					if item.IsLeveled() {
						return true
					}
				case *gurps.Equipment:
					return true
				case *gurps.Skill:
					if !item.Container() {
						return true
					}
				case *gurps.Spell:
					if !item.Container() {
						return true
					}
				}
			}
		}
		return false
	}, func() {
		var entity *gurps.Entity
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Advantage:
					if item.IsLeveled() {
						levels := increment(*item.Levels)
						item.Levels = &levels
						entity = item.Entity
					}
				case *gurps.Equipment:
					item.Quantity = increment(item.Quantity)
					entity = item.Entity
				case *gurps.Skill:
					if !item.Container() {
						item.Points = increment(item.Points)
						entity = item.Entity
					}
				case *gurps.Spell:
					if !item.Container() {
						item.Points = increment(item.Points)
						entity = item.Entity
					}
				}
			}
		}
		if entity != nil {
			entity.Recalculate()
			widget.MarkModified(p)
		}
	})
}

func (p *PageList) installDecrementHandler() {
	p.installPerformHandlers(menus.DecrementItemID, func() bool {
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Advantage:
					if item.IsLeveled() && *item.Levels > 0 {
						return true
					}
				case *gurps.Equipment:
					if item.Quantity > 0 {
						return true
					}
				case *gurps.Skill:
					if !item.Container() && item.Points > 0 {
						return true
					}
				case *gurps.Spell:
					if !item.Container() && item.Points > 0 {
						return true
					}
				}
			}
		}
		return false
	}, func() {
		var entity *gurps.Entity
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Advantage:
					if item.IsLeveled() {
						levels := decrement(*item.Levels)
						item.Levels = &levels
						entity = item.Entity
					}
				case *gurps.Equipment:
					item.Quantity = decrement(item.Quantity).Max(0)
					entity = item.Entity
				case *gurps.Skill:
					if !item.Container() {
						item.Points = decrement(item.Points).Max(0)
						entity = item.Entity
					}
				case *gurps.Spell:
					if !item.Container() {
						item.Points = decrement(item.Points).Max(0)
						entity = item.Entity
					}
				}
			}
		}
		if entity != nil {
			entity.Recalculate()
			widget.MarkModified(p)
		}
	})
}

func increment(value fixed.F64d4) fixed.F64d4 {
	return value.Trunc() + fixed.F64d4One
}

func decrement(value fixed.F64d4) fixed.F64d4 {
	v := value.Trunc()
	if v == value {
		v -= fixed.F64d4One
	}
	return v
}

func (p *PageList) installIncrementUsesHandler() {
	p.installPerformHandlers(menus.IncrementUsesItemID, func() bool {
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.Uses < e.MaxUses {
					return true
				}
			}
		}
		return false
	}, func() {
		updated := false
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.Uses < e.MaxUses {
					e.Uses = xmath.MinInt(e.Uses+1, e.MaxUses)
					updated = true
				}
			}
		}
		if updated {
			widget.MarkModified(p)
		}
	})
}

func (p *PageList) installDecrementUsesHandler() {
	p.installPerformHandlers(menus.DecrementUsesItemID, func() bool {
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.MaxUses > 0 && e.Uses > 0 {
					return true
				}
			}
		}
		return false
	}, func() {
		updated := false
		for _, row := range p.table.SelectedRows() {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.MaxUses > 0 && e.Uses > 0 {
					e.Uses = xmath.MaxInt(e.Uses-1, 0)
					updated = true
				}
			}
		}
		if updated {
			widget.MarkModified(p)
		}
	})
}

// Sync the underlying data.
func (p *PageList) Sync() {
	rows := p.table.SelectedRows()
	selection := make(map[node.Node]bool, len(rows))
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			selection[n.Data()] = true
		}
	}
	p.table.SetTopLevelRows(p.topLevelRowsCallback(p.table))
	if len(selection) != 0 {
		_, indexes := p.collectRowMappings(0, make([]int, 0, len(selection)), selection, p.table.TopLevelRows())
		if len(indexes) != 0 {
			p.table.SelectByIndex(indexes...)
		}
	}
}

func (p *PageList) collectRowMappings(index int, indexes []int, selection map[node.Node]bool, rows []unison.TableRowData) (updatedIndex int, updatedIndexes []int) {
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			if selection[n.Data()] {
				indexes = append(indexes, index)
			}
		}
		index++
		if row.IsOpen() {
			index, indexes = p.collectRowMappings(index, indexes, selection, row.ChildRows())
		}
	}
	return index, indexes
}
