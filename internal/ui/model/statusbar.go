// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/icons"
	zone "github.com/lrstanley/bubblezone"
	"github.com/lrstanley/hangar-ui/internal/api"
	"github.com/lrstanley/hangar-ui/internal/types"
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

	loadingText string
	spinner     spinner.Model

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
		Target: api.Manager.ActiveName(),
		URL:    api.Manager.Active().URL(),
		Logo:   "hangar-ui",

		spinner: spinner.New(),
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

	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = m.baseStyle.Copy().
		Foreground(types.Theme.StatusBarKeyDescFg)

	return m
}

func (m *StatusBar) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			if zone.Get("statusbar_target").InBounds(msg) {
				return m, tea.Batch(
					m.spinner.Tick,
					types.MsgAsCmd(types.ViewChangeMsg{View: types.ViewTargets}),
					types.MsgAsCmd(types.FocusChangeMsg{View: types.ViewTargets}),
				)
			} else if zone.Get("statusbar_help").InBounds(msg) {
				return m, tea.Batch(
					m.spinner.Tick,
					types.MsgAsCmd(types.ViewChangeMsg{View: types.ViewHelp}),
					types.MsgAsCmd(types.FocusChangeMsg{View: types.ViewHelp}),
				)
			}
		}
	case types.FlyEvent:
		if msg == types.FlyActiveTargetUpdated {
			m.Target = api.Manager.ActiveName()
			m.URL = api.Manager.Active().URL()
		}
	case types.LoadingMsg:
		m.loadingText = msg.Text
		return m, m.spinner.Tick
	case types.CancelLoadingMsg:
		m.loadingText = ""
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, m.spinner.Tick
}

func (m *StatusBar) View() string {
	target := m.targetStyle.Render(m.Target)
	url := m.urlStyle.Render(m.URL)
	logo := m.logoStyle.Render(m.Logo)
	loading := ""

	if m.loadingText != "" {
		loading = m.descStyle.Render(" ") + m.spinner.View() + m.descStyle.Render(m.loadingText)
	}

	help := ""
	bindings := m.keys.ShortHelp()
	helpWidth := m.Width - x.WMulti(target, url, logo, loading) - 2

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

	return x.X(
		0,
		zone.Mark("statusbar_target", target),
		loading,
		zone.Mark("statusbar_help", help),
		url,
		logo,
	)
}
