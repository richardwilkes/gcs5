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
	"sort"

	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

var (
	currentBindings              = make(KeyBindings)
	factoryBindings              = make(map[string]*Binding)
	_               json.Omitter = KeyBindings{}
)

// KeyBindings holds a set of key bindings.
type KeyBindings map[string]unison.KeyBinding

// Binding holds a single key binding.
type Binding struct {
	ID         string
	KeyBinding unison.KeyBinding
	Action     *unison.Action
}

// RegisterKeyBinding register a keybinding.
func RegisterKeyBinding(id string, action *unison.Action) {
	if _, exists := factoryBindings[id]; exists {
		return
	}
	delete(currentBindings, id)
	factoryBindings[id] = &Binding{
		ID:         id,
		KeyBinding: action.KeyBinding,
		Action:     action,
	}
}

// CurrentBindings returns a sorted list with the current bindings. Note that only the ID and Action field are valid.
func CurrentBindings() []*Binding {
	list := make([]*Binding, 0, len(factoryBindings))
	for _, v := range factoryBindings {
		list = append(list, &Binding{
			ID:     v.ID,
			Action: v.Action,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		if txt.NaturalLess(list[i].Action.Title, list[j].Action.Title, true) {
			return true
		}
		if list[i].Action.Title != list[j].Action.Title {
			return false
		}
		return list[i].ID < list[j].ID
	})
	return list
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
		if info, ok := factoryBindings[k]; ok && v != info.KeyBinding {
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
	data := make(map[string]unison.KeyBinding, len(currentBindings))
	for k, v := range currentBindings {
		if info, ok := factoryBindings[k]; ok && info.KeyBinding != v {
			data[k] = v
		}
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *KeyBindings) UnmarshalJSON(data []byte) error {
	m := make(map[string]unison.KeyBinding, len(currentBindings))
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	kb := make(KeyBindings)
	for k, v := range factoryBindings {
		if current, ok := m[k]; ok && current != v.KeyBinding {
			kb[k] = current
		} else {
			kb[k] = v.KeyBinding
		}
	}
	*b = kb
	return nil
}

// MakeCurrent applies these key bindings to the current key bindings set.
func (b KeyBindings) MakeCurrent() {
	var actions []*unison.Action
	for k, v := range factoryBindings {
		current, ok := b[k]
		if !ok {
			current = v.KeyBinding
		}
		if current != v.KeyBinding {
			currentBindings[k] = current
		} else {
			delete(currentBindings, k)
		}
		if v.Action.KeyBinding != current {
			v.Action.KeyBinding = current
			actions = append(actions, v.Action)
		}
	}
	if len(actions) != 0 {
		factory := unison.DefaultMenuFactory()
		for _, w := range unison.Windows() {
			if bar := factory.BarForWindowNoCreate(w); bar != nil {
				for _, a := range actions {
					if item := bar.Item(a.ID); item != nil {
						item.SetKeyBinding(a.KeyBinding)
					}
				}
				if factory.BarIsPerWindow() {
					break
				}
			}
		}
	}
}

// Set the binding for the given ID.
func (b KeyBindings) Set(id string, binding unison.KeyBinding) {
	if f, ok := factoryBindings[id]; ok {
		if f.KeyBinding != binding {
			b[id] = binding
		} else {
			delete(b, id)
		}
	}
}

// Reset to factory defaults.
func (b *KeyBindings) Reset() {
	*b = make(KeyBindings)
}

// ResetOne resets one font by ID to factory defaults.
func (b KeyBindings) ResetOne(id string) {
	delete(currentBindings, id)
}
