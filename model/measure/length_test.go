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

	"github.com/richardwilkes/gcs/model/measure"
	"github.com/stretchr/testify/assert"
)

func TestRealLengthConversion(t *testing.T) {
	assert.Equal(t, `0.25in`, measure.Length{Length: 0.25, Units: measure.Inch}.String())
	assert.Equal(t, float32(18), measure.Length{Length: 0.25, Units: measure.Inch}.Pixels())
	assert.Equal(t, `1in`, measure.Length{Length: 1, Units: measure.Inch}.String())
	assert.Equal(t, float32(72), measure.Length{Length: 1, Units: measure.Inch}.Pixels())
	assert.Equal(t, `15in`, measure.Length{Length: 15, Units: measure.Inch}.String())
	assert.Equal(t, float32(1080), measure.Length{Length: 15, Units: measure.Inch}.Pixels())
	assert.Equal(t, "1cm", measure.Length{Length: 1, Units: measure.Centimeter}.String())
	assert.Equal(t, float32(28.3464566929), measure.Length{Length: 1, Units: measure.Centimeter}.Pixels())
	assert.Equal(t, "1mm", measure.Length{Length: 1, Units: measure.Millimeter}.String())
	assert.Equal(t, float32(2.8346456693), measure.Length{Length: 1, Units: measure.Millimeter}.Pixels())
}
