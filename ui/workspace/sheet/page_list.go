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
	"github.com/richardwilkes/gcs/ui/workspace/editors"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

var _ widget.Syncer = &PageList{}

// PageList holds a list for a sheet page.
type PageList struct {
	unison.Panel
	tableHeader   *unison.TableHeader
	table         *unison.Table
	provider      editors.TableProvider
	canPerformMap map[int]func() bool
	performMap    map[int]func()
}

// NewAdvantagesPageList creates the advantages page list.
func NewAdvantagesPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, editors.NewAdvantagesProvider(provider, true))
	p.installToggleDisabledHandler(owner)
	p.installIncrementLevelHandler(owner)
	p.installDecrementLevelHandler(owner)
	return p
}

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, editors.NewEquipmentProvider(provider, true, true))
	p.installToggleEquippedHandler(owner)
	p.installIncrementQuantityHandler(owner)
	p.installDecrementQuantityHandler(owner)
	p.installIncrementUsesHandler(owner)
	p.installDecrementUsesHandler(owner)
	p.installIncrementTechLevelHandler(owner)
	p.installDecrementTechLevelHandler(owner)
	p.installConvertToContainerHandler(owner)
	return p
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, editors.NewEquipmentProvider(provider, true, false))
	p.installIncrementQuantityHandler(owner)
	p.installDecrementQuantityHandler(owner)
	p.installIncrementUsesHandler(owner)
	p.installDecrementUsesHandler(owner)
	p.installIncrementTechLevelHandler(owner)
	p.installDecrementTechLevelHandler(owner)
	p.installConvertToContainerHandler(owner)
	return p
}

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, editors.NewSkillsProvider(provider, true))
	p.installIncrementPointsHandler(owner)
	p.installDecrementPointsHandler(owner)
	p.installIncrementSkillHandler(owner)
	p.installDecrementSkillHandler(owner)
	p.installIncrementTechLevelHandler(owner)
	p.installDecrementTechLevelHandler(owner)
	return p
}

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(owner widget.Rebuildable, provider gurps.SpellListProvider) *PageList {
	p := newPageList(owner, editors.NewSpellsProvider(provider, true))
	p.installIncrementPointsHandler(owner)
	p.installDecrementPointsHandler(owner)
	p.installIncrementSkillHandler(owner)
	p.installDecrementSkillHandler(owner)
	return p
}

// NewNotesPageList creates the notes page list.
func NewNotesPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	return newPageList(owner, editors.NewNotesProvider(provider, true))
}

// NewConditionalModifiersPageList creates the conditional modifiers page list.
func NewConditionalModifiersPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, editors.NewConditionalModifiersProvider(entity))
}

// NewReactionsPageList creates the reaction modifiers page list.
func NewReactionsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, editors.NewReactionModifiersProvider(entity))
}

// NewMeleeWeaponsPageList creates the melee weapons page list.
func NewMeleeWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, editors.NewWeaponsProvider(entity, weapon.Melee))
}

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, editors.NewWeaponsProvider(entity, weapon.Ranged))
}

func newPageList(owner widget.Rebuildable, provider editors.TableProvider) *PageList {
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
	widget.TableSetupColumnSizes(p.table, headers)
	p.table.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	widget.TableInstallStdCallbacks(p.table)
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
	p.tableHeader = widget.TableCreateHeader(p.table, headers)
	p.tableHeader.BackgroundInk = theme.HeaderColor
	p.tableHeader.DividerInk = theme.HeaderColor
	p.tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, unison.Insets{Bottom: 1}, false)
	p.tableHeader.SetBorder(p.tableHeader.HeaderBorder)
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
	canOpenPageRefFunc := editors.NewCanOpenPageRefFunc(p.table)
	p.installPerformHandlers(constants.OpenOnePageReferenceItemID, canOpenPageRefFunc, editors.NewOpenPageRefFunc(p.table))
	p.installPerformHandlers(constants.OpenEachPageReferenceItemID, canOpenPageRefFunc,
		editors.NewOpenEachPageRefFunc(p.table))
}

