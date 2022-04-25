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

// AttributeChoice holds a single attribute choice.
type AttributeChoice struct {
	Key   string
	Title string
}

// AttributeChoices collects the available choices for attributes for the given entity, or nil.
func AttributeChoices(entity *Entity) []*AttributeChoice {
	list := AttributeDefsFor(entity).List()
	choices := make([]*AttributeChoice, len(list)+1)
	choices[0] = &AttributeChoice{
		Key:   "10",
		Title: "10",
	}
	for i, def := range list {
		choices[i+1] = &AttributeChoice{
			Key:   def.DefID,
			Title: def.Name,
		}
	}
	return choices
}

func (c *AttributeChoice) String() string {
	return c.Title
}
