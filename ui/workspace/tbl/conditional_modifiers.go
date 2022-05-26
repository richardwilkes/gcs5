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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var conditionalModifierColMap = map[int]int{
	0: gurps.ConditionalModifierValueColumn,
	1: gurps.ConditionalModifierDescriptionColumn,
}

type condModProvider struct {
	provider gurps.ConditionalModifierListProvider
}

// NewConditionalModifiersProvider creates a new table provider for conditional modifiers.
func NewConditionalModifiersProvider(provider gurps.ConditionalModifierListProvider) TableProvider {
	return &condModProvider{
		provider: provider,
	}
}

func (p *condModProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *condModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(conditionalModifierColMap); i++ {
		switch conditionalModifierColMap[i] {
		case gurps.ConditionalModifierValueColumn:
			headers = append(headers, NewHeader("±", i18n.Text("Modifier"), true))
		case gurps.ConditionalModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Condition"), "", true))
		default:
			jot.Fatalf(1, "invalid conditional modifier column: %d", conditionalModifierColMap[i])
		}
	}
	return headers
}

func (p *condModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.ConditionalModifiers()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, conditionalModifierColMap, one, true))
	}
	return rows
}

func (p *condModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *condModProvider) HierarchyColumnIndex() int {
	return -1
}

func (p *condModProvider) ExcessWidthColumnIndex() int {
	for k, v := range conditionalModifierColMap {
		if v == gurps.ConditionalModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *condModProvider) OpenEditor(_ widget.Rebuildable, _ *unison.Table) {
}

func (p *condModProvider) CreateItem(_ widget.Rebuildable, _ *unison.Table, _ bool) {
}
