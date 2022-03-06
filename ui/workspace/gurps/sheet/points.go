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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// PointsPanel holds the contents of the points block on the sheet.
type PointsPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPointsPanel creates a new points panel.
func NewPointsPanel(entity *gurps.Entity) *PointsPanel {
	p := &PointsPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.EndAlignment,
		VAlign: unison.FillAlignment,
		VSpan:  2,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: fmt.Sprintf(i18n.Text("%s PointsPanel"), p.entity.TotalPoints.String())}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	field := widget.NewNumericPageField(entity.UnspentPoints(), fixed.F64d4FromInt(-999999), fixed.F64d4FromInt(999999),
		func(v fixed.F64d4) {
			if v != entity.UnspentPoints() {
				entity.SetUnspentPoints(v)
				MarkModified(p)
			}
		})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	field.Tooltip = unison.NewTooltipWithText(i18n.Text("PointsPanel earned but not yet spent"))
	p.AddChild(field)
	p.AddChild(widget.NewPageLabel(i18n.Text("Unspent")))

	ad, disad, race, quirk := entity.AdvantagePoints()
	p.AddChild(widget.NewNonEditablePageFieldEnd(race.String(), i18n.Text("Total points spent on a racial package")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Race")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(entity.AttributePoints().String(),
		i18n.Text("Total points spent on attributes")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Attributes")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(ad.String(), i18n.Text("Total points spent on advantages")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Advantages")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(disad.String(), i18n.Text("Total points spent on disadvantages")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Disadvantages")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(quirk.String(), i18n.Text("Total points spent on quirks")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Quirks")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(entity.SkillPoints().String(), i18n.Text("Total points spent on skills")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Skills")))

	p.AddChild(widget.NewNonEditablePageFieldEnd(entity.SpellPoints().String(), i18n.Text("Total points spent on spells")))
	p.AddChild(widget.NewPageLabel(i18n.Text("Spells")))

	return p
}
