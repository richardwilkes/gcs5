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

package sheet

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// PageList holds a list for a sheet page.
type PageList struct {
	unison.Panel
	tableHeader *unison.TableHeader
	table       *unison.Table
}

// NewPageList creates a new list for a sheet page.
func NewPageList(entity *gurps.Entity, columnHeaders []unison.TableColumnHeader) *PageList {
	p := &PageList{
		table: unison.NewTable(),
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.StartAlignment,
		HGrab:  true,
	})
	p.table.Padding.Top = 0
	p.table.Padding.Bottom = 0
	p.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range p.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom32.Size{})
		p.table.ColumnSizes[i].AutoMinimum = pref.Width
		p.table.ColumnSizes[i].AutoMaximum = 800
		p.table.ColumnSizes[i].Minimum = pref.Width
		p.table.ColumnSizes[i].Maximum = 10000
	}
	p.table.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.tableHeader = unison.NewTableHeader(p.table, columnHeaders...)
	p.tableHeader.BackgroundInk = theme.HeaderColor
	p.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fixed.F64d4FromString(s1); err == nil {
			var n2 fixed.F64d4
			if n2, err = fixed.F64d4FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	p.tableHeader.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.table.SizeColumnsToFit(true)
	p.AddChild(p.tableHeader)
	p.AddChild(p.table)
	return p
}

// NewPageListHeader creates a new page list header.
func NewPageListHeader(title, tooltip string) *unison.DefaultTableColumnHeader {
	header := unison.NewTableColumnHeader(title)
	header.Font = theme.PageLabelPrimaryFont
	header.OnBackgroundInk = theme.OnHeaderColor
	if tooltip != "" {
		header.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return header
}

// NewPageReferenceHeader creates a new page reference page header.
func NewPageReferenceHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`A reference to the book and page the item appears on
e.g. B22 would refer to "Basic Set", page 22`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  res.BookmarkSVG,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

// NewEquippedHeader creates a new equipped page header.
func NewEquippedHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`Whether this piece of equipment is equipped or just carried. Items that are not equipped do not apply any features they may normally contribute to the character.`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  res.CircledCheckSVG,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

// NewMoneyHeader creates a new money page header.
func NewMoneyHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`The value of one of these pieces of equipment`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  res.CoinsSVG,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

// NewExtendedMoneyHeader creates a new extended money page header.
func NewExtendedMoneyHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`The value of all of these pieces of equipment, plus the value of any contained equipment`))
	baseline := header.Font.Baseline()
	header.Drawable = &widget.DrawableSVGPair{
		Left:  res.StackSVG,
		Right: res.CoinsSVG,
		Size:  geom32.NewSize(baseline*2+4, baseline),
	}
	return header
}

// NewWeightHeader creates a new weight page header.
func NewWeightHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`The weight of one of these pieces of equipment`))
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  res.WeightSVG,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

// NewExtendedWeightHeader creates a new extended weight page header.
func NewExtendedWeightHeader() unison.TableColumnHeader {
	header := NewPageListHeader("", i18n.Text(`The weight of all of these pieces of equipment, plus the weight of any contained equipment`))
	baseline := header.Font.Baseline()
	header.Drawable = &widget.DrawableSVGPair{
		Left:  res.StackSVG,
		Right: res.WeightSVG,
		Size:  geom32.NewSize(baseline*2+4, baseline),
	}
	return header
}
