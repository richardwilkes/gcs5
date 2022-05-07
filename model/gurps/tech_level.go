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
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// ExtractTechLevel extracts the first number it finds in the string and returns that as the tech level. The start and
// end (inclusive) indexes within the string where the number resided are returned, but will be -1 if the string didn't
// contain a resolvable number. The returned tech level will be clamped to the range 0 to 12.
func ExtractTechLevel(str string) (techLevel fxp.Int, start, end int) {
	var buffer strings.Builder
	decimal := true
	looking := true
outer:
	for i, ch := range str {
		isDigit := ch >= '0' && ch <= '9'
		switch {
		case looking:
			if isDigit {
				start = i
				end = i
				buffer.WriteRune(ch)
				looking = false
			}
		case isDigit:
			end = i
			buffer.WriteRune(ch)
		case decimal && ch == '.':
			end = i
			buffer.WriteRune(ch)
			decimal = false
		default:
			break outer
		}
	}
	if buffer.Len() == 0 {
		return 0, -1, -1
	}
	var err error
	if techLevel, err = f64.FromString[fxp.DP](buffer.String()); err == nil {
		return techLevel.Max(0).Min(fxp.Twelve), start, end
	}
	return 0, -1, -1
}

// ReplaceTechLevel replaces the tech level (as found by a call to ExtractTechLevel) with a new value.
func ReplaceTechLevel(str string, value fxp.Int) string {
	if _, start, end := ExtractTechLevel(str); start != -1 {
		var buffer strings.Builder
		if start > 0 {
			buffer.WriteString(str[:start])
		}
		buffer.WriteString(value.String())
		if end < len(str)-1 {
			buffer.WriteString(str[end+1:])
		}
		str = buffer.String()
	}
	return str
}

// AdjustTechLevel returns a new string with the adjusted tech level.
func AdjustTechLevel(str string, delta fxp.Int) (s string, changed bool) {
	tl, start, _ := ExtractTechLevel(str)
	newTL := (tl + delta).Max(0).Min(fxp.Twelve)
	if tl == newTL {
		return str, false
	}
	if start == -1 {
		return newTL.String() + str, true
	}
	return ReplaceTechLevel(str, newTL), true
}
