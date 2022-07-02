// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package ui

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
	"github.com/lrstanley/hangar-ui/internal/ui/view"
	"github.com/lrstanley/hangar-ui/internal/x"
	"github.com/muesli/termenv"
)

const (
	viewBorderHeight = 1
	viewBorderWidth  = 1
)

// App is the interface that wraps the types.App interface.
var _ types.App = &App{}

type App struct {
	// Core.
	cli  *clix.CLI[types.Flags]
	keys *model.KeyMap

	height int
	width  int

	// App-level models.
	commandbar *model.CommandBar
	statusbar  *model.StatusBar

	// State related items.
	focused  types.Viewable
	active   types.Viewable
	previous types.Viewable
	views    map[types.Viewable]view.View
}

func New(_ context.Context, cli *clix.CLI[types.Flags]) *App {
	// See: https://github.com/charmbracelet/lipgloss/issues/73
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	a := &App{
		cli: cli,

		focused:  types.ViewRoot,
		active:   types.ViewRoot,
		previous: types.ViewRoot,
		views:    map[types.Viewable]view.View{},
	}

	a.keys = model.NewKeyMap(a)

	a.commandbar = model.NewCommandBar(a)
	a.statusbar = model.NewStatusBar(a, a.keys)

	a.views[types.ViewRoot] = view.NewRoot(a)
	a.views[types.ViewHelp] = view.NewHelp(a, a.keys)

	// Send initial sizes to all views.
	vh, vw := a.getViewSize()
	for _, v := range a.views {
		_, _ = v.Update(tea.WindowSizeMsg{
			Height: vh,
			Width:  vw,
		})
	}
	return a
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("App.Update: %#v", msg)
	var cmd tea.Cmd

	// Only do updates if we're initialized, or it's a window size message, which
	// allows us to initialize.
	if _, ok := msg.(tea.WindowSizeMsg); !ok && !a.isInitialized() {
		return a, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width

		// Repurpose the message, and adjust for app view sizes, then propagate
		// to children.
		msg.Height, msg.Width = a.getViewSize()
		return a.propagateMessage(msg)
	case tea.KeyMsg:
		cmdFocused := a.IsFocused(types.ViewCommandBar)

		switch {
		case key.Matches(msg, types.KeyCmdFilter) && !cmdFocused:
			a.SetFocused(types.ViewCommandBar)
			_, _ = a.commandbar.Update(model.MsgCmdFilter)
			return a, nil
		case key.Matches(msg, types.KeyCmdInvoke) && !cmdFocused:
			a.SetFocused(types.ViewCommandBar)
			_, _ = a.commandbar.Update(model.MsgCmdInvoke)
			return a, nil
		case key.Matches(msg, types.KeyQuit):
			return a, tea.Quit
		case key.Matches(msg, types.KeyHelp) && !cmdFocused:
			a.SetActive(types.ViewHelp, true)
			return a, nil
		}

		if cmdFocused {
			_, cmd = a.commandbar.Update(msg)
			return a, cmd
		}

		// Don't propagate key messages to anything but the active view.
		_, cmd = a.views[a.active].Update(msg)
		return a, cmd
	case tea.MouseMsg:
		// Explicitly handle some specific mouse events.
		switch msg.Type {
		case tea.MouseWheelUp, tea.MouseWheelDown:
			_, cmd = a.views[a.active].Update(msg)
			return a, cmd
		default:
			// Check to see if the mouse is over the commandbar, or statusbar.

			if msg.Y < a.commandbar.Height {
				a.SetFocused(types.ViewCommandBar)
				_, cmd = a.commandbar.Update(msg)
				return a, cmd
			}

			if msg.Y >= a.height-a.statusbar.Height {
				a.SetFocused(types.ViewStatusBar)
				_, cmd = a.statusbar.Update(msg)
				return a, cmd
			}

			if a.IsFocused(types.ViewCommandBar) { // If the command bar was focused, unfocus it.
				a.SetFocused(a.active)
				a.commandbar.Update(model.MsgNone)
			}

			minYBounds := a.commandbar.Height + viewBorderHeight
			maxYBounds := a.height - a.statusbar.Height - viewBorderHeight - 1
			minXBounds := viewBorderWidth
			maxXBounds := a.width - viewBorderWidth - 1

			if msg.Y >= minYBounds && msg.Y <= maxYBounds && msg.X >= minXBounds && msg.X <= maxXBounds {
				// Don't propagate mouse events to anything but the active view.
				msg.X -= minXBounds
				msg.Y -= minYBounds
				a.SetFocused(a.active)
				_, cmd = a.views[a.active].Update(msg)
				return a, cmd
			}

			return a, nil
		}
	case types.ViewMsg: // A message for a specific view, propagated from a child.
		_, cmd = a.views[msg.View].Update(msg.Msg)
		return a, cmd
	}

	return a.propagateMessage(msg)
}

func (a *App) propagateMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	componentMsg := msg

	// Update app components here.
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		componentMsg = tea.WindowSizeMsg{Height: a.commandbar.Height, Width: a.width}
	}
	_, cmd = a.commandbar.Update(componentMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if _, ok := msg.(tea.WindowSizeMsg); ok {
		componentMsg = tea.WindowSizeMsg{Height: a.statusbar.Height, Width: a.width}
	}
	_, cmd = a.statusbar.Update(componentMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// All other components.
	for _, v := range a.views {
		if _, cmd = v.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return a, tea.Batch(cmds...)
	}

	return a, nil
}

func (a *App) getViewSize() (h, w int) {
	h = a.height -
		a.commandbar.Height -
		a.statusbar.Height -
		(2 * viewBorderHeight)

	w = a.width -
		(2 * viewBorderWidth)

	return h, w
}

func (a *App) isInitialized() bool {
	return a.height > 0 && a.width > 0
}

func (a *App) View() string {
	if !a.isInitialized() {
		return ""
	}

	v := a.views[a.active].View()
	s := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBackground(types.Theme.ViewBorderBg).
		BorderForeground(types.Theme.ViewBorderInactiveFg)

	if a.IsFocused(a.active) {
		s = s.BorderForeground(types.Theme.ViewBorderActiveFg)
	}

	return x.Y(lipgloss.Top, a.commandbar.View(), s.Render(v), a.statusbar.View())
}

func (a *App) SetFocused(v types.Viewable) {
	if a.focused != v {
		_, _ = a.Update(types.FocusChangeMsg{View: v})
	}

	a.focused = v
}

func (a *App) IsFocused(v types.Viewable) bool {
	return a.focused == v
}

func (a *App) SetActive(v types.Viewable, focused bool) {
	if a.active == v {
		a.Back(focused)
		return
	}

	a.previous = a.active
	a.active = v
	_, _ = a.Update(types.ViewChangeMsg{View: v})

	if focused {
		a.SetFocused(v)
	}
}

func (a *App) Active() types.Viewable {
	return a.active
}

func (a *App) Previous() types.Viewable {
	return a.previous
}

func (a *App) Back(focused bool) {
	a.previous = types.ViewRoot
	a.active = types.ViewRoot
	_, _ = a.Update(types.ViewChangeMsg{View: a.active})

	if focused {
		a.SetFocused(types.ViewRoot)
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}
