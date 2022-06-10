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
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/gcs/ui/workspace/editors"
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
	_ widget.DockableKind          = &Template{}
	_ unison.TabCloser             = &Template{}
)

// Template holds the view for a GURPS character template.
type Template struct {
	unison.Panel
	path              string
	undoMgr           *unison.UndoManager
	scroll            *unison.ScrollPanel
	template          *gurps.Template
	crc               uint64
	scale             int
	content           *templateContent
	scaleField        *widget.PercentageField
	Lists             [listCount]*PageList
	needsSaveAsPrompt bool
}

// NewTemplateFromFile loads a GURPS template file and creates a new unison.Dockable for it.
func NewTemplateFromFile(filePath string) (unison.Dockable, error) {
	template, err := gurps.NewTemplateFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	t := NewTemplate(filePath, template)
	t.needsSaveAsPrompt = false
	return t, nil
}

// NewTemplate creates a new unison.Dockable for GURPS template files.
func NewTemplate(filePath string, template *gurps.Template) *Template {
	d := &Template{
		path:              filePath,
		undoMgr:           unison.NewUndoManager(200, func(err error) { jot.Error(err) }),
		scroll:            unison.NewScrollPanel(),
		template:          template,
		scale:             settings.Global().General.InitialSheetUIScale,
		crc:               template.CRC64(),
		needsSaveAsPrompt: true,
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	d.scroll.SetContent(d.createContent(), unison.UnmodifiedBehavior, unison.UnmodifiedBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, theme.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	scaleTitle := i18n.Text("Scale")
	d.scaleField = widget.NewPercentageField(scaleTitle, func() int { return d.scale }, func(v int) {
		d.scale = v
		d.applyScale()
	}, gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax, false, false)
	d.scaleField.SetMarksModified(false)
	d.scaleField.Tooltip = unison.NewTooltipWithText(scaleTitle)

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1},
		false), unison.NewEmptyBorder(unison.StdInsets())))
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

	d.applyScale()

	d.InstallCmdHandlers(constants.SaveItemID, func(_ any) bool { return d.Modified() }, func(_ any) { d.save(false) })
	d.InstallCmdHandlers(constants.SaveAsItemID, unison.AlwaysEnabled, func(_ any) { d.save(true) })
	d.installNewItemCmdHandlers(constants.NewTraitItemID, constants.NewTraitContainerItemID, traitsListIndex)
	d.installNewItemCmdHandlers(constants.NewSkillItemID, constants.NewSkillContainerItemID, skillsListIndex)
	d.installNewItemCmdHandlers(constants.NewTechniqueItemID, -1, skillsListIndex)
	d.installNewItemCmdHandlers(constants.NewSpellItemID, constants.NewSpellContainerItemID, spellsListIndex)
	d.installNewItemCmdHandlers(constants.NewRitualMagicSpellItemID, -1, spellsListIndex)
	d.installNewItemCmdHandlers(constants.NewCarriedEquipmentItemID,
		constants.NewCarriedEquipmentContainerItemID, carriedEquipmentListIndex)
	d.installNewItemCmdHandlers(constants.NewOtherEquipmentItemID,
		constants.NewOtherEquipmentContainerItemID, otherEquipmentListIndex)
	d.installNewItemCmdHandlers(constants.NewNoteItemID, constants.NewNoteContainerItemID, notesListIndex)
	d.InstallCmdHandlers(constants.AddNaturalAttacksItemID, unison.AlwaysEnabled, func(_ any) {
		editors.InsertItem[*gurps.Trait](d, d.Lists[traitsListIndex].table, gurps.NewNaturalAttacks(nil, nil),
			func(target, parent *gurps.Trait) { target.Parent = parent },
			func(target *gurps.Trait) []*gurps.Trait { return target.Children },
			func(target *gurps.Trait, children []*gurps.Trait) { target.Children = children },
			d.template.TraitList, d.template.SetTraitList, d.Lists[traitsListIndex].provider.RowData,
			func(target *gurps.Trait) uuid.UUID { return target.ID })
	})

	return d
}

