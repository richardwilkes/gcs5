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

package criteria

// String holds the criteria for matching a string.
type String struct {
	Compare   StringCompareType `json:"compare,omitempty"`
	Qualifier string            `json:"qualifier,omitempty"`
}

// Normalize the data.
func (s *String) Normalize() {
	if s.Compare = s.Compare.EnsureValid(); s.Compare == Any {
		s.Qualifier = ""
	}
}
