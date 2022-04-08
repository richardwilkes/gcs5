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

package gid

// Field undo IDs
const (
	FieldNone = iota
	FieldName
	FieldTitle
	FieldOrganization
	FieldPlayer
	FieldGender
	FieldAge
	FieldBirthday
	FieldReligion
	FieldTechLevel
	FieldHair
	FieldEyes
	FieldSkin
	FieldHand
	FieldSize
	FieldHeight
	FieldWeight
	FieldDefaultPlayerName
	FieldDefaultTechLevel
	FieldGCalcKey
	FieldImageExportResolution
	FieldPageOffset
	FieldFontSize
	FieldInitialPoints
	FieldTooltipDelay
	FieldTooltipDismissal
	FieldUnspentPoints
)

// Base field undo IDs where the field count is dynamic.
const (
	FieldPrimaryAttributeBase = 1000 * (iota + 1)
	FieldSecondaryAttributeBase
	FieldPointPoolCurrentBase
	FieldPointPoolMaximumBase
)
