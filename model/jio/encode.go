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

package jio

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

// Save the data as JSON to the given path. Parent directories will be created automatically, if needed.
func Save(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return errs.Wrap(err)
	}
	if err := safe.WriteFileWithMode(path, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		return errs.Wrap(encoder.Encode(data))
	}, 0o640); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}
