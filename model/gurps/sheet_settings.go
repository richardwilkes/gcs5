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
	"github.com/richardwilkes/gcs/model/unit/length"
	"github.com/richardwilkes/gcs/model/unit/weight"
)

const (
	sheetSettingsDefaultLengthUnitsKey         = "default_length_units"
	sheetSettingsDefaultWeightUnitsKey         = "default_weight_units"
	sheetSettingsUserDescriptionDisplayKey     = "user_description_display"
	sheetSettingsModifiersDisplayKey           = "modifiers_display"
	sheetSettingsNotesDisplayKey               = "notes_display"
	sheetSettingsSkillLevelAdjDisplayKey       = "skill_level_adj_display"
	sheetSettingsDamageProgressionKey          = "damage_progression"
	sheetSettingsUseMultiplicativeModifiersKey = "use_multiplicative_modifiers"
	sheetSettingsUseModifyingDicePlusAddsKey   = "use_modifying_dice_plus_adds"
	sheetSettingsShowCollegeInSheetSpellsKey   = "show_college_in_sheet_spells"
	sheetSettingsShowDifficultyKey             = "show_difficulty"
	sheetSettingsShowAdvantageModifierAdjKey   = "show_advantage_modifier_adj"
	sheetSettingsShowEquipmentModifierAdjKey   = "show_equipment_modifier_adj"
	sheetSettingsShowSpellAdjKey               = "show_spell_adj"
	sheetSettingsUseTitleInFooterKey           = "use_title_in_footer"
	sheetSettingsPageKey                       = "page"
	sheetSettingsBlockLayoutKey                = "block_layout"
	sheetSettingsAttributesKey                 = "attributes"
	sheetSettingsHitLocationsKey               = "hit_locations"
)

// SheetSettings holds sheet settings.
type SheetSettings struct {
	DefaultLengthUnits         length.Units
	DefaultWeightUnits         weight.Units
	UserDescriptionDisplay     string
	ModifiersDisplay           string
	NotesDisplay               string
	SkillLevelAdjDisplay       string
	DamageProgression          DamageProgression
	Page                       PageSettings
	BlockLayout                []string
	Attributes                 map[string]*AttributeDef
	HitLocations               *BodyType
	UseMultiplicativeModifiers bool
	UseModifyingDicePlusAdds   bool
	ShowCollegeInSheetSpells   bool
	ShowDifficulty             bool
	ShowAdvantageModifierAdj   bool
	ShowEquipmentModifierAdj   bool
	ShowSpellAdj               bool
	UseTitleInFooter           bool
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
		UserDescriptionDisplay: "tooltip", // TODO: Use type
		ModifiersDisplay:       "inline",  // TODO: Use type
		NotesDisplay:           "inline",  // TODO: Use type
		SkillLevelAdjDisplay:   "tooltip", // TODO: Use type
		DamageProgression:      BasicSet,
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
		Attributes:   FactoryAttributeDefs(),
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

// NewSheetSettingsFromJSON creates a new SheetSettings from a JSON object.
func NewSheetSettingsFromJSON(data map[string]interface{}) *SheetSettings {
	s := FactorySheetSettings()
	// TODO: Implement
	return s
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (s *SheetSettings) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	s.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (s *SheetSettings) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	// TODO: Implement
	encoder.EndObject()
}
