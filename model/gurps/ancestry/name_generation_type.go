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

package ancestry

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible NameGenerationType values.
const (
	Simple NameGenerationType = iota
	MarkovChain
)

type nameGenerationTypeData struct {
	Key    string
	String string
}

// NameGenerationType holds the type of a name generation technique.
type NameGenerationType uint8

var nameGenerationTypeValues = []*nameGenerationTypeData{
	{
		Key:    "simple",
		String: i18n.Text("Simple"),
	},
	{
		Key:    "markov_chain",
		String: i18n.Text("Markov Chain"),
	},
}

// NameGenerationTypeFromKey extracts a NameGenerationType from a key.
func NameGenerationTypeFromKey(key string) NameGenerationType {
	for i, one := range nameGenerationTypeValues {
		if strings.EqualFold(key, one.Key) {
			return NameGenerationType(i)
		}
	}
	return 0
}

// EnsureValid returns the first NameGenerationType if this NameGenerationType is not a known value.
func (a NameGenerationType) EnsureValid() NameGenerationType {
	if int(a) < len(nameGenerationTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this NameGenerationType.
func (a NameGenerationType) Key() string {
	return nameGenerationTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a NameGenerationType) String() string {
	return nameGenerationTypeValues[a.EnsureValid()].String
}
