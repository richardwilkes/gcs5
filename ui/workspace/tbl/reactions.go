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

var reactionsColMap = map[int]int{
	0: gurps.ConditionalModifierValueColumn,
	1: gurps.ConditionalModifierDescriptionColumn,
}

type reactionModProvider struct {
	provider gurps.ReactionModifierListProvider
}

// NewReactionModifiersProvider creates a new table provider for reaction modifiers.
func NewReactionModifiersProvider(provider gurps.ReactionModifierListProvider) TableProvider {
	return &reactionModProvider{
		provider: provider,
	}
}

func (p *reactionModProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *reactionModProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(reactionsColMap); i++ {
		switch conditionalModifierColMap[i] {
		case gurps.ConditionalModifierValueColumn:
			headers = append(headers, NewHeader("±", i18n.Text("Modifier"), true))
		case gurps.ConditionalModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Reaction"), "", true))
		default:
			jot.Fatalf(1, "invalid reaction modifier column: %d", reactionsColMap[i])
		}
	}
	return headers
}

func (p *reactionModProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.Reactions()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, conditionalModifierColMap, one, true))
	}
	return rows
}

func (p *reactionModProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *reactionModProvider) HierarchyColumnIndex() int {
	return -1
}

func (p *reactionModProvider) ExcessWidthColumnIndex() int {
	for k, v := range conditionalModifierColMap {
		if v == gurps.ConditionalModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *reactionModProvider) OpenEditor(_ widget.Rebuildable, _ *unison.Table) {
}

func (p *reactionModProvider) CreateItem(_ widget.Rebuildable, _ *unison.Table, _ bool) {
}
