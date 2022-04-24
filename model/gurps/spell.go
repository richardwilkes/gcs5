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
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
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

// SpellItem holds the Spell data that only exists in non-containers.
type SpellItem struct {
	TechLevel         *string             `json:"tech_level,omitempty"`
	Difficulty        AttributeDifficulty `json:"difficulty"`
	College           CollegeList         `json:"college,omitempty"`
	PowerSource       string              `json:"power_source,omitempty"`
	Class             string              `json:"spell_class,omitempty"`
	Resist            string              `json:"resist,omitempty"`
	CastingCost       string              `json:"casting_cost,omitempty"`
	MaintenanceCost   string              `json:"maintenance_cost,omitempty"`
	CastingTime       string              `json:"casting_time,omitempty"`
	Duration          string              `json:"duration,omitempty"`
	RitualSkillName   string              `json:"base_skill,omitempty"`
	RitualPrereqCount int                 `json:"prereq_count,omitempty"`
	Points            f64d4.Int           `json:"points,omitempty"`
	Prereq            *PrereqList         `json:"prereqs,omitempty"`
	Weapons           []*Weapon           `json:"weapons,omitempty"`
}

// SpellContainer holds the Spell data that only exists in containers.
type SpellContainer struct {
	Children []*Spell `json:"children,omitempty"`
	Open     bool     `json:"open,omitempty"`
}

// SpellData holds the Spell data that is written to disk.
type SpellData struct {
	Type            string    `json:"type"`
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name,omitempty"`
	PageRef         string    `json:"reference,omitempty"`
	LocalNotes      string    `json:"notes,omitempty"`
	VTTNotes        string    `json:"vtt_notes,omitempty"`
	Tags            []string  `json:"categories,omitempty"` // TODO: use tags key instead
	*SpellItem      `json:",omitempty"`
	*SpellContainer `json:",omitempty"`
}

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
	s := Spell{
		SpellData: SpellData{
			Type: gid.Spell,
			ID:   id.NewUUID(),
			Name: i18n.Text("Spell"),
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		s.Type += commonContainerKeyPostfix
		s.SpellContainer = &SpellContainer{Open: true}
	} else {
		s.SpellItem = &SpellItem{
			Difficulty: AttributeDifficulty{
				Attribute:  AttributeIDFor(entity, gid.Intelligence),
				Difficulty: skill.Hard,
			},
			PowerSource: i18n.Text("Arcane"),
			Class:       i18n.Text("Regular"),
			CastingCost: "1",
			CastingTime: "1 sec",
			Duration:    "Instant",
			Points:      f64d4.One,
			Prereq:      NewPrereqList(),
		}
	}
	s.UpdateLevel()
	return &s
}

// NewRitualMagicSpell creates a new Ritual Magic Spell.
func NewRitualMagicSpell(entity *Entity, parent *Spell) *Spell {
	s := NewSpell(entity, parent, false)
	s.Type = gid.RitualMagicSpell
	s.RitualSkillName = "Ritual Magic"
	s.Points = 0
	s.UpdateLevel()
	return s
}

