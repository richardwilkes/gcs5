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
	"github.com/richardwilkes/gcs/unit/length"
	"github.com/richardwilkes/gcs/unit/weight"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// Standard height and width for the portrait
const (
	PortraitHeight = 96
	PortraitWidth  = 3 * PortraitHeight / 4
)

// BaseProfile holds the base profile information.
type BaseProfile struct {
	Name         string `json:"name,omitempty"`
	TechLevel    string `json:"tech_level,omitempty"`
	SizeModifier int    `json:"SM,omitempty"`
	PortraitData []byte `json:"portrait,omitempty"`
	portrait     *unison.Image
}

// NPCProfile holds the profile information for an NPC.
type NPCProfile struct {
	BaseProfile  `json:",inline"`
	Title        string        `json:"title,omitempty"`
	Organization string        `json:"organization,omitempty"`
	Religion     string        `json:"religion,omitempty"`
	Age          string        `json:"age,omitempty"`
	Eyes         string        `json:"eyes,omitempty"`
	Hair         string        `json:"hair,omitempty"`
	Skin         string        `json:"skin,omitempty"`
	Handedness   string        `json:"handedness,omitempty"`
	Gender       string        `json:"gender,omitempty"`
	Height       length.Length `json:"height,omitempty"`
	Weight       weight.Weight `json:"weight,omitempty"`
}

// PCProfile holds the profile information for a PC.
type PCProfile struct {
	NPCProfile `json:",inline"`
	PlayerName string `json:"player_name,omitempty"`
	Birthday   string `json:"birthday,omitempty"`
}

// Portrait returns the portrait image, if there is one.
func (p *BaseProfile) Portrait() *unison.Image {
	if p.portrait == nil && len(p.PortraitData) > 0 {
		var err error
		p.portrait, err = unison.NewImageFromBytes(p.PortraitData, 0.5)
		if err != nil {
			jot.Error(errs.NewWithCause("unable to load portrait data", err))
			p.portrait = nil
			p.PortraitData = nil
		}
	}
	return p.portrait
}
