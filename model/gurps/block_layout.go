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
	"github.com/richardwilkes/toolbox/txt"
)

// Valid block layout keys
const (
	BlockLayoutReactionsKey            = "reactions"
	BlockLayoutConditionalModifiersKey = "conditional_modifiers"
	BlockLayoutMeleeKey                = "melee"
	BlockLayoutRangedKey               = "ranged"
	BlockLayoutTraitsKey               = "advantages"
	BlockLayoutSkillsKey               = "skills"
	BlockLayoutSpellsKey               = "spells"
	BlockLayoutEquipmentKey            = "equipment"
	BlockLayoutOtherEquipmentKey       = "other_equipment"
	BlockLayoutNotesKey                = "notes"
)

var allBlockLayoutKeys = []string{
	BlockLayoutReactionsKey,
	BlockLayoutConditionalModifiersKey,
	BlockLayoutMeleeKey,
	BlockLayoutRangedKey,
	BlockLayoutTraitsKey,
	BlockLayoutSkillsKey,
	BlockLayoutSpellsKey,
	BlockLayoutEquipmentKey,
	BlockLayoutOtherEquipmentKey,
	BlockLayoutNotesKey,
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

// NewBlockLayoutFromString creates a new BlockLayout from an input string.
func NewBlockLayoutFromString(str string) (blockLayout *BlockLayout, inputWasValid bool) {
	var layout []string
	remaining := CreateFullKeySet()
	inputWasValid = true
	for _, line := range strings.Split(strings.ToLower(str), "\n") {
		var parts []string
		for _, part := range strings.Split(txt.CollapseSpaces(line), " ") {
			if part == "" {
				continue
			}
			if len(parts) > 1 {
				inputWasValid = false
				break
			}
			if remaining[part] {
				delete(remaining, part)
				parts = append(parts, part)
			} else {
				inputWasValid = false
			}
		}
		if len(parts) != 0 {
			layout = append(layout, strings.Join(parts, " "))
		}
	}
	if len(remaining) != 0 {
		for _, k := range allBlockLayoutKeys {
			if remaining[k] {
				layout = append(layout, k)
			}
		}
	}
	return &BlockLayout{Layout: layout}, inputWasValid
}

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (b *BlockLayout) EnsureValidity() {
	var layout []string
	remaining := CreateFullKeySet()
	for _, line := range b.Layout {
		var parts []string
		for _, part := range strings.Split(strings.ToLower(txt.CollapseSpaces(line)), " ") {
			if remaining[part] {
				delete(remaining, part)
				parts = append(parts, part)
				if len(parts) > 1 {
					break
				}
			}
		}
		if len(parts) != 0 {
			layout = append(layout, strings.Join(parts, " "))
		}
	}
	if len(remaining) != 0 {
		for _, k := range allBlockLayoutKeys {
			if remaining[k] {
				layout = append(layout, k)
			}
		}
	}
	b.Layout = layout
}

// ByRow breaks the layout down into rows.
func (b *BlockLayout) ByRow() [][]string {
	var layout [][]string
	remaining := CreateFullKeySet()
	for _, line := range b.Layout {
		var parts []string
		for _, part := range strings.Split(strings.ToLower(txt.CollapseSpaces(line)), " ") {
			if remaining[part] {
				delete(remaining, part)
				parts = append(parts, part)
			}
		}
		if len(parts) != 0 {
			layout = append(layout, parts)
		}
	}
	if len(remaining) != 0 {
		for _, k := range allBlockLayoutKeys {
			if remaining[k] {
				layout = append(layout, []string{k})
			}
		}
	}
	return layout
}

func (b *BlockLayout) String() string {
	var buffer strings.Builder
	for _, row := range b.ByRow() {
		buffer.WriteString(strings.Join(row, " "))
		buffer.WriteByte('\n')
	}
	return strings.TrimSpace(buffer.String())
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
		BlockLayoutReactionsKey + " " + BlockLayoutConditionalModifiersKey,
		BlockLayoutMeleeKey,
		BlockLayoutRangedKey,
		BlockLayoutTraitsKey + " " + BlockLayoutSkillsKey,
		BlockLayoutSpellsKey,
		BlockLayoutEquipmentKey,
		BlockLayoutOtherEquipmentKey,
		BlockLayoutNotesKey,
	}
}

// CreateFullKeySet creates a map that contains each of the possible block layout keys.
func CreateFullKeySet() map[string]bool {
	m := make(map[string]bool)
	for _, one := range allBlockLayoutKeys {
		m[one] = true
	}
	return m
}

// HTMLGridTemplate returns the text for the HTML grid layout.
func (b *BlockLayout) HTMLGridTemplate() string {
	var buffer strings.Builder
	remaining := CreateFullKeySet()
	for _, line := range b.Layout {
		parts := strings.Split(strings.ToLower(txt.CollapseSpaces(line)), " ")
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
