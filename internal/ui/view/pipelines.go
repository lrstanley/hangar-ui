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
	"github.com/dustin/go-humanize"
	"github.com/evertras/bubble-table/table"
	"github.com/lrstanley/hangar-ui/internal/api"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui/model"
)

const (
	colPipelineID             = "id"
	colPipelineName           = "name"
	colPipelineInstanceVars   = "instance_vars"
	colPipelinePaused         = "paused"
	colPipelinePublic         = "public"
	colPipelineArchived       = "archived"
	colPipelineTeam           = "team"
	colPipelineLastUpdated    = "last_updated"
	colPipelineLastUpdatedRaw = "last_updated_raw"
)

type Pipelines struct {
	*Base
	model model.Table

	pipelineCache api.PipelineListMsg
}

func NewPipelines(app types.App) *Pipelines {
	v := &Pipelines{
		Base: &Base{
			app:    app,
			is:     types.ViewPipelines,
			logger: log.WithField("src", "pipelines"),
		},
		model: model.NewTable(app, types.ViewPipelines, []table.Column{
			table.NewColumn(colPipelineID, "ID", 4),
			table.NewFlexColumn(colPipelineName, "Name", 5).WithFiltered(true),
			table.NewFlexColumn(colPipelineInstanceVars, "Instance Vars", 4).WithFiltered(true),
			table.NewColumn(colPipelinePaused, "Pause", 5),
			table.NewColumn(colPipelinePublic, "Public", 6),
			table.NewColumn(colPipelineArchived, "Archive", 7),
			table.NewFlexColumn(colPipelineTeam, "Team", 4).WithFiltered(true),
			table.NewFlexColumn(colPipelineLastUpdated, "Last Updated", 2),
		}, colPipelineName),
	}

	return v
}

func (v *Pipelines) UpdateRows() {
	var rows []table.Row
	var row table.RowData

	for _, data := range v.pipelineCache.Pipelines {
		row = table.RowData{
			colPipelineID:             table.NewStyledCell(data.ID, lipgloss.NewStyle().Align(lipgloss.Right)),
			colPipelineName:           data.Name,
			colPipelineInstanceVars:   data.InstanceVars.String(),
			colPipelinePaused:         v.model.Checkmark(data.Paused),
			colPipelinePublic:         v.model.Checkmark(data.Public),
			colPipelineArchived:       v.model.Checkmark(data.Archived),
			colPipelineTeam:           data.TeamName,
			colPipelineLastUpdated:    humanize.Time(time.Unix(data.LastUpdated, 0)),
			colPipelineLastUpdatedRaw: data.LastUpdated,
		}

		rows = append(rows, table.NewRow(row))
	}

	v.model.UpdateRows(rows)
}

func (v *Pipelines) Init() tea.Cmd {
	return tea.Batch(
		v.model.Init(),
		api.Manager.QueryPipelines,
	)
}

func (v *Pipelines) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return v, api.Manager.QueryPipelines
		case key.Matches(msg, types.KeySortName):
			v.model.Sort(colPipelineName)
			return v, nil
		case key.Matches(msg, types.KeySortTime):
			v.model.Sort(colPipelineLastUpdatedRaw)
			return v, nil
		}
	case types.ViewChangeMsg:
		if msg.View == v.is {
			return v, api.Manager.QueryPipelines
		}
	case api.PipelineListMsg:
		v.pipelineCache = msg
		v.UpdateRows()

		if v.Focused() {
			return v, types.DelayCmd(10*time.Second, api.Manager.QueryPipelines)
		}
		return v, nil
	}

	var cmd tea.Cmd
	v.model, cmd = v.model.Update(msg)
	return v, cmd
}

func (v *Pipelines) View() string {
	return v.model.View()
}
