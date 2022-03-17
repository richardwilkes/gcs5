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

package tbl

import (
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// NewHeader creates a new list header.
func NewHeader(title, tooltip string, forPage bool) *unison.DefaultTableColumnHeader {
	header := unison.NewTableColumnHeader(title)
	if forPage {
		header.Font = theme.PageLabelPrimaryFont
		header.OnBackgroundInk = theme.OnHeaderColor
	}
	if tooltip != "" {
		header.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return header
}

// NewSVGHeader creates a new list header with an SVG image as its content rather than text.
func NewSVGHeader(svg *unison.SVG, tooltip string, forPage bool) *unison.DefaultTableColumnHeader {
	header := NewHeader("", tooltip, forPage)
	baseline := header.Font.Baseline()
	header.Drawable = &unison.DrawableSVG{
		SVG:  svg,
		Size: geom32.NewSize(baseline, baseline),
	}
	return header
}

// NewSVGPairHeader creates a new list header with a pair of SVG images as its content rather than text.
func NewSVGPairHeader(leftSVG, rightSVG *unison.SVG, tooltip string, forPage bool) *unison.DefaultTableColumnHeader {
	header := NewHeader("", tooltip, forPage)
	baseline := header.Font.Baseline()
	header.Drawable = &widget.DrawableSVGPair{
		Left:  leftSVG,
		Right: rightSVG,
		Size:  geom32.NewSize(baseline*2+4, baseline),
	}
	return header
}

// NewPageRefHeader creates a new page reference header.
func NewPageRefHeader(forPage bool) *unison.DefaultTableColumnHeader {
	return NewSVGHeader(res.BookmarkSVG,
		i18n.Text(`A reference to the book and page the item appears on e.g. B22 would refer to "Basic Set", page 22`),
		forPage)
}

// NewEquippedHeader creates a new equipped header.
func NewEquippedHeader(forPage bool) *unison.DefaultTableColumnHeader {
	return NewSVGHeader(res.CircledCheckSVG,
		i18n.Text(`Whether this piece of equipment is equipped or just carried. Items that are not equipped do not apply any features they may normally contribute to the character.`),
		forPage)
}

// NewMoneyHeader creates a new money header.
func NewMoneyHeader(forPage bool) *unison.DefaultTableColumnHeader {
	return NewSVGHeader(res.CoinsSVG,
		i18n.Text(`The value of one of these pieces of equipment`),
		forPage)
}

// NewExtendedMoneyHeader creates a new extended money page header.
func NewExtendedMoneyHeader(forPage bool) *unison.DefaultTableColumnHeader {
	return NewSVGPairHeader(res.StackSVG, res.CoinsSVG,
		i18n.Text(`The value of all of these pieces of equipment, plus the value of any contained equipment`), forPage)
}

// NewWeightHeader creates a new weight page header.
func NewWeightHeader(forPage bool) unison.TableColumnHeader {
	return NewSVGHeader(res.WeightSVG,
		i18n.Text(`The weight of one of these pieces of equipment`),
		forPage)
}

// NewExtendedWeightHeader creates a new extended weight page header.
func NewExtendedWeightHeader(forPage bool) unison.TableColumnHeader {
	return NewSVGPairHeader(res.StackSVG, res.WeightSVG,
		i18n.Text(`The weight of all of these pieces of equipment, plus the weight of any contained equipment`), forPage)
}
