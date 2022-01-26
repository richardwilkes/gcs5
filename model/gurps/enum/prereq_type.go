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

package enum

import "strings"

// Possible PrereqType values.
const (
	AdvantagePrereq PrereqType = iota
	AttributePrereq
	ContainedQuantityPrereq
	ContainedWeightPrereq
	PrereqList
	SkillPrereq
	SpellPrereq
)

// PrereqType holds the type of a Feature.
type PrereqType uint8

// PrereqTypeFromString extracts a PrereqType from a string.
func PrereqTypeFromString(str string) PrereqType {
	for one := AdvantagePrereq; one <= SpellPrereq; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return PrereqList
}

// Key returns the key used to represent this PrereqType.
func (b PrereqType) Key() string {
	switch b {
	case AdvantagePrereq:
		return "advantage_prereq"
	case AttributePrereq:
		return "attribute_prereq"
	case ContainedQuantityPrereq:
		return "contained_quantity_prereq"
	case ContainedWeightPrereq:
		return "contained_weight_prereq"
	case SkillPrereq:
		return "skill_prereq"
	case SpellPrereq:
		return "spell_prereq"
	default: // PrereqList
		return "prereq_list"
	}
}
