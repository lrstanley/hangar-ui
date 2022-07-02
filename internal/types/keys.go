// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import "github.com/charmbracelet/bubbles/key"

var (
	KeyCancel = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	KeyHelp = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	KeyQuit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	KeyUp = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "go up"),
	)
	KeyDown = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "go down"),
	)
	KeyLeft = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go left"),
	)
	KeyRight = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "go right"),
	)
	KeyEnter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)

	// Specifically for the command bar.
	KeyCmdFilter = key.NewBinding(
		key.WithKeys("/", "ctrl+f"),
		key.WithHelp("/", "filter results"),
	)
	KeyCmdInvoke = key.NewBinding(
		key.WithKeys(":"),
		key.WithHelp(":", "run command"),
	)
)
