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
	"bytes"
	"compress/gzip"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	traitListColMap = map[int]int{
		0: gurps.TraitDescriptionColumn,
		1: gurps.TraitPointsColumn,
		2: gurps.TraitTagsColumn,
		3: gurps.TraitReferenceColumn,
	}
	traitPageColMap = map[int]int{
		0: gurps.TraitDescriptionColumn,
		1: gurps.TraitPointsColumn,
		2: gurps.TraitReferenceColumn,
	}
)

type traitsProvider struct {
	table    *unison.Table[*Node[*gurps.Trait]]
	colMap   map[int]int
	provider gurps.TraitListProvider
	forPage  bool
}

// NewTraitsProvider creates a new table provider for traits.
func NewTraitsProvider(provider gurps.TraitListProvider, forPage bool) widget.TableProvider[*Node[*gurps.Trait]] {
	p := &traitsProvider{
		provider: provider,
		forPage:  forPage,
	}
	if forPage {
		p.colMap = traitPageColMap
	} else {
		p.colMap = traitListColMap
	}
	return p
}

func (p *traitsProvider) SetTable(table *unison.Table[*Node[*gurps.Trait]]) {
	p.table = table
}

func (p *traitsProvider) RootRowCount() int {
	return len(p.provider.TraitList())
}

func (p *traitsProvider) RootRows() []*Node[*gurps.Trait] {
	data := p.provider.TraitList()
	rows := make([]*Node[*gurps.Trait], 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode[*gurps.Trait](p.table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *traitsProvider) SetRootRows(rows []*Node[*gurps.Trait]) {
	p.provider.SetTraitList(ExtractNodeDataFromList(rows))
}

func (p *traitsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *traitsProvider) DragKey() string {
	return gid.TraitModifier
}

func (p *traitsProvider) DragSVG() *unison.SVG {
	return res.GCSTraitsSVG
}

func (p *traitsProvider) DropShouldMoveData(from, to *unison.Table[*Node[*gurps.Trait]]) bool {
	return from == to
}

func (p *traitsProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Trait"), i18n.Text("Traits")
}

func (p *traitsProvider) Headers() []unison.TableColumnHeader[*Node[*gurps.Trait]] {
	var headers []unison.TableColumnHeader[*Node[*gurps.Trait]]
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.TraitDescriptionColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Trait"), "", p.forPage))
		case gurps.TraitPointsColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		case gurps.TraitTagsColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Tags"), "", p.forPage))
		case gurps.TraitReferenceColumn:
			headers = append(headers, NewPageRefHeader[*gurps.Trait](p.forPage))
		default:
			jot.Fatalf(1, "invalid trait column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *traitsProvider) SyncHeader(_ []unison.TableColumnHeader[*Node[*gurps.Trait]]) {
}

func (p *traitsProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.TraitDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *traitsProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *traitsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Trait]]) {
	OpenEditor[*gurps.Trait](table, func(item *gurps.Trait) { EditTrait(owner, item) })
}

func (p *traitsProvider) CreateItem(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Trait]], variant widget.ItemVariant) {
	item := gurps.NewTrait(p.Entity(), nil, variant == widget.ContainerItemVariant)
	InsertItem[*gurps.Trait](owner, table, item, p.provider.TraitList, p.provider.SetTraitList,
		func(_ *unison.Table[*Node[*gurps.Trait]]) []*Node[*gurps.Trait] { return p.RootRows() })
	EditTrait(owner, item)
}

func (p *traitsProvider) DuplicateSelection(table *unison.Table[*Node[*gurps.Trait]]) {
	duplicateTableSelection(table, p.provider.TraitList(),
		func(nodes []*gurps.Trait) { p.provider.SetTraitList(nodes) },
		func(node *gurps.Trait) *[]*gurps.Trait { return &node.Children })
}

func (p *traitsProvider) DeleteSelection(table *unison.Table[*Node[*gurps.Trait]]) {
	deleteTableSelection(table, p.provider.TraitList(),
		func(nodes []*gurps.Trait) { p.provider.SetTraitList(nodes) },
		func(node *gurps.Trait) *[]*gurps.Trait { return &node.Children })
}

func (p *traitsProvider) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	if err := json.NewEncoder(gz).Encode(p.provider.TraitList()); err != nil {
		return nil, errs.Wrap(err)
	}
	if err := gz.Close(); err != nil {
		return nil, errs.Wrap(err)
	}
	return buffer.Bytes(), nil
}

func (p *traitsProvider) Deserialize(data []byte) error {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return errs.Wrap(err)
	}
	var rows []*gurps.Trait
	if err = json.NewDecoder(gz).Decode(&rows); err != nil {
		return errs.Wrap(err)
	}
	p.provider.SetTraitList(rows)
	return nil
}
