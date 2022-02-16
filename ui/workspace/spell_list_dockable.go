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
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// SpellListDockable holds the view for a spell list.
type SpellListDockable struct {
	unison.Panel
	path   string
	scroll *unison.ScrollPanel
	table  *unison.Table
}

// NewSpellListDockable creates a new SpellListDockable for spell list files.
func NewSpellListDockable(filePath string) (*SpellListDockable, error) {
	spells, err := gurps.NewSpellsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := &SpellListDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	d.Self = d

	d.table.ColumnSizes = make([]unison.ColumnSize, spellColumnCount)
	rows := make([]unison.TableRowData, 0, len(spells))
	for _, one := range spells {
		rows = append(rows, NewSpellNode(d, one))
	}
	d.table.SetTopLevelRows(rows)
	d.table.SizeColumnsToFit(true)

	d.scroll.MouseWheelMultiplier = 4
	d.scroll.SetColumnHeader(unison.NewTableHeader(d.table,
		unison.NewTableColumnHeader(i18n.Text("Spell"), ""),
		unison.NewTableColumnHeader(i18n.Text("Resist"), i18n.Text("Resistance")),
		unison.NewTableColumnHeader(i18n.Text("Class"), ""),
		unison.NewTableColumnHeader(i18n.Text("College"), ""),
		unison.NewTableColumnHeader(i18n.Text("Cost"), i18n.Text("The mana cost to cast the spell")),
		unison.NewTableColumnHeader(i18n.Text("Maintain"), i18n.Text("The mana cost to maintain the spell")),
		unison.NewTableColumnHeader(i18n.Text("Time"), i18n.Text("The time required to cast the spell")),
		unison.NewTableColumnHeader(i18n.Text("Duration"), ""),
		unison.NewTableColumnHeader(i18n.Text("Diff"), i18n.Text("Difficulty")),
		unison.NewTableColumnHeader(i18n.Text("Category"), ""),
		newPageReferenceHeader(),
	))
	d.scroll.SetContent(d.table, unison.FillBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	d.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})
	d.AddChild(d.scroll)
	return d, nil
}

// TitleIcon implements FileBackedDockable
func (d *SpellListDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements FileBackedDockable
func (d *SpellListDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements FileBackedDockable
func (d *SpellListDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements FileBackedDockable
func (d *SpellListDockable) BackingFilePath() string {
	return d.path
}

// Modified implements FileBackedDockable
func (d *SpellListDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *SpellListDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *SpellListDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
