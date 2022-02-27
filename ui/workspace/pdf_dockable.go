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
	"time"

	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/pdf"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const (
	minPDFDockableScale                = 25
	maxPDFDockableScale                = 300
	deltaPDFDockableScale              = 10
	maxElapsedRenderTimeWithoutOverlay = time.Millisecond * 250
	renderTimeSlop                     = time.Millisecond * 10
)

var (
	_ FileBackedDockable = &PDFDockable{}
	_ unison.TabCloser   = &PDFDockable{}
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
	backButton         *unison.Button
	forwardButton      *unison.Button
	firstPageButton    *unison.Button
	previousPageButton *unison.Button
	nextPageButton     *unison.Button
	lastPageButton     *unison.Button
	page               *pdf.Page
	link               *pdf.Link
	rolloverRect       geom32.Rect
	scale              int
	historyPos         int
	history            []int
	noUpdate           bool
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
	d.FocusChangeInHierarchyCallback = d.focusChangeInHierarchy
	d.GainedFocusCallback = d.pdf.RequestRenderPriority
	d.SetLayout(&unison.FlexLayout{Columns: 1})

	d.docPanel = unison.NewPanel()
	d.docPanel.SetSizer(d.docSizer)
	d.docPanel.DrawCallback = d.draw
	d.docPanel.MouseDownCallback = d.mouseDown
	d.docPanel.MouseMoveCallback = d.mouseMove
	d.docPanel.MouseUpCallback = d.mouseUp
	d.docPanel.SetFocusable(true)

	d.scroll = unison.NewScrollPanel()
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.scroll.SetContent(d.docPanel, unison.FillBehavior)
	d.scroll.ContentView().DrawOverCallback = d.drawOverlay

	d.backButton = unison.NewSVGButton(icons.BackSVG())
	d.backButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Back"))
	d.backButton.ClickCallback = func() { d.Back() }

	d.forwardButton = unison.NewSVGButton(icons.ForwardSVG())
	d.forwardButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Forward"))
	d.forwardButton.ClickCallback = func() { d.Forward() }

	d.firstPageButton = unison.NewSVGButton(icons.FirstSVG())
	d.firstPageButton.Tooltip = unison.NewTooltipWithText(i18n.Text("First Page"))
	d.firstPageButton.ClickCallback = func() { d.LoadPage(0) }

	d.previousPageButton = unison.NewSVGButton(icons.PreviousSVG())
	d.previousPageButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Previous Page"))
	d.previousPageButton.ClickCallback = func() { d.LoadPage(d.pdf.MostRecentPageNumber() - 1) }

	d.nextPageButton = unison.NewSVGButton(icons.NextSVG())
	d.nextPageButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Next Page"))
	d.nextPageButton.ClickCallback = func() { d.LoadPage(d.pdf.MostRecentPageNumber() + 1) }

	d.lastPageButton = unison.NewSVGButton(icons.LastSVG())
	d.lastPageButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Last Page"))
	d.lastPageButton.ClickCallback = func() { d.LoadPage(d.pdf.PageCount() - 1) }

	pageLabel := unison.NewLabel()
	pageLabel.Font = unison.DefaultFieldTheme.Font
	pageLabel.Text = i18n.Text("Page")

	d.pageNumberField = unison.NewField()
	d.pageNumberField.MinimumTextWidth = d.pageNumberField.Font.SimpleWidth(strconv.Itoa(d.pdf.PageCount() * 10))
	d.pageNumberField.ModifiedCallback = func() {
		if d.noUpdate {
			return
		}
		if pageNum, e := strconv.Atoi(d.pageNumberField.Text()); e == nil && pageNum > 0 && pageNum <= d.pdf.PageCount() {
			d.LoadPage(pageNum - 1)
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
	d.scaleField.Tooltip = unison.NewTooltipWithText(i18n.Text("Scale"))
	d.scaleField.MinimumTextWidth = d.scaleField.Font.SimpleWidth(strconv.Itoa(maxPDFDockableScale) + "%")
	d.scaleField.SetText(strconv.Itoa(d.scale) + "%")
	d.scaleField.ModifiedCallback = func() {
		if d.noUpdate {
			return
		}
		if s, e := strconv.Atoi(strings.TrimRight(d.scaleField.Text(), "%")); e == nil && s >= minPDFDockableScale && s <= maxPDFDockableScale {
			d.scale = s
			d.LoadPage(d.pdf.MostRecentPageNumber())
		}
	}
	d.scaleField.ValidateCallback = func() bool {
		if s, e := strconv.Atoi(strings.TrimRight(d.scaleField.Text(), "%")); e != nil || s < minPDFDockableScale || s > maxPDFDockableScale {
			return false
		}
		return true
	}

	d.searchField = widget.NewSearchField()
	pageSearch := i18n.Text("Page Search")
	d.searchField.Watermark = pageSearch
	d.searchField.Tooltip = unison.NewTooltipWithText(pageSearch)
	d.searchField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	existingCallback := d.searchField.ModifiedCallback
	d.searchField.ModifiedCallback = func() {
		if d.noUpdate {
			return
		}
		d.LoadPage(d.pdf.MostRecentPageNumber())
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
	toolbar.AddChild(d.backButton)
	toolbar.AddChild(d.forwardButton)
	toolbar.AddChild(unison.NewPanel())
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

	d.LoadPage(0)

	return d, nil
}

// ClearHistory clears the existing history.
func (d *PDFDockable) ClearHistory() {
	d.history = nil
	d.historyPos = 0
	d.backButton.SetEnabled(false)
	d.forwardButton.SetEnabled(false)
}

// SetSearchText sets the search text and updates the display.
func (d *PDFDockable) SetSearchText(text string) {
	d.searchField.SetText(text)
}

// Back moves back in history one step.
func (d *PDFDockable) Back() {
	if d.historyPos > 0 {
		d.historyPos--
		d.LoadPage(d.history[d.historyPos])
	}
}

// Forward moves forward in history one step.
func (d *PDFDockable) Forward() {
	if d.historyPos < len(d.history)-1 {
		d.historyPos++
		d.LoadPage(d.history[d.historyPos])
	}
}

// LoadPage loads the specified page.
func (d *PDFDockable) LoadPage(pageNumber int) {
	d.pdf.LoadPage(pageNumber, float32(d.scale)/100, d.searchField.Text())
	d.MarkForRedraw()
}

func (d *PDFDockable) pageLoaded() {
	d.noUpdate = true
	defer func() { d.noUpdate = false }()

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

	pageNumber := d.page.PageNumber
	if d.history == nil {
		d.history = append(d.history, pageNumber)
		d.historyPos = 0
	} else if d.history[d.historyPos] != pageNumber {
		d.historyPos++
		if d.historyPos < len(d.history) {
			if d.history[d.historyPos] != pageNumber {
				d.history[d.historyPos] = pageNumber
				d.history = d.history[:d.historyPos+1]
			}
		} else {
			d.history = append(d.history, pageNumber)
		}
	}
	lastPageNumber := d.pdf.PageCount() - 1
	d.backButton.SetEnabled(d.historyPos > 0)
	d.forwardButton.SetEnabled(d.historyPos < len(d.history)-1)
	d.firstPageButton.SetEnabled(pageNumber != 0)
	d.previousPageButton.SetEnabled(pageNumber > 0)
	d.nextPageButton.SetEnabled(pageNumber < lastPageNumber)
	d.lastPageButton.SetEnabled(pageNumber != lastPageNumber)

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
			d.LoadPage(d.link.PageNumber)
			// TODO: Use d.link.PageX & PageY to ensure location is scrolled into place
		} else if err := desktop.OpenBrowser(d.link.URI); err != nil {
			unison.ErrorDialogWithError(i18n.Text("Unable to open link"), err)
		}
	}
	return true
}

func (d *PDFDockable) focusChangeInHierarchy(_, _ *unison.Panel) {
	d.pdf.RequestRenderPriority()
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
		d.LoadPage(0)
	case unison.KeyEnd:
		d.LoadPage(d.pdf.PageCount() - 1)
	case unison.KeyLeft, unison.KeyUp:
		d.LoadPage(d.pdf.MostRecentPageNumber() - 1)
	case unison.KeyRight, unison.KeyDown:
		d.LoadPage(d.pdf.MostRecentPageNumber() + 1)
	default:
		return false
	}
	if d.scale != scale {
		d.scale = scale
		f := d.scaleField.ModifiedCallback
		d.scaleField.ModifiedCallback = nil
		d.scaleField.SetText(strconv.Itoa(d.scale) + "%")
		d.scaleField.ModifiedCallback = f
		d.LoadPage(d.pdf.MostRecentPageNumber())
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
	if d.page.Image != nil {
		r := geom32.Rect{Size: d.page.Image.LogicalSize()}
		gc.DrawRect(r, unison.White.Paint(gc, r, unison.Fill))
		gc.DrawImageInRect(d.page.Image, r, nil, nil)
		if len(d.page.Matches) != 0 {
			p := unison.NewPaint()
			p.SetStyle(unison.Fill)
			p.SetBlendMode(unison.ModulateBlendMode)
			p.SetColor(theme.PDFMarkerHighlightColor.GetColor())
			for _, match := range d.page.Matches {
				gc.DrawRect(match, p)
			}
		}
		if d.link != nil {
			p := unison.NewPaint()
			p.SetStyle(unison.Fill)
			p.SetBlendMode(unison.ModulateBlendMode)
			p.SetColor(theme.PDFLinkHighlightColor.GetColor())
			gc.DrawRect(d.rolloverRect, p)
		}
	}
}

func (d *PDFDockable) drawOverlay(gc *unison.Canvas, dirty geom32.Rect) {
	if d.page != nil && d.page.Error != nil {
		d.drawOverlayMsg(gc, dirty, fmt.Sprintf("%s", d.page.Error), true) //nolint:gocritic // I want the extra processing %s does in this case
	}
	if finished, pageNumber, requested := d.pdf.RenderingFinished(); !finished {
		if waitFor := maxElapsedRenderTimeWithoutOverlay - time.Since(requested); waitFor > renderTimeSlop {
			unison.InvokeTaskAfter(d.MarkForRedraw, waitFor)
		} else {
			d.drawOverlayMsg(gc, dirty, fmt.Sprintf(i18n.Text("Rendering page %d…"), pageNumber+1), false)
		}
	}
}

func (d *PDFDockable) drawOverlayMsg(gc *unison.Canvas, dirty geom32.Rect, msg string, forError bool) {
	var fgInk, bgInk unison.Ink
	var icon unison.Drawable
	font := unison.SystemFont.Face().Font(24)
	baseline := font.Baseline()
	if forError {
		fgInk = unison.OnErrorColor
		bgInk = unison.ErrorColor.GetColor().SetAlphaIntensity(0.7)
		icon = &unison.DrawableSVG{
			SVG:  unison.CircledExclamationSVG(),
			Size: geom32.NewSize(baseline, baseline),
		}
	} else {
		fgInk = unison.OnContentColor
		bgInk = unison.ContentColor.GetColor().SetAlphaIntensity(0.7)
	}
	decoration := &unison.TextDecoration{
		Font:  font,
		Paint: fgInk.Paint(gc, dirty, unison.Fill),
	}
	text := unison.NewText(msg, decoration)
	r := d.scroll.ContentView().ContentRect(false)
	cy := r.CenterY()
	width := text.Width()
	height := text.Height()
	var iconSize geom32.Size
	if icon != nil {
		iconSize = icon.LogicalSize()
		width += iconSize.Width + unison.StdHSpacing
		if height < iconSize.Height {
			height = iconSize.Height
		}
	}
	backWidth := width + 40
	backHeight := height + 40
	r.X += (r.Width - backWidth) / 2
	if forError {
		r.Y = cy - (backHeight + unison.StdVSpacing)
	} else {
		r.Y = cy + unison.StdVSpacing
	}
	r.Width = backWidth
	r.Height = backHeight
	gc.DrawRoundedRect(r, 10, 10, bgInk.Paint(gc, dirty, unison.Fill))
	x := r.X + (r.Width-width)/2
	if icon != nil {
		icon.DrawInRect(gc, geom32.NewRect(x, r.Y+(r.Height-iconSize.Height)/2, iconSize.Width, iconSize.Height), nil, decoration.Paint)
		x += iconSize.Width + unison.StdHSpacing
	}
	text.Draw(gc, x, r.Y+(r.Height-height)/2+baseline)
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
	return fs.BaseName(d.path)
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
