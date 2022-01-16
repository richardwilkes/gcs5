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

package pdf

type params struct {
	sequence   int
	pageNumber int
	search     string
	scale      float32
}

func (p *params) sameAs(pageNumber int, scale float32, search string) bool {
	return p.pageNumber == pageNumber && p.scale == scale && p.search == search
}
