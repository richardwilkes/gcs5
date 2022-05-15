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
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/toolbox/i18n"
)

// AttributeChoice holds a single attribute choice.
type AttributeChoice struct {
	Key   string
	Title string
}

// AttributeChoices collects the available choices for attributes for the given entity, or nil.
func AttributeChoices(entity *Entity, prefix string, includeSkills bool) []*AttributeChoice {
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix = prefix + " "
	}
	list := AttributeDefsFor(entity).List()
	extra := 1
	if includeSkills {
		extra += 3
	}
	choices := make([]*AttributeChoice, len(list)+extra)
	i := 0
	choices[i] = &AttributeChoice{
		Key:   "10",
		Title: prefix + "10",
	}
	i++
	if includeSkills {
		choices[i] = &AttributeChoice{
			Key:   gid.Skill,
			Title: prefix + i18n.Text("Skill"),
		}
		i++
		choices[i] = &AttributeChoice{
			Key:   gid.Parry,
			Title: prefix + i18n.Text("Parry"),
		}
		i++
		choices[i] = &AttributeChoice{
			Key:   gid.Block,
			Title: prefix + i18n.Text("Block"),
		}
		i++
	}
	for j, def := range list {
		choices[j+i] = &AttributeChoice{
			Key:   def.DefID,
			Title: prefix + def.Name,
		}
	}
	return choices
}

func (c *AttributeChoice) String() string {
	return c.Title
}
