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

import (
	"github.com/richardwilkes/gcs/model/encoding"
)

const (
	blockLayoutLayoutKey               = "layout"
	blockLayoutReactionsKey            = "reactions"
	blockLayoutConditionalModifiersKey = "conditional_modifiers"
	blockLayoutMeleeKey                = "melee"
	blockLayoutRangedKey               = "ranged"
	blockLayoutAdvantagesKey           = "advantages"
	blockLayoutSkillsKey               = "skills"
	blockLayoutSpellsKey               = "spells"
	blockLayoutEquipmentKey            = "equipment"
	blockLayoutOtherEquipmentKey       = "other_equipment"
	blockLayoutNotesKey                = "notes"
)

// BlockLayout holds the layout arrangement for the top-level blocks on a character sheet.
type BlockLayout struct {
	Layout []string
}

// FactoryBlockLayout returns a new BlockLayout with factory defaults.
func FactoryBlockLayout() *BlockLayout {
	return &BlockLayout{
		Layout: []string{
			blockLayoutReactionsKey + " " + blockLayoutConditionalModifiersKey,
			blockLayoutMeleeKey,
			blockLayoutRangedKey,
			blockLayoutAdvantagesKey + " " + blockLayoutSkillsKey,
			blockLayoutSpellsKey,
			blockLayoutEquipmentKey,
			blockLayoutOtherEquipmentKey,
			blockLayoutNotesKey,
		},
	}
}

// NewBlockLayoutFromJSON creates a new BlockLayout from a JSON object.
func NewBlockLayoutFromJSON(data map[string]interface{}) *BlockLayout {
	l := FactoryBlockLayout()
	array := encoding.Array(data[blockLayoutLayoutKey])
	list := make([]string, 0, len(array))
	for _, one := range array {
		if str := encoding.String(one); str != "" {
			list = append(list, str)
		}
	}
	if len(list) != 0 {
		l.Layout = list
	}
	return l
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (l *BlockLayout) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	l.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (l *BlockLayout) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.Key(blockLayoutLayoutKey)
	encoder.StartArray()
	for _, one := range l.Layout {
		encoder.String(one)
	}
	encoder.EndArray()
	encoder.EndObject()
}
