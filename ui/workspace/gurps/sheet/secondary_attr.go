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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// SecondaryAttrPanel holds the contents of the secondary attributes block on the sheet.
type SecondaryAttrPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewSecondaryAttrPanel creates a new secondary attributes panel.
func NewSecondaryAttrPanel(entity *gurps.Entity) *SecondaryAttrPanel {
	p := &SecondaryAttrPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Secondary Attributes")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	return p
}
