// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
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
			app: app,
			is:  types.ViewHelp,
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
	log.Printf("Help.Update: %#v", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
		v.model.Height = msg.Height - v.model.Style.GetVerticalFrameSize()
		v.model.Width = msg.Width - v.model.Style.GetHorizontalFrameSize()
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
		Height(v.height).
		Width(v.width).
		Padding(0, 1).
		Background(types.Theme.Bg)

	return s.Render(out)
}
