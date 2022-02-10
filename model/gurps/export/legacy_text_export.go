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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/errs"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

type legacyExporter struct {
	entity             *gurps.Entity
	template           []byte
	pos                int
	exportPath         string
	onlyCategories     map[string]bool
	excludedCategories map[string]bool
	keyBuffer          bytes.Buffer
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
			ex.keyBuffer.WriteByte(ch)
		default:
			if !ex.enhancedKeyParsing || ch != '@' {
				ex.pos--
			}
			if err = ex.emitKey(); err != nil {
				return err
			}
			ex.keyBuffer.Reset()
			lookForKeyMarker = true
		}
	}
	if ex.keyBuffer.Len() != 0 {
		if err = ex.emitKey(); err != nil {
			return err
		}
	}
	return nil
}

func (ex *legacyExporter) emitKey() error {
	key := ex.keyBuffer.String()
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
	case "HIT_LOCATION_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.SheetSettings.HitLocations.Locations)))
	case "ADVANTAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeByAdvantageCategories)
	case "ADVANTAGES_ALL_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeAdvantagesAndPerks)
	case "ADVANTAGES_ONLY_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeAdvantages)
	case "DISADVANTAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeDisadvantages)
	case "DISADVANTAGES_ALL_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeDisadvantagesAndQuirks)
	case "QUIRKS_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeQuirks)
	case "PERKS_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includePerks)
	case "LANGUAGES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeLanguages)
	case "CULTURAL_FAMILIARITIES_LOOP_COUNT":
		ex.writeAdvantageLoopCount(ex.includeCulturalFamiliarities)
	case "SKILLS_LOOP_COUNT":
		count := 0
		gurps.TraverseSkills(func(_ *gurps.Skill) bool {
			count++
			return false
		}, ex.entity.Skills...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "SPELLS_LOOP_COUNT":
		count := 0
		gurps.TraverseSpells(func(_ *gurps.Spell) bool {
			count++
			return false
		}, ex.entity.Spells...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "MELEE_LOOP_COUNT", "HIERARCHICAL_MELEE_LOOP_COUNT":
		// TODO: Is the hierarchical one right? It is what the old code did... but doesn't seem right
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.EquippedWeapons(weapon.Melee))))
	case "RANGED_LOOP_COUNT", "HIERARCHICAL_RANGED_LOOP_COUNT":
		// TODO: Is the hierarchical one right? It is what the old code did... but doesn't seem right
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.EquippedWeapons(weapon.Ranged))))
	case "EQUIPMENT_LOOP_COUNT":
		count := 0
		gurps.TraverseEquipment(func(eqp *gurps.Equipment) bool {
			if ex.includeByCategories(eqp.Categories) {
				count++
			}
			return false
		}, ex.entity.CarriedEquipment...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "OTHER_EQUIPMENT_LOOP_COUNT":
		count := 0
		gurps.TraverseEquipment(func(eqp *gurps.Equipment) bool {
			if ex.includeByCategories(eqp.Categories) {
				count++
			}
			return false
		}, ex.entity.OtherEquipment...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "NOTES_LOOP_COUNT":
		count := 0
		gurps.TraverseNotes(func(_ *gurps.Note) bool {
			count++
			return false
		}, ex.entity.Notes...)
		ex.writeEncodedText(strconv.Itoa(count))
	case "REACTION_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.Reactions())))
	case "CONDITIONAL_MODIFIERS_LOOP_COUNT":
		ex.writeEncodedText(strconv.Itoa(len(ex.entity.ConditionalModifiers())))
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
	case "ENCUMBRANCE_LOOP_START":
		ex.processEncumbranceLoop(ex.extractUpToMarker("ENCUMBRANCE_LOOP_END"))
	case "HIT_LOCATION_LOOP_START":
	case "ADVANTAGES_LOOP_START":
	case "ADVANTAGES_ALL_LOOP_START":
	case "ADVANTAGES_ONLY_LOOP_START":
	case "DISADVANTAGES_LOOP_START":
	case "DISADVANTAGES_ALL_LOOP_START":
	case "QUIRKS_LOOP_START":
	case "PERKS_LOOP_START":
	case "LANGUAGES_LOOP_START":
	case "CULTURAL_FAMILIARITIES_LOOP_START":
	case "SKILLS_LOOP_START":
	case "SPELLS_LOOP_START":
	case "MELEE_LOOP_START":
	case "HIERARCHICAL_MELEE_LOOP_START":
	case "RANGED_LOOP_START":
	case "HIERARCHICAL_RANGED_LOOP_START":
	case "EQUIPMENT_LOOP_START":
	case "OTHER_EQUIPMENT_LOOP_START":
	case "NOTES_LOOP_START":
	case "REACTION_LOOP_START":
	case "CONDITIONAL_MODIFIERS_LOOP_START":
	case "PRIMARY_ATTRIBUTE_LOOP_START":
	case "SECONDARY_ATTRIBUTE_LOOP_START":
	case "POINT_POOL_LOOP_START":
	case "CONTINUE_ID", "CAMPAIGN", "OPTIONS_CODE":
		// No-op
	default:
		switch {
		case strings.HasPrefix(key, "ONLY_CATEGORIES_"):
		case strings.HasPrefix(key, "EXCLUDE_CATEGORIES_"):
		case strings.HasPrefix(key, "COLOR_"):
		default:
			/*
				if (!processAttributeKeys(out, gurpsCharacter, key)) {
					writeEncodedText(out, String.format(UNIDENTIFIED_KEY, key));
				}
			*/
		}
		/*
		   if (!checkForLoopKeys(in, out, key)) {
		       if (key.startsWith(KEY_ONLY_CATEGORIES)) {
		           setOnlyCategories(key);
		       } else if (key.startsWith(KEY_EXCLUDE_CATEGORIES)) {
		           setExcludeCategories(key);
		       } else if (key.startsWith(KEY_COLOR_PREFIX)) {
		           String colorKey = key.substring(KEY_COLOR_PREFIX.length()).toLowerCase();
		           for (ThemeColor one : Colors.ALL) {
		               if (colorKey.equals(one.getKey())) {
		                   out.write(Colors.encodeToHex(one));
		               }
		           }
		       } else if (!processAttributeKeys(out, gurpsCharacter, key)) {
		           writeEncodedText(out, String.format(UNIDENTIFIED_KEY, key));
		       }
		   }
		*/
	}
	return nil
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
	// TODO: Implement
}
