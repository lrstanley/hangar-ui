// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/offset"
	"github.com/lrstanley/hangar-ui/internal/x"
)

const (
	navBarPadding    = 1
	navBarItemMargin = 1
)

type NavBar struct {
	*Base

	views []types.Viewable

	activeStyle   lipgloss.Style
	inactiveStyle lipgloss.Style
}

func NewNavBar(app types.App, views []types.Viewable) *NavBar {
	m := &NavBar{
		Base: &Base{
			app:    app,
			is:     types.ViewNavigation,
			Height: 1,
			logger: log.WithField("src", "navbar"),
		},
		views: views,
	}

	m.activeStyle = lipgloss.NewStyle().
		Foreground(types.Theme.NavActiveFg).
		Background(types.Theme.NavActiveBg).
		Padding(0, 1).
		MarginRight(navBarItemMargin).
		MarginBackground(types.Theme.Bg)

	m.inactiveStyle = m.activeStyle.Copy().
		Foreground(types.Theme.NavInactiveFg).
		Background(types.Theme.NavInactiveBg)

	return m
}

func (m *NavBar) Init() tea.Cmd {
	return nil
}

func (m *NavBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debugf("msg: %#v", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			for _, v := range m.views {
				if offset.GetArea(string(v)).InBounds(msg) {
					m.app.SetActive(v, true)
					return m, nil
				}
			}
		}
	}

	return m, nil
}

func (m *NavBar) View() string {
	m.buf.Reset()
	s := lipgloss.NewStyle().
		Width(m.Width).
		MaxWidth(m.Width).
		Padding(0, navBarPadding).
		Background(types.Theme.Bg)

	active := m.app.Active()
	var style lipgloss.Style

	for i, v := range m.views {
		if v == active {
			style = m.activeStyle.Copy()
		} else {
			style = m.inactiveStyle.Copy()
		}

		if i+1 == len(m.views) {
			break
		}

		m.buf.WriteString(offset.AreaID(string(v), style.Render(string(v))))
	}

	lastV := string(m.views[len(m.views)-1])

	return s.Render(m.buf.String() + x.PlaceX(
		m.Width-x.W(m.buf.String())-(2*navBarPadding),
		x.Right,
		offset.AreaID(lastV, style.Margin(0).Render(lastV)),
		lipgloss.WithWhitespaceBackground(types.Theme.Bg),
	))
}
