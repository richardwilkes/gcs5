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
	"strings"

	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Fraction holds a fraction value.
type Fraction struct {
	Numerator   f64d4.Int
	Denominator f64d4.Int
}

// NewFractionFromString creates a new fractional value from a string.
func NewFractionFromString(s string) Fraction {
	parts := strings.SplitN(s, "/", 2)
	f := Fraction{
		Numerator:   f64d4.FromStringForced(strings.TrimSpace(parts[0])),
		Denominator: f64d4.One,
	}
	if len(parts) > 1 {
		f.Denominator = f64d4.FromStringForced(strings.TrimSpace(parts[1]))
	}
	return f
}

// Normalize the fraction, eliminating any division by zero.
func (f *Fraction) Normalize() {
	if f.Denominator == 0 {
		f.Numerator = 0
		f.Denominator = f64d4.One
	} else if f.Denominator < 0 {
		f.Numerator = f.Numerator.Mul(NegOne)
		f.Denominator = f.Denominator.Mul(NegOne)
	}
}

// Value returns the computed value.
func (f Fraction) Value() f64d4.Int {
	return f.Numerator.Div(f.Denominator)
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive
func (f Fraction) StringWithSign() string {
	s := f.Numerator.StringWithSign()
	if f.Denominator == f64d4.One {
		return s
	}
	return s + "/" + f.Denominator.String()
}

func (f Fraction) String() string {
	s := f.Numerator.String()
	if f.Denominator == f64d4.One {
		return s
	}
	return s + "/" + f.Denominator.String()
}

// MarshalJSON implements json.Marshaler.
func (f Fraction) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Fraction) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}
	*f = NewFractionFromString(s)
	return nil
}
