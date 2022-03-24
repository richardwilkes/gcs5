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

package undo

import "github.com/richardwilkes/unison"

var _ unison.UndoEdit = &Edit[int]{}

type Edit[T any] struct {
	ID          int
	EditName    string
	EditCost    int
	UndoFunc    func(*Edit[T])
	RedoFunc    func(*Edit[T])
	AbsorbFunc  func(*Edit[T], unison.UndoEdit) bool
	ReleaseFunc func(*Edit[T])
	BeforeData  T
	AfterData   T
}

func (e *Edit[T]) Name() string {
	return e.EditName
}

func (e *Edit[T]) Cost() int {
	return e.EditCost
}

func (e *Edit[T]) Undo() {
	if e.UndoFunc != nil {
		e.UndoFunc(e)
	}
}

func (e *Edit[T]) Redo() {
	if e.RedoFunc != nil {
		e.RedoFunc(e)
	}
}

func (e *Edit[T]) Absorb(other unison.UndoEdit) bool {
	if e.AbsorbFunc != nil {
		return e.AbsorbFunc(e, other)
	}
	return false
}

func (e *Edit[T]) Release() {
	if e.ReleaseFunc != nil {
		e.ReleaseFunc(e)
	}
}
