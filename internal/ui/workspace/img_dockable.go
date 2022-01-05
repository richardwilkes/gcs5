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

const (
	minImageDockableScale   = 0.1
	maxImageDockableScale   = 10
	imageDockableScaleDelta = 0.1
)

var (
	_ unison.Dockable  = &ImageDockable{}
	_ unison.TabCloser = &ImageDockable{}
)

type ImageDockable struct {
	unison.Panel
	path       string
	img        *unison.Image
	imgPanel   *unison.Panel
	scroll     *unison.ScrollPanel
	dragStart  geom32.Point
	dragOrigin geom32.Point
	scale      float32
	inDrag     bool
}

func NewImageDockable(filePath string) (*ImageDockable, error) {
	img, err := unison.NewImageFromFilePathOrURL(filePath, 1)
	if err != nil {
		return nil, err
	}
	d := &ImageDockable{
		path:  filePath,
		img:   img,
		scale: 1,
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{Columns: 1})
	d.imgPanel = unison.NewPanel()
	d.imgPanel.SetSizer(d.imageSizer)
	d.imgPanel.DrawCallback = d.draw
	d.imgPanel.MouseDownCallback = d.mouseDown
	d.imgPanel.MouseDragCallback = d.mouseDrag
	d.imgPanel.MouseUpCallback = d.mouseUp
	d.imgPanel.UpdateCursorCallback = d.updateCursor
	d.imgPanel.MouseWheelCallback = d.mouseWheel
	d.KeyDownCallback = d.keyDown
	d.scroll = unison.NewScrollPanel()
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.SetContent(d.imgPanel, unison.FillBehavior)
	d.AddChild(d.scroll)
	return d, nil
}

func (d *ImageDockable) updateCursor(_ geom32.Point) *unison.Cursor {
	if d.inDrag {
		return unison.ClosedHandCursor()
	}
	return unison.OpenHandCursor()
}

func (d *ImageDockable) mouseDown(where geom32.Point, _, _ int, _ unison.Modifiers) bool {
	d.dragStart = d.imgPanel.PointToRoot(where)
	d.dragOrigin.X, d.dragOrigin.Y = d.scroll.Position()
	d.inDrag = true
	d.RequestFocus()
	d.UpdateCursorNow()
	return true
}

func (d *ImageDockable) mouseDrag(where geom32.Point, _ int, _ unison.Modifiers) bool {
	pt := d.dragStart
	pt.Subtract(d.imgPanel.PointToRoot(where))
	d.scroll.SetPosition(d.dragOrigin.X+pt.X, d.dragOrigin.Y+pt.Y)
	return true
}

func (d *ImageDockable) mouseUp(_ geom32.Point, _ int, _ unison.Modifiers) bool {
	d.inDrag = false
	d.UpdateCursorNow()
	return true
}

func (d *ImageDockable) mouseWheel(_, delta geom32.Point, _ unison.Modifiers) bool {
	d.scale += delta.Y * imageDockableScaleDelta
	if d.scale < minImageDockableScale {
		d.scale = minImageDockableScale
	} else if d.scale > maxImageDockableScale {
		d.scale = maxImageDockableScale
	}
	d.scroll.MarkForLayoutAndRedraw()
	return true
}

func (d *ImageDockable) keyDown(keyCode unison.KeyCode, _ unison.Modifiers, _ bool) bool {
	var scale float32
	switch keyCode {
	case unison.Key1:
		scale = 1
	case unison.Key2:
		scale = 2
	case unison.Key3:
		scale = 3
	case unison.Key4:
		scale = 4
	case unison.Key5:
		scale = 5
	case unison.Key6:
		scale = 6
	case unison.Key7:
		scale = 7
	case unison.Key8:
		scale = 8
	case unison.Key9:
		scale = 9
	case unison.Key0:
		scale = maxImageDockableScale
	case unison.KeyMinus:
		scale -= imageDockableScaleDelta
		if scale < minImageDockableScale {
			scale = minImageDockableScale
		}
	case unison.KeyEqual:
		scale += imageDockableScaleDelta
		if scale > maxImageDockableScale {
			scale = maxImageDockableScale
		}
	default:
		return false
	}
	if d.scale != scale {
		d.scale = scale
		d.scroll.MarkForLayoutAndRedraw()
	}
	return true
}

func (d *ImageDockable) imageSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	pref = d.img.Size()
	pref.Width *= d.scale
	pref.Height *= d.scale
	return geom32.NewSize(50, 50), pref, unison.MaxSize(pref)
}

func (d *ImageDockable) draw(gc *unison.Canvas, dirty geom32.Rect) {
	gc.DrawRect(dirty, unison.ContentColor.Paint(gc, dirty, unison.Fill))
	size := d.img.Size()
	gc.DrawImageInRect(d.img, geom32.NewRect(0, 0, size.Width*d.scale, size.Height*d.scale), nil, nil)
}

func (d *ImageDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

func (d *ImageDockable) Title() string {
	return xfs.BaseName(d.path)
}

func (d *ImageDockable) Tooltip() string {
	return d.path
}

func (d *ImageDockable) Modified() bool {
	return false
}

func (d *ImageDockable) MayAttemptClose() bool {
	return true
}

func (d *ImageDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
