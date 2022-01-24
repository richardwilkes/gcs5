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
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/txt"
)

// KeyBindings holds overrides for key bindings, usually for menu items.
type KeyBindings struct {
	bindings map[string]string
}

// NewKeyBindings creates a new, empty, KeyBindings object.
func NewKeyBindings() *KeyBindings {
	return &KeyBindings{bindings: make(map[string]string)}
}

// NewKeyBindingsFromJSON creates a new KeyBindings from a JSON object.
func NewKeyBindingsFromJSON(data map[string]interface{}) *KeyBindings {
	p := NewKeyBindings()
	for k, v := range data {
		if s := strings.TrimSpace(encoding.String(v)); s != "" {
			p.bindings[k] = s
		}
	}
	return p
}

// Empty implements encoding.Empty.
func (p *KeyBindings) Empty() bool {
	return len(p.bindings) == 0
}

// ToJSON emits this object as JSON.
func (p *KeyBindings) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	keys := make([]string, 0, len(p.bindings))
	for k := range p.bindings {
		keys = append(keys, k)
	}
	txt.SortStringsNaturalAscending(keys)
	for _, k := range keys {
		encoder.KeyedString(k, p.bindings[k], false, false)
	}
	encoder.EndObject()
}

// Count returns the number of key binding overrides being tracked.
func (p *KeyBindings) Count() int {
	return len(p.bindings)
}

// Override returns the override for the key or an empty string.
func (p *KeyBindings) Override(key string) string {
	return p.bindings[key]
}

// SetOverride sets the override for the key. Pass in an empty string to remove an override.
func (p *KeyBindings) SetOverride(key, override string) {
	override = strings.TrimSpace(override)
	if override == "" {
		delete(p.bindings, key)
	} else {
		p.bindings[key] = override
	}
}
