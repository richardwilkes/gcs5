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

package encoding

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// LoadJSON data from the specified path.
func LoadJSON(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errs.NewWithCause(path, err)
	}
	return loadJSON(f, path)
}

// LoadJSONFromFS data from the specified filesystem path.
func LoadJSONFromFS(fsys fs.FS, path string) (interface{}, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, errs.NewWithCause(path, err)
	}
	return loadJSON(f, path)
}

func loadJSON(r io.ReadCloser, path string) (interface{}, error) {
	defer xio.CloseIgnoringErrors(r)
	var data interface{}
	decoder := json.NewDecoder(bufio.NewReader(r))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, errs.NewWithCause(path, err)
	}
	return data, nil
}

// Bool extracts a bool value from the data. Supports bool, string, json.Number as input. All other data types return
// false.
func Bool(data interface{}) bool {
	switch d := data.(type) {
	case bool:
		return d
	case string, json.Number:
		return d != "0"
	default:
		return false
	}
}

// Number returns a fixed.F64d4 from the data. Supports json.Number, string, bool as input. All other data types return
// 0.
func Number(data interface{}) fixed.F64d4 {
	switch d := data.(type) {
	case json.Number:
		return fixed.F64d4FromStringForced(string(d))
	case string:
		return fixed.F64d4FromStringForced(d)
	case bool:
		if d {
			return fixed.F64d4FromInt64(1)
		}
		return 0
	default:
		return 0
	}
}

// String returns a string from the data. Supports string, json.Number, bool as input. All other data types return "".
func String(data interface{}) string {
	switch d := data.(type) {
	case string:
		return d
	case json.Number:
		return string(d)
	case bool:
		if d {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// Array returns a slice from the data or nil if the data isn't an array.
func Array(data interface{}) []interface{} {
	if a, ok := data.([]interface{}); ok {
		return a
	}
	return nil
}

// Object returns a map from the data or nil if the data isn't an object.
func Object(data interface{}) map[string]interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		return m
	}
	return nil
}
