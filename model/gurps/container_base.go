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

package gurps

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/i18n"
)

// ContainerKeyPostfix is the key postfix used to identify containers.
const ContainerKeyPostfix = "_container"

// ContainerBase holds the type and ID of the data.
type ContainerBase[T node.Node] struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	IsOpen   bool      `json:"open,omitempty"`     // Container only
	Children []T       `json:"children,omitempty"` // Container only
}

func newContainerBase[T node.Node](typeKey string, isContainer bool) ContainerBase[T] {
	if isContainer {
		typeKey += ContainerKeyPostfix
	}
	return ContainerBase[T]{
		ID:     id.NewUUID(),
		Type:   typeKey,
		IsOpen: isContainer,
	}
}

// UUID returns the UUID of this data.
func (c *ContainerBase[T]) UUID() uuid.UUID {
	return c.ID
}

// Container returns true if this is a container.
func (c *ContainerBase[T]) Container() bool {
	return strings.HasSuffix(c.Type, ContainerKeyPostfix)
}

func (c *ContainerBase[T]) kind(base string) string {
	if c.Container() {
		return fmt.Sprintf(i18n.Text("%s Container"), base)
	}
	return base
}

// Open returns true if this node is currently open.
func (c *ContainerBase[T]) Open() bool {
	return c.IsOpen && c.Container()
}

// SetOpen sets the current open state for this node.
func (c *ContainerBase[T]) SetOpen(open bool) {
	c.IsOpen = open && c.Container()
}

// NodeChildren returns the children of this node, if any.
func (c *ContainerBase[T]) NodeChildren() []node.Node {
	if c.Container() {
		children := make([]node.Node, len(c.Children))
		for i, child := range c.Children {
			children[i] = child
		}
		return children
	}
	return nil
}

func (c *ContainerBase[T]) clearUnusedFields() {
	if !c.Container() {
		c.Children = nil
		c.IsOpen = false
	}
}
