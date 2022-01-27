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

package paper

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Orientation values.
const (
	Portrait Orientation = iota
	Landscape
)

type orientationData struct {
	Key        string
	String     string
	Dimensions func(size Size) (width, height Length)
}

// Orientation holds the orientation of the page.
type Orientation uint8

var orientationValues = []*orientationData{
	{
		Key:        "portrait",
		String:     i18n.Text("Portrait"),
		Dimensions: func(size Size) (width, height Length) { return size.Dimensions() },
	},
	{
		Key:    "landscape",
		String: i18n.Text("Landscape"),
		Dimensions: func(size Size) (width, height Length) {
			width, height = size.Dimensions()
			return height, width
		},
	},
}

// OrientationFromString extracts a Orientation from a key.
func OrientationFromString(key string) Orientation {
	for i, one := range orientationValues {
		if strings.EqualFold(key, one.Key) {
			return Orientation(i)
		}
	}
	return 0
}

// EnsureValid returns the first Orientation if this Orientation is not a known value.
func (o Orientation) EnsureValid() Orientation {
	if int(o) < len(orientationValues) {
		return o
	}
	return 0
}

// Key returns the key used to represent this ThresholdOp.
func (o Orientation) Key() string {
	return orientationValues[o.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (o Orientation) String() string {
	return orientationValues[o.EnsureValid()].String
}

// Dimensions returns the paper dimensions after orienting the paper.
func (o Orientation) Dimensions(size Size) (width, height Length) {
	return orientationValues[o.EnsureValid()].Dimensions(size)
}
