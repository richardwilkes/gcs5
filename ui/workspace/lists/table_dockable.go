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

package lists

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/crc"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/unison"
)

var (
	_ workspace.FileBackedDockable = &TableDockable{}
	_ unison.UndoManagerProvider   = &TableDockable{}
	_ widget.ModifiableRoot        = &TableDockable{}
	_ widget.Rebuildable           = &TableDockable{}
	_ widget.DockableKind          = &TableDockable{}
	_ unison.TabCloser             = &TableDockable{}
)

// TableDockable holds the view for a file that contains a (potentially hierarchical) list of data.
type TableDockable struct {
	unison.Panel
	path              string
	extension         string
	undoMgr           *unison.UndoManager
	provider          tbl.TableProvider
	saver             func(path string) error
	canCreateIDs      map[int]bool
	hierarchyButton   *unison.Button
	sizeToFitButton   *unison.Button
	scale             int
	scaleField        *widget.PercentageField
	backButton        *unison.Button
	forwardButton     *unison.Button
	searchField       *unison.Field
	matchesLabel      *unison.Label
	scroll            *unison.ScrollPanel
	tableHeader       *unison.TableHeader
	table             *unison.Table
	crc               uint64
	searchResult      []unison.TableRowData
	searchIndex       int
	needsSaveAsPrompt bool
}

type advantageListProvider struct {
	advantages []*gurps.Advantage
}

func (p *advantageListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *advantageListProvider) AdvantageList() []*gurps.Advantage {
	return p.advantages
}

func (p *advantageListProvider) SetAdvantageList(list []*gurps.Advantage) {
	p.advantages = list
}

