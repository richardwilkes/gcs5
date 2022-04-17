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

package editors

import (
	"fmt"

	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

const editorGroup = "editors"

var (
	_ unison.Dockable            = &Editor{}
	_ unison.TabCloser           = &Editor{}
	_ widget.ModifiableRoot      = &Editor{}
	_ unison.UndoManagerProvider = &Editor{}
)

// Editor provides the base editor functionality.
type Editor struct {
	unison.Panel
	TabTitle           string
	IsModifiedCallback func() bool
	ApplyCallback      func()
	undoMgr            *unison.UndoManager
	applyButton        *unison.Button
	cancelButton       *unison.Button
	promptForSave      bool
}

// Setup the editor and display it.
func (e *Editor) Setup(ws *workspace.Workspace, dc *unison.DockContainer, initContent func(*unison.Panel)) {
	e.undoMgr = unison.NewUndoManager(100, func(err error) { jot.Error(err) })
	e.SetLayout(&unison.FlexLayout{Columns: 1})
	e.AddChild(e.createToolbar())
	content := unison.NewPanel()
	content.SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(unison.StdHSpacing * 2)))
	initContent(content)
	scroller := unison.NewScrollPanel()
	scroller.SetContent(content, unison.FollowBehavior, unison.FillBehavior)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	e.AddChild(scroller)
	e.promptForSave = true
	if dc != nil && dc.Group == editorGroup {
		dc.Stack(e, -1)
	} else if dc = ws.DocumentDock.ContainerForGroup(editorGroup); dc != nil {
		dc.Stack(e, -1)
	} else {
		ws.DocumentDock.DockTo(e, nil, unison.RightSide)
		if dc = unison.DockContainerFor(e); dc != nil && dc.Group == "" {
			dc.Group = editorGroup
		}
	}
}

func (e *Editor) createToolbar() unison.Paneler {
	toolbar := unison.NewPanel()
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(unison.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	e.applyButton = unison.NewSVGButton(res.CheckmarkSVG)
	e.applyButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Apply Changes"))
	e.applyButton.SetEnabled(false)
	e.applyButton.ClickCallback = func() {
		if e.ApplyCallback != nil {
			e.ApplyCallback()
		}
		e.promptForSave = false
		e.AttemptClose()
	}
	toolbar.AddChild(e.applyButton)
	e.cancelButton = unison.NewSVGButton(res.NotSVG)
	e.cancelButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Discard Changes"))
	e.cancelButton.SetEnabled(false)
	e.cancelButton.ClickCallback = func() {
		e.promptForSave = false
		e.AttemptClose()
	}
	toolbar.AddChild(e.cancelButton)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})
	return toolbar
}

// TitleIcon implements unison.Dockable
func (e *Editor) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  res.GCSNotesSVG,
		Size: suggestedSize,
	}
}

// Title implements unison.Dockable
func (e *Editor) Title() string {
	return e.TabTitle
}

// Tooltip implements unison.Dockable
func (e *Editor) Tooltip() string {
	return ""
}

// Modified implements unison.Dockable
func (e *Editor) Modified() bool {
	if e.IsModifiedCallback == nil {
		return false
	}
	modified := e.IsModifiedCallback()
	e.applyButton.SetEnabled(modified)
	e.cancelButton.SetEnabled(modified)
	return modified
}

// MarkModified implements widget.ModifiableRoot.
func (e *Editor) MarkModified() {
	if dc := unison.DockContainerFor(e); dc != nil {
		dc.UpdateTitle(e)
	}
}

// MayAttemptClose implements unison.TabCloser
func (e *Editor) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (e *Editor) AttemptClose() {
	if e.promptForSave && e.ApplyCallback != nil && e.IsModifiedCallback != nil && e.IsModifiedCallback() {
		if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Save changes made to\n%s?"), e.Title()), "") == unison.ModalResponseOK {
			e.ApplyCallback()
		}
	}
	if dc := unison.DockContainerFor(e); dc != nil {
		dc.Close(e)
	}
}

// UndoManager implements undo.Provider
func (e *Editor) UndoManager() *unison.UndoManager {
	return e.undoMgr
}
