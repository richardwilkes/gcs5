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
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
)

// BodyType holds a set of hit locations.
type BodyType struct {
	BodyTypeStorage `json:",inline"`
	owningLocation  *HitLocation
}

// BodyTypeStorage defines the current BodyType data format.
type BodyTypeStorage struct {
	Name      string         `json:"name"`
	Roll      dice.Dice      `json:"roll"`
	Locations []*HitLocation `json:"locations,omitempty"`
}

// FactoryBodyType returns a new copy of the default factory BodyType.
func FactoryBodyType() *BodyType {
	var bodyType BodyType
	jot.FatalIfErr(xfs.LoadJSONFromFS(embeddedFS, "data/body_types/Humanoid.body", &bodyType))
	return &bodyType
}

// FactoryBodyTypes returns the list of the known factory BodyTypes.
func FactoryBodyTypes() []*BodyType {
	entries, err := embeddedFS.ReadDir("data/body_types")
	jot.FatalIfErr(err)
	list := make([]*BodyType, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if path.Ext(name) == ".body" {
			var bodyType BodyType
			jot.FatalIfErr(xfs.LoadJSONFromFS(embeddedFS, "data/body_types/"+name, &bodyType))
			list = append(list, &bodyType)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return txt.NaturalLess(list[i].Name, list[j].Name, true)
	})
	return list
}

// UnmarshalJSON implements json.Unmarshaler. Loads the current format as well as older variants.
func (b *BodyType) UnmarshalJSON(data []byte) error {
	var variants struct {
		BodyTypeStorage `json:",inline"`
		Type            string          `json:"type"`
		HitLocations    BodyTypeStorage `json:"hit_locations"`
	}
	if err := json.Unmarshal(data, &variants); err != nil {
		return err
	}
	if variants.Type == "hit_locations" {
		b.BodyTypeStorage = variants.HitLocations
	} else {
		b.BodyTypeStorage = variants.BodyTypeStorage
	}
	return nil
}

// SaveTo saves the BodyType data to the specified file. 'entity' may be nil.
func (b *BodyType) SaveTo(filePath string, entity *Entity) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	b.calc(entity)
	return xfs.SaveJSONWithMode(filePath, b, true, 0o640)
}

func (b *BodyType) calc(entity *Entity) {
	b.updateRollRanges()
	b.updateDR(entity)
}

func (b *BodyType) updateRollRanges() {
	start := b.Roll.Minimum(false)
	for _, location := range b.Locations {
		start = location.updateRollRange(start)
	}
}

func (b *BodyType) updateDR(entity *Entity) {
	for _, location := range b.Locations {
		location.updateDR(entity)
	}
}
