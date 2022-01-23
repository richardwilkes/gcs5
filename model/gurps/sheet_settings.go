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
	"github.com/richardwilkes/gcs/model/enums/display"
	"github.com/richardwilkes/gcs/model/enums/dmg"
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
	Page                       *PageSettings
	BlockLayout                []string
	Attributes                 *AttributeDefs
	HitLocations               *BodyType
	DamageProgression          dmg.Progression
	UserDescriptionDisplay     display.Option
	ModifiersDisplay           display.Option
	NotesDisplay               display.Option
	SkillLevelAdjDisplay       display.Option
	UseMultiplicativeModifiers bool
	UseModifyingDicePlusAdds   bool
	ShowCollegeInSheetSpells   bool
	ShowDifficulty             bool
	ShowAdvantageModifierAdj   bool
	ShowEquipmentModifierAdj   bool
	ShowSpellAdj               bool
	UseTitleInFooter           bool
}

// FactorySheetSettings returns a new SheetSettings with factory defaults.
func FactorySheetSettings() *SheetSettings {
	return &SheetSettings{
		DefaultLengthUnits:     length.FeetAndInches,
		DefaultWeightUnits:     weight.Pound,
		Page:                   FactoryPageSettings(),
		BlockLayout:            FactoryBlockLayout(),
		Attributes:             FactoryAttributeDefs(),
		HitLocations:           FactoryBodyType(),
		DamageProgression:      dmg.BasicSet,
		UserDescriptionDisplay: display.Tooltip,
		ModifiersDisplay:       display.Inline,
		NotesDisplay:           display.Inline,
		SkillLevelAdjDisplay:   display.Tooltip,
		ShowSpellAdj:           true,
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
	encoder.KeyedString(sheetSettingsDefaultLengthUnitsKey, s.DefaultLengthUnits)

	/*
	   w.keyValue(KEY_DEFAULT_LENGTH_UNITS, Enums.toId(mDefaultLengthUnits));
	   w.keyValue(KEY_DEFAULT_WEIGHT_UNITS, Enums.toId(mDefaultWeightUnits));
	   w.keyValue(KEY_USER_DESCRIPTION_DISPLAY, Enums.toId(mUserDescriptionDisplay));
	   w.keyValue(KEY_MODIFIERS_DISPLAY, Enums.toId(mModifiersDisplay));
	   w.keyValue(KEY_NOTES_DISPLAY, Enums.toId(mNotesDisplay));
	   w.keyValue(KEY_SKILL_LEVEL_ADJ_DISPLAY, Enums.toId(mSkillLevelAdjustmentsDisplay));
	   w.keyValue(KEY_USE_MULTIPLICATIVE_MODIFIERS, mUseMultiplicativeModifiers);
	   w.keyValue(KEY_USE_MODIFYING_DICE_PLUS_ADDS, mUseModifyingDicePlusAdds);
	   w.keyValue(KEY_DAMAGE_PROGRESSION, Enums.toId(mDamageProgression));
	   w.keyValue(KEY_USE_SIMPLE_METRIC_CONVERSIONS, mUseSimpleMetricConversions);
	   w.keyValue(KEY_SHOW_COLLEGE_IN_SPELLS, mShowCollegeInSpells);
	   w.keyValue(KEY_SHOW_DIFFICULTY, mShowDifficulty);
	   w.keyValue(KEY_SHOW_ADVANTAGE_MODIFIER_ADJ, mShowAdvantageModifierAdj);
	   w.keyValue(KEY_SHOW_EQUIPMENT_MODIFIER_ADJ, mShowEquipmentModifierAdj);
	   w.keyValue(KEY_SHOW_SPELL_ADJ, mShowSpellAdj);
	   w.keyValue(KEY_USE_TITLE_IN_FOOTER, mUseTitleInFooter);
	   w.key(KEY_PAGE);
	   mPageSettings.toJSON(w);
	   w.key(KEY_BLOCK_LAYOUT);
	   w.startArray();
	   for (String one : mBlockLayout) {
	       w.value(one);
	   }
	   w.endArray();
	   if (full) {
	       w.key(KEY_ATTRIBUTES);
	       AttributeDef.writeOrdered(w, mAttributes);
	       w.key(KEY_HIT_LOCATIONS);
	       mHitLocations.toJSON(w, mCharacter);
	   }
	*/

	encoder.EndObject()
}
