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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	meleeWeaponColMap = map[int]int{
		0: gurps.WeaponDescriptionColumn,
		1: gurps.WeaponUsageColumn,
		2: gurps.WeaponSLColumn,
		3: gurps.WeaponParryColumn,
		4: gurps.WeaponBlockColumn,
		5: gurps.WeaponDamageColumn,
		6: gurps.WeaponReachColumn,
		7: gurps.WeaponSTColumn,
	}
	rangedWeaponColMap = map[int]int{
		0:  gurps.WeaponDescriptionColumn,
		1:  gurps.WeaponUsageColumn,
		2:  gurps.WeaponSLColumn,
		3:  gurps.WeaponAccColumn,
		4:  gurps.WeaponDamageColumn,
		5:  gurps.WeaponRangeColumn,
		6:  gurps.WeaponRoFColumn,
		7:  gurps.WeaponShotsColumn,
		8:  gurps.WeaponBulkColumn,
		9:  gurps.WeaponRecoilColumn,
		10: gurps.WeaponSTColumn,
	}
)

type weaponsProvider struct {
	colMap     map[int]int
	provider   gurps.WeaponListProvider
	weaponType weapon.Type
	forPage    bool
}

// NewWeaponsProvider creates a new table provider for weapons.
func NewWeaponsProvider(provider gurps.WeaponListProvider, weaponType weapon.Type, forPage bool) TableProvider {
	p := &weaponsProvider{
		provider:   provider,
		weaponType: weaponType,
		forPage:    forPage,
	}
	if weaponType == weapon.Melee {
		p.colMap = meleeWeaponColMap
	} else {
		p.colMap = rangedWeaponColMap
	}
	return p
}

func (p *weaponsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *weaponsProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.WeaponDescriptionColumn:
			headers = append(headers, NewHeader(p.weaponType.String(), "", p.forPage))
		case gurps.WeaponUsageColumn:
			headers = append(headers, NewHeader(i18n.Text("Usage"), "", p.forPage))
		case gurps.WeaponSLColumn:
			headers = append(headers, NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), p.forPage))
		case gurps.WeaponParryColumn:
			headers = append(headers, NewHeader(i18n.Text("Parry"), "", p.forPage))
		case gurps.WeaponBlockColumn:
			headers = append(headers, NewHeader(i18n.Text("Block"), "", p.forPage))
		case gurps.WeaponDamageColumn:
			headers = append(headers, NewHeader(i18n.Text("Damage"), "", p.forPage))
		case gurps.WeaponReachColumn:
			headers = append(headers, NewHeader(i18n.Text("Reach"), "", p.forPage))
		case gurps.WeaponSTColumn:
			headers = append(headers, NewHeader(i18n.Text("ST"), i18n.Text("Minimum Strength"), p.forPage))
		case gurps.WeaponAccColumn:
			headers = append(headers, NewHeader(i18n.Text("Acc"), i18n.Text("Accuracy Bonus"), p.forPage))
		case gurps.WeaponRangeColumn:
			headers = append(headers, NewHeader(i18n.Text("Range"), "", p.forPage))
		case gurps.WeaponRoFColumn:
			headers = append(headers, NewHeader(i18n.Text("RoF"), i18n.Text("Rate of Fire"), p.forPage))
		case gurps.WeaponShotsColumn:
			headers = append(headers, NewHeader(i18n.Text("Shots"), "", p.forPage))
		case gurps.WeaponBulkColumn:
			headers = append(headers, NewHeader(i18n.Text("Bulk"), "", p.forPage))
		case gurps.WeaponRecoilColumn:
			headers = append(headers, NewHeader(i18n.Text("Recoil"), "", p.forPage))
		default:
			jot.Fatalf(1, "invalid weapon column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *weaponsProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.Weapons(p.weaponType)
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *weaponsProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *weaponsProvider) HierarchyColumnIndex() int {
	return -1
}

func (p *weaponsProvider) ExcessWidthColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.WeaponDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *weaponsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	if !p.forPage {
		OpenEditor[*gurps.Weapon](table, func(item *gurps.Weapon) { EditWeapon(owner, item) })
	}
}

func (p *weaponsProvider) CreateItem(_ widget.Rebuildable, _ *unison.Table, _ ItemVariant) {
}

func (p *weaponsProvider) DeleteSelection(table *unison.Table) {
	if !p.forPage {
		deleteTableSelection(table, p.provider.Weapons(p.weaponType),
			func(nodes []*gurps.Weapon) { p.provider.SetWeapons(p.weaponType, nodes) },
			func(node *gurps.Weapon) **gurps.Weapon {
				var dummy *gurps.Weapon
				return &dummy
			},
			func(node *gurps.Weapon) *[]*gurps.Weapon {
				var dummy []*gurps.Weapon
				return &dummy
			})
	}
}
