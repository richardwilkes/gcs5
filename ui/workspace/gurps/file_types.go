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
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/sheet"
	"github.com/richardwilkes/unison"
)

// RegisterFileTypes registers GCS file types.
func RegisterFileTypes() {
	registerExportableGCSFileInfo(".gcs", res.GCSSheet, sheet.NewSheetDockable)
	registerGCSFileInfo(".gct", res.GCSTemplate, NewTemplateDockable)
	registerGCSFileInfo(".adq", res.GCSAdvantages, NewAdvantageListDockable)
	registerGCSFileInfo(".adm", res.GCSAdvantageModifiers, NewAdvantageModifierListDockable)
	registerGCSFileInfo(".eqp", res.GCSEquipment, NewEquipmentListDockable)
	registerGCSFileInfo(".eqm", res.GCSEquipmentModifiers, NewEquipmentModifierListDockable)
	registerGCSFileInfo(".skl", res.GCSSkills, NewSkillListDockable)
	registerGCSFileInfo(".spl", res.GCSSpells, NewSpellListDockable)
	registerGCSFileInfo(".not", res.GCSNotes, NewNoteListDockable)
}

func registerGCSFileInfo(ext string, svg *unison.SVG, loader func(filePath string) (unison.Dockable, error)) {
	library.FileInfo{
		Extension: ext,
		SVG:       svg,
		Load:      loader,
		IsGCSData: true,
	}.Register()
}

func registerExportableGCSFileInfo(ext string, svg *unison.SVG, loader func(filePath string) (unison.Dockable, error)) {
	library.FileInfo{
		Extension:    ext,
		SVG:          svg,
		Load:         loader,
		IsGCSData:    true,
		IsExportable: true,
	}.Register()
}
