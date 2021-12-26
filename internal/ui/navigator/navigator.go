package navigator

import (
	"github.com/richardwilkes/gcs/internal/settings"
	"github.com/richardwilkes/unison"
)

// Navigator holds the workspace navigation panel.
type Navigator struct {
	unison.Panel
	scroll *unison.ScrollPanel
	table  *unison.Table
}

// NewNavigator creates a new workspace navigation panel.
func NewNavigator() *Navigator {
	n := &Navigator{
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	n.Self = n

	n.table.ColumnSizes = make([]unison.ColumnSize, 1)
	rows := make([]unison.TableRowData, 0, len(settings.Global.Libraries))
	for _, one := range settings.Global.Libraries {
		rows = append(rows, NewLibraryNode(n, one))
	}
	n.table.SetTopLevelRows(rows)
	n.table.SizeColumnsToFit(true)

	n.scroll.MouseWheelMultiplier = 2
	n.scroll.SetContent(n.table, unison.FillBehavior)
	n.scroll.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	n.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})
	n.AddChild(n.scroll)
	return n
}

func (n *Navigator) adjustTableSize() {
	n.table.SyncToModel()
	n.table.SizeColumnsToFit(true)
}
