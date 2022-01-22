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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/internal/pdf"
	"github.com/richardwilkes/gcs/internal/ui/icons"
	"github.com/richardwilkes/gcs/internal/ui/search"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const (
	minPDFDockableScale   = 25
	maxPDFDockableScale   = 300
	deltaPDFDockableScale = 10
)

var (
	_                 FileBackedDockable = &PDFDockable{}
	_                 unison.TabCloser   = &PDFDockable{}
	pdfMatchHighlight *unison.Paint
	pdfLinkHighlight  *unison.Paint
)

// PDFDockable holds the view for a PDF file.
type PDFDockable struct {
	unison.Panel
	path               string
	pdf                *pdf.PDF
	scroll             *unison.ScrollPanel
	docPanel           *unison.Panel
	pageNumberField    *unison.Field
	scaleField         *unison.Field
	searchField        *unison.Field
	matchesLabel       *unison.Label
	firstPageButton    *unison.Button
	previousPageButton *unison.Button
	nextPageButton     *unison.Button
	lastPageButton     *unison.Button
	page               *pdf.Page
	link               *pdf.Link
	rolloverRect       geom32.Rect
	scale              int
}

// NewPDFDockable creates a new FileBackedDockable for PDF files.
func NewPDFDockable(filePath string) (*PDFDockable, error) {
	d := &PDFDockable{
		path:  filePath,
		scale: 100,
	}
	d.Self = d
	var err error
	if d.pdf, err = pdf.New(filePath, func() {
		unison.InvokeTask(d.pageLoaded)
	}); err != nil {
		return nil, err
	}
	d.KeyDownCallback = d.keyDown
	d.SetLayout(&unison.FlexLayout{Columns: 1})

	d.docPanel = unison.NewPanel()
	d.docPanel.SetSizer(d.docSizer)
	d.docPanel.DrawCallback = d.draw
	d.docPanel.MouseDownCallback = d.mouseDown
	d.docPanel.MouseMoveCallback = d.mouseMove
	d.docPanel.MouseUpCallback = d.mouseUp
	d.docPanel.SetFocusable(true)

	d.scroll = unison.NewScrollPanel()
	d.scroll.MouseWheelMultiplier = 4
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.SetContent(d.docPanel, unison.FillBehavior)

	d.firstPageButton = icons.NewIconButton(icons.FirstSVG(), 16)
	d.firstPageButton.ClickCallback = func() { d.loadPage(0) }

	d.previousPageButton = icons.NewIconButton(icons.PreviousSVG(), 16)
	d.previousPageButton.ClickCallback = func() { d.loadPage(d.pdf.MostRecentPageNumber() - 1) }

	d.nextPageButton = icons.NewIconButton(icons.NextSVG(), 16)
	d.nextPageButton.ClickCallback = func() { d.loadPage(d.pdf.MostRecentPageNumber() + 1) }

	d.lastPageButton = icons.NewIconButton(icons.LastSVG(), 16)
	d.lastPageButton.ClickCallback = func() { d.loadPage(d.pdf.PageCount() - 1) }

	pageLabel := unison.NewLabel()
	pageLabel.Font = unison.DefaultFieldTheme.Font
	pageLabel.Text = i18n.Text("Page")

	d.pageNumberField = unison.NewField()
	d.pageNumberField.MinimumTextWidth = d.pageNumberField.Font.Width(strconv.Itoa(d.pdf.PageCount() * 10))
	d.pageNumberField.ModifiedCallback = func() {
		if pageNum, e := strconv.Atoi(d.pageNumberField.Text()); e == nil && pageNum > 0 && pageNum <= d.pdf.PageCount() {
			d.loadPage(pageNum - 1)
		}
	}
	d.pageNumberField.ValidateCallback = func() bool {
		pageNum, e := strconv.Atoi(d.pageNumberField.Text())
		if e != nil || pageNum < 1 || pageNum > d.pdf.PageCount() {
			return false
		}
		return true
	}

	ofLabel := unison.NewLabel()
	ofLabel.Font = unison.DefaultFieldTheme.Font
	ofLabel.Text = fmt.Sprintf(i18n.Text("of %d"), d.pdf.PageCount())

	d.scaleField = unison.NewField()
	d.scaleField.MinimumTextWidth = d.scaleField.Font.Width(strconv.Itoa(maxPDFDockableScale) + "%")
	d.scaleField.SetText(strconv.Itoa(d.scale) + "%")
	d.scaleField.ModifiedCallback = func() {
		if s, e := strconv.Atoi(strings.TrimRight(d.scaleField.Text(), "%")); e == nil && s >= minPDFDockableScale && s <= maxPDFDockableScale {
			d.scale = s
			d.loadPage(d.pdf.MostRecentPageNumber())
		}
	}
	d.scaleField.ValidateCallback = func() bool {
		if s, e := strconv.Atoi(strings.TrimRight(d.scaleField.Text(), "%")); e != nil || s < minPDFDockableScale || s > maxPDFDockableScale {
			return false
		}
		return true
	}

	d.searchField = search.NewField()
	d.searchField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	existingCallback := d.searchField.ModifiedCallback
	d.searchField.ModifiedCallback = func() {
		d.loadPage(d.pdf.MostRecentPageNumber())
		existingCallback()
	}

	d.matchesLabel = unison.NewLabel()
	d.matchesLabel.Text = "-"
	d.matchesLabel.Tooltip = unison.NewTooltipWithText(i18n.Text("Number of matches found"))

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(d.firstPageButton)
	toolbar.AddChild(d.previousPageButton)
	toolbar.AddChild(d.nextPageButton)
	toolbar.AddChild(d.lastPageButton)
	toolbar.AddChild(unison.NewPanel())
	toolbar.AddChild(pageLabel)
	toolbar.AddChild(d.pageNumberField)
	toolbar.AddChild(ofLabel)
	toolbar.AddChild(unison.NewPanel())
	toolbar.AddChild(d.scaleField)
	toolbar.AddChild(unison.NewPanel())
	toolbar.AddChild(d.searchField)
	toolbar.AddChild(d.matchesLabel)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	d.AddChild(toolbar)
	d.AddChild(d.scroll)

	d.loadPage(0)

	return d, nil
}

