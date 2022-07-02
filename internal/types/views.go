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
	ViewHelp        Viewable = "help"
	ViewAbout       Viewable = "about"
	SubViewSomeItem Viewable = "someitem"
)

type ViewChangeMsg struct {
	View Viewable
}

type FocusChangeMsg struct {
	View Viewable
}

type ViewMsg struct {
	View Viewable
	Msg  tea.Msg
}
