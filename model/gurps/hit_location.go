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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio"
)

// HitLocationData holds the Hitlocation data that gets written to disk.
type HitLocationData struct {
	Entity      *Entity   `json:"-"`
	LocID       string    `json:"id"`
	ChoiceName  string    `json:"choice_name"`
	TableName   string    `json:"table_name"`
	Slots       int       `json:"slots"`
	HitPenalty  int       `json:"hit_penalty"`
	DRBonus     int       `json:"dr_bonus"`
	Description string    `json:"description"`
	SubTable    *BodyType `json:"sub_table,omitempty"`
}

// HitLocation holds a single hit location.
type HitLocation struct {
	HitLocationData
	RollRange   string
	owningTable *BodyType
}

// MarshalJSON implements json.Marshaler.
func (h *HitLocation) MarshalJSON() ([]byte, error) {
	type calc struct {
		RollRange string         `json:"roll_range"`
		DR        map[string]int `json:"dr,omitempty"`
	}
	data := struct {
		HitLocationData
		Calc calc `json:"calc"`
	}{
		HitLocationData: h.HitLocationData,
		Calc: calc{
			RollRange: h.RollRange,
		},
	}
	if h.Entity != nil {
		data.Calc.DR = h.DR(h.Entity, nil, nil)
		if _, exists := data.Calc.DR[feature.All]; !exists {
			data.Calc.DR[feature.All] = 0
		}
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *HitLocation) UnmarshalJSON(data []byte) error {
	h.HitLocationData = HitLocationData{}
	if err := json.Unmarshal(data, &h.HitLocationData); err != nil {
		return err
	}
	if h.SubTable != nil {
		h.SubTable.SetOwningLocation(h)
	}
	return nil
}

// ID returns the ID.
func (h *HitLocation) ID() string {
	return h.LocID
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (h *HitLocation) SetID(value string) {
	h.LocID = id.Sanitize(value, false, ReservedIDs...)
}

// DR computes the DR coverage for this HitLocation. If 'tooltip' isn't nil, the buffer will be updated with details on
// how the DR was calculated. If 'drMap' isn't nil, it will be returned.
func (h *HitLocation) DR(entity *Entity, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if h.DRBonus != 0 {
		drMap[feature.All] += h.DRBonus
		if tooltip != nil {
			fmt.Fprintf(tooltip, "\n%s [%+d against %s attacks]", h.ChoiceName, h.DRBonus, feature.All)
		}
	}
	drMap = entity.AddDRBonusesFor(feature.HitLocationPrefix+h.LocID, tooltip, drMap)
	if h.owningTable != nil && h.owningTable.owningLocation != nil {
		drMap = h.owningTable.owningLocation.DR(entity, tooltip, drMap)
	}
	if tooltip != nil && len(drMap) != 0 {
		keys := make([]string, 0, len(drMap))
		for k := range drMap {
			keys = append(keys, k)
		}
		txt.SortStringsNaturalAscending(keys)
		base := drMap[feature.All]
		var buffer bytes.Buffer
		buffer.WriteByte('\n')
		for _, k := range keys {
			value := drMap[k]
			if !strings.EqualFold(feature.All, k) {
				value += base
			}
			fmt.Fprintf(&buffer, "\n%d against %s attacks", value, k)
		}
		buffer.WriteByte('\n')
		tooltip.Insert(0, buffer.Bytes())
	}
	return drMap
}

// DisplayDR returns the DR for this location, formatted as a string.
func (h *HitLocation) DisplayDR(entity *Entity, tooltip *xio.ByteBuffer) string {
	drMap := h.DR(entity, tooltip, nil)
	all, exists := drMap[feature.All]
	if !exists {
		drMap[feature.All] = 0
	}
	keys := make([]string, 0, len(drMap))
	keys = append(keys, feature.All)
	for k := range drMap {
		if k != feature.All {
			keys = append(keys, k)
		}
	}
	txt.SortStringsNaturalAscending(keys[1:])
	var buffer strings.Builder
	for _, k := range keys {
		dr := drMap[k]
		if k != feature.All {
			dr += all
		}
		if buffer.Len() != 0 {
			buffer.WriteByte('/')
		}
		buffer.WriteString(strconv.Itoa(dr))
	}
	return buffer.String()
}

// SetSubTable sets the BodyType as a sub-table.
func (h *HitLocation) SetSubTable(bodyType *BodyType) {
	if bodyType == nil && h.SubTable != nil {
		h.SubTable.SetOwningLocation(nil)
	}
	if h.SubTable = bodyType; h.SubTable != nil {
		h.SubTable.SetOwningLocation(h)
	}
}

func (h *HitLocation) populateMap(m map[string]*HitLocation) {
	m[h.LocID] = h
	if h.SubTable != nil {
		h.SubTable.populateMap(m)
	}
}

func (h *HitLocation) updateRollRange(start int) int {
	switch h.Slots {
	case 0:
		h.RollRange = "-"
	case 1:
		h.RollRange = strconv.Itoa(start)
	default:
		h.RollRange = fmt.Sprintf("%d-%d", start, start+h.Slots-1)
	}
	if h.SubTable != nil {
		h.SubTable.updateRollRanges()
	}
	return start + h.Slots
}
