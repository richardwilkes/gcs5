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
	"time"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	wsettings "github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/unison"
)

const (
	reactionsListIndex = iota
	conditionalModifiersListIndex
	meleeWeaponsListIndex
	rangedWeaponsListIndex
	advantagesListIndex
	skillsListIndex
	spellsListIndex
	carriedEquipmentListIndex
	otherEquipmentListIndex
	notesListIndex
	listCount
)

var (
	_ workspace.FileBackedDockable = &Sheet{}
	_ unison.UndoManagerProvider   = &Sheet{}
	_ widget.ModifiableRoot        = &Sheet{}
	_ widget.Rebuildable           = &Sheet{}
	_ widget.DockableKind          = &Sheet{}
)

// Sheet holds the view for a GURPS character sheet.
type Sheet struct {
	unison.Panel
	path               string
	undoMgr            *unison.UndoManager
	scroll             *unison.ScrollPanel
	entity             *gurps.Entity
	crc                uint64
	scale              int
	scaleField         *widget.PercentageField
	pages              *unison.Panel
	PortraitPanel      *PortraitPanel
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
	Lists              [listCount]*PageList
	rebuild            bool
	full               bool
}

// ActiveEntity returns the currently active entity.
func ActiveEntity() *gurps.Entity {
	d := workspace.ActiveDockable()
	if d == nil {
		return nil
	}
	if s, ok := d.(*Sheet); ok {
		return s.Entity()
	}
	return nil
}

// NewSheetFromFile loads a GURPS character sheet file and creates a new unison.Dockable for it.
func NewSheetFromFile(filePath string) (unison.Dockable, error) {
	entity, err := gurps.NewEntityFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewSheet(filePath, entity), nil
}

// NewSheet creates a new unison.Dockable for GURPS character sheet files.
func NewSheet(filePath string, entity *gurps.Entity) unison.Dockable {
	s := &Sheet{
		path:    filePath,
		undoMgr: unison.NewUndoManager(200, func(err error) { jot.Error(err) }),
		scroll:  unison.NewScrollPanel(),
		entity:  entity,
		crc:     entity.CRC64(),
		scale:   settings.Global().General.InitialSheetUIScale,
		pages:   unison.NewPanel(),
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
	s.pages.AddChild(s.createTopBlock())
	s.createLists()
	s.scroll.SetContent(s.pages, unison.UnmodifiedBehavior, unison.UnmodifiedBehavior)
	s.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	s.scroll.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	sheetSettingsButton := unison.NewSVGButton(res.SettingsSVG)
	sheetSettingsButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sheet Settings"))
	sheetSettingsButton.ClickCallback = func() { wsettings.ShowSheetSettings(s.entity) }

	s.scaleField = widget.NewPercentageField(func() int { return s.scale }, func(v int) {
		s.scale = v
		s.applyScale()
	}, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax)
	s.scaleField.SetMarksModified(false)
	s.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1},
		false), unison.NewEmptyBorder(unison.StdInsets())))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(sheetSettingsButton)
	toolbar.AddChild(s.scaleField)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	s.AddChild(toolbar)
	s.AddChild(s.scroll)

	s.applyScale()
	return s
}

// DockableKind implements widget.DockableKind
func (s *Sheet) DockableKind() string {
	return widget.SheetDockableKind
}

// Entity returns the entity this is displaying information for.
func (s *Sheet) Entity() *gurps.Entity {
	return s.entity
}

// UndoManager implements undo.Provider
func (s *Sheet) UndoManager() *unison.UndoManager {
	return s.undoMgr
}

func (s *Sheet) applyScale() {
	s.pages.SetScale(float32(s.scale) / 100)
	s.scroll.Sync()
}

// TitleIcon implements workspace.FileBackedDockable
func (s *Sheet) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(s.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (s *Sheet) Title() string {
	return fs.BaseName(s.path)
}

func (s *Sheet) String() string {
	return s.Title()
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
	return s.crc != s.entity.CRC64()
}

type closeWithEntity interface {
	unison.TabCloser
	CloseWithEntity(entity *gurps.Entity) bool
}

// MayAttemptClose implements unison.TabCloser
func (s *Sheet) MayAttemptClose() bool {
	allow := true
	for _, wnd := range unison.Windows() {
		if ws := workspace.FromWindow(wnd); ws != nil {
			ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
				for _, d := range dc.Dockables() {
					if fe, ok := d.(closeWithEntity); ok && fe.CloseWithEntity(s.entity) {
						if !fe.MayAttemptClose() {
							allow = false
							return true
						}
					}
				}
				return false
			})
		}
	}
	return allow
}

// AttemptClose implements unison.TabCloser
func (s *Sheet) AttemptClose() {
	for _, wnd := range unison.Windows() {
		if ws := workspace.FromWindow(wnd); ws != nil {
			ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
				for _, d := range dc.Dockables() {
					if fe, ok := d.(closeWithEntity); ok && fe.CloseWithEntity(s.entity) {
						fe.AttemptClose()
					}
				}
				return false
			})
		}
	}
	if dc := unison.DockContainerFor(s); dc != nil {
		dc.Close(s)
	}
}

func (s *Sheet) createTopBlock() *Page {
	p := NewPage(s.entity)
	p.AddChild(s.createFirstRow())
	p.AddChild(s.createSecondRow())
	return p
}

