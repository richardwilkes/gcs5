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
	registerExportableGCSFileInfo(".gcs", res.GCSSheet, sheet.NewSheetFromFile)
	registerGCSFileInfo(".gct", []string{".gct"}, res.GCSTemplate, NewTemplateFromFile)
	groupWith := []string{".adq", ".adm", ".eqp", ".eqm", ".skl", ".spl", ".not"}
	registerGCSFileInfo(".adq", groupWith, res.GCSAdvantages, NewAdvantageTableDockableFromFile)
	registerGCSFileInfo(".adm", groupWith, res.GCSAdvantageModifiers, NewAdvantageModifierTableDockableFromFile)
	registerGCSFileInfo(".eqp", groupWith, res.GCSEquipment, NewEquipmentTableDockableFromFile)
	registerGCSFileInfo(".eqm", groupWith, res.GCSEquipmentModifiers, NewEquipmentModifierTableDockableFromFile)
	registerGCSFileInfo(".skl", groupWith, res.GCSSkills, NewSkillTableDockableFromFile)
	registerGCSFileInfo(".spl", groupWith, res.GCSSpells, NewSpellTableDockableFromFile)
	registerGCSFileInfo(".not", groupWith, res.GCSNotes, NewNoteTableDockableFromFile)
}

func registerGCSFileInfo(ext string, groupWith []string, svg *unison.SVG, loader func(filePath string) (unison.Dockable, error)) {
	library.FileInfo{
		Extension:             ext,
		ExtensionsToGroupWith: groupWith,
		SVG:                   svg,
		Load:                  loader,
		IsGCSData:             true,
	}.Register()
}

func registerExportableGCSFileInfo(ext string, svg *unison.SVG, loader func(filePath string) (unison.Dockable, error)) {
	library.FileInfo{
		Extension:             ext,
		ExtensionsToGroupWith: []string{ext},
		SVG:                   svg,
		Load:                  loader,
		IsGCSData:             true,
		IsExportable:          true,
	}.Register()
}
