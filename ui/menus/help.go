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
	"fmt"

	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

var (
	// SponsorGCSDevelopment opens the web site for sponsoring GCS development.
	SponsorGCSDevelopment *unison.Action
	// MakeDonation opens the web site for make a donation.
	MakeDonation *unison.Action
	// UpdateApp opens the web site for GCS updates.
	UpdateApp *unison.Action
	// ReleaseNotes opens the release notes.
	ReleaseNotes *unison.Action
	// License opens the license.
	License *unison.Action
	// WebSite opens the GCS web site.
	WebSite *unison.Action
	// MailingList opens the GCS mailing list site.
	MailingList *unison.Action
)

func registerHelpMenuActions() {
	SponsorGCSDevelopment = &unison.Action{
		ID:    constants.SponsorGCSDevelopmentItemID,
		Title: fmt.Sprintf(i18n.Text("Sponsor %s Development"), cmdline.AppName),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://github.com/sponsors/richardwilkes")
		},
	}
	MakeDonation = &unison.Action{
		ID:    constants.MakeDonationItemID,
		Title: fmt.Sprintf(i18n.Text("Make a One-time Donation for %s Development"), cmdline.AppName),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://paypal.me/GURPSCharacterSheet")
		},
	}
	UpdateApp = &unison.Action{
		ID: constants.UpdateAppItemID,
		EnabledCallback: func(action *unison.Action, mi any) bool {
			var releases []library.Release
			action.Title, releases = library.AppUpdateResult()
			if menuItem, ok := mi.(unison.MenuItem); ok {
				menuItem.SetTitle(action.Title)
			}
			return releases != nil
		},
		ExecuteCallback: func(_ *unison.Action, _ any) { library.AppUpdate() },
	}
	ReleaseNotes = &unison.Action{
		ID:    constants.ReleaseNotesItemID,
		Title: i18n.Text("Release Notes"),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://github.com/richardwilkes/gcs/releases")
		},
	}
	License = &unison.Action{
		ID:    constants.ReleaseNotesItemID,
		Title: i18n.Text("License"),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://github.com/richardwilkes/gcs/blob/master/LICENSE")
		},
	}
	WebSite = &unison.Action{
		ID:    constants.WebSiteItemID,
		Title: i18n.Text("Web Site"),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://" + constants.WebSiteDomain)
		},
	}
	MailingList = &unison.Action{
		ID:    constants.MailingListItemID,
		Title: i18n.Text("Mailing Lists"),
		ExecuteCallback: func(_ *unison.Action, _ any) {
			showWebPage("https://groups.io/g/gcs")
		},
	}
}

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

func showWebPage(uri string) {
	if err := desktop.Open(uri); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to open link"), err)
	}
}
