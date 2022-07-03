// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
)

const (
	MsgCmdFilter Msg = iota + 1
	MsgCmdInvoke
)

type CommandBar struct {
	*Base

	input       textinput.Model
	method      Msg
	style       lipgloss.Style
	prefixStyle lipgloss.Style
}

func NewCommandBar(app types.App) *CommandBar {
	m := &CommandBar{
		Base: &Base{
			app:    app,
			is:     types.ViewCommandBar,
			Height: 3,
		},
		input: textinput.New(),
	}

	m.input.Placeholder = "[/] to filter, [:] to search commands"
	m.input.PlaceholderStyle = m.input.PlaceholderStyle.Background(types.Theme.Bg).Foreground(types.Theme.InputPlaceholderFg)
	m.input.PromptStyle = m.input.PromptStyle.Background(types.Theme.Bg).Foreground(types.Theme.InputFg)
	m.input.TextStyle = m.input.TextStyle.Background(types.Theme.Bg).Foreground(types.Theme.InputFg)
	m.input.CursorStyle = m.input.CursorStyle.Background(types.Theme.Bg).Foreground(types.Theme.InputCursorFg)

	m.style = lipgloss.NewStyle().
		Background(types.Theme.Bg).
		BorderBackground(types.Theme.ViewBorderBg).
		Border(lipgloss.RoundedBorder())

	m.prefixStyle = lipgloss.NewStyle().
		Background(types.Theme.Bg).
		Foreground(types.Theme.Fg)

	return m
}

func (m *CommandBar) toggle() {
	if m.method == MsgCmdFilter {
		m.method = MsgCmdInvoke
	} else {
		m.method = MsgCmdFilter
	}
}

func (m *CommandBar) Init() tea.Cmd {
	return nil
}

func (m *CommandBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("CommandBar.Update: %#v", msg)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style.Width(msg.Width - 2)
		m.Width = msg.Width
		m.input.Width = msg.Width - 8
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseRight: // TODO: right click to change mode, and left to change cursor position?
			m.app.SetFocused(m.is)
			m.toggle()
		case tea.MouseLeft:
			m.app.SetFocused(m.is)

			msg.X -= 6 // Account for border, command type, cursor prefix, etc.
			msg.Y -= 1

			m.input.SetCursor(msg.X)
			return m, nil
		}
	case Msg:
		switch msg {
		case MsgCmdFilter, MsgCmdInvoke:
			m.method = msg
			_ = m.input.Reset()
			m.input.Focus()
		case MsgNone:
			m.method = msg
			m.input.Blur()
		}

		return m, nil
	case types.FocusChangeMsg:
		if m.is != msg.View {
			m.input.Blur()
			return m, nil
		}

		m.toggle()
		m.input.Focus()
	case tea.KeyMsg:
		if !m.input.Focused() {
			return m, nil
		}

		switch {
		case key.Matches(msg, types.KeyCmdFilter) && m.input.Value() == "":
			// If they're in the input, then try and use / or : again, just change
			// the method.
			m.method = MsgCmdFilter
			return m, nil
		case key.Matches(msg, types.KeyCmdInvoke) && m.input.Value() == "":
			// If they're in the input, then try and use / or : again, just change
			// the method.
			m.method = MsgCmdInvoke
			return m, nil
		case key.Matches(msg, types.KeyCancel):
			_ = m.input.Reset()
			m.method = MsgNone
			m.input.Blur()
			m.app.SetFocused(m.app.Active())
			return m, nil
		case key.Matches(msg, types.KeyEnter):
			if m.method == MsgCmdInvoke {
				// TODO: switch to command view.
				// Also forward up/down/enter when active with the command view,
				// to select the necessary command.
				m.method = MsgNone
				_ = m.input.Reset()
				m.app.SetFocused(m.app.Active())
			}

			m.input.Blur()
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m *CommandBar) View() string {
	input := m.input.View()
	s := m.style.Copy()

	if m.Focused() {
		s = s.BorderForeground(types.Theme.ViewBorderActiveFg)
	} else {
		s = s.BorderForeground(types.Theme.ViewBorderInactiveFg)
	}

	var prefix string

	switch m.method {
	case MsgCmdFilter:
		prefix = "?"
	case MsgCmdInvoke:
		prefix = "!"
	default:
		prefix = "-"
	}

	return s.Render(m.prefixStyle.Render("["+prefix+"]") + input)
}
