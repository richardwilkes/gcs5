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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var (
	_ unison.Dockable  = &SettingsDockable{}
	_ unison.TabCloser = &SettingsDockable{}
)

// SettingsDockable holds common settings dockable data.
type SettingsDockable struct {
	unison.Panel
	TabTitle  string
	Extension string
	Loader    func(fileSystem fs.FS, filePath string) error
	Saver     func(filePath string) error
	Resetter  func()
}

// Activate attempts to locate an existing dockable that 'matcher' returns true for. If found, it will have been
// activated and focused.
func Activate(matcher func(d unison.Dockable) bool) (ws *Workspace, dc *unison.DockContainer, found bool) {
	if ws = Any(); ws == nil {
		jot.Error("no workspace available")
		return nil, nil, false
	}
	if focus := ws.Window.Focus(); focus != nil {
		if focusedDC := unison.DockContainerFor(focus); focusedDC != nil && focusedDC.Dock == ws.DocumentDock.Dock {
			dc = focusedDC
		}
	}
	ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(container *unison.DockContainer) bool {
		for _, one := range container.Dockables() {
			if matcher(one) {
				found = true
				container.SetCurrentDockable(one)
				container.AcquireFocus()
				return true
			}
			if dc == nil {
				dc = container
			}
		}
		return false
	})
	return ws, dc, found
}

// Setup the dockable and display it.
func (d *SettingsDockable) Setup(ws *Workspace, dc *unison.DockContainer, addToStartToolbar, addToEndToolbar, initContent func(*unison.Panel)) {
	d.SetLayout(&unison.FlexLayout{Columns: 1})
	d.AddChild(d.createToolbar(addToStartToolbar, addToEndToolbar))
	content := unison.NewPanel()
	content.SetBorder(unison.NewEmptyBorder(geom32.NewUniformInsets(unison.StdHSpacing * 2)))
	initContent(content)
	scroller := unison.NewScrollPanel()
	scroller.SetContent(content, unison.FillBehavior)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	d.AddChild(scroller)
	if dc != nil {
		dc.Stack(d, -1)
	} else {
		ws.DocumentDock.DockTo(d, nil, unison.LeftSide)
	}
}

// TitleIcon implements unison.Dockable
func (d *SettingsDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  icons.SettingsSVG(),
		Size: suggestedSize,
	}
}

// Title implements unison.Dockable
func (d *SettingsDockable) Title() string {
	return d.TabTitle
}

// Tooltip implements unison.Dockable
func (d *SettingsDockable) Tooltip() string {
	return ""
}

// Modified implements unison.Dockable
func (d *SettingsDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *SettingsDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *SettingsDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func (d *SettingsDockable) createToolbar(addToStartToolbar, addToEndToolbar func(*unison.Panel)) *unison.Panel {
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
	if addToStartToolbar != nil {
		addToStartToolbar(toolbar)
	}
	index := len(toolbar.Children())
	if addToEndToolbar != nil {
		addToEndToolbar(toolbar)
	}
	if d.Resetter != nil {
		b := unison.NewSVGButton(icons.ResetSVG())
		b.Tooltip = unison.NewTooltipWithText(i18n.Text("Reset"))
		b.ClickCallback = d.handleReset
		toolbar.AddChild(b)
	}
	if d.Loader != nil || d.Saver != nil {
		b := unison.NewSVGButton(icons.MenuSVG())
		b.ClickCallback = func() { d.showMenu(b) }
		toolbar.AddChild(b)
	}
	if len(toolbar.Children()) != index {
		spacer := unison.NewPanel()
		spacer.SetLayoutData(&unison.FlexLayoutData{HGrab: true})
		toolbar.AddChildAtIndex(spacer, index)
	}
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})
	return toolbar
}

func (d *SettingsDockable) handleReset() {
	if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Are you sure you want to reset the\n%s?"), d.TabTitle), "") == unison.ModalResponseOK {
		d.Resetter()
	}
}

func (d *SettingsDockable) showMenu(b *unison.Button) {
	f := unison.DefaultMenuFactory()
	id := unison.ContextMenuIDFlag
	m := f.NewMenu(id, "", nil)
	id++
	if d.Loader != nil {
		m.InsertItem(-1, f.NewItem(id, i18n.Text("Import…"), 0, 0, nil, d.handleImport))
		id++
	}
	if d.Saver != nil {
		m.InsertItem(-1, f.NewItem(id, i18n.Text("Export…"), 0, 0, nil, d.handleExport))
		id++
	}
	if d.Loader != nil {
		libraries := settings.Global().Libraries()
		sets := library.ScanForNamedFileSets(nil, "", d.Extension, false, libraries)
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
	}
	m.Popup(b.RectToRoot(b.ContentRect(true)), 0)
}

func (d *SettingsDockable) insertFileToLoad(m unison.Menu, id int, ref *library.NamedFileRef) {
	m.InsertItem(-1, m.Factory().NewItem(id, ref.Name, 0, 0, nil, func(_ unison.MenuItem) {
		d.doLoad(ref.FileSystem, ref.FilePath)
	}))
}

func (d *SettingsDockable) doLoad(fileSystem fs.FS, filePath string) {
	if err := d.Loader(fileSystem, filePath); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to load ")+d.TabTitle, err)
	}
}

func (d *SettingsDockable) handleImport(_ unison.MenuItem) {
	dialog := unison.NewOpenDialog()
	dialog.SetResolvesAliases(true)
	dialog.SetAllowedExtensions(d.Extension)
	dialog.SetAllowsMultipleSelection(false)
	dialog.SetCanChooseDirectories(false)
	dialog.SetCanChooseFiles(true)
	if dialog.RunModal() {
		p := dialog.Path()
		d.doLoad(os.DirFS(filepath.Dir(p)), filepath.Base(p))
	}
}

func (d *SettingsDockable) handleExport(_ unison.MenuItem) {
	dialog := unison.NewSaveDialog()
	dialog.SetAllowedExtensions(d.Extension)
	if dialog.RunModal() {
		if err := d.Saver(dialog.Path()); err != nil {
			unison.ErrorDialogWithError(i18n.Text("Unable to save ")+d.TabTitle, err)
		}
	}
}
