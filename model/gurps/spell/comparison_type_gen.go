// Code generated from "enum.go.tmpl" - DO NOT EDIT.

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

package spell

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible values.
const (
	Name ComparisonType = iota
	Category
	College
	CollegeCount
	Any
	LastComparisonType = Any
)

var (
	// AllComparisonType holds all possible values.
	AllComparisonType = []ComparisonType{
		Name,
		Category,
		College,
		CollegeCount,
		Any,
	}
	comparisonTypeData = []struct {
		key    string
		string string
	}{
		{
			key:    "name",
			string: i18n.Text("Name"),
		},
		{
			key:    "category",
			string: i18n.Text("Category"),
		},
		{
			key:    "college",
			string: i18n.Text("College"),
		},
		{
			key:    "college_count",
			string: i18n.Text("College Count"),
		},
		{
			key:    "any",
			string: i18n.Text("Any"),
		},
	}
)

// ComparisonType holds the type of a comparison.
type ComparisonType byte

// EnsureValid ensures this is of a known value.
func (enum ComparisonType) EnsureValid() ComparisonType {
	if enum <= LastComparisonType {
		return enum
	}
	return 0
}

// Key returns the key used in serialization.
func (enum ComparisonType) Key() string {
	return comparisonTypeData[enum.EnsureValid()].key
}

// String implements fmt.Stringer.
func (enum ComparisonType) String() string {
	return comparisonTypeData[enum.EnsureValid()].string
}

// ExtractComparisonType extracts the value from a string.
func ExtractComparisonType(str string) ComparisonType {
	for i, one := range comparisonTypeData {
		if strings.EqualFold(one.key, str) {
			return ComparisonType(i)
		}
	}
	return 0
}

// MarshalText implements the encoding.TextMarshaler interface.
func (enum ComparisonType) MarshalText() (text []byte, err error) {
	return []byte(enum.Key()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (enum *ComparisonType) UnmarshalText(text []byte) error {
	*enum = ExtractComparisonType(string(text))
	return nil
}
