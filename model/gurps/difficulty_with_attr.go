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
	"encoding/json"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
)

// AttributeDifficulty holds an attribute ID and a difficulty.
type AttributeDifficulty struct {
	Attribute  string
	Difficulty skill.Difficulty
}

func (a *AttributeDifficulty) String() string {
	return a.Attribute + "/" + string(a.Difficulty)
}

// MarshalJSON implements json.Marshaler.
func (a *AttributeDifficulty) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttributeDifficulty) UnmarshalJSON(data []byte) error {
	parts := strings.SplitN(string(data), "/", 2)
	a.Attribute = strings.TrimSpace(parts[0])
	a.Difficulty = skill.A
	if len(parts) == 2 {
		text := strings.TrimSpace(parts[1])
		for _, one := range skill.AllDifficulties {
			if strings.EqualFold(string(one), text) {
				a.Difficulty = one
				break
			}
		}
	}
	return nil
}

// Normalize the data. Should be called after loading from disk or the user.
func (a *AttributeDifficulty) Normalize(entity *Entity) {
	a.Difficulty = a.Difficulty.EnsureValid()
	text := strings.TrimSpace(a.Attribute)
	if text == "" {
		text = DefaultAttributeIDFor(entity)
	}
	var attr *AttributeDef
	list := AttributeDefsFor(entity).List()
	for _, one := range list {
		if strings.EqualFold(one.ID(), text) {
			attr = one
			break
		}
	}
	if attr == nil {
		for _, one := range list {
			if strings.EqualFold(one.Name, text) {
				attr = one
				break
			}
		}
	}
	if attr != nil {
		text = attr.ID()
	}
	a.Attribute = id.Sanitize(text, true)
}
