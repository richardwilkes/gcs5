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

package widget

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/richardwilkes/gcs/ui/undo"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// IntegerField holds the value for an integer field.
type IntegerField struct {
	*unison.Field
	undoID    int
	undoTitle string
	get       func() int
	set       func(int)
	minimum   int
	maximum   int
	showSign  bool
	inUndo    bool
}

// NewIntegerField creates a new field that holds an integer.
func NewIntegerField(undoID int, undoTitle string, get func() int, set func(int), min, max int, showSign bool) *IntegerField {
	f := &IntegerField{
		Field:     unison.NewField(),
		undoID:    undoID,
		undoTitle: undoTitle,
		get:       get,
		set:       set,
		minimum:   min,
		maximum:   max,
		showSign:  showSign,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.RuneTypedCallback = f.runeTyped
	if min != math.MinInt && max != math.MaxInt {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(f.formatted(min)), f.Font.SimpleWidth(f.formatted(max)))
	}
	f.Sync()
	return f
}

func (f *IntegerField) formatted(value int) string {
	if f.showSign {
		return fmt.Sprintf("%+d", value)
	}
	return strconv.Itoa(value)
}

func (f *IntegerField) trimmed(text string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), "+"))
}

func (f *IntegerField) validate() bool {
	v, err := strconv.Atoi(f.trimmed(f.Text()))
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid integer"))
		return false
	}
	if f.minimum != math.MinInt && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be at least %s"), f.formatted(f.minimum)))
		return false
	}
	if f.maximum != math.MaxInt && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be no more than %s"), f.formatted(f.maximum)))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *IntegerField) modified() {
	text := f.Text()
	if !f.inUndo && f.undoID != undo.NoneID {
		if mgr := undo.Manager(f); mgr != nil {
			mgr.Add(&undo.Edit[string]{
				ID:       f.undoID,
				EditName: f.undoTitle,
				EditCost: 1,
				UndoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.BeforeData, true) },
				RedoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.AfterData, true) },
				AbsorbFunc: func(e *undo.Edit[string], other unison.UndoEdit) bool {
					if e2, ok := other.(*undo.Edit[string]); ok && e2.ID == f.undoID {
						e.AfterData = e2.AfterData
						return true
					}
					return false
				},
				BeforeData: f.formatted(f.get()),
				AfterData:  text,
			})
		}
	}
	if v, err := strconv.Atoi(f.trimmed(text)); err == nil &&
		(f.minimum == math.MinInt || v >= f.minimum) &&
		(f.maximum == math.MaxInt || v <= f.maximum) && f.get() != v {
		f.set(v)
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
}

func (f *IntegerField) setWithoutUndo(text string, focus bool) {
	f.inUndo = true
	f.SetText(text)
	f.inUndo = false
	if focus {
		f.RequestFocus()
		f.SelectAll()
	}
}

// Sync the field to the current value.
func (f *IntegerField) Sync() {
	value := f.get()
	if f.minimum != math.MinInt && value < f.minimum {
		value = f.minimum
	} else if f.maximum != math.MaxInt && value > f.maximum {
		value = f.maximum
	}
	f.setWithoutUndo(f.formatted(value), false)
}

func (f *IntegerField) runeTyped(ch rune) bool {
	if !unicode.IsControl(ch) {
		if f.minimum >= 0 && ch == '-' {
			unison.Beep()
			return false
		}
		if text := f.trimmed(string(f.RunesIfPasted([]rune{ch}))); text != "-" {
			if _, err := strconv.Atoi(text); err != nil {
				unison.Beep()
				return false
			}
		}
	}
	return f.DefaultRuneTyped(ch)
}
