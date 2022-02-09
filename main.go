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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/richardwilkes/gcs/internal/ui"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/jotrotate"
)

func main() {
	cmdline.AppName = "GCS"
	cmdline.AppCmdName = "gcs"
	cmdline.License = "Mozilla Public License, version 2.0"
	cmdline.CopyrightYears = fmt.Sprintf("1998-%d", time.Now().Year())
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.gcs"
	if cmdline.AppVersion == "" {
		cmdline.AppVersion = "0.0"
	}
	cl := cmdline.New(true)
	fileList := jotrotate.ParseAndSetup(cl)

	settings.Global() // Here to force early initialization
	processDir("../gcs_master_library/Library/Home Brew/Characters")
	atexit.Exit(0)

	ui.Start(fileList) // Never returns
}

func processDir(dir string) {
	const convertedDir = "converted/"
	entries, err := os.ReadDir(dir)
	jot.FatalIfErr(err)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		switch filepath.Ext(name) {
		case ".adq":
			var adq []*gurps.Advantage
			adq, err = gurps.NewAdvantagesFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveAdvantages(adq, convertedDir+name))
		case ".adm":
			var adm []*gurps.AdvantageModifier
			adm, err = gurps.NewAdvantageModifiersFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveAdvantageModifiers(adm, convertedDir+name))
		case ".eqm":
			var eqm []*gurps.EquipmentModifier
			eqm, err = gurps.NewEquipmentModifiersFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveEquipmentModifiers(eqm, convertedDir+name))
		case ".eqp":
			var eqp []*gurps.Equipment
			eqp, err = gurps.NewEquipmentFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveEquipment(eqp, convertedDir+name))
		case ".not":
			var not []*gurps.Note
			not, err = gurps.NewNotesFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveNotes(not, convertedDir+name))
		case ".skl":
			var skl []*gurps.Skill
			skl, err = gurps.NewSkillsFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveSkills(skl, convertedDir+name))
		case ".spl":
			var spl []*gurps.Spell
			spl, err = gurps.NewSpellsFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gurps.SaveSpells(spl, convertedDir+name))
		case ".ghl":
			var ghl *gurps.BodyType
			ghl, err = gurps.NewBodyTypeFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(ghl.Save(convertedDir + name))
		case ".gas":
			var gas *gurps.AttributeDefs
			gas, err = gurps.NewAttributeDefsFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(gas.Save(convertedDir + name))
		case ".ancestry":
			var anc *ancestry.Ancestry
			anc, err = ancestry.NewAncestoryFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(anc.Save(convertedDir + name))
		case ".gcs":
			var entity *gurps.Entity
			entity, err = gurps.NewEntityFromFile(os.DirFS(dir), name)
			jot.FatalIfErr(err)
			jot.FatalIfErr(entity.Save(convertedDir + name))
		// case ".gct":
		default:
			fmt.Println("skipping " + name)
			continue
		}
		var m map[string]interface{}
		jot.FatalIfErr(jio.LoadFromFile(context.Background(), dir+"/"+name, &m))
		jot.FatalIfErr(jio.SaveToFile(context.Background(), convertedDir+"orig-sorted-"+name, m))
		m = make(map[string]interface{})
		jot.FatalIfErr(jio.LoadFromFile(context.Background(), convertedDir+name, &m))
		jot.FatalIfErr(jio.SaveToFile(context.Background(), convertedDir+"/sorted-"+name, m))
	}
}
