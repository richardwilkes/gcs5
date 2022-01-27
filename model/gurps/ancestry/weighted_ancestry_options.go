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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/rand"
)

// WeightedAncestryOptions is a string that has a weight associated with it.
type WeightedAncestryOptions struct {
	Weight int
	Value  *Options
}

// NewWeightedAncestryOptionsFromJSON creates a new WeightedAncestryOptions from a JSON object.
func NewWeightedAncestryOptionsFromJSON(data map[string]interface{}) *WeightedAncestryOptions {
	return &WeightedAncestryOptions{
		Weight: int(encoding.Number(data[optionWeightKey]).AsInt64()),
		Value:  NewOptionsFromJSON(encoding.Object(data[optionValueKey])),
	}
}

// ToJSON emits this object as JSON.
func (o *WeightedAncestryOptions) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedNumber(optionWeightKey, fixed.F64d4FromInt64(int64(o.Weight)), false)
	encoding.ToKeyedJSON(o.Value, optionValueKey, encoder)
	encoder.EndObject()
}

// Valid returns true if this option has a valid weight.
func (o *WeightedAncestryOptions) Valid() bool {
	return o.Weight > 0
}

// WeightedAncestryOptionsFromJSON creates a slice of options from a JSON array.
func WeightedAncestryOptionsFromJSON(array []interface{}) []*WeightedAncestryOptions {
	if len(array) == 0 {
		return nil
	}
	options := make([]*WeightedAncestryOptions, len(array))
	for i, one := range array {
		options[i] = NewWeightedAncestryOptionsFromJSON(encoding.Object(one))
	}
	return options
}

// WeightedAncestryOptionsToJSON emits the options as JSON.
func WeightedAncestryOptionsToJSON(key string, options []*WeightedAncestryOptions, encoder *encoding.JSONEncoder) {
	if len(options) != 0 {
		encoder.Key(key)
		encoder.StartArray()
		for _, one := range options {
			one.ToJSON(encoder)
		}
		encoder.EndArray()
	}
}

// ChooseWeightedAncestryOptions selects a string option from the available set.
func ChooseWeightedAncestryOptions(options []*WeightedAncestryOptions) *Options {
	total := 0
	for _, one := range options {
		total += one.Weight
	}
	if total > 0 {
		choice := 1 + rand.NewCryptoRand().Intn(total)
		for _, one := range options {
			choice -= one.Weight
			if choice < 1 {
				return one.Value
			}
		}
	}
	return nil
}
