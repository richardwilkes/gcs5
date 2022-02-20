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

package pageref

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const extension = ".refs"

type wndData struct {
	wnd         *unison.Window
	resetButton *unison.Button
	menuButton  *unison.Button
	content     *unison.Panel
}

var data *wndData

// Show the Page Reference Mappings window.
func Show() {
	if data == nil {
		wnd, err := unison.NewWindow(i18n.Text("Page Reference Mappings"))
		if err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to open Page Reference Mappings"), err.Error())
			return
		}
		wnd.WillCloseCallback = func() { data = nil }
		data = &wndData{
			wnd: wnd,
		}
		content := data.wnd.Content()
		content.SetLayout(&unison.FlexLayout{Columns: 1})
		content.AddChild(data.createToolbar())
		content.AddChild(data.createContent())
		data.wnd.Pack()
		data.content.RequestFocus()
	}
	data.wnd.ToFront()
}

func (d *wndData) createToolbar() *unison.Panel {
	toolbar := unison.NewPanel()
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	spacer := unison.NewPanel()
	spacer.SetLayoutData(&unison.FlexLayoutData{HGrab: true})
	toolbar.AddChild(spacer)
	d.resetButton = unison.NewSVGButton(icons.ResetSVG())
	d.resetButton.ClickCallback = d.reset
	toolbar.AddChild(d.resetButton)
	d.menuButton = unison.NewSVGButton(icons.MenuSVG())
	d.menuButton.ClickCallback = d.showMenu
	toolbar.AddChild(d.menuButton)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})
	return toolbar
}

func (d *wndData) createContent() unison.Paneler {
	d.content = unison.NewPanel()
	d.content.SetLayout(&unison.FlexLayout{
		Columns:  4,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.content.SetBorder(unison.NewEmptyBorder(geom32.NewUniformInsets(unison.StdHSpacing * 2)))
	d.sync()

	scroller := unison.NewScrollPanel()
	//scroller.SetColumnHeader(header)
	scroller.SetContent(d.content, unison.FillBehavior)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	return scroller
}

func (d *wndData) reset() {
	settings.Global().PageRefs = settings.PageRefs{}
	d.sync()
}

func (d *wndData) sync() {
	d.content.RemoveAllChildren()
	for _, one := range settings.Global().PageRefs.List() {
		d.createIDField(one)
		d.createOffsetField(one)
		d.createNameField(one)
		d.createTrashField(one)
	}
	d.wnd.MarkForRedraw()
}

func (d *wndData) createIDField(ref *settings.PageRef) {
	p := unison.NewLabel()
	p.Text = ref.ID
	p.HAlign = unison.MiddleAlignment
	p.OnBackgroundInk = unison.DefaultTooltipTheme.Label.OnBackgroundInk
	p.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.NewUniformInsets(1), false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    1,
			Left:   unison.StdHSpacing,
			Bottom: 1,
			Right:  unison.StdHSpacing,
		})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.DefaultTooltipTheme.BackgroundInk.Paint(gc, rect, unison.Fill))
		p.DefaultDraw(gc, rect)
	}
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(p)
}

func (d *wndData) createOffsetField(ref *settings.PageRef) {
	p := widget.NewSignedIntegerField(ref.Offset, -9999, 9999, func(v int) {
		ref.Offset = v
		settings.Global().PageRefs.Set(ref)
	})
	p.Tooltip = unison.NewTooltipWithText(i18n.Text(`If your PDF is opening up to the wrong page when opening
page references, enter an offset here to compensate.`))
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(p)
}

func (d *wndData) createNameField(ref *settings.PageRef) {
	p := unison.NewLabel()
	p.Text = filepath.Base(ref.Path)
	p.Tooltip = unison.NewTooltipWithText(ref.Path)
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.content.AddChild(p)
}

func (d *wndData) createTrashField(ref *settings.PageRef) {
	b := unison.NewSVGButton(icons.TrashSVG())
	b.ClickCallback = func() {
		settings.Global().PageRefs.Remove(ref.ID)
		parent := b.Parent()
		index := parent.IndexOfChild(b)
		for i := index; i > index-4; i-- {
			parent.RemoveChildAtIndex(i)
		}
		parent.MarkForLayoutAndRedraw()
	}
	b.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.MiddleAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(b)
}

func (d *wndData) showMenu() {
	f := unison.DefaultMenuFactory()
	id := unison.ContextMenuIDFlag
	m := f.NewMenu(id, "", nil)
	id++
	m.InsertItem(-1, f.NewItem(id, i18n.Text("Import…"), 0, 0, nil, d.handleImport))
	id++
	m.InsertItem(-1, f.NewItem(id, i18n.Text("Export…"), 0, 0, nil, d.handleExport))
	id++
	libraries := settings.Global().Libraries()
	sets := library.ScanForNamedFileSets(nil, "", extension, false, libraries)
	if len(sets) != 0 {
		m.InsertSeparator(-1, false)
		for _, lib := range sets {
			m.InsertItem(-1, f.NewItem(id, lib.Name, 0, 0, func(_ unison.MenuItem) bool { return false }, nil))
			id++
			for _, one := range lib.List {
				d.insertFileToLoad(m, id, one)
				id++
			}
		}
	}
	m.Popup(d.menuButton.RectToRoot(d.menuButton.ContentRect(true)), 0)
}

func (d *wndData) insertFileToLoad(m unison.Menu, id int, ref *library.NamedFileRef) {
	m.InsertItem(-1, m.Factory().NewItem(id, ref.Name, 0, 0, nil, func(_ unison.MenuItem) {
		d.load(ref.FileSystem, ref.FilePath)
	}))
}

func (d *wndData) load(fileSystem fs.FS, filePath string) {
	s, err := settings.NewPageRefsFromFS(fileSystem, filePath)
	if err != nil {
		unison.ErrorDialogWithMessage(i18n.Text("Unable to load page reference mappings"), err.Error())
		return
	}
	settings.Global().PageRefs = *s
	d.sync()
}

func (d *wndData) handleImport(_ unison.MenuItem) {
	dialog := unison.NewOpenDialog()
	dialog.SetResolvesAliases(true)
	dialog.SetAllowedExtensions(extension)
	dialog.SetAllowsMultipleSelection(false)
	dialog.SetCanChooseDirectories(false)
	dialog.SetCanChooseFiles(true)
	if dialog.RunModal() {
		p := dialog.Path()
		d.load(os.DirFS(filepath.Dir(p)), filepath.Base(p))
	}
}

func (d *wndData) handleExport(_ unison.MenuItem) {
	dialog := unison.NewSaveDialog()
	dialog.SetAllowedExtensions(extension)
	if dialog.RunModal() {
		if err := settings.Global().PageRefs.Save(dialog.Path()); err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to save page reference mappings"), err.Error())
		}
	}
}
