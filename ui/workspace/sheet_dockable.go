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
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// SheetDockable holds the view for a GURPS character sheet.
type SheetDockable struct {
	unison.Panel
	path   string
	scroll *unison.ScrollPanel
	entity *gurps.Entity
}

// NewSheetDockable creates a new SheetDockable for GURPS character sheet files.
func NewSheetDockable(filePath string) (*SheetDockable, error) {
	entity, err := gurps.NewEntityFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	d := &SheetDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		entity: entity,
	}
	d.Self = d

	label := unison.NewLabel()
	label.Text = "Not yet implemented…"
	label.HAlign = unison.MiddleAlignment
	label.VAlign = unison.MiddleAlignment
	label.Font = unison.LabelFont.Face().Font(24)
	d.scroll.SetContent(label, unison.FillBehavior)

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
func (d *SheetDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements FileBackedDockable
func (d *SheetDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements FileBackedDockable
func (d *SheetDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements FileBackedDockable
func (d *SheetDockable) BackingFilePath() string {
	return d.path
}

// Modified implements FileBackedDockable
func (d *SheetDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *SheetDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *SheetDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
