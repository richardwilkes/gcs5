package navigator

import (
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

func createNodeCell(ext, title string) *unison.Panel {
	size := unison.LabelFont.ResolvedFont().Size() + 5
	info, ok := fileTypes[ext]
	if !ok {
		info, ok = fileTypes["file"]
	}
	label := unison.NewLabel()
	label.Text = title
	label.Drawable = &unison.DrawableSVG{
		SVG:  info.svg,
		Size: geom32.NewSize(size, size),
	}
	return label.AsPanel()
}
