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
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/unison"
)

// NewAdvantageListDockable creates a new unison.Dockable for advantage list files.
func NewAdvantageListDockable(filePath string) (unison.Dockable, error) {
	advantages, err := gurps.NewAdvantagesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewListFileDockable(filePath, tbl.NewAdvantageTableHeaders(false), tbl.NewAdvantageRowData(advantages, false)), nil
}

func stringSliceContains(strs []string, text string) bool {
	for _, s := range strs {
		if strings.Contains(strings.ToLower(s), text) {
			return true
		}
	}
	return false
}
