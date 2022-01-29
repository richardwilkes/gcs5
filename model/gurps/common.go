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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/id"
)

const (
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
	ID        uuid.UUID
	Name      string
	PageRef   string
	Notes     string
	VTTNotes  string
	Container bool
	Open      bool
}

// FromJSON loads common data from a JSON object.
func (c *Common) FromJSON(typeKey string, data map[string]interface{}) {
	c.Container = encoding.String(data[commonTypeKey]) == typeKey+commonContainerKeyPostfix
	c.ID = id.ParseOrNewUUID(encoding.String(data[commonIDKey]))
	c.Name = encoding.String(data[commonNameKey])
	c.PageRef = encoding.String(data[commonPageRefKey])
	c.Notes = encoding.String(data[commonNotesKey])
	c.VTTNotes = encoding.String(data[commonVTTNotesKey])
	if c.Container {
		c.Open = encoding.Bool(data[commonOpenKey])
	}
}

// ToInlineJSON emits this object as JSON.
func (c *Common) ToInlineJSON(typeKey string, encoder *encoding.JSONEncoder) {
	typeString := typeKey
	if c.Container {
		typeString += commonContainerKeyPostfix
	}
	encoder.KeyedString(commonTypeKey, typeString, false, false)
	encoder.KeyedString(commonIDKey, c.ID.String(), false, false)
	encoder.KeyedString(commonNameKey, c.Name, true, true)
	encoder.KeyedString(commonPageRefKey, c.PageRef, true, true)
	encoder.KeyedString(commonNotesKey, c.Notes, true, true)
	encoder.KeyedString(commonVTTNotesKey, c.VTTNotes, true, true)
	if c.Container {
		encoder.KeyedBool(commonOpenKey, c.Open, true)
	}
}

// FeaturesListFromJSON loads a features list from a JSON object.
func FeaturesListFromJSON(data map[string]interface{}) []*Feature {
	array := encoding.Array(data[commonFeaturesKey])
	if len(array) == 0 {
		return nil
	}
	features := make([]*Feature, len(array))
	for i, one := range array {
		features[i] = NewFeatureFromJSON(encoding.Object(one))
	}
	return features
}

// FeaturesListToJSON emits the features list as JSON.
func FeaturesListToJSON(features []*Feature, encoder *encoding.JSONEncoder) {
	if len(features) != 0 {
		encoder.Key(commonFeaturesKey)
		encoder.StartArray()
		for _, one := range features {
			one.ToJSON(encoder)
		}
		encoder.EndArray()
	}
}

// WeaponsListFromJSON loads a weapons list from a JSON object.
func WeaponsListFromJSON(data map[string]interface{}) []*Weapon {
	array := encoding.Array(data[commonWeaponsKey])
	if len(array) == 0 {
		return nil
	}
	weapons := make([]*Weapon, len(array))
	for i, one := range array {
		weapons[i] = NewWeaponFromJSON(encoding.Object(one))
	}
	return weapons
}

// WeaponsListToJSON emits the weapons list as JSON.
func WeaponsListToJSON(weapons []*Weapon, encoder *encoding.JSONEncoder) {
	if len(weapons) != 0 {
		encoder.Key(commonWeaponsKey)
		encoder.StartArray()
		for _, one := range weapons {
			one.ToJSON(encoder)
		}
		encoder.EndArray()
	}
}

// SkillDefaultsListFromJSON loads a SkillDefault list from a JSON object.
func SkillDefaultsListFromJSON(data map[string]interface{}) []*SkillDefault {
	array := encoding.Array(data[commonSkillDefaultsKey])
	if len(array) == 0 {
		return nil
	}
	skillDefaults := make([]*SkillDefault, len(array))
	for i, one := range array {
		skillDefaults[i] = NewSkillDefaultFromJSON(false, encoding.Object(one))
	}
	return skillDefaults
}

// SkillDefaultsListToJSON emits the weapons list as JSON.
func SkillDefaultsListToJSON(skillDefaults []*SkillDefault, encoder *encoding.JSONEncoder) {
	if len(skillDefaults) != 0 {
		encoder.Key(commonSkillDefaultsKey)
		encoder.StartArray()
		for _, one := range skillDefaults {
			one.ToJSON(false, encoder)
		}
		encoder.EndArray()
	}
}

// AdvantageModifiersListFromJSON loads an advantage modifiers list from a JSON object.
func AdvantageModifiersListFromJSON(key string, data map[string]interface{}) []*AdvantageModifier {
	array := encoding.Array(data[key])
	if len(array) == 0 {
		return nil
	}
	modifiers := make([]*AdvantageModifier, len(array))
	for i, one := range array {
		modifiers[i] = NewAdvantageModifierFromJSON(encoding.Object(one))
	}
	return modifiers
}

// AdvantageModifiersListToJSON emits the advantage modifiers list as JSON.
func AdvantageModifiersListToJSON(key string, modifiers []*AdvantageModifier, encoder *encoding.JSONEncoder) {
	if len(modifiers) != 0 {
		encoder.Key(key)
		encoder.StartArray()
		for _, one := range modifiers {
			one.ToJSON(encoder)
		}
		encoder.EndArray()
	}
}

// StringListFromJSON loads a string list from a JSON object.
func StringListFromJSON(key string, omitEmptyEntries bool, data map[string]interface{}) []string {
	array := encoding.Array(data[key])
	if len(array) == 0 {
		return nil
	}
	list := make([]string, 0, len(array))
	for _, one := range array {
		s := encoding.String(one)
		if omitEmptyEntries && s == "" {
			continue
		}
		list = append(list, s)
	}
	return list
}

// StringListToJSON emits the string list as JSON.
func StringListToJSON(key string, list []string, encoder *encoding.JSONEncoder) {
	if len(list) != 0 {
		encoder.Key(key)
		encoder.StartArray()
		for _, one := range list {
			encoder.String(one)
		}
		encoder.EndArray()
	}
}
