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

package tbl

type ColumnMapper struct {
	list map[int]int
	page map[int]int
}

func (m *ColumnMapper) Map(index int, forPage bool) int {
	var exists bool
	if forPage {
		index, exists = m.page[index]
	} else {
		index, exists = m.list[index]
	}
	if exists {
		return index
	}
	return -1
}

func (m *ColumnMapper) Size(forPage bool) int {
	if forPage {
		return len(m.page)
	}
	return len(m.list)
}
