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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/widget/ntable"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type adjustRawPointsList[T gurps.NodeConstraint[T]] struct {
	Owner widget.Rebuildable
	List  []*rawPointsAdjuster[T]
}

func (a *adjustRawPointsList[T]) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	a.Finish()
}

func (a *adjustRawPointsList[T]) Finish() {
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type rawPointsAdjuster[T gurps.NodeConstraint[T]] struct {
	Target gurps.RawPointsAdjuster[T]
	Points fxp.Int
}

func newRawPointsAdjuster[T gurps.NodeConstraint[T]](target gurps.RawPointsAdjuster[T]) *rawPointsAdjuster[T] {
	return &rawPointsAdjuster[T]{
		Target: target,
		Points: target.RawPoints(),
	}
}

func (a *rawPointsAdjuster[T]) Apply() {
	a.Target.SetRawPoints(a.Points)
}

func canAdjustRawPoints[T gurps.NodeConstraint[T]](table *unison.Table[*ntable.Node[T]], increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if provider, ok := interface{}(row.Data()).(gurps.RawPointsAdjuster[T]); ok && !provider.Container() {
			if increment || provider.RawPoints() > 0 {
				return true
			}
		}
	}
	return false
}

func adjustRawPoints[T gurps.NodeConstraint[T]](owner widget.Rebuildable, table *unison.Table[*ntable.Node[T]], increment bool) {
	before := &adjustRawPointsList[T]{Owner: owner}
	after := &adjustRawPointsList[T]{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if provider, ok := interface{}(row.Data()).(gurps.RawPointsAdjuster[T]); ok {
			if increment || provider.RawPoints() > 0 {
				before.List = append(before.List, newRawPointsAdjuster[T](provider))
				rawPts := provider.RawPoints()
				pts := rawPts.Trunc()
				if increment {
					pts += fxp.One
				} else if rawPts == pts {
					pts -= fxp.One
				}
				provider.SetRawPoints(pts.Max(0))
				after.List = append(after.List, newRawPointsAdjuster[T](provider))
			}
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			var name string
			if increment {
				name = i18n.Text("Increment Points")
			} else {
				name = i18n.Text("Decrement Points")
			}
			mgr.Add(&unison.UndoEdit[*adjustRawPointsList[T]]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit *unison.UndoEdit[*adjustRawPointsList[T]]) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit *unison.UndoEdit[*adjustRawPointsList[T]]) { edit.AfterData.Apply() },
				BeforeData: before,
				AfterData:  after,
			})
		}
		before.Finish()
	}
}
