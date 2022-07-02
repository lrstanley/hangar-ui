// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/hangar-ui/internal/types"
)

type Base struct {
	app types.App
	is  types.Viewable

	height int
	width  int
}

func (v *Base) Focused() bool {
	return v.app.IsFocused(v.is)
}

func (v *Base) Active() bool {
	return v.app.Active() == v.is
}

type View interface {
	tea.Model

	Active() bool
}
