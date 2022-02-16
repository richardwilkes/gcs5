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
	spellDescriptionColumn = iota
	spellResistColumn
	spellClassColumn
	spellCollegeColumn
	spellCastCostColumn
	spellMaintainCostColumn
	spellCastTimeColumn
	spellDurationColumn
	spellDifficultyColumn
	spellCategoryColumn
	spellReferenceColumn
	spellColumnCount
)

var _ unison.TableRowData = &SpellNode{}

// SpellNode holds a spell in the spell list.
type SpellNode struct {
	dockable  *SpellListDockable
	spell     *gurps.Spell
	children  []unison.TableRowData
	cellCache []*cellCache
}

// NewSpellNode creates a new SpellNode.
func NewSpellNode(dockable *SpellListDockable, spell *gurps.Spell) *SpellNode {
	n := &SpellNode{
		dockable:  dockable,
		spell:     spell,
		cellCache: make([]*cellCache, spellColumnCount),
	}
	return n
}

// CanHaveChildRows always returns true.
func (n *SpellNode) CanHaveChildRows() bool {
	return n.spell.Container()
}

// ChildRows returns the children of this node.
func (n *SpellNode) ChildRows() []unison.TableRowData {
	if n.spell.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.spell.Children))
		for i, one := range n.spell.Children {
			n.children[i] = NewSpellNode(n.dockable, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *SpellNode) CellDataForSort(index int) string {
	switch index {
	case spellDescriptionColumn:
		text := n.spell.Description()
		secondary := n.spell.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case spellResistColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.Resist
	case spellClassColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.Class
	case spellCollegeColumn:
		if n.spell.Container() {
			return ""
		}
		return strings.Join(n.spell.College, ", ")
	case spellCastCostColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.CastingCost
	case spellMaintainCostColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.MaintenanceCost
	case spellCastTimeColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.CastingTime
	case spellDurationColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.Duration
	case spellDifficultyColumn:
		if n.spell.Container() {
			return ""
		}
		return n.spell.Difficulty.Description(n.spell.Entity)
	case spellCategoryColumn:
		return strings.Join(n.spell.Categories, ", ")
	case spellReferenceColumn:
		return n.spell.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *SpellNode) ColumnCell(row, col int, selected bool) unison.Paneler {
	width := n.dockable.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].matches(width, data) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].panel.Children() {
			child.Self.(*unison.Label).LabelTheme.OnBackgroundInk = color
		}
		return n.cellCache[col].panel
	}
	p := &unison.Panel{}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	if col == spellDescriptionColumn {
		createAndAddCellLabel(p, width, n.spell.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.spell.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			createAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	} else {
		createAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &cellCache{
		width: width,
		data:  data,
		panel: p,
	}
	return p
}

// IsOpen returns true if this node should display its children.
func (n *SpellNode) IsOpen() bool {
	return n.spell.Container() && n.spell.Open
}

// SetOpen sets the current open state for this node.
func (n *SpellNode) SetOpen(open bool) {
	if n.spell.Container() && open != n.spell.Open {
		n.spell.Open = open
		n.dockable.table.SyncToModel()
		n.dockable.table.SizeColumnsToFit(true)
	}
}
