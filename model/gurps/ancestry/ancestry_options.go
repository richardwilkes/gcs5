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
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/eval/f64d4eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	optionsNameKey           = "name"
	optionsHeightFormulaKey  = "height_formula"
	optionsWeightFormulaKey  = "weight_formula"
	optionsAgeFormulaKey     = "age_formula"
	optionsHairKey           = "hair_options"
	optionsEyeKey            = "eye_options"
	optionsSkinKey           = "skin_options"
	optionsHandednessKey     = "handedness_options"
	optionsNameGeneratorsKey = "name_generators"
)

const (
	defaultHeight     = 64
	defaultWeight     = 140
	defaultAge        = 18
	defaultHair       = "Brown"
	defaultEye        = "Brown"
	defaultSkin       = "Brown"
	defaultHandedness = "Right"
)

// Options holds options that may be randomized for an Entity's ancestry.
type Options struct {
	Name              string
	HeightFormula     string
	WeightFormula     string
	AgeFormula        string
	HairOptions       []*StringOption
	EyeOptions        []*StringOption
	SkinOptions       []*StringOption
	HandednessOptions []*StringOption
	NameGenerators    []string
}

// NewOptionsFromJSON creates a new Options from a JSON object.
func NewOptionsFromJSON(data map[string]interface{}) *Options {
	o := &Options{
		Name:              encoding.String(data[optionsNameKey]),
		HeightFormula:     encoding.String(data[optionsHeightFormulaKey]),
		WeightFormula:     encoding.String(data[optionsWeightFormulaKey]),
		AgeFormula:        encoding.String(data[optionsAgeFormulaKey]),
		HairOptions:       StringOptionsFromJSON(encoding.Array(data[optionsHairKey])),
		EyeOptions:        StringOptionsFromJSON(encoding.Array(data[optionsEyeKey])),
		SkinOptions:       StringOptionsFromJSON(encoding.Array(data[optionsSkinKey])),
		HandednessOptions: StringOptionsFromJSON(encoding.Array(data[optionsHandednessKey])),
	}
	array := encoding.Array(data[optionsNameGeneratorsKey])
	if len(array) != 0 {
		o.NameGenerators = make([]string, len(data))
		for i, one := range array {
			o.NameGenerators[i] = encoding.String(one)
		}
	}
	return o
}

// ToJSON emits this object as JSON.
func (o *Options) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(optionsNameKey, o.Name, true, true)
	encoder.KeyedString(optionsHeightFormulaKey, o.HeightFormula, true, true)
	encoder.KeyedString(optionsWeightFormulaKey, o.WeightFormula, true, true)
	encoder.KeyedString(optionsAgeFormulaKey, o.AgeFormula, true, true)
	StringOptionsToJSON(optionsHairKey, o.HairOptions, encoder)
	StringOptionsToJSON(optionsEyeKey, o.EyeOptions, encoder)
	StringOptionsToJSON(optionsSkinKey, o.SkinOptions, encoder)
	StringOptionsToJSON(optionsHandednessKey, o.HandednessOptions, encoder)
	if len(o.NameGenerators) != 0 {
		encoder.Key(optionsNameGeneratorsKey)
		encoder.StartArray()
		for _, one := range o.NameGenerators {
			encoder.String(one)
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}

// RandomHeight returns a randomized height.
func (o *Options) RandomHeight(resolver eval.VariableResolver) measure.Length {
	if result, err := f64d4eval.NewEvaluator(resolver, true).Evaluate(o.HeightFormula); err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", o.HeightFormula))
	} else if value, ok := result.(fixed.F64d4); ok {
		return measure.Length(value)
	}
	return measure.LengthFromInt64(defaultHeight, measure.Inch)
}

// RandomWeight returns a randomized weight.
func (o *Options) RandomWeight(resolver eval.VariableResolver) measure.Weight {
	if result, err := f64d4eval.NewEvaluator(resolver, true).Evaluate(o.WeightFormula); err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", o.WeightFormula))
	} else if value, ok := result.(fixed.F64d4); ok {
		return measure.Weight(value)
	}
	return measure.WeightFromInt64(defaultWeight, measure.Pound)
}

// RandomAge returns a randomized age.
func (o *Options) RandomAge(resolver eval.VariableResolver) int {
	if result, err := f64d4eval.NewEvaluator(resolver, true).Evaluate(o.AgeFormula); err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", o.AgeFormula))
	} else if value, ok := result.(fixed.F64d4); ok {
		return int(value.AsInt64())
	}
	return defaultAge
}

// RandomHair returns a randomized hair.
func (o *Options) RandomHair() string {
	if choice := ChooseStringOption(o.HairOptions); choice != "" {
		return choice
	}
	return defaultHair
}

// RandomEye returns a randomized eye.
func (o *Options) RandomEye() string {
	if choice := ChooseStringOption(o.EyeOptions); choice != "" {
		return choice
	}
	return defaultEye
}

// RandomSkin returns a randomized skin.
func (o *Options) RandomSkin() string {
	if choice := ChooseStringOption(o.SkinOptions); choice != "" {
		return choice
	}
	return defaultSkin
}

// RandomHandedness returns a randomized handedness.
func (o *Options) RandomHandedness() string {
	if choice := ChooseStringOption(o.HandednessOptions); choice != "" {
		return choice
	}
	return defaultHandedness
}

// RandomName returns a randomized name.
func (o *Options) RandomName(nameGeneratorRefs []*NameGeneratorRef) string {
	m := make(map[string]*NameGeneratorRef)
	for _, one := range nameGeneratorRefs {
		m[one.FileRef.Name] = one
	}
	var buffer strings.Builder
	for _, one := range o.NameGenerators {
		if ref, ok := m[one]; ok {
			if generator, err := ref.Generator(); err != nil {
				jot.Error(err)
			} else {
				if name := strings.TrimSpace(generator.Generate()); name != "" {
					if buffer.Len() != 0 {
						buffer.WriteByte(' ')
					}
					buffer.WriteString(name)
				}
			}
		}
	}
	return buffer.String()
}
