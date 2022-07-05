// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"math/rand"
	"strconv"

	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/x"
)

const (
	colKeyTargetName    = "name"
	colKeyTargetURL     = "url"
	colKeyTargetTeam    = "team"
	colKeyTargetExpires = "expires"
)

type Targets struct {
	*Base

	model table.Model

	baseStyle lipgloss.Style
}

func NewTargets(app types.App) *Targets {
	m := &Targets{
		Base: &Base{
			app:    app,
			is:     types.ViewTargets,
			logger: log.WithField("src", "targets"),
		},
		model: table.New([]table.Column{
			table.NewFlexColumn(colKeyTargetName, "Name", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetURL, "URL", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetTeam, "Team", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetExpires, "Expires", 1),
		}).HighlightStyle(lipgloss.NewStyle().Background(types.Theme.NavActiveBg).Foreground(types.Theme.NavActiveBg)).
			BorderRounded().
			SortByAsc(colKeyTargetName).
			WithPageSize(1).
			WithPaginationWrapping(false).
			Focused(true).WithHighlightedRow(1).Filtered(true),
	}

	m.baseStyle = lipgloss.NewStyle().
		Foreground(types.Theme.Fg).
		Background(types.Theme.Bg).
		BorderBackground(types.Theme.ViewBorderBg).
		BorderForeground(types.Theme.ViewBorderInactiveFg)

	var rows []table.Row

	testing := []string{"BubbleTea", "Example", "Another thing", "This is a test"}

	for i := 0; i < 100; i++ {
		rows = append(rows, table.NewRow(table.RowData{
			colKeyTargetName:    testing[rand.Intn(len(testing))] + strconv.Itoa(i),
			colKeyTargetURL:     "https://bubbletea.com",
			colKeyTargetTeam:    "BubbleTea " + strconv.Itoa(i),
			colKeyTargetExpires: "Never",
		}))
	}

	m.model = m.model.WithRows(rows)

	return m
}

func (v *Targets) Init() tea.Cmd { return nil }

func (v *Targets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	v.logger.Debugf("msg: %#v", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width

		v.model = v.model.WithTargetWidth(msg.Width)
		v.model = v.model.WithMaxTotalWidth(msg.Width)
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft, tea.MouseRight:
			v.app.SetFocused(v.is)
		case tea.MouseWheelUp:
			v.model = v.model.PageUp()
		case tea.MouseWheelDown:
			v.model = v.model.PageDown()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, types.KeyCancel):
			v.app.Back(true)
			return v, nil
		}
	// TODO: https://github.com/Evertras/bubble-table/issues/116
	case types.FilterMsg:
		v.model = v.model.WithFilterInputValue(msg.Filter)
	}

	var cmd tea.Cmd
	v.model, cmd = v.model.Update(msg)
	return v, cmd
}

func (v *Targets) View() string {
	s := lipgloss.NewStyle().
		Width(v.width).
		Height(v.height).
		MaxHeight(v.height).
		MaxWidth(v.width).
		Background(types.Theme.Bg)

	// - Top/bottom borders +
	// - Top footer border +
	// - Row bottom border +
	// - Header & header footer == 6.
	v.model = v.model.WithPageSize(v.height - 6)
	rows := len(v.model.GetVisibleRows())
	v.model = v.model.WithStaticFooter(x.Expand(v.height - rows - 6))

	baseStyle := v.baseStyle.Copy()
	if v.Focused() {
		baseStyle.BorderForeground(types.Theme.ViewBorderActiveFg)
	}

	// TODO: when pagination is triggered, our padded footer doesn't get used.

	return s.Render(v.model.WithBaseStyle(baseStyle).View())
}
