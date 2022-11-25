// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/concourse/concourse/atc"
)

type PipelineListMsg struct {
	Pipelines []atc.Pipeline
	Error     error
}

func (c *apiManager) QueryPipelines() tea.Msg {
	defer c.Loading("fetching pipelines")()

	p, err := c.Active().Client().ListPipelines()

	c.logger.WithFields(log.Fields{
		"pipelines": len(p),
		"error":     err,
	}).Debug("queried pipeline list")

	return PipelineListMsg{Pipelines: p, Error: err}
}
