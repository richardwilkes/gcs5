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

package gurps

import "C"
import (
	"strings"

	"github.com/google/uuid"
)

const (
	commonCalcKey             = "calc"
	commonCategoriesKey       = "categories"
	commonChildrenKey         = "children"
	commonContainerKeyPostfix = "_container"
	commonDisabledKey         = "disabled"
	commonFeaturesKey         = "features"
	commonIDKey               = "id"
	commonModifiersKey        = "modifiers"
	commonNameKey             = "name"
	commonNotesKey            = "notes"
	commonOpenKey             = "open"
	commonPageRefKey          = "reference"
	commonSkillDefaultsKey    = "defaults"
	commonTypeKey             = "type"
	commonVTTNotesKey         = "vtt_notes"
	commonWeaponsKey          = "weapons"
)

// Common data most of the top-level objects share.
type Common struct {
	Type     string    `json:"type"`
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name,omitempty"`
	PageRef  string    `json:"reference,omitempty"`
	Notes    string    `json:"notes,omitempty"`
	VTTNotes string    `json:"vtt_notes,omitempty"`
	Open     bool      `json:"open,omitempty"`
}

func (c *Common) Container() bool {
	return strings.HasSuffix(c.Type, commonContainerKeyPostfix)
}

// FillWithNameableKeys adds any nameable keys found in this Common to the provided map.
func (c *Common) FillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(c.Name, nameables)
	ExtractNameables(c.Notes, nameables)
	ExtractNameables(c.VTTNotes, nameables)
}

// ApplyNameableKeys replaces any nameable keys found in this Common with the corresponding values in the provided map.
func (c *Common) ApplyNameableKeys(nameables map[string]string) {
	c.Name = ApplyNameables(c.Name, nameables)
	c.Notes = ApplyNameables(c.Notes, nameables)
	c.VTTNotes = ApplyNameables(c.VTTNotes, nameables)
}
