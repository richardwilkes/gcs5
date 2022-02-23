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
	"context"
	"io/fs"

	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
)

var (
	// CurrentBindings holds the current key bindings.
	CurrentBindings = make(KeyBindings)
	// FactoryBindings holds the original key bindings before any modifications.
	FactoryBindings              = make(KeyBindings)
	binders                      = make(map[string]func(id, binding string))
	_               json.Omitter = KeyBindings{}
)

// KeyBindings holds a set of key bindings.
type KeyBindings map[string]string

// RegisterKeyBinding register a keybinding.
func RegisterKeyBinding(id, binding string, f func(id, binding string)) {
	binders[id] = f
	if _, exists := FactoryBindings[id]; exists {
		return
	}
	FactoryBindings[id] = binding
	CurrentBindings[id] = binding
}

// NewKeyBindingsFromFS creates a new set of key bindings from a file. Any missing values will be filled in with
// defaults.
func NewKeyBindingsFromFS(fileSystem fs.FS, filePath string) (KeyBindings, error) {
	var b KeyBindings
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// ShouldOmit implements json.Omitter.
func (b KeyBindings) ShouldOmit() bool {
	for k, v := range b {
		if FactoryBindings[k] != v {
			return false
		}
	}
	return true
}

// Save writes the Fonts to the file as JSON.
func (b KeyBindings) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, b)
}

// MarshalJSON implements json.Marshaler.
func (b KeyBindings) MarshalJSON() ([]byte, error) {
	data := make(map[string]string, len(CurrentBindings))
	for k, v := range CurrentBindings {
		if FactoryBindings[k] != v {
			data[k] = v
		}
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *KeyBindings) UnmarshalJSON(data []byte) error {
	m := make(map[string]string, len(CurrentBindings))
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	bindings := make(KeyBindings)
	for k, v := range FactoryBindings {
		if current, ok := m[k]; ok {
			bindings[k] = current
		} else {
			bindings[k] = v
		}
	}
	*b = bindings
	return nil
}

// MakeCurrent applies these key bindings to the current key bindings set.
func (b KeyBindings) MakeCurrent() {
	for k, v := range FactoryBindings {
		if current, ok := b[k]; ok {
			v = current
		}
		CurrentBindings[k] = v
		if f, ok := binders[k]; ok {
			f(k, v)
		}
	}
}

// Reset to factory defaults.
func (b *KeyBindings) Reset() {
	*b = make(KeyBindings)
	for k, v := range FactoryBindings {
		(*b)[k] = v
	}
}

// ResetOne resets one font by ID to factory defaults.
func (b KeyBindings) ResetOne(id string) {
	b[id] = FactoryBindings[id]
}
