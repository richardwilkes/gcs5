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

type LastDirs struct {
	Map map[string]string
}

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

// ToKeyedJSON emits this object as JSON with the specified key, but only if not empty.
func (d *LastDirs) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	if len(d.Map) != 0 {
		encoder.Key(key)
		d.ToJSON(encoder)
	}
}

// ToJSON emits this object as JSON.
func (d *LastDirs) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	for k, v := range d.Map {
		encoder.KeyedString(k, v, false, false)
	}
	encoder.EndObject()
}
