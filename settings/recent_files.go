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

package settings

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	recentFilesMaxKey   = "max"
	recentFilesPathsKey = "paths"
)

// RecentFiles holds a list of recently opened files.
type RecentFiles struct {
	max   int
	paths []string
}

// NewRecentFiles creates a new, empty, RecentFiles object.
func NewRecentFiles() *RecentFiles {
	return &RecentFiles{max: 20}
}

// NewRecentFilesFromJSON creates a new RecentFiles from a JSON object.
func NewRecentFilesFromJSON(data map[string]interface{}) *RecentFiles {
	p := NewRecentFiles()
	if v, ok := data[recentFilesMaxKey]; ok {
		p.max = xmath.MaxInt(int(encoding.Number(v).AsInt64()), 0)
	} else {
		p.max = 20
	}
	if p.max > 0 {
		array := encoding.Array(data[recentFilesPathsKey])
		size := xmath.MinInt(len(array), p.max)
		if size != 0 {
			p.paths = make([]string, 0, size)
			for i, v := range array {
				if i == size {
					break
				}
				p.paths = append(p.paths, encoding.String(v))
			}
		}
	}
	return p
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (p *RecentFiles) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	p.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (p *RecentFiles) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedNumber(recentFilesMaxKey, fixed.F64d4FromInt64(int64(p.max)), false)
	list := p.List()
	if len(list) != 0 {
		encoder.Key(recentFilesPathsKey)
		encoder.StartArray()
		for _, one := range list {
			encoder.String(one)
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}

// Max returns the maximum number of recently opened files to track.
func (p *RecentFiles) Max() int {
	return p.max
}

// SetMax sets the maximum number of recently opened files to track.
func (p *RecentFiles) SetMax(max int) {
	if max < 0 {
		max = 0
	}
	if p.max != max {
		p.max = max
		list := p.List()
		if len(list) > p.max {
			paths := make([]string, p.max)
			copy(paths, list)
			p.paths = paths
		}
	}
}

// List returns the current list of recently opened files. Files that are no longer readable for any reason are omitted.
func (p *RecentFiles) List() []string {
	list := make([]string, 0, len(p.paths))
	for _, one := range p.paths {
		if fs.FileIsReadable(one) {
			list = append(list, one)
		}
	}
	if len(list) != len(p.paths) {
		p.paths = make([]string, len(list))
		copy(p.paths, list)
	}
	return list
}

// Clear the list of recently opened files.
func (p *RecentFiles) Clear() {
	p.paths = nil
}

// Add a file path to the list of recently opened files.
func (p *RecentFiles) Add(filePath string) {
	ext := path.Ext(filePath)
	if runtime.GOOS == toolbox.MacOS || runtime.GOOS == toolbox.WindowsOS {
		ext = strings.ToLower(ext)
	}
	for _, one := range library.AcceptableExtensions() {
		if one == ext {
			full, err := filepath.Abs(filePath)
			if err != nil {
				return
			}
			if fs.FileIsReadable(full) {
				for i, f := range p.paths {
					if f == full {
						copy(p.paths[i:], p.paths[i+1:])
						p.paths[len(p.paths)-1] = ""
						p.paths = p.paths[:len(p.paths)-1]
						break
					}
				}
				p.paths = append(p.paths, "")
				copy(p.paths[1:], p.paths)
				p.paths[0] = full
				if len(p.paths) > p.max {
					p.paths = p.paths[:p.max]
				}
			}
			return
		}
	}
}
