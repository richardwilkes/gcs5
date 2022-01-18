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
	"encoding/json"
	"testing"

	"github.com/richardwilkes/gcs/unit/weight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type embeddedWeight struct {
	Field weight.Weight
}

func TestWeightConversion(t *testing.T) {
	assert.Equal(t, "1#", weight.FromInt64(1, weight.Pound).Format(weight.Pound))
	assert.Equal(t, "15#", weight.FromInt64(15, weight.Pound).Format(weight.Pound))
	assert.Equal(t, "0.5kg", weight.FromInt64(1, weight.Pound).Format(weight.Kilogram))
	assert.Equal(t, "7.5kg", weight.FromInt64(15, weight.Pound).Format(weight.Kilogram))

	w, err := weight.FromString("1", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "1#", w.String())
	w, err = weight.FromString("1", weight.Kilogram)
	assert.NoError(t, err)
	assert.Equal(t, "2#", w.String())
	w, err = weight.FromString("22.34 lb", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34#", w.String())
	w, err = weight.FromString(" +22.34   lb  ", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34#", w.String())
	w, err = weight.FromString("0.5kg", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "0.5kg", w.Format(weight.Kilogram))
	w, err = weight.FromString(" 15.25 kg ", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "15.25kg", w.Format(weight.Kilogram))
}

func TestWeightJSON(t *testing.T) {
	inc := weight.FromFloat64(1.0/3.0, weight.Pound)
	max := weight.FromInt64(5, weight.Pound)
	for i := weight.Weight(0); i <= max; i += inc {
		e1 := embeddedWeight{Field: i}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embeddedWeight
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
