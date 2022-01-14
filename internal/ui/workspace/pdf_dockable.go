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
	"fmt"
	"os"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/pdf"
	"github.com/richardwilkes/toolbox/errs"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const (
	minPDFDockableScale   = 0.5
	maxPDFDockableScale   = 3
	PDFDockableScaleDelta = 0.1
)

var (
	_ FileBackedDockable = &PDFDockable{}
	_ unison.TabCloser   = &PDFDockable{}
)

type PDFDockable struct {
	unison.Panel
	path             string
	doc              *pdf.Document
	pageCount        int
	pageNumber       int
	loadedPageNumber int
	toc              []*pdf.TOCEntry
	docPanel         *unison.Panel
	page             *pdf.RenderedPage
	pageError        string
	pageImg          *unison.Image
	scroll           *unison.ScrollPanel
	rolloverRect     geom32.Rect
	link             *pdf.PageLink
	scale            float32
	loadedScale      float32
	hidpiScale       float32
}

func NewPDFDockable(filePath string) (*PDFDockable, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	var doc *pdf.Document
	if doc, err = pdf.New(data, 0); err != nil {
		return nil, errs.Wrap(err)
	}
	d := &PDFDockable{
		path:             filePath,
		doc:              doc,
		pageCount:        doc.PageCount(),
		loadedPageNumber: -1,
		scale:            1,
		hidpiScale:       0.5, // TODO: Use monitor info to set this?
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{Columns: 1})
	d.docPanel = unison.NewPanel()
	d.docPanel.SetSizer(d.docSizer)
	d.docPanel.DrawCallback = d.draw
	d.docPanel.MouseDownCallback = d.mouseDown
	d.docPanel.MouseMoveCallback = d.mouseMove
	d.docPanel.MouseUpCallback = d.mouseUp
	d.docPanel.UpdateCursorCallback = d.updateCursor
	d.KeyDownCallback = d.keyDown
	d.scroll = unison.NewScrollPanel()
	d.scroll.MouseWheelMultiplier = 4
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.SetContent(d.docPanel, unison.FillBehavior)
	d.AddChild(d.scroll)
	return d, nil
}

func (d *PDFDockable) updateCursor(_ geom32.Point) *unison.Cursor {
	// TODO: change cursor when over clickable link
	return unison.ArrowCursor()
}

func (d *PDFDockable) overLink(where geom32.Point) (rect geom32.Rect, link *pdf.PageLink) {
	if d.page != nil && d.page.Links != nil {
		for _, link = range d.page.Links {
			rect.X = float32(link.Bounds.Min.X) * d.hidpiScale
			rect.Y = float32(link.Bounds.Min.Y) * d.hidpiScale
			rect.Width = float32(link.Bounds.Dx()) * d.hidpiScale
			rect.Height = float32(link.Bounds.Dy()) * d.hidpiScale
			if rect.ContainsPoint(where) {
				return rect, link
			}
		}
	}
	return rect, nil
}

func (d *PDFDockable) mouseDown(where geom32.Point, _, _ int, _ unison.Modifiers) bool {
	d.RequestFocus()
	d.UpdateCursorNow()
	return true
}

func (d *PDFDockable) mouseMove(where geom32.Point, mod unison.Modifiers) bool {
	r, link := d.overLink(where)
	if r != d.rolloverRect || link != d.link {
		d.rolloverRect = r
		d.link = link
		d.MarkForRedraw()
	}
	return true
}

func (d *PDFDockable) mouseUp(where geom32.Point, button int, _ unison.Modifiers) bool {
	r, link := d.overLink(where)
	if r != d.rolloverRect || link != d.link {
		d.rolloverRect = r
		d.link = link
		d.MarkForRedraw()
	}
	d.UpdateCursorNow()
	if button == unison.ButtonLeft && d.link != nil {
		if d.link.PageNumber >= 0 {
			d.pageNumber = d.link.PageNumber
			d.MarkForRedraw()
			// TODO: Use d.link.PageX & PageY to ensure location is scrolled into place
		}
	}
	return true
}

