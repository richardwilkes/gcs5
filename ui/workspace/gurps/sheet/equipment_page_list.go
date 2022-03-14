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

package sheet

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// NewCarriedEquipmentPageList creates the carried equipment page list.
func NewCarriedEquipmentPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, equipmentTableColumnHeaders(entity, true))
}

// NewOtherEquipmentPageList creates the other equipment page list.
func NewOtherEquipmentPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, equipmentTableColumnHeaders(entity, false))
}

func equipmentTableColumnHeaders(entity *gurps.Entity, carried bool) []unison.TableColumnHeader {
	var list []unison.TableColumnHeader
	if carried {
		list = append(list, NewEquippedHeader())
	}
	list = append(list, NewPageListHeader(i18n.Text("Qty"), i18n.Text("Quantity")))
	if carried {
		list = append(list, NewPageListHeader(fmt.Sprintf(i18n.Text("Carried Equipment (%s; $%s)"),
			entity.WeightCarried(false).String(), entity.WealthCarried().String()), ""))
	} else {
		list = append(list, NewPageListHeader(fmt.Sprintf(i18n.Text("Other Equipment ($%s)"),
			entity.WealthNotCarried().String()), ""))
	}
	return append(list,
		NewPageListHeader(i18n.Text("Uses"), i18n.Text("The number of uses remaining")),
		NewMoneyHeader(),
		NewWeightHeader(),
		NewExtendedMoneyHeader(),
		NewExtendedWeightHeader(),
		NewPageReferenceHeader(),
	)
}
