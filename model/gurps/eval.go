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
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// InstallEvaluatorFunctions installs additional functions for the evaluator.
func InstallEvaluatorFunctions(m map[string]eval.Function) {
	m["advantage_level"] = evalAdvantageLevel
	m["dice"] = evalDice
	m["roll"] = evalRoll
	m["signed"] = evalSigned
	m["ssrt"] = evalSSRT
	m["ssrt_to_yards"] = evalSSRTYards
}

func evalToBool(e *eval.Evaluator, arguments string) (bool, error) {
	evaluated, err := e.EvaluateNew(arguments)
	if err != nil {
		return false, err
	}
	switch a := evaluated.(type) {
	case bool:
		return a, nil
	case fxp.Int:
		return a != 0, nil
	case string:
		return txt.IsTruthy(a), nil
	default:
		return false, nil
	}
}

func evalToNumber(e *eval.Evaluator, arguments string) (fxp.Int, error) {
	evaluated, err := e.EvaluateNew(arguments)
	if err != nil {
		return 0, err
	}
	return eval.FixedFrom[fxp.DP](evaluated)
}

func evalToString(e *eval.Evaluator, arguments string) (string, error) {
	v, err := e.EvaluateNew(arguments)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", v), nil
}

func evalAdvantageLevel(e *eval.Evaluator, arguments string) (interface{}, error) {
	entity, ok := e.Resolver.(*Entity)
	if !ok || entity.Type != datafile.PC {
		return fxp.NegOne, nil
	}
	arguments = strings.Trim(arguments, `"`)
	levels := fxp.NegOne
	TraverseAdvantages(func(adq *Advantage) bool {
		if strings.EqualFold(adq.Name, arguments) {
			if adq.IsLeveled() {
				levels = adq.Levels
			}
			return true
		}
		return false
	}, true, entity.Advantages...)
	return levels, nil
}

func evalDice(e *eval.Evaluator, arguments string) (interface{}, error) {
	var argList []int
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		n, err := evalToNumber(e, arg)
		if err != nil {
			return nil, err
		}
		argList = append(argList, f64.As[fxp.DP, int](n))
	}
	var d *dice.Dice
	switch len(argList) {
	case 1:
		d = &dice.Dice{
			Count:      1,
			Sides:      argList[0],
			Multiplier: 1,
		}
	case 2:
		d = &dice.Dice{
			Count:      argList[0],
			Sides:      argList[1],
			Multiplier: 1,
		}
	case 3:
		d = &dice.Dice{
			Count:      argList[0],
			Sides:      argList[1],
			Modifier:   argList[2],
			Multiplier: 1,
		}
	case 4:
		d = &dice.Dice{
			Count:      argList[0],
			Sides:      argList[1],
			Modifier:   argList[2],
			Multiplier: argList[3],
		}
	default:
		return nil, errs.New("invalid dice specification")
	}
	return d.String(), nil
}

func evalRoll(e *eval.Evaluator, arguments string) (interface{}, error) {
	if strings.IndexByte(arguments, '(') != -1 {
		var err error
		if arguments, err = evalToString(e, arguments); err != nil {
			return nil, err
		}
	}
	return f64.From[fxp.DP](dice.New(arguments).Roll(false)), nil
}

func evalSigned(e *eval.Evaluator, arguments string) (interface{}, error) {
	n, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return n.StringWithSign(), nil
}

func evalSSRT(e *eval.Evaluator, arguments string) (interface{}, error) {
	// Takes 3 args: length (number), units (string), flag (bool) indicating for size (true) or speed/range (false)
	var arg string
	arg, arguments = eval.NextArg(arguments)
	n, err := evalToString(e, arg)
	if err != nil {
		return nil, err
	}
	arg, arguments = eval.NextArg(arguments)
	var units string
	if units, err = evalToString(e, arg); err != nil {
		return nil, err
	}
	arg, _ = eval.NextArg(arguments)
	var wantSize bool
	if wantSize, err = evalToBool(e, arg); err != nil {
		return nil, err
	}
	var length measure.Length
	if length, err = measure.LengthFromString(n+" "+units, measure.Yard); err != nil {
		return nil, err
	}
	result := yardsToValue(length, wantSize)
	if !wantSize {
		result = -result
	}
	return f64.From[fxp.DP](result), nil
}

