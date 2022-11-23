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
	"github.com/concourse/concourse/fly/rc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/hangar-ui/internal/types"
)

var Root *apiClient

type apiClient struct {
	ctx      context.Context
	logger   log.Interface
	config   *clix.CLI[types.Flags]
	signaler chan types.FlyMsg

	cancelFn func()
	wg       sync.WaitGroup

	targets           types.Atomic[rc.Targets]
	currentTarget     types.Atomic[rc.Target]
	currentTargetName types.Atomic[string]
}

func NewAPIClient(ctx context.Context, config *clix.CLI[types.Flags]) {
	Root = &apiClient{
		logger:   log.WithField("src", "api-client"),
		config:   config,
		signaler: make(chan types.FlyMsg, 50),
	}

	Root.ctx, Root.cancelFn = context.WithCancel(ctx)

	Root.UpdateTargets()

	if config.Flags.Target != "" {
		if err := Root.SetActive(config.Flags.Target); err != nil {
			Root.logger.WithError(err).WithField("target", config.Flags.Target).Fatal("failed to configure target")
		}
	} else {
		targets := Root.TargetNames()
		if len(targets) == 0 {
			Root.logger.Fatal("no targets found, please setup one with `fly -t <target> login`")
			return
		}

		if err := Root.SetActive(targets[0]); err != nil {
			Root.logger.WithError(err).WithField("target", targets[0]).Fatal("failed to configure target")
			return
		}
	}

	Root.wg.Add(1)
	go Root.Watcher()
}

func (c *apiClient) HandleSignal(cb func(types.FlyMsg)) {
	var cmd types.FlyMsg
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
func (c *apiClient) UpdateTargets() {
	targets, err := rc.LoadTargets()
	if err != nil {
		c.logger.WithError(err).Fatal("failed to load targets from flyrc")
		return
	}

	c.targets.Store(targets)
	c.signaler <- types.FlyTargetsUpdated
}

func (c *apiClient) Targets() rc.Targets {
	return c.targets.Load()
}

// TargetNames returns a []string of target names (as one would pass into --target)
func (c *apiClient) TargetNames() []string {
	var names []string
	for k := range c.targets.Load() {
		names = append(names, string(k))
	}

	sort.Strings(names)
	return names
}

// Active returns the active target.
func (c *apiClient) Active() rc.Target {
	return c.currentTarget.Load()
}

// ActiveName returns the name of the active target.
func (c *apiClient) ActiveName() string {
	return c.currentTargetName.Load()
}

// Client returns a concourse.Client for the active target.
func (c *apiClient) Client() concourse.Client {
	return c.currentTarget.Load().Client()
}

// SetActive sets the active target to the given target name. This will fail if
// the target credentials are outdated.
func (c *apiClient) SetActive(targetName string) error {
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
func (c *apiClient) Close() {
	c.cancelFn()
	c.wg.Wait()
}

// Watcher is a background worker that will periodically update the API client
// using the flyrc, and other forms of configuration.
func (c *apiClient) Watcher() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			c.UpdateTargets()
		}
	}
}
