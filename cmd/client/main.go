// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/api"
	"github.com/lrstanley/hangar-ui/internal/types"
	"github.com/lrstanley/hangar-ui/internal/ui"
)

var (
	logger log.Interface

	cli = &clix.CLI[types.Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/hangar-ui", "master", "https://liam.sh"),
	}
)

func main() {
	cli.LoggerConfig.Quiet = true
	cli.Parse()

	logger = cli.Logger
	ctx := log.NewContext(context.Background(), logger)

	if !cli.LoggerConfig.Quiet && cli.LoggerConfig.Path == "" {
		fmt.Println("log path is required if logging is enabled")
	}

	types.SetTheme("default")

	api.NewAPIClient(ctx, cli)
	prog := tea.NewProgram(ui.New(ctx, cli), tea.WithAltScreen(), tea.WithMouseCellMotion())

	go api.Manager.HandleMsg(func(cmd types.FlyMsg) {
		logger.WithField("msg", fmt.Sprintf("%T", cmd)).WithField("cmd", cmd).Debug("received message from api client")
		prog.Send(cmd)
	})

	if _, err := prog.Run(); err != nil {
		logger.WithError(err).Fatal("failed to start hangar-ui")
	}
}
