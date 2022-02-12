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

package workspace

import (
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/unison"
)

const (
	skillDescriptionColumn = iota
	skillDifficultyColumn
	skillCategoryColumn
	skillReferenceColumn
	skillColumnCount
)

var _ unison.TableRowData = &SkillNode{}

// SkillNode holds a skill in the skill list.
type SkillNode struct {
	dockable *SkillListDockable
	skill    *gurps.Skill
	children []unison.TableRowData
}

// NewSkillNode creates a new SkillNode.
func NewSkillNode(dockable *SkillListDockable, skill *gurps.Skill) *SkillNode {
	n := &SkillNode{
		dockable: dockable,
		skill:    skill,
	}
	return n
}

// CanHaveChildRows always returns true.
func (n *SkillNode) CanHaveChildRows() bool {
	return n.skill.Container()
}

// ChildRows returns the children of this node.
func (n *SkillNode) ChildRows() []unison.TableRowData {
	if n.skill.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.children))
		for i, one := range n.skill.Children {
			n.children[i] = NewSkillNode(n.dockable, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *SkillNode) CellDataForSort(index int) string {
	switch index {
	case skillDescriptionColumn:
		text := n.skill.Description()
		secondary := n.skill.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case skillDifficultyColumn:
		if n.skill.Container() {
			return ""
		}
		return n.skill.Difficulty.String()
	case skillCategoryColumn:
		return strings.Join(n.skill.Categories, ", ")
	case skillReferenceColumn:
		return n.skill.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *SkillNode) ColumnCell(index int, selected bool) unison.Paneler {
	if index == skillDescriptionColumn {
		p := &unison.Panel{}
		p.Self = p
		p.SetLayout(&unison.FlexLayout{Columns: 1})
		p.AddChild(createCellLabel(n.skill.Description(), selected).AsPanel())
		if text := n.skill.SecondaryText(); strings.TrimSpace(text) != "" {
			secondary := createCellLabel(text, selected)
			desc := secondary.Font.Descriptor()
			desc.Size--
			secondary.Font = desc.Font()
			p.AddChild(secondary.AsPanel())
		}
		return p
	}
	return createCellLabel(n.CellDataForSort(index), selected).AsPanel()
}

func createCellLabel(text string, selected bool) *unison.Label {
	label := unison.NewLabel()
	label.Text = text
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	return label
}

// IsOpen returns true if this node should display its children.
func (n *SkillNode) IsOpen() bool {
	return n.skill.Container() && n.skill.Open
}

// SetOpen sets the current open state for this node.
func (n *SkillNode) SetOpen(open bool) {
	if n.skill.Container() && open != n.skill.Open {
		n.skill.Open = open
		n.dockable.table.SyncToModel()
		n.dockable.table.SizeColumnsToFit(true)
	}
}