// MarshalJSON implements json.Marshaler.
func (s *Spell) MarshalJSON() ([]byte, error) {
	if s.Container() {
		s.SpellItem = nil
	} else {
		s.SpellContainer = nil
		if s.LevelData.Level > 0 {
			type calc struct {
				Level              f64d4.Int `json:"level"`
				RelativeSkillLevel string    `json:"rsl"`
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
	}
	return json.Marshal(&s.SpellData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Spell) UnmarshalJSON(data []byte) error {
	s.SpellData = SpellData{}
	if err := json.Unmarshal(data, &s.SpellData); err != nil {
		return err
	}
	if s.Container() {
		if s.SpellContainer == nil {
			s.SpellContainer = &SpellContainer{}
		}
		for _, one := range s.Children {
			one.Parent = s
		}
	} else {
		if s.SpellItem == nil {
			s.SpellItem = &SpellItem{}
		}
		if s.Prereq == nil {
			s.Prereq = NewPrereqList()
		}
	}
	return nil
}

// UUID returns the UUID of this data.
func (s *Spell) UUID() uuid.UUID {
	return s.ID
}

// Kind returns the kind of data.
func (s *Spell) Kind() string {
	if s.Container() {
		return i18n.Text("Spell Container")
	}
	return i18n.Text("Spell")
}

// Container returns true if this is a container.
func (s *Spell) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}

// Open returns true if this node is currently open.
func (s *Spell) Open() bool {
	if s.Container() {
		return s.SpellContainer.Open
	}
	return false
}

// SetOpen sets the current open state for this node.
func (s *Spell) SetOpen(open bool) {
	if s.Container() {
		s.SpellContainer.Open = open
	}
}

// NodeChildren returns the children of this node, if any.
func (s *Spell) NodeChildren() []node.Node {
	if s.Container() {
		children := make([]node.Node, len(s.Children))
		for i, child := range s.Children {
			children[i] = child
		}
		return children
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
			data.Primary = s.LevelAsString()
			data.Alignment = unison.EndAlignment
		}
	case SpellRelativeLevelColumn:
		if !s.Container() {
			data.Type = node.Text
			rsl := s.AdjustedRelativeLevel()
			if rsl == f64d4.Min {
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
		data.Primary = s.AdjustedPoints().String()
		data.Alignment = unison.EndAlignment
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
	case rsl == math.MinInt:
		return "-"
	case s.Type != gid.RitualMagicSpell:
		return ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
	default:
		return rsl.StringWithSign()
	}
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Spell) AdjustedRelativeLevel() f64d4.Int {
	if s.Container() {
		return f64d4.Min
	}
	if s.Entity != nil && s.Level() > 0 {
		return s.LevelData.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return f64d4.Min
}

// UpdateLevel updates the level of the spell, returning true if it has changed.
func (s *Spell) UpdateLevel() bool {
	saved := s.LevelData
	if s.Type == gid.RitualMagicSpell {
		s.LevelData = s.calculateRitualMagicLevel()
	} else {
		s.LevelData = s.calculateLevel()
	}
	return saved != s.LevelData
}

// LevelAsString returns the level as a string.
func (s *Spell) LevelAsString() string {
	if s.Container() {
		return ""
	}
	level := s.Level().Trunc()
	if level <= 0 {
		return "-"
	}
	return level.String()
}

// Level returns the computed level without updating it.
func (s *Spell) Level() f64d4.Int {
	if s.Type == gid.RitualMagicSpell {
		return s.calculateRitualMagicLevel().Level
	}
	return s.calculateLevel().Level
}

// IncrementSkillLevel adds enough points to increment the skill level to the next level.
func (s *Spell) IncrementSkillLevel() {
	if !s.Container() {
		basePoints := s.Points.Trunc() + f64d4.One
		maxPoints := basePoints
		if s.Difficulty.Difficulty == skill.Wildcard {
			maxPoints += fxp.Twelve
		} else {
			maxPoints += fxp.Four
		}
		oldLevel := s.Level()
		for points := basePoints; points < maxPoints; points += f64d4.One {
			s.Points = points
			if s.Level() > oldLevel {
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
		oldLevel := s.Level()
		for points := basePoints; points >= minPoints; points -= f64d4.One {
			s.Points = points
			if s.Level() < oldLevel {
				break
			}
		}
		if s.Points > 0 {
			oldLevel = s.Level()
			for s.Points > 0 {
				s.Points -= f64d4.One
				if s.Level() != oldLevel {
					s.Points += f64d4.One
					break
				}
			}
		}
	}
}

func (s *Spell) calculateLevel() skill.Level {
	var tooltip xio.ByteBuffer
	relativeLevel := s.Difficulty.Difficulty.BaseRelativeLevel()
	level := fxp.NegOne
	if s.Entity != nil {
		pts := s.Points
		level = s.Entity.ResolveAttributeCurrent(s.Difficulty.Attribute)
		if s.Difficulty.Difficulty == skill.Wildcard {
			pts = pts.Div(fxp.Three).Trunc()
		}
		switch {
		case pts < f64d4.One:
			level = fxp.NegOne
			relativeLevel = 0
		case pts == f64d4.One:
		// relativeLevel is preset to this point value
		case pts < fxp.Four:
			relativeLevel += f64d4.One
		default:
			relativeLevel += f64d4.One + pts.Div(fxp.Four).Trunc()
		}
		if level != f64d4.One {
			relativeLevel += s.BestCollegeSpellBonus(s.Tags, s.College, &tooltip)
			relativeLevel += s.SpellBonusesFor(feature.SpellPowerSourceID, s.PowerSource, s.Tags, &tooltip)
			relativeLevel += s.SpellBonusesFor(feature.SpellNameID, s.Name, s.Tags, &tooltip)
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

func (s *Spell) calculateRitualMagicLevel() skill.Level {
	var skillLevel skill.Level
	if len(s.College) == 0 {
		skillLevel = s.determineSkillLevelForCollege("")
	} else {
		for _, college := range s.College {
			possible := s.determineSkillLevelForCollege(college)
			if skillLevel.Level < possible.Level {
				skillLevel = possible
			}
		}
	}
	if s.Entity != nil {
		tooltip := &xio.ByteBuffer{}
		tooltip.WriteString(skillLevel.Tooltip)
		levels := s.BestCollegeSpellBonus(s.Tags, s.College, tooltip)
		levels += s.SpellBonusesFor(feature.SpellPowerSourceID, s.PowerSource, s.Tags, tooltip)
		levels += s.SpellBonusesFor(feature.SpellNameID, s.Name, s.Tags, tooltip)
		levels = levels.Trunc()
		skillLevel.Level += levels
		skillLevel.RelativeLevel += levels
		skillLevel.Tooltip = tooltip.String()
	}
	return skillLevel
}

func (s *Spell) determineSkillLevelForCollege(college string) skill.Level {
	def := &SkillDefault{
		DefaultType:    gid.Skill,
		Name:           s.RitualSkillName,
		Specialization: college,
		Modifier:       f64d4.FromInt(-s.RitualPrereqCount),
	}
	if college == "" {
		def.Name = ""
	}
	var limit f64d4.Int
	skillLevel := CalculateTechniqueLevel(s.Entity, s.Name, college, s.Tags, def, s.Difficulty.Difficulty, s.AdjustedPoints(), false, &limit)
	// CalculateTechniqueLevel() does not add the default skill modifier to the relative level, only to the final level
	skillLevel.RelativeLevel += def.Modifier
	def.Specialization = ""
	def.Modifier -= fxp.Six
	fallback := CalculateTechniqueLevel(s.Entity, s.Name, college, s.Tags, def, s.Difficulty.Difficulty, s.AdjustedPoints(), false, &limit)
	fallback.RelativeLevel += def.Modifier
	if skillLevel.Level >= fallback.Level {
		return skillLevel
	}
	return fallback
}

// BestCollegeSpellBonus returns the best college spell bonus for this spell.
func (s *Spell) BestCollegeSpellBonus(tags, colleges []string, tooltip *xio.ByteBuffer) f64d4.Int {
	best := f64d4.Min
	var bestTooltip string
	for _, college := range colleges {
		var buffer *xio.ByteBuffer
		if tooltip != nil {
			buffer = &xio.ByteBuffer{}
		}
		if pts := s.SpellBonusesFor(feature.SpellCollegeID, college, tags, buffer); best < pts {
			best = pts
			if buffer != nil {
				bestTooltip = buffer.String()
			}
		}
	}
	if tooltip != nil {
		tooltip.WriteString(bestTooltip)
	}
	if best == f64d4.Min {
		best = 0
	}
	return best
}

// SpellBonusesFor returns the bonus for this spell.
func (s *Spell) SpellBonusesFor(featureID, qualifier string, tags []string, tooltip *xio.ByteBuffer) f64d4.Int {
	level := s.Entity.BonusFor(featureID, tooltip)
	level += s.Entity.BonusFor(featureID+"/"+strings.ToLower(qualifier), tooltip)
	level += s.Entity.SpellComparedBonusFor(featureID+"*", qualifier, tags, tooltip)
	return level
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
	level := s.Level()
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
		adj := (level - fxp.Fifteen).Div(fxp.Five).AsInt()
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
func (s *Spell) AdjustedPoints() f64d4.Int {
	if s.Container() {
		var total f64d4.Int
		for _, one := range s.Children {
			total += one.AdjustedPoints()
		}
		return total
	}
	points := s.Points
	if s.Entity != nil && s.Entity.Type == datafile.PC {
		points += s.bestCollegeSpellPointBonus(nil)
		points += s.Entity.SpellPointBonusesFor(feature.SpellPowerSourcePointsID, s.PowerSource, s.Tags, nil)
		points += s.Entity.SpellPointBonusesFor(feature.SpellPointsID, s.Name, s.Tags, nil)
		points = points.Max(0)
	}
	return points
}

func (s *Spell) bestCollegeSpellPointBonus(tooltip *xio.ByteBuffer) f64d4.Int {
	best := f64d4.Min
	bestTooltip := ""
	for _, college := range s.College {
		var buffer *xio.ByteBuffer
		if tooltip != nil {
			buffer = &xio.ByteBuffer{}
		}
		points := s.Entity.SpellPointBonusesFor(feature.SpellCollegePointsID, college, s.Tags, buffer)
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
	if best == f64d4.Min {
		best = 0
	}
	return best
}
