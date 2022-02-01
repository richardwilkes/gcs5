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
	"encoding/json"
	"io/fs"
	"sort"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
)

// AttributeDefs holds a set of AttributeDef objects.
type AttributeDefs struct {
	Set map[string]*AttributeDef
}

// ResolveAttributeName returns the name of the attribute, if possible.
func ResolveAttributeName(entity *Entity, attribute string) string {
	if def := AttributeDefsFor(entity).Set[attribute]; def != nil {
		return def.Name
	}
	return attribute
}

// AttributeDefsFor returns the AttributeDefs for the given Entity, or the global settings if the Entity is nil.
func AttributeDefsFor(entity *Entity) *AttributeDefs {
	return SheetSettingsFor(entity).Attributes
}

// DefaultAttributeIDFor returns the default attribute ID to use for the given Entity, which may be nil.
func DefaultAttributeIDFor(entity *Entity) string {
	list := AttributeDefsFor(entity).List()
	if len(list) != 0 {
		return list[0].ID()
	}
	return "st"
}

// AttributeIDFor looks up the preferred ID and if it cannot be found, falls back to a default. 'entity' may be nil.
func AttributeIDFor(entity *Entity, preferred string) string {
	defs := AttributeDefsFor(entity)
	if _, exists := defs.Set[preferred]; exists {
		return preferred
	}
	if list := defs.List(); len(list) != 0 {
		return list[0].ID()
	}
	return "st"
}

// FactoryAttributeDefs returns the factory AttributeDef set.
func FactoryAttributeDefs() *AttributeDefs {
	defs, err := NewAttributeDefsFromFile(embeddedFS, "data/standard.attr")
	jot.FatalIfErr(err)
	return defs
}

// NewAttributeDefsFromFile loads an AttributeDef set from a file.
func NewAttributeDefsFromFile(fsys fs.FS, filePath string) (*AttributeDefs, error) {
	var a AttributeDefs
	if err := xfs.LoadJSONFromFS(fsys, filePath, &a); err != nil {
		// Check for older formats
		var old struct {
			Attributes        *AttributeDefs `json:"attributes"`
			AttributeSettings *AttributeDefs `json:"attribute_settings"`
		}
		if err = xfs.LoadJSONFromFS(fsys, filePath, &a); err != nil {
			return nil, err
		}
		if old.Attributes != nil {
			return old.Attributes, nil
		}
		if old.AttributeSettings != nil {
			return old.AttributeSettings, nil
		}
		return nil, errs.New("invalid attribute definitions file: " + filePath)
	}
	return &a, nil
}

// Save writes the AttributeDefs to the file as JSON.
func (a *AttributeDefs) Save(filePath string) error {
	return xfs.SaveJSON(filePath, a, true)
}

// MarshalJSON implements json.Marshaler.
func (a *AttributeDefs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.List())
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttributeDefs) UnmarshalJSON(data []byte) error {
	var list []*AttributeDef
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	a.Set = make(map[string]*AttributeDef, len(list))
	for i, one := range list {
		one.Order = i + 1
		a.Set[one.ID()] = one
	}
	return nil
}

// List returns the map of AttributeDef objects as an ordered list.
func (a *AttributeDefs) List() []*AttributeDef {
	list := make([]*AttributeDef, 0, len(a.Set))
	for _, v := range a.Set {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Order < list[j].Order })
	return list
}
