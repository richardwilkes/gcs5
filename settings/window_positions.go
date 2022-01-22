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
	"github.com/richardwilkes/toolbox/txt"
)

// WindowPositions holds the WindowPosition data for recently used windows.
type WindowPositions struct {
	Map map[string]*WindowPosition
}

// NewWindowPositions creates a new, empty, WindowPositions object.
func NewWindowPositions() *WindowPositions {
	return &WindowPositions{Map: make(map[string]*WindowPosition)}
}

// NewWindowPositionsFromJSON creates a new WindowPositions from a JSON object.
func NewWindowPositionsFromJSON(data map[string]interface{}) *WindowPositions {
	p := NewWindowPositions()
	for k, v := range data {
		p.Map[k] = NewWindowPositionFromJSON(encoding.Object(v))
	}
	return p
}

// ToKeyedJSON emits this object as JSON with the specified key, but only if not empty.
func (p *WindowPositions) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	if len(p.Map) != 0 {
		encoder.Key(key)
		p.ToJSON(encoder)
	}
}

// ToJSON emits this object as JSON.
func (p *WindowPositions) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	keys := make([]string, 0, len(p.Map))
	for k := range p.Map {
		keys = append(keys, k)
	}
	txt.SortStringsNaturalAscending(keys)
	for _, k := range keys {
		encoder.Key(k)
		p.Map[k].ToJSON(encoder)
	}
	encoder.EndObject()
}
