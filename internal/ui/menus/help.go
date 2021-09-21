package menus

import (
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func setupHelpMenu(bar unison.Menu) {
	f := bar.Factory()
	m := bar.Menu(unison.HelpMenuID)
	m.InsertItem(-1, SponsorGCSDevelopment.NewMenuItem(f))
	m.InsertItem(-1, MakeDonation.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, UpdateApp.NewMenuItem(f))
	m.InsertItem(-1, ReleaseNotes.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, WebSite.NewMenuItem(f))
}

// SponsorGCSDevelopment opens the web site for sponsoring GCS development.
var SponsorGCSDevelopment = &unison.Action{
	ID:    SponsorGCSDevelopmentItemID,
	Title: i18n.Text("Sponsor GCS Development"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://github.com/sponsors/richardwilkes")
	},
}

// MakeDonation opens the web site for make a donation.
var MakeDonation = &unison.Action{
	ID:    MakeDonationItemID,
	Title: i18n.Text("Make A Donation For GCS Development"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://paypal.me/GURPSCharacterSheet")
	},
}

// UpdateApp opens the web site for GCS updates.
var UpdateApp = &unison.Action{
	ID:              UpdateAppItemID,
	Title:           i18n.Text("Checking for GCS updates…"),
	ExecuteCallback: unimplemented,
}

// ReleaseNotes opens the release notes.
var ReleaseNotes = &unison.Action{
	ID:    ReleaseNotesItemID,
	Title: i18n.Text("Release Notes"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://github.com/richardwilkes/gcs/releases")
	},
}

// License opens the license.
var License = &unison.Action{
	ID:    ReleaseNotesItemID,
	Title: i18n.Text("License"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://github.com/richardwilkes/gcs/blob/master/LICENSE")
	},
}

// WebSite opens the GCS web site.
var WebSite = &unison.Action{
	ID:    WebSiteItemID,
	Title: i18n.Text("Web Site"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://gurpscharactersheet.com")
	},
}

// MailingList opens the GCS mailing list site.
var MailingList = &unison.Action{
	ID:    MailingListItemID,
	Title: i18n.Text("Mailing Lists"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		desktop.OpenBrowser("https://groups.io/g/gcs")
	},
}
