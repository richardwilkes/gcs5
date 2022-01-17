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
	"path"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
	"gopkg.in/yaml.v3"
)

// BodyTypeCurrentVersion holds the current data format version.
const BodyTypeCurrentVersion = 2

type bodyTypeFile struct {
	Type     string    `json:"type" yaml:"type"`
	Version  int       `json:"version" yaml:"version"`
	BodyType *BodyType `json:"hit_locations" yaml:"hit_locations"`
}

// BodyType holds a set of hit locations.
type BodyType struct {
	ID             string         `json:"id" yaml:"id"`
	Name           string         `json:"name" yaml:"name"`
	Roll           dice.Dice      `json:"roll" yaml:"roll"`
	Locations      []*HitLocation `json:"locations" yaml:"locations,omitempty"`
	owningLocation *HitLocation
}

// FactoryBodyType returns a new copy of the factory BodyType.
func FactoryBodyType() *BodyType {
	b, err := LoadBodyType(embeddedFS, "embedded/body_types/Humanoid.yaml")
	jot.FatalIfErr(err)
	return b
}

// LoadBodyType creates a BodyType from a file.
func LoadBodyType(fsys fs.FS, filePath string) (*BodyType, error) {
	f, err := fsys.Open(filePath)
	if err != nil {
		return nil, errs.NewWithCause("unable to open body type file", err)
	}
	defer xio.CloseIgnoringErrors(f)
	var content bodyTypeFile
	switch strings.ToLower(path.Ext(filePath)) {
	case ".json", ".ghl":
		if err = json.NewDecoder(f).Decode(&content); err != nil {
			return nil, errs.NewWithCause("unable to read body type file: "+filePath, err)
		}
	case ".yaml":
		if err = yaml.NewDecoder(f).Decode(&content); err != nil {
			return nil, errs.NewWithCause("unable to read body type file: "+filePath, err)
		}
	default:
		return nil, errs.New("unexpected file extension: " + filePath)
	}
	if content.BodyType == nil {
		return nil, errs.New("no hit locations in file: " + filePath)
	}
	return content.BodyType, nil
}

// SaveTo saves the BodyType data to the specified file. If 'calc' is true, then the calculation fields for third parties
// will be filled in. 'entity' may be nil and is only used if 'calc' is true.
func (h *BodyType) SaveTo(filePath string, calc bool, entity *Entity) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		if calc {
			h.FillCalc(entity)
			defer h.ClearCalc()
		} else {
			h.ClearCalc()
		}
		content := &bodyTypeFile{
			Type:     "hit_locations",
			Version:  BodyTypeCurrentVersion,
			BodyType: h,
		}
		switch strings.ToLower(path.Ext(filePath)) {
		case ".json", ".ghl":
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			return encoder.Encode(&content)
		case ".yaml":
			encoder := yaml.NewEncoder(w)
			encoder.SetIndent(2)
			if e := encoder.Encode(&content); e != nil {
				return e
			}
			return encoder.Close()
		default:
			return errs.New("unexpected file extension")
		}
	}, 0o640); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return nil
}

func (h *BodyType) updateRollRanges() {
	start := h.Roll.Minimum(false)
	for _, location := range h.Locations {
		start = location.updateRollRange(start)
	}
}

func (h *BodyType) updateDR(entity *Entity) {
	for _, location := range h.Locations {
		location.updateDR(entity)
	}
}

// FillCalc fills in the calculation fields for third parties. 'entity' may be nil.
func (h *BodyType) FillCalc(entity *Entity) {
	h.ClearCalc()
	h.updateRollRanges()
	h.updateDR(entity)
}

// ClearCalc clears the calculation fields.
func (h *BodyType) ClearCalc() {
	for _, loc := range h.Locations {
		loc.clearCalc()
	}
}
