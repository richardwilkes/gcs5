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
	enum2 "github.com/richardwilkes/gcs/model/enum"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/gcs/model/gurps/measure"
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

// GlobalSheetSettingsProvider must be initialized prior to using this package. It provides access to the global sheet
// settings that should be used when an Entity is not available to provide them.
var GlobalSheetSettingsProvider func() *SheetSettings

// SheetSettings holds sheet settings.
type SheetSettings struct {
	Page                       *PageSettings
	BlockLayout                *BlockLayout
	Attributes                 *AttributeDefs
	HitLocations               *BodyType
	DamageProgression          enum.DamageProgression
	DefaultLengthUnits         measure.LengthUnits
	DefaultWeightUnits         measure.WeightUnits
	UserDescriptionDisplay     enum2.DisplayOption
	ModifiersDisplay           enum2.DisplayOption
	NotesDisplay               enum2.DisplayOption
	SkillLevelAdjDisplay       enum2.DisplayOption
	UseMultiplicativeModifiers bool
	UseModifyingDicePlusAdds   bool
	ShowCollegeInSheetSpells   bool
	ShowDifficulty             bool
	ShowAdvantageModifierAdj   bool
	ShowEquipmentModifierAdj   bool
	ShowSpellAdj               bool
	UseTitleInFooter           bool
}

// SheetSettingsFor returns the SheetSettings for the given Entity, or the global settings if the Entity is nil.
func SheetSettingsFor(entity *Entity) *SheetSettings {
	if entity == nil {
		return GlobalSheetSettingsProvider()
	}
	return entity.SheetSettings
}

// FactorySheetSettings returns a new SheetSettings with factory defaults.
func FactorySheetSettings() *SheetSettings {
	return &SheetSettings{
		Page:                   FactoryPageSettings(),
		BlockLayout:            FactoryBlockLayout(),
		Attributes:             FactoryAttributeDefs(),
		HitLocations:           FactoryBodyType(),
		DamageProgression:      enum.BasicSet,
		DefaultLengthUnits:     measure.FeetAndInches,
		DefaultWeightUnits:     measure.Pound,
		UserDescriptionDisplay: enum2.Tooltip,
		ModifiersDisplay:       enum2.Inline,
		NotesDisplay:           enum2.Inline,
		SkillLevelAdjDisplay:   enum2.Tooltip,
		ShowSpellAdj:           true,
	}
}

// NewSheetSettingsFromJSON creates a new SheetSettings from a JSON object.
func NewSheetSettingsFromJSON(data map[string]interface{}, entity *Entity) *SheetSettings {
	s := FactorySheetSettings()
	s.Page = NewPageSettingsFromJSON(encoding.Object(data[sheetSettingsPageKey]))
	s.BlockLayout = NewBlockLayoutFromJSON(encoding.Object(data[sheetSettingsBlockLayoutKey]))
	if entity != nil {
		s.Attributes = NewAttributeDefsFromJSON(encoding.Array(data[sheetSettingsAttributesKey]))
		s.HitLocations = NewBodyTypeFromJSON(encoding.Object(data[sheetSettingsHitLocationsKey]))
	}
	s.DamageProgression = enum.DamageProgressionFromString(encoding.String(data[sheetSettingsDamageProgressionKey]))
	s.DefaultLengthUnits = measure.LengthUnitsFromString(encoding.String(data[sheetSettingsDefaultLengthUnitsKey]))
	s.DefaultWeightUnits = measure.WeightUnitsFromString(encoding.String(data[sheetSettingsDefaultWeightUnitsKey]))
	s.UserDescriptionDisplay = enum2.DisplayOptionFromString(encoding.String(data[sheetSettingsUserDescriptionDisplayKey]), s.UserDescriptionDisplay)
	s.ModifiersDisplay = enum2.DisplayOptionFromString(encoding.String(data[sheetSettingsModifiersDisplayKey]), s.ModifiersDisplay)
	s.NotesDisplay = enum2.DisplayOptionFromString(encoding.String(data[sheetSettingsNotesDisplayKey]), s.NotesDisplay)
	s.SkillLevelAdjDisplay = enum2.DisplayOptionFromString(encoding.String(data[sheetSettingsSkillLevelAdjDisplayKey]), s.SkillLevelAdjDisplay)
	s.UseMultiplicativeModifiers = encoding.Bool(data[sheetSettingsUseMultiplicativeModifiersKey])
	s.UseModifyingDicePlusAdds = encoding.Bool(data[sheetSettingsUseModifyingDicePlusAddsKey])
	s.ShowCollegeInSheetSpells = encoding.Bool(data[sheetSettingsShowCollegeInSheetSpellsKey])
	s.ShowDifficulty = encoding.Bool(data[sheetSettingsShowDifficultyKey])
	s.ShowAdvantageModifierAdj = encoding.Bool(data[sheetSettingsShowAdvantageModifierAdjKey])
	s.ShowEquipmentModifierAdj = encoding.Bool(data[sheetSettingsShowEquipmentModifierAdjKey])
	s.ShowSpellAdj = encoding.Bool(data[sheetSettingsShowSpellAdjKey])
	s.UseTitleInFooter = encoding.Bool(data[sheetSettingsUseTitleInFooterKey])
	return s
}

