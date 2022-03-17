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

type CellData struct {
	Type      CellType
	Checked   bool
	Alignment unison.Alignment
	Primary   string
	Secondary string
}

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