func evalSSRTYards(e *eval.Evaluator, arguments string) (interface{}, error) {
	v, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return valueToYards(f64.As[fxp.DP, int](v)), nil
}

func yardsToValue(length measure.Length, allowNegative bool) int {
	inches := fxp.Int(length)
	yards := inches.Div(fxp.ThirtySix)
	if allowNegative {
		switch {
		case inches <= fxp.One.Div(fxp.Five):
			return -15
		case inches <= fxp.One.Div(fxp.Three):
			return -14
		case inches <= fxp.Half:
			return -13
		case inches <= fxp.Two.Div(fxp.Three):
			return -12
		case inches <= fxp.One:
			return -11
		case inches <= fxp.OneAndAHalf:
			return -10
		case inches <= fxp.Two:
			return -9
		case inches <= fxp.Three:
			return -8
		case inches <= fxp.Five:
			return -7
		case inches <= fxp.Eight:
			return -6
		}
		feet := inches.Div(fxp.Twelve)
		switch {
		case feet <= fxp.One:
			return -5
		case feet <= fxp.OneAndAHalf:
			return -4
		case feet <= fxp.Two:
			return -3
		case yards <= fxp.One:
			return -2
		case yards <= fxp.OneAndAHalf:
			return -1
		}
	}
	if yards <= fxp.Two {
		return 0
	}
	amt := 0
	for yards > fxp.Ten {
		yards = yards.Div(fxp.Ten)
		amt += 6
	}
	switch {
	case yards > fxp.Seven:
		return amt + 4
	case yards > fxp.Five:
		return amt + 3
	case yards > fxp.Three:
		return amt + 2
	case yards > fxp.Two:
		return amt + 1
	case yards > fxp.OneAndAHalf:
		return amt
	default:
		return amt - 1
	}
}

func valueToYards(value int) fxp.Int {
	if value < -15 {
		value = -15
	}
	switch value {
	case -15:
		return fxp.One.Div(fxp.Five).Div(fxp.ThirtySix)
	case -14:
		return fxp.One.Div(fxp.Three).Div(fxp.ThirtySix)
	case -13:
		return fxp.Half.Div(fxp.ThirtySix)
	case -12:
		return fxp.Two.Div(fxp.Three).Div(fxp.ThirtySix)
	case -11:
		return fxp.One.Div(fxp.ThirtySix)
	case -10:
		return fxp.OneAndAHalf.Div(fxp.ThirtySix)
	case -9:
		return fxp.Two.Div(fxp.ThirtySix)
	case -8:
		return fxp.Three.Div(fxp.ThirtySix)
	case -7:
		return fxp.Five.Div(fxp.ThirtySix)
	case -6:
		return fxp.Eight.Div(fxp.ThirtySix)
	case -5:
		return fxp.One.Div(fxp.Three)
	case -4:
		return fxp.OneAndAHalf.Div(fxp.Three)
	case -3:
		return fxp.Two.Div(fxp.Three)
	case -2:
		return fxp.One
	case -1:
		return fxp.OneAndAHalf
	case 0:
		return fxp.Two
	case 1:
		return fxp.Three
	case 2:
		return fxp.Five
	case 3:
		return fxp.Seven
	}
	value -= 4
	multiplier := fxp.One
	for i := 0; i < value/6; i++ {
		multiplier = multiplier.Mul(fxp.Ten)
	}
	var v fxp.Int
	switch value % 6 {
	case 0:
		v = fxp.Ten
	case 1:
		v = fxp.Fifteen
	case 2:
		v = fxp.Twenty
	case 3:
		v = fxp.Thirty
	case 4:
		v = fxp.Fifty
	case 5:
		v = fxp.Seventy
	}
	return v.Mul(multiplier)
}
