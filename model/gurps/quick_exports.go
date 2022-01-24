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
	"sort"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	quickExportsMaxKey     = "max"
	quickExportsExportsKey = "exports"
)

// QuickExports holds a list containing information about previous exports.
type QuickExports struct {
	max  int
	info []*ExportInfo
}

// NewQuickExports creates a new, empty, QuickExports object.
func NewQuickExports() *QuickExports {
	return &QuickExports{max: 20}
}

// NewQuickExportsFromJSON creates a new QuickExports from a JSON object.
func NewQuickExportsFromJSON(data map[string]interface{}) *QuickExports {
	q := NewQuickExports()
	if v, ok := data[quickExportsMaxKey]; ok {
		q.max = xmath.MaxInt(int(encoding.Number(v).AsInt64()), 0)
	} else {
		q.max = 20
	}
	if q.max > 0 {
		array := encoding.Array(data[quickExportsExportsKey])
		count := xmath.MinInt(len(array), q.max)
		for i, one := range array {
			if i == count {
				break
			}
			q.info = append(q.info, NewExportInfoFromJSON(encoding.Object(one)))
		}
	}
	return q
}

// Empty implements encoding.Empty.
func (q *QuickExports) Empty() bool {
	return len(q.info) == 0
}

// ToJSON emits this object as JSON.
func (q *QuickExports) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedNumber(quickExportsMaxKey, fixed.F64d4FromInt64(int64(q.max)), false)
	sort.Slice(q.info, func(i, j int) bool { return q.info[i].LastUsed > q.info[j].LastUsed })
	if len(q.info) > q.max {
		list := make([]*ExportInfo, q.max)
		copy(list, q.info)
		q.info = list
	}
	if len(q.info) != 0 {
		encoder.Key(quickExportsExportsKey)
		encoder.StartArray()
		for _, data := range q.info {
			data.ToJSON(encoder)
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}
