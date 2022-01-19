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

package gurps_test

import (
	_ "embed"
	"testing"

	"github.com/goccy/go-json"
	"github.com/richardwilkes/gcs/gurps"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed profile_test.json
var profileTestSample []byte

func TestProfileImage(t *testing.T) {
	var profile gurps.PCProfile
	assert.NoError(t, json.Unmarshal(profileTestSample, &profile))
	assert.NotEmpty(t, profile.Name)
	assert.NotEmpty(t, profile.TechLevel)
	assert.NotEmpty(t, profile.SizeModifier)
	assert.NotEmpty(t, profile.PortraitData)
	assert.NotEmpty(t, profile.Title)
	assert.NotEmpty(t, profile.Organization)
	assert.NotEmpty(t, profile.Religion)
	assert.NotEmpty(t, profile.Age)
	assert.NotEmpty(t, profile.Eyes)
	assert.NotEmpty(t, profile.Hair)
	assert.NotEmpty(t, profile.Skin)
	assert.NotEmpty(t, profile.Handedness)
	assert.NotEmpty(t, profile.Gender)
	assert.NotEmpty(t, profile.Height)
	assert.NotEmpty(t, profile.Weight)
	assert.NotEmpty(t, profile.PlayerName)
	assert.NotEmpty(t, profile.Birthday)

	data, err := json.MarshalIndent(&profile, "", "  ")
	assert.NoError(t, err)

	var roundTripProfile gurps.PCProfile
	assert.NoError(t, json.Unmarshal(data, &roundTripProfile))
	assert.Equal(t, profile, roundTripProfile)

	img := profile.Portrait()
	require.NotNil(t, img)
	assert.Equal(t, geom32.NewSize(gurps.PortraitWidth, gurps.PortraitHeight), img.LogicalSize())
}
