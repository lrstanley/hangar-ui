// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/concourse/concourse/fly/rc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/lrstanley/hangar-ui/internal/types"
)

type apiClient struct {
	app    types.App
	ctx    context.Context
	logger log.Interface

	cancelFn func()
	wg       sync.WaitGroup

	targets       types.Atomic[rc.Targets]
	currentTarget types.Atomic[rc.Target]
}

func NewAPIClient(ctx context.Context, app types.App) *apiClient {
	c := &apiClient{
		app:    app,
		logger: log.WithField("src", "api-client"),
	}

	c.ctx, c.cancelFn = context.WithCancel(ctx)

	c.wg.Add(1)
	go c.Watcher()

	return c
}

func (c *apiClient) Active() rc.Target {
	return c.currentTarget.Load()
}

func (c *apiClient) Client() concourse.Client {
	return c.currentTarget.Load().Client()
}

func (c *apiClient) SetActive(targetName string) error {
	target, err := rc.LoadTarget(rc.TargetName(targetName), false)
	if err != nil {
		return err
	}

	c.currentTarget.Store(target)
	return nil
}

func (c *apiClient) Close() {
	c.cancelFn()
	c.wg.Wait()
}

func (c *apiClient) Watcher() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(10 * time.Second):
			c.targets.Load()
		}
	}
}
