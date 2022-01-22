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
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	pageRefPathKey   = "path"
	pageRefOffsetKey = "offset"
)

// PageRef holds a path to a file and an offset for all page references within that file.
type PageRef struct {
	ID     string
	Path   string
	Offset int
}

// NewPageRefFromJSON creates a new PageRef from a JSON object.
func NewPageRefFromJSON(id string, data map[string]interface{}) *PageRef {
	return &PageRef{
		ID:     id,
		Path:   encoding.String(data[pageRefPathKey]),
		Offset: int(encoding.Number(data[pageRefOffsetKey]).AsInt64()),
	}
}

// ToJSON emits this object as JSON.
func (p *PageRef) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(pageRefPathKey, p.Path, false, false)
	encoder.KeyedNumber(pageRefOffsetKey, fixed.F64d4FromInt64(int64(p.Offset)), true)
	encoder.EndObject()
}
