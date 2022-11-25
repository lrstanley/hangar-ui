// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package view

import (
	"time"

	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/lrstanley/hangar-ui/internal/api"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
)

const (
	colKeyTargetName     = "name"
	colKeyClusterName    = "cluster_name"
	colKeyTargetURL      = "url"
	colKeyTargetTeam     = "team"
	colKeyTargetInsecure = "insecure"
	colKeyClusterVersion = "cluster_version"
)

type Targets struct {
	*Base
	model model.Table

	targetInfoCache api.TargetInfoMsg
}

func NewTargets(app types.App) *Targets {
	v := &Targets{
		Base: &Base{
			app:    app,
			is:     types.ViewTargets,
			logger: log.WithField("src", "targets"),
		},
		model: model.NewTable(app, types.ViewTargets, []table.Column{
			table.NewFlexColumn(colKeyTargetName, "Target Name", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyClusterName, "Cluster Name", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetURL, "API URL", 3).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetTeam, "Team", 2).WithFiltered(true),
			table.NewFlexColumn(colKeyTargetInsecure, "Insecure", 1),
			table.NewFlexColumn(colKeyClusterVersion, "Version", 1),
		}, colKeyTargetName),
	}

	return v
}

func (v *Targets) UpdateRows() {
	var rows []table.Row
	var row table.RowData

	for _, data := range v.targetInfoCache {
		row = table.RowData{
			colKeyTargetName:     string(data.TargetName),
			colKeyTargetURL:      data.Target.URL(),
			colKeyTargetTeam:     data.Target.Team().Name(),
			colKeyTargetInsecure: table.NewStyledCell("false", lipgloss.NewStyle().Foreground(types.Theme.SuccessFg)),
		}

		if data.TargetName == api.Manager.ActiveName() {
			row[colKeyTargetName] = string(data.TargetName) + " (active)"
		}

		if data.Target.TLSConfig() == nil || data.Target.TLSConfig().InsecureSkipVerify {
			row[colKeyTargetInsecure] = table.NewStyledCell("true", lipgloss.NewStyle().Foreground(types.Theme.FailureFg))
		}

		if data.Error == nil {
			row[colKeyClusterName] = data.Info.ClusterName
			row[colKeyClusterVersion] = data.Info.Version
		} else {
			row[colKeyClusterName] = table.NewStyledCell(data.Error.Error(), lipgloss.NewStyle().Foreground(types.Theme.FailureFg))
		}

		if data.TargetName == api.Manager.ActiveName() {
			rows = append(rows, table.NewRow(row).WithStyle(lipgloss.NewStyle().Bold(true)))
		} else {
			rows = append(rows, table.NewRow(row))
		}
	}

	v.model.UpdateRows(rows)
}

func (v *Targets) Init() tea.Cmd {
	return tea.Batch(
		v.model.Init(),
		api.Manager.QueryTargetInfo,
	)
}

func (v *Targets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, types.KeyEnter):
			api.Manager.SetActive(v.model.SelectedRow().Data[colKeyTargetName].(string))
			v.UpdateRows()
			return v, nil
		case key.Matches(msg, types.KeyRefresh):
			return v, api.Manager.QueryTargetInfo
		}
	case types.ViewChangeMsg:
		if msg.View == v.is {
			return v, api.Manager.QueryTargetInfo
		}
	case api.TargetInfoMsg:
		v.targetInfoCache = msg
		v.UpdateRows()

		if v.Focused() {
			return v, types.DelayCmd(10*time.Second, api.Manager.QueryTargetInfo)
		}
		return v, nil
	}

	var cmd tea.Cmd
	v.model, cmd = v.model.Update(msg)
	return v, cmd
}

func (v *Targets) View() string {
	return v.model.View()
}
