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
	"context"
	"io/fs"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/toolbox/errs"
)

const templateTypeKey = "template"

var _ ListProvider = &Template{}

// Template holds the GURPS Template data that is written to disk.
type Template struct {
	Type       string       `json:"type"`
	Version    int          `json:"version"`
	ID         uuid.UUID    `json:"id"`
	Advantages []*Advantage `json:"advantages,omitempty"`
	Skills     []*Skill     `json:"skills,omitempty"`
	Spells     []*Spell     `json:"spells,omitempty"`
	Equipment  []*Equipment `json:"equipment,omitempty"`
	Notes      []*Note      `json:"notes,omitempty"`
}

// NewTemplateFromFile loads a Template from a file.
func NewTemplateFromFile(fileSystem fs.FS, filePath string) (*Template, error) {
	var template Template
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &template); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if template.Type != templateTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(template.Version); err != nil {
		return nil, err
	}
	return &template, nil
}

// NewTemplate creates a new Template.
func NewTemplate() *Template {
	template := &Template{
		Type: templateTypeKey,
		ID:   id.NewUUID(),
	}
	return template
}

// Save the Template to a file as JSON.
func (t *Template) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, t)
}

// AdvantageList implements ListProvider
func (t *Template) AdvantageList() []*Advantage {
	return t.Advantages
}

// CarriedEquipmentList implements ListProvider
func (t *Template) CarriedEquipmentList() []*Equipment {
	return t.Equipment
}

// OtherEquipmentList implements ListProvider
func (t *Template) OtherEquipmentList() []*Equipment {
	return nil
}

// SkillList implements ListProvider
func (t *Template) SkillList() []*Skill {
	return t.Skills
}

// SpellList implements ListProvider
func (t *Template) SpellList() []*Spell {
	return t.Spells
}

// NoteList implements ListProvider
func (t *Template) NoteList() []*Note {
	return t.Notes
}
