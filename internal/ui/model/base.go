// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/hangar-ui/internal/types"
)

type Base struct {
	tea.Model
	app types.App
	is  types.Viewable

	Height int
	Width  int

	buf strings.Builder
}

func (v *Base) Focused() bool {
	return v.app.IsFocused(v.is)
}

type Msg int

const MsgNone Msg = -1
