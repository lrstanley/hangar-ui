// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
)

type Targets struct {
	*Base
}

func NewTargets(app types.App) *Targets {
	return &Targets{
		Base: &Base{
			app: app,
			is:  types.ViewTargets,
		},
	}
}

func (v *Targets) Init() tea.Cmd { return nil }

func (v *Targets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Targets.Update: %#v", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft, tea.MouseRight:
			v.app.SetFocused(v.is)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, types.KeyCancel):
			v.app.Back(true)
			return v, nil
		}
	}
	return v, nil
}

func (v *Targets) View() string {
	s := lipgloss.NewStyle().
		Width(v.width - 2). // -2 for border
		Height(v.height - 2).
		MaxHeight(v.height).
		MaxWidth(v.width).
		Background(types.Theme.Bg).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(types.Theme.ViewBorderBg).
		BorderForeground(types.Theme.ViewBorderInactiveFg)

	if v.Focused() {
		s = s.BorderForeground(types.Theme.ViewBorderActiveFg)
	}

	return s.Render("// TARGETS")
}
