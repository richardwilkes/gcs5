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

package gurps

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps/feature"
)

type featureMapKeyer interface {
	featureMapKey() string
}

type nameableKeysProvider interface {
	fillWithNameableKeys(nameables map[string]string)
	applyNameableKeys(nameables map[string]string)
}

// Feature holds data that affects another object.
type Feature struct {
	Type   feature.Type `json:"type"`
	Self   interface{}  `json:"-"`
	Parent fmt.Stringer `json:"-"`
}

// FeatureMapKey returns the key used for matching within the feature map.
func (f *Feature) FeatureMapKey() string {
	if t, ok := f.Self.(featureMapKeyer); ok {
		return t.featureMapKey()
	}
	return string(f.Type)
}

// FillWithNameableKeys fills the map with nameable keys.
func (f *Feature) FillWithNameableKeys(nameables map[string]string) {
	if t, ok := f.Self.(nameableKeysProvider); ok {
		t.fillWithNameableKeys(nameables)
	}
}

// ApplyNameableKeys applies the nameable keys to this object.
func (f *Feature) ApplyNameableKeys(nameables map[string]string) {
	if t, ok := f.Self.(nameableKeysProvider); ok {
		t.applyNameableKeys(nameables)
	}
}
