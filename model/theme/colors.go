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

package theme

import (
	"context"
	"io/fs"

	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// Additional colors over and above what unison provides by default.
var (
	HeaderColor                = &unison.ThemeColor{Light: unison.RGB(43, 43, 43), Dark: unison.RGB(64, 64, 64)}
	OnHeaderColor              = &unison.ThemeColor{Light: unison.White, Dark: unison.Silver}
	EditableBorderColor        = &unison.ThemeColor{Light: unison.Silver, Dark: unison.RGB(96, 96, 96)}
	EditableBorderFocusedColor = &unison.ThemeColor{Light: unison.RGB(0, 0, 192), Dark: unison.RGB(0, 102, 102)}
	AccentColor                = &unison.ThemeColor{Light: unison.RGB(0, 102, 102), Dark: unison.RGB(100, 153, 153)}
	SearchListColor            = &unison.ThemeColor{Light: unison.LightCyan, Dark: unison.RGB(0, 43, 43)}
	OnSearchListColor          = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(204, 204, 204)}
	PageColor                  = &unison.ThemeColor{Light: unison.White, Dark: unison.RGB(16, 16, 16)}
	OnPageColor                = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(160, 160, 160)}
	PageVoidColor              = &unison.ThemeColor{Light: unison.Grey, Dark: unison.Black}
	MarkerColor                = &unison.ThemeColor{Light: unison.RGB(252, 242, 196), Dark: unison.RGB(0, 51, 0)}
	OnMarkerColor              = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(221, 221, 221)}
	OverloadedColor            = &unison.ThemeColor{Light: unison.RGB(192, 64, 64), Dark: unison.RGB(115, 37, 37)}
	OnOverloadedColor          = &unison.ThemeColor{Light: unison.White, Dark: unison.RGB(221, 221, 221)}
	HintColor                  = &unison.ThemeColor{Light: unison.Grey, Dark: unison.RGB(64, 64, 64)}
	LinkColor                  = &unison.ThemeColor{Light: unison.SpringGreen, Dark: unison.SpringGreen}
	OnLinkColor                = &unison.ThemeColor{Light: unison.Black, Dark: unison.Black}
	PDFLinkHighlightColor      = &unison.ThemeColor{Light: unison.SpringGreen, Dark: unison.SpringGreen}
	PDFMarkerHighlightColor    = &unison.ThemeColor{Light: unison.Yellow, Dark: unison.Yellow}
)

