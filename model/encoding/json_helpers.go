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

package encoding

// JSONer defines the method for objects that can turn themselves into JSON.
type JSONer interface {
	ToJSON(encoder *JSONEncoder)
}

// Empty returns true if the object is empty.
type Empty interface {
	Empty() bool
}

// ToKeyedJSON adds a key and emits the object, unless it is empty.
func ToKeyedJSON(obj JSONer, key string, encoder *JSONEncoder) {
	if obj == nil {
		return
	}
	if empty, ok := obj.(Empty); ok && empty.Empty() {
		return
	}
	encoder.Key(key)
	obj.ToJSON(encoder)
}
