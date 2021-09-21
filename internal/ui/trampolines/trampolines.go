package trampolines

import "github.com/richardwilkes/unison"

// These functions are here to break what would otherwise be circular dependencies.

// MenuSetup sets up the menus for the given window.
var MenuSetup func(wnd *unison.Window)