// NewAdvantageTableDockableFromFile loads a list of advantages from a file and creates a new unison.Dockable for them.
func NewAdvantageTableDockableFromFile(filePath string) (unison.Dockable, error) {
	advantages, err := gurps.NewAdvantagesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewAdvantageTableDockable(filePath, advantages)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewAdvantageTableDockable creates a new unison.Dockable for advantage list files.
func NewAdvantageTableDockable(filePath string, advantages []*gurps.Advantage) *TableDockable {
	provider := &advantageListProvider{advantages: advantages}
	return NewTableDockable(filePath, library.AdvantagesExt, tbl.NewAdvantagesProvider(provider, false),
		func(path string) error { return gurps.SaveAdvantages(provider.AdvantageList(), path) },
		constants.NewAdvantageItemID, constants.NewAdvantageContainerItemID)
}

type advantageModifierListProvider struct {
	modifiers []*gurps.AdvantageModifier
}

func (p *advantageModifierListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *advantageModifierListProvider) AdvantageModifierList() []*gurps.AdvantageModifier {
	return p.modifiers
}

func (p *advantageModifierListProvider) SetAdvantageModifierList(list []*gurps.AdvantageModifier) {
	p.modifiers = list
}

// NewAdvantageModifierTableDockableFromFile loads a list of advantage modifiers from a file and creates a new
// unison.Dockable for them.
func NewAdvantageModifierTableDockableFromFile(filePath string) (unison.Dockable, error) {
	modifiers, err := gurps.NewAdvantageModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewAdvantageModifierTableDockable(filePath, modifiers)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewAdvantageModifierTableDockable creates a new unison.Dockable for advantage modifier list files.
func NewAdvantageModifierTableDockable(filePath string, modifiers []*gurps.AdvantageModifier) *TableDockable {
	provider := &advantageModifierListProvider{modifiers: modifiers}
	return NewTableDockable(filePath, library.AdvantageModifiersExt,
		tbl.NewAdvantageModifiersProvider(provider),
		func(path string) error { return gurps.SaveAdvantageModifiers(provider.AdvantageModifierList(), path) },
		constants.NewAdvantageModifierItemID, constants.NewAdvantageContainerModifierItemID)
}

type equipmentListProvider struct {
	carried []*gurps.Equipment
	other   []*gurps.Equipment
}

func (p *equipmentListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *equipmentListProvider) CarriedEquipmentList() []*gurps.Equipment {
	return p.carried
}

func (p *equipmentListProvider) SetCarriedEquipmentList(list []*gurps.Equipment) {
	p.carried = list
}

func (p *equipmentListProvider) OtherEquipmentList() []*gurps.Equipment {
	return p.other
}

func (p *equipmentListProvider) SetOtherEquipmentList(list []*gurps.Equipment) {
	p.other = list
}

// NewEquipmentTableDockableFromFile loads a list of equipment from a file and creates a new unison.Dockable for them.
func NewEquipmentTableDockableFromFile(filePath string) (unison.Dockable, error) {
	equipment, err := gurps.NewEquipmentFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewEquipmentTableDockable(filePath, equipment)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewEquipmentTableDockable creates a new unison.Dockable for equipment list files.
func NewEquipmentTableDockable(filePath string, equipment []*gurps.Equipment) *TableDockable {
	provider := &equipmentListProvider{other: equipment}
	return NewTableDockable(filePath, library.EquipmentExt, tbl.NewEquipmentProvider(provider, false, false),
		func(path string) error { return gurps.SaveEquipment(provider.OtherEquipmentList(), path) },
		constants.NewCarriedEquipmentItemID, constants.NewCarriedEquipmentContainerItemID)
}

type equipmentModifierListProvider struct {
	modifiers []*gurps.EquipmentModifier
}

func (p *equipmentModifierListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *equipmentModifierListProvider) EquipmentModifierList() []*gurps.EquipmentModifier {
	return p.modifiers
}

func (p *equipmentModifierListProvider) SetEquipmentModifierList(list []*gurps.EquipmentModifier) {
	p.modifiers = list
}

// NewEquipmentModifierTableDockableFromFile loads a list of equipment modifiers from a file and creates a new
// unison.Dockable for them.
func NewEquipmentModifierTableDockableFromFile(filePath string) (unison.Dockable, error) {
	modifiers, err := gurps.NewEquipmentModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewEquipmentModifierTableDockable(filePath, modifiers)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewEquipmentModifierTableDockable creates a new unison.Dockable for equipment modifier list files.
func NewEquipmentModifierTableDockable(filePath string, modifiers []*gurps.EquipmentModifier) *TableDockable {
	provider := &equipmentModifierListProvider{modifiers: modifiers}
	return NewTableDockable(filePath, library.EquipmentModifiersExt,
		tbl.NewEquipmentModifiersProvider(provider),
		func(path string) error { return gurps.SaveEquipmentModifiers(provider.EquipmentModifierList(), path) },
		constants.NewEquipmentModifierItemID, constants.NewEquipmentContainerModifierItemID)
}

type skillListProvider struct {
	skills []*gurps.Skill
}

func (p *skillListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *skillListProvider) SkillList() []*gurps.Skill {
	return p.skills
}

func (p *skillListProvider) SetSkillList(list []*gurps.Skill) {
	p.skills = list
}

// NewSkillTableDockableFromFile loads a list of skills from a file and creates a new unison.Dockable for them.
func NewSkillTableDockableFromFile(filePath string) (unison.Dockable, error) {
	skills, err := gurps.NewSkillsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewSkillTableDockable(filePath, skills)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewSkillTableDockable creates a new unison.Dockable for skill list files.
func NewSkillTableDockable(filePath string, skills []*gurps.Skill) *TableDockable {
	provider := &skillListProvider{skills: skills}
	return NewTableDockable(filePath, library.SkillsExt, tbl.NewSkillsProvider(provider, false),
		func(path string) error { return gurps.SaveSkills(provider.SkillList(), path) },
		constants.NewSkillItemID, constants.NewSkillContainerItemID, constants.NewTechniqueItemID)
}

type spellListProvider struct {
	spells []*gurps.Spell
}

func (p *spellListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *spellListProvider) SpellList() []*gurps.Spell {
	return p.spells
}

func (p *spellListProvider) SetSpellList(list []*gurps.Spell) {
	p.spells = list
}

// NewSpellTableDockableFromFile loads a list of spells from a file and creates a new unison.Dockable for them.
func NewSpellTableDockableFromFile(filePath string) (unison.Dockable, error) {
	spells, err := gurps.NewSpellsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewSpellTableDockable(filePath, spells)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewSpellTableDockable creates a new unison.Dockable for spell list files.
func NewSpellTableDockable(filePath string, spells []*gurps.Spell) *TableDockable {
	provider := &spellListProvider{spells: spells}
	return NewTableDockable(filePath, library.SpellsExt, tbl.NewSpellsProvider(provider, false),
		func(path string) error { return gurps.SaveSpells(provider.SpellList(), path) },
		constants.NewSpellItemID, constants.NewSpellContainerItemID, constants.NewRitualMagicSpellItemID)
}

type noteListProvider struct {
	notes []*gurps.Note
}

func (p *noteListProvider) Entity() *gurps.Entity {
	return nil
}

func (p *noteListProvider) NoteList() []*gurps.Note {
	return p.notes
}

func (p *noteListProvider) SetNoteList(list []*gurps.Note) {
	p.notes = list
}

// NewNoteTableDockableFromFile loads a list of notes from a file and creates a new unison.Dockable for them.
func NewNoteTableDockableFromFile(filePath string) (unison.Dockable, error) {
	notes, err := gurps.NewNotesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewNoteTableDockable(filePath, notes)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewNoteTableDockable creates a new unison.Dockable for note list files.
func NewNoteTableDockable(filePath string, notes []*gurps.Note) *TableDockable {
	provider := &noteListProvider{notes: notes}
	return NewTableDockable(filePath, library.NotesExt, tbl.NewNotesProvider(provider, false),
		func(path string) error { return gurps.SaveNotes(provider.NoteList(), path) },
		constants.NewNoteItemID, constants.NewNoteContainerItemID)
}

// NewTableDockable creates a new TableDockable for list data files.
func NewTableDockable(filePath, extension string, provider tbl.TableProvider, saver func(path string) error, canCreateIDs ...int) *TableDockable {
	d := &TableDockable{
		path:              filePath,
		extension:         extension,
		undoMgr:           unison.NewUndoManager(200, func(err error) { jot.Error(err) }),
		provider:          provider,
		saver:             saver,
		canCreateIDs:      make(map[int]bool),
		scroll:            unison.NewScrollPanel(),
		table:             unison.NewTable(),
		scale:             settings.Global().General.InitialListUIScale,
		needsSaveAsPrompt: true,
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{Columns: 1})

	for _, id := range canCreateIDs {
		d.canCreateIDs[id] = true
	}

	headers := provider.Headers()
	d.table.ColumnSizes = make([]unison.ColumnSize, len(headers))
	for i := range d.table.ColumnSizes {
		_, pref, _ := headers[i].AsPanel().Sizes(unison.Size{})
		d.table.ColumnSizes[i].AutoMinimum = pref.Width
		d.table.ColumnSizes[i].AutoMaximum = 800
		d.table.ColumnSizes[i].Minimum = pref.Width
		d.table.ColumnSizes[i].Maximum = 10000
	}
	d.table.HierarchyColumnIndex = provider.HierarchyColumnIndex()
	d.table.SetTopLevelRows(provider.RowData(d.table))
	d.table.SizeColumnsToFit(true)
	mouseDownCallback := d.table.MouseDownCallback
	d.table.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		d.table.RequestFocus()
		return mouseDownCallback(where, button, clickCount, mod)
	}
	d.table.SelectionDoubleClickCallback = func() {
		if enabled, _ := d.table.CanPerformCmd(nil, constants.OpenEditorItemID); enabled {
			d.table.PerformCmd(nil, constants.OpenEditorItemID)
		}
	}

	d.tableHeader = unison.NewTableHeader(d.table, headers...)
	d.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fxp.FromString(s1); err == nil {
			var n2 fxp.Int
			if n2, err = fxp.FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	d.scroll.SetColumnHeader(d.tableHeader)
	d.scroll.SetContent(d.table, unison.FillBehavior, unison.FillBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	d.hierarchyButton = unison.NewSVGButton(res.HierarchySVG)
	d.hierarchyButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Opens/closes all hierarchical rows"))
	d.hierarchyButton.ClickCallback = d.toggleHierarchy

	d.sizeToFitButton = unison.NewSVGButton(res.SizeToFitSVG)
	d.sizeToFitButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sets the width of each column to fit its contents"))
	d.sizeToFitButton.ClickCallback = d.sizeToFit

	scaleTitle := i18n.Text("Scale")
	d.scaleField = widget.NewPercentageField(scaleTitle, func() int { return d.scale }, func(v int) {
		d.scale = v
		d.applyScale()
	}, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax, false, false)
	d.scaleField.Tooltip = unison.NewTooltipWithText(scaleTitle)

	d.backButton = unison.NewSVGButton(res.BackSVG)
	d.backButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Previous Match"))
	d.backButton.ClickCallback = d.previousMatch
	d.backButton.SetEnabled(false)

	d.forwardButton = unison.NewSVGButton(res.ForwardSVG)
	d.forwardButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Next Match"))
	d.forwardButton.ClickCallback = d.nextMatch
	d.forwardButton.SetEnabled(false)

	d.searchField = unison.NewField()
	search := i18n.Text("Search")
	d.searchField.Watermark = search
	d.searchField.Tooltip = unison.NewTooltipWithText(search)
	d.searchField.ModifiedCallback = d.searchModified
	d.searchField.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
		if keyCode == unison.KeyReturn || keyCode == unison.KeyNumPadEnter {
			if mod.ShiftDown() {
				d.previousMatch()
			} else {
				d.nextMatch()
			}
			return true
		}
		return d.searchField.DefaultKeyDown(keyCode, mod, repeat)
	}
	d.searchField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})

	d.matchesLabel = unison.NewLabel()
	d.matchesLabel.Text = "-"
	d.matchesLabel.Tooltip = unison.NewTooltipWithText(i18n.Text("Number of matches found"))

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1},
		false), unison.NewEmptyBorder(unison.StdInsets())))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(d.hierarchyButton)
	toolbar.AddChild(d.sizeToFitButton)
	toolbar.AddChild(d.scaleField)
	toolbar.AddChild(d.backButton)
	toolbar.AddChild(d.forwardButton)
	toolbar.AddChild(d.searchField)
	toolbar.AddChild(d.matchesLabel)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	d.AddChild(toolbar)
	d.AddChild(d.scroll)

	d.applyScale()

	d.CanPerformCmdCallback = d.canPerformCmd
	d.PerformCmdCallback = d.performCmd

	d.crc = d.crc64()
	return d
}

// UndoManager implements undo.Provider
func (d *TableDockable) UndoManager() *unison.UndoManager {
	return d.undoMgr
}

// DockableKind implements widget.DockableKind
func (d *TableDockable) DockableKind() string {
	return widget.ListDockableKind
}

func (d *TableDockable) applyScale() {
	s := float32(d.scale) / 100
	d.tableHeader.SetScale(s)
	d.table.SetScale(s)
	d.scroll.Sync()
}

// TitleIcon implements workspace.FileBackedDockable
func (d *TableDockable) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (d *TableDockable) Title() string {
	return fs.BaseName(d.path)
}

func (d *TableDockable) String() string {
	return d.Title()
}

// Tooltip implements workspace.FileBackedDockable
func (d *TableDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (d *TableDockable) BackingFilePath() string {
	return d.path
}

// Modified implements workspace.FileBackedDockable
func (d *TableDockable) Modified() bool {
	return d.crc != d.crc64()
}

// MarkModified implements widget.ModifiableRoot.
func (d *TableDockable) MarkModified() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.UpdateTitle(d)
	}
}

