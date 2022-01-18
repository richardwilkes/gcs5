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
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

// The data format versions for BodyType.
const (
	BodyTypeCurrentVersion = 3
	BodyTypeJavaVersion    = 2
)

// BodyType holds a set of hit locations.
type BodyType struct {
	BodyTypeStorage
	owningLocation *HitLocation
}

// BodyTypeStorage defines the current BodyType data format.
type BodyTypeStorage struct {
	Version   int            `json:"version"`
	ID        string         `json:"id"`
	Name      string         `json:"name,omitempty"`
	Roll      dice.Dice      `json:"roll"`
	Locations []*HitLocation `json:"locations,omitempty"`
}

// FactoryBodyType returns a new copy of the factory BodyType.
func FactoryBodyType() *BodyType {
	b, err := LoadBodyType(embeddedFS, "data/body_types/Humanoid.body")
	jot.FatalIfErr(err)
	return b
}

// LoadBodyType creates a BodyType from a file.
func LoadBodyType(fileSystem fs.FS, filePath string) (*BodyType, error) {
	f, err := fileSystem.Open(filePath)
	if err != nil {
		return nil, errs.NewWithCause("unable to open body type file", err)
	}
	defer xio.CloseIgnoringErrors(f)
	var bt BodyType
	if err = json.NewDecoder(f).Decode(&bt); err != nil {
		return nil, errs.NewWithCause("invalid body type file: "+filePath, err)
	}
	return &bt, nil
}

// MarshalJSON implements json.Marshaler. Sets the current version on the output.
func (b *BodyType) MarshalJSON() ([]byte, error) {
	b.Version = BodyTypeCurrentVersion
	return json.Marshal(&b.BodyTypeStorage)
}

// UnmarshalJSON implements json.Unmarshaler. Loads the current format as well as older variants.
func (b *BodyType) UnmarshalJSON(data []byte) error {
	var variants struct {
		BodyTypeStorage `json:",inline"` // v3+, except for the Version field, which is v0+
		Type            string           `json:"type"`          // v0-2
		HitLocations    BodyTypeStorage  `json:"hit_locations"` // v0-2
	}
	if err := json.Unmarshal(data, &variants); err != nil {
		return err
	}
	if variants.Version <= BodyTypeJavaVersion && variants.Type == "hit_locations" {
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
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		b.FillCalc(entity)
		defer b.ClearCalc()
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(b)
	}, 0o640); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return nil
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

// FillCalc fills in the calculation fields for third parties. 'entity' may be nil.
func (b *BodyType) FillCalc(entity *Entity) {
	b.ClearCalc()
	b.updateRollRanges()
	b.updateDR(entity)
}

// ClearCalc clears the calculation fields.
func (b *BodyType) ClearCalc() {
	for _, loc := range b.Locations {
		loc.clearCalc()
	}
}
