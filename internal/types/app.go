// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App interface {
	tea.Model

	IsFocused(v Viewable) bool
	Active() Viewable
	Previous() Viewable
}
