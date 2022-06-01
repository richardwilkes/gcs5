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
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

var (
	_ Node              = &Skill{}
	_ TechLevelProvider = &Skill{}
)

// Columns that can be used with the skill method .CellData()
const (
	SkillDescriptionColumn = iota
	SkillDifficultyColumn
	SkillTagsColumn
	SkillReferenceColumn
	SkillLevelColumn
	SkillRelativeLevelColumn
	SkillPointsColumn
)

const skillListTypeKey = "skill_list"

// Skill holds the data for a skill.
type Skill struct {
	SkillData
	Entity            *Entity
	Parent            *Skill
	LevelData         skill.Level
	UnsatisfiedReason string
}

type skillListData struct {
	Type    string   `json:"type"`
	Version int      `json:"version"`
	Rows    []*Skill `json:"rows"`
}

// NewSkillsFromFile loads an Skill list from a file.
func NewSkillsFromFile(fileSystem fs.FS, filePath string) ([]*Skill, error) {
	var data skillListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != skillListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveSkills writes the Skill list to the file as JSON.
func SaveSkills(skills []*Skill, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &skillListData{
		Type:    skillListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    skills,
	})
}

// NewSkill creates a new Skill.
func NewSkill(entity *Entity, parent *Skill, container bool) *Skill {
	return newSkill(entity, parent, gid.Skill, container)
}

// NewTechnique creates a new technique (i.e. a specialized use of a Skill). All parameters may be nil or empty.
func NewTechnique(entity *Entity, parent *Skill, skillName string) *Skill {
	if skillName == "" {
		skillName = i18n.Text("Skill")
	}
	s := newSkill(entity, parent, gid.Technique, false)
	s.TechniqueDefault = &SkillDefault{
		DefaultType: gid.Skill,
		Name:        skillName,
	}
	return s
}

func newSkill(entity *Entity, parent *Skill, typeKey string, container bool) *Skill {
	s := Skill{
		SkillData: SkillData{
			ContainerBase: newContainerBase[*Skill](typeKey, container),
		},
		Entity: entity,
		Parent: parent,
	}
	if !container {
		s.Difficulty.Attribute = AttributeIDFor(entity, gid.Dexterity)
		s.Difficulty.Difficulty = skill.Average
		s.Points = fxp.One
	}
	s.Name = s.Kind()
	return &s
}

