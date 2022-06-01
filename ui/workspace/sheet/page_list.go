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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

var _ widget.Syncer = &PageList{}

// PageList holds a list for a sheet page.
type PageList struct {
	unison.Panel
	tableHeader   *unison.TableHeader
	table         *unison.Table
	provider      tbl.TableProvider
	canPerformMap map[int]func() bool
	performMap    map[int]func()
}

// NewAdvantagesPageList creates the advantages page list.
func NewAdvantagesPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, tbl.NewAdvantagesProvider(provider, true))
	p.installToggleStateHandler()
	p.installIncrementHandler()
	p.installDecrementHandler()
	return p
}

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, tbl.NewEquipmentProvider(provider, true, true))
	p.installToggleStateHandler()
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	p.installConvertToContainerHandler(owner)
	return p
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, tbl.NewEquipmentProvider(provider, true, false))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementUsesHandler()
	p.installDecrementUsesHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	p.installConvertToContainerHandler(owner)
	return p
}

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, tbl.NewSkillsProvider(provider, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementSkillHandler()
	p.installDecrementSkillHandler()
	p.installIncrementTechLevelHandler()
	p.installDecrementTechLevelHandler()
	return p
}

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(owner widget.Rebuildable, provider gurps.SpellListProvider) *PageList {
	p := newPageList(owner, tbl.NewSpellsProvider(provider, true))
	p.installIncrementHandler()
	p.installDecrementHandler()
	p.installIncrementSkillHandler()
	p.installDecrementSkillHandler()
	return p
}

// NewNotesPageList creates the notes page list.
func NewNotesPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	return newPageList(owner, tbl.NewNotesProvider(provider, true))
}

// NewConditionalModifiersPageList creates the conditional modifiers page list.
func NewConditionalModifiersPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, tbl.NewConditionalModifiersProvider(entity))
}

// NewReactionsPageList creates the reaction modifiers page list.
func NewReactionsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, tbl.NewReactionModifiersProvider(entity))
}

// NewMeleeWeaponsPageList creates the melee weapons page list.
func NewMeleeWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, tbl.NewWeaponsProvider(entity, weapon.Melee))
}

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, tbl.NewWeaponsProvider(entity, weapon.Ranged))
}