var (
	// CurrentColors holds the current theme.
	CurrentColors = []*ThemedColor{
		{ID: "background", Title: i18n.Text("Background"), Color: unison.BackgroundColor},
		{ID: "on_background", Title: i18n.Text("On Background"), Color: unison.OnBackgroundColor},
		{ID: "content", Title: i18n.Text("Content"), Color: unison.ContentColor},
		{ID: "on_content", Title: i18n.Text("On Content"), Color: unison.OnContentColor},
		{ID: "banding", Title: i18n.Text("Banding"), Color: unison.BandingColor},
		{ID: "on_banding", Title: i18n.Text("On Banding"), Color: unison.OnBandingColor},
		{ID: "header", Title: i18n.Text("Header"), Color: HeaderColor},
		{ID: "on_header", Title: i18n.Text("On Header"), Color: OnHeaderColor},
		{ID: "tab_focused", Title: i18n.Text("Focused Tab"), Color: unison.TabFocusedColor},
		{ID: "on_tab_focused", Title: i18n.Text("On Focused Tab"), Color: unison.OnTabFocusedColor},
		{ID: "tab_current", Title: i18n.Text("Current Tab"), Color: unison.TabCurrentColor},
		{ID: "on_tab_current", Title: i18n.Text("On Current Tab"), Color: unison.OnTabCurrentColor},
		{ID: "editable", Title: i18n.Text("Editable"), Color: unison.EditableColor},
		{ID: "on_editable", Title: i18n.Text("On Editable"), Color: unison.OnEditableColor},
		{ID: "editable_border", Title: i18n.Text("Editable Border"), Color: EditableBorderColor},
		{ID: "editable_border_focused", Title: i18n.Text("Focused Editable Border"), Color: EditableBorderFocusedColor},
		{ID: "selection", Title: i18n.Text("Selection"), Color: unison.SelectionColor},
		{ID: "on_selection", Title: i18n.Text("On Selection"), Color: unison.OnSelectionColor},
		{ID: "inactive_selection", Title: i18n.Text("Inactive Selection"), Color: unison.InactiveSelectionColor},
		{ID: "on_inactive_selection", Title: i18n.Text("On Inactive Selection"), Color: unison.OnInactiveSelectionColor},
		{ID: "indirect_selection", Title: i18n.Text("Indirect Selection"), Color: unison.IndirectSelectionColor},
		{ID: "on_indirect_selection", Title: i18n.Text("On Indirect Selection"), Color: unison.OnIndirectSelectionColor},
		{ID: "scroll", Title: i18n.Text("Scroll"), Color: unison.ScrollColor},
		{ID: "scroll_rollover", Title: i18n.Text("Scroll Rollover"), Color: unison.ScrollRolloverColor},
		{ID: "scroll_edge", Title: i18n.Text("Scroll Edge"), Color: unison.ScrollEdgeColor},
		{ID: "accent", Title: i18n.Text("Accent"), Color: AccentColor},
		{ID: "control", Title: i18n.Text("Control"), Color: unison.ControlColor},
		{ID: "on_control", Title: i18n.Text("On Control"), Color: unison.OnControlColor},
		{ID: "control_pressed", Title: i18n.Text("Pressed Control"), Color: unison.ControlPressedColor},
		{ID: "on_control_pressed", Title: i18n.Text("On Pressed Control"), Color: unison.OnControlPressedColor},
		{ID: "control_edge", Title: i18n.Text("Control Edge"), Color: unison.ControlEdgeColor},
		{ID: "divider", Title: i18n.Text("Divider"), Color: unison.DividerColor},
		{ID: "icon_button", Title: i18n.Text("Icon Button"), Color: unison.IconButtonColor},
		{ID: "icon_button_rollover", Title: i18n.Text("Icon Button Rollover"), Color: unison.IconButtonRolloverColor},
		{ID: "icon_button_pressed", Title: i18n.Text("Pressed Icon Button"), Color: unison.IconButtonPressedColor},
		{ID: "drop_area", Title: i18n.Text("Drop Area"), Color: unison.DropAreaColor},
		{ID: "tooltip", Title: i18n.Text("Tooltip"), Color: unison.TooltipColor},
		{ID: "on_tooltip", Title: i18n.Text("On Tooltip"), Color: unison.OnTooltipColor},
		{ID: "search_list", Title: i18n.Text("Search List"), Color: SearchListColor},
		{ID: "on_search_list", Title: i18n.Text("On Search List"), Color: OnSearchListColor},
		{ID: "marker", Title: i18n.Text("Marker"), Color: MarkerColor},
		{ID: "on_marker", Title: i18n.Text("On Marker"), Color: OnMarkerColor},
		{ID: "error", Title: i18n.Text("Error"), Color: unison.ErrorColor},
		{ID: "on_error", Title: i18n.Text("On Error"), Color: unison.OnErrorColor},
		{ID: "warning", Title: i18n.Text("Warning"), Color: unison.WarningColor},
		{ID: "on_warning", Title: i18n.Text("On Warning"), Color: unison.OnWarningColor},
		{ID: "overloaded", Title: i18n.Text("Overloaded"), Color: OverloadedColor},
		{ID: "on_overloaded", Title: i18n.Text("On Overloaded"), Color: OnOverloadedColor},
		{ID: "page", Title: i18n.Text("Page"), Color: PageColor},
		{ID: "on_page", Title: i18n.Text("On Page"), Color: OnPageColor},
		{ID: "page_void", Title: i18n.Text("Page Void"), Color: PageVoidColor},
		{ID: "hint", Title: i18n.Text("Hint"), Color: HintColor},
		{ID: "link", Title: i18n.Text("Link"), Color: LinkColor},
		{ID: "on_link", Title: i18n.Text("On Link"), Color: OnLinkColor},
		{ID: "pdf_link", Title: i18n.Text("PDF Link Highlight"), Color: PDFLinkHighlightColor},
		{ID: "pdf_marker", Title: i18n.Text("PDF Marker Highlight"), Color: PDFMarkerHighlightColor},
	}
	// FactoryColors holds the original theme before any modifications.
	FactoryColors []*ThemedColor
)

func init() {
	FactoryColors = make([]*ThemedColor, len(CurrentColors))
	for i, c := range CurrentColors {
		FactoryColors[i] = &ThemedColor{
			ID:    c.ID,
			Title: c.Title,
			Color: &unison.ThemeColor{
				Light: c.Color.Light,
				Dark:  c.Color.Dark,
			},
		}
	}
}

// ThemedColor holds a themed color.
type ThemedColor struct {
	ID    string
	Title string
	Color *unison.ThemeColor
}

// Colors holds a set of themed colors.
type Colors struct {
	data map[string]*unison.ThemeColor // Just here for serialization
}

// NewColorsFromFS creates a new set of colors from a file. Any missing values will be filled in with defaults.
func NewColorsFromFS(fileSystem fs.FS, filePath string) (*Colors, error) {
	var c Colors
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Save writes the Colors to the file as JSON.
func (c *Colors) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, c)
}

// MarshalJSON implements json.Marshaler.
func (c *Colors) MarshalJSON() ([]byte, error) {
	c.data = make(map[string]*unison.ThemeColor, len(CurrentColors))
	for _, one := range CurrentColors {
		c.data[one.ID] = one.Color
	}
	return json.Marshal(&c.data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Colors) UnmarshalJSON(data []byte) error {
	c.data = make(map[string]*unison.ThemeColor, len(CurrentColors))
	if err := json.Unmarshal(data, &c.data); err != nil {
		return err
	}
	for _, one := range FactoryColors {
		if _, ok := c.data[one.ID]; !ok {
			c.data[one.ID] = one.Color
		}
	}
	return nil
}

// MakeCurrent applies these colors to the current theme color set and updates all windows.
func (c *Colors) MakeCurrent() {
	for _, one := range CurrentColors {
		if v, ok := c.data[one.ID]; ok {
			*one.Color = *v
		}
	}
	unison.ThemeChanged()
}

// Reset to factory defaults.
func (c *Colors) Reset() {
	c.data = make(map[string]*unison.ThemeColor, len(CurrentColors))
	for _, one := range FactoryColors {
		c.data[one.ID] = one.Color
	}
}

// ResetOne resets one color by ID to factory defaults.
func (c *Colors) ResetOne(id string) {
	for _, v := range FactoryColors {
		if v.ID == id {
			c.data[id] = v.Color
			break
		}
	}
}