// MarshalJSON implements json.Marshaler.
func (s *Skill) MarshalJSON() ([]byte, error) {
	s.ClearUnusedFieldsForType()
	if s.Container() || s.LevelData.Level <= 0 {
		return json.Marshal(&s.SkillData)
	}
	type calc struct {
		Level              fxp.Int `json:"level"`
		RelativeSkillLevel string  `json:"rsl"`
	}
	return json.Marshal(&struct {
		SkillData
		Calc calc `json:"calc"`
	}{
		SkillData: s.SkillData,
		Calc: calc{
			Level:              s.LevelData.Level,
			RelativeSkillLevel: s.RelativeLevel(),
		},
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Skill) UnmarshalJSON(data []byte) error {
	var localData struct {
		SkillData
		// Old data fields
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(data, &localData); err != nil {
		return err
	}
	localData.ClearUnusedFieldsForType()
	s.SkillData = localData.SkillData
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
func (s *Skill) CellData(column int, data *CellData) {
	switch column {
	case SkillDescriptionColumn:
		data.Type = Text
		data.Primary = s.Description()
		data.Secondary = s.SecondaryText()
		data.UnsatisfiedReason = s.UnsatisfiedReason
	case SkillDifficultyColumn:
		if !s.Container() {
			data.Type = Text
			data.Primary = s.Difficulty.Description(s.Entity)
		}
	case SkillTagsColumn:
		data.Type = Text
		data.Primary = CombineTags(s.Tags)
	case SkillReferenceColumn, PageRefCellAlias:
		data.Type = PageRef
		data.Primary = s.PageRef
		data.Secondary = s.Name
	case SkillLevelColumn:
		if !s.Container() {
			data.Type = Text
			level := s.CalculateLevel()
			data.Primary = level.LevelAsString(s.Container())
			if level.Tooltip != "" {
				data.Tooltip = IncludesModifiersFrom + ":" + level.Tooltip
			}
			data.Alignment = unison.EndAlignment
		}
	case SkillRelativeLevelColumn:
		if !s.Container() {
			data.Type = Text
			data.Primary = FormatRelativeSkill(s.Entity, s.Type, s.Difficulty, s.AdjustedRelativeLevel())
		}
	case SkillPointsColumn:
		data.Type = Text
		var tooltip xio.ByteBuffer
		data.Primary = s.AdjustedPoints(&tooltip).String()
		data.Alignment = unison.EndAlignment
		if tooltip.Len() != 0 {
			data.Tooltip = IncludesModifiersFrom + ":" + tooltip.String()
		}
	}
}

// FormatRelativeSkill formats the relative skill for display.
func FormatRelativeSkill(entity *Entity, typ string, difficulty AttributeDifficulty, rsl fxp.Int) string {
	switch {
	case rsl == fxp.Min:
		return "-"
	case strings.HasPrefix(typ, gid.Skill) || strings.HasPrefix(typ, gid.Spell):
		s := ResolveAttributeName(entity, difficulty.Attribute)
		if rsl != 0 {
			s += rsl.StringWithSign()
		}
		return s
	default:
		return rsl.StringWithSign()
	}
}

// Depth returns the number of parents this node has.
func (s *Skill) Depth() int {
	count := 0
	p := s.Parent
	for p != nil {
		count++
		p = p.Parent
	}
	return count
}

// OwningEntity returns the owning Entity.
func (s *Skill) OwningEntity() *Entity {
	return s.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (s *Skill) SetOwningEntity(entity *Entity) {
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

// DefaultSkill returns the skill currently defaulted to, or nil.
func (s *Skill) DefaultSkill() *Skill {
	if s.Entity == nil {
		return nil
	}
	if strings.HasPrefix(s.Type, gid.Technique) {
		return s.Entity.BaseSkill(s.TechniqueDefault, true)
	}
	return s.Entity.BaseSkill(s.DefaultedFrom, true)
}

// Notes implements WeaponOwner.
func (s *Skill) Notes() string {
	return s.LocalNotes
}

// ModifierNotes returns the notes due to modifiers.
func (s *Skill) ModifierNotes() string {
	if strings.HasPrefix(s.Type, gid.Technique) {
		return i18n.Text("Default: ") + s.TechniqueDefault.FullName(s.Entity) + s.TechniqueDefault.ModifierAsString()
	}
	if s.Difficulty.Difficulty != skill.Wildcard {
		defSkill := s.DefaultSkill()
		if defSkill != nil && s.DefaultedFrom != nil {
			return i18n.Text("Default: ") + defSkill.String() + s.DefaultedFrom.ModifierAsString()
		}
	}
	return ""
}

// FeatureList returns the list of Features.
func (s *Skill) FeatureList() feature.Features {
	return s.Features
}

// TagList returns the list of tags.
func (s *Skill) TagList() []string {
	return s.Tags
}

// Description implements WeaponOwner.
func (s *Skill) Description() string {
	return s.String()
}

// SecondaryText returns the less important information that should be displayed with the description.
func (s *Skill) SecondaryText() string {
	var buffer strings.Builder
	prefs := SheetSettingsFor(s.Entity)
	if prefs.ModifiersDisplay.Inline() {
		text := s.ModifierNotes()
		if strings.TrimSpace(text) != "" {
			buffer.WriteString(text)
		}
	}
	if prefs.NotesDisplay.Inline() {
		text := s.Notes()
		if strings.TrimSpace(text) != "" {
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			buffer.WriteString(text)
		}
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

func (s *Skill) String() string {
	var buffer strings.Builder
	buffer.WriteString(s.Name)
	if !s.Container() {
		if s.TechLevel != nil {
			buffer.WriteString("/TL")
			buffer.WriteString(*s.TechLevel)
		}
		if s.Specialization != "" {
			buffer.WriteString(" (")
			buffer.WriteString(s.Specialization)
			buffer.WriteByte(')')
		}
	}
	return buffer.String()
}

// RelativeLevel returns the adjusted relative level as a string.
func (s *Skill) RelativeLevel() string {
	if s.Container() || s.LevelData.Level <= 0 {
		return ""
	}
	rsl := s.AdjustedRelativeLevel()
	switch {
	case rsl == fxp.Min:
		return "-"
	case strings.HasPrefix(s.Type, gid.Skill):
		return ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
	default:
		return rsl.StringWithSign()
	}
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Skill) AdjustedRelativeLevel() fxp.Int {
	if s.Container() {
		return fxp.Min
	}
	if s.Entity != nil && s.LevelData.Level > 0 {
		if strings.HasPrefix(s.Type, gid.Technique) {
			return s.LevelData.RelativeLevel + s.TechniqueDefault.Modifier
		}
		return s.LevelData.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fxp.Min
}

// AdjustedPoints returns the points, adjusted for any bonuses.
func (s *Skill) AdjustedPoints(tooltip *xio.ByteBuffer) fxp.Int {
	if s.Container() {
		var total fxp.Int
		for _, one := range s.Children {
			total += one.AdjustedPoints(tooltip)
		}
		return total
	}
	return AdjustedPointsForNonContainerSkillOrTechnique(s.Entity, s.Points, s.Name, s.Specialization, s.Tags, tooltip)
}

// AdjustedPointsForNonContainerSkillOrTechnique returns the points, adjusted for any bonuses.
func AdjustedPointsForNonContainerSkillOrTechnique(entity *Entity, points fxp.Int, name, specialization string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	if entity != nil && entity.Type == datafile.PC {
		points += entity.SkillPointComparedBonusFor(feature.SkillPointsID+"*", name, specialization, tags, tooltip)
		points += entity.BonusFor(feature.SkillPointsID+"/"+strings.ToLower(name), tooltip)
		points = points.Max(0)
	}
	return points
}

// IncrementSkillLevel adds enough points to increment the skill level to the next level.
func (s *Skill) IncrementSkillLevel() {
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
func (s *Skill) DecrementSkillLevel() {
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

// CalculateLevel returns the computed level without updating it.
func (s *Skill) CalculateLevel() skill.Level {
	points := s.AdjustedPoints(nil)
	if strings.HasPrefix(s.Type, gid.Skill) {
		return CalculateSkillLevel(s.Entity, s.Name, s.Specialization, s.Tags, s.DefaultedFrom, s.Difficulty, points,
			s.EncumbrancePenaltyMultiplier)
	}
	return CalculateTechniqueLevel(s.Entity, s.Name, s.Specialization, s.Tags, s.TechniqueDefault,
		s.Difficulty.Difficulty, points, true, s.TechniqueLimitModifier)
}

// CalculateSkillLevel returns the calculated level for a skill.
func CalculateSkillLevel(entity *Entity, name, specialization string, tags []string, def *SkillDefault, difficulty AttributeDifficulty, points, encumbrancePenaltyMultiplier fxp.Int) skill.Level {
	var tooltip xio.ByteBuffer
	relativeLevel := difficulty.Difficulty.BaseRelativeLevel()
	level := entity.ResolveAttributeCurrent(difficulty.Attribute)
	if level != fxp.Min {
		if difficulty.Difficulty == skill.Wildcard {
			points = points.Div(fxp.Three)
		} else if def != nil && def.Points > 0 {
			points += def.Points
		}
		points = points.Trunc()
		switch {
		case points == fxp.One:
			// relativeLevel is preset to this point value
		case points > fxp.One && points < fxp.Four:
			relativeLevel += fxp.One
		case points >= fxp.Four:
			relativeLevel += fxp.One + points.Div(fxp.Four).Trunc()
		case difficulty.Difficulty != skill.Wildcard && def != nil && def.Points < 0:
			relativeLevel = def.AdjLevel - level
		default:
			level = fxp.Min
			relativeLevel = 0
		}
		if level != fxp.Min {
			level += relativeLevel
			if difficulty.Difficulty != skill.Wildcard && def != nil && level < def.AdjLevel {
				level = def.AdjLevel
			}
			if entity != nil {
				bonus := entity.SkillComparedBonusFor(feature.SkillNameID+"*", name, specialization, tags, &tooltip)
				level += bonus
				relativeLevel += bonus
				bonus = entity.BonusFor(feature.SkillNameID+"/"+strings.ToLower(name), &tooltip)
				level += bonus
				relativeLevel += bonus
				bonus = entity.EncumbranceLevel(true).Penalty().Mul(encumbrancePenaltyMultiplier)
				level += bonus
				if bonus != 0 {
					fmt.Fprintf(&tooltip, i18n.Text("\nEncumbrance [%s]"), bonus.StringWithSign())
				}
			}
		}
	}
	return skill.Level{
		Level:         level,
		RelativeLevel: relativeLevel,
		Tooltip:       tooltip.String(),
	}
}

// CalculateTechniqueLevel returns the calculated level for a technique.
func CalculateTechniqueLevel(entity *Entity, name, specialization string, tags []string, def *SkillDefault, difficulty skill.Difficulty, points fxp.Int, requirePoints bool, limitModifier *fxp.Int) skill.Level {
	var tooltip xio.ByteBuffer
	var relativeLevel fxp.Int
	level := fxp.Min
	if entity != nil {
		if def.DefaultType == gid.Skill {
			if sk := entity.BaseSkill(def, requirePoints); sk != nil {
				level = sk.CalculateLevel().Level
			}
		} else {
			// Take the modifier back out, as we wanted the base, not the final value.
			level = def.SkillLevelFast(entity, true, nil, false) - def.Modifier
		}
		if level != fxp.Min {
			baseLevel := level
			level += def.Modifier
			if difficulty == skill.Hard {
				points -= fxp.One
			}
			if points > 0 {
				relativeLevel = points
			}
			if level != fxp.Min {
				relativeLevel += entity.BonusFor(feature.SkillNameID+"/"+strings.ToLower(name), &tooltip)
				relativeLevel += entity.SkillComparedBonusFor(feature.SkillNameID+"*", name, specialization, tags, &tooltip)
				level += relativeLevel
			}
			if limitModifier != nil {
				if max := baseLevel + *limitModifier; level > max {
					relativeLevel -= level - max
					level = max
				}
			}
		}
	}
	return skill.Level{
		Level:         level,
		RelativeLevel: relativeLevel,
		Tooltip:       tooltip.String(),
	}
}

// UpdateLevel updates the level of the skill, returning true if it has changed.
func (s *Skill) UpdateLevel() bool {
	saved := s.LevelData
	s.DefaultedFrom = s.bestDefaultWithPoints(nil)
	s.LevelData = s.CalculateLevel()
	return saved != s.LevelData
}

func (s *Skill) bestDefaultWithPoints(excluded *SkillDefault) *SkillDefault {
	if strings.HasPrefix(s.Type, gid.Technique) {
		return nil
	}
	best := s.bestDefault(excluded)
	if best != nil {
		baseLine := (s.Entity.ResolveAttributeCurrent(s.Difficulty.Attribute) + s.Difficulty.Difficulty.BaseRelativeLevel()).Trunc()
		level := best.Level.Trunc()
		best.AdjLevel = level
		switch {
		case level == baseLine:
			best.Points = fxp.One
		case level == baseLine+fxp.One:
			best.Points = fxp.Two
		case level > baseLine+fxp.One:
			best.Points = fxp.Four.Mul(level - (baseLine + fxp.One))
		default:
			best.Points = -level.Max(0)
		}
	}
	return best
}

func (s *Skill) bestDefault(excluded *SkillDefault) *SkillDefault {
	if s.Entity == nil || len(s.Defaults) == 0 {
		return nil
	}
	excludes := make(map[string]bool)
	excludes[s.String()] = true
	var bestDef *SkillDefault
	best := fxp.Min
	for _, def := range s.Defaults {
		// For skill-based defaults, prune out any that already use a default that we are involved with
		if def.Equivalent(excluded) || s.inDefaultChain(def, make(map[*Skill]bool)) {
			continue
		}
		level := def.SkillLevel(s.Entity, true, excludes, strings.HasPrefix(s.Type, gid.Skill))
		if def.SkillBased() {
			if other := s.Entity.BestSkillNamed(def.Name, def.Specialization, true, excludes); other != nil {
				level -= s.Entity.SkillComparedBonusFor(feature.SkillNameID+"*", def.Name, def.Specialization, s.Tags, nil)
				level -= s.Entity.BonusFor(feature.SkillNameID+"/"+strings.ToLower(def.Name), nil)
			}
		}
		if best < level {
			best = level
			bestDef = def.CloneWithoutLevelOrPoints()
			bestDef.Level = level
		}
	}
	return bestDef
}

func (s *Skill) inDefaultChain(def *SkillDefault, lookedAt map[*Skill]bool) bool {
	if s.Entity == nil || def == nil || !def.SkillBased() {
		return false
	}
	hadOne := false
	for _, one := range s.Entity.SkillNamed(def.Name, def.Specialization, true, nil) {
		if one == s {
			return true
		}
		if _, has := lookedAt[one]; !has {
			lookedAt[one] = true
			if s.inDefaultChain(one.DefaultedFrom, lookedAt) {
				return true
			}
		}
		hadOne = true
	}
	return !hadOne
}

// TechniqueSatisfied returns true if the Technique is satisfied.
func (s *Skill) TechniqueSatisfied(tooltip *xio.ByteBuffer, prefix string) bool {
	if strings.HasPrefix(s.Type, gid.Skill) || !s.TechniqueDefault.SkillBased() {
		return true
	}
	sk := s.Entity.BestSkillNamed(s.TechniqueDefault.Name, s.TechniqueDefault.Specialization, false, nil)
	satisfied := sk != nil && (strings.HasPrefix(sk.Type, gid.Technique) || sk.Points > 0)
	if !satisfied && tooltip != nil {
		tooltip.WriteString(prefix)
		if sk == nil {
			tooltip.WriteString(i18n.Text("Requires a skill named "))
		} else {
			tooltip.WriteString(i18n.Text("Requires at least 1 point in the skill named "))
		}
		tooltip.WriteString(s.TechniqueDefault.FullName(s.Entity))
	}
	return satisfied
}

// TL implements TechLevelProvider.
func (s *Skill) TL() string {
	if s.TechLevel != nil {
		return *s.TechLevel
	}
	return ""
}

// SetTL implements TechLevelProvider.
func (s *Skill) SetTL(tl string) {
	if s.TechLevel != nil {
		*s.TechLevel = tl
	}
}
