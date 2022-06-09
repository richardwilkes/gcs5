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
	tableHeader *unison.TableHeader
	table       *unison.Table
	provider    editors.TableProvider
}

// NewTraitsPageList creates the traits page list.
func NewTraitsPageList(owner widget.Rebuildable, provider gurps.ListProvider) *PageList {
	p := newPageList(owner, editors.NewTraitsProvider(provider, true))
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
	return newPageList(nil, editors.NewWeaponsProvider(entity, weapon.Melee, true))
}

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return newPageList(nil, editors.NewWeaponsProvider(entity, weapon.Ranged, true))
}

func newPageList(owner widget.Rebuildable, provider editors.TableProvider) *PageList {
	p := &PageList{
		table:    unison.NewTable(),
		provider: provider,
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
	singular, plural := p.provider.ItemNames()
	p.table.InstallDragSupport(p.provider.DragSVG(), p.provider.DragKey(), singular, plural)
	if owner != nil {
		p.table.InstallDropSupport(p.provider.DragKey(), widget.StdDropCallback)
		p.InstallCmdHandlers(constants.OpenEditorItemID,
			func(_ any) bool { return p.table.HasSelection() },
			func(_ any) { p.provider.OpenEditor(owner, p.table) })
		p.InstallCmdHandlers(unison.DeleteItemID,
			func(_ any) bool { return p.table.HasSelection() },
			func(_ any) { p.provider.DeleteSelection(p.table) })
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

func (p *PageList) installOpenPageReferenceHandlers() {
	p.InstallCmdHandlers(constants.OpenOnePageReferenceItemID,
		func(_ any) bool { return editors.CanOpenPageRef(p.table) },
		func(_ any) { editors.OpenPageRef(p.table) })
	p.InstallCmdHandlers(constants.OpenEachPageReferenceItemID,
		func(_ any) bool { return editors.CanOpenPageRef(p.table) },
		func(_ any) { editors.OpenEachPageRef(p.table) })
}

func (p *PageList) installToggleDisabledHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.ToggleStateItemID,
		func(_ any) bool { return canToggleDisabled(p.table) },
		func(_ any) { toggleDisabled(owner, p.table) })
}

func (p *PageList) installToggleEquippedHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.ToggleStateItemID,
		func(_ any) bool { return canToggleEquipped(p.table) },
		func(_ any) { toggleEquipped(owner, p.table) })
}

func (p *PageList) installIncrementPointsHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementItemID,
		func(_ any) bool { return canAdjustRawPoints(p.table, true) },
		func(_ any) { adjustRawPoints(owner, p.table, true) })
}

func (p *PageList) installDecrementPointsHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementItemID,
		func(_ any) bool { return canAdjustRawPoints(p.table, false) },
		func(_ any) { adjustRawPoints(owner, p.table, false) })
}

func (p *PageList) installIncrementLevelHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementItemID,
		func(_ any) bool { return canAdjustTraitLevel(p.table, true) },
		func(_ any) { adjustTraitLevel(owner, p.table, true) })
}

func (p *PageList) installDecrementLevelHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementItemID,
		func(_ any) bool { return canAdjustTraitLevel(p.table, false) },
		func(_ any) { adjustTraitLevel(owner, p.table, false) })
}

func (p *PageList) installIncrementQuantityHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementItemID,
		func(_ any) bool { return canAdjustQuantity(p.table, true) },
		func(_ any) { adjustQuantity(owner, p.table, true) })
}

func (p *PageList) installDecrementQuantityHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementItemID,
		func(_ any) bool { return canAdjustQuantity(p.table, false) },
		func(_ any) { adjustQuantity(owner, p.table, false) })
}

func (p *PageList) installIncrementUsesHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementUsesItemID,
		func(_ any) bool { return canAdjustUses(p.table, 1) },
		func(_ any) { adjustUses(owner, p.table, 1) })
}

func (p *PageList) installDecrementUsesHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementUsesItemID,
		func(_ any) bool { return canAdjustUses(p.table, -1) },
		func(_ any) { adjustUses(owner, p.table, -1) })
}

func (p *PageList) installIncrementSkillHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementSkillLevelItemID,
		func(_ any) bool { return canAdjustSkillLevel(p.table, true) },
		func(_ any) { adjustSkillLevel(owner, p.table, true) })
}

func (p *PageList) installDecrementSkillHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementSkillLevelItemID,
		func(_ any) bool { return canAdjustSkillLevel(p.table, false) },
		func(_ any) { adjustSkillLevel(owner, p.table, false) })
}

func (p *PageList) installIncrementTechLevelHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.IncrementTechLevelItemID,
		func(_ any) bool { return canAdjustTechLevel(p.table, fxp.One) },
		func(_ any) { adjustTechLevel(owner, p.table, fxp.One) })
}

func (p *PageList) installDecrementTechLevelHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.DecrementTechLevelItemID,
		func(_ any) bool { return canAdjustTechLevel(p.table, -fxp.One) },
		func(_ any) { adjustTechLevel(owner, p.table, -fxp.One) })
}

func (p *PageList) installConvertToContainerHandler(owner widget.Rebuildable) {
	p.InstallCmdHandlers(constants.ConvertToContainerItemID,
		func(_ any) bool { return canConvertToContainer(p.table) },
		func(_ any) { convertToContainer(owner, p.table) })
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
	return editors.RecordTableSelection(p.table)
}

// ApplySelection locates the rows with the given UUIDs and selects them, replacing any existing selection.
func (p *PageList) ApplySelection(selection map[uuid.UUID]bool) {
	if p != nil {
		editors.ApplyTableSelection(p.table, selection)
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

// CreateItem calls CreateItem on the contained TableProvider.
func (p *PageList) CreateItem(owner widget.Rebuildable, variant editors.ItemVariant) {
	p.provider.CreateItem(owner, p.table, variant)
}
