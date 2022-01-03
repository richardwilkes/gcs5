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

package library

import (
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

// Version holds a standard version value.
type Version struct {
	Major  int
	Minor  int
	BugFix int
}

// MinimumLibraryVersion is the oldest version of the library data that can be loaded.
func MinimumLibraryVersion() Version {
	return Version{}
}

// IncompatibleFutureLibraryVersion is the newest version at which the library data can no longer be loaded.
func IncompatibleFutureLibraryVersion() Version {
	return Version{Major: 4}
}

// VersionFromString parses the text for a version.
func VersionFromString(text string) Version {
	var v Version
	if err := v.UnmarshalText([]byte(text)); err != nil {
		jot.Warn(err)
	}
	return v
}

// Less returns true if this Version is less than the other.
func (v Version) Less(other Version) bool {
	if v.Major < other.Major {
		return true
	}
	if v.Major > other.Major {
		return false
	}
	if v.Minor < other.Minor {
		return true
	}
	if v.Minor > other.Minor {
		return false
	}
	return v.BugFix < other.BugFix
}

func (v Version) String() string {
	var buffer strings.Builder
	buffer.WriteString(strconv.FormatInt(int64(v.Major), 10))
	buffer.WriteByte('.')
	buffer.WriteString(strconv.FormatInt(int64(v.Minor), 10))
	if v.BugFix != 0 {
		buffer.WriteByte('.')
		buffer.WriteString(strconv.FormatInt(int64(v.BugFix), 10))
	}
	return buffer.String()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (v Version) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (v *Version) UnmarshalText(text []byte) error {
	v.Major = 0
	v.Minor = 0
	v.BugFix = 0
	if len(text) > 0 {
		parts := strings.SplitN(string(text), ".", 3)
		switch len(parts) {
		case 3:
			bugfix, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return errs.NewWithCausef(err, "unable to parse '%s'", parts[2])
			}
			v.BugFix = int(bugfix)
			fallthrough
		case 2:
			minor, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return errs.NewWithCausef(err, "unable to parse '%s'", parts[1])
			}
			v.Minor = int(minor)
			fallthrough
		case 1:
			major, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return errs.NewWithCausef(err, "unable to parse '%s'", parts[0])
			}
			v.Major = int(major)
		}
	}
	return nil
}