func (d *PDFDockable) loadPage(pageNumber int) {
	d.pdf.LoadPage(pageNumber, float32(d.scale)/100, d.searchField.Text())
	pageNumber = d.pdf.MostRecentPageNumber()
	lastPageNumber := d.pdf.PageCount() - 1
	d.firstPageButton.SetEnabled(pageNumber != 0)
	d.previousPageButton.SetEnabled(pageNumber > 0)
	d.nextPageButton.SetEnabled(pageNumber < lastPageNumber)
	d.lastPageButton.SetEnabled(pageNumber != lastPageNumber)
}

func (d *PDFDockable) pageLoaded() {
	d.page = d.pdf.CurrentPage()

	pageText := ""
	if d.page.PageNumber >= 0 {
		pageText = strconv.Itoa(d.page.PageNumber + 1)
	}
	if pageText != d.pageNumberField.Text() {
		d.pageNumberField.SetText(pageText)
		d.pageNumberField.Parent().MarkForLayoutAndRedraw()
	}

	scaleText := strconv.Itoa(d.scale) + "%"
	if scaleText != d.scaleField.Text() {
		d.scaleField.SetText(scaleText)
		d.scaleField.Parent().MarkForLayoutAndRedraw()
	}

	matchText := "-"
	if d.searchField.Text() != "" {
		matchText = strconv.Itoa(len(d.page.Matches))
	}
	if matchText != d.matchesLabel.Text {
		d.matchesLabel.Text = matchText
		d.matchesLabel.Parent().MarkForLayoutAndRedraw()
	}

	d.docPanel.MarkForLayoutAndRedraw()
	d.scroll.MarkForLayoutAndRedraw()
	d.link = nil
}

func (d *PDFDockable) overLink(where geom32.Point) (rect geom32.Rect, link *pdf.Link) {
	if d.page != nil && d.page.Links != nil {
		for _, link = range d.page.Links {
			if link.Bounds.ContainsPoint(where) {
				return link.Bounds, link
			}
		}
	}
	return rect, nil
}

func (d *PDFDockable) checkForLinkAt(where geom32.Point) {
	r, link := d.overLink(where)
	if r != d.rolloverRect || link != d.link {
		d.rolloverRect = r
		d.link = link
		d.MarkForRedraw()
	}
}

func (d *PDFDockable) mouseDown(_ geom32.Point, _, _ int, _ unison.Modifiers) bool {
	d.RequestFocus()
	return true
}

func (d *PDFDockable) mouseMove(where geom32.Point, _ unison.Modifiers) bool {
	d.checkForLinkAt(where)
	return true
}

func (d *PDFDockable) mouseUp(where geom32.Point, button int, _ unison.Modifiers) bool {
	d.checkForLinkAt(where)
	if button == unison.ButtonLeft && d.link != nil {
		if d.link.PageNumber >= 0 {
			d.loadPage(d.link.PageNumber)
			// TODO: Use d.link.PageX & PageY to ensure location is scrolled into place
		} else if err := desktop.OpenBrowser(d.link.URI); err != nil {
			unison.ErrorDialogWithError(i18n.Text("Unable to open link"), err)
		}
	}
	return true
}

