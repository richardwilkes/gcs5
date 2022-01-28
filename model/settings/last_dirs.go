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
	"github.com/richardwilkes/gcs/model/encoding"
)

// LastDirs holds the last directory used for a given key.
type LastDirs struct {
	Map map[string]string
}

// NewLastDirs creates a new, empty, LastDirs.
func NewLastDirs() *LastDirs {
	return &LastDirs{Map: make(map[string]string)}
}

// NewLastDirsFromJSON creates a new LastDirs from a JSON object.
func NewLastDirsFromJSON(data map[string]interface{}) *LastDirs {
	d := NewLastDirs()
	for k, v := range data {
		d.Map[k] = encoding.String(v)
	}
	return d
}

// Empty implements encoding.Empty.
func (d *LastDirs) Empty() bool {
	return len(d.Map) == 0
}

// ToJSON emits this object as JSON.
func (d *LastDirs) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	for k, v := range d.Map {
		encoder.KeyedString(k, v, false, false)
	}
	encoder.EndObject()
}