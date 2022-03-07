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
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var _ workspace.FileBackedDockable = &Sheet{}

// Sheet holds the view for a GURPS character sheet.
type Sheet struct {
	unison.Panel
	path               string
	scroll             *unison.ScrollPanel
	entity             *gurps.Entity
	scaleField         *widget.PercentageField
	pages              *unison.Panel
	PortaitPanel       *PortraitPanel
	IdentityPanel      *IdentityPanel
	MiscPanel          *MiscPanel
	DescriptionPanel   *DescriptionPanel
	PointsPanel        *PointsPanel
	PrimaryAttrPanel   *PrimaryAttrPanel
	SecondaryAttrPanel *SecondaryAttrPanel
	PointPoolsPanel    *PointPoolsPanel
	BodyPanel          *BodyPanel
	EncumbrancePanel   *EncumbrancePanel
	LiftingPanel       *LiftingPanel
	DamagePanel        *DamagePanel
}

// NewSheet creates a new unison.Dockable for GURPS character sheet files.
func NewSheet(filePath string) (unison.Dockable, error) {
	entity, err := gurps.NewEntityFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	d := &Sheet{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		entity: entity,
		pages:  unison.NewPanel(),
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	d.pages.SetLayout(&unison.FlexLayout{
		Columns:  1,
		VSpacing: 1,
	})
	d.pages.AddChild(d.createFirstPage())
	d.scroll.SetContent(d.pages, unison.UnmodifiedBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	scale := settings.Global().General.InitialSheetUIScale
	d.scaleField = widget.NewPercentageField(scale, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax, d.applyScale)
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

	d.applyScale(scale)
	return d, nil
}

func (d *Sheet) applyScale(scale int) {
	d.pages.SetScale(float32(scale) / 100)
	d.scroll.Sync()
}

// TitleIcon implements workspace.FileBackedDockable
func (d *Sheet) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (d *Sheet) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements workspace.FileBackedDockable
func (d *Sheet) Tooltip() string {
	return d.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (d *Sheet) BackingFilePath() string {
	return d.path
}

// Modified implements workspace.FileBackedDockable
func (d *Sheet) Modified() bool {
	return d.MiscPanel.Modified
}

// MayAttemptClose implements unison.TabCloser
func (d *Sheet) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *Sheet) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func (d *Sheet) createFirstPage() *Page {
	p := NewPage(d.entity)
	p.AddChild(d.createFirstRow())
	p.AddChild(d.createSecondRow())
	return p
}

func (d *Sheet) createFirstRow() *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns:  4,
		HSpacing: 1,
		VSpacing: 1,
		HAlign:   unison.FillAlignment,
		VAlign:   unison.FillAlignment,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})

	d.PortaitPanel = NewPortraitPanel(d.entity)
	d.IdentityPanel = NewIdentityPanel(d.entity)
	d.MiscPanel = NewMiscPanel(d.entity)
	d.DescriptionPanel = NewDescriptionPanel(d.entity)
	d.PointsPanel = NewPointsPanel(d.entity)

	p.AddChild(d.PortaitPanel)
	p.AddChild(d.IdentityPanel)
	p.AddChild(d.MiscPanel)
	p.AddChild(d.PointsPanel)
	p.AddChild(d.DescriptionPanel)
	return p
}

func (d *Sheet) createSecondRow() *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns:  4,
		HSpacing: 1,
		VSpacing: 1,
		HAlign:   unison.FillAlignment,
		VAlign:   unison.FillAlignment,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})

	d.PrimaryAttrPanel = NewPrimaryAttrPanel(d.entity)
	d.SecondaryAttrPanel = NewSecondaryAttrPanel(d.entity)
	d.PointPoolsPanel = NewPointPoolsPanel(d.entity)
	d.BodyPanel = NewBodyPanel(d.entity)
	d.EncumbrancePanel = NewEncumbrancePanel(d.entity)
	d.LiftingPanel = NewLiftingPanel(d.entity)
	d.DamagePanel = NewDamagePanel(d.entity)

	endWrapper := unison.NewPanel()
	endWrapper.SetLayout(&unison.FlexLayout{
		Columns:  1,
		VSpacing: 1,
	})
	endWrapper.SetLayoutData(&unison.FlexLayoutData{
		VSpan:  3,
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	endWrapper.AddChild(d.EncumbrancePanel)
	endWrapper.AddChild(d.LiftingPanel)

	p.AddChild(d.PrimaryAttrPanel)
	p.AddChild(d.SecondaryAttrPanel)
	p.AddChild(d.BodyPanel)
	p.AddChild(endWrapper)
	p.AddChild(d.DamagePanel)
	p.AddChild(d.PointPoolsPanel)

	return p
}

func drawBandedBackground(p unison.Paneler, gc *unison.Canvas, rect geom32.Rect, start, step int) {
	gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	children := p.AsPanel().Children()
	for i := start; i < len(children); i += step {
		var ink unison.Ink
		if ((i-start)/step)&1 == 1 {
			ink = unison.BandingColor
		} else {
			ink = unison.ContentColor
		}
		r := children[i].FrameRect()
		r.X = rect.X
		r.Width = rect.Width
		gc.DrawRect(r, ink.Paint(gc, r, unison.Fill))
	}
}
