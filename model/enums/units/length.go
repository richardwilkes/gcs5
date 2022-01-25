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

// Possible Length values.
const (
	Millimeter Length = iota
	Centimeter
	Inch
)

// Length holds the real-world length unit type.
type Length uint8

// Key returns the key used to represent this GURPSLength.
func (l Length) Key() string {
	switch l {
	case Millimeter:
		return "mm"
	case Centimeter:
		return "cm"
	default: // RealInch
		return "in"
	}
}
