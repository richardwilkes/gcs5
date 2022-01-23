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

package length_test

import (
	"testing"

	"github.com/richardwilkes/gcs/model/enums/units"
	"github.com/richardwilkes/gcs/model/unit/length"
	"github.com/stretchr/testify/assert"
)

func TestGURPSLengthConversion(t *testing.T) {
	assert.Equal(t, `1"`, length.GURPSFromInt64(1, units.Inch).Format(units.FeetAndInches))
	assert.Equal(t, `1'3"`, length.GURPSFromInt64(15, units.Inch).Format(units.FeetAndInches))
	assert.Equal(t, "2.5cm", length.GURPSFromFloat64(2.5, units.Centimeter).Format(units.Centimeter))
	assert.Equal(t, "37.5cm", length.GURPSFromFloat64(37.5, units.Centimeter).Format(units.Centimeter))

	w, err := length.GURPSFromString("1", units.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `1"`, w.Format(units.FeetAndInches))
	w, err = length.GURPSFromString(`6'         2"`, units.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `6'2"`, w.Format(units.FeetAndInches))
	w, err = length.GURPSFromString(" +32   yd  ", units.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "96'", w.Format(units.FeetAndInches))
	w, err = length.GURPSFromString("0.5m", units.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "50cm", w.Format(units.Centimeter))
	w, err = length.GURPSFromString("1cm", units.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "1cm", w.Format(units.Centimeter))
}
