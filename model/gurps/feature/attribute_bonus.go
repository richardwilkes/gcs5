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

package feature

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/toolbox/xio"
)

var _ Bonus = &AttributeBonus{}

// AttributeBonus holds the data for a bonus to an attribute.
type AttributeBonus struct {
	Type       Type                      `json:"type"`
	Limitation attribute.BonusLimitation `json:"limitation,omitempty"`
	Parent     fmt.Stringer              `json:"-"`
	Attribute  string                    `json:"attribute"`
	LeveledAmount
}

// NewAttributeBonus creates a new AttributeBonus.
func NewAttributeBonus(attrID string) *AttributeBonus {
	return &AttributeBonus{
		Type:          AttributeBonusType,
		Attribute:     attrID,
		Limitation:    attribute.None,
		LeveledAmount: LeveledAmount{Amount: fxp.One},
	}
}

// FeatureMapKey implements Feature.
func (a *AttributeBonus) FeatureMapKey() string {
	key := AttributeIDPrefix + a.Attribute
	if a.Limitation != attribute.None {
		key += "." + string(a.Limitation)
	}
	return key
}

// FillWithNameableKeys implements Feature.
func (a *AttributeBonus) FillWithNameableKeys(_ map[string]string) {
}

// ApplyNameableKeys implements Feature.
func (a *AttributeBonus) ApplyNameableKeys(_ map[string]string) {
}

// AddToTooltip implements Bonus.
func (a *AttributeBonus) AddToTooltip(buffer *xio.ByteBuffer) {
	basicAddToTooltip(a.Parent, &a.LeveledAmount, buffer)
}
