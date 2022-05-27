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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/editors"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	advantageListColMap = map[int]int{
		0: gurps.AdvantageDescriptionColumn,
		1: gurps.AdvantagePointsColumn,
		2: gurps.AdvantageTagsColumn,
		3: gurps.AdvantageReferenceColumn,
	}
	advantagePageColMap = map[int]int{
		0: gurps.AdvantageDescriptionColumn,
		1: gurps.AdvantagePointsColumn,
		2: gurps.AdvantageReferenceColumn,
	}
)

type advantagesProvider struct {
	colMap   map[int]int
	provider gurps.AdvantageListProvider
	forPage  bool
}

// NewAdvantagesProvider creates a new table provider for advantages.
func NewAdvantagesProvider(provider gurps.AdvantageListProvider, forPage bool) TableProvider {
	p := &advantagesProvider{
		provider: provider,
		forPage:  forPage,
	}
	if forPage {
		p.colMap = advantagePageColMap
	} else {
		p.colMap = advantageListColMap
	}
	return p
}

func (p *advantagesProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *advantagesProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.AdvantageDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Advantage / Disadvantage"), "", p.forPage))
		case gurps.AdvantagePointsColumn:
			headers = append(headers, NewHeader(i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		case gurps.AdvantageTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", p.forPage))
		case gurps.AdvantageReferenceColumn:
			headers = append(headers, NewPageRefHeader(p.forPage))
		default:
			jot.Fatalf(1, "invalid advantage column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *advantagesProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.AdvantageList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *advantagesProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *advantagesProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.AdvantageDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *advantagesProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *advantagesProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	OpenEditor[*gurps.Advantage](table, func(item *gurps.Advantage) { editors.EditAdvantage(owner, item) })
}

func (p *advantagesProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant) {
	CreateItem[*gurps.Advantage](owner, p.Entity(), table, variant == ContainerItemVariant, gurps.NewAdvantage,
		func(target *gurps.Advantage) []*gurps.Advantage { return target.Children },
		func(target *gurps.Advantage, children []*gurps.Advantage) { target.Children = children },
		p.provider.AdvantageList, p.provider.SetAdvantageList, p.RowData,
		func(target *gurps.Advantage) uuid.UUID { return target.ID })
}
