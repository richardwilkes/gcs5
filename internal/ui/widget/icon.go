package widget

import (
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Icon provides a simple icon widget that draws the path it is given.
type Icon struct {
	unison.Panel
	Path *unison.Path
	Size geom32.Size
	Ink  unison.Ink
}

// NewIcon creates a new Icon.
func NewIcon() *Icon {
	icon := &Icon{}
	icon.Self = icon
	icon.SetSizer(func(_ geom32.Size) (min, pref, max geom32.Size) {
		return icon.Size, icon.Size, icon.Size
	})
	icon.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		if icon.Path != nil {
			gc.DrawPath(icon.Path, unison.ChooseInk(icon.Ink, unison.OnBackgroundColor).Paint(gc, rect, unison.Fill))
		}
	}
	return icon
}
