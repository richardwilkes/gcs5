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
	"io/fs"
	"sort"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

// AttributeDefs holds a set of AttributeDef objects.
type AttributeDefs struct {
	Set map[string]*AttributeDef
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

// FactoryAttributeDefs returns the factory AttributeDef set.
func FactoryAttributeDefs() *AttributeDefs {
	defs, err := NewAttributeDefsFromFile(embeddedFS, "data/standard.attr")
	jot.FatalIfErr(err)
	return defs
}

// NewAttributeDefsFromFile loads an AttributeDef set from a file.
func NewAttributeDefsFromFile(fsys fs.FS, filePath string) (*AttributeDefs, error) {
	data, err := encoding.LoadJSONFromFS(fsys, filePath)
	if err != nil {
		return nil, err
	}
	// Check for older formats
	if obj := encoding.Object(data); obj != nil {
		var exists bool
		if data, exists = obj["attributes"]; !exists {
			if data, exists = obj["attribute_settings"]; !exists {
				return nil, errs.New("invalid attribute definitions file: " + filePath)
			}
		}
	}
	return NewAttributeDefsFromJSON(encoding.Array(data)), nil
}

// NewAttributeDefsFromJSON creates a new AttributeDefs from a JSON object.
func NewAttributeDefsFromJSON(data []interface{}) *AttributeDefs {
	a := &AttributeDefs{Set: make(map[string]*AttributeDef)}
	for i, one := range encoding.Array(data) {
		def := NewAttributeDefFromJSON(encoding.Object(one), i+1)
		a.Set[def.ID()] = def
	}
	return a
}

// Save writes the AttributeDefs to the file as JSON.
func (a *AttributeDefs) Save(filePath string) error {
	return encoding.SaveJSON(filePath, true, a.ToJSON)
}

// ToJSON emits this object as JSON.
func (a *AttributeDefs) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartArray()
	for _, def := range a.List() {
		def.ToJSON(encoder)
	}
	encoder.EndArray()
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
