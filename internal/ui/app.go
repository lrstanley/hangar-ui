// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package ui

import (
	"context"

	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
	"github.com/lrstanley/hangar-ui/internal/ui/view"
	"github.com/lrstanley/hangar-ui/internal/x"
	"github.com/muesli/termenv"
)

// App is the interface that wraps the types.App interface.
var _ types.App = &App{}

type App struct {
	// Core.
	cli    *clix.CLI[types.Flags]
	logger log.Interface
	keys   *model.KeyMap

	height int
	width  int

	// App-level models.
	commandbar *model.CommandBar
	navbar     *model.NavBar
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

	zone.NewGlobal()

	a := &App{
		cli:    cli,
		logger: log.WithField("src", "app"),

		focused:  types.ViewRoot,
		active:   types.ViewRoot,
		previous: types.ViewRoot,
		views:    map[types.Viewable]view.View{},
	}

	a.keys = model.NewKeyMap(a)

	a.commandbar = model.NewCommandBar(a)
	a.navbar = model.NewNavBar(a, []types.Viewable{
		types.ViewRoot,
		types.ViewTargets,
		types.ViewHelp,
	})
	a.statusbar = model.NewStatusBar(a, a.keys)

	a.views[types.ViewRoot] = view.NewRoot(a)
	a.views[types.ViewHelp] = view.NewHelp(a, a.keys)
	a.views[types.ViewTargets] = view.NewTargets(a)

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
	a.logger.Debugf("msg: %#v", msg)

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
			if a.IsFocused(types.ViewHelp) {
				a.Back(true)
			} else {
				a.SetActive(types.ViewHelp, true)
			}
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
		switch msg.Type {
		case tea.MouseWheelUp, tea.MouseWheelDown:
			_, cmd = a.views[a.active].Update(msg)
			return a, cmd
		}
		return a.propagateMessage(msg)
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
		componentMsg = tea.WindowSizeMsg{Height: a.navbar.Height, Width: a.width}
	}
	_, cmd = a.navbar.Update(componentMsg)
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
	switch msg.(type) {
	case tea.MouseMsg:
		if _, cmd = a.views[a.active].Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	default:
		for _, v := range a.views {
			if _, cmd = v.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
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
		a.navbar.Height -
		a.statusbar.Height

	w = a.width

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

	return zone.Scan(x.Y(
		lipgloss.Top,
		a.commandbar.View(),
		a.navbar.View(),
		v, a.statusbar.View(),
	))
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
	if v != a.active {
		_, _ = a.Update(types.ViewChangeMsg{View: v})
	}

	a.previous = a.active
	a.active = v

	if focused {
		a.SetFocused(v)
	}

	if a.previous == types.ViewHelp {
		a.previous = types.ViewRoot
	}
}

func (a *App) Active() types.Viewable {
	return a.active
}

func (a *App) Previous() types.Viewable {
	return a.previous
}

func (a *App) Back(focused bool) {
	a.active, a.previous = a.previous, a.active

	if a.active == a.previous {
		a.active = types.ViewRoot
	}

	if focused {
		a.SetFocused(a.active)
	}
	_, _ = a.Update(types.ViewChangeMsg{View: a.active})
}

func (a *App) Init() tea.Cmd {
	return nil
}
