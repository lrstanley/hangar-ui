// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/apex/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/concourse/concourse/fly/rc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/types"
)

var Manager *apiManager

type apiManager struct {
	ctx      context.Context
	logger   log.Interface
	config   *clix.CLI[types.Flags]
	signaler chan tea.Msg

	cancelFn func()
	wg       sync.WaitGroup

	targets           types.Atomic[rc.Targets]
	currentTarget     types.Atomic[rc.Target]
	currentTargetName types.Atomic[string]
}

func NewAPIClient(ctx context.Context, config *clix.CLI[types.Flags]) {
	Manager = &apiManager{
		logger:   log.WithField("src", "api-client"),
		config:   config,
		signaler: make(chan tea.Msg, 50),
	}

	Manager.ctx, Manager.cancelFn = context.WithCancel(ctx)

	Manager.UpdateTargets()

	if config.Flags.Target != "" {
		if err := Manager.SetActive(config.Flags.Target); err != nil {
			Manager.logger.WithError(err).WithField("target", config.Flags.Target).Fatal("failed to configure target")
		}
	} else {
		targets := Manager.TargetNames()
		if len(targets) == 0 {
			Manager.logger.Fatal("no targets found, please setup one with `fly -t <target> login`")
			return
		}

		if err := Manager.SetActive(targets[0]); err != nil {
			Manager.logger.WithError(err).WithField("target", targets[0]).Fatal("failed to configure target")
			return
		}
	}

	Manager.wg.Add(1)
	go Manager.Watcher()
}

// HandleMsg allows the user to handle a message from the signaler channel.
func (c *apiManager) HandleMsg(cb func(tea.Msg)) {
	var cmd tea.Msg
	for {
		select {
		case <-c.ctx.Done():
			return
		case cmd = <-c.signaler:
			cb(cmd)
		}
	}
}

// UpdateTargets updates the targets list from the fly config file (flyrc).
func (c *apiManager) UpdateTargets() {
	targets, err := rc.LoadTargets()
	if err != nil {
		c.logger.WithError(err).Fatal("failed to load targets from flyrc")
		return
	}

	c.targets.Store(targets)
	c.signaler <- types.FlyTargetsUpdated
}

// Targets returns the current known targets list.
func (c *apiManager) Targets() rc.Targets {
	return c.targets.Load()
}

// TargetNames returns a []string of target names (as one would pass into --target)
func (c *apiManager) TargetNames() []string {
	var names []string
	for k := range c.targets.Load() {
		names = append(names, string(k))
	}

	sort.Strings(names)
	return names
}

// Active returns the active target.
func (c *apiManager) Active() rc.Target {
	return c.currentTarget.Load()
}

// ActiveName returns the name of the active target.
func (c *apiManager) ActiveName() string {
	return c.currentTargetName.Load()
}

// Client returns a concourse.Client for the active target.
func (c *apiManager) Client() concourse.Client {
	return c.currentTarget.Load().Client()
}

// SetActive sets the active target to the given target name. This will fail if
// the target credentials are outdated.
func (c *apiManager) SetActive(targetName string) error {
	c.logger.WithField("target", targetName).Debug("setting active target")

	target, err := rc.LoadTarget(rc.TargetName(targetName), false)
	if err != nil {
		return err
	}

	c.currentTarget.Store(target)
	c.currentTargetName.Store(targetName)

	c.signaler <- types.FlyActiveTargetUpdated

	return nil
}

// Close closes the api client, and signals to stop the background worker.
func (c *apiManager) Close() {
	c.cancelFn()
	c.wg.Wait()
}

// Watcher is a background worker that will periodically update the API client
// using the flyrc, and other forms of configuration.
func (c *apiManager) Watcher() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			c.UpdateTargets()
		}
	}
}
