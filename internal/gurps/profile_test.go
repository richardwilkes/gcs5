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
	"encoding/json"
	"testing"

	"github.com/richardwilkes/gcs/internal/gurps"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed profile_test.json
var profileTestSample []byte

func TestProfileImage(t *testing.T) {
	var profile gurps.Profile
	assert.NoError(t, json.Unmarshal(profileTestSample, &profile))
	img := profile.Portrait()
	require.NotNil(t, img)
	assert.Equal(t, geom32.NewSize(gurps.PortraitWidth, gurps.PortraitHeight), img.LogicalSize())

	data, err := json.MarshalIndent(&profile, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, profileTestSample, data)
}
