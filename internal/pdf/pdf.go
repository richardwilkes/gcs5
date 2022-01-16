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

package pdf

import (
	"image"
	"os"
	"sync"

	"github.com/richardwilkes/pdf"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// TOC holds a table of contents entry.
type TOC struct {
	Title        string
	PageNumber   int
	PageLocation geom32.Point
	Children     []*TOC
}

// Page holds a rendered PDF page.
type Page struct {
	Error      error
	PageNumber int
	Image      *unison.Image
	TOC        []*TOC
	Links      []*Link
	Matches    []geom32.Rect
}

// Link holds a single link on a page. If PageNumber if >= 0, then this is an internal link and the URI will be empty.
type Link struct {
	Bounds       geom32.Rect
	PageNumber   int
	PageLocation geom32.Point
	URI          string
}

type params struct {
	sequence   int
	pageNumber int
	search     string
	scale      float32
}

// PDF holds a PDF page renderer.
type PDF struct {
	MaxSearchMatches   int
	PPI                int
	DisplayScaleAdjust float32
	doc                *pdf.Document
	pageCount          int
	pageLoadedCallback func()
	lock               sync.RWMutex
	page               *Page
	lastRequest        *params
	rendering          *params
	sequence           int
}

// New creates a new PDF page renderer.
func New(filePath string, pageLoadedCallback func()) (*PDF, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	var doc *pdf.Document
	if doc, err = pdf.New(data, 0); err != nil {
		return nil, errs.Wrap(err)
	}
	display := unison.PrimaryDisplay()
	return &PDF{
		MaxSearchMatches:   100,
		PPI:                display.PPI(),
		DisplayScaleAdjust: 1 / display.ScaleX,
		doc:                doc,
		pageCount:          doc.PageCount(),
		pageLoadedCallback: pageLoadedCallback,
	}, nil
}

// PageCount returns the total page count.
func (p *PDF) PageCount() int {
	return p.pageCount
}

// CurrentPage returns the currently rendered page.
func (p *PDF) CurrentPage() *Page {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.page
}

// MostRecentPageNumber returns the most recent page number that has been asked to be rendered.
func (p *PDF) MostRecentPageNumber() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.rendering != nil {
		return p.rendering.pageNumber
	}
	if p.page != nil {
		return p.page.PageNumber
	}
	return 0
}

// LoadPage requests the given page to be loaded and rendered at the specified scale.
func (p *PDF) LoadPage(pageNumber int, scale float32, search string) {
	if pageNumber < 0 || pageNumber >= p.pageCount {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.lastRequest != nil && p.lastRequest.pageNumber == pageNumber && p.lastRequest.scale == scale &&
		p.lastRequest.search == search {
		return
	}
	p.sequence++
	p.lastRequest = &params{
		sequence:   p.sequence,
		pageNumber: pageNumber,
		scale:      scale,
		search:     search,
	}
	if p.rendering == nil {
		p.rendering = p.lastRequest
		go p.render(p.rendering)
	}
}

func (p *PDF) render(state *params) {
	if p.shouldAbortRender() {
		return
	}

	dpi := int(state.scale * float32(p.PPI) / p.DisplayScaleAdjust)
	toc := p.doc.TableOfContents(dpi)
	if p.shouldAbortRender() {
		return
	}

	page, err := p.doc.RenderPage(state.pageNumber, dpi, p.MaxSearchMatches, state.search)
	if err != nil {
		p.errorDuringRender(state.pageNumber, err)
		return
	}
	if p.shouldAbortRender() {
		return
	}

	var img *unison.Image
	img, err = unison.NewImageFromPixels(page.Image.Rect.Dx(), page.Image.Rect.Dy(), page.Image.Pix, p.DisplayScaleAdjust)
	if err != nil {
		p.errorDuringRender(state.pageNumber, err)
		return
	}
	p.lock.Lock()
	if p.rendering.sequence != p.lastRequest.sequence {
		p.rendering = p.lastRequest
		go p.render(p.rendering)
		p.lock.Unlock()
		return
	}
	p.rendering = nil
	p.page = &Page{
		PageNumber: state.pageNumber,
		Image:      img,
		TOC:        p.convertTOCEntries(toc),
		Links:      p.convertLinks(page.Links),
		Matches:    p.convertMatches(page.SearchHits),
	}
	p.lock.Unlock()
	p.pageLoadedCallback()
}

func (p *PDF) shouldAbortRender() bool {
	p.lock.RLock()
	abort := p.rendering.sequence != p.lastRequest.sequence
	p.lock.RUnlock()
	if abort {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.rendering = p.lastRequest
		go p.render(p.rendering)
	}
	return abort
}

func (p *PDF) errorDuringRender(pageNumber int, err error) {
	p.lock.Lock()
	if p.rendering.sequence != p.lastRequest.sequence {
		p.rendering = p.lastRequest
		go p.render(p.rendering)
		p.lock.Unlock()
		return
	}
	p.page = &Page{
		Error:      err,
		PageNumber: pageNumber,
	}
	p.lock.Unlock()
	p.pageLoadedCallback()
}

func (p *PDF) convertTOCEntries(entries []*pdf.TOCEntry) []*TOC {
	if len(entries) == 0 {
		return nil
	}
	toc := make([]*TOC, len(entries))
	for i, entry := range entries {
		toc[i] = &TOC{
			Title:        entry.Title,
			PageNumber:   entry.PageNumber,
			PageLocation: p.pointFromPagePoint(entry.PageX, entry.PageY),
			Children:     p.convertTOCEntries(entry.Children),
		}
	}
	return toc
}

func (p *PDF) convertLinks(pageLinks []*pdf.PageLink) []*Link {
	if len(pageLinks) == 0 {
		return nil
	}
	links := make([]*Link, len(pageLinks))
	for i, link := range pageLinks {
		links[i] = &Link{
			Bounds:       p.rectFromPageRect(link.Bounds),
			PageNumber:   link.PageNumber,
			PageLocation: p.pointFromPagePoint(link.PageX, link.PageY),
			URI:          link.URI,
		}
	}
	return links
}

func (p *PDF) convertMatches(hits []image.Rectangle) []geom32.Rect {
	if len(hits) == 0 {
		return nil
	}
	matches := make([]geom32.Rect, len(hits))
	for i, hit := range hits {
		matches[i] = p.rectFromPageRect(hit)
	}
	return matches
}

func (p *PDF) pointFromPagePoint(x, y int) geom32.Point {
	return geom32.NewPoint(float32(x)*p.DisplayScaleAdjust, float32(y)*p.DisplayScaleAdjust)
}

func (p *PDF) rectFromPageRect(r image.Rectangle) geom32.Rect {
	return geom32.NewRect(float32(r.Min.X)*p.DisplayScaleAdjust, float32(r.Min.Y)*p.DisplayScaleAdjust, float32(r.Dx())*p.DisplayScaleAdjust,
		float32(r.Dy())*p.DisplayScaleAdjust)
}