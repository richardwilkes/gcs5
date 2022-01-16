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

package workspace

import (
	"github.com/richardwilkes/gcs/internal/library"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var (
	_ FileBackedDockable = &placeholder{}
	_ unison.TabCloser   = &placeholder{}
)

type placeholder struct {
	unison.Panel
	path string
}

func newPlaceholder(filePath string) *placeholder {
	p := &placeholder{
		path: filePath,
	}
	p.Self = p
	return p
}

func (p *placeholder) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(p.path).SVG,
		Size: suggestedSize,
	}
}

func (p *placeholder) Title() string {
	return xfs.BaseName(p.path)
}

func (p *placeholder) Tooltip() string {
	return p.path
}

func (p *placeholder) BackingFilePath() string {
	return p.path
}

func (p *placeholder) Modified() bool {
	return false
}

func (p *placeholder) MayAttemptClose() bool {
	return true
}

func (p *placeholder) AttemptClose() {
	if dc := unison.DockContainerFor(p); dc != nil {
		dc.Close(p)
	}
}
