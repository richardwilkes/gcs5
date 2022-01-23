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

package units

// Possible RealLength values.
const (
	RealMillimeter RealLength = iota
	RealCentimeter
	RealInch
)

// RealLength holds the real-world length unit type.
type RealLength uint8

// Key returns the key used to represent this GURPSLength.
func (l RealLength) Key() string {
	switch l {
	case RealMillimeter:
		return "mm"
	case RealCentimeter:
		return "cm"
	default: // RealInch
		return "in"
	}
}
