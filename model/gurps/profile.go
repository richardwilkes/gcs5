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
	"encoding/base64"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/enums/units"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/unison"
)

const (
	profilePlayerNameKey   = "player_name"
	profileNameKey         = "name"
	profileTitleKey        = "title"
	profileOrganizationKey = "organization"
	profileReligionKey     = "religion"
	profileAgeKey          = "age"
	profileBirthdayKey     = "birthday"
	profileEyesKey         = "eyes"
	profileHairKey         = "hair"
	profileSkinKey         = "skin"
	profileHandednessKey   = "handedness"
	profileGenderKey       = "gender"
	profileTechLevelKey    = "tech_level"
	profileHeightKey       = "height"
	profileWeightKey       = "weight"
	profileSizeModifierKey = "AdjustedSizeModifier"
	profilePortraitKey     = "portrait"
)

// Standard height and width for the portrait
const (
	PortraitHeight = 96
	PortraitWidth  = 3 * PortraitHeight / 4
)

// Profile holds the profile information for an NPC.
type Profile struct {
	PlayerName        string
	Name              string
	Title             string
	Organization      string
	Religion          string
	Age               string
	Birthday          string
	Eyes              string
	Hair              string
	Skin              string
	Handedness        string
	Gender            string
	TechLevel         string
	PortraitData      string
	portrait          *unison.Image
	Height            measure.Length
	Weight            measure.Weight
	SizeModifier      fixed.F64d4
	SizeModifierBonus fixed.F64d4
}

// NewProfileFromJSON creates a new Profile from a JSON object.
func NewProfileFromJSON(data map[string]interface{}) *Profile {
	p := &Profile{
		PlayerName:   encoding.String(data[profilePlayerNameKey]),
		Name:         encoding.String(data[profileNameKey]),
		Title:        encoding.String(data[profileTitleKey]),
		Organization: encoding.String(data[profileOrganizationKey]),
		Religion:     encoding.String(data[profileReligionKey]),
		Age:          encoding.String(data[profileAgeKey]),
		Birthday:     encoding.String(data[profileBirthdayKey]),
		Eyes:         encoding.String(data[profileEyesKey]),
		Hair:         encoding.String(data[profileHairKey]),
		Skin:         encoding.String(data[profileSkinKey]),
		Handedness:   encoding.String(data[profileHandednessKey]),
		Gender:       encoding.String(data[profileGenderKey]),
		TechLevel:    encoding.String(data[profileTechLevelKey]),
		Height:       measure.LengthFromStringForced(encoding.String(data[profileHeightKey]), units.FeetAndInches),
		Weight:       measure.WeightFromStringForced(encoding.String(data[profileWeightKey]), units.Pound),
		SizeModifier: encoding.Number(data[profileSizeModifierKey]),
		PortraitData: encoding.String(data[profilePortraitKey]),
	}
	return p
}

// ToJSON emits this object as JSON.
func (p *Profile) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(profilePlayerNameKey, p.PlayerName, true, true)
	encoder.KeyedString(profileNameKey, p.Name, true, true)
	encoder.KeyedString(profileTitleKey, p.Title, true, true)
	encoder.KeyedString(profileOrganizationKey, p.Organization, true, true)
	encoder.KeyedString(profileReligionKey, p.Religion, true, true)
	encoder.KeyedString(profileAgeKey, p.Age, true, true)
	encoder.KeyedString(profileBirthdayKey, p.Birthday, true, true)
	encoder.KeyedString(profileEyesKey, p.Eyes, true, true)
	encoder.KeyedString(profileHairKey, p.Hair, true, true)
	encoder.KeyedString(profileSkinKey, p.Skin, true, true)
	encoder.KeyedString(profileHandednessKey, p.Handedness, true, true)
	encoder.KeyedString(profileGenderKey, p.Gender, true, true)
	encoder.KeyedString(profileTechLevelKey, p.TechLevel, true, true)
	if p.Height != 0 {
		encoder.KeyedString(profileHeightKey, p.Height.String(), false, false)
	}
	if p.Weight != 0 {
		encoder.KeyedString(profileWeightKey, p.Weight.String(), false, false)
	}
	encoder.KeyedNumber(profileSizeModifierKey, p.SizeModifier, true)
	encoder.KeyedString(profilePortraitKey, p.PortraitData, true, true)
	encoder.EndObject()
}

// Portrait returns the portrait image, if there is one.
func (p *Profile) Portrait() *unison.Image {
	if p.portrait == nil && p.PortraitData != "" {
		buffer, err := base64.RawStdEncoding.DecodeString(p.PortraitData)
		if err != nil {
			jot.Error(errs.NewWithCause("unable to decode portrait data", err))
			p.portrait = nil
			p.PortraitData = ""
			return nil
		}
		if p.portrait, err = unison.NewImageFromBytes(buffer, 0.5); err != nil {
			jot.Error(errs.NewWithCause("unable to load portrait data", err))
			p.portrait = nil
			p.PortraitData = ""
			return nil
		}
	}
	return p.portrait
}

// AdjustedSizeModifier returns the adjusted size modifier.
func (p *Profile) AdjustedSizeModifier() fixed.F64d4 {
	return (p.SizeModifier + p.SizeModifierBonus).Trunc()
}
