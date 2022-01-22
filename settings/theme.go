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

package settings

import (
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/unison"
)

const (
	themeColorsKey           = "colors"
	themeFontsKey            = "fonts"
	fontDescriptorFamilyKey  = "family"
	fontDescriptorSizeKey    = "size"
	fontDescriptorWeightKey  = "weight"
	fontDescriptorSpacingKey = "spacing"
	fontDescriptorSlantKey   = "slant"
)

// Theme holds colors and fonts used in the UI.
type Theme struct {
	Colors map[string]unison.Color
	Fonts  map[string]unison.FontDescriptor
}

// NewTheme creates a new Theme.
func NewTheme() *Theme {
	return &Theme{
		Colors: make(map[string]unison.Color),
		Fonts:  make(map[string]unison.FontDescriptor),
	}
}

// NewThemeFromJSON creates a new Theme from a JSON object.
func NewThemeFromJSON(data map[string]interface{}) *Theme {
	p := NewTheme()
	for k, v := range encoding.Object(data[themeColorsKey]) {
		if c, err := unison.ColorDecode(encoding.String(v)); err != nil {
			p.Colors[k] = c
		}
	}
	for k, v := range encoding.Object(data[themeFontsKey]) {
		one := encoding.Object(v)
		p.Fonts[k] = unison.FontDescriptor{
			Family:  encoding.String(one[fontDescriptorFamilyKey]),
			Size:    float32(encoding.Number(one[fontDescriptorSizeKey]).AsFloat64()),
			Weight:  unison.WeightFromString(encoding.String(one[fontDescriptorWeightKey])),
			Spacing: unison.SpacingFromString(encoding.String(one[fontDescriptorWeightKey])),
			Slant:   unison.SlantFromString(encoding.String(one[fontDescriptorWeightKey])),
		}
	}
	return p
}

// ToKeyedJSON emits this object as JSON with the specified key, but only if not empty.
func (p *Theme) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	if len(p.Colors) != 0 || len(p.Fonts) != 0 {
		encoder.Key(key)
		p.ToJSON(encoder)
	}
}

// ToJSON emits this object as JSON.
func (p *Theme) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()

	encoder.Key(themeColorsKey)
	encoder.StartObject()
	keys := make([]string, 0, len(p.Colors))
	for k := range p.Colors {
		keys = append(keys, k)
	}
	txt.SortStringsNaturalAscending(keys)
	for _, k := range keys {
		encoder.KeyedString(k, p.Colors[k].String(), false, false)
	}
	encoder.EndObject()

	encoder.Key(themeFontsKey)
	encoder.StartObject()
	keys = make([]string, 0, len(p.Fonts))
	for k := range p.Fonts {
		keys = append(keys, k)
	}
	txt.SortStringsNaturalAscending(keys)
	for _, k := range keys {
		encoder.Key(k)
		encoder.StartObject()
		desc := p.Fonts[k]
		encoder.KeyedString(fontDescriptorFamilyKey, desc.Family, false, false)
		encoder.KeyedNumber(fontDescriptorSizeKey, fixed.F64d4FromFloat64(float64(desc.Size)), false)
		encoder.KeyedString(fontDescriptorWeightKey, desc.Weight.String(), true, true)
		encoder.KeyedString(fontDescriptorSpacingKey, desc.Spacing.String(), true, true)
		encoder.KeyedString(fontDescriptorSlantKey, desc.Slant.String(), true, true)
		encoder.EndObject()
	}
	encoder.EndObject()

	encoder.EndObject()
}
