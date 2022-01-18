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
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/xio"
)

// HitLocationPrefix is the prefix used on all hit locations for DR bonuses.
const HitLocationPrefix = "hit_location."

// HitLocation holds a single hit location.
type HitLocation struct {
	ID          string          `json:"id"`
	ChoiceName  string          `json:"choice_name"`
	TableName   string          `json:"table_name"`
	Slots       int             `json:"slots,omitempty"`
	HitPenalty  int             `json:"hit_penalty,omitempty"`
	DRBonus     int             `json:"dr_bonus,omitempty"`
	Description string          `json:"description,omitempty"`
	Calc        HitLocationCalc `json:"calc"`
	SubTable    *BodyType       `json:"sub_table,omitempty"`
	owningTable *BodyType
}

// HitLocationCalc holds values GCS calculates for a HitLocation, but that we want to be present in any json output so
// that other uses of the data don't have to replicate the code to calculate it.
type HitLocationCalc struct {
	RollRange string         `json:"roll_range"`
	DR        map[string]int `json:"dr,omitempty"`
}

// DR computes the DR coverage for this HitLocation. If 'tooltip' isn't nil, the buffer will be updated with details on
// how the DR was calculated. If 'drMap' isn't nil, it will be returned.
func (h *HitLocation) DR(entity *Entity, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if h.DRBonus != 0 {
		drMap[All] += h.DRBonus
		if tooltip != nil {
			fmt.Fprintf(tooltip, "\n%s [%+d against %s attacks]", h.ChoiceName, h.DRBonus, All)
		}
	}
	drMap = entity.AddDRBonusesFor(HitLocationPrefix+h.ID, tooltip, drMap)
	if h.owningTable != nil && h.owningTable.owningLocation != nil {
		drMap = h.owningTable.owningLocation.DR(entity, tooltip, drMap)
	}
	if tooltip != nil && len(drMap) != 0 {
		keys := make([]string, 0, len(drMap))
		for k := range drMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		base := drMap[All]
		var buffer bytes.Buffer
		buffer.WriteByte('\n')
		for _, k := range keys {
			value := drMap[k]
			if !strings.EqualFold(All, k) {
				value += base
			}
			fmt.Fprintf(&buffer, "\n%d against %s attacks", value, k)
		}
		buffer.WriteByte('\n')
		tooltip.Insert(0, buffer.Bytes())
	}
	return drMap
}

func (h *HitLocation) updateRollRange(start int) int {
	switch h.Slots {
	case 0:
		h.Calc.RollRange = "-"
	case 1:
		h.Calc.RollRange = strconv.Itoa(start)
	default:
		h.Calc.RollRange = fmt.Sprintf("%d-%d", start, start+h.Slots-1)
	}
	if h.SubTable != nil {
		h.SubTable.updateRollRanges()
	}
	return start + h.Slots
}

func (h *HitLocation) updateDR(entity *Entity) {
	h.Calc.DR = nil
	if entity != nil {
		h.Calc.DR = h.DR(entity, nil, nil)
		if _, exists := h.Calc.DR[All]; !exists {
			h.Calc.DR[All] = 0
		}
	}
	if h.SubTable != nil {
		h.SubTable.updateDR(entity)
	}
}
