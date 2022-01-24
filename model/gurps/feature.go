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

import "github.com/richardwilkes/gcs/model/encoding"

const featureTypeKey = "type"

// Feature defines the methods that all features must have.
type Feature interface {
	// ToJSON emits this Feature as JSON.
	ToJSON(encoder *encoding.JSONEncoder)
	// CloneFeature creates a clone of this Feature.
	CloneFeature() Feature
	// DataType returns the data type key for this Feature.
	DataType() string
	// FeatureKey returns the key used in the Feature map for things this Feature applies to.
	FeatureKey() string
	// FillWithNameableKeys adds any nameable keys found in this Feature to the provided map.
	FillWithNameableKeys(set map[string]bool)
	// ApplyNameableKeys replaces any nameable keys found in this Feature with the corresponding values in the provided map.
	ApplyNameableKeys(nameables map[string]string)
}
