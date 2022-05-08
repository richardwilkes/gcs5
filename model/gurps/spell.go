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
	"fmt"
	"io/fs"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

var _ node.Node = &Spell{}

// Columns that can be used with the spell method .CellData()
const (
	SpellDescriptionColumn = iota
	SpellResistColumn
	SpellClassColumn
	SpellCollegeColumn
	SpellCastCostColumn
	SpellMaintainCostColumn
	SpellCastTimeColumn
	SpellDurationColumn
	SpellDifficultyColumn
	SpellTagsColumn
	SpellReferenceColumn
	SpellLevelColumn
	SpellRelativeLevelColumn
	SpellPointsColumn
	SpellDescriptionForPageColumn
)

const spellListTypeKey = "spell_list"

// Spell holds the data for a spell.
type Spell struct {
	SpellData
	Entity            *Entity
	Parent            *Spell
	LevelData         skill.Level
	UnsatisfiedReason string
	Satisfied         bool
}

type spellListData struct {
	Type    string   `json:"type"`
	Version int      `json:"version"`
	Rows    []*Spell `json:"rows"`
}

// NewSpellsFromFile loads an Spell list from a file.
func NewSpellsFromFile(fileSystem fs.FS, filePath string) ([]*Spell, error) {
	var data spellListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != spellListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveSpells writes the Spell list to the file as JSON.
func SaveSpells(spells []*Spell, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &spellListData{
		Type:    spellListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    spells,
	})
}

// NewSpell creates a new Spell.
func NewSpell(entity *Entity, parent *Spell, container bool) *Spell {
	s := newSpell(entity, parent, gid.Spell, container)
	s.UpdateLevel()
	return s
}

// NewRitualMagicSpell creates a new Ritual Magic Spell.
func NewRitualMagicSpell(entity *Entity, parent *Spell) *Spell {
	s := newSpell(entity, parent, gid.RitualMagicSpell, false)
	s.RitualSkillName = "Ritual Magic"
	s.Points = 0
	s.UpdateLevel()
	return s
}

func newSpell(entity *Entity, parent *Spell, typeKey string, container bool) *Spell {
	s := Spell{
		SpellData: SpellData{
			ContainerBase: newContainerBase[*Spell](typeKey, container),
		},
		Entity: entity,
		Parent: parent,
	}
	if !container {
		s.Difficulty.Attribute = AttributeIDFor(entity, gid.Intelligence)
		s.Difficulty.Difficulty = skill.Hard
		s.PowerSource = i18n.Text("Arcane")
		s.Class = i18n.Text("Regular")
		s.CastingCost = "1"
		s.CastingTime = "1 sec"
		s.Duration = "Instant"
		s.Points = fxp.One
	}
	s.Name = s.Kind()
	return &s
}

