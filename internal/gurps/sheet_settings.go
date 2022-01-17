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
	"github.com/richardwilkes/gcs/unit/length"
	"github.com/richardwilkes/gcs/unit/weight"
)

// SheetSettings holds sheet settings.
type SheetSettings struct {
	DefaultLengthUnits         length.Units `json:"default_length_units"`
	DefaultWeightUnits         weight.Units `json:"default_weight_units"`
	UserDescriptionDisplay     string       `json:"user_description_display"`
	ModifiersDisplay           string       `json:"modifiers_display"`
	NotesDisplay               string       `json:"notes_display"`
	SkillLevelAdjDisplay       string       `json:"skill_level_adj_display"`
	DamageProgression          string       `json:"damage_progression"`
	UseMultiplicativeModifiers bool         `json:"use_multiplicative_modifiers"`
	UseModifyingDicePlusAdds   bool         `json:"use_modifying_dice_plus_adds"`
	ShowCollegeInSheetSpells   bool         `json:"show_college_in_sheet_spells"`
	ShowDifficulty             bool         `json:"show_difficulty"`
	ShowAdvantageModifierAdj   bool         `json:"show_advantage_modifier_adj"`
	ShowEquipmentModifierAdj   bool         `json:"show_equipment_modifier_adj"`
	ShowSpellAdj               bool         `json:"show_spell_adj"`
	UseTitleInFooter           bool         `json:"use_title_in_footer"`
	Page                       PageSettings `json:"page"`
	BlockLayout                []string     `json:"block_layout"`
	Attributes                 []*Attribute `json:"attributes"`
	HitLocations               *BodyType    `json:"hit_locations"`
}

// PageSettings holds page settings.
type PageSettings struct {
	PaperSize    string        `json:"paper_size"`
	TopMargin    length.Length `json:"top_margin"`
	LeftMargin   length.Length `json:"left_margin"`
	BottomMargin length.Length `json:"bottom_margin"`
	RightMargin  length.Length `json:"right_margin"`
	Orientation  string        `json:"orientation"`
}

// FactorySheetSettings returns a new SheetSettings will factory defaults.
func FactorySheetSettings() *SheetSettings {
	return &SheetSettings{
		DefaultLengthUnits:     length.FeetAndInches,
		DefaultWeightUnits:     weight.Pound,
		UserDescriptionDisplay: "tooltip",   // TODO: Use type
		ModifiersDisplay:       "inline",    // TODO: Use type
		NotesDisplay:           "inline",    // TODO: Use type
		SkillLevelAdjDisplay:   "tooltip",   // TODO: Use type
		DamageProgression:      "basic_set", // TODO: Use type
		ShowSpellAdj:           true,
		Page: PageSettings{
			PaperSize:    "na-letter", // TODO: Use type
			TopMargin:    length.FromFloat64(0.25, length.Inch),
			LeftMargin:   length.FromFloat64(0.25, length.Inch),
			BottomMargin: length.FromFloat64(0.25, length.Inch),
			RightMargin:  length.FromFloat64(0.25, length.Inch),
			Orientation:  "portrait", // TODO: Use type
		},
		BlockLayout:  FactoryBlockLayout(),
		Attributes:   FactoryAttributes(),
		HitLocations: FactoryBodyType(),
	}
}

// FactoryBlockLayout returns the block layout factory settings.
func FactoryBlockLayout() []string {
	return []string{
		// TODO: Use constants
		"reactions conditional_modifiers",
		"melee",
		"ranged",
		"advantages skills",
		"spells",
		"equipment",
		"other_equipment",
		"notes",
	}
}
