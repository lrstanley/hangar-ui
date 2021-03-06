// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/offset"
)

type Root struct {
	*Base
}

func NewRoot(app types.App) *Root {
	return &Root{
		Base: &Base{
			app:    app,
			is:     types.ViewRoot,
			logger: log.WithField("src", "root"),
		},
	}
}

func (v *Root) Init() tea.Cmd { return nil }

func (v *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	v.logger.Debugf("msg: %#v", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
	case tea.MouseMsg:
		if !offset.Get(string(v.is)).InBounds(msg) {
			return v, nil
		}

		switch msg.Type {
		case tea.MouseLeft, tea.MouseRight:
			v.app.SetFocused(v.is)
		}
	}
	return v, nil
}

func (v *Root) View() string {
	s := lipgloss.NewStyle().
		Width(v.width-2). // 2 for border
		Height(v.height-2).
		MaxHeight(v.height).
		MaxWidth(v.width).
		Padding(0, 1).
		Background(types.Theme.Bg).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(types.Theme.ViewBorderBg).
		BorderForeground(types.Theme.ViewBorderInactiveFg)

	if v.Focused() {
		s = s.BorderForeground(types.Theme.ViewBorderActiveFg)
	}

	return offset.ID(string(v.is), s.Render("// ROOT"))
}
