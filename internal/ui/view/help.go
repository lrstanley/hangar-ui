// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
	"github.com/lrstanley/hangar-ui/internal/ui/offset"
)

type Help struct {
	*Base

	keys  *model.KeyMap
	model viewport.Model

	titleStyle    lipgloss.Style
	keyStyle      lipgloss.Style
	keyInnerStyle lipgloss.Style
	descStyle     lipgloss.Style
}

func NewHelp(app types.App, keys *model.KeyMap) *Help {
	v := &Help{
		Base: &Base{
			app:    app,
			is:     types.ViewHelp,
			logger: log.WithField("src", "help"),
		},
		keys:  keys,
		model: viewport.New(0, 0),
	}

	v.titleStyle = lipgloss.NewStyle().
		Background(types.Theme.TitleBg).
		Foreground(types.Theme.TitleFg).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false)

	v.keyStyle = lipgloss.NewStyle().
		Background(types.Theme.Bg).
		Foreground(types.Theme.Fg)

	v.keyInnerStyle = lipgloss.NewStyle().
		Background(types.Theme.Bg).
		Foreground(types.Theme.TitleFg).
		Bold(true)

	v.descStyle = lipgloss.NewStyle().
		Background(types.Theme.Bg).
		Foreground(types.Theme.Fg)

	v.generateHelp()

	return v
}

func (v *Help) generateHelp() {
	var buf strings.Builder

	for view, bindings := range v.keys.Binds {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}

		buf.WriteString(v.titleStyle.Render(string(view)) + "\n")

		for _, binding := range bindings {
			buf.WriteString(
				v.keyStyle.Copy().Width(12).Render(
					v.keyStyle.Render("<")+
						v.keyInnerStyle.Render(binding.Help().Key)+
						v.keyStyle.Render(">"),
				) +
					v.descStyle.Render(binding.Help().Desc) +
					"\n",
			)
		}
	}
	v.model.SetContent(buf.String())
}

func (v *Help) Init() tea.Cmd { return nil }

func (v *Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	v.logger.Debugf("msg: %#v", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
		v.model.Height = msg.Height - v.model.Style.GetVerticalFrameSize() - 2 // 2 for border
		v.model.Width = msg.Width - v.model.Style.GetHorizontalFrameSize() - 2 // 2 for border
	case tea.MouseMsg:
		if !offset.Get(string(v.is)).InBounds(msg) {
			return v, nil
		}

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
	case types.ViewChangeMsg:
		v.generateHelp()
	}

	var cmd tea.Cmd
	v.model, cmd = v.model.Update(msg)
	return v, cmd
}

func (v *Help) View() string {
	out := v.model.View()

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

	return offset.ID(string(v.is), s.Render(out))
}
