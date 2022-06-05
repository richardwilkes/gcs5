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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	equipmentModifierColMap = map[int]int{
		0: gurps.EquipmentModifierDescriptionColumn,
		1: gurps.EquipmentModifierTechLevelColumn,
		2: gurps.EquipmentModifierCostColumn,
		3: gurps.EquipmentModifierWeightColumn,
		4: gurps.EquipmentModifierTagsColumn,
		5: gurps.EquipmentModifierReferenceColumn,
	}
	equipmentModifierInEditorColMap = map[int]int{
		0: gurps.EquipmentModifierEnabledColumn,
		1: gurps.EquipmentModifierDescriptionColumn,
		2: gurps.EquipmentModifierTechLevelColumn,
		3: gurps.EquipmentModifierCostColumn,
		4: gurps.EquipmentModifierWeightColumn,
		5: gurps.EquipmentModifierTagsColumn,
		6: gurps.EquipmentModifierReferenceColumn,
	}
)

type eqpModProvider struct {
	colMap   map[int]int
	provider gurps.EquipmentModifierListProvider
}

// NewEquipmentModifiersProvider creates a new table provider for equipment modifiers.
func NewEquipmentModifiersProvider(provider gurps.EquipmentModifierListProvider, forEditor bool) TableProvider {
	p := &eqpModProvider{
		provider: provider,
	}
	if forEditor {
		p.colMap = equipmentModifierInEditorColMap
	} else {
		p.colMap = equipmentModifierColMap
	}
	return p
}

func (p *eqpModProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *eqpModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.EquipmentModifierEnabledColumn:
			headers = append(headers, NewEnabledHeader(false))
		case gurps.EquipmentModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Equipment Modifier"), "", false))
		case gurps.EquipmentModifierTechLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), false))
		case gurps.EquipmentModifierCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost Adjustment"), "", false))
		case gurps.EquipmentModifierWeightColumn:
			headers = append(headers, NewHeader(i18n.Text("Weight Adjustment"), "", false))
		case gurps.EquipmentModifierTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", false))
		case gurps.EquipmentModifierReferenceColumn:
			headers = append(headers, NewPageRefHeader(false))
		default:
			jot.Fatalf(1, "invalid equipment modifier column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *eqpModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.EquipmentModifierList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, false))
	}
	return rows
}

func (p *eqpModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *eqpModProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.EquipmentModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *eqpModProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *eqpModProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	OpenEditor[*gurps.EquipmentModifier](table, func(item *gurps.EquipmentModifier) {
		EditEquipmentModifier(owner, item)
	})
}

func (p *eqpModProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant) {
	item := gurps.NewEquipmentModifier(p.Entity(), nil, variant == ContainerItemVariant)
	InsertItem[*gurps.EquipmentModifier](owner, table, item,
		func(target, parent *gurps.EquipmentModifier) {},
		func(target *gurps.EquipmentModifier) []*gurps.EquipmentModifier { return target.Children },
		func(target *gurps.EquipmentModifier, children []*gurps.EquipmentModifier) { target.Children = children },
		p.provider.EquipmentModifierList, p.provider.SetEquipmentModifierList, p.RowData,
		func(target *gurps.EquipmentModifier) uuid.UUID { return target.ID })
	EditEquipmentModifier(owner, item)
}

func (p *eqpModProvider) DeleteSelection(table *unison.Table) {
	deleteTableSelection(table, p.provider.EquipmentModifierList(),
		func(nodes []*gurps.EquipmentModifier) { p.provider.SetEquipmentModifierList(nodes) },
		func(node *gurps.EquipmentModifier) **gurps.EquipmentModifier { return &node.Parent },
		func(node *gurps.EquipmentModifier) *[]*gurps.EquipmentModifier { return &node.Children })
}
