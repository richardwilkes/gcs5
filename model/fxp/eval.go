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

package fxp

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/eval/f64d4eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// EvaluateToNumber evaluates the provided expression and returns a number.
func EvaluateToNumber(expression string, resolver eval.VariableResolver) fixed.F64d4 {
	result, err := f64d4eval.NewEvaluator(resolver, true).Evaluate(expression)
	if err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", expression))
		return 0
	}
	if value, ok := result.(fixed.F64d4); ok {
		return value
	}
	if str, ok := result.(string); ok {
		var value fixed.F64d4
		if value, err = fixed.F64d4FromString(str); err == nil {
			return value
		}
	}
	jot.Warn(errs.Newf("unable to resolve '%s' to a number", expression))
	return 0

}
