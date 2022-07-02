// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"io"
	"log"

	llog "github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui"
)

var (
	logger llog.Interface

	cli = &clix.CLI[types.Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/hangar-ui", "master", "https://liam.sh"),
	}
)

func main() {
	cli.Parse()
	logger = cli.Logger
	ctx := llog.NewContext(context.Background(), logger)

	if cli.Debug {
		// TODO: map to log.Logger.
		if f, err := tea.LogToFile("debug.log", "hangar-ui"); err != nil {
			logger.WithError(err).Fatal("failed to log to file")
		} else {
			defer f.Close()
		}
	} else {
		log.SetOutput(io.Discard)
	}

	types.SetTheme("default")

	app := ui.New(ctx, cli)

	if err := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		logger.WithError(err).Fatal("failed to start hangar-ui")
	}
}