func newPageList(owner widget.Rebuildable, provider tbl.TableProvider) *PageList {
	p := &PageList{
		table:         unison.NewTable(),
		provider:      provider,
		canPerformMap: make(map[int]func() bool),
		performMap:    make(map[int]func()),
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetBorder(unison.NewLineBorder(theme.HeaderColor, 0, unison.NewUniformInsets(1), false))
	p.table.DividerInk = theme.HeaderColor
	p.table.MinimumRowHeight = theme.PageFieldPrimaryFont.LineHeight()
	p.table.Padding.Top = 0
	p.table.Padding.Bottom = 0
	p.table.HierarchyColumnIndex = provider.HierarchyColumnIndex()
	p.table.HierarchyIndent = theme.PageFieldPrimaryFont.LineHeight()
	p.table.PreventUserColumnResize = true
	headers := provider.Headers()
	p.table.ColumnSizes = make([]unison.ColumnSize, len(headers))
	for i := range p.table.ColumnSizes {
		_, pref, _ := headers[i].AsPanel().Sizes(unison.Size{})
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
	p.table.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		p.table.RequestFocus()
		return p.table.DefaultMouseDown(where, button, clickCount, mod)
	}
	p.table.FrameChangeCallback = func() {
		p.table.SizeColumnsToFitWithExcessIn(p.provider.ExcessWidthColumnIndex())
	}
	p.table.CanPerformCmdCallback = func(_ any, id int) (enabled, handled bool) {
		if f, ok := p.canPerformMap[id]; ok {
			return f(), true
		}
		return false, false
	}
	p.table.PerformCmdCallback = func(_ any, id int) bool {
		if f, ok := p.performMap[id]; ok {
			f()
			return true
		}
		return false
	}
	p.table.SelectionDoubleClickCallback = func() {
		if enabled, _ := p.table.CanPerformCmd(nil, constants.OpenEditorItemID); enabled {
			p.table.PerformCmd(nil, constants.OpenEditorItemID)
		}
	}
	p.tableHeader = unison.NewTableHeader(p.table, headers...)
	p.tableHeader.BackgroundInk = theme.HeaderColor
	p.tableHeader.DividerInk = theme.HeaderColor
	p.tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, unison.Insets{Bottom: 1}, false)
	p.tableHeader.SetBorder(p.tableHeader.HeaderBorder)
	p.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fxp.FromString(s1); err == nil {
			var n2 fxp.Int
			if n2, err = fxp.FromString(s2); err == nil {
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
	p.tableHeader.DrawCallback = func(gc *unison.Canvas, dirty unison.Rect) {
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
	p.table.SetTopLevelRows(p.provider.RowData(p.table))
	p.AddChild(p.tableHeader)
	p.AddChild(p.table)
	if owner != nil {
		p.installPerformHandlers(constants.OpenEditorItemID,
			func() bool { return p.table.HasSelection() },
			func() { p.provider.OpenEditor(owner, p.table) })
	}
	p.installOpenPageReferenceHandlers()
	_, pref, _ := p.tableHeader.Sizes(geom.Size[float32]{})
	p.SetLayoutData(&unison.FlexLayoutData{
		MinSize: geom.NewSize(0, pref.Height*2),
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
		HGrab:   true,
		VGrab:   true,
	})
	return p
}

func (p *PageList) installPerformHandlers(id int, can func() bool, do func()) {
	p.canPerformMap[id] = can
	p.performMap[id] = do
}

func (p *PageList) installOpenPageReferenceHandlers() {
	canOpenPageRefFunc := tbl.NewCanOpenPageRefFunc(p.table)
	p.installPerformHandlers(constants.OpenOnePageReferenceItemID, canOpenPageRefFunc, tbl.NewOpenPageRefFunc(p.table))
	p.installPerformHandlers(constants.OpenEachPageReferenceItemID, canOpenPageRefFunc,
		tbl.NewOpenEachPageRefFunc(p.table))
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
						item.Levels = increment(item.Levels)
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
					if item.IsLeveled() && item.Levels > 0 {
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
						item.Levels = decrement(item.Levels).Max(0)
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

func increment(value fxp.Int) fxp.Int {
	return value.Trunc() + fxp.One
}

func decrement(value fxp.Int) fxp.Int {
	v := value.Trunc()
	if v == value {
		v -= fxp.One
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
					if item.TechLevel, changed = gurps.AdjustTechLevel(item.TechLevel, fxp.One); changed {
						entity = item.Entity
					}
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fxp.One); changed {
							entity = item.Entity
						}
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, fxp.One); changed {
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
					if item.TechLevel, changed = gurps.AdjustTechLevel(item.TechLevel, -fxp.One); changed {
						entity = item.Entity
					}
				case *gurps.Skill:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, -fxp.One); changed {
							entity = item.Entity
						}
					}
				case *gurps.Spell:
					if !item.Container() && item.TechLevel != nil {
						if *item.TechLevel, changed = gurps.AdjustTechLevel(*item.TechLevel, -fxp.One); changed {
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

func (p *PageList) installConvertToContainerHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.ConvertToContainerItemID,
		func() bool { return canConvertToContainer(p.table) },
		func() { convertToContainer(owner, p.table) })
}

// SelectedNodes returns the set of selected nodes. If 'minimal' is true, then children of selected rows that may also
// be selected are not returned, just the topmost row that is selected in any given hierarchy.
func (p *PageList) SelectedNodes(minimal bool) []*tbl.Node {
	if p == nil {
		return nil
	}
	rows := p.table.SelectedRows(minimal)
	selection := make([]*tbl.Node, 0, len(rows))
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			selection = append(selection, n)
		}
	}
	return selection
}

// RecordSelection collects the currently selected row UUIDs.
func (p *PageList) RecordSelection() map[uuid.UUID]bool {
	if p == nil {
		return nil
	}
	rows := p.table.SelectedRows(false)
	selection := make(map[uuid.UUID]bool, len(rows))
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			selection[n.Data().UUID()] = true
		}
	}
	return selection
}

// ApplySelection locates the rows with the given UUIDs and selects them, replacing any existing selection.
func (p *PageList) ApplySelection(selection map[uuid.UUID]bool) {
	p.table.ClearSelection()
	if len(selection) != 0 {
		_, indexes := p.collectRowMappings(0, make([]int, 0, len(selection)), selection, p.table.TopLevelRows())
		if len(indexes) != 0 {
			p.table.SelectByIndex(indexes...)
		}
	}
}

// Sync the underlying data.
func (p *PageList) Sync() {
	p.provider.SyncHeader(p.tableHeader.ColumnHeaders)
	selection := p.RecordSelection()
	p.table.SetTopLevelRows(p.provider.RowData(p.table))
	p.ApplySelection(selection)
	p.table.NeedsLayout = true
	p.NeedsLayout = true
	if parent := p.Parent(); parent != nil {
		parent.NeedsLayout = true
	}
}

func (p *PageList) collectRowMappings(index int, indexes []int, selection map[uuid.UUID]bool, rows []unison.TableRowData) (updatedIndex int, updatedIndexes []int) {
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			if selection[n.Data().UUID()] {
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

// CreateItem calls CreateItem on the contained TableProvider.
func (p *PageList) CreateItem(owner widget.Rebuildable, variant tbl.ItemVariant) {
	p.provider.CreateItem(owner, p.table, variant)
}
