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

// SkillListDockable holds the view for a skill list.
type SkillListDockable struct {
	unison.Panel
	path   string
	scroll *unison.ScrollPanel
	table  *unison.Table
}

// NewSkillListDockable creates a new SkillListDockable for skill list files.
func NewSkillListDockable(filePath string) (*SkillListDockable, error) {
	skills, err := gurps.NewSkillsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := &SkillListDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	d.Self = d

	d.table.ColumnSizes = make([]unison.ColumnSize, skillColumnCount)
	rows := make([]unison.TableRowData, 0, len(skills))
	for _, one := range skills {
		rows = append(rows, NewSkillNode(d, one))
	}
	d.table.SetTopLevelRows(rows)
	d.table.SizeColumnsToFit(true)

	d.scroll.MouseWheelMultiplier = 4
	d.scroll.SetColumnHeader(unison.NewTableHeader(d.table,
		unison.NewTableColumnHeader(i18n.Text("Skill / Technique"), ""),
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
func (d *SkillListDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements FileBackedDockable
func (d *SkillListDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements FileBackedDockable
func (d *SkillListDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements FileBackedDockable
func (d *SkillListDockable) BackingFilePath() string {
	return d.path
}

// Modified implements FileBackedDockable
func (d *SkillListDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *SkillListDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *SkillListDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
