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
	"bytes"
	"context"
	"encoding/binary"
	"hash"
	"hash/crc64"
	"io/fs"
	"sort"

	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

const attributeSettingsListTypeKey = "attribute_settings"

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
	return gid.Strength
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
	return gid.Strength
}

// FactoryAttributeDefs returns the factory AttributeDef set.
func FactoryAttributeDefs() *AttributeDefs {
	defs, err := NewAttributeDefsFromFile(embeddedFS, "data/standard.attr")
	jot.FatalIfErr(err)
	return defs
}

type attributeDefsData struct {
	Type    string         `json:"type"`
	Version int            `json:"version"`
	Rows    *AttributeDefs `json:"rows"`
}

// NewAttributeDefsFromFile loads an AttributeDef set from a file.
func NewAttributeDefsFromFile(fileSystem fs.FS, filePath string) (*AttributeDefs, error) {
	var data struct {
		attributeDefsData
		OldKey1 *AttributeDefs `json:"attribute_settings"`
		OldKey2 *AttributeDefs `json:"attributes"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != attributeSettingsListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	var defs *AttributeDefs
	if data.attributeDefsData.Rows != nil {
		defs = data.attributeDefsData.Rows
	}
	if defs == nil && data.OldKey1 != nil {
		defs = data.OldKey1
	}
	if defs == nil && data.OldKey2 != nil {
		defs = data.OldKey2
	}
	if defs == nil {
		defs = FactoryAttributeDefs()
	} else {
		defs.EnsureValidity()
	}
	return defs, nil
}

// Save writes the AttributeDefs to the file as JSON.
func (a *AttributeDefs) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &attributeDefsData{
		Type:    attributeSettingsListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    a,
	})
}

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (a *AttributeDefs) EnsureValidity() {
	// TODO: Implement validity check
}

// MarshalJSON implements json.Marshaler.
func (a *AttributeDefs) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	e := json.NewEncoder(&buffer)
	e.SetEscapeHTML(false)
	err := e.Encode(a.List())
	return buffer.Bytes(), err
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

// Clone a copy of this.
func (a *AttributeDefs) Clone() *AttributeDefs {
	clone := &AttributeDefs{Set: make(map[string]*AttributeDef)}
	for k, v := range a.Set {
		clone.Set[k] = v.Clone()
	}
	return clone
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

// CRC64 calculates a CRC-64 for this data.
func (a *AttributeDefs) CRC64() uint64 {
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	a.crc64(h)
	return h.Sum64()
}

func (a *AttributeDefs) crc64(h hash.Hash64) {
	var buffer [8]byte
	binary.LittleEndian.PutUint64(buffer[:], uint64(len(a.Set)))
	h.Write(buffer[:])
	for _, one := range a.List() {
		one.crc64(h)
	}
}