// MayAttemptClose implements unison.TabCloser
func (d *TableDockable) MayAttemptClose() bool {
	return workspace.MayAttemptCloseOfGroup(d)
}

// AttemptClose implements unison.TabCloser
func (d *TableDockable) AttemptClose() bool {
	if !workspace.CloseGroup(d) {
		return false
	}
	if d.Modified() {
		switch unison.YesNoCancelDialog(fmt.Sprintf(i18n.Text("Save changes made to\n%s?"), d.Title()), "") {
		case unison.ModalResponseDiscard:
		case unison.ModalResponseOK:
			if !d.save(false) {
				return false
			}
		case unison.ModalResponseCancel:
			return false
		}
	}
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
	return true
}

func (d *TableDockable) save(forceSaveAs bool) bool {
	success := false
	if forceSaveAs || d.needsSaveAsPrompt {
		success = workspace.SaveDockableAs(d, d.extension, d.saver, func(path string) {
			d.crc = d.crc64()
			d.path = path
		})
	} else {
		success = workspace.SaveDockable(d, d.saver, func() { d.crc = d.crc64() })
	}
	if success {
		d.needsSaveAsPrompt = false
	}
	return success
}

func (d *TableDockable) toggleHierarchy() {
	first := true
	open := false
	for _, row := range d.table.TopLevelRows() {
		if row.CanHaveChildRows() {
			if first {
				first = false
				open = !row.IsOpen()
			}
			setRowOpen(row, open)
		}
	}
	d.table.SyncToModel()
}

