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
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
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

var (
	_ unison.TableRowData = &SpellNode{}
	_ tbl.Matcher         = &SpellNode{}
)

// SpellNode holds a spell in the spell list.
type SpellNode struct {
	table     *unison.Table
	parent    *SpellNode
	spell     *gurps.Spell
	children  []unison.TableRowData
	cellCache []*tbl.CellCache
}

// NewSpellListDockable creates a new unison.Dockable for spell list files.
func NewSpellListDockable(filePath string) (unison.Dockable, error) {
	spells, err := gurps.NewSpellsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewTableDockable(filePath, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Spell"), "", false),
		tbl.NewHeader(i18n.Text("Resist"), i18n.Text("Resistance"), false),
		tbl.NewHeader(i18n.Text("Class"), "", false),
		tbl.NewHeader(i18n.Text("College"), "", false),
		tbl.NewHeader(i18n.Text("Diff"), i18n.Text("Difficulty"), false),
		tbl.NewHeader(i18n.Text("Maintain"), i18n.Text("The mana cost to maintain the spell"), false),
		tbl.NewHeader(i18n.Text("Time"), i18n.Text("The time required to cast the spell"), false),
		tbl.NewHeader(i18n.Text("Duration"), "", false),
		tbl.NewHeader(i18n.Text("Cost"), i18n.Text("The mana cost to cast the spell"), false),
		tbl.NewHeader(i18n.Text("Category"), "", false),
		tbl.NewPageRefHeader(false),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(spells))
		for _, one := range spells {
			rows = append(rows, NewSpellNode(table, nil, one))
		}
		return rows
	}), nil
}

// NewSpellNode creates a new SpellNode.
func NewSpellNode(table *unison.Table, parent *SpellNode, spell *gurps.Spell) *SpellNode {
	n := &SpellNode{
		table:     table,
		parent:    parent,
		spell:     spell,
		cellCache: make([]*tbl.CellCache, spellColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *SpellNode) ParentRow() unison.TableRowData {
	return n.parent
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
			n.children[i] = NewSpellNode(n.table, n, one)
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
	width := n.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].Matches(width, data) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].Panel.Children() {
			child.Self.(*unison.Label).LabelTheme.OnBackgroundInk = color
		}
		return n.cellCache[col].Panel
	}
	p := &unison.Panel{}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	switch col {
	case spellDescriptionColumn:
		tbl.CreateAndAddCellLabel(p, width, n.spell.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.spell.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			tbl.CreateAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case spellReferenceColumn:
		tbl.CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.spell.Name, unison.DefaultLabelTheme.Font, selected)
	default:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &tbl.CellCache{
		Width: width,
		Data:  data,
		Panel: p,
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
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *SpellNode) Match(text string) bool {
	return strings.Contains(strings.ToLower(n.spell.Description()), text) ||
		strings.Contains(strings.ToLower(n.spell.SecondaryText()), text) ||
		(!n.spell.Container() &&
			(strings.Contains(strings.ToLower(n.spell.Resist), text) ||
				strings.Contains(strings.ToLower(n.spell.Class), text) ||
				strings.Contains(strings.ToLower(n.spell.CastingCost), text) ||
				strings.Contains(strings.ToLower(n.spell.MaintenanceCost), text) ||
				strings.Contains(strings.ToLower(n.spell.CastingTime), text) ||
				strings.Contains(strings.ToLower(n.spell.Duration), text) ||
				strings.Contains(strings.ToLower(n.spell.Difficulty.Description(n.spell.Entity)), text) ||
				stringSliceContains(n.spell.College, text))) ||
		strings.Contains(strings.ToLower(n.spell.PageRef), text) ||
		stringSliceContains(n.spell.Categories, text)
}
