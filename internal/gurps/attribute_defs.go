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
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/internal/id"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
)

// AttributeDefs holds a slice of AttributeDef.
type AttributeDefs []*AttributeDef

// FactoryAttributeDefs returns the attribute factory settings.
func FactoryAttributeDefs() AttributeDefs {
	var defs AttributeDefs
	jot.FatalIfErr(xfs.LoadJSONFromFS(embeddedFS, "data/standard.attr", &defs))
	return defs
}

// UnmarshalJSON implements json.Unmarshaler. Loads the current format as well as older variants.
func (a *AttributeDefs) UnmarshalJSON(data []byte) error {
	var current []*AttributeDef
	if err := json.Unmarshal(data, &current); err != nil {
		var variants struct {
			JavaVersion []*AttributeDef `json:"attributes"`
		}
		if err2 := json.Unmarshal(data, &variants); err2 != nil {
			return err
		}
		*a = variants.JavaVersion
	} else {
		*a = current
	}
	set := make(map[string]bool)
	for _, one := range *a {
		one.ID = id.Sanitize(one.ID, false, ReservedAttributeDefIDs...)
		if set[one.ID] {
			return errs.New("duplicate ID in attributes: " + one.ID)
		}
		set[one.ID] = true
	}
	return nil
}

// SaveTo saves the AttributeDefs data to the specified file.
func (a AttributeDefs) SaveTo(filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return xfs.SaveJSONWithMode(filePath, a, true, 0o640)
}
