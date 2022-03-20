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

// NewWeaponTableHeaders creates a new set of table column headers for weapons.
func NewWeaponTableHeaders(melee bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	if melee {
		headers = append(headers, NewHeader(i18n.Text("Melee Weapon"), "", true))
	} else {
		headers = append(headers, NewHeader(i18n.Text("Ranged Weapon"), "", true))
	}
	headers = append(headers,
		NewHeader(i18n.Text("Usage"), "", true),
		NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), true),
	)
	if melee {
		headers = append(headers,
			NewHeader(i18n.Text("Parry"), "", true),
			NewHeader(i18n.Text("Block"), "", true),
			NewHeader(i18n.Text("Damage"), "", true),
			NewHeader(i18n.Text("Reach"), "", true),
		)
	} else {
		headers = append(headers,
			NewHeader(i18n.Text("Acc"), i18n.Text("Accuracy Bonus"), true),
			NewHeader(i18n.Text("Damage"), "", true),
			NewHeader(i18n.Text("Range"), "", true),
			NewHeader(i18n.Text("RoF"), i18n.Text("Rate of Fire"), true),
			NewHeader(i18n.Text("Shots"), "", true),
			NewHeader(i18n.Text("Bulk"), "", true),
			NewHeader(i18n.Text("Recoil"), "", true),
		)
	}
	return append(headers, NewHeader(i18n.Text("ST"), i18n.Text("Minimum Strength"), true))
}

// NewWeaponRowData creates a new table data provider function for weapons.
func NewWeaponRowData(topLevelRowProvider func() []*gurps.Weapon, melee bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		var colMap map[int]int
		if melee {
			colMap = meleeWeaponColMap
		} else {
			colMap = rangedWeaponColMap
		}
		data := topLevelRowProvider()
		rows := make([]unison.TableRowData, 0, len(data))
		for _, one := range data {
			rows = append(rows, NewNode(table, nil, colMap, one, true))
		}
		return rows
	}
}