func (d *PDFDockable) keyDown(keyCode unison.KeyCode, _ unison.Modifiers, _ bool) bool {
	scale := d.scale
	switch keyCode {
	case unison.Key1:
		scale = 1
	case unison.Key2:
		scale = 2
	case unison.Key3:
		scale = 3
	case unison.KeyMinus:
		scale -= PDFDockableScaleDelta
		if scale < minPDFDockableScale {
			scale = minPDFDockableScale
		}
	case unison.KeyEqual:
		scale += PDFDockableScaleDelta
		if scale > maxPDFDockableScale {
			scale = maxPDFDockableScale
		}
	case unison.KeyLeft:
		if d.pageNumber > 0 {
			d.pageNumber--
			d.MarkForRedraw()
		}
	case unison.KeyRight:
		if d.pageNumber < d.pageCount-1 {
			d.pageNumber++
			d.MarkForRedraw()
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

func (d *PDFDockable) dpi() int {
	// Using a baseline dpi of 96, since that's what most software does nowadays, rather than 72.
	// Multiplied by 2 to support crisp rendering on high-dpi displays.
	// TODO: Consider using the scale from the monitor instead of just multiplying by 2.
	return int(d.scale * 96 * 2)
}

func (d *PDFDockable) ensurePageIsLoaded() {
	dpi := d.dpi()
	if d.scale != d.loadedScale {
		d.toc = d.doc.TableOfContents(dpi)
	}
	if d.pageNumber != d.loadedPageNumber || d.scale != d.loadedScale {
		d.pageImg = nil
		d.page = nil
		d.link = nil
		var err error
		if d.page, err = d.doc.RenderPage(d.pageNumber, dpi, 0, ""); err != nil {
			d.pageImg = nil
			d.pageError = err.Error()
		} else {
			if d.pageImg, err = unison.NewImageFromPixels(d.page.Image.Rect.Dx(), d.page.Image.Rect.Dy(),
				d.page.Image.Pix, d.hidpiScale); err != nil {
				d.pageImg = nil
				d.pageError = err.Error()
			} else {
				d.pageError = ""
			}
		}
	}
	d.loadedPageNumber = d.pageNumber
	d.loadedScale = d.scale
}

func (d *PDFDockable) docSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	d.ensurePageIsLoaded()
	if d.pageImg != nil {
		pref = d.pageImg.LogicalSize()
	} else {
		// TODO: Size based on error message
		pref.Width = 400
		pref.Height = 300
	}
	return geom32.NewSize(50, 50), pref, unison.MaxSize(pref)
}

func (d *PDFDockable) draw(gc *unison.Canvas, dirty geom32.Rect) {
	d.ensurePageIsLoaded()
	gc.DrawRect(dirty, unison.ContentColor.Paint(gc, dirty, unison.Fill))
	if d.pageImg != nil {
		r := geom32.Rect{Size: d.pageImg.LogicalSize()}
		gc.DrawRect(r, unison.White.Paint(gc, r, unison.Fill))
		gc.DrawImageInRect(d.pageImg, r, nil, nil)
		if d.link != nil {
			gc.DrawRect(d.rolloverRect, unison.GreenYellow.SetAlphaIntensity(0.3).Paint(gc, d.rolloverRect, unison.Fill))
		}
	} else {
		// TODO: Show error message
		fmt.Println("d.pageImg was nil")
	}
}

func (d *PDFDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

func (d *PDFDockable) Title() string {
	return xfs.BaseName(d.path)
}

func (d *PDFDockable) Tooltip() string {
	return d.path
}

func (d *PDFDockable) BackingFilePath() string {
	return d.path
}

func (d *PDFDockable) Modified() bool {
	return false
}

func (d *PDFDockable) MayAttemptClose() bool {
	return true
}

func (d *PDFDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