// MarshalJSON implements json.Marshaler.
func (s *Spell) MarshalJSON() ([]byte, error) {
	s.ClearUnusedFieldsForType()
	if s.Container() || s.LevelData.Level <= 0 {
		return json.Marshal(&s.SpellData)
	}
	type calc struct {
		Level              fxp.Int `json:"level"`
		RelativeSkillLevel string  `json:"rsl"`
	}
	return json.Marshal(&struct {
		SpellData
		Calc calc `json:"calc"`
	}{
		SpellData: s.SpellData,
		Calc: calc{
			Level:              s.LevelData.Level,
			RelativeSkillLevel: s.RelativeLevel(),
		},
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Spell) UnmarshalJSON(data []byte) error {
	var localData struct {
		SpellData
		// Old data fields
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(data, &localData); err != nil {
		return err
	}
	localData.ClearUnusedFieldsForType()
	s.SpellData = localData.SpellData
	s.Tags = convertOldCategoriesToTags(s.Tags, localData.Categories)
	slices.Sort(s.Tags)
	if s.Container() {
		for _, one := range s.Children {
			one.Parent = s
		}
	}
	return nil
}

// CellData returns the cell data information for the given column.
func (s *Spell) CellData(column int, data *node.CellData) {
	switch column {
	case SpellDescriptionColumn:
		data.Type = node.Text
		data.Primary = s.Description()
		data.Secondary = s.SecondaryText()
	case SpellResistColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.Resist
		}
	case SpellClassColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.Class
		}
	case SpellCollegeColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = strings.Join(s.College, ", ")
		}
	case SpellCastCostColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.CastingCost
		}
	case SpellMaintainCostColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.MaintenanceCost
		}
	case SpellCastTimeColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.CastingTime
		}
	case SpellDurationColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.Duration
		}
	case SpellDifficultyColumn:
		if !s.Container() {
			data.Type = node.Text
			data.Primary = s.Difficulty.Description(s.Entity)
		}
	case SpellTagsColumn:
		data.Type = node.Text
		data.Primary = CombineTags(s.Tags)
	case SpellReferenceColumn:
		data.Type = node.PageRef
		data.Primary = s.PageRef
		data.Secondary = s.Name
	case SpellLevelColumn:
		if !s.Container() {
			data.Type = node.Text
			level := s.CalculateLevel()
			data.Primary = level.LevelAsString(s.Container())
			if level.Tooltip != "" {
				data.Tooltip = IncludesModifiersFrom + ":" + level.Tooltip
			}
			data.Alignment = unison.EndAlignment
		}
	case SpellRelativeLevelColumn:
		if !s.Container() {
			data.Type = node.Text
			rsl := s.AdjustedRelativeLevel()
			if rsl == fxp.Min {
				data.Primary = "-"
			} else {
				data.Primary = ResolveAttributeName(s.Entity, s.Difficulty.Attribute)
				if rsl != 0 {
					data.Primary += rsl.StringWithSign()
				}
			}
		}
	case SpellPointsColumn:
		data.Type = node.Text
		var tooltip xio.ByteBuffer
		data.Primary = s.AdjustedPoints(&tooltip).String()
		data.Alignment = unison.EndAlignment
		if tooltip.Len() != 0 {
			data.Tooltip = IncludesModifiersFrom + ":" + tooltip.String()
		}
	case SpellDescriptionForPageColumn:
		data.Type = node.Text
		data.Primary = s.Description()
		data.Secondary = s.SecondaryText()
		if !s.Container() {
			var buffer strings.Builder
			addPartToBuffer(&buffer, i18n.Text("Resistance"), s.Resist)
			addPartToBuffer(&buffer, i18n.Text("Class"), s.Class)
			addPartToBuffer(&buffer, i18n.Text("Cost"), s.CastingCost)
			addPartToBuffer(&buffer, i18n.Text("Maintain"), s.MaintenanceCost)
			addPartToBuffer(&buffer, i18n.Text("Time"), s.CastingTime)
			addPartToBuffer(&buffer, i18n.Text("Duration"), s.Duration)
			if buffer.Len() != 0 {
				if data.Secondary == "" {
					data.Secondary = buffer.String()
				} else {
					data.Secondary += "\n" + buffer.String()
				}
			}
		}
	}
}

func addPartToBuffer(buffer *strings.Builder, label, content string) {
	if content != "" && content != "-" {
		if buffer.Len() != 0 {
			buffer.WriteString("; ")
		}
		buffer.WriteString(label)
		buffer.WriteString(": ")
		buffer.WriteString(content)
	}
}

// Depth returns the number of parents this node has.
func (s *Spell) Depth() int {
	count := 0
	p := s.Parent
	for p != nil {
		count++
		p = p.Parent
	}
	return count
}

