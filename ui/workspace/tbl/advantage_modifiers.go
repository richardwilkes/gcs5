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

var advantageModifierColMap = map[int]int{
	0: gurps.AdvantageModifierDescriptionColumn,
	1: gurps.AdvantageModifierCostColumn,
	2: gurps.AdvantageModifierTagsColumn,
	3: gurps.AdvantageModifierReferenceColumn,
}

type advModProvider struct {
	provider gurps.AdvantageModifierListProvider
}

// NewAdvantageModifiersProvider creates a new table provider for advantage modifiers.
func NewAdvantageModifiersProvider(provider gurps.AdvantageModifierListProvider) TableProvider {
	return &advModProvider{
		provider: provider,
	}
}

func (p *advModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(advantageModifierColMap); i++ {
		switch advantageModifierColMap[i] {
		case gurps.AdvantageModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Advantage Modifier"), "", false))
		case gurps.AdvantageModifierCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost Modifier"), "", false))
		case gurps.AdvantageModifierTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", false))
		case gurps.AdvantageModifierReferenceColumn:
			headers = append(headers, NewPageRefHeader(false))
		default:
			jot.Fatalf(1, "invalid advantage modifier column: %d", advantageModifierColMap[i])
		}
	}
	return headers
}

func (p *advModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.AdvantageModifierList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, advantageModifierColMap, one, false))
	}
	return rows
}

func (p *advModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *advModProvider) HierarchyColumnIndex() int {
	for k, v := range advantageModifierColMap {
		if v == gurps.AdvantageModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *advModProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}