func setRowOpen(row unison.TableRowData, open bool) {
	row.SetOpen(open)
	for _, child := range row.ChildRows() {
		if child.CanHaveChildRows() {
			setRowOpen(child, open)
		}
	}
}

func (d *TableDockable) sizeToFit() {
	d.table.SizeColumnsToFit(true)
	d.table.MarkForRedraw()
}

func (d *TableDockable) searchModified() {
	d.searchIndex = 0
	d.searchResult = nil
	text := strings.ToLower(d.searchField.Text())
	for _, row := range d.table.TopLevelRows() {
		d.search(text, row)
	}
	d.adjustForMatch()
}

func (d *TableDockable) search(text string, row unison.TableRowData) {
	if matcher, ok := row.(tbl.Matcher); ok {
		if matcher.Match(text) {
			d.searchResult = append(d.searchResult, row)
		}
	}
	if row.CanHaveChildRows() {
		for _, child := range row.ChildRows() {
			d.search(text, child)
		}
	}
}

func (d *TableDockable) previousMatch() {
	if d.searchIndex > 0 {
		d.searchIndex--
		d.adjustForMatch()
	}
}

func (d *TableDockable) nextMatch() {
	if d.searchIndex < len(d.searchResult)-1 {
		d.searchIndex++
		d.adjustForMatch()
	}
}

