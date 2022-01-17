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
	"github.com/richardwilkes/gcs/internal/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const workspaceClientDataKey = "workspace"

// Workspace holds the data necessary to track the Workspace.
type Workspace struct {
	Window       *unison.Window
	TopDock      *unison.Dock
	Navigator    *Navigator
	DocumentDock *DocumentDock
}

// FromWindow returns the Workspace associated with the given Window, or nil.
func FromWindow(wnd *unison.Window) *Workspace {
	if wnd != nil {
		if data, ok := wnd.ClientData()[workspaceClientDataKey]; ok {
			if w, ok2 := data.(*Workspace); ok2 {
				return w
			}
		}
	}
	return nil
}

// NewWorkspace creates a new Workspace for the given Window.
func NewWorkspace(wnd *unison.Window) *Workspace {
	w := &Workspace{
		Window:       wnd,
		TopDock:      unison.NewDock(),
		Navigator:    newNavigator(),
		DocumentDock: NewDocumentDock(),
	}
	wnd.SetContent(w.TopDock)
	w.TopDock.DockTo(w.Navigator, nil, unison.LeftSide)
	w.TopDock.RootDockLayout().SetDividerPosition(300)
	dc := unison.DockContainerFor(w.Navigator)
	w.TopDock.DockTo(w.DocumentDock, dc, unison.RightSide)
	dc.SetCurrentDockable(w.Navigator)
	wnd.ClientData()[workspaceClientDataKey] = w
	wnd.WillCloseCallback = w.willClose
	return w
}

func (w *Workspace) willClose() {
	globalSettings := gurps.Global()
	globalSettings.LibraryExplorer.OpenRowKeys = w.Navigator.DisclosedPaths()
	if err := globalSettings.Save(); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to save global settings"), err)
	}
}
