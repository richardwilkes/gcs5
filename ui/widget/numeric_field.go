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
	"strings"
	"unicode"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/unison"
)

// NumericField holds the value for a numeric field.
type NumericField struct {
	*unison.Field
	undoID     int64
	undoTitle  string
	get        func() fxp.Int
	set        func(fxp.Int)
	minimum    fxp.Int
	maximum    fxp.Int
	noMinWidth bool
	inUndo     bool
}

// NewNumericField creates a new field that holds a fixed-point number.
func NewNumericField(undoTitle string, get func() fxp.Int, set func(fxp.Int), min, max fxp.Int, noMinWidth bool) *NumericField {
	f := &NumericField{
		Field:      unison.NewField(),
		undoID:     unison.NextUndoID(),
		undoTitle:  undoTitle,
		get:        get,
		set:        set,
		minimum:    min,
		noMinWidth: noMinWidth,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.RuneTypedCallback = f.runeTyped
	f.SetMaximum(max)
	f.Sync()
	return f
}

// SetMaximum sets the maximum value allowed.
func (f *NumericField) SetMaximum(maximum fxp.Int) {
	f.maximum = maximum
	if !f.noMinWidth && f.minimum != fxp.Min && f.maximum != fxp.Max {
		f.MinimumTextWidth = xmath.Max(f.Font.SimpleWidth((f.minimum.Trunc() + fxp.One - 1).String()),
			f.Font.SimpleWidth((f.maximum.Trunc() + fxp.One - 1).String()))
	}
}

func (f *NumericField) trimmed(text string) string {
	return strings.TrimSpace(text)
}

func (f *NumericField) validate() bool {
	v, err := fxp.FromString(f.trimmed(f.Text()))
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid number"))
		return false
	}
	if f.minimum != fxp.Min && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Number must be at least %s"), f.minimum.String()))
		return false
	}
	if f.maximum != fxp.Max && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Number must be no more than %s"), f.maximum.String()))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *NumericField) modified() {
	text := f.Text()
	if !f.inUndo && f.undoID != unison.NoUndoID {
		if mgr := unison.UndoManagerFor(f); mgr != nil {
			mgr.Add(&unison.UndoEdit[string]{
				ID:       f.undoID,
				EditName: f.undoTitle,
				EditCost: 1,
				UndoFunc: func(e *unison.UndoEdit[string]) { f.setWithoutUndo(e.BeforeData, true) },
				RedoFunc: func(e *unison.UndoEdit[string]) { f.setWithoutUndo(e.AfterData, true) },
				AbsorbFunc: func(e *unison.UndoEdit[string], other unison.Undoable) bool {
					if e2, ok := other.(*unison.UndoEdit[string]); ok && e2.ID == f.undoID {
						e.AfterData = e2.AfterData
						return true
					}
					return false
				},
				BeforeData: f.get().String(),
				AfterData:  text,
			})
		}
	}
	if v, err := fxp.FromString(f.trimmed(text)); err == nil &&
		(f.minimum == fxp.Min || v >= f.minimum) &&
		(f.maximum == fxp.Max || v <= f.maximum) && f.get() != v {
		f.set(v)
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
}

func (f *NumericField) setWithoutUndo(text string, focus bool) {
	f.inUndo = true
	f.SetText(text)
	f.inUndo = false
	if focus {
		f.RequestFocus()
		f.SelectAll()
	}
}

// Sync the field to the current value.
func (f *NumericField) Sync() {
	value := f.get()
	if f.minimum != fxp.Min && value < f.minimum {
		value = f.minimum
	} else if f.maximum != fxp.Max && value > f.maximum {
		value = f.maximum
	}
	f.setWithoutUndo(value.String(), false)
}

func (f *NumericField) runeTyped(ch rune) bool {
	if !unicode.IsControl(ch) {
		if f.minimum >= 0 && ch == '-' {
			unison.Beep()
			return false
		}
		if text := f.trimmed(string(f.RunesIfPasted([]rune{ch}))); text != "-" {
			if _, err := fxp.FromString(text); err != nil {
				unison.Beep()
				return false
			}
		}
	}
	return f.DefaultRuneTyped(ch)
}
