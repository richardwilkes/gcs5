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

package tbl

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/editors"
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
	colMap   map[int]int
	provider gurps.EquipmentListProvider
	forPage  bool
	carried  bool
}

// NewEquipmentProvider creates a new table provider for equipment. 'carried' is only relevant if 'forPage' is true.
func NewEquipmentProvider(provider gurps.EquipmentListProvider, forPage, carried bool) TableProvider {
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

func (p *equipmentProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *equipmentProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.EquipmentEquippedColumn:
			headers = append(headers, NewEquippedHeader(p.forPage))
		case gurps.EquipmentQuantityColumn:
			headers = append(headers, NewHeader(i18n.Text("#"), i18n.Text("Quantity"), p.forPage))
		case gurps.EquipmentDescriptionColumn:
			headers = append(headers, NewHeader(p.descriptionText(), "", p.forPage))
		case gurps.EquipmentUsesColumn:
			headers = append(headers, NewHeader(i18n.Text("Uses"), i18n.Text("The number of uses remaining"), p.forPage))
		case gurps.EquipmentMaxUsesColumn:
			headers = append(headers, NewHeader(i18n.Text("Uses"), i18n.Text("The maximum number of uses"), p.forPage))
		case gurps.EquipmentTLColumn:
			headers = append(headers, NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), p.forPage))
		case gurps.EquipmentLCColumn:
			headers = append(headers, NewHeader(i18n.Text("LC"), i18n.Text("Legality Class"), p.forPage))
		case gurps.EquipmentCostColumn:
			headers = append(headers, NewMoneyHeader(p.forPage))
		case gurps.EquipmentExtendedCostColumn:
			headers = append(headers, NewExtendedMoneyHeader(p.forPage))
		case gurps.EquipmentWeightColumn:
			headers = append(headers, NewWeightHeader(p.forPage))
		case gurps.EquipmentExtendedWeightColumn:
			headers = append(headers, NewExtendedWeightHeader(p.forPage))
		case gurps.EquipmentTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", p.forPage))
		case gurps.EquipmentReferenceColumn:
			headers = append(headers, NewPageRefHeader(p.forPage))
		default:
			jot.Fatalf(1, "invalid equipment column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *equipmentProvider) RowData(table *unison.Table) []unison.TableRowData {
	var data []*gurps.Equipment
	if p.carried {
		data = p.provider.CarriedEquipmentList()
	} else {
		data = p.provider.OtherEquipmentList()
	}
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *equipmentProvider) SyncHeader(headers []unison.TableColumnHeader) {
	if p.forPage {
		for i := 0; i < len(carriedEquipmentPageColMap); i++ {
			if carriedEquipmentPageColMap[i] == gurps.EquipmentDescriptionColumn {
				if header, ok2 := headers[i].(*PageTableColumnHeader); ok2 {
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

func (p *equipmentProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	for _, row := range table.SelectedRows(false) {
		if node, ok := row.(*Node); ok {
			var e *gurps.Equipment
			if e, ok = node.Data().(*gurps.Equipment); ok {
				editors.EditEquipment(owner, e, p.carried)
			}
		}
	}
}

func (p *equipmentProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, container bool) {
	// TODO: Implement
}
