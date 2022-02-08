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
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/settings/display"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

// SettingsProvider must be initialized prior to using this package. It provides access to settings to avoid circular
// references.
var SettingsProvider interface {
	GeneralSettings() *settings.General
	SheetSettings() *SheetSettings
	Libraries() library.Libraries
}

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

// SheetSettingsData holds the SheetSettings data that is written to disk.
type SheetSettingsData struct {
	Page                       *settings.Page              `json:"page,omitempty"`
	BlockLayout                []string                    `json:"block_layout,omitempty"`
	Attributes                 *AttributeDefs              `json:"attributes,omitempty"`
	HitLocations               *BodyType                   `json:"hit_locations,omitempty"`
	DamageProgression          attribute.DamageProgression `json:"damage_progression"`
	DefaultLengthUnits         measure.LengthUnits         `json:"default_length_units"`
	DefaultWeightUnits         measure.WeightUnits         `json:"default_weight_units"`
	UserDescriptionDisplay     display.Option              `json:"user_description_display"`
	ModifiersDisplay           display.Option              `json:"modifiers_display"`
	NotesDisplay               display.Option              `json:"notes_display"`
	SkillLevelAdjDisplay       display.Option              `json:"skill_level_adj_display"`
	UseMultiplicativeModifiers bool                        `json:"use_multiplicative_modifiers,omitempty"`
	UseModifyingDicePlusAdds   bool                        `json:"use_modifying_dice_plus_adds,omitempty"`
	ShowCollegeInSheetSpells   bool                        `json:"show_college_in_sheet_spells,omitempty"`
	ShowDifficulty             bool                        `json:"show_difficulty,omitempty"`
	ShowAdvantageModifierAdj   bool                        `json:"show_advantage_modifier_adj,omitempty"`
	ShowEquipmentModifierAdj   bool                        `json:"show_equipment_modifier_adj,omitempty"`
	ShowSpellAdj               bool                        `json:"show_spell_adj,omitempty"`
	UseTitleInFooter           bool                        `json:"use_title_in_footer,omitempty"`
}

// SheetSettings holds sheet settings.
type SheetSettings struct {
	SheetSettingsData
	Entity *Entity `json:"-"`
}

// SheetSettingsFor returns the SheetSettings for the given Entity, or the global settings if the Entity is nil.
func SheetSettingsFor(entity *Entity) *SheetSettings {
	if entity == nil {
		if SettingsProvider == nil {
			jot.Fatal(1, errs.New("SettingsProvider has not been set yet"))
		}
		return SettingsProvider.SheetSettings()
	}
	return entity.SheetSettings
}

// FactorySheetSettings returns a new SheetSettings with factory defaults.
func FactorySheetSettings(entity *Entity) *SheetSettings {
	return &SheetSettings{
		SheetSettingsData: SheetSettingsData{
			Page:                   settings.NewPage(),
			BlockLayout:            FactoryBlockLayout(),
			Attributes:             FactoryAttributeDefs(),
			HitLocations:           FactoryBodyType(),
			DamageProgression:      attribute.BasicSet,
			DefaultLengthUnits:     measure.FeetAndInches,
			DefaultWeightUnits:     measure.Pound,
			UserDescriptionDisplay: display.Tooltip,
			ModifiersDisplay:       display.Inline,
			NotesDisplay:           display.Inline,
			SkillLevelAdjDisplay:   display.Tooltip,
			ShowSpellAdj:           true,
		},
		Entity: entity,
	}
}

// FactoryBlockLayout returns the factory block layout setting.
func FactoryBlockLayout() []string {
	return []string{
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

// MarshalJSON implements json.Marshaler.
func (s *SheetSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(&s.SheetSettingsData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SheetSettings) UnmarshalJSON(data []byte) error {
	s.SheetSettingsData = SheetSettingsData{}
	if err := json.Unmarshal(data, &s.SheetSettingsData); err != nil {
		return err
	}
	return nil
}

// Clone creates a copy of this.
func (s *SheetSettings) Clone(entity *Entity) *SheetSettings {
	clone := *s
	clone.Page = s.Page.Clone()
	clone.BlockLayout = make([]string, len(s.BlockLayout))
	copy(clone.BlockLayout, s.BlockLayout)
	clone.Attributes = s.Attributes.Clone()
	clone.HitLocations = s.HitLocations.Clone(entity, nil)
	return &clone
}
