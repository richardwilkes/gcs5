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
	"path"
	"sort"
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
)

const (
	bodyTypeNameKey      = "name"
	bodyTypeRollKey      = "roll"
	bodyTypeLocationsKey = "locations"
)

// BodyType holds a set of hit locations.
type BodyType struct {
	Name           string
	Roll           *dice.Dice
	locations      []*HitLocation
	owningLocation *HitLocation
	locationLookup map[string]*HitLocation
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
func NewBodyTypeFromFile(fsys fs.FS, filePath string) (*BodyType, error) {
	data, err := encoding.LoadJSONFromFS(fsys, filePath)
	if err != nil {
		return nil, err
	}
	obj := encoding.Object(data)
	// Check for older formats
	var exists bool
	if data, exists = obj["hit_locations"]; exists {
		obj = encoding.Object(data)
	}
	if obj == nil {
		return nil, errs.New("invalid body type definition file: " + filePath)
	}
	return NewBodyTypeFromJSON(obj), nil
}

// NewBodyTypeFromJSON creates a new BodyType from a JSON object.
func NewBodyTypeFromJSON(data map[string]interface{}) *BodyType {
	a := &BodyType{
		Name: encoding.String(data[bodyTypeNameKey]),
		Roll: dice.New(encoding.String(data[bodyTypeRollKey])),
	}
	array := encoding.Array(data[bodyTypeLocationsKey])
	if len(array) != 0 {
		a.locations = make([]*HitLocation, 0, len(array))
		for _, one := range array {
			a.AddLocation(NewHitLocationFromJSON(encoding.Object(one)))
		}
	}
	a.Update()
	return a
}

// Save writes the BodyType to the file as JSON.
func (b *BodyType) Save(filePath string) error {
	return encoding.SaveJSON(filePath, true, func(encoder *encoding.JSONEncoder) {
		b.ToJSON(encoder, nil)
	})
}

// ToJSON emits this object as JSON.
func (b *BodyType) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	encoder.KeyedString(bodyTypeNameKey, b.Name, true, true)
	encoder.KeyedString(bodyTypeRollKey, b.Roll.String(), false, false)
	if len(b.locations) != 0 {
		encoder.Key(bodyTypeLocationsKey)
		encoder.StartArray()
		for _, location := range b.locations {
			location.ToJSON(encoder, entity)
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}

// Update the role ranges and populate the lookup map.
func (b *BodyType) Update() {
	b.updateRollRanges()
	b.locationLookup = make(map[string]*HitLocation)
	b.populateMap(b.locationLookup)
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
	for _, location := range b.locations {
		start = location.updateRollRange(start)
	}
}

func (b *BodyType) populateMap(m map[string]*HitLocation) {
	for _, location := range b.locations {
		location.populateMap(m)
	}
}

// AddLocation adds a HitLocation to the end of list.
func (b *BodyType) AddLocation(loc *HitLocation) {
	b.locations = append(b.locations, loc)
	loc.owningTable = b
}

// RemoveLocation removes a HitLocation.
func (b *BodyType) RemoveLocation(loc *HitLocation) {
	for i, one := range b.locations {
		if one == loc {
			copy(b.locations[i:], b.locations[i+1:])
			b.locations[len(b.locations)-1] = nil
			b.locations = b.locations[:len(b.locations)-1]
			loc.owningTable = nil
		}
	}
}

// UniqueHitLocations returns the list of unique hit locations.
func (b *BodyType) UniqueHitLocations() []*HitLocation {
	if len(b.locationLookup) == 0 {
		b.Update()
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
func (b *BodyType) LookupLocationByID(idStr string) *HitLocation {
	if len(b.locationLookup) == 0 {
		b.Update()
	}
	return b.locationLookup[idStr]
}
