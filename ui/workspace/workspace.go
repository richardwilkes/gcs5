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
	"path"
	"strings"
	"time"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
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

// ShowUnableToLocateWorkspaceError displays an error dialog.
func ShowUnableToLocateWorkspaceError() {
	unison.ErrorDialogWithMessage(i18n.Text("Unable to locate workspace"), "")
}

// Activate attempts to locate an existing dockable that 'matcher' returns true for. If found, it will have been
// activated and focused.
func Activate(matcher func(d unison.Dockable) bool) (ws *Workspace, dc *unison.DockContainer, found bool) {
	if ws = Any(); ws == nil {
		jot.Error("no workspace available")
		return nil, nil, false
	}
	dc = ws.CurrentlyFocusedDockContainer()
	ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(container *unison.DockContainer) bool {
		for _, one := range container.Dockables() {
			if matcher(one) {
				found = true
				container.SetCurrentDockable(one)
				container.AcquireFocus()
				return true
			}
			if dc == nil {
				dc = container
			}
		}
		return false
	})
	return ws, dc, found
}

// ActiveDockable returns the currently active dockable in the active window.
func ActiveDockable() unison.Dockable {
	ws := FromWindow(unison.ActiveWindow())
	if ws == nil {
		return nil
	}
	dc := ws.CurrentlyFocusedDockContainer()
	if dc == nil {
		return nil
	}
	return dc.CurrentDockable()
}

// FromWindowOrAny first calls FromWindow(wnd) and if that fails to find a Workspace, then calls Any().
func FromWindowOrAny(wnd *unison.Window) *Workspace {
	ws := FromWindow(wnd)
	if ws == nil {
		ws = Any()
	}
	return ws
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

// Any first tries to return the workspace for the active window. If that fails, then it looks for any available
// workspace and returns that.
func Any() *Workspace {
	if ws := FromWindow(unison.ActiveWindow()); ws != nil {
		return ws
	}
	for _, wnd := range unison.Windows() {
		if ws := FromWindow(wnd); ws != nil {
			return ws
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
	dc := unison.DockContainerFor(w.Navigator)
	w.TopDock.DockTo(w.DocumentDock, dc, unison.RightSide)
	dc.SetCurrentDockable(w.Navigator)
	wnd.ClientData()[workspaceClientDataKey] = w
	wnd.WillCloseCallback = w.willClose
	// On some platforms, this needs to be done after a delay... but we do it without the delay, too, so that
	// well-behaved platforms don't flash
	w.TopDock.RootDockLayout().SetDividerPosition(settings.Global().LibraryExplorer.DividerPosition)
	unison.InvokeTaskAfter(func() {
		w.TopDock.RootDockLayout().SetDividerPosition(settings.Global().LibraryExplorer.DividerPosition)
	}, time.Millisecond)
	return w
}

func (w *Workspace) willClose() {
	globalSettings := settings.Global()
	globalSettings.LibraryExplorer.OpenRowKeys = w.Navigator.DisclosedPaths()
	if err := globalSettings.Save(); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to save global settings"), err)
	}
}

// CurrentlyFocusedDockContainer returns the currently focused DockContainer, if any.
func (w *Workspace) CurrentlyFocusedDockContainer() *unison.DockContainer {
	if focus := w.Window.Focus(); focus != nil {
		if dc := unison.DockContainerFor(focus); dc != nil && dc.Dock == w.DocumentDock.Dock {
			return dc
		}
	}
	return nil
}

// LocateFileBackedDockable searches for a FileBackedDockable with the given path.
func (w *Workspace) LocateFileBackedDockable(filePath string) FileBackedDockable {
	var dockable FileBackedDockable
	w.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
		for _, one := range dc.Dockables() {
			if fbd, ok := one.(FileBackedDockable); ok {
				if filePath == fbd.BackingFilePath() {
					dockable = fbd
					return true
				}
			}
		}
		return false
	})
	return dockable
}

// LocateDockContainerForExtension searches for the first FileBackedDockable with the given extension and returns its
// DockContainer.
func (w *Workspace) LocateDockContainerForExtension(ext ...string) *unison.DockContainer {
	var extDC *unison.DockContainer
	w.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
		if DockContainerHoldsExtension(dc, ext...) {
			extDC = dc
			return true
		}
		return false
	})
	return extDC
}

// DockContainerHoldsExtension returns true if an immediate child of the given DockContainer has a FileBackedDockable
// with the given extension.
func DockContainerHoldsExtension(dc *unison.DockContainer, ext ...string) bool {
	for _, one := range dc.Dockables() {
		if fbd, ok := one.(FileBackedDockable); ok {
			fbdExt := path.Ext(fbd.BackingFilePath())
			for _, e := range ext {
				if strings.EqualFold(fbdExt, e) {
					return true
				}
			}
		}
	}
	return false
}

// MayAttemptCloseOfDockable returns true if the given Dockable and any grouped Dockables associated with it may be
// closed.
func MayAttemptCloseOfDockable(d unison.Dockable) bool {
	allow := true
	for _, wnd := range unison.Windows() {
		if ws := FromWindow(wnd); ws != nil {
			ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
				for _, other := range dc.Dockables() {
					if fe, ok := other.(widget.GroupedCloser); ok && fe.CloseWithGroup(d) {
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

// AttemptCloseOfDockable attempts to close the given Dockable and any grouped Dockables associated with it.
func AttemptCloseOfDockable(d unison.Dockable) bool {
	allow := true
	for _, wnd := range unison.Windows() {
		if ws := FromWindow(wnd); ws != nil {
			ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
				for _, other := range dc.Dockables() {
					if fe, ok := other.(widget.GroupedCloser); ok && fe.CloseWithGroup(d) {
						if !fe.MayAttemptClose() {
							allow = false
							return true
						}
						if !fe.AttemptClose() {
							allow = false
							return true
						}
					}
				}
				return false
			})
		}
	}
	if !allow {
		return false
	}
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
	return true
}
