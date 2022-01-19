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
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

// ContextKey is used to insert values into the a context.
type ContextKey int

const (
	// EntityCtxKey holds an *Entity.
	EntityCtxKey ContextKey = iota
)

// SaveWithEntity inserts the provided entity into the context and writes the formatted JSON version of the data to the
// path. Creates any intermediate directories required.
func SaveWithEntity(filePath string, data interface{}, entity *Entity) error {
	ctx := context.Background()
	if entity != nil {
		ctx = context.WithValue(ctx, EntityCtxKey, entity)
	}
	return Save(ctx, filePath, data)
}

// Save writes the formatted JSON version of the data to the path. Creates any intermediate directories required.
func Save(ctx context.Context, filePath string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return errs.Wrap(encoder.EncodeContext(ctx, data))
	}, 0o640); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return nil
}
