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
	"time"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
)

const (
	windowPositionXKey           = "x"
	windowPositionYKey           = "y"
	windowPositionWidthKey       = "width"
	windowPositionHeightKey      = "height"
	windowPositionLastUpdatedKey = "last_updated"
)

// WindowPosition holds a window's last known frame and when the frame's size or position was last altered.
type WindowPosition struct {
	Frame       geom32.Rect
	LastUpdated time.Time
}

// NewWindowPositionFromJSON creates a new WindowPosition from a JSON object.
func NewWindowPositionFromJSON(data map[string]interface{}) *WindowPosition {
	p := &WindowPosition{
		Frame: geom32.NewRect(float32(encoding.Number(data[windowPositionXKey]).AsFloat64()),
			float32(encoding.Number(data[windowPositionYKey]).AsFloat64()),
			float32(encoding.Number(data[windowPositionWidthKey]).AsFloat64()),
			float32(encoding.Number(data[windowPositionHeightKey]).AsFloat64())),
	}
	if err := p.LastUpdated.UnmarshalText([]byte(encoding.String(data[windowPositionLastUpdatedKey]))); err != nil {
		p.LastUpdated = time.Now()
	}
	return p
}

// ToJSON emits this object as JSON.
func (p *WindowPosition) ToJSON(encoder *encoding.JSONEncoder) {
	lastUpdate, _ := p.LastUpdated.MarshalText() //nolint:errcheck // An empty string is ok on error
	encoder.StartObject()
	encoder.KeyedNumber(windowPositionXKey, fixed.F64d4FromFloat64(float64(p.Frame.X)), false)
	encoder.KeyedNumber(windowPositionYKey, fixed.F64d4FromFloat64(float64(p.Frame.Y)), false)
	encoder.KeyedNumber(windowPositionWidthKey, fixed.F64d4FromFloat64(float64(p.Frame.Width)), false)
	encoder.KeyedNumber(windowPositionHeightKey, fixed.F64d4FromFloat64(float64(p.Frame.Height)), false)
	encoder.KeyedString(windowPositionLastUpdatedKey, string(lastUpdate), true, true)
	encoder.EndObject()
}
