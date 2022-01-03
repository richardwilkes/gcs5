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
	"encoding/json"
	"testing"

	"github.com/richardwilkes/gcs/unit/length"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type embeddedLength struct {
	Field length.Length
}

func TestLengthConversion(t *testing.T) {
	assert.Equal(t, `1"`, length.FromInt64(1, length.Inch).Format(length.FeetAndInches))
	assert.Equal(t, `1'3"`, length.FromInt64(15, length.Inch).Format(length.FeetAndInches))
	assert.Equal(t, "2.5 cm", length.FromFloat64(2.5, length.Centimeter).Format(length.Centimeter))
	assert.Equal(t, "37.5 cm", length.FromFloat64(37.5, length.Centimeter).Format(length.Centimeter))

	w, err := length.FromString("1", length.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `1"`, w.Format(length.FeetAndInches))
	w, err = length.FromString(`6'2"`, length.Inch)
	assert.NoError(t, err)
	assert.Equal(t, `6'2"`, w.Format(length.FeetAndInches))
	w, err = length.FromString(" +32yd  ", length.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "96'", w.Format(length.FeetAndInches))
	w, err = length.FromString("0.5m", length.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "50 cm", w.Format(length.Centimeter))
	w, err = length.FromString("1cm", length.Inch)
	assert.NoError(t, err)
	assert.Equal(t, "1 cm", w.Format(length.Centimeter))
}

func TestHeightJSON(t *testing.T) {
	inc := length.FromFloat64(1.0/3.0, length.Inch)
	max := length.FromInt64(5, length.Inch)
	for i := length.Length(0); i <= max; i += inc {
		e1 := embeddedLength{Field: i}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embeddedLength
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
