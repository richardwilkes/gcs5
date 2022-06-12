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

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	equipmentListColMap = map[int]int{
		0: gurps.EquipmentDescriptionColumn,
		1: gurps.EquipmentMaxUsesColumn,
		2: gurps.EquipmentTLColumn,
		3: gurps.EquipmentLCColumn,
		4: gurps.EquipmentCostColumn,
		5: gurps.EquipmentWeightColumn,
		6: gurps.EquipmentTagsColumn,
		7: gurps.EquipmentReferenceColumn,
	}
	carriedEquipmentPageColMap = map[int]int{
		0:  gurps.EquipmentEquippedColumn,
		1:  gurps.EquipmentQuantityColumn,
		2:  gurps.EquipmentDescriptionColumn,
		3:  gurps.EquipmentUsesColumn,
		4:  gurps.EquipmentTLColumn,
		5:  gurps.EquipmentLCColumn,
		6:  gurps.EquipmentCostColumn,
		7:  gurps.EquipmentWeightColumn,
		8:  gurps.EquipmentExtendedCostColumn,
		9:  gurps.EquipmentExtendedWeightColumn,
		10: gurps.EquipmentReferenceColumn,
	}
	otherEquipmentPageColMap = map[int]int{
		0: gurps.EquipmentQuantityColumn,
		1: gurps.EquipmentDescriptionColumn,
		2: gurps.EquipmentUsesColumn,
		3: gurps.EquipmentTLColumn,
		4: gurps.EquipmentLCColumn,
		5: gurps.EquipmentCostColumn,
		6: gurps.EquipmentWeightColumn,
		7: gurps.EquipmentExtendedCostColumn,
		8: gurps.EquipmentExtendedWeightColumn,
		9: gurps.EquipmentReferenceColumn,
	}
)

type equipmentProvider struct {
	table    *unison.Table[*Node[*gurps.Equipment]]
	colMap   map[int]int
	provider gurps.EquipmentListProvider
	forPage  bool
	carried  bool
}

// NewEquipmentProvider creates a new table provider for equipment. 'carried' is only relevant if 'forPage' is true.
func NewEquipmentProvider(provider gurps.EquipmentListProvider, forPage, carried bool) widget.TableProvider[*Node[*gurps.Equipment]] {
	p := &equipmentProvider{
		provider: provider,
		forPage:  forPage,
		carried:  carried,
	}
	if forPage {
		if carried {
			p.colMap = carriedEquipmentPageColMap
		} else {
			p.colMap = otherEquipmentPageColMap
		}
	} else {
		p.colMap = equipmentListColMap
	}
	return p
}

func (p *equipmentProvider) SetTable(table *unison.Table[*Node[*gurps.Equipment]]) {
	p.table = table
}

func (p *equipmentProvider) RootRowCount() int {
	return len(p.equipmentList())
}

