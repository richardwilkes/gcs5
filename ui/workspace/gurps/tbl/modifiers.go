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
	"github.com/richardwilkes/unison"
)

var (
	advantageModifierColMap = map[int]int{
		0: gurps.AdvantageModifierDescriptionColumn,
		1: gurps.AdvantageModifierCostColumn,
		2: gurps.AdvantageModifierCategoryColumn,
		3: gurps.AdvantageModifierReferenceColumn,
	}
	equipmentModifierColMap = map[int]int{
		0: gurps.EquipmentModifierDescriptionColumn,
		1: gurps.EquipmentModifierTechLevelColumn,
		2: gurps.EquipmentModifierCostColumn,
		3: gurps.EquipmentModifierWeightColumn,
		4: gurps.EquipmentModifierCategoryColumn,
		5: gurps.EquipmentModifierReferenceColumn,
	}
	conditionalModifierColMap = map[int]int{
		0: gurps.ConditionalModifierValueColumn,
		1: gurps.ConditionalModifierDescriptionColumn,
	}
)

// NewAdvantageModifierTableHeaders creates a new set of table column headers for advantage modifiers.
func NewAdvantageModifierTableHeaders() []unison.TableColumnHeader {
	return []unison.TableColumnHeader{
		NewHeader(i18n.Text("Advantage Modifier"), "", false),
		NewHeader(i18n.Text("Cost Modifier"), "", false),
		NewHeader(i18n.Text("Category"), "", false),
		NewPageRefHeader(false),
	}
}

// NewAdvantageModifierRowData creates a new table data provider function for advantage modifiers.
func NewAdvantageModifierRowData(topLevelData []*gurps.AdvantageModifier) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, advantageModifierColMap, one, false))
		}
		return rows
	}
}

// NewEquipmentModifierTableHeaders creates a new set of table column headers for equipment modifiers.
func NewEquipmentModifierTableHeaders() []unison.TableColumnHeader {
	return []unison.TableColumnHeader{
		NewHeader(i18n.Text("Equipment Modifier"), "", false),
		NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), false),
		NewHeader(i18n.Text("Cost Adjustment"), "", false),
		NewHeader(i18n.Text("Weight Adjustment"), "", false),
		NewHeader(i18n.Text("Category"), "", false),
		NewPageRefHeader(false),
	}
}

// NewEquipmentModifierRowData creates a new table data provider function for equipment modifiers.
func NewEquipmentModifierRowData(topLevelData []*gurps.EquipmentModifier) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, equipmentModifierColMap, one, false))
		}
		return rows
	}
}

// NewConditionalModifierTableHeaders creates a new set of table column headers for conditional modifiers.
func NewConditionalModifierTableHeaders(descTitle string) []unison.TableColumnHeader {
	return []unison.TableColumnHeader{
		NewHeader(i18n.Text("Modifier"), "", true),
		NewHeader(descTitle, "", true),
	}
}

// NewConditionalModifierRowData creates a new table data provider function for conditional modifiers.
func NewConditionalModifierRowData(topLevelData []*gurps.ConditionalModifier) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, conditionalModifierColMap, one, true))
		}
		return rows
	}
}
