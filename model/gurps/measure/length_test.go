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

package measure_test

import (
	"testing"

	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/stretchr/testify/assert"
)

func TestGURPSLengthConversion(t *testing.T) {
	assert.Equal(t, `1"`, measure.LengthFromInt64(1, measure.Inch).Format(measure.FeetAndInches))
	assert.Equal(t, `1'3"`, measure.LengthFromInt64(15, measure.Inch).Format(measure.FeetAndInches))
	assert.Equal(t, "2.5cm", measure.LengthFromStringForced("2.5", measure.Centimeter).Format(measure.Centimeter))
	assert.Equal(t, "37.5cm", measure.LengthFromStringForced("37.5", measure.Centimeter).Format(measure.Centimeter))

	w, err := measure.LengthFromString("1", measure.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `1"`, w.Format(measure.FeetAndInches))
	w, err = measure.LengthFromString(`6'         2"`, measure.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `6'2"`, w.Format(measure.FeetAndInches))
	w, err = measure.LengthFromString(" +32   yd  ", measure.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "96'", w.Format(measure.FeetAndInches))
	w, err = measure.LengthFromString("0.5m", measure.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "50cm", w.Format(measure.Centimeter))
	w, err = measure.LengthFromString("1cm", measure.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "1cm", w.Format(measure.Centimeter))
}
