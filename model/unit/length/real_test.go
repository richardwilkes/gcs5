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

func TestRealLengthConversion(t *testing.T) {
	assert.Equal(t, `0.25in`, length.Real{Length: 0.25, Units: units.RealInch}.String())
	assert.Equal(t, float32(18), length.Real{Length: 0.25, Units: units.RealInch}.Pixels())
	assert.Equal(t, `1in`, length.Real{Length: 1, Units: units.RealInch}.String())
	assert.Equal(t, float32(72), length.Real{Length: 1, Units: units.RealInch}.Pixels())
	assert.Equal(t, `15in`, length.Real{Length: 15, Units: units.RealInch}.String())
	assert.Equal(t, float32(1080), length.Real{Length: 15, Units: units.RealInch}.Pixels())
	assert.Equal(t, "1cm", length.Real{Length: 1, Units: units.RealCentimeter}.String())
	assert.Equal(t, float32(28.3464566929), length.Real{Length: 1, Units: units.RealCentimeter}.Pixels())
	assert.Equal(t, "1mm", length.Real{Length: 1, Units: units.RealMillimeter}.String())
	assert.Equal(t, float32(2.8346456693), length.Real{Length: 1, Units: units.RealMillimeter}.Pixels())
}
