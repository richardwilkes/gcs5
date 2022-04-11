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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var equipmentModifierColMap = map[int]int{
	0: gurps.EquipmentModifierDescriptionColumn,
	1: gurps.EquipmentModifierTechLevelColumn,
	2: gurps.EquipmentModifierCostColumn,
	3: gurps.EquipmentModifierWeightColumn,
	4: gurps.EquipmentModifierCategoryColumn,
	5: gurps.EquipmentModifierReferenceColumn,
}

type eqpModProvider struct {
	provider gurps.EquipmentModifierListProvider
}

// NewEquipmentModifiersProvider creates a new table provider for equipment modifiers.
func NewEquipmentModifiersProvider(provider gurps.EquipmentModifierListProvider) TableProvider {
	return &eqpModProvider{
		provider: provider,
	}
}

func (p *eqpModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(equipmentModifierColMap); i++ {
		switch equipmentModifierColMap[i] {
		case gurps.EquipmentModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Equipment Modifier"), "", false))
		case gurps.EquipmentModifierTechLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), false))
		case gurps.EquipmentModifierCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost Adjustment"), "", false))
		case gurps.EquipmentModifierWeightColumn:
			headers = append(headers, NewHeader(i18n.Text("Weight Adjustment"), "", false))
		case gurps.EquipmentModifierCategoryColumn:
			headers = append(headers, NewHeader(i18n.Text("Category"), "", false))
		case gurps.EquipmentModifierReferenceColumn:
			headers = append(headers, NewPageRefHeader(false))
		default:
			jot.Fatalf(1, "invalid equipment modifier column: %d", equipmentModifierColMap[i])
		}
	}
	return headers
}

func (p *eqpModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.EquipmentModifierList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, equipmentModifierColMap, one, false))
	}
	return rows
}

func (p *eqpModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *eqpModProvider) HierarchyColumnIndex() int {
	for k, v := range equipmentModifierColMap {
		if v == gurps.EquipmentModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *eqpModProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}
