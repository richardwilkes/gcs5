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
	"strings"

	"github.com/richardwilkes/json"
)

const (
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

var allBlockLayoutKeys = []string{
	blockLayoutReactionsKey,
	blockLayoutConditionalModifiersKey,
	blockLayoutMeleeKey,
	blockLayoutRangedKey,
	blockLayoutAdvantagesKey,
	blockLayoutSkillsKey,
	blockLayoutSpellsKey,
	blockLayoutEquipmentKey,
	blockLayoutOtherEquipmentKey,
	blockLayoutNotesKey,
}

// BlockLayout holds the sheet's block layout.
type BlockLayout struct {
	Layout []string
}

// NewBlockLayout creates a new default BlockLayout.
func NewBlockLayout() *BlockLayout {
	var b BlockLayout
	b.Reset()
	return &b
}

// MarshalJSON implements json.Marshaler.
func (b *BlockLayout) MarshalJSON() ([]byte, error) {
	return json.Marshal(&b.Layout)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockLayout) UnmarshalJSON(data []byte) error {
	b.Layout = nil
	if err := json.Unmarshal(data, &b.Layout); err != nil {
		return err
	}
	if len(b.Layout) == 0 {
		b.Reset()
	}
	return nil
}

// Clone this data.
func (b *BlockLayout) Clone() *BlockLayout {
	clone := *b
	clone.Layout = make([]string, len(b.Layout))
	copy(clone.Layout, b.Layout)
	return &clone
}

// Reset returns the BlockLayout to factory settings.
func (b *BlockLayout) Reset() {
	b.Layout = []string{
		blockLayoutReactionsKey + " " + blockLayoutConditionalModifiersKey,
		blockLayoutMeleeKey,
		blockLayoutRangedKey,
		blockLayoutAdvantagesKey + " " + blockLayoutSkillsKey,
		blockLayoutSpellsKey,
		blockLayoutEquipmentKey,
		blockLayoutOtherEquipmentKey,
		blockLayoutNotesKey,
	}
}

// CreateFullKeySet creates a map that contains each of the possible block layout keys.
func (b *BlockLayout) CreateFullKeySet() map[string]bool {
	return map[string]bool{
		blockLayoutReactionsKey:            true,
		blockLayoutConditionalModifiersKey: true,
		blockLayoutMeleeKey:                true,
		blockLayoutRangedKey:               true,
		blockLayoutAdvantagesKey:           true,
		blockLayoutSkillsKey:               true,
		blockLayoutSpellsKey:               true,
		blockLayoutEquipmentKey:            true,
		blockLayoutOtherEquipmentKey:       true,
		blockLayoutNotesKey:                true,
	}
}

// HTMLGridTemplate returns the text for the HTML grid layout.
func (b *BlockLayout) HTMLGridTemplate() string {
	var buffer strings.Builder
	remaining := b.CreateFullKeySet()
	for _, line := range b.Layout {
		parts := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
		if parts[0] != "" && remaining[parts[0]] {
			delete(remaining, parts[0])
			if len(parts) > 1 && remaining[parts[1]] {
				delete(remaining, parts[1])
				appendToGridTemplate(&buffer, parts[0], parts[1])
				continue
			}
			appendToGridTemplate(&buffer, parts[0], parts[0])
		}
	}
	for _, k := range allBlockLayoutKeys {
		if remaining[k] {
			appendToGridTemplate(&buffer, k, k)
		}
	}
	return buffer.String()
}

func appendToGridTemplate(buffer *strings.Builder, left, right string) {
	buffer.WriteByte('"')
	buffer.WriteString(left)
	buffer.WriteByte(' ')
	buffer.WriteString(right)
	buffer.WriteByte('"')
	buffer.WriteByte('\n')
}
