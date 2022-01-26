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

package settings

import (
	"sort"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
)

// PageRefs holds a set of page references.
type PageRefs struct {
	refs map[string]*PageRef
}

// NewPageRefs creates a new, empty, PageRefs object.
func NewPageRefs() *PageRefs {
	return &PageRefs{refs: make(map[string]*PageRef)}
}

// NewPageRefsFromJSON creates a new PageRefs from a JSON object.
func NewPageRefsFromJSON(data map[string]interface{}) *PageRefs {
	p := NewPageRefs()
	for k, v := range data {
		p.refs[k] = NewPageRefFromJSON(k, encoding.Object(v))
	}
	return p
}

// Empty implements encoding.Empty.
func (p *PageRefs) Empty() bool {
	return len(p.refs) == 0
}

// ToJSON emits this object as JSON.
func (p *PageRefs) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	keys := make([]string, 0, len(p.refs))
	for k := range p.refs {
		keys = append(keys, k)
	}
	txt.SortStringsNaturalAscending(keys)
	for _, k := range keys {
		encoder.Key(k)
		p.refs[k].ToJSON(encoder)
	}
	encoder.EndObject()
}

// Count returns the number of page references being tracked.
func (p *PageRefs) Count() int {
	return len(p.refs)
}

// Add a PageRef.
func (p *PageRefs) Add(id string, ref *PageRef) {
	p.refs[id] = ref
}

// Remove a PageRef.
func (p *PageRefs) Remove(id string) {
	delete(p.refs, id)
}

// Lookup the PageRef for the given ID. If not found or if the path it points to isn't a readable file, returns nil.
func (p *PageRefs) Lookup(id string) *PageRef {
	if ref, ok := p.refs[id]; ok && fs.FileIsReadable(ref.Path) {
		return ref
	}
	return nil
}

// List returns the current PageRef list.
func (p *PageRefs) List() []*PageRef {
	list := make([]*PageRef, 0, len(p.refs))
	for _, ref := range p.refs {
		list = append(list, ref)
	}
	sort.Slice(list, func(i, j int) bool { return txt.NaturalLess(list[i].ID, list[j].ID, true) })
	return list
}
