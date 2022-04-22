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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/unison"
)

var (
	_ workspace.FileBackedDockable = &Template{}
	_ unison.UndoManagerProvider   = &Template{}
	_ widget.ModifiableRoot        = &Template{}
	_ widget.Rebuildable           = &Template{}
)

// Template holds the view for a GURPS character template.
type Template struct {
	unison.Panel
	path       string
	undoMgr    *unison.UndoManager
	scroll     *unison.ScrollPanel
	template   *gurps.Template
	crc        uint64
	scale      int
	content    *templateContent
	scaleField *widget.PercentageField
	Lists      [listCount]*PageList
	rebuild    bool
	full       bool
}

// NewTemplateFromFile loads a GURPS template file and creates a new unison.Dockable for it.
func NewTemplateFromFile(filePath string) (unison.Dockable, error) {
	template, err := gurps.NewTemplateFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewTemplate(filePath, template), nil
}

// NewTemplate creates a new unison.Dockable for GURPS template files.
func NewTemplate(filePath string, template *gurps.Template) unison.Dockable {
	t := &Template{
		path:     filePath,
		undoMgr:  unison.NewUndoManager(200, func(err error) { jot.Error(err) }),
		scroll:   unison.NewScrollPanel(),
		template: template,
		scale:    settings.Global().General.InitialSheetUIScale,
		crc:      template.CRC64(),
	}
	t.Self = t
	t.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	t.scroll.SetContent(t.createContent(), unison.UnmodifiedBehavior, unison.UnmodifiedBehavior)
	t.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	t.scroll.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	t.scaleField = widget.NewPercentageField(func() int { return t.scale }, func(v int) {
		t.scale = v
		t.applyScale()
	}, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax)
	t.scaleField.SetMarksModified(false)
	t.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(unison.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(t.scaleField)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	t.AddChild(toolbar)
	t.AddChild(t.scroll)

	t.applyScale()
	return t
}

func (d *Template) applyScale() {
	d.scroll.Content().AsPanel().SetScale(float32(d.scale) / 100)
	d.scroll.Sync()
}

// UndoManager implements undo.Provider
func (d *Template) UndoManager() *unison.UndoManager {
	return d.undoMgr
}

// TitleIcon implements workspace.FileBackedDockable
func (d *Template) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (d *Template) Title() string {
	return fs.BaseName(d.path)
}

func (d *Template) String() string {
	return d.Title()
}

// Tooltip implements workspace.FileBackedDockable
func (d *Template) Tooltip() string {
	return d.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (d *Template) BackingFilePath() string {
	return d.path
}

// Modified implements workspace.FileBackedDockable
func (d *Template) Modified() bool {
	return d.crc != d.template.CRC64()
}

// MarkModified implements widget.ModifiableRoot.
func (d *Template) MarkModified() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.UpdateTitle(d)
	}
}

// MayAttemptClose implements unison.TabCloser
func (d *Template) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *Template) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func (d *Template) createContent() unison.Paneler {
	d.content = newTemplateContent()
	d.createLists()
	return d.content
}

func (d *Template) createLists() {
	d.content.RemoveAllChildren()
	for _, col := range settings.Global().Sheet.BlockLayout.ByRow() {
		rowPanel := unison.NewPanel()
		for _, c := range col {
			switch c {
			case gurps.BlockLayoutAdvantagesKey:
				d.Lists[advantagesListIndex] = NewAdvantagesPageList(d, d.template)
				rowPanel.AddChild(d.Lists[advantagesListIndex])
			case gurps.BlockLayoutSkillsKey:
				d.Lists[skillsListIndex] = NewSkillsPageList(d.template)
				rowPanel.AddChild(d.Lists[skillsListIndex])
			case gurps.BlockLayoutSpellsKey:
				d.Lists[spellsListIndex] = NewSpellsPageList(d.template)
				rowPanel.AddChild(d.Lists[spellsListIndex])
			case gurps.BlockLayoutEquipmentKey:
				d.Lists[carriedEquipmentListIndex] = NewCarriedEquipmentPageList(d.template)
				rowPanel.AddChild(d.Lists[carriedEquipmentListIndex])
			case gurps.BlockLayoutNotesKey:
				d.Lists[notesListIndex] = NewNotesPageList(d, d.template)
				rowPanel.AddChild(d.Lists[notesListIndex])
			}
		}
		if len(rowPanel.Children()) != 0 {
			rowPanel.SetLayout(&unison.FlexLayout{
				Columns:      len(rowPanel.Children()),
				HSpacing:     1,
				HAlign:       unison.FillAlignment,
				EqualColumns: true,
			})
			rowPanel.SetLayoutData(&unison.FlexLayoutData{
				HAlign: unison.FillAlignment,
				VAlign: unison.StartAlignment,
				HGrab:  true,
			})
			d.content.AddChild(rowPanel)
		}
	}
	d.content.ApplyPreferredSize()
}

// SheetSettingsUpdated implements gurps.SheetSettingsResponder.
func (d *Template) SheetSettingsUpdated(entity *gurps.Entity, blockLayout bool) {
	if entity == nil {
		d.MarkForRebuild(blockLayout)
	}
}

// MarkForRebuild implements widget.Rebuildable.
func (d *Template) MarkForRebuild(full bool) {
	if full {
		d.full = full
	}
	if !d.rebuild {
		d.rebuild = true
		unison.InvokeTaskAfter(func() {
			doFull := d.full
			d.rebuild = false
			d.full = false
			if doFull {
				selMap := make([]map[uuid.UUID]bool, listCount)
				for i, one := range d.Lists {
					selMap[i] = one.RecordSelection()
				}
				defer func() {
					for i, one := range d.Lists {
						one.ApplySelection(selMap[i])
					}
				}()
				d.createLists()
			}
			widget.DeepSync(d)
		}, 50*time.Millisecond)
	}
}
