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

package weight_test

import (
	"testing"

	"github.com/richardwilkes/gcs/model/enums/units"
	"github.com/richardwilkes/gcs/model/unit/weight"
	"github.com/stretchr/testify/assert"
)

func TestWeightConversion(t *testing.T) {
	assert.Equal(t, "1#", weight.FromInt64(1, units.Pound).Format(units.Pound))
	assert.Equal(t, "15#", weight.FromInt64(15, units.Pound).Format(units.Pound))
	assert.Equal(t, "0.5kg", weight.FromInt64(1, units.Pound).Format(units.Kilogram))
	assert.Equal(t, "7.5kg", weight.FromInt64(15, units.Pound).Format(units.Kilogram))

	w, err := weight.FromString("1", units.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "1#", w.String())
	w, err = weight.FromString("1", units.Kilogram)
	assert.NoError(t, err)
	assert.Equal(t, "2#", w.String())
	w, err = weight.FromString("22.34 lb", units.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34#", w.String())
	w, err = weight.FromString(" +22.34   lb  ", units.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34#", w.String())
	w, err = weight.FromString("0.5kg", units.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "0.5kg", w.Format(units.Kilogram))
	w, err = weight.FromString(" 15.25 kg ", units.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "15.25kg", w.Format(units.Kilogram))
}
