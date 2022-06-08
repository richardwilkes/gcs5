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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	traitModifierColMap = map[int]int{
		0: gurps.TraitModifierDescriptionColumn,
		1: gurps.TraitModifierCostColumn,
		2: gurps.TraitModifierTagsColumn,
		3: gurps.TraitModifierReferenceColumn,
	}
	traitModifierInEditorColMap = map[int]int{
		0: gurps.TraitModifierEnabledColumn,
		1: gurps.TraitModifierDescriptionColumn,
		2: gurps.TraitModifierCostColumn,
		3: gurps.TraitModifierTagsColumn,
		4: gurps.TraitModifierReferenceColumn,
	}
)

type traitModifierProvider struct {
	colMap   map[int]int
	provider gurps.TraitModifierListProvider
}

// NewTraitModifiersProvider creates a new table provider for trait modifiers.
func NewTraitModifiersProvider(provider gurps.TraitModifierListProvider, forEditor bool) TableProvider {
	p := &traitModifierProvider{
		provider: provider,
	}
	if forEditor {
		p.colMap = traitModifierInEditorColMap
	} else {
		p.colMap = traitModifierColMap
	}
	return p
}

func (p *traitModifierProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *traitModifierProvider) DragKey() string {
	return gid.TraitModifier
}

func (p *traitModifierProvider) DragSVG() *unison.SVG {
	return res.GCSTraitModifiersSVG
}

func (p *traitModifierProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Trait Modifier"), i18n.Text("Trait Modifiers")
}

func (p *traitModifierProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.TraitModifierEnabledColumn:
			headers = append(headers, NewEnabledHeader(false))
		case gurps.TraitModifierDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Trait Modifier"), "", false))
		case gurps.TraitModifierCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost Modifier"), "", false))
		case gurps.TraitModifierTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", false))
		case gurps.TraitModifierReferenceColumn:
			headers = append(headers, NewPageRefHeader(false))
		default:
			jot.Fatalf(1, "invalid trait modifier column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *traitModifierProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.TraitModifierList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, false))
	}
	return rows
}

func (p *traitModifierProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *traitModifierProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.TraitModifierDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *traitModifierProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *traitModifierProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	OpenEditor[*gurps.TraitModifier](table, func(item *gurps.TraitModifier) {
		EditTraitModifier(owner, item)
	})
}

func (p *traitModifierProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant) {
	item := gurps.NewTraitModifier(p.Entity(), nil, variant == ContainerItemVariant)
	InsertItem[*gurps.TraitModifier](owner, table, item,
		func(target, parent *gurps.TraitModifier) { target.Parent = parent },
		func(target *gurps.TraitModifier) []*gurps.TraitModifier { return target.Children },
		func(target *gurps.TraitModifier, children []*gurps.TraitModifier) { target.Children = children },
		p.provider.TraitModifierList, p.provider.SetTraitModifierList, p.RowData,
		func(target *gurps.TraitModifier) uuid.UUID { return target.ID })
	EditTraitModifier(owner, item)
}

func (p *traitModifierProvider) DeleteSelection(table *unison.Table) {
	deleteTableSelection(table, p.provider.TraitModifierList(),
		func(nodes []*gurps.TraitModifier) { p.provider.SetTraitModifierList(nodes) },
		func(node *gurps.TraitModifier) **gurps.TraitModifier { return &node.Parent },
		func(node *gurps.TraitModifier) *[]*gurps.TraitModifier { return &node.Children })
}
