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
	mark               int
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
		mark:               -1,
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
				ex.mark = ex.pos
			} else {
				ex.out.WriteByte(ch)
			}
		case ch == '_' || (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
			ex.keyBuffer.WriteByte(ch)
			ex.mark = ex.pos
		default:
			if !ex.enhancedKeyParsing || ch != '@' {
				ex.pos = ex.mark
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
	switch ex.keyBuffer.String() {
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
		/*
		       case KEY_TIRED_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "fp", I18n.text("Tired"));
		           break;
		       case KEY_FP_COLLAPSE_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "fp", I18n.text("Collapse"));
		           break;
		       case KEY_UNCONSCIOUS_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "fp", I18n.text("Unconscious"));
		           break;
		       case KEY_REELING_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", I18n.text("Reeling"));
		           break;
		       case KEY_HP_COLLAPSE_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", I18n.text("Collapse"));
		           break;
		       case KEY_DEATH_CHECK_1_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", String.format(I18n.text("Dying #%d"), Integer.valueOf(1)));
		           break;
		       case KEY_DEATH_CHECK_2_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", String.format(I18n.text("Dying #%d"), Integer.valueOf(2)));
		           break;
		       case KEY_DEATH_CHECK_3_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", String.format(I18n.text("Dying #%d"), Integer.valueOf(3)));
		           break;
		       case KEY_DEATH_CHECK_4_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", String.format(I18n.text("Dying #%d"), Integer.valueOf(4)));
		           break;
		       case KEY_DEAD_DEPRECATED:
		           deprecatedWritePointPoolThreshold(out, gurpsCharacter, "hp", I18n.text("Dead"));
		           break;
		       case KEY_BASIC_LIFT:
		           writeEncodedText(out, gurpsCharacter.getBasicLift().toString());
		           break;
		       case KEY_ONE_HANDED_LIFT:
		           writeEncodedText(out, gurpsCharacter.getOneHandedLift().toString());
		           break;
		       case KEY_TWO_HANDED_LIFT:
		           writeEncodedText(out, gurpsCharacter.getTwoHandedLift().toString());
		           break;
		       case KEY_SHOVE:
		           writeEncodedText(out, gurpsCharacter.getShoveAndKnockOver().toString());
		           break;
		       case KEY_RUNNING_SHOVE:
		           writeEncodedText(out, gurpsCharacter.getRunningShoveAndKnockOver().toString());
		           break;
		       case KEY_CARRY_ON_BACK:
		           writeEncodedText(out, gurpsCharacter.getCarryOnBack().toString());
		           break;
		       case KEY_SHIFT_SLIGHTLY:
		           writeEncodedText(out, gurpsCharacter.getShiftSlightly().toString());
		           break;
		       case KEY_CARRIED_WEIGHT:
		           writeEncodedText(out, EquipmentColumn.getDisplayWeight(gurpsCharacter, gurpsCharacter.getWeightCarried(false)));
		           break;
		       case KEY_CARRIED_VALUE:
		           writeEncodedText(out, "$" + gurpsCharacter.getWealthCarried().toLocalizedString());
		           break;
		       case KEY_OTHER_VALUE:
		           writeEncodedText(out, "$" + gurpsCharacter.getWealthNotCarried().toLocalizedString());
		           break;
		       case KEY_ALL_NOTES_COMBINED:
		           StringBuilder buffer = new StringBuilder();
		           for (Note note : gurpsCharacter.getNotesIterator()) {
		               if (!buffer.isEmpty()) {
		                   buffer.append("\n\n");
		               }
		               buffer.append(note.getDescription());
		           }
		           writeEncodedText(out, buffer.toString());
		           break;
		       case KEY_RACE_DEPRECATED:
		           // Use ancestry instead
		           break;
		       case KEY_BODY_TYPE:
		           writeEncodedText(out, gurpsCharacter.getSheetSettings().getHitLocations().getName());
		           break;
		       default:
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
		           break;
		   }
		*/
	case "CONTINUE_ID", "CAMPAIGN", "OPTIONS_CODE":
		// No-op
	}
	return nil
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
