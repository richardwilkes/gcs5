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

package menus

import (
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// WebSiteDomain holds the web site domain for GCS.
const WebSiteDomain = "gurpscharactersheet.com"

func setupHelpMenu(bar unison.Menu) {
	f := bar.Factory()
	m := bar.Menu(unison.HelpMenuID)
	m.InsertItem(-1, SponsorGCSDevelopment.NewMenuItem(f))
	m.InsertItem(-1, MakeDonation.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, UpdateApp.NewMenuItem(f))
	m.InsertItem(-1, ReleaseNotes.NewMenuItem(f))
	m.InsertItem(-1, License.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, WebSite.NewMenuItem(f))
	m.InsertItem(-1, MailingList.NewMenuItem(f))
}

// SponsorGCSDevelopment opens the web site for sponsoring GCS development.
var SponsorGCSDevelopment = &unison.Action{
	ID:    SponsorGCSDevelopmentItemID,
	Title: i18n.Text("Sponsor GCS Development"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		showWebPage("https://github.com/sponsors/richardwilkes")
	},
}

// MakeDonation opens the web site for make a donation.
var MakeDonation = &unison.Action{
	ID:    MakeDonationItemID,
	Title: i18n.Text("Make A Donation For GCS Development"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		showWebPage("https://paypal.me/GURPSCharacterSheet")
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
		showWebPage("https://github.com/richardwilkes/gcs/releases")
	},
}

// License opens the license.
var License = &unison.Action{
	ID:    ReleaseNotesItemID,
	Title: i18n.Text("License"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		showWebPage("https://github.com/richardwilkes/gcs/blob/master/LICENSE")
	},
}

// WebSite opens the GCS web site.
var WebSite = &unison.Action{
	ID:    WebSiteItemID,
	Title: i18n.Text("Web Site"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		showWebPage("https://" + WebSiteDomain)
	},
}

// MailingList opens the GCS mailing list site.
var MailingList = &unison.Action{
	ID:    MailingListItemID,
	Title: i18n.Text("Mailing Lists"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		showWebPage("https://groups.io/g/gcs")
	},
}

func showWebPage(uri string) {
	if err := desktop.OpenBrowser(uri); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to open link"), err)
	}
}
