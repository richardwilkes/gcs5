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
	"sort"
	"strings"

	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Attributes holds a set of Attribute objects.
type Attributes struct {
	Set map[string]*Attribute
}

// NewAttributes creates a new Attributes.
func NewAttributes(entity *Entity) *Attributes {
	a := &Attributes{Set: make(map[string]*Attribute)}
	i := 0
	for attrID := range entity.SheetSettings.Attributes.Set {
		a.Set[attrID] = NewAttribute(entity, attrID, i)
		i++
	}
	return a
}

// MarshalJSON implements json.Marshaler.
func (a *Attributes) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	e := json.NewEncoder(&buffer)
	e.SetEscapeHTML(false)
	err := e.Encode(a.List())
	return buffer.Bytes(), err
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attributes) UnmarshalJSON(data []byte) error {
	var list []*Attribute
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	a.Set = make(map[string]*Attribute, len(list))
	for i, one := range list {
		one.Order = i
		a.Set[one.ID()] = one
	}
	return nil
}

// Clone a copy of this.
func (a *Attributes) Clone(entity *Entity) *Attributes {
	clone := &Attributes{Set: make(map[string]*Attribute)}
	for k, v := range a.Set {
		clone.Set[k] = v.Clone(entity)
	}
	return clone
}

// List returns the map of Attribute objects as an ordered list.
func (a *Attributes) List() []*Attribute {
	list := make([]*Attribute, 0, len(a.Set))
	for _, v := range a.Set {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	return list
}

// Cost returns the points spent for the specified Attribute.
func (a *Attributes) Cost(attrID string) f64d4.Int {
	if attr, ok := a.Set[attrID]; ok {
		return attr.PointCost()
	}
	return 0
}

// Current resolves the given attribute ID to its current value, or f64d4.Min.
func (a *Attributes) Current(attrID string) f64d4.Int {
	if attr, ok := a.Set[attrID]; ok {
		return attr.Current()
	}
	if v, err := f64d4.FromString(attrID); err == nil {
		return v
	}
	return f64d4.Min
}

// Maximum resolves the given attribute ID to its maximum value, or f64d4.Min.
func (a *Attributes) Maximum(attrID string) f64d4.Int {
	if attr, ok := a.Set[attrID]; ok {
		return attr.Maximum()
	}
	if v, err := f64d4.FromString(attrID); err == nil {
		return v
	}
	return f64d4.Min
}

// PoolThreshold resolves the given attribute ID and state to the value for its pool threshold, or f64d4.Min.
func (a *Attributes) PoolThreshold(attrID, state string) f64d4.Int {
	if attr, ok := a.Set[attrID]; ok {
		if def := attr.AttributeDef(); def != nil {
			for _, one := range def.Thresholds {
				if strings.EqualFold(one.State, state) {
					return one.Threshold(attr.Maximum())
				}
			}
		}
	}
	return f64d4.Min
}
