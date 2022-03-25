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
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom"
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
func NewAdvantagesPageList(provider gurps.ListProvider) *PageList {
	p := newPageList(tbl.NewAdvantageTableHeaders(true), 0, 0,
		tbl.NewAdvantageRowData(func() []*gurps.Advantage { return provider.AdvantageList() }, true))
	p.installToggleStateHandler()
	p.installIncrementHandler()
	p.installDecrementHandler()
	return p
}

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(provider gurps.ListProvider) *PageList {
	p := newPageList(tbl.NewEquipmentTableHeaders(provider, true, true), 2, 2,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return provider.CarriedEquipmentList() }, true, true))
	p.installToggleStateHandler()
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	return p
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(provider gurps.ListProvider) *PageList {
	p := newPageList(tbl.NewEquipmentTableHeaders(provider, true, false), 1, 1,
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return provider.OtherEquipmentList() }, true, false))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	return p
}

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(provider gurps.SkillListProvider) *PageList {
	p := newPageList(tbl.NewSkillTableHeaders(provider, true), 0, 0, tbl.NewSkillRowData(provider, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementSkillHandler()
	p.installDecrementSkillHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	return p
}

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(provider gurps.SpellListProvider) *PageList {
	p := newPageList(tbl.NewSpellTableHeaders(provider, true), 0, 0, tbl.NewSpellRowData(provider, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementSkillHandler()
	p.installDecrementSkillHandler()
	return p
}

// NewNotesPageList creates the notes page list.
func NewNotesPageList(provider gurps.ListProvider) *PageList {
	return newPageList(tbl.NewNoteTableHeaders(true), 0, 0,
		tbl.NewNoteRowData(func() []*gurps.Note { return provider.NoteList() }, true))
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
	p.SetBorder(unison.NewLineBorder(theme.HeaderColor, 0, geom.NewUniformInsets[float32](1), false))
	p.table.DividerInk = theme.HeaderColor
	p.table.MinimumRowHeight = theme.PageFieldPrimaryFont.LineHeight()
	p.table.Padding.Top = 0
	p.table.Padding.Bottom = 0
	p.table.HierarchyColumnIndex = hierarchyColumnIndex
	p.table.HierarchyIndent = theme.PageFieldPrimaryFont.LineHeight()
	p.table.PreventUserColumnResize = true
	p.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range p.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom.Size[float32]{})
		pref.Width += p.table.Padding.Left + p.table.Padding.Right
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
	p.table.MouseDownCallback = func(where geom.Point[float32], button, clickCount int, mod unison.Modifiers) bool {
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
	p.tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, geom.Insets[float32]{Bottom: 1}, false)
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
	p.tableHeader.DrawCallback = func(gc *unison.Canvas, dirty geom.Rect[float32]) {
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

func (p *PageList) installToggleStateHandler() {
	p.installPerformHandlers(constants.ToggleStateItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch n.Data().(type) {
				case *gurps.Advantage:
					return true
				case *gurps.Equipment:
					return true
				}
			}
		}
		return false
	}, func() {
		var entity *gurps.Entity
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Advantage:
					item.Disabled = !item.Disabled
					entity = item.Entity
				case *gurps.Equipment:
					item.Equipped = !item.Equipped
					entity = item.Entity
				}
			}
		}
		if entity != nil {
			entity.Recalculate()
			widget.MarkModified(p)
		}
	})
}

func (p *PageList) installIncrementHandler() {
	p.installPerformHandlers(constants.IncrementItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
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
		for _, row := range p.table.SelectedRows(false) {
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
	p.installPerformHandlers(constants.DecrementItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
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
		for _, row := range p.table.SelectedRows(false) {
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
	p.installPerformHandlers(constants.IncrementUsesItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
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
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.Uses < e.MaxUses {
					e.Uses = xmath.Min(e.Uses+1, e.MaxUses)
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
	p.installPerformHandlers(constants.DecrementUsesItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
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
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				var e *gurps.Equipment
				if e, ok = n.Data().(*gurps.Equipment); ok && e.MaxUses > 0 && e.Uses > 0 {
					e.Uses = xmath.Max(e.Uses-1, 0)
					updated = true
				}
			}
		}
		if updated {
			widget.MarkModified(p)
		}
	})
}

func (p *PageList) installIncrementSkillHandler() {
	p.installPerformHandlers(constants.IncrementSkillLevelItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
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
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Skill:
					if !item.Container() {
						item.IncrementSkillLevel()
						entity = item.Entity
					}
				case *gurps.Spell:
					if !item.Container() {
						item.IncrementSkillLevel()
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

func (p *PageList) installDecrementSkillHandler() {
	p.installPerformHandlers(constants.DecrementSkillLevelItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
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
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Skill:
					if !item.Container() {
						item.DecrementSkillLevel()
						entity = item.Entity
					}
				case *gurps.Spell:
					if !item.Container() {
						item.DecrementSkillLevel()
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

func (p *PageList) installIncrementTechLevelHandler() {
	p.installPerformHandlers(constants.IncrementTechLevelItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Equipment:
					return true
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						return true
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						return true
					}
				}
			}
		}
		return false
	}, func() {
		var entity *gurps.Entity
		var changed bool
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Equipment:
					if item.TechLevel, changed = gurps.AdjustTechLevel(item.TechLevel, fixed.F64d4One); changed {
						entity = item.Entity
					}
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fixed.F64d4One); changed {
							entity = item.Entity
						}
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fixed.F64d4One); changed {
							entity = item.Entity
						}
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

func (p *PageList) installDecrementTechLevelHandler() {
	p.installPerformHandlers(constants.DecrementTechLevelItemID, func() bool {
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Equipment:
					return true
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						return true
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						return true
					}
				}
			}
		}
		return false
	}, func() {
		var entity *gurps.Entity
		var changed bool
		for _, row := range p.table.SelectedRows(false) {
			if n, ok := row.(*tbl.Node); ok {
				switch item := n.Data().(type) {
				case *gurps.Equipment:
					if item.TechLevel, changed = gurps.AdjustTechLevel(item.TechLevel, fxp.NegOne); changed {
						entity = item.Entity
					}
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fxp.NegOne); changed {
							entity = item.Entity
						}
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fxp.NegOne); changed {
							entity = item.Entity
						}
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

// Sync the underlying data.
func (p *PageList) Sync() {
	rows := p.table.SelectedRows(false)
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
	p.table.NeedsLayout = true
	p.NeedsLayout = true
	if parent := p.Parent(); parent != nil {
		parent.NeedsLayout = true
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
