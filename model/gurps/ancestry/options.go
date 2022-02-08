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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
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
	Name              string          `json:"name,omitempty"`
	HeightFormula     string          `json:"height_formula,omitempty"`
	WeightFormula     string          `json:"weight_formula,omitempty"`
	AgeFormula        string          `json:"age_formula,omitempty"`
	HairOptions       []*StringOption `json:"hair_options,omitempty"`
	EyeOptions        []*StringOption `json:"eye_options,omitempty"`
	SkinOptions       []*StringOption `json:"skin_options,omitempty"`
	HandednessOptions []*StringOption `json:"handedness_options,omitempty"`
	NameGenerators    []string        `json:"name_generators,omitempty"`
}

// RandomHeight returns a randomized height.
func (o *Options) RandomHeight(resolver eval.VariableResolver) measure.Length {
	value := fxp.EvaluateToNumber(o.HeightFormula, resolver)
	if value <= 0 {
		return measure.LengthFromInt(defaultHeight, measure.Inch)
	}
	return measure.Length(value)
}

// RandomWeight returns a randomized weight.
func (o *Options) RandomWeight(resolver eval.VariableResolver) measure.Weight {
	value := fxp.EvaluateToNumber(o.WeightFormula, resolver)
	if value <= 0 {
		return measure.WeightFromInt(defaultWeight, measure.Pound)
	}
	return measure.Weight(value)
}

// RandomAge returns a randomized age.
func (o *Options) RandomAge(resolver eval.VariableResolver) int {
	value := fxp.EvaluateToNumber(o.AgeFormula, resolver).Trunc()
	if value <= 0 {
		value = fixed.F64d4FromInt(defaultAge)
	}
	return value.AsInt()
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
