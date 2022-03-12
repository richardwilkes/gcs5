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

	s := &Sheet{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		entity: entity,
		pages:  unison.NewPanel(),
	}
	s.Self = s
	s.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	s.pages.SetLayout(&unison.FlexLayout{
		Columns:  1,
		VSpacing: 1,
	})
	s.pages.AddChild(s.createFirstPage())
	s.scroll.SetContent(s.pages, unison.UnmodifiedBehavior)
	s.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	s.scroll.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	scale := settings.Global().General.InitialSheetUIScale
	s.scaleField = widget.NewPercentageField(scale, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax, s.applyScale)
	s.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))

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
	toolbar.AddChild(s.scaleField)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	s.AddChild(toolbar)
	s.AddChild(s.scroll)

	s.applyScale(scale)
	return s, nil
}

func (s *Sheet) applyScale(scale int) {
	s.pages.SetScale(float32(scale) / 100)
	s.scroll.Sync()
}

// TitleIcon implements workspace.FileBackedDockable
func (s *Sheet) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(s.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (s *Sheet) Title() string {
	return fs.BaseName(s.path)
}

// Tooltip implements workspace.FileBackedDockable
func (s *Sheet) Tooltip() string {
	return s.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (s *Sheet) BackingFilePath() string {
	return s.path
}

// Modified implements workspace.FileBackedDockable
func (s *Sheet) Modified() bool {
	return s.MiscPanel.Modified
}

// MayAttemptClose implements unison.TabCloser
func (s *Sheet) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (s *Sheet) AttemptClose() {
	if dc := unison.DockContainerFor(s); dc != nil {
		dc.Close(s)
	}
}

func (s *Sheet) createFirstPage() *Page {
	p := NewPage(s.entity)
	p.AddChild(s.createFirstRow())
	p.AddChild(s.createSecondRow())
	return p
}

func (s *Sheet) createFirstRow() *unison.Panel {
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

	s.PortaitPanel = NewPortraitPanel(s.entity)
	s.IdentityPanel = NewIdentityPanel(s.entity)
	s.MiscPanel = NewMiscPanel(s.entity)
	s.DescriptionPanel = NewDescriptionPanel(s.entity)
	s.PointsPanel = NewPointsPanel(s.entity)

	p.AddChild(s.PortaitPanel)
	p.AddChild(s.IdentityPanel)
	p.AddChild(s.MiscPanel)
	p.AddChild(s.PointsPanel)
	p.AddChild(s.DescriptionPanel)
	return p
}

func (s *Sheet) createSecondRow() *unison.Panel {
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

	s.PrimaryAttrPanel = NewPrimaryAttrPanel(s.entity)
	s.SecondaryAttrPanel = NewSecondaryAttrPanel(s.entity)
	s.PointPoolsPanel = NewPointPoolsPanel(s.entity)
	s.BodyPanel = NewBodyPanel(s.entity)
	s.EncumbrancePanel = NewEncumbrancePanel(s.entity)
	s.LiftingPanel = NewLiftingPanel(s.entity)
	s.DamagePanel = NewDamagePanel(s.entity)

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
	endWrapper.AddChild(s.EncumbrancePanel)
	endWrapper.AddChild(s.LiftingPanel)

	p.AddChild(s.PrimaryAttrPanel)
	p.AddChild(s.SecondaryAttrPanel)
	p.AddChild(s.BodyPanel)
	p.AddChild(endWrapper)
	p.AddChild(s.DamagePanel)
	p.AddChild(s.PointPoolsPanel)

	return p
}

// MarkModified updates the modification timestamp and marks the entity as modified.
func MarkModified(p unison.Paneler) {
	panel := p.AsPanel()
	for panel != nil {
		panel.NeedsLayout = true
		if um, ok := panel.Self.(*Sheet); ok {
			um.MiscPanel.MarkModified()
			um.Sync()
			break
		}
		panel = panel.Parent()
	}
}

// Sync the panel to the current data.
func (s *Sheet) Sync() {
	s.PortaitPanel.Sync()
	s.IdentityPanel.Sync()
	s.MiscPanel.Sync()
	s.DescriptionPanel.Sync()
	s.PointsPanel.Sync()
	s.PrimaryAttrPanel.Sync()
	s.SecondaryAttrPanel.Sync()
	s.PointPoolsPanel.Sync()
	s.BodyPanel.Sync()
	s.EncumbrancePanel.Sync()
	s.LiftingPanel.Sync()
	s.DamagePanel.Sync()
	// TODO: call update on any embedded lists
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
