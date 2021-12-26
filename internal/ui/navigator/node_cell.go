package navigator

import (
	"github.com/richardwilkes/gcs/internal/ui/widget"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

func createNodeCell(ext, title string) *unison.Panel {
	panel := unison.NewPanel()
	panelLayout := &unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	}
	panel.SetLayout(panelLayout)
	s := unison.LabelFont.ResolvedFont().Size() + 5
	size := geom32.NewSize(s, s)
	info, ok := fileTypes[ext]
	if !ok {
		info, ok = fileTypes["file"]
	}
	if ok {
		p := info.svg.PathForSize(size)
		svgSize := info.svg.Size()
		if svgSize.Width != svgSize.Height {
			p = p.NewTranslatedPt(info.svg.OffsetToCenterWithinScaledSize(size))
		}
		icon := widget.NewIcon()
		icon.Path = p
		icon.Size = size
		icon.SetLayoutData(&unison.FlexLayoutData{
			HSpan: 1,
			VSpan: 1,
		})
		panel.AddChild(icon)
	} else {
		panelLayout.Columns = 1
	}
	label := unison.NewLabel()
	label.Text = title
	label.SetLayoutData(&unison.FlexLayoutData{
		HSpan: 1,
		VSpan: 1,
		HGrab: true,
	})
	panel.AddChild(label)
	panel.NeedsLayout = true
	return panel
}
