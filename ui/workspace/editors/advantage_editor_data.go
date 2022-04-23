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

package editors

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

var _ editorData[*gurps.Advantage] = &advantageEditorData{}

type advantageEditorData struct {
	name           string
	notes          string
	vttNotes       string
	userDesc       string
	tags           string
	pageRef        string
	ancestry       string
	basePoints     f64d4.Int
	pointsPerLevel f64d4.Int
	levels         f64d4.Int
	cr             advantage.SelfControlRoll
	crAdj          gurps.SelfControlRollAdj
	containerType  advantage.ContainerType
	roundCostDown  bool
	disabled       bool
	mental         bool
	physical       bool
	social         bool
	exotic         bool
	supernatural   bool
}

func (d *advantageEditorData) From(advantage *gurps.Advantage) {
	d.name = advantage.Name
	d.notes = advantage.LocalNotes
	d.vttNotes = advantage.VTTNotes
	d.userDesc = advantage.UserDesc
	d.tags = gurps.CombineTags(advantage.Categories)
	d.pageRef = advantage.PageRef
	d.cr = advantage.CR
	d.crAdj = advantage.CRAdj
	d.disabled = advantage.Disabled
	if advantage.Container() {
		d.containerType = advantage.ContainerType
		d.ancestry = advantage.Ancestry
	} else {
		d.mental = advantage.Mental
		d.physical = advantage.Physical
		d.social = advantage.Social
		d.exotic = advantage.Exotic
		d.supernatural = advantage.Supernatural
		d.roundCostDown = advantage.RoundCostDown
		d.basePoints = advantage.BasePoints
		d.pointsPerLevel = advantage.PointsPerLevel
		d.levels = advantage.Levels
	}
}

func (d *advantageEditorData) Apply(advantage *gurps.Advantage) {
	advantage.Name = d.name
	advantage.LocalNotes = d.notes
	advantage.VTTNotes = d.vttNotes
	advantage.UserDesc = d.userDesc
	advantage.Categories = gurps.ExtractTags(d.tags)
	advantage.PageRef = d.pageRef
	advantage.CR = d.cr
	advantage.CRAdj = d.crAdj
	advantage.Disabled = d.disabled
	if advantage.Container() {
		advantage.ContainerType = d.containerType
		advantage.Ancestry = d.ancestry
	} else {
		advantage.Mental = d.mental
		advantage.Physical = d.physical
		advantage.Social = d.social
		advantage.Exotic = d.exotic
		advantage.Supernatural = d.supernatural
		advantage.RoundCostDown = d.roundCostDown
		advantage.BasePoints = d.basePoints
		advantage.PointsPerLevel = d.pointsPerLevel
		if d.pointsPerLevel != 0 {
			advantage.Levels = d.levels
		} else {
			advantage.Levels = 0
		}
	}
}
