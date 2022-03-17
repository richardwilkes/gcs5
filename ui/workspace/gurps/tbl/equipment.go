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
	"github.com/richardwilkes/toolbox/i18n"
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
		6: gurps.EquipmentCategoryColumn,
		7: gurps.EquipmentReferenceColumn,
	}
	carriedEquipmentPageColMap = map[int]int{
		0: gurps.EquipmentEquippedColumn,
		1: gurps.EquipmentQuantityColumn,
		2: gurps.EquipmentDescriptionColumn,
		3: gurps.EquipmentUsesColumn,
		4: gurps.EquipmentCostColumn,
		5: gurps.EquipmentWeightColumn,
		6: gurps.EquipmentExtendedCostColumn,
		7: gurps.EquipmentExtendedWeightColumn,
		8: gurps.EquipmentReferenceColumn,
	}
	otherEquipmentPageColMap = map[int]int{
		0: gurps.EquipmentQuantityColumn,
		1: gurps.EquipmentDescriptionColumn,
		2: gurps.EquipmentUsesColumn,
		3: gurps.EquipmentCostColumn,
		4: gurps.EquipmentWeightColumn,
		5: gurps.EquipmentExtendedCostColumn,
		6: gurps.EquipmentExtendedWeightColumn,
		7: gurps.EquipmentReferenceColumn,
	}
)

// NewEquipmentTableHeaders creates a new set of table column headers for equipment. 'carried' is only relevant if
// 'forPage' is true.
func NewEquipmentTableHeaders(entity *gurps.Entity, forPage, carried bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	if forPage {
		if carried {
			headers = append(headers, NewEquippedHeader(true))
		}
		headers = append(headers, NewHeader(i18n.Text("Qty"), i18n.Text("Quantity"), true))
		if carried {
			headers = append(headers, NewHeader(fmt.Sprintf(i18n.Text("Carried Equipment (%s; $%s)"),
				entity.WeightCarried(false).String(), entity.WealthCarried().String()), "", true))
		} else {
			headers = append(headers, NewHeader(fmt.Sprintf(i18n.Text("Other Equipment ($%s)"),
				entity.WealthNotCarried().String()), "", true))
		}
		headers = append(headers,
			NewHeader(i18n.Text("Uses"), i18n.Text("The number of uses remaining"), true),
			NewMoneyHeader(true),
			NewWeightHeader(true),
			NewExtendedMoneyHeader(true),
			NewExtendedWeightHeader(true),
		)
	} else {
		headers = append(headers,
			NewHeader(i18n.Text("Equipment"), "", false),
			NewHeader(i18n.Text("Uses"), i18n.Text("The maximum number of uses"), false),
			NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), false),
			NewHeader(i18n.Text("LC"), i18n.Text("Legality Class"), false),
			NewMoneyHeader(false),
			NewWeightHeader(false),
			NewHeader(i18n.Text("Category"), "", false),
		)
	}
	return append(headers, NewPageRefHeader(forPage))
}

// NewEquipmentRowData creates a new table data provider function for equipment. 'carried' is only relevant if 'forPage'
// is true.
func NewEquipmentRowData(topLevelData []*gurps.Equipment, forPage, carried bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		var colMap map[int]int
		if forPage {
			if carried {
				colMap = carriedEquipmentPageColMap
			} else {
				colMap = otherEquipmentPageColMap
			}
		} else {
			colMap = equipmentListColMap
		}
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, colMap, one, forPage))
		}
		return rows
	}
}