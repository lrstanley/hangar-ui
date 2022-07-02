// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App interface {
	tea.Model

	SetFocused(v Viewable)
	IsFocused(v Viewable) bool
	SetActive(v Viewable, focused bool)
	Active() Viewable
	Previous() Viewable
	Back(focused bool)
}