func (p *PageList) installToggleDisabledHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.ToggleStateItemID,
		func() bool { return canToggleDisabled(p.table) },
		func() { toggleDisabled(owner, p.table) })
}

func (p *PageList) installToggleEquippedHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.ToggleStateItemID,
		func() bool { return canToggleEquipped(p.table) },
		func() { toggleEquipped(owner, p.table) })
}

func (p *PageList) installIncrementPointsHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementItemID,
		func() bool { return canAdjustRawPoints(p.table, true) },
		func() { adjustRawPoints(owner, p.table, true) })
}

func (p *PageList) installDecrementPointsHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementItemID,
		func() bool { return canAdjustRawPoints(p.table, false) },
		func() { adjustRawPoints(owner, p.table, false) })
}

func (p *PageList) installIncrementLevelHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementItemID,
		func() bool { return canAdjustAdvantageLevel(p.table, true) },
		func() { adjustAdvantageLevel(owner, p.table, true) })
}

func (p *PageList) installDecrementLevelHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementItemID,
		func() bool { return canAdjustAdvantageLevel(p.table, false) },
		func() { adjustAdvantageLevel(owner, p.table, false) })
}

func (p *PageList) installIncrementQuantityHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementItemID,
		func() bool { return canAdjustQuantity(p.table, true) },
		func() { adjustQuantity(owner, p.table, true) })
}

func (p *PageList) installDecrementQuantityHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementItemID,
		func() bool { return canAdjustQuantity(p.table, false) },
		func() { adjustQuantity(owner, p.table, false) })
}

func (p *PageList) installIncrementUsesHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementUsesItemID,
		func() bool { return canAdjustUses(p.table, 1) },
		func() { adjustUses(owner, p.table, 1) })
}

func (p *PageList) installDecrementUsesHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementUsesItemID,
		func() bool { return canAdjustUses(p.table, -1) },
		func() { adjustUses(owner, p.table, -1) })
}

func (p *PageList) installIncrementSkillHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementSkillLevelItemID,
		func() bool { return canAdjustSkillLevel(p.table, true) },
		func() { adjustSkillLevel(owner, p.table, true) })
}

func (p *PageList) installDecrementSkillHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementSkillLevelItemID,
		func() bool { return canAdjustSkillLevel(p.table, false) },
		func() { adjustSkillLevel(owner, p.table, false) })
}

func (p *PageList) installIncrementTechLevelHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.IncrementTechLevelItemID,
		func() bool { return canAdjustTechLevel(p.table, fxp.One) },
		func() { adjustTechLevel(owner, p.table, fxp.One) })
}

func (p *PageList) installDecrementTechLevelHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.DecrementTechLevelItemID,
		func() bool { return canAdjustTechLevel(p.table, -fxp.One) },
		func() { adjustTechLevel(owner, p.table, -fxp.One) })
}

func (p *PageList) installConvertToContainerHandler(owner widget.Rebuildable) {
	p.installPerformHandlers(constants.ConvertToContainerItemID,
		func() bool { return canConvertToContainer(p.table) },
		func() { convertToContainer(owner, p.table) })
}

// SelectedNodes returns the set of selected nodes. If 'minimal' is true, then children of selected rows that may also
// be selected are not returned, just the topmost row that is selected in any given hierarchy.
func (p *PageList) SelectedNodes(minimal bool) []*editors.Node {
	if p == nil {
		return nil
	}
	rows := p.table.SelectedRows(minimal)
	selection := make([]*editors.Node, 0, len(rows))
	for _, row := range rows {
		if n, ok := row.(*editors.Node); ok {
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
		if n, ok := row.(*editors.Node); ok {
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
		if n, ok := row.(*editors.Node); ok {
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
func (p *PageList) CreateItem(owner widget.Rebuildable, variant editors.ItemVariant) {
	p.provider.CreateItem(owner, p.table, variant)
}
