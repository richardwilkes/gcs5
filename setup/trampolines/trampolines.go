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

package trampolines

import "github.com/richardwilkes/unison"

// These functions are here to break what would otherwise be circular dependencies.

// MenuSetup sets up the menus for the given window.
var MenuSetup func(wnd *unison.Window)

// OpenFile attempts to open the given file path in the given window, which should contain a workspace. May pass nil for
// wnd to let it pick the first such window it discovers.
var OpenFile func(wnd *unison.Window, filePath string) (dockable unison.Dockable, wasOpen bool)

// OpenPageReference opens the given page reference in the given window, which should contain a workspace. May pass nil
// for wnd to let it pick the first such window it discovers.
var OpenPageReference func(wnd *unison.Window, ref, highlight string)
