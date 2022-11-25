// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import tea "github.com/charmbracelet/bubbletea"

type Viewable string

const (
	ViewRoot        Viewable = "main"
	ViewCommandBar  Viewable = "commandbar"
	ViewStatusBar   Viewable = "statusbar"
	ViewNavigation  Viewable = "navigation"
	ViewHelp        Viewable = "help"
	ViewTargets     Viewable = "targets"
	ViewAbout       Viewable = "about"
	SubViewSomeItem Viewable = "someitem"
)

// ViewChangeMsg is sent when the primary view is changed (not necessarily focused).
type ViewChangeMsg struct {
	View Viewable
}

// FocusChangeMsg is sent when the focused view changes.
type FocusChangeMsg struct {
	View Viewable
}

// ViewMsg is a message that is sent to a specific view (if available).
type ViewMsg struct {
	View Viewable
	Msg  tea.Msg
}

type AppBackMsg struct {
	Focused bool
}