func (d *TableDockable) adjustForMatch() {
	d.backButton.SetEnabled(d.searchIndex != 0)
	d.forwardButton.SetEnabled(len(d.searchResult) != 0 && d.searchIndex != len(d.searchResult)-1)
	if len(d.searchResult) != 0 {
		d.matchesLabel.Text = fmt.Sprintf(i18n.Text("%d of %d"), d.searchIndex+1, len(d.searchResult))
		row := d.searchResult[d.searchIndex]
		d.table.DiscloseRow(row, false)
		d.table.ClearSelection()
		rowIndex := d.table.RowToIndex(row)
		d.table.SelectByIndex(rowIndex)
		d.table.ScrollRowIntoView(rowIndex)
	} else {
		d.matchesLabel.Text = "-"
	}
	d.matchesLabel.Parent().MarkForLayoutAndRedraw()
}

// Rebuild implements widget.Rebuildable.
func (d *TableDockable) Rebuild(_ bool) {
	d.table.SyncToModel()
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.UpdateTitle(d)
	}
}

func (d *TableDockable) canPerformCmd(_ any, id int) (enabled, handled bool) {
	switch id {
	case constants.OpenEditorItemID:
		return d.table.HasSelection(), true
	case constants.OpenOnePageReferenceItemID,
		constants.OpenEachPageReferenceItemID:
		return tbl.CanOpenPageRef(d.table), true
	case constants.SaveItemID:
		return d.Modified(), true
	case constants.SaveAsItemID:
		return true, true
	default:
		if d.canCreateIDs[id] {
			return true, true
		}
		return false, false
	}
}

func (d *TableDockable) performCmd(_ any, id int) bool {
	switch id {
	case constants.OpenEditorItemID:
		d.provider.OpenEditor(d, d.table)
	case constants.OpenOnePageReferenceItemID:
		tbl.OpenPageRef(d.table)
	case constants.OpenEachPageReferenceItemID:
		tbl.OpenEachPageRef(d.table)
	case constants.SaveItemID:
		d.save(false)
	case constants.SaveAsItemID:
		d.save(true)
	case constants.NewAdvantageItemID,
		constants.NewAdvantageModifierItemID,
		constants.NewSkillItemID,
		constants.NewSpellItemID,
		constants.NewCarriedEquipmentItemID,
		constants.NewOtherEquipmentItemID,
		constants.NewEquipmentModifierItemID,
		constants.NewNoteItemID:
		d.provider.CreateItem(d, d.table, tbl.NoItemVariant)
	case constants.NewAdvantageContainerItemID,
		constants.NewAdvantageContainerModifierItemID,
		constants.NewSkillContainerItemID,
		constants.NewSpellContainerItemID,
		constants.NewCarriedEquipmentContainerItemID,
		constants.NewOtherEquipmentContainerItemID,
		constants.NewEquipmentContainerModifierItemID,
		constants.NewNoteContainerItemID:
		d.provider.CreateItem(d, d.table, tbl.ContainerItemVariant)
	case constants.NewTechniqueItemID,
		constants.NewRitualMagicSpellItemID:
		d.provider.CreateItem(d, d.table, tbl.AlternateItemVariant)
	default:
		return false
	}
	return true
}

func (d *TableDockable) crc64() uint64 {
	var buffer bytes.Buffer
	rows := d.provider.RowData(d.table)
	data := make([]any, 0, len(rows))
	for _, row := range rows {
		if n, ok := row.(*tbl.Node); ok {
			data = append(data, n.Data())
		}
	}
	if err := jio.Save(context.Background(), &buffer, data); err != nil {
		return 0
	}
	return crc.Bytes(0, buffer.Bytes())
}
