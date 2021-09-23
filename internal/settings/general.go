package settings

import (
	"os"
	"os/user"
)

// General settings
type General struct {
	DefaultPlayerName           string  `json:"default_player_name"`
	DefaultTechLevel            string  `json:"default_tech_level"`
	PDFViewer                   string  `json:"pdf_viewer"`
	InitialPoints               int     `json:"initial_points"`
	TooltipTimeout              int     `json:"tooltip_timeout"`
	ImageResolution             int     `json:"image_resolution"`
	InitialUIScale              float32 `json:"initial_ui_scale"`
	AutoFillProfile             bool    `json:"auto_fill_profile"`
	IncludeUnspentPointsInTotal bool    `json:"include_unspent_points_in_total"`
}

// NewGeneral return new general settings.
func NewGeneral() *General {
	var name string
	if u, err := user.Current(); err != nil {
		name = os.Getenv("USER")
	} else {
		name = u.Name
	}
	return &General{
		DefaultPlayerName:           name,
		DefaultTechLevel:            "3",
		PDFViewer:                   "", // TODO: get default for platform
		InitialPoints:               250,
		TooltipTimeout:              60,
		ImageResolution:             200,
		InitialUIScale:              1.25,
		AutoFillProfile:             true,
		IncludeUnspentPointsInTotal: true,
	}
}
