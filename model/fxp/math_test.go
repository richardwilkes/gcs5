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

package fxp_test

import (
	"testing"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

func TestMod(t *testing.T) {
	twoAndAHalf := fixed.F64d4FromStringForced("2.5")
	r := fxp.Mod(fixed.F64d4FromStringForced("12"), twoAndAHalf)
	if r != fxp.Two {
		t.Errorf("%v != %v", r, fxp.Two)
	}

	r = fxp.Mod(fixed.F64d4FromStringForced("21.5"), twoAndAHalf)
	if r != fxp.OneAndAHalf {
		t.Errorf("%v != %v", r, fxp.OneAndAHalf)
	}

	pointTwo := fixed.F64d4FromStringForced("0.2")
	pointThree := fixed.F64d4FromStringForced("0.3")
	r = fxp.Mod(fxp.Half, pointThree)
	if r != pointTwo {
		t.Errorf("%v != %v", r, pointTwo)
	}

	negPointTwo := fixed.F64d4FromStringForced("-0.2")
	r = fxp.Mod(fxp.NegHalf, pointThree)
	if r != negPointTwo {
		t.Errorf("%v != %v", r, negPointTwo)
	}
}
