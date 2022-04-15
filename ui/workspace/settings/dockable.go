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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

const settingsGroup = "settings"

var (
	_ unison.Dockable  = &Dockable{}
	_ unison.TabCloser = &Dockable{}
)

// Dockable holds common settings dockable data.
type Dockable struct {
	unison.Panel
	TabTitle  string
	Extension string
	Loader    func(fileSystem fs.FS, filePath string) error
	Saver     func(filePath string) error
	Resetter  func()
}

// ActiveDockable returns the currently active dockable in the active window.
func ActiveDockable() unison.Dockable {
	ws := workspace.FromWindow(unison.ActiveWindow())
	if ws == nil {
		return nil
	}
	dc := ws.CurrentlyFocusedDockContainer()
	if dc == nil {
		return nil
	}
	return dc.CurrentDockable()
}

// Activate attempts to locate an existing dockable that 'matcher' returns true for. If found, it will have been
// activated and focused.
func Activate(matcher func(d unison.Dockable) bool) (ws *workspace.Workspace, dc *unison.DockContainer, found bool) {
	if ws = workspace.Any(); ws == nil {
		jot.Error("no workspace available")
		return nil, nil, false
	}
	dc = ws.CurrentlyFocusedDockContainer()
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
func (d *Dockable) Setup(ws *workspace.Workspace, dc *unison.DockContainer, addToStartToolbar, addToEndToolbar, initContent func(*unison.Panel)) {
	d.SetLayout(&unison.FlexLayout{Columns: 1})
	d.AddChild(d.createToolbar(addToStartToolbar, addToEndToolbar))
	content := unison.NewPanel()
	content.SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(unison.StdHSpacing * 2)))
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
	if dc != nil && dc.Group == settingsGroup {
		dc.Stack(d, -1)
	} else if dc = ws.DocumentDock.ContainerForGroup(settingsGroup); dc != nil {
		dc.Stack(d, -1)
	} else {
		ws.DocumentDock.DockTo(d, nil, unison.RightSide)
		if dc = unison.DockContainerFor(d); dc != nil && dc.Group == "" {
			dc.Group = settingsGroup
		}
	}
}

// TitleIcon implements unison.Dockable
func (d *Dockable) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  res.SettingsSVG,
		Size: suggestedSize,
	}
}

// Title implements unison.Dockable
func (d *Dockable) Title() string {
	return d.TabTitle
}

// Tooltip implements unison.Dockable
func (d *Dockable) Tooltip() string {
	return ""
}

// Modified implements unison.Dockable
func (d *Dockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *Dockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *Dockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}

func (d *Dockable) createToolbar(addToStartToolbar, addToEndToolbar func(*unison.Panel)) *unison.Panel {
	toolbar := unison.NewPanel()
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(unison.Insets{
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
		b := unison.NewSVGButton(res.ResetSVG)
		b.Tooltip = unison.NewTooltipWithText(i18n.Text("Reset"))
		b.ClickCallback = d.handleReset
		toolbar.AddChild(b)
	}
	if d.Loader != nil || d.Saver != nil {
		b := unison.NewSVGButton(res.MenuSVG)
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

func (d *Dockable) handleReset() {
	if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Are you sure you want to reset the\n%s?"), d.TabTitle), "") == unison.ModalResponseOK {
		d.Resetter()
	}
}

func (d *Dockable) showMenu(b *unison.Button) {
	f := unison.DefaultMenuFactory()
	id := unison.ContextMenuIDFlag
	m := f.NewMenu(id, "", nil)
	id++
	if d.Loader != nil {
		m.InsertItem(-1, f.NewItem(id, i18n.Text("Import…"), unison.KeyBinding{}, nil, d.handleImport))
		id++
	}
	if d.Saver != nil {
		m.InsertItem(-1, f.NewItem(id, i18n.Text("Export…"), unison.KeyBinding{}, nil, d.handleExport))
		id++
	}
	if d.Loader != nil {
		libraries := settings.Global().Libraries()
		sets := library.ScanForNamedFileSets(nil, "", d.Extension, false, libraries)
		if len(sets) != 0 {
			m.InsertSeparator(-1, false)
			for _, lib := range sets {
				m.InsertItem(-1, f.NewItem(id, lib.Name, unison.KeyBinding{},
					func(_ unison.MenuItem) bool { return false }, nil))
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

func (d *Dockable) insertFileToLoad(m unison.Menu, id int, ref *library.NamedFileRef) {
	m.InsertItem(-1, m.Factory().NewItem(id, ref.Name, unison.KeyBinding{}, nil, func(_ unison.MenuItem) {
		d.doLoad(ref.FileSystem, ref.FilePath)
	}))
}

func (d *Dockable) doLoad(fileSystem fs.FS, filePath string) {
	if err := d.Loader(fileSystem, filePath); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to load ")+d.TabTitle, err)
	}
}

func (d *Dockable) handleImport(_ unison.MenuItem) {
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

func (d *Dockable) handleExport(_ unison.MenuItem) {
	dialog := unison.NewSaveDialog()
	dialog.SetAllowedExtensions(d.Extension)
	if dialog.RunModal() {
		if err := d.Saver(dialog.Path()); err != nil {
			unison.ErrorDialogWithError(i18n.Text("Unable to save ")+d.TabTitle, err)
		}
	}
}
