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
	"context"
	"io/fs"

	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/library"
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

// SheetSettingsData holds the SheetSettings data that is written to disk.
type SheetSettingsData struct {
	Page                       *settings.Page              `json:"page,omitempty"`
	BlockLayout                *BlockLayout                `json:"block_layout,omitempty"`
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
func FactorySheetSettings() *SheetSettings {
	return &SheetSettings{
		SheetSettingsData: SheetSettingsData{
			Page:                   settings.NewPage(),
			BlockLayout:            NewBlockLayout(),
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
	}
}

// NewSheetSettingsFromFile loads new settings from a file.
func NewSheetSettingsFromFile(fileSystem fs.FS, filePath string) (*SheetSettings, error) {
	var data struct {
		SheetSettings `json:",inline"`
		OldLocation   *SheetSettings `json:"sheet_settings"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, err
	}
	var s *SheetSettings
	if data.OldLocation != nil {
		s = data.OldLocation
	} else {
		ss := data.SheetSettings
		s = &ss
	}
	s.EnsureValidity()
	return s, nil
}

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (s *SheetSettings) EnsureValidity() {
	if s.Page == nil {
		s.Page = settings.NewPage()
	} else {
		s.Page.EnsureValidity()
	}
	if s.BlockLayout == nil {
		s.BlockLayout = NewBlockLayout()
	} else {
		s.BlockLayout.EnsureValidity()
	}
	if s.Attributes == nil {
		s.Attributes = FactoryAttributeDefs()
	} else {
		s.Attributes.EnsureValidity()
	}
	if s.HitLocations == nil {
		s.HitLocations = FactoryBodyType()
	} else {
		s.HitLocations.EnsureValidity()
	}
	s.DamageProgression = s.DamageProgression.EnsureValid()
	s.DefaultLengthUnits = s.DefaultLengthUnits.EnsureValid()
	s.DefaultWeightUnits = s.DefaultWeightUnits.EnsureValid()
	s.UserDescriptionDisplay = s.UserDescriptionDisplay.EnsureValid()
	s.ModifiersDisplay = s.ModifiersDisplay.EnsureValid()
	s.NotesDisplay = s.NotesDisplay.EnsureValid()
	s.SkillLevelAdjDisplay = s.SkillLevelAdjDisplay.EnsureValid()
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
	s.EnsureValidity()
	return nil
}

// Clone creates a copy of this.
func (s *SheetSettings) Clone(entity *Entity) *SheetSettings {
	clone := *s
	clone.Page = s.Page.Clone()
	clone.BlockLayout = s.BlockLayout.Clone()
	clone.Attributes = s.Attributes.Clone()
	clone.HitLocations = s.HitLocations.Clone(entity, nil)
	return &clone
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (s *SheetSettings) SetOwningEntity(entity *Entity) {
	s.Entity = entity
	s.HitLocations.Update(entity)
}

// Save writes the settings to the file as JSON.
func (s *SheetSettings) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, s)
}
