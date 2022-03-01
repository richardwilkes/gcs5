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
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/sheet"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var (
	_ workspace.FileBackedDockable = &SheetDockable{}
	_ sheet.ScaleProvider          = &SheetDockable{}
)

// SheetDockable holds the view for a GURPS character sheet.
type SheetDockable struct {
	unison.Panel
	path       string
	scroll     *unison.ScrollPanel
	entity     *gurps.Entity
	scaleField *widget.PercentageField
	scale      int
}

// NewSheetDockable creates a new unison.Dockable for GURPS character sheet files.
func NewSheetDockable(filePath string) (unison.Dockable, error) {
	entity, err := gurps.NewEntityFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	d := &SheetDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		entity: entity,
		scale:  settings.Global().General.InitialUIScale,
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	pages := unison.NewPanel()
	pages.SetLayout(&unison.FlexLayout{
		Columns:  1,
		VSpacing: 1,
	})
	pages.AddChild(d.createFirstPage())
	d.scroll.SetContent(pages, unison.UnmodifiedBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	d.scaleField = widget.NewPercentageField(d.scale, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax, func(v int) {
		d.scale = v
		r := pages.FrameRect()
		_, r.Size, _ = pages.Sizes(geom32.Size{})
		pages.SetFrameRect(r)
	})
	d.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(d.scaleField)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	d.AddChild(toolbar)
	d.AddChild(d.scroll)
	return d, nil
}

// CurrentScale implements sheet.ScaleProvider
func (d *SheetDockable) CurrentScale() float32 {
	return float32(d.scale) / 100.0
}

// TitleIcon implements workspace.FileBackedDockable
func (d *SheetDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (d *SheetDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements workspace.FileBackedDockable
func (d *SheetDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (d *SheetDockable) BackingFilePath() string {
	return d.path
}

// Modified implements workspace.FileBackedDockable
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

func (d *SheetDockable) createFirstPage() *sheet.Page {
	p := sheet.NewPage(d.entity)
	p.AddChild(sheet.NewPortrait(d.entity))
	return p
}
