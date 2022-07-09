// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/icons"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/offset"
	"github.com/lrstanley/hangar-ui/internal/x"
)

const (
	helpSeparator = " • "
	helpEllipsis  = "…"
)

var icon = icons.IconSet["console"].GetGlyph()

type StatusBar struct {
	*Base

	keys *KeyMap

	Target string
	URL    string
	Logo   string

	baseStyle   lipgloss.Style
	targetStyle lipgloss.Style
	urlStyle    lipgloss.Style
	logoStyle   lipgloss.Style
	descStyle   lipgloss.Style

	separator string
}

func NewStatusBar(app types.App, keys *KeyMap) *StatusBar {
	m := &StatusBar{
		Base: &Base{
			app:    app,
			is:     types.ViewStatusBar,
			Height: 1,
			logger: log.WithField("src", "statusbar"),
		},
		keys:   keys,
		Target: "target",
		URL:    "concourse.example.com",
		Logo:   "hangar-ui",
	}

	m.baseStyle = lipgloss.NewStyle().
		Foreground(types.Theme.StatusBarFg).
		Background(types.Theme.StatusBarBg)

	m.targetStyle = lipgloss.NewStyle().Padding(0, 1).
		Background(types.Theme.StatusBarTargetBg).
		Foreground(types.Theme.StatusBarTargetFg)

	m.urlStyle = lipgloss.NewStyle().Padding(0, 1).
		Background(types.Theme.StatusBarURLBg).
		Foreground(types.Theme.StatusBarURLFg)

	m.logoStyle = lipgloss.NewStyle().Padding(0, 1).
		Background(types.Theme.StatusBarLogoBg).
		Foreground(types.Theme.StatusBarLogoFg).Bold(true)

	m.descStyle = m.baseStyle.Copy().
		Foreground(types.Theme.StatusBarKeyDescFg)

	m.separator = m.baseStyle.Copy().
		Foreground(types.Theme.StatusBarTargetBg).
		Render(helpSeparator)

	return m
}

func (m *StatusBar) Init() tea.Cmd {
	return nil
}

func (m *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debugf("msg: %#v", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			if offset.GetArea("statusbar_target").InBounds(msg) {
				m.app.SetActive(types.ViewTargets, true)
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *StatusBar) View() string {
	target := m.targetStyle.Render(m.Target)
	url := m.urlStyle.Render(m.URL)
	logo := m.logoStyle.Render(m.Logo)

	help := ""
	bindings := m.keys.ShortHelp()
	helpWidth := m.Width - x.WMulti(target, url, logo) - 2

	var totalWidth, w int
	var str, tail string
	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if helpWidth > 0 && i > 0 {
			sep = m.separator
		}

		str = sep + m.baseStyle.Render("<"+kb.Help().Key+">") +
			m.baseStyle.Render(" ") +
			m.descStyle.Render(kb.Help().Desc)

		w = x.W(str)

		if helpWidth > 0 && totalWidth+w > helpWidth {
			// If there's room for an ellipsis, print that.
			tail = m.baseStyle.Render(" " + helpEllipsis)

			if totalWidth+x.W(tail) < helpWidth {
				help += tail
			}

			break
		}

		totalWidth += w
		help += str
	}

	help = m.baseStyle.Copy().Width(helpWidth+2).Align(lipgloss.Right).Padding(0, 1).Render(help)

	return x.X(0, offset.AreaID("statusbar_target", target), help, url, logo)
}
