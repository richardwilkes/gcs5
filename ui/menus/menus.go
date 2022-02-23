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

package menus

import (
	"sync"

	"github.com/richardwilkes/gcs/ui/about"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

var registerKeyBindingsOnce sync.Once

// Setup the menu bar for the window.
func Setup(wnd *unison.Window) {
	registerKeyBindingsOnce.Do(func() {
		registerFileMenuActions()
		registerEditMenuActions()
		registerItemMenuActions()
		registerLibraryMenuActions()
		registerSettingsMenuActions()
	})
	unison.DefaultMenuFactory().BarForWindow(wnd, func(bar unison.Menu) {
		unison.InsertStdMenus(bar, about.Show, nil, nil)
		std := bar.Item(unison.PreferencesItemID)
		std.Menu().RemoveItem(std.Index())
		setupFileMenu(bar)
		setupEditMenu(bar)
		i := bar.Item(unison.EditMenuID).Index() + 1
		f := bar.Factory()
		bar.InsertMenu(i, createItemMenu(f))
		i++
		bar.InsertMenu(i, f.NewMenu(LibraryMenuID, i18n.Text("Library"), updateLibraryMenu))
		i++
		bar.InsertMenu(i, createSettingsMenu(f))
		setupHelpMenu(bar)
	})
}

// TODO: Implement each call site
func unimplemented(a *unison.Action, _ interface{}) {
	unison.ErrorDialogWithMessage("Unimplemented Action:", a.Title)
}
