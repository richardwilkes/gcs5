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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

// TableDockable holds the view for a file that contains a (potentially hierarchical) list of data.
type TableDockable struct {
	unison.Panel
	path            string
	lockButton      *unison.Button
	hierarchyButton *unison.Button
	sizeToFitButton *unison.Button
	scale           int
	scaleField      *widget.PercentageField
	backButton      *unison.Button
	forwardButton   *unison.Button
	searchField     *unison.Field
	matchesLabel    *unison.Label
	scroll          *unison.ScrollPanel
	tableHeader     *unison.TableHeader
	table           *unison.Table
	searchResult    []unison.TableRowData
	searchIndex     int
	locked          bool
}

// NewAdvantageTableDockableFromFile loads a list of advantages from a file and creates a new unison.Dockable for them.
func NewAdvantageTableDockableFromFile(filePath string) (unison.Dockable, error) {
	advantages, err := gurps.NewAdvantagesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewAdvantageTableDockable(filePath, advantages), nil
}

// NewAdvantageTableDockable creates a new unison.Dockable for advantage list files.
func NewAdvantageTableDockable(filePath string, advantages []*gurps.Advantage) unison.Dockable {
	return NewTableDockable(filePath, tbl.NewAdvantageTableHeaders(false),
		tbl.NewAdvantageRowData(func() []*gurps.Advantage { return advantages }, false))
}

// NewAdvantageModifierTableDockableFromFile loads a list of advantage modifiers from a file and creates a new
// unison.Dockable for them.
func NewAdvantageModifierTableDockableFromFile(filePath string) (unison.Dockable, error) {
	modifiers, err := gurps.NewAdvantageModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewAdvantageModifierTableDockable(filePath, modifiers), nil
}

// NewAdvantageModifierTableDockable creates a new unison.Dockable for advantage modifier list files.
func NewAdvantageModifierTableDockable(filePath string, modifiers []*gurps.AdvantageModifier) unison.Dockable {
	return NewTableDockable(filePath, tbl.NewAdvantageModifierTableHeaders(),
		tbl.NewAdvantageModifierRowData(func() []*gurps.AdvantageModifier { return modifiers }))
}

// NewEquipmentTableDockableFromFile loads a list of equipment from a file and creates a new unison.Dockable for them.
func NewEquipmentTableDockableFromFile(filePath string) (unison.Dockable, error) {
	equipment, err := gurps.NewEquipmentFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewEquipmentTableDockable(filePath, equipment), nil
}

// NewEquipmentTableDockable creates a new unison.Dockable for equipment list files.
func NewEquipmentTableDockable(filePath string, equipment []*gurps.Equipment) unison.Dockable {
	return NewTableDockable(filePath, tbl.NewEquipmentTableHeaders(nil, false, false),
		tbl.NewEquipmentRowData(func() []*gurps.Equipment { return equipment }, false, false))
}

// NewEquipmentModifierTableDockableFromFile loads a list of equipment modifiers from a file and creates a new
// unison.Dockable for them.
func NewEquipmentModifierTableDockableFromFile(filePath string) (unison.Dockable, error) {
	modifiers, err := gurps.NewEquipmentModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewEquipmentModifierTableDockable(filePath, modifiers), nil
}

// NewEquipmentModifierTableDockable creates a new unison.Dockable for equipment modifier list files.
func NewEquipmentModifierTableDockable(filePath string, modifiers []*gurps.EquipmentModifier) unison.Dockable {
	return NewTableDockable(filePath, tbl.NewEquipmentModifierTableHeaders(),
		tbl.NewEquipmentModifierRowData(func() []*gurps.EquipmentModifier { return modifiers }))
}

type skillListProvider struct {
	skills []*gurps.Skill
}

func (p *skillListProvider) SkillList() []*gurps.Skill {
	return p.skills
}

// NewSkillTableDockableFromFile loads a list of skills from a file and creates a new unison.Dockable for them.
func NewSkillTableDockableFromFile(filePath string) (unison.Dockable, error) {
	skills, err := gurps.NewSkillsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewSkillTableDockable(filePath, skills), nil
}

// NewSkillTableDockable creates a new unison.Dockable for skill list files.
func NewSkillTableDockable(filePath string, skills []*gurps.Skill) unison.Dockable {
	p := &skillListProvider{skills: skills}
	return NewTableDockable(filePath, tbl.NewSkillTableHeaders(p, false), tbl.NewSkillRowData(p, false))
}

type spellListProvider struct {
	spells []*gurps.Spell
}

