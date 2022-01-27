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

package ancestry

import (
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/eval"
)

const (
	ancestryCommonOptionsKey = "common_options"
	ancestryGenderOptionsKey = "gender_options"
)

// Ancestry holds details necessary to generate ancestry-specific customizations.
type Ancestry struct {
	CommonOptions *Options
	GenderOptions []*WeightedAncestryOptions
}

// AvailableAncestries scans the libraries and returns the available ancestries.
func AvailableAncestries(libraries *library.Libraries) []*library.NamedFileSet {
	return library.ScanForNamedFileSets(embeddedFS, "data", ".ancestry", libraries)
}

// NewAncestryFromJSON creates a new Ancestry from a JSON object.
func NewAncestryFromJSON(data map[string]interface{}) *Ancestry {
	a := &Ancestry{
		GenderOptions: WeightedAncestryOptionsFromJSON(encoding.Array(data[ancestryGenderOptionsKey])),
	}
	obj := encoding.Object(data[ancestryCommonOptionsKey])
	if len(obj) != 0 {
		a.CommonOptions = NewOptionsFromJSON(obj)
	}
	return a
}

// ToJSON emits this object as JSON.
func (a *Ancestry) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	if a.CommonOptions != nil {
		encoding.ToKeyedJSON(a.CommonOptions, ancestryCommonOptionsKey, encoder)
	}
	WeightedAncestryOptionsToJSON(ancestryGenderOptionsKey, a.GenderOptions, encoder)
	encoder.EndObject()
}

// RandomGender returns a randomized gender.
func (a *Ancestry) RandomGender() string {
	if choice := ChooseWeightedAncestryOptions(a.GenderOptions); choice != nil {
		return choice.Name
	}
	return ""
}

// GenderedOptions returns the options for the specified gender, or nil.
func (a *Ancestry) GenderedOptions(gender string) *Options {
	gender = strings.TrimSpace(gender)
	for _, one := range a.GenderOptions {
		if strings.EqualFold(one.Value.Name, gender) {
			return one.Value
		}
	}
	return nil
}

// RandomHeight returns a randomized height.
func (a *Ancestry) RandomHeight(resolver eval.VariableResolver, gender string) measure.Length {
	if options := a.GenderedOptions(gender); options != nil && options.HeightFormula != "" {
		return options.RandomHeight(resolver)
	}
	if a.CommonOptions != nil && a.CommonOptions.HeightFormula != "" {
		return a.CommonOptions.RandomHeight(resolver)
	}
	return measure.LengthFromInt64(defaultHeight, measure.Inch)
}

// RandomWeight returns a randomized weight.
func (a *Ancestry) RandomWeight(resolver eval.VariableResolver, gender string) measure.Weight {
	if options := a.GenderedOptions(gender); options != nil && options.WeightFormula != "" {
		return options.RandomWeight(resolver)
	}
	if a.CommonOptions != nil && a.CommonOptions.WeightFormula != "" {
		return a.CommonOptions.RandomWeight(resolver)
	}
	return measure.WeightFromInt64(defaultWeight, measure.Pound)
}

// RandomAge returns a randomized age.
func (a *Ancestry) RandomAge(resolver eval.VariableResolver, gender string) int {
	if options := a.GenderedOptions(gender); options != nil && options.AgeFormula != "" {
		return options.RandomAge(resolver)
	}
	if a.CommonOptions != nil && a.CommonOptions.AgeFormula != "" {
		return a.CommonOptions.RandomAge(resolver)
	}
	return defaultAge
}

// RandomHair returns a randomized hair.
func (a *Ancestry) RandomHair(gender string) string {
	if options := a.GenderedOptions(gender); options != nil && len(options.HairOptions) != 0 {
		return options.RandomHair()
	}
	if a.CommonOptions != nil && len(a.CommonOptions.HairOptions) != 0 {
		return a.CommonOptions.RandomHair()
	}
	return defaultHair
}

// RandomEye returns a randomized eye.
func (a *Ancestry) RandomEye(gender string) string {
	if options := a.GenderedOptions(gender); options != nil && len(options.EyeOptions) != 0 {
		return options.RandomEye()
	}
	if a.CommonOptions != nil && len(a.CommonOptions.EyeOptions) != 0 {
		return a.CommonOptions.RandomEye()
	}
	return defaultEye
}

// RandomSkin returns a randomized skin.
func (a *Ancestry) RandomSkin(gender string) string {
	if options := a.GenderedOptions(gender); options != nil && len(options.SkinOptions) != 0 {
		return options.RandomSkin()
	}
	if a.CommonOptions != nil && len(a.CommonOptions.SkinOptions) != 0 {
		return a.CommonOptions.RandomSkin()
	}
	return defaultSkin
}

// RandomHandedness returns a randomized handedness.
func (a *Ancestry) RandomHandedness(gender string) string {
	if options := a.GenderedOptions(gender); options != nil && len(options.HandednessOptions) != 0 {
		return options.RandomHandedness()
	}
	if a.CommonOptions != nil && len(a.CommonOptions.HandednessOptions) != 0 {
		return a.CommonOptions.RandomHandedness()
	}
	return defaultHandedness
}

// RandomName returns a randomized name.
func (a *Ancestry) RandomName(nameGeneratorRefs []*NameGeneratorRef, gender string) string {
	if options := a.GenderedOptions(gender); options != nil && len(options.NameGenerators) != 0 {
		return options.RandomName(nameGeneratorRefs)
	}
	if a.CommonOptions != nil && len(a.CommonOptions.NameGenerators) != 0 {
		return a.CommonOptions.RandomName(nameGeneratorRefs)
	}
	return ""
}
