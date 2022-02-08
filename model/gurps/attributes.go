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

	"github.com/richardwilkes/json"
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
