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
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/dsettings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type dockable struct {
	dsettings.Dockable
	content *unison.Panel
}

// Show the Page Reference Mappings.
func Show() {
	ws, dc, found := dsettings.Activate(func(d unison.Dockable) bool {
		_, ok := d.(*dockable)
		return ok
	})
	if !found && ws != nil {
		d := &dockable{}
		d.Self = d
		d.TabTitle = i18n.Text("Page Reference Mappings")
		d.Extension = ".refs"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
		if len(d.content.Children()) > 1 {
			d.content.Children()[1].RequestFocus()
		}
	}
}

func (d *dockable) initContent(content *unison.Panel) {
	d.content = content
	d.content.SetLayout(&unison.FlexLayout{
		Columns:  4,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.sync()
}

func (d *dockable) reset() {
	settings.Global().PageRefs = settings.PageRefs{}
	d.sync()
}

func (d *dockable) sync() {
	d.content.RemoveAllChildren()
	for _, one := range settings.Global().PageRefs.List() {
		d.createIDField(one)
		d.createOffsetField(one)
		d.createNameField(one)
		d.createTrashField(one)
	}
	d.MarkForRedraw()
}

func (d *dockable) createIDField(ref *settings.PageRef) {
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

func (d *dockable) createOffsetField(ref *settings.PageRef) {
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

func (d *dockable) createNameField(ref *settings.PageRef) {
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

func (d *dockable) createTrashField(ref *settings.PageRef) {
	b := unison.NewSVGButton(icons.TrashSVG())
	b.ClickCallback = func() {
		if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Are you sure you want to remove\n%s (%s)?"), ref.ID,
			filepath.Base(ref.Path)), "") == unison.ModalResponseOK {
			settings.Global().PageRefs.Remove(ref.ID)
			parent := b.Parent()
			index := parent.IndexOfChild(b)
			for i := index; i > index-4; i-- {
				parent.RemoveChildAtIndex(i)
			}
			parent.MarkForLayoutAndRedraw()
		}
	}
	b.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.MiddleAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(b)
}

func (d *dockable) load(fileSystem fs.FS, filePath string) error {
	s, err := settings.NewPageRefsFromFS(fileSystem, filePath)
	if err != nil {
		return err
	}
	settings.Global().PageRefs = *s
	d.sync()
	return nil
}

func (d *dockable) save(filePath string) error {
	return settings.Global().PageRefs.Save(filePath)
}
