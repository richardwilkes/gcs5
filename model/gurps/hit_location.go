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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// HitLocationPrefix is the prefix used on all hit locations for DR bonuses.
const HitLocationPrefix = "hit_location."

const (
	hitLocationIDKey            = "id"
	hitLocationChoiceNameKey    = "choice_name"
	hitLocationTableNameKey     = "table_name"
	hitLocationSlotsKey         = "slots"
	hitLocationHitPenaltyKey    = "hit_penalty"
	hitLocationDRBonusKey       = "dr_bonus"
	hitLocationDescriptionKey   = "description"
	hitLocationSubTableKey      = "sub_table"
	hitLocationCalcRollRangeKey = "roll_range"
	hitLocationCalcDRKey        = "dr"
)

// HitLocation holds a single hit location.
type HitLocation struct {
	id          string
	ChoiceName  string
	TableName   string
	RollRange   string
	Slots       int
	HitPenalty  int
	DRBonus     int
	Description string
	SubTable    *BodyType
	owningTable *BodyType
}

// NewHitLocationFromJSON creates a new HitLocation from a JSON object.
func NewHitLocationFromJSON(data map[string]interface{}) *HitLocation {
	h := &HitLocation{
		ChoiceName:  encoding.String(data[hitLocationChoiceNameKey]),
		TableName:   encoding.String(data[hitLocationTableNameKey]),
		Slots:       int(encoding.Number(data[hitLocationSlotsKey]).AsInt64()),
		HitPenalty:  int(encoding.Number(data[hitLocationHitPenaltyKey]).AsInt64()),
		DRBonus:     int(encoding.Number(data[hitLocationDRBonusKey]).AsInt64()),
		Description: encoding.String(data[hitLocationDescriptionKey]),
	}
	h.SetID(encoding.String(data[hitLocationIDKey]))
	if obj := encoding.Object(data[hitLocationSubTableKey]); obj != nil {
		h.SetSubTable(NewBodyTypeFromJSON(obj))
	}
	return h
}

// ToJSON emits this object as JSON.
func (h *HitLocation) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	encoder.KeyedString(hitLocationIDKey, h.id, false, false)
	encoder.KeyedString(hitLocationChoiceNameKey, h.ChoiceName, true, true)
	encoder.KeyedString(hitLocationTableNameKey, h.TableName, true, true)
	encoder.KeyedNumber(hitLocationSlotsKey, fixed.F64d4FromInt64(int64(h.Slots)), true)
	encoder.KeyedNumber(hitLocationHitPenaltyKey, fixed.F64d4FromInt64(int64(h.HitPenalty)), true)
	encoder.KeyedNumber(hitLocationDRBonusKey, fixed.F64d4FromInt64(int64(h.DRBonus)), true)
	encoder.KeyedString(hitLocationDescriptionKey, h.Description, true, true)
	if h.SubTable != nil {
		encoder.Key(hitLocationSubTableKey)
		h.SubTable.ToJSON(encoder, entity)
	}

	// Emit the calculated values for third parties
	encoder.Key(calcKey)
	encoder.StartObject()
	encoder.KeyedString(hitLocationCalcRollRangeKey, h.RollRange, false, false)
	if entity != nil {
		drMap := h.DR(entity, nil, nil)
		if _, exists := drMap[All]; !exists {
			drMap[All] = 0
		}
		keys := make([]string, 0, len(drMap))
		for k := range drMap {
			keys = append(keys, k)
		}
		txt.SortStringsNaturalAscending(keys)
		encoder.Key(hitLocationCalcDRKey)
		encoder.StartObject()
		for _, k := range keys {
			encoder.KeyedNumber(k, fixed.F64d4FromInt64(int64(drMap[k])), false)
		}
		encoder.EndObject()
	}
	encoder.EndObject()

	encoder.EndObject()
}

// ID returns the ID.
func (h *HitLocation) ID() string {
	return h.id
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (h *HitLocation) SetID(value string) {
	h.id = id.Sanitize(value, false, ReservedIDs...)
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
	drMap = entity.AddDRBonusesFor(HitLocationPrefix+h.id, tooltip, drMap)
	if h.owningTable != nil && h.owningTable.owningLocation != nil {
		drMap = h.owningTable.owningLocation.DR(entity, tooltip, drMap)
	}
	if tooltip != nil && len(drMap) != 0 {
		keys := make([]string, 0, len(drMap))
		for k := range drMap {
			keys = append(keys, k)
		}
		txt.SortStringsNaturalAscending(keys)
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

// DisplayDR returns the DR for this location, formatted as a string.
func (h *HitLocation) DisplayDR(entity *Entity, tooltip *xio.ByteBuffer) string {
	drMap := h.DR(entity, tooltip, nil)
	all, exists := drMap[All]
	if !exists {
		drMap[All] = 0
	}
	keys := make([]string, 0, len(drMap))
	keys = append(keys, All)
	for k := range drMap {
		if k != All {
			keys = append(keys, k)
		}
	}
	txt.SortStringsNaturalAscending(keys[1:])
	var buffer strings.Builder
	for _, k := range keys {
		dr := drMap[k]
		if k != All {
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
	m[h.id] = h
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
