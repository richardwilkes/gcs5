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

package export

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/images"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	descriptionKey        = "DESCRIPTION"
	descriptionPrimaryKey = "DESCRIPTION_PRIMARY"
	parentIDKey           = "PARENT_ID"
	pointsKey             = "POINTS"
	refKey                = "REF"
	satisfiedKey          = "SATISFIED"
	styleIndentWarningKey = "STYLE_INDENT_WARNING"
	typeKey               = "TYPE"
)

type legacyExporter struct {
	entity             *gurps.Entity
	template           []byte
	pos                int
	exportPath         string
	onlyCategories     map[string]bool
	excludedCategories map[string]bool
	out                *bufio.Writer
	encodeText         bool
	enhancedKeyParsing bool
}

// LegacyExport performs the text template export function that matches the old Java code base.
func LegacyExport(entity *gurps.Entity, templatePath, exportPath string) (err error) {
	ex := &legacyExporter{
		entity:             entity,
		exportPath:         exportPath,
		onlyCategories:     make(map[string]bool),
		excludedCategories: make(map[string]bool),
		encodeText:         true,
	}
	if ex.template, err = os.ReadFile(templatePath); err != nil {
		return errs.Wrap(err)
	}
	var out *os.File
	if out, err = os.Create(exportPath); err != nil {
		return errs.Wrap(err)
	}
	ex.out = bufio.NewWriter(out)
	defer func() { //nolint:gosec // Yes, this is safe
		if flushErr := ex.out.Flush(); flushErr != nil && err == nil {
			err = errs.Wrap(flushErr)
		}
		if closeErr := out.Close(); closeErr != nil && err == nil {
			err = errs.Wrap(closeErr)
		}
	}()
	lookForKeyMarker := true
	var keyBuffer bytes.Buffer
	for ex.pos < len(ex.template) {
		ch := ex.template[ex.pos]
		ex.pos++
		switch {
		case lookForKeyMarker:
			if ch == '@' {
				lookForKeyMarker = false
			} else {
				ex.out.WriteByte(ch)
			}
		case ch == '_' || (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
			keyBuffer.WriteByte(ch)
		default:
			if !ex.enhancedKeyParsing || ch != '@' {
				ex.pos--
			}
			if err = ex.emitKey(keyBuffer.String()); err != nil {
				return err
			}
			keyBuffer.Reset()
			lookForKeyMarker = true
		}
	}
	if keyBuffer.Len() != 0 {
		if err = ex.emitKey(keyBuffer.String()); err != nil {
			return err
		}
	}
	return nil
}

func (ex *legacyExporter) emitKey(key string) error {
	switch key {
	case "GRID_TEMPLATE":
		ex.out.WriteString(ex.entity.SheetSettings.BlockLayout.HTMLGridTemplate())
	case "ENCODING_OFF":
		ex.encodeText = false
	case "ENHANCED_KEY_PARSING":
		ex.enhancedKeyParsing = true
	case "PORTRAIT":
		portraitData := ex.entity.Profile.PortraitData
		if len(portraitData) == 0 {
			portraitData = images.DefaultPortraitData
		}
		filePath := filepath.Join(filepath.Dir(ex.exportPath), xfs.TrimExtension(filepath.Base(ex.exportPath))+".png")
		if err := os.WriteFile(filePath, portraitData, 0o640); err != nil {
			return errs.Wrap(err)
		}
		ex.out.WriteString(url.PathEscape(filePath))
	case "PORTRAIT_EMBEDDED":
		portraitData := ex.entity.Profile.PortraitData
		if len(portraitData) == 0 {
			portraitData = images.DefaultPortraitData
		}
		ex.out.WriteString("data:image/png;base64,")
		ex.out.WriteString(base64.URLEncoding.EncodeToString(portraitData))
	case "NAME":
		ex.writeEncodedText(ex.entity.Profile.Name)
	case "TITLE":
		ex.writeEncodedText(ex.entity.Profile.Title)
	case "ORGANIZATION":
		ex.writeEncodedText(ex.entity.Profile.Organization)
	case "RELIGION":
		ex.writeEncodedText(ex.entity.Profile.Religion)
	case "PLAYER":
		ex.writeEncodedText(ex.entity.Profile.PlayerName)
	case "CREATED_ON":
		ex.writeEncodedText(ex.entity.CreatedOn.String())
	case "MODIFIED_ON":
		ex.writeEncodedText(ex.entity.ModifiedOn.String())
	case "TOTAL_POINTS":
		if settings.Global().General.IncludeUnspentPointsInTotal {
			ex.writeEncodedText(ex.entity.TotalPoints.String())
		} else {
			ex.writeEncodedText(ex.entity.SpentPoints().String())
		}
	case "ATTRIBUTE_POINTS":
		ex.writeEncodedText(ex.entity.AttributePoints().String())
	case "ST_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Strength).String())
	case "DX_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Dexterity).String())
	case "IQ_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Intelligence).String())
	case "HT_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Health).String())
	case "PERCEPTION_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Perception).String())
	case "WILL_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.Will).String())
	case "FP_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.FatiguePoints).String())
	case "HP_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.HitPoints).String())
	case "BASIC_SPEED_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.BasicSpeed).String())
	case "BASIC_MOVE_POINTS":
		ex.writeEncodedText(ex.entity.Attributes.Cost(gid.BasicMove).String())
	case "ADVANTAGE_POINTS":
		pts, _, _, _ := ex.entity.AdvantagePoints()
		ex.writeEncodedText(pts.String())
	case "DISADVANTAGE_POINTS":
		_, pts, _, _ := ex.entity.AdvantagePoints()
		ex.writeEncodedText(pts.String())
	case "QUIRK_POINTS":
		_, _, _, pts := ex.entity.AdvantagePoints()
		ex.writeEncodedText(pts.String())
	case "RACE_POINTS":
		_, _, pts, _ := ex.entity.AdvantagePoints()
		ex.writeEncodedText(pts.String())
	case "SKILL_POINTS":
		ex.writeEncodedText(ex.entity.SkillPoints().String())
	case "SPELL_POINTS":
		ex.writeEncodedText(ex.entity.SpellPoints().String())
	case "UNSPENT_POINTS", "EARNED_POINTS":
		ex.writeEncodedText(ex.entity.UnspentPoints().String())
	case "HEIGHT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultLengthUnits.Format(ex.entity.Profile.Height))
	case "WEIGHT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.Profile.Weight))
	case "GENDER":
		ex.writeEncodedText(ex.entity.Profile.Gender)
	case "HAIR":
		ex.writeEncodedText(ex.entity.Profile.Hair)
	case "EYES":
		ex.writeEncodedText(ex.entity.Profile.Eyes)
	case "AGE":
		ex.writeEncodedText(ex.entity.Profile.Age)
	case "SIZE":
		ex.writeEncodedText(ex.entity.Profile.AdjustedSizeModifier().StringWithSign())
	case "SKIN":
		ex.writeEncodedText(ex.entity.Profile.Skin)
	case "BIRTHDAY":
		ex.writeEncodedText(ex.entity.Profile.Birthday)
	case "TL":
		ex.writeEncodedText(ex.entity.Profile.TechLevel)
	case "HAND":
		ex.writeEncodedText(ex.entity.Profile.Handedness)
	case "ST":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Strength).String())
	case "DX":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Dexterity).String())
	case "IQ":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Intelligence).String())
	case "HT":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Health).String())
	case "FP":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.FatiguePoints).String())
	case "BASIC_FP":
		ex.writeEncodedText(ex.entity.Attributes.Maximum(gid.FatiguePoints).String())
	case "HP":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.HitPoints).String())
	case "BASIC_HP":
		ex.writeEncodedText(ex.entity.Attributes.Maximum(gid.HitPoints).String())
	case "WILL":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Will).String())
	case "FRIGHT_CHECK":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.FrightCheck).String())
	case "BASIC_SPEED":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.BasicSpeed).String())
	case "BASIC_MOVE":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.BasicMove).String())
	case "PERCEPTION":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Perception).String())
	case "VISION":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Vision).String())
	case "HEARING":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Hearing).String())
	case "TASTE_SMELL":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.TasteSmell).String())
	case "TOUCH":
		ex.writeEncodedText(ex.entity.Attributes.Current(gid.Touch).String())
	case "THRUST":
		ex.writeEncodedText(ex.entity.Thrust().String())
	case "SWING":
		ex.writeEncodedText(ex.entity.Swing().String())
	case "GENERAL_DR":
		dr := 0
		if torso := ex.entity.SheetSettings.HitLocations.LookupLocationByID(ex.entity, gid.Torso); torso != nil {
			dr = torso.DR(ex.entity, nil, nil)[gid.All]
		}
		ex.writeEncodedText(strconv.Itoa(dr))
	case "CURRENT_DODGE":
		ex.writeEncodedText(strconv.Itoa(ex.entity.Dodge(ex.entity.EncumbranceLevel(false))))
	case "CURRENT_MOVE":
		ex.writeEncodedText(strconv.Itoa(ex.entity.Move(ex.entity.EncumbranceLevel(false))))
	case "BEST_CURRENT_PARRY":
		ex.writeEncodedText(ex.bestWeaponDefense(func(w *gurps.Weapon) string { return w.ResolvedParry(nil) }))
	case "BEST_CURRENT_BLOCK":
		ex.writeEncodedText(ex.bestWeaponDefense(func(w *gurps.Weapon) string { return w.ResolvedBlock(nil) }))
	case "TIRED":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.FatiguePoints, "tired").String())
	case "FP_COLLAPSE":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.FatiguePoints, "collapse").String())
	case "UNCONSCIOUS":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.FatiguePoints, "unconscious").String())
	case "REELING":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "reeling").String())
	case "HP_COLLAPSE":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "collapse").String())
	case "DEATH_CHECK_1":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "dying #1").String())
	case "DEATH_CHECK_2":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "dying #2").String())
	case "DEATH_CHECK_3":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "dying #3").String())
	case "DEATH_CHECK_4":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "dying #4").String())
	case "DEAD":
		ex.writeEncodedText(ex.entity.Attributes.PoolThreshold(gid.HitPoints, "dead").String())
	case "BASIC_LIFT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.BasicLift()))
	case "ONE_HANDED_LIFT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.OneHandedLift()))
	case "TWO_HANDED_LIFT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.TwoHandedLift()))
	case "SHOVE":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.ShoveAndKnockOver()))
	case "RUNNING_SHOVE":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.RunningShoveAndKnockOver()))
	case "CARRY_ON_BACK":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.CarryOnBack()))
	case "SHIFT_SLIGHTLY":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.ShiftSlightly()))
	case "CARRIED_WEIGHT":
		ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.WeightCarried(false)))
	case "CARRIED_VALUE":
		ex.writeEncodedText("$" + ex.entity.WealthCarried().String())
	case "OTHER_EQUIPMENT_VALUE":
		ex.writeEncodedText("$" + ex.entity.WealthNotCarried().String())
	case "NOTES":
		needBlanks := false
		gurps.TraverseNotes(func(n *gurps.Note) bool {
			if needBlanks {
				ex.out.WriteString("\n\n")
			} else {
				needBlanks = true
			}
			ex.writeEncodedText(n.Text)
			return false
		}, ex.entity.Notes...)
	case "RACE":
		ex.writeEncodedText(ex.entity.Ancestry().Name)
	case "BODY_TYPE":
		ex.writeEncodedText(ex.entity.SheetSettings.HitLocations.Name)
	case "ENCUMBRANCE_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(datafile.AllEncumbrance)))
	case "ENCUMBRANCE_LOOP_START":
		ex.processEncumbranceLoop(ex.extractUpToMarker("ENCUMBRANCE_LOOP_END"))
	case "HIT_LOCATION_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.SheetSettings.HitLocations.Locations)))
	case "HIT_LOCATION_LOOP_START":
		ex.processHitLocationLoop(ex.extractUpToMarker("HIT_LOCATION_LOOP_END"))
	case "ADVANTAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeByAdvantageCategories)
	case "ADVANTAGES_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("ADVANTAGES_LOOP_END"), ex.includeByAdvantageCategories)
	case "ADVANTAGES_ALL_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeAdvantagesAndPerks)
	case "ADVANTAGES_ALL_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("ADVANTAGES_ALL_LOOP_END"), ex.includeAdvantagesAndPerks)
	case "ADVANTAGES_ONLY_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeAdvantages)
	case "ADVANTAGES_ONLY_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("ADVANTAGES_ONLY_LOOP_END"), ex.includeAdvantages)
	case "DISADVANTAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeDisadvantages)
	case "DISADVANTAGES_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("DISADVANTAGES_LOOP_END"), ex.includeDisadvantages)
	case "DISADVANTAGES_ALL_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeDisadvantagesAndQuirks)
	case "DISADVANTAGES_ALL_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("DISADVANTAGES_ALL_LOOP_END"), ex.includeDisadvantagesAndQuirks)
	case "QUIRKS_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeQuirks)
	case "QUIRKS_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("QUIRKS_LOOP_END"), ex.includeQuirks)
	case "PERKS_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includePerks)
	case "PERKS_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("PERKS_LOOP_END"), ex.includePerks)
	case "LANGUAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeLanguages)
	case "LANGUAGES_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("LANGUAGES_LOOP_END"), ex.includeLanguages)
	case "CULTURAL_FAMILIARITIES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeCulturalFamiliarities)
	case "CULTURAL_FAMILIARITIES_LOOP_START":
		ex.processAdvantagesLoop(ex.extractUpToMarker("CULTURAL_FAMILIARITIES_LOOP_END"), ex.includeCulturalFamiliarities)
	case "SKILLS_LOOP_COUNT":
		count := 0
		gurps.TraverseSkills(func(_ *gurps.Skill) bool {
			count++
			return false
		}, ex.entity.Skills...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "SKILLS_LOOP_START":
		ex.processSkillsLoop(ex.extractUpToMarker("SKILLS_LOOP_END"))
	case "SPELLS_LOOP_COUNT":
		count := 0
		gurps.TraverseSpells(func(_ *gurps.Spell) bool {
			count++
			return false
		}, ex.entity.Spells...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "SPELLS_LOOP_START":
		ex.processSpellsLoop(ex.extractUpToMarker("SPELLS_LOOP_END"))
	case "MELEE_LOOP_COUNT", "HIERARCHICAL_MELEE_LOOP_COUNT":
		// TODO: Is the hierarchical one right? It is what the old code did... but doesn't seem right
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.EquippedWeapons(weapon.Melee))))
	case "MELEE_LOOP_START":
		// TODO: Implement
	case "HIERARCHICAL_MELEE_LOOP_START":
		// TODO: Implement
	case "RANGED_LOOP_COUNT", "HIERARCHICAL_RANGED_LOOP_COUNT":
		// TODO: Is the hierarchical one right? It is what the old code did... but doesn't seem right
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.EquippedWeapons(weapon.Ranged))))
	case "RANGED_LOOP_START":
		// TODO: Implement
	case "HIERARCHICAL_RANGED_LOOP_START":
		// TODO: Implement
	case "EQUIPMENT_LOOP_COUNT":
		count := 0
		gurps.TraverseEquipment(func(eqp *gurps.Equipment) bool {
			if ex.includeByCategories(eqp.Categories) {
				count++
			}
			return false
		}, ex.entity.CarriedEquipment...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "EQUIPMENT_LOOP_START":
		// TODO: Implement
	case "OTHER_EQUIPMENT_LOOP_COUNT":
		count := 0
		gurps.TraverseEquipment(func(eqp *gurps.Equipment) bool {
			if ex.includeByCategories(eqp.Categories) {
				count++
			}
			return false
		}, ex.entity.OtherEquipment...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "OTHER_EQUIPMENT_LOOP_START":
		// TODO: Implement
	case "NOTES_LOOP_COUNT":
		count := 0
		gurps.TraverseNotes(func(_ *gurps.Note) bool {
			count++
			return false
		}, ex.entity.Notes...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "NOTES_LOOP_START":
		// TODO: Implement
	case "REACTION_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.Reactions())))
	case "REACTION_LOOP_START":
		// TODO: Implement
	case "CONDITIONAL_MODIFIERS_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.ConditionalModifiers())))
	case "CONDITIONAL_MODIFIERS_LOOP_START":
		// TODO: Implement
	case "PRIMARY_ATTRIBUTE_LOOP_COUNT":
		count := 0
		for _, def := range ex.entity.SheetSettings.Attributes.List() {
			if def.Type != attribute.Pool && def.Primary() {
				if _, exists := ex.entity.Attributes.Set[def.DefID]; exists {
					count++
				}
			}
		}
		ex.writeEncodedText(strconv.Itoa(count))
	case "PRIMARY_ATTRIBUTE_LOOP_START":
		// TODO: Implement
	case "SECONDARY_ATTRIBUTE_LOOP_COUNT":
		count := 0
		for _, def := range ex.entity.SheetSettings.Attributes.List() {
			if def.Type != attribute.Pool && !def.Primary() {
				if _, exists := ex.entity.Attributes.Set[def.DefID]; exists {
					count++
				}
			}
		}
		ex.writeEncodedText(strconv.Itoa(count))
	case "SECONDARY_ATTRIBUTE_LOOP_START":
		// TODO: Implement
	case "POINT_POOL_LOOP_COUNT":
		count := 0
		for _, def := range ex.entity.SheetSettings.Attributes.List() {
			if def.Type == attribute.Pool {
				if _, exists := ex.entity.Attributes.Set[def.DefID]; exists {
					count++
				}
			}
		}
		ex.writeEncodedText(strconv.Itoa(count))
	case "POINT_POOL_LOOP_START":
		// TODO: Implement
	case "CONTINUE_ID", "CAMPAIGN", "OPTIONS_CODE":
		// No-op
	default:
		switch {
		case strings.HasPrefix(key, "ONLY_CATEGORIES_"):
			splitIntoMap(key, "ONLY_CATEGORIES_", ex.onlyCategories)
		case strings.HasPrefix(key, "EXCLUDE_CATEGORIES_"):
			splitIntoMap(key, "EXCLUDE_CATEGORIES_", ex.excludedCategories)
		case strings.HasPrefix(key, "COLOR_"):
			// TODO: Implement
			/*
				String colorKey = key.substring(KEY_COLOR_PREFIX.length()).toLowerCase();
				for (ThemeColor one : Colors.ALL) {
					if (colorKey.equals(one.getKey())) {
						out.write(Colors.encodeToHex(one));
					}
				}
			*/
		default:
			attrKey := strings.ToLower(key)
			if attr, ok := ex.entity.Attributes.Set[attrKey]; ok {
				if def := attr.AttributeDef(); def != nil {
					ex.writeEncodedText(attr.Maximum().String())
					return nil
				}
			}
			switch {
			case strings.HasSuffix(attrKey, "_name"):
				if attr, ok := ex.entity.Attributes.Set[attrKey[:len(attrKey)-len("_name")]]; ok {
					if def := attr.AttributeDef(); def != nil {
						ex.writeEncodedText(def.Name)
						return nil
					}
				}
			case strings.HasSuffix(attrKey, "_full_name"):
				if attr, ok := ex.entity.Attributes.Set[attrKey[:len(attrKey)-len("_full_name")]]; ok {
					if def := attr.AttributeDef(); def != nil {
						ex.writeEncodedText(def.ResolveFullName())
						return nil
					}
				}
			case strings.HasSuffix(attrKey, "_combined_name"):
				if attr, ok := ex.entity.Attributes.Set[attrKey[:len(attrKey)-len("_combined_name")]]; ok {
					if def := attr.AttributeDef(); def != nil {
						ex.writeEncodedText(def.CombinedName())
						return nil
					}
				}
			case strings.HasSuffix(attrKey, "_points"):
				if attr, ok := ex.entity.Attributes.Set[attrKey[:len(attrKey)-len("_points")]]; ok {
					if def := attr.AttributeDef(); def != nil {
						ex.writeEncodedText(attr.PointCost().String())
						return nil
					}
				}
			case strings.HasSuffix(attrKey, "_current"):
				if attr, ok := ex.entity.Attributes.Set[attrKey[:len(attrKey)-len("_current")]]; ok {
					if def := attr.AttributeDef(); def != nil {
						ex.writeEncodedText(attr.Current().String())
						return nil
					}
				}
			}
			ex.unidentifiedKey(key)
		}
	}
	return nil
}

func (ex *legacyExporter) unidentifiedKey(key string) {
	ex.writeEncodedText(fmt.Sprintf(`Unidentified key: %q`, key))
}

func splitIntoMap(in, prefix string, m map[string]bool) {
	for _, one := range strings.Split(in[len(prefix):], "_") {
		if one != "" {
			m[one] = true
		}
	}
}

func (ex *legacyExporter) extractUpToMarker(marker string) []byte {
	remaining := ex.template[ex.pos:]
	i := bytes.Index(remaining, []byte(marker))
	if i == -1 {
		ex.pos = len(ex.template)
		return remaining
	}
	buffer := ex.template[ex.pos : ex.pos+i]
	ex.pos += i + len(marker)
	if ex.enhancedKeyParsing && ex.pos < len(ex.template) && ex.template[ex.pos] == '@' {
		ex.pos++
	}
	return buffer
}

func (ex *legacyExporter) writeEncodedText(text string) {
	if ex.encodeText {
		for _, ch := range text {
			switch ch {
			case '<':
				ex.out.WriteString("&lt;")
			case '>':
				ex.out.WriteString("&gt;")
			case '&':
				ex.out.WriteString("&amp;")
			case '"':
				ex.out.WriteString("&quot;")
			case '\'':
				ex.out.WriteString("&apos;")
			case '\n':
				ex.out.WriteString("<br>")
			default:
				if ch >= ' ' && ch <= '~' {
					ex.out.WriteRune(ch)
				} else {
					ex.out.WriteString("&#")
					ex.out.WriteString(strconv.Itoa(int(ch)))
					ex.out.WriteByte(';')
				}
			}
		}
	} else {
		ex.out.WriteString(text)
	}
}

func (ex *legacyExporter) bestWeaponDefense(f func(weapon *gurps.Weapon) string) string {
	best := "-"
	bestValue := fixed.F64d4Min
	for _, w := range ex.entity.EquippedWeapons(weapon.Melee) {
		if s := f(w); s != "" && !strings.EqualFold(s, "no") {
			if v, rem := fxp.Extract(s); v != 0 || rem != s {
				if bestValue < v {
					bestValue = v
					best = s
				}
			}
		}
	}
	return best
}

func (ex *legacyExporter) writeAdvantageLoopCount(f func(*gurps.Advantage) bool) {
	count := 0
	gurps.TraverseAdvantages(func(adq *gurps.Advantage) bool {
		if f(adq) {
			count++
		}
		return false
	}, true, ex.entity.Advantages...)
	ex.writeEncodedText(strconv.Itoa(count))
}

func (ex *legacyExporter) includeByCategories(categories []string) bool {
	for cat := range ex.onlyCategories {
		if gurps.HasCategory(cat, categories) {
			return true
		}
	}
	if len(ex.onlyCategories) != 0 {
		return false
	}
	for cat := range ex.excludedCategories {
		if gurps.HasCategory(cat, categories) {
			return false
		}
	}
	return true
}

func (ex *legacyExporter) includeByAdvantageCategories(adq *gurps.Advantage) bool {
	return ex.includeByCategories(adq.Categories)
}

func (ex *legacyExporter) includeAdvantages(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() > fxp.One && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includePerks(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() == fxp.One && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeAdvantagesAndPerks(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() > 0 && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeDisadvantages(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() < fxp.NegOne && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeQuirks(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() == fxp.NegOne && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeDisadvantagesAndQuirks(adq *gurps.Advantage) bool {
	return adq.AdjustedPoints() < 0 && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeLanguages(adq *gurps.Advantage) bool {
	return gurps.HasCategory("Language", adq.Categories) && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) includeCulturalFamiliarities(adq *gurps.Advantage) bool {
	return strings.HasPrefix(strings.ToLower(adq.Name), "cultural familiarity (") && ex.includeByAdvantageCategories(adq)
}

func (ex *legacyExporter) processEncumbranceLoop(buffer []byte) {
	for _, enc := range datafile.AllEncumbrance {
		ex.processBuffer(buffer, func(key string) {
			switch key {
			case "CURRENT_MARKER":
				if enc == ex.entity.EncumbranceLevel(false) {
					ex.writeEncodedText("current")
				}
			case "CURRENT_MARKER_1":
				if enc == ex.entity.EncumbranceLevel(false) {
					ex.writeEncodedText("1")
				}
			case "CURRENT_MARKER_BULLET":
				if enc == ex.entity.EncumbranceLevel(false) {
					ex.writeEncodedText("•")
				}
			case "LEVEL":
				if enc == ex.entity.EncumbranceLevel(false) {
					ex.writeEncodedText("• ")
				}
				fallthrough
			case "LEVEL_NO_MARKER":
				ex.writeEncodedText(fmt.Sprintf("%s (%s)", enc.String(), (-enc.Penalty()).String()))
			case "LEVEL_ONLY":
				ex.writeEncodedText((-enc.Penalty()).String())
			case "MAX_LOAD":
				ex.writeEncodedText(ex.entity.SheetSettings.DefaultWeightUnits.Format(ex.entity.MaximumCarry(enc)))
			case "MOVE":
				ex.writeEncodedText(strconv.Itoa(ex.entity.Move(enc)))
			case "DODGE":
				ex.writeEncodedText(strconv.Itoa(ex.entity.Dodge(enc)))
			default:
				ex.unidentifiedKey(key)
			}
		})
	}
}

func (ex *legacyExporter) processHitLocationLoop(buffer []byte) {
	for i, location := range ex.entity.SheetSettings.HitLocations.Locations {
		ex.processBuffer(buffer, func(key string) {
			switch key {
			case "ID":
				ex.writeEncodedText(strconv.Itoa(i))
			case "ROLL":
				ex.writeEncodedText(location.RollRange)
			case "WHERE":
				ex.writeEncodedText(location.TableName)
			case "PENALTY":
				ex.writeEncodedText(strconv.Itoa(location.HitPenalty))
			case "DR":
				ex.writeEncodedText(location.DisplayDR(ex.entity, nil))
			case "DR_TOOLTIP":
				var tooltip xio.ByteBuffer
				location.DisplayDR(ex.entity, &tooltip)
				ex.writeEncodedText(tooltip.String())
			case "EQUIPMENT":
				ex.writeEncodedText(strings.Join(ex.hitLocationEquipment(location), ", "))
			case "EQUIPMENT_FORMATTED":
				for _, one := range ex.hitLocationEquipment(location) {
					ex.out.WriteString("<p>")
					ex.writeEncodedText(one)
					ex.out.WriteString("</p>\n")
				}
			default:
				ex.unidentifiedKey(key)
			}
		})
	}
}

func (ex *legacyExporter) hitLocationEquipment(location *gurps.HitLocation) []string {
	var list []string
	gurps.TraverseEquipment(func(eqp *gurps.Equipment) bool {
		if eqp.Equipped {
			for _, f := range eqp.Features {
				if bonus, ok := f.(*feature.DRBonus); ok {
					if strings.EqualFold(location.LocID, bonus.Location) {
						list = append(list, eqp.Name)
					}
				}
			}
		}
		return false
	}, ex.entity.CarriedEquipment...)
	return list
}

func (ex *legacyExporter) processAdvantagesLoop(buffer []byte, f func(*gurps.Advantage) bool) {
	gurps.TraverseAdvantages(func(adq *gurps.Advantage) bool {
		if f(adq) {
			ex.processBuffer(buffer, func(key string) {
				switch key {
				case "ID":
					ex.writeEncodedText(adq.ID.String())
				case parentIDKey:
					if adq.Parent != nil {
						ex.writeEncodedText(adq.Parent.ID.String())
					}
				case typeKey:
					if adq.Container() {
						ex.writeEncodedText(strings.ToUpper(adq.ContainerType.Key()))
					} else {
						ex.writeEncodedText("ITEM")
					}
				case pointsKey:
					ex.writeEncodedText(adq.AdjustedPoints().String())
				case descriptionKey:
					ex.writeEncodedText(adq.String())
					ex.writeNote(adq.ModifierNotes())
					ex.writeNote(adq.Notes())
				case descriptionPrimaryKey:
					ex.writeEncodedText(adq.String())
				case "DESCRIPTION_USER":
					ex.writeEncodedText(adq.UserDesc)
				case "DESCRIPTION_USER_FORMATTED":
					if adq.UserDesc != "" {
						for _, one := range strings.Split(adq.UserDesc, "\n") {
							ex.out.WriteString("<p>")
							ex.writeEncodedText(one)
							ex.out.WriteString("</p>\n")
						}
					}
				case refKey:
					ex.writeEncodedText(adq.PageRef)
				case styleIndentWarningKey:
					ex.handleStyleIndentWarning(adq.Depth(), adq.Satisfied)
				case satisfiedKey:
					ex.handleSatisfied(adq.Satisfied)
				default:
					switch {
					case strings.HasPrefix(key, "DESCRIPTION_MODIFIER_NOTES"):
						ex.writeWithOptionalParens(key, adq.ModifierNotes())
					case strings.HasPrefix(key, "DESCRIPTION_NOTES"):
						ex.writeWithOptionalParens(key, adq.Notes())
					case strings.HasPrefix(key, "MODIFIER_NOTES_FOR_"):
						if mod := adq.ActiveModifierFor(key[len("MODIFIER_NOTES_FOR_"):]); mod != nil {
							ex.writeEncodedText(mod.Notes)
						}
					case strings.HasPrefix(key, "PREFIX_DEPTHx"):
						ex.handlePrefixDepth(key, adq.Depth())
					default:
						ex.unidentifiedKey(key)
					}
				}
			})
		}
		return false
	}, true, ex.entity.Advantages...)
}

func (ex *legacyExporter) processSkillsLoop(buffer []byte) {
	gurps.TraverseSkills(func(s *gurps.Skill) bool {
		ex.processBuffer(buffer, func(key string) {
			switch key {
			case "ID":
				ex.writeEncodedText(s.ID.String())
			case parentIDKey:
				if s.Parent != nil {
					ex.writeEncodedText(s.Parent.ID.String())
				}
			case typeKey:
				if s.Container() {
					ex.writeEncodedText("GROUP")
				} else {
					ex.writeEncodedText("ITEM")
				}
			case pointsKey:
				ex.writeEncodedText(s.AdjustedPoints().String())
			case descriptionKey:
				ex.writeEncodedText(s.String())
				ex.writeNote(s.ModifierNotes())
				ex.writeNote(s.Notes())
			case descriptionPrimaryKey:
				ex.writeEncodedText(s.String())
			case "SL":
				ex.writeEncodedText(s.LevelAsString())
			case "RSL":
				ex.writeEncodedText(s.RelativeLevel())
			case "DIFFICULTY":
				if !s.Container() {
					ex.writeEncodedText(s.Difficulty.String())
				}
			case refKey:
				ex.writeEncodedText(s.PageRef)
			case styleIndentWarningKey:
				ex.handleStyleIndentWarning(s.Depth(), s.Satisfied)
			case satisfiedKey:
				ex.handleSatisfied(s.Satisfied)
			default:
				switch {
				case strings.HasPrefix(key, "DESCRIPTION_MODIFIER_NOTES"):
					ex.writeWithOptionalParens(key, s.ModifierNotes())
				case strings.HasPrefix(key, "DESCRIPTION_NOTES"):
					ex.writeWithOptionalParens(key, s.Notes())
				case strings.HasPrefix(key, "PREFIX_DEPTHx"):
					ex.handlePrefixDepth(key, s.Depth())
				default:
					ex.unidentifiedKey(key)
				}
			}
		})
		return false
	}, ex.entity.Skills...)
}

func (ex *legacyExporter) processSpellsLoop(buffer []byte) {
	gurps.TraverseSpells(func(s *gurps.Spell) bool {
		ex.processBuffer(buffer, func(key string) {
			switch key {
			case "ID":
				ex.writeEncodedText(s.ID.String())
			case parentIDKey:
				if s.Parent != nil {
					ex.writeEncodedText(s.Parent.ID.String())
				}
			case typeKey:
				if s.Container() {
					ex.writeEncodedText("GROUP")
				} else {
					ex.writeEncodedText("ITEM")
				}
			case pointsKey:
				ex.writeEncodedText(s.AdjustedPoints().String())
			case descriptionKey:
				ex.writeEncodedText(s.String())
				ex.writeNote(s.Notes())
				ex.writeNote(s.Rituals())
			case descriptionPrimaryKey:
				ex.writeEncodedText(s.String())
			case "SL":
				ex.writeEncodedText(s.LevelAsString())
			case "RSL":
				ex.writeEncodedText(s.RelativeLevel())
			case "DIFFICULTY":
				if !s.Container() {
					ex.writeEncodedText(s.Difficulty.String())
				}
			case "CLASS":
				ex.writeEncodedText(s.Class)
			case "COLLEGE":
				ex.writeEncodedText(strings.Join(s.College, ", "))
			case "MANA_CAST":
				ex.writeEncodedText(s.CastingCost)
			case "MANA_MAINTAIN":
				ex.writeEncodedText(s.MaintenanceCost)
			case "TIME_CAST":
				ex.writeEncodedText(s.CastingTime)
			case "DURATION":
				ex.writeEncodedText(s.Duration)
			case "RESIST":
				ex.writeEncodedText(s.Resist)
			case refKey:
				ex.writeEncodedText(s.PageRef)
			case styleIndentWarningKey:
				ex.handleStyleIndentWarning(s.Depth(), s.Satisfied)
			case satisfiedKey:
				ex.handleSatisfied(s.Satisfied)
			default:
				switch {
				case strings.HasPrefix(key, "DESCRIPTION_NOTES"):
					notes := s.Notes()
					rituals := s.Rituals()
					if rituals != "" {
						if strings.TrimSpace(notes) != "" {
							notes += "; "
						}
						notes += rituals
					}
					ex.writeWithOptionalParens(key, notes)
				case strings.HasPrefix(key, "PREFIX_DEPTHx"):
					ex.handlePrefixDepth(key, s.Depth())
				default:
					ex.unidentifiedKey(key)
				}
			}
		})
		return false
	}, ex.entity.Spells...)
}

func (ex *legacyExporter) handleStyleIndentWarning(depth int, satisfied bool) {
	var style strings.Builder
	if depth > 0 {
		fmt.Fprintf(&style, "padding-left: %dpx;", depth*12)
	}
	if !satisfied {
		style.WriteString("color: red;")
	}
	if style.Len() != 0 {
		ex.out.WriteString(` style="`)
		ex.out.WriteString(style.String())
		ex.out.WriteString(`" `)
	}
}

func (ex *legacyExporter) handleSatisfied(satisfied bool) {
	if satisfied {
		ex.writeEncodedText("Y")
	} else {
		ex.writeEncodedText("N")
	}
}

func (ex *legacyExporter) handlePrefixDepth(key string, depth int) {
	if amt, err := strconv.Atoi(key[len("PREFIX_DEPTHx"):]); err != nil {
		ex.unidentifiedKey(key)
	} else {
		ex.writeEncodedText(strconv.Itoa(amt * depth))
	}
}

func (ex *legacyExporter) writeNote(note string) {
	if strings.TrimSpace(note) != "" {
		ex.out.WriteString(`<div class="note">`)
		ex.writeEncodedText(note)
		ex.out.WriteString("</div>")
	}
}

func (ex *legacyExporter) writeWithOptionalParens(key, text string) {
	if strings.TrimSpace(text) != "" {
		var pre, post string
		switch {
		case strings.HasSuffix(key, "_PAREN"):
			pre = " ("
			post = ")"
		case strings.HasSuffix(key, "_BRACKET"):
			pre = " ["
			post = "]"
		case strings.HasSuffix(key, "_CURLY"):
			pre = " {"
			post = "}"
		}
		ex.out.WriteString(pre)
		ex.out.WriteString(text)
		ex.out.WriteString(post)
	}
}

func (ex *legacyExporter) processBuffer(buffer []byte, f func(key string)) {
	var keyBuffer bytes.Buffer
	lookForKeyMarker := true
	i := 0
	for i < len(buffer) {
		ch := buffer[i]
		i++
		switch {
		case lookForKeyMarker:
			if ch == '@' {
				lookForKeyMarker = false
			} else {
				ex.out.WriteByte(ch)
			}
		case ch == '_' || (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
			keyBuffer.WriteByte(ch)
		default:
			if !ex.enhancedKeyParsing || ch != '@' {
				i--
			}
			f(keyBuffer.String())
			keyBuffer.Reset()
			lookForKeyMarker = true
		}
	}
}