func (p *equipmentProvider) RootRows() []*Node[*gurps.Equipment] {
	data := p.equipmentList()
	rows := make([]*Node[*gurps.Equipment], 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode[*gurps.Equipment](p.table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *equipmentProvider) SetRootRows(rows []*Node[*gurps.Equipment]) {
	list := ExtractNodeDataFromList(rows)
	if p.carried {
		p.provider.SetCarriedEquipmentList(list)
	} else {
		p.provider.SetOtherEquipmentList(list)
	}
}

func (p *equipmentProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *equipmentProvider) DragKey() string {
	return gid.Equipment
}

func (p *equipmentProvider) DragSVG() *unison.SVG {
	return res.GCSEquipmentSVG
}

func (p *equipmentProvider) DropShouldMoveData(drop *unison.TableDrop[*Node[*gurps.Equipment]]) bool {
	// Within same table?
	if drop.Table == drop.TableDragData.Table {
		return true
	}
	// Within same dockable?
	dockable := unison.Ancestor[unison.Dockable](drop.Table)
	if dockable != nil && dockable == unison.Ancestor[unison.Dockable](drop.TableDragData.Table) {
		return true
	}
	return false
}

func (p *equipmentProvider) DropCopyRow(drop *unison.TableDrop[*Node[*gurps.Equipment]], row *Node[*gurps.Equipment]) *Node[*gurps.Equipment] {
	eqp := ExtractFromRowData[*gurps.Equipment](row).Clone(p.provider.Entity(), nil)
	return NewNode[*gurps.Equipment](drop.Table, nil, p.colMap, eqp, p.forPage)
}

func (p *equipmentProvider) DropSetRowChildren(_ *unison.TableDrop[*Node[*gurps.Equipment]], row *Node[*gurps.Equipment], children []*Node[*gurps.Equipment]) {
	list := make([]*gurps.Equipment, 0, len(children))
	for _, child := range children {
		list = append(list, ExtractFromRowData[*gurps.Equipment](child))
	}
	if row == nil {
		if p.carried {
			p.provider.SetCarriedEquipmentList(list)
		} else {
			p.provider.SetOtherEquipmentList(list)
		}
	} else {
		ExtractFromRowData[*gurps.Equipment](row).Children = list
		row.children = nil
	}
}

func (p *equipmentProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Equipment Item"), i18n.Text("Equipment Items")
}

func (p *equipmentProvider) Headers() []unison.TableColumnHeader[*Node[*gurps.Equipment]] {
	var headers []unison.TableColumnHeader[*Node[*gurps.Equipment]]
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.EquipmentEquippedColumn:
			headers = append(headers, NewEquippedHeader[*gurps.Equipment](p.forPage))
		case gurps.EquipmentQuantityColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("#"), i18n.Text("Quantity"), p.forPage))
		case gurps.EquipmentDescriptionColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](p.descriptionText(), "", p.forPage))
		case gurps.EquipmentUsesColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("Uses"), i18n.Text("The number of uses remaining"), p.forPage))
		case gurps.EquipmentMaxUsesColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("Uses"), i18n.Text("The maximum number of uses"), p.forPage))
		case gurps.EquipmentTLColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("TL"), i18n.Text("Tech Level"), p.forPage))
		case gurps.EquipmentLCColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("LC"), i18n.Text("Legality Class"), p.forPage))
		case gurps.EquipmentCostColumn:
			headers = append(headers, NewMoneyHeader[*gurps.Equipment](p.forPage))
		case gurps.EquipmentExtendedCostColumn:
			headers = append(headers, NewExtendedMoneyHeader[*gurps.Equipment](p.forPage))
		case gurps.EquipmentWeightColumn:
			headers = append(headers, NewWeightHeader[*gurps.Equipment](p.forPage))
		case gurps.EquipmentExtendedWeightColumn:
			headers = append(headers, NewExtendedWeightHeader[*gurps.Equipment](p.forPage))
		case gurps.EquipmentTagsColumn:
			headers = append(headers, NewHeader[*gurps.Equipment](i18n.Text("Tags"), "", p.forPage))
		case gurps.EquipmentReferenceColumn:
			headers = append(headers, NewPageRefHeader[*gurps.Equipment](p.forPage))
		default:
			jot.Fatalf(1, "invalid equipment column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *equipmentProvider) SyncHeader(headers []unison.TableColumnHeader[*Node[*gurps.Equipment]]) {
	if p.forPage {
		for i := 0; i < len(carriedEquipmentPageColMap); i++ {
			if carriedEquipmentPageColMap[i] == gurps.EquipmentDescriptionColumn {
				if header, ok2 := headers[i].(*PageTableColumnHeader[*gurps.Equipment]); ok2 {
					header.Label.Text = p.descriptionText()
				}
				break
			}
		}
	}
}

func (p *equipmentProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.EquipmentDescriptionColumn {
			return k
		}
	}
	return -1
}

func (p *equipmentProvider) ExcessWidthColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.EquipmentDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *equipmentProvider) descriptionText() string {
	title := i18n.Text("Equipment")
	if p.forPage {
		if entity, ok := p.provider.(*gurps.Entity); ok {
			if p.carried {
				title = fmt.Sprintf(i18n.Text("Carried Equipment (%s; $%s)"),
					entity.SheetSettings.DefaultWeightUnits.Format(entity.WeightCarried(false)),
					entity.WealthCarried().String())
			} else {
				title = fmt.Sprintf(i18n.Text("Other Equipment ($%s)"), entity.WealthNotCarried().String())
			}
		}
	}
	return title
}

func (p *equipmentProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Equipment]]) {
	OpenEditor[*gurps.Equipment](table, func(item *gurps.Equipment) { EditEquipment(owner, item, p.carried) })
}

func (p *equipmentProvider) CreateItem(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Equipment]], variant widget.ItemVariant) {
	topListFunc := p.provider.OtherEquipmentList
	setTopListFunc := p.provider.SetOtherEquipmentList
	if p.carried {
		topListFunc = p.provider.CarriedEquipmentList
		setTopListFunc = p.provider.SetCarriedEquipmentList
	}
	item := gurps.NewEquipment(p.Entity(), nil, variant == widget.ContainerItemVariant)
	InsertItem[*gurps.Equipment](owner, table, item,
		func(target *gurps.Equipment) []*gurps.Equipment { return target.Children },
		func(target *gurps.Equipment, children []*gurps.Equipment) { target.Children = children },
		topListFunc, setTopListFunc,
		func(_ *unison.Table[*Node[*gurps.Equipment]]) []*Node[*gurps.Equipment] { return p.RootRows() },
		func(target *gurps.Equipment) uuid.UUID { return target.ID })
	EditEquipment(owner, item, p.carried)
}

func (p *equipmentProvider) DeleteSelection(table *unison.Table[*Node[*gurps.Equipment]]) {
	list := p.equipmentList()
	var setList func([]*gurps.Equipment)
	if p.carried {
		setList = p.provider.SetCarriedEquipmentList
	} else {
		setList = p.provider.SetOtherEquipmentList
	}
	deleteTableSelection(table, list,
		func(nodes []*gurps.Equipment) { setList(nodes) },
		func(node *gurps.Equipment) *[]*gurps.Equipment { return &node.Children })
}

func (p *equipmentProvider) equipmentList() []*gurps.Equipment {
	if p.carried {
		return p.provider.CarriedEquipmentList()
	}
	return p.provider.OtherEquipmentList()
}