func (d *PDFDockable) keyDown(keyCode unison.KeyCode, _ unison.Modifiers, _ bool) bool {
	scale := d.scale
	switch keyCode {
	case unison.KeyQ:
		scale = 25
	case unison.KeyH:
		scale = 50
	case unison.KeyT:
		scale = 75
	case unison.Key1:
		scale = 100
	case unison.Key2:
		scale = 200
	case unison.Key3:
		scale = 300
	case unison.KeyMinus:
		scale -= deltaPDFDockableScale
		if scale < minPDFDockableScale {
			scale = minPDFDockableScale
		}
	case unison.KeyEqual:
		scale += deltaPDFDockableScale
		if scale > maxPDFDockableScale {
			scale = maxPDFDockableScale
		}
	case unison.KeyHome:
		d.loadPage(0)
	case unison.KeyEnd:
		d.loadPage(d.pdf.PageCount() - 1)
	case unison.KeyLeft, unison.KeyUp:
		d.loadPage(d.pdf.MostRecentPageNumber() - 1)
	case unison.KeyRight, unison.KeyDown:
		d.loadPage(d.pdf.MostRecentPageNumber() + 1)
	default:
		return false
	}
	if d.scale != scale {
		d.scale = scale
		f := d.scaleField.ModifiedCallback
		d.scaleField.ModifiedCallback = nil
		d.scaleField.SetText(strconv.Itoa(d.scale) + "%")
		d.scaleField.ModifiedCallback = f
		d.loadPage(d.pdf.MostRecentPageNumber())
	}
	return true
}

func (d *PDFDockable) docSizer(_ geom32.Size) (min, pref, max geom32.Size) {
	if d.page == nil || d.page.Error != nil {
		pref.Width = 400
		pref.Height = 300
	} else {
		pref = d.page.Image.LogicalSize()
	}
	return geom32.NewSize(50, 50), pref, unison.MaxSize(pref)
}

func (d *PDFDockable) draw(gc *unison.Canvas, dirty geom32.Rect) {
	gc.DrawRect(dirty, unison.ContentColor.Paint(gc, dirty, unison.Fill))
	if d.page == nil {
		return
	}
	if d.page.Error != nil {
		r := d.docPanel.ContentRect(false)
		r.Inset(geom32.NewUniformInsets(unison.StdHSpacing))
		unison.DrawLabel(gc, r, unison.MiddleAlignment, unison.StartAlignment, fmt.Sprintf("%s", d.page.Error), //nolint:gocritic // Want the special handling %s provides
			unison.SystemFont, unison.OnContentColor, unison.DefaultDialogTheme.ErrorIcon, unison.LeftSide,
			unison.StdHSpacing, false)
		return
	}
	r := geom32.Rect{Size: d.page.Image.LogicalSize()}
	gc.DrawRect(r, unison.White.Paint(gc, r, unison.Fill))
	gc.DrawImageInRect(d.page.Image, r, nil, nil)
	for _, match := range d.page.Matches {
		gc.DrawRect(match, getPDFMatchHighlightPaint())
	}
	if d.link != nil {
		gc.DrawRect(d.rolloverRect, getPDFLinkHighlightPaint())
	}
}

// TitleIcon implements FileBackedDockable
func (d *PDFDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements FileBackedDockable
func (d *PDFDockable) Title() string {
	return xfs.BaseName(d.path)
}

// Tooltip implements FileBackedDockable
func (d *PDFDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements FileBackedDockable
func (d *PDFDockable) BackingFilePath() string {
	return d.path
}

// Modified implements FileBackedDockable
func (d *PDFDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *PDFDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *PDFDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func getPDFMatchHighlightPaint() *unison.Paint {
	if pdfMatchHighlight == nil {
		pdfMatchHighlight = unison.NewPaint()
		pdfMatchHighlight.SetStyle(unison.Fill)
		pdfMatchHighlight.SetBlendMode(unison.ModulateBlendMode)
		pdfMatchHighlight.SetColor(unison.Yellow)
	}
	return pdfMatchHighlight
}

func getPDFLinkHighlightPaint() *unison.Paint {
	if pdfLinkHighlight == nil {
		pdfLinkHighlight = unison.NewPaint()
		pdfLinkHighlight.SetStyle(unison.Fill)
		pdfLinkHighlight.SetBlendMode(unison.ModulateBlendMode)
		pdfLinkHighlight.SetColor(unison.GreenYellow)
	}
	return pdfLinkHighlight
}
