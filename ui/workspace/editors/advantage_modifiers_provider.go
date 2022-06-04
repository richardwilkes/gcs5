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
	advantageModifierColMap = map[int]int{
		0: gurps.AdvantageModifierDescriptionColumn,
		1: gurps.AdvantageModifierCostColumn,
		2: gurps.AdvantageModifierTagsColumn,
		3: gurps.AdvantageModifierReferenceColumn,
	}
	advantageModifierInEditorColMap = map[int]int{
		0: gurps.AdvantageModifierEnabledColumn,
		1: gurps.AdvantageModifierDescriptionColumn,
		2: gurps.AdvantageModifierCostColumn,
		3: gurps.AdvantageModifierTagsColumn,
		4: gurps.AdvantageModifierReferenceColumn,
	}
)

type advModProvider struct {
	colMap   map[int]int
	provider gurps.AdvantageModifierListProvider
}

// NewAdvantageModifiersProvider creates a new table provider for advantage modifiers.
func NewAdvantageModifiersProvider(provider gurps.AdvantageModifierListProvider, forEditor bool) TableProvider {
	p := &advModProvider{
		provider: provider,
	}
	if forEditor {
		p.colMap = advantageModifierInEditorColMap
	} else {
		p.colMap = advantageModifierColMap
	}
	return p
}

func (p *advModProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *advModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.AdvantageModifierEnabledColumn:
			headers = append(headers, NewEnabledHeader(false))
		case gurps.AdvantageModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Trait Modifier"), "", false))
		case gurps.AdvantageModifierCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost Modifier"), "", false))
		case gurps.AdvantageModifierTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", false))
		case gurps.AdvantageModifierReferenceColumn:
			headers = append(headers, NewPageRefHeader(false))
		default:
			jot.Fatalf(1, "invalid advantage modifier column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *advModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.AdvantageModifierList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, false))
	}
	return rows
}

func (p *advModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *advModProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.AdvantageModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *advModProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *advModProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	OpenEditor[*gurps.AdvantageModifier](table, func(item *gurps.AdvantageModifier) {
		EditAdvantageModifier(owner, item)
	})
}

func (p *advModProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant) {
	item := gurps.NewAdvantageModifier(p.Entity(), nil, variant == ContainerItemVariant)
	InsertItem[*gurps.AdvantageModifier](owner, table, item,
		func(target, parent *gurps.AdvantageModifier) {},
		func(target *gurps.AdvantageModifier) []*gurps.AdvantageModifier { return target.Children },
		func(target *gurps.AdvantageModifier, children []*gurps.AdvantageModifier) { target.Children = children },
		p.provider.AdvantageModifierList, p.provider.SetAdvantageModifierList, p.RowData,
		func(target *gurps.AdvantageModifier) uuid.UUID { return target.ID })
	EditAdvantageModifier(owner, item)
}