func (s *Sheet) createFirstRow() *unison.Panel {
	s.PortraitPanel = NewPortraitPanel(s.entity)
	s.IdentityPanel = NewIdentityPanel(s.entity)
	s.MiscPanel = NewMiscPanel(s.entity)
	s.DescriptionPanel = NewDescriptionPanel(s.entity)
	s.PointsPanel = NewPointsPanel(s.entity)

	right := unison.NewPanel()
	right.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 1,
		VSpacing: 1,
		HAlign:   unison.FillAlignment,
		VAlign:   unison.FillAlignment,
	})

	right.AddChild(s.IdentityPanel)
	right.AddChild(s.MiscPanel)
	right.AddChild(s.PointsPanel)
	right.AddChild(s.DescriptionPanel)

	p := unison.NewPanel()
	p.SetLayout(&portraitLayout{
		portrait: s.PortraitPanel,
		rest:     right,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.AddChild(s.PortraitPanel)
	p.AddChild(right)

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

func (s *Sheet) createLists() {
	children := s.pages.Children()
	if len(children) == 0 {
		return
	}
	page, ok := children[0].Self.(*Page)
	if !ok {
		return
	}
	children = page.Children()
	if len(children) < 2 {
		return
	}
	for i := len(children) - 1; i > 1; i-- {
		page.RemoveChildAtIndex(i)
	}
	// Add the various blocks, based on the layout preference.
	for _, col := range s.entity.SheetSettings.BlockLayout.ByRow() {
		rowPanel := unison.NewPanel()
		rowPanel.SetLayout(&unison.FlexLayout{
			Columns:  len(col),
			HSpacing: 1,
			HAlign:   unison.FillAlignment,
			VAlign:   unison.FillAlignment,
		})
		rowPanel.SetLayoutData(&unison.FlexLayoutData{
			HAlign: unison.FillAlignment,
			VAlign: unison.StartAlignment,
			HGrab:  true,
		})
		for _, c := range col {
			switch c {
			case gurps.BlockLayoutReactionsKey:
				s.Lists[reactionsListIndex] = NewReactionsPageList(s.entity)
				rowPanel.AddChild(s.Lists[reactionsListIndex])
			case gurps.BlockLayoutConditionalModifiersKey:
				s.Lists[conditionalModifiersListIndex] = NewConditionalModifiersPageList(s.entity)
				rowPanel.AddChild(s.Lists[conditionalModifiersListIndex])
			case gurps.BlockLayoutMeleeKey:
				s.Lists[meleeWeaponsListIndex] = NewMeleeWeaponsPageList(s.entity)
				rowPanel.AddChild(s.Lists[meleeWeaponsListIndex])
			case gurps.BlockLayoutRangedKey:
				s.Lists[rangedWeaponsListIndex] = NewRangedWeaponsPageList(s.entity)
				rowPanel.AddChild(s.Lists[rangedWeaponsListIndex])
			case gurps.BlockLayoutAdvantagesKey:
				s.Lists[advantagesListIndex] = NewAdvantagesPageList(s, s.entity)
				rowPanel.AddChild(s.Lists[advantagesListIndex])
			case gurps.BlockLayoutSkillsKey:
				s.Lists[skillsListIndex] = NewSkillsPageList(s, s.entity)
				rowPanel.AddChild(s.Lists[skillsListIndex])
			case gurps.BlockLayoutSpellsKey:
				s.Lists[spellsListIndex] = NewSpellsPageList(s.entity)
				rowPanel.AddChild(s.Lists[spellsListIndex])
			case gurps.BlockLayoutEquipmentKey:
				s.Lists[carriedEquipmentListIndex] = NewCarriedEquipmentPageList(s.entity)
				rowPanel.AddChild(s.Lists[carriedEquipmentListIndex])
			case gurps.BlockLayoutOtherEquipmentKey:
				s.Lists[otherEquipmentListIndex] = NewOtherEquipmentPageList(s.entity)
				rowPanel.AddChild(s.Lists[otherEquipmentListIndex])
			case gurps.BlockLayoutNotesKey:
				s.Lists[notesListIndex] = NewNotesPageList(s, s.entity)
				rowPanel.AddChild(s.Lists[notesListIndex])
			}
		}
		page.AddChild(rowPanel)
	}
	page.ApplyPreferredSize()
}

// MarkModified implements widget.ModifiableRoot.
func (s *Sheet) MarkModified() {
	s.MiscPanel.UpdateModified()
	widget.DeepSync(s)
	if dc := unison.DockContainerFor(s); dc != nil {
		dc.UpdateTitle(s)
	}
}

// SheetSettingsUpdated implements gurps.SheetSettingsResponder.
func (s *Sheet) SheetSettingsUpdated(entity *gurps.Entity, blockLayout bool) {
	if s.entity == entity {
		s.MarkForRebuild(blockLayout)
	}
}

// MarkForRebuild implements widget.Rebuildable.
func (s *Sheet) MarkForRebuild(full bool) {
	if full {
		s.full = full
	}
	if !s.rebuild {
		s.rebuild = true
		unison.InvokeTaskAfter(func() {
			doFull := s.full
			s.rebuild = false
			s.full = false
			s.entity.Recalculate()
			if doFull {
				selMap := make([]map[uuid.UUID]bool, listCount)
				for i, one := range s.Lists {
					if one != nil {
						selMap[i] = one.RecordSelection()
					}
				}
				defer func() {
					for i, one := range s.Lists {
						if one != nil {
							one.ApplySelection(selMap[i])
						}
					}
				}()
				s.createLists()
			}
			widget.DeepSync(s)
		}, 50*time.Millisecond)
	}
}

func drawBandedBackground(p unison.Paneler, gc *unison.Canvas, rect unison.Rect, start, step int) {
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