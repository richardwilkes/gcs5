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
	"context"
	"embed"
	"encoding/binary"
	"hash"
	"hash/crc64"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
)

const bodyTypeListTypeKey = "body_type"

//go:embed data
var embeddedFS embed.FS

// BodyType holds a set of hit locations.
type BodyType struct {
	Name           string         `json:"name,omitempty"`
	Roll           *dice.Dice     `json:"roll"`
	Locations      []*HitLocation `json:"locations,omitempty"`
	owningLocation *HitLocation
	locationLookup map[string]*HitLocation
}

type bodyTypeListData struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	*BodyType
}

// FactoryBodyType returns a new copy of the default factory BodyType.
func FactoryBodyType() *BodyType {
	bodyType, err := NewBodyTypeFromFile(embeddedFS, "data/body_types/Humanoid.body")
	jot.FatalIfErr(err)
	return bodyType
}

// FactoryBodyTypes returns the list of the known factory BodyTypes.
func FactoryBodyTypes() []*BodyType {
	entries, err := embeddedFS.ReadDir("data/body_types")
	jot.FatalIfErr(err)
	list := make([]*BodyType, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if path.Ext(name) == ".body" {
			var bodyType *BodyType
			bodyType, err = NewBodyTypeFromFile(embeddedFS, "data/body_types/"+name)
			jot.FatalIfErr(err)
			list = append(list, bodyType)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return txt.NaturalLess(list[i].Name, list[j].Name, true)
	})
	return list
}

// NewBodyTypeFromFile loads an BodyType from a file.
func NewBodyTypeFromFile(fileSystem fs.FS, filePath string) (*BodyType, error) {
	var data struct {
		bodyTypeListData
		HitLocations []*HitLocation `json:"hit_locations"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != bodyTypeListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	if data.Locations == nil {
		data.BodyType.Locations = data.HitLocations
	}
	data.BodyType.EnsureValidity()
	data.BodyType.Update(nil)
	return data.BodyType, nil
}

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (b *BodyType) EnsureValidity() {
	// TODO: Implement validity check
}

// Clone a copy of this.
func (b *BodyType) Clone(entity *Entity, owningLocation *HitLocation) *BodyType {
	clone := &BodyType{
		Name:           b.Name,
		Roll:           dice.New(b.Roll.String()),
		Locations:      make([]*HitLocation, len(b.Locations)),
		owningLocation: owningLocation,
	}
	for i, one := range b.Locations {
		clone.Locations[i] = one.Clone(entity, clone)
	}
	clone.Update(entity)
	return clone
}

// Save writes the BodyType to the file as JSON.
func (b *BodyType) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &bodyTypeListData{
		Type:     bodyTypeListTypeKey,
		Version:  gid.CurrentDataVersion,
		BodyType: b,
	})
}

// Update the role ranges and populate the lookup map.
func (b *BodyType) Update(entity *Entity) {
	b.updateRollRanges()
	b.locationLookup = make(map[string]*HitLocation)
	b.populateMap(entity, b.locationLookup)
}

// SetOwningLocation sets the owning HitLocation.
func (b *BodyType) SetOwningLocation(loc *HitLocation) {
	b.owningLocation = loc
	if loc != nil {
		b.Name = ""
	}
}

func (b *BodyType) updateRollRanges() {
	start := b.Roll.Minimum(false)
	for _, location := range b.Locations {
		start = location.updateRollRange(start)
	}
}

func (b *BodyType) populateMap(entity *Entity, m map[string]*HitLocation) {
	for _, location := range b.Locations {
		location.populateMap(entity, m)
	}
}

// AddLocation adds a HitLocation to the end of list.
func (b *BodyType) AddLocation(loc *HitLocation) {
	b.Locations = append(b.Locations, loc)
	loc.owningTable = b
}

// RemoveLocation removes a HitLocation.
func (b *BodyType) RemoveLocation(loc *HitLocation) {
	for i, one := range b.Locations {
		if one == loc {
			copy(b.Locations[i:], b.Locations[i+1:])
			b.Locations[len(b.Locations)-1] = nil
			b.Locations = b.Locations[:len(b.Locations)-1]
			loc.owningTable = nil
		}
	}
}

// UniqueHitLocations returns the list of unique hit locations.
func (b *BodyType) UniqueHitLocations(entity *Entity) []*HitLocation {
	if len(b.locationLookup) == 0 {
		b.Update(entity)
	}
	locations := make([]*HitLocation, 0, len(b.locationLookup))
	for _, v := range b.locationLookup {
		locations = append(locations, v)
	}
	sort.Slice(locations, func(i, j int) bool {
		if txt.NaturalLess(locations[i].ChoiceName, locations[j].ChoiceName, false) {
			return true
		}
		if strings.EqualFold(locations[i].ChoiceName, locations[j].ChoiceName) {
			return txt.NaturalLess(locations[i].ID(), locations[j].ID(), false)
		}
		return false
	})
	return locations
}

// LookupLocationByID returns the HitLocation that matches the given ID.
func (b *BodyType) LookupLocationByID(entity *Entity, idStr string) *HitLocation {
	if len(b.locationLookup) == 0 {
		b.Update(entity)
	}
	return b.locationLookup[idStr]
}

// CRC64 calculates a CRC-64 for this data.
func (b *BodyType) CRC64() uint64 {
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	b.crc64(h)
	return h.Sum64()
}

func (b *BodyType) crc64(h hash.Hash64) {
	h.Write([]byte(b.Name))
	h.Write([]byte(b.Roll.String()))
	var buffer [8]byte
	binary.LittleEndian.PutUint64(buffer[:], uint64(len(b.Locations)))
	h.Write(buffer[:])
	for _, loc := range b.Locations {
		loc.crc64(h)
	}
}