func (p *spellListProvider) SpellList() []*gurps.Spell {
	return p.spells
}

// NewSpellTableDockableFromFile loads a list of spells from a file and creates a new unison.Dockable for them.
func NewSpellTableDockableFromFile(filePath string) (unison.Dockable, error) {
	spells, err := gurps.NewSpellsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewSpellTableDockable(filePath, spells), nil
}

// NewSpellTableDockable creates a new unison.Dockable for spell list files.
func NewSpellTableDockable(filePath string, spells []*gurps.Spell) unison.Dockable {
	p := &spellListProvider{spells: spells}
	return NewTableDockable(filePath, tbl.NewSpellTableHeaders(p, false), tbl.NewSpellRowData(p, false))
}

// NewNoteTableDockableFromFile loads a list of notes from a file and creates a new unison.Dockable for them.
func NewNoteTableDockableFromFile(filePath string) (unison.Dockable, error) {
	notes, err := gurps.NewNotesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewNoteTableDockable(filePath, notes), nil
}

// NewNoteTableDockable creates a new unison.Dockable for note list files.
func NewNoteTableDockable(filePath string, notes []*gurps.Note) unison.Dockable {
	return NewTableDockable(filePath, tbl.NewNoteTableHeaders(false),
		tbl.NewNoteRowData(func() []*gurps.Note { return notes }, false))
}

// NewTableDockable creates a new TableDockable for list data files.
func NewTableDockable(filePath string, columnHeaders []unison.TableColumnHeader, topLevelRows func(table *unison.Table) []unison.TableRowData) *TableDockable {
	d := &TableDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
		scale:  settings.Global().General.InitialListUIScale,
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{Columns: 1})

	d.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range d.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom.Size[float32]{})
		d.table.ColumnSizes[i].AutoMinimum = pref.Width
		d.table.ColumnSizes[i].AutoMaximum = 800
		d.table.ColumnSizes[i].Minimum = pref.Width
		d.table.ColumnSizes[i].Maximum = 10000
	}
	d.table.SetTopLevelRows(topLevelRows(d.table))
	d.table.SizeColumnsToFit(true)

	d.tableHeader = unison.NewTableHeader(d.table, columnHeaders...)
	d.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := f64d4.FromString(s1); err == nil {
			var n2 f64d4.Int
			if n2, err = f64d4.FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	d.scroll.SetColumnHeader(d.tableHeader)
	d.scroll.SetContent(d.table, unison.FillBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	d.lockButton = unison.NewSVGButton(res.LockSVG)
	d.toggleLock()
	d.lockButton.ClickCallback = d.toggleLock

	d.hierarchyButton = unison.NewSVGButton(res.HierarchySVG)
	d.hierarchyButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Opens/closes all hierarchical rows"))
	d.hierarchyButton.ClickCallback = d.toggleHierarchy

	d.sizeToFitButton = unison.NewSVGButton(res.SizeToFitSVG)
	d.sizeToFitButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sets the width of each column to fit its contents"))
	d.sizeToFitButton.ClickCallback = d.sizeToFit

	d.scaleField = widget.NewPercentageField(func() int { return d.scale }, func(v int) {
		d.scale = v
		d.applyScale()
	}, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax)
	d.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))

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
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom.Insets[float32]{Bottom: 1}, false),
		unison.NewEmptyBorder(geom.Insets[float32]{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(d.lockButton)
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
	return d
}

func (d *TableDockable) applyScale() {
	s := float32(d.scale) / 100
	d.tableHeader.SetScale(s)
	d.table.SetScale(s)
	d.scroll.Sync()
}

// TitleIcon implements workspace.FileBackedDockable
func (d *TableDockable) TitleIcon(suggestedSize geom.Size[float32]) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (d *TableDockable) Title() string {
	return fs.BaseName(d.path)
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
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *TableDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *TableDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func (d *TableDockable) toggleLock() {
	d.locked = !d.locked
	if dsvg, ok := d.lockButton.Drawable.(*unison.DrawableSVG); ok {
		if d.locked {
			dsvg.SVG = res.LockSVG
			d.lockButton.Tooltip = unison.NewTooltipWithSecondaryText(i18n.Text("Locked"), i18n.Text("Click to enable editing"))
		} else {
			dsvg.SVG = res.UnlockedSVG
			d.lockButton.Tooltip = unison.NewTooltipWithSecondaryText(i18n.Text("Unlocked"), i18n.Text("Click to disable editing"))
		}
	}
	d.lockButton.MarkForRedraw()
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