func (d *Template) installNewItemCmdHandlers(itemID, containerID, listIndex int) {
	variant := editors.NoItemVariant
	if containerID == -1 {
		variant = editors.AlternateItemVariant
	} else {
		d.InstallCmdHandlers(containerID, unison.AlwaysEnabled,
			func(_ any) { d.Lists[listIndex].CreateItem(d, editors.ContainerItemVariant) })
	}
	d.InstallCmdHandlers(itemID, unison.AlwaysEnabled, func(_ any) { d.Lists[listIndex].CreateItem(d, variant) })
}

// DockableKind implements widget.DockableKind
func (d *Template) DockableKind() string {
	return widget.TemplateDockableKind
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
	if dc := unison.Ancestor[*unison.DockContainer](d); dc != nil {
		dc.UpdateTitle(d)
	}
}

// MayAttemptClose implements unison.TabCloser
func (d *Template) MayAttemptClose() bool {
	return workspace.MayAttemptCloseOfGroup(d)
}

// AttemptClose implements unison.TabCloser
func (d *Template) AttemptClose() bool {
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
	if dc := unison.Ancestor[*unison.DockContainer](d); dc != nil {
		dc.Close(d)
	}
	return true
}

func (d *Template) createContent() unison.Paneler {
	d.content = newTemplateContent()
	d.createLists()
	return d.content
}

func (d *Template) save(forceSaveAs bool) bool {
	success := false
	if forceSaveAs || d.needsSaveAsPrompt {
		success = workspace.SaveDockableAs(d, library.TemplatesExt, d.template.Save, func(path string) {
			d.crc = d.template.CRC64()
			d.path = path
		})
	} else {
		success = workspace.SaveDockable(d, d.template.Save, func() { d.crc = d.template.CRC64() })
	}
	if success {
		d.needsSaveAsPrompt = false
	}
	return success
}

func (d *Template) createLists() {
	h, v := d.scroll.Position()
	refocusOn := -1
	if wnd := d.Window(); wnd != nil {
		if focus := wnd.Focus(); focus != nil {
			for i, one := range d.Lists {
				if one.table.Self == focus.Self {
					refocusOn = i
					break
				}
			}
		}
	}
	d.content.RemoveAllChildren()
	for _, col := range settings.Global().Sheet.BlockLayout.ByRow() {
		rowPanel := unison.NewPanel()
		for _, c := range col {
			switch c {
			case gurps.BlockLayoutTraitsKey:
				d.Lists[traitsListIndex] = NewTraitsPageList(d, d.template)
				rowPanel.AddChild(d.Lists[traitsListIndex])
			case gurps.BlockLayoutSkillsKey:
				d.Lists[skillsListIndex] = NewSkillsPageList(d, d.template)
				rowPanel.AddChild(d.Lists[skillsListIndex])
			case gurps.BlockLayoutSpellsKey:
				d.Lists[spellsListIndex] = NewSpellsPageList(d, d.template)
				rowPanel.AddChild(d.Lists[spellsListIndex])
			case gurps.BlockLayoutEquipmentKey:
				d.Lists[carriedEquipmentListIndex] = NewCarriedEquipmentPageList(d, d.template)
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
	if refocusOn != -1 {
		d.Lists[refocusOn].table.RequestFocus()
	}
	d.scroll.SetPosition(h, v)
}

// SheetSettingsUpdated implements gurps.SheetSettingsResponder.
func (d *Template) SheetSettingsUpdated(entity *gurps.Entity, blockLayout bool) {
	if entity == nil {
		d.Rebuild(blockLayout)
	}
}

// Rebuild implements widget.Rebuildable.
func (d *Template) Rebuild(full bool) {
	if full {
		selMap := make([]map[uuid.UUID]bool, listCount)
		for i, one := range d.Lists {
			if one != nil {
				selMap[i] = one.RecordSelection()
			}
		}
		defer func() {
			for i, one := range d.Lists {
				if one != nil {
					one.ApplySelection(selMap[i])
				}
			}
		}()
		d.createLists()
	}
	widget.DeepSync(d)
	if dc := unison.Ancestor[*unison.DockContainer](d); dc != nil {
		dc.UpdateTitle(d)
	}
}