// ToJSON emits this object as JSON.
func (s *SheetSettings) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	encoding.ToKeyedJSON(s.Page, sheetSettingsPageKey, encoder)
	encoding.ToKeyedJSON(s.BlockLayout, sheetSettingsBlockLayoutKey, encoder)
	if entity != nil {
		encoding.ToKeyedJSON(s.Attributes, sheetSettingsAttributesKey, encoder)
		ToKeyedJSON(s.HitLocations, sheetSettingsHitLocationsKey, encoder, nil)
	}
	encoder.KeyedString(sheetSettingsDamageProgressionKey, s.DamageProgression.Key(), false, false)
	encoder.KeyedString(sheetSettingsDefaultLengthUnitsKey, s.DefaultLengthUnits.Key(), false, false)
	encoder.KeyedString(sheetSettingsDefaultWeightUnitsKey, s.DefaultWeightUnits.Key(), false, false)
	encoder.KeyedString(sheetSettingsUserDescriptionDisplayKey, s.UserDescriptionDisplay.Key(), false, false)
	encoder.KeyedString(sheetSettingsModifiersDisplayKey, s.ModifiersDisplay.Key(), false, false)
	encoder.KeyedString(sheetSettingsNotesDisplayKey, s.NotesDisplay.Key(), false, false)
	encoder.KeyedString(sheetSettingsSkillLevelAdjDisplayKey, s.SkillLevelAdjDisplay.Key(), false, false)
	encoder.KeyedBool(sheetSettingsUseMultiplicativeModifiersKey, s.UseMultiplicativeModifiers, true)
	encoder.KeyedBool(sheetSettingsUseModifyingDicePlusAddsKey, s.UseModifyingDicePlusAdds, true)
	encoder.KeyedBool(sheetSettingsShowCollegeInSheetSpellsKey, s.ShowCollegeInSheetSpells, true)
	encoder.KeyedBool(sheetSettingsShowDifficultyKey, s.ShowDifficulty, true)
	encoder.KeyedBool(sheetSettingsShowAdvantageModifierAdjKey, s.ShowAdvantageModifierAdj, true)
	encoder.KeyedBool(sheetSettingsShowEquipmentModifierAdjKey, s.ShowEquipmentModifierAdj, true)
	encoder.KeyedBool(sheetSettingsShowSpellAdjKey, s.ShowSpellAdj, true)
	encoder.KeyedBool(sheetSettingsUseTitleInFooterKey, s.UseTitleInFooter, true)
	encoder.EndObject()
}
