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

package fxp

import (
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// DP is an alias for the fixed-point decimal places configuration we are using.
type DP = fixed.D4

// Int is an alias for the fixed-point type we are using.
type Int = f64.Int[DP]

// Common values that can be reused.
var (
	Min           = Int(f64.Min)
	Max           = Int(f64.Max)
	NegEighty     = f64.From[DP](-80)
	NegFour       = f64.From[DP](-4)
	NegThree      = f64.From[DP](-3)
	NegTwo        = f64.From[DP](-2)
	NegOne        = f64.From[DP](-1)
	NegPointEight = f64.FromStringForced[DP]("-0.8")
	Half          = f64.FromStringForced[DP]("0.5")
	One           = f64.From[DP](1)
	OneAndAHalf   = f64.FromStringForced[DP]("1.5")
	Two           = f64.From[DP](2)
	Three         = f64.From[DP](3)
	Four          = f64.From[DP](4)
	Five          = f64.From[DP](5)
	Six           = f64.From[DP](6)
	Seven         = f64.From[DP](7)
	Eight         = f64.From[DP](8)
	Nine          = f64.From[DP](9)
	Ten           = f64.From[DP](10)
	Twelve        = f64.From[DP](12)
	Fifteen       = f64.From[DP](15)
	Nineteen      = f64.From[DP](19)
	Twenty        = f64.From[DP](20)
	TwentyFour    = f64.From[DP](24)
	ThirtySix     = f64.From[DP](36)
	Thirty        = f64.From[DP](30)
	Forty         = f64.From[DP](40)
	Fifty         = f64.From[DP](50)
	Seventy       = f64.From[DP](70)
	Eighty        = f64.From[DP](80)
	NinetyNine    = f64.From[DP](99)
	Hundred       = f64.From[DP](100)
	Thousand      = f64.From[DP](1000)
	MaxBasePoints = f64.From[DP](999999)
)