// RelativeLevel returns the adjusted relative level as a string.
func (s *Spell) RelativeLevel() string {
	if s.Container() || s.LevelData.Level <= 0 {
		return ""
	}
	rsl := s.AdjustedRelativeLevel()
	switch {
	case rsl == fxp.Min:
		return "-"
	case s.Type != gid.RitualMagicSpell:
		return ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
	default:
		return rsl.StringWithSign()
	}
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Spell) AdjustedRelativeLevel() fxp.Int {
	if s.Container() {
		return fxp.Min
	}
	if s.Entity != nil && s.CalculateLevel().Level > 0 {
		return s.LevelData.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fxp.Min
}

// UpdateLevel updates the level of the spell, returning true if it has changed.
func (s *Spell) UpdateLevel() bool {
	saved := s.LevelData
	if strings.HasPrefix(s.Type, gid.Spell) {
		s.LevelData = CalculateSpellLevel(s.Entity, s.Name, s.PowerSource, s.College, s.Tags, s.Difficulty,
			s.AdjustedPoints(nil))
	} else {
		s.LevelData = CalculateRitualMagicSpellLevel(s.Entity, s.Name, s.PowerSource, s.RitualSkillName,
			s.RitualPrereqCount, s.College, s.Tags, s.Difficulty, s.AdjustedPoints(nil))
	}
	return saved != s.LevelData
}

// CalculateLevel returns the computed level without updating it.
func (s *Spell) CalculateLevel() skill.Level {
	if strings.HasPrefix(s.Type, gid.Spell) {
		return CalculateSpellLevel(s.Entity, s.Name, s.PowerSource, s.College, s.Tags, s.Difficulty,
			s.AdjustedPoints(nil))
	}
	return CalculateRitualMagicSpellLevel(s.Entity, s.Name, s.PowerSource, s.RitualSkillName, s.RitualPrereqCount,
		s.College, s.Tags, s.Difficulty, s.AdjustedPoints(nil))
}

// IncrementSkillLevel adds enough points to increment the skill level to the next level.
func (s *Spell) IncrementSkillLevel() {
	if !s.Container() {
		basePoints := s.Points.Trunc() + fxp.One
		maxPoints := basePoints
		if s.Difficulty.Difficulty == skill.Wildcard {
			maxPoints += fxp.Twelve
		} else {
			maxPoints += fxp.Four
		}
		oldLevel := s.CalculateLevel().Level
		for points := basePoints; points < maxPoints; points += fxp.One {
			s.Points = points
			s.UpdateLevel()
			if s.CalculateLevel().Level > oldLevel {
				break
			}
		}
	}
}

// DecrementSkillLevel removes enough points to decrement the skill level to the previous level.
func (s *Spell) DecrementSkillLevel() {
	if !s.Container() && s.Points > 0 {
		basePoints := s.Points.Trunc()
		minPoints := basePoints
		if s.Difficulty.Difficulty == skill.Wildcard {
			minPoints -= fxp.Twelve
		} else {
			minPoints -= fxp.Four
		}
		minPoints = minPoints.Max(0)
		oldLevel := s.CalculateLevel().Level
		for points := basePoints; points >= minPoints; points -= fxp.One {
			s.Points = points
			s.UpdateLevel()
			if s.CalculateLevel().Level < oldLevel {
				break
			}
		}
		if s.Points > 0 {
			oldLevel = s.CalculateLevel().Level
			for s.Points > 0 {
				s.Points -= fxp.One
				s.UpdateLevel()
				if s.CalculateLevel().Level != oldLevel {
					s.Points += fxp.One
					break
				}
			}
		}
	}
}

// CalculateSpellLevel returns the calculated spell level.
func CalculateSpellLevel(entity *Entity, name, powerSource string, colleges, tags []string, difficulty AttributeDifficulty, pts fxp.Int) skill.Level {
	var tooltip xio.ByteBuffer
	relativeLevel := difficulty.Difficulty.BaseRelativeLevel()
	level := fxp.Min
	if entity != nil {
		pts = pts.Trunc()
		level = entity.ResolveAttributeCurrent(difficulty.Attribute)
		if difficulty.Difficulty == skill.Wildcard {
			pts = pts.Div(fxp.Three).Trunc()
		}
		switch {
		case pts < fxp.One:
			level = fxp.Min
			relativeLevel = 0
		case pts == fxp.One:
		// relativeLevel is preset to this point value
		case pts < fxp.Four:
			relativeLevel += fxp.One
		default:
			relativeLevel += fxp.One + pts.Div(fxp.Four).Trunc()
		}
		if level != fxp.Min {
			relativeLevel += entity.BestCollegeSpellBonus(tags, colleges, &tooltip)
			relativeLevel += entity.SpellBonusesFor(feature.SpellPowerSourceID, powerSource, tags, &tooltip)
			relativeLevel += entity.SpellBonusesFor(feature.SpellNameID, name, tags, &tooltip)
			relativeLevel = relativeLevel.Trunc()
			level += relativeLevel
		}
	}
	return skill.Level{
		Level:         level,
		RelativeLevel: relativeLevel,
		Tooltip:       tooltip.String(),
	}
}

// CalculateRitualMagicSpellLevel returns the calculated spell level.
func CalculateRitualMagicSpellLevel(entity *Entity, name, powerSource, ritualSkillName string, ritualPrereqCount int, colleges, tags []string, difficulty AttributeDifficulty, points fxp.Int) skill.Level {
	var skillLevel skill.Level
	if len(colleges) == 0 {
		skillLevel = determineRitualMagicSkillLevelForCollege(entity, name, "", ritualSkillName, ritualPrereqCount,
			tags, difficulty, points)
	} else {
		for _, college := range colleges {
			possible := determineRitualMagicSkillLevelForCollege(entity, name, college, ritualSkillName,
				ritualPrereqCount, tags, difficulty, points)
			if skillLevel.Level < possible.Level {
				skillLevel = possible
			}
		}
	}
	if entity != nil {
		tooltip := &xio.ByteBuffer{}
		tooltip.WriteString(skillLevel.Tooltip)
		levels := entity.BestCollegeSpellBonus(tags, colleges, tooltip)
		levels += entity.SpellBonusesFor(feature.SpellPowerSourceID, powerSource, tags, tooltip)
		levels += entity.SpellBonusesFor(feature.SpellNameID, name, tags, tooltip)
		levels = levels.Trunc()
		skillLevel.Level += levels
		skillLevel.RelativeLevel += levels
		skillLevel.Tooltip = tooltip.String()
	}
	return skillLevel
}

func determineRitualMagicSkillLevelForCollege(entity *Entity, name, college, ritualSkillName string, ritualPrereqCount int, tags []string, difficulty AttributeDifficulty, points fxp.Int) skill.Level {
	def := &SkillDefault{
		DefaultType:    gid.Skill,
		Name:           ritualSkillName,
		Specialization: college,
		Modifier:       fxp.From(-ritualPrereqCount),
	}
	if college == "" {
		def.Name = ""
	}
	var limit fxp.Int
	skillLevel := CalculateTechniqueLevel(entity, name, college, tags, def, difficulty.Difficulty, points, false, &limit)
	// CalculateTechniqueLevel() does not add the default skill modifier to the relative level, only to the final level
	skillLevel.RelativeLevel += def.Modifier
	def.Specialization = ""
	def.Modifier -= fxp.Six
	fallback := CalculateTechniqueLevel(entity, name, college, tags, def, difficulty.Difficulty, points, false, &limit)
	fallback.RelativeLevel += def.Modifier
	if skillLevel.Level >= fallback.Level {
		return skillLevel
	}
	return fallback
}

// RitualMagicSatisfied returns true if the Ritual Magic Spell is satisfied.
func (s *Spell) RitualMagicSatisfied(tooltip *xio.ByteBuffer, prefix string) bool {
	if s.Type != gid.RitualMagicSpell {
		return true
	}
	if len(s.College) == 0 {
		if tooltip != nil {
			tooltip.WriteString(prefix)
			tooltip.WriteString(i18n.Text("Must be assigned to a college"))
		}
		return false
	}
	for _, college := range s.College {
		if s.Entity.BestSkillNamed(s.RitualSkillName, college, false, nil) != nil {
			return true
		}
	}
	if s.Entity.BestSkillNamed(s.RitualSkillName, "", false, nil) != nil {
		return true
	}
	if tooltip != nil {
		tooltip.WriteString(prefix)
		tooltip.WriteString(i18n.Text("Requires a skill named "))
		tooltip.WriteString(s.RitualSkillName)
		tooltip.WriteString(" (")
		tooltip.WriteString(s.College[0])
		tooltip.WriteByte(')')
		for _, college := range s.College[1:] {
			tooltip.WriteString(i18n.Text(" or "))
			tooltip.WriteString(s.RitualSkillName)
			tooltip.WriteString(" (")
			tooltip.WriteString(college)
			tooltip.WriteByte(')')
		}
	}
	return false
}

// OwningEntity returns the owning Entity.
func (s *Spell) OwningEntity() *Entity {
	return s.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (s *Spell) SetOwningEntity(entity *Entity) {
	s.Entity = entity
	if s.Container() {
		for _, child := range s.Children {
			child.SetOwningEntity(entity)
		}
	} else {
		for _, w := range s.Weapons {
			w.SetOwner(s)
		}
	}
}

// Notes implements WeaponOwner.
func (s *Spell) Notes() string {
	return s.LocalNotes
}

// Rituals returns the rituals required to cast the spell.
func (s *Spell) Rituals() string {
	if s.Container() || !(s.Entity != nil && s.Entity.Type == datafile.PC && s.Entity.SheetSettings.ShowSpellAdj) {
		return ""
	}
	level := s.CalculateLevel().Level
	switch {
	case level < fxp.Ten:
		return i18n.Text("Ritual: need both hands and feet free and must speak; Time: 2x")
	case level < fxp.Fifteen:
		return i18n.Text("Ritual: speak quietly and make a gesture")
	case level < fxp.Twenty:
		ritual := i18n.Text("Ritual: speak a word or two OR make a small gesture")
		if strings.Contains(strings.ToLower(s.Class), "blocking") {
			return ritual
		}
		return ritual + i18n.Text("; Cost: -1")
	default:
		adj := fxp.As[int]((level - fxp.Fifteen).Div(fxp.Five))
		class := strings.ToLower(s.Class)
		time := ""
		if !strings.Contains(class, "missile") {
			time = fmt.Sprintf(i18n.Text("; Time: x1/%d, rounded up, min 1 sec"), 1<<adj)
		}
		cost := ""
		if !strings.Contains(class, "blocking") {
			cost = fmt.Sprintf(i18n.Text("; Cost: -%d"), adj+1)
		}
		return i18n.Text("Ritual: none") + time + cost
	}
}

// FeatureList returns the list of Features.
func (s *Spell) FeatureList() feature.Features {
	return nil
}

// TagList returns the list of tags.
func (s *Spell) TagList() []string {
	return s.Tags
}

// Description implements WeaponOwner.
func (s *Spell) Description() string {
	return s.String()
}

// SecondaryText returns the less important information that should be displayed with the description.
func (s *Spell) SecondaryText() string {
	var buffer strings.Builder
	prefs := SheetSettingsFor(s.Entity)
	if prefs.NotesDisplay.Inline() {
		text := s.Notes()
		if strings.TrimSpace(text) != "" {
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			buffer.WriteString(text)
		}
	}
	if rituals := s.Rituals(); rituals != "" {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(rituals)
	}
	if prefs.SkillLevelAdjDisplay.Inline() {
		if s.LevelData.Tooltip != "" && s.LevelData.Tooltip != NoAdditionalModifiers {
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			levelTooltip := strings.ReplaceAll(strings.TrimSpace(s.LevelData.Tooltip), "\n", ", ")
			if strings.HasPrefix(levelTooltip, IncludesModifiersFrom+",") {
				levelTooltip = IncludesModifiersFrom + ":" + levelTooltip[len(IncludesModifiersFrom)+1:]
			}
			buffer.WriteString(levelTooltip)
		}
	}
	return buffer.String()
}

func (s *Spell) String() string {
	var buffer strings.Builder
	buffer.WriteString(s.Name)
	if !s.Container() {
		if s.TechLevel != nil {
			buffer.WriteString("/TL")
			buffer.WriteString(*s.TechLevel)
		}
	}
	return buffer.String()
}

// AdjustedPoints returns the points, adjusted for any bonuses.
func (s *Spell) AdjustedPoints(tooltip *xio.ByteBuffer) fxp.Int {
	if s.Container() {
		var total fxp.Int
		for _, one := range s.Children {
			total += one.AdjustedPoints(tooltip)
		}
		return total
	}
	return AdjustedPointsForNonContainerSpell(s.Entity, s.Points, s.Name, s.PowerSource, s.College, s.Tags, tooltip)
}

// AdjustedPointsForNonContainerSpell returns the points, adjusted for any bonuses.
func AdjustedPointsForNonContainerSpell(entity *Entity, points fxp.Int, name, powerSource string, colleges, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	if entity != nil && entity.Type == datafile.PC {
		points += bestCollegeSpellPointBonus(entity, colleges, tags, tooltip)
		points += entity.SpellPointBonusesFor(feature.SpellPowerSourcePointsID, powerSource, tags, tooltip)
		points += entity.SpellPointBonusesFor(feature.SpellPointsID, name, tags, tooltip)
		points = points.Max(0)
	}
	return points
}

func bestCollegeSpellPointBonus(entity *Entity, colleges, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	best := fxp.Min
	bestTooltip := ""
	for _, college := range colleges {
		var buffer *xio.ByteBuffer
		if tooltip != nil {
			buffer = &xio.ByteBuffer{}
		}
		points := entity.SpellPointBonusesFor(feature.SpellCollegePointsID, college, tags, buffer)
		if best < points {
			best = points
			if buffer != nil {
				bestTooltip = buffer.String()
			}
		}
	}
	if tooltip != nil {
		tooltip.WriteString(bestTooltip)
	}
	if best == fxp.Min {
		best = 0
	}
	return best
}
