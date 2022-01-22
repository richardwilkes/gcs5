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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	exportInfoFilePathKey     = "file_path"
	exportInfoTemplatePathKey = "template_path"
	exportInfoExportPathKey   = "export_path"
	exportInfoLastUsedKey     = "last_used"
)

// ExportInfo holds information about a recent export so that it can be redone quickly.
type ExportInfo struct {
	FilePath     string
	TemplatePath string
	ExportPath   string
	LastUsed     int64
}

// NewExportInfoFromJSON creates a new ExportInfo from a JSON object.
func NewExportInfoFromJSON(data map[string]interface{}) *ExportInfo {
	return &ExportInfo{
		FilePath:     encoding.String(data[exportInfoFilePathKey]),
		TemplatePath: encoding.String(data[exportInfoTemplatePathKey]),
		ExportPath:   encoding.String(data[exportInfoExportPathKey]),
		LastUsed:     encoding.Number(data[exportInfoLastUsedKey]).AsInt64(),
	}
}

// ToJSON emits this object as JSON.
func (e *ExportInfo) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(exportInfoFilePathKey, e.FilePath, false, false)
	encoder.KeyedString(exportInfoTemplatePathKey, e.TemplatePath, false, false)
	encoder.KeyedString(exportInfoExportPathKey, e.ExportPath, false, false)
	encoder.KeyedNumber(exportInfoLastUsedKey, fixed.F64d4FromInt64(e.LastUsed), false)
	encoder.EndObject()
}
