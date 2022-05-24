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

package node

import "github.com/richardwilkes/unison"

// CellData holds data for creating a cell's visual representation.
type CellData struct {
	Type              CellType
	Disabled          bool
	Checked           bool
	Alignment         unison.Alignment
	Primary           string
	Secondary         string
	Tooltip           string
	UnsatisfiedReason string
}

// ForSort returns a string that can be used to sort or search against for this data.
func (c *CellData) ForSort() string {
	switch c.Type {
	case Text:
		if c.Secondary != "" {
			return c.Primary + "\n" + c.Secondary
		}
		return c.Primary
	case Toggle:
		if c.Checked {
			return "√"
		}
	case PageRef:
		return c.Primary
	}
	return ""
}
