// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"sync"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/fly/rc"
)

type Target struct {
	Name string

	mu sync.RWMutex

	config rc.Target

	info      atc.Info
	infoError error
}

func (t *Target) fetchInfo() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.info, t.infoError = t.config.Client().GetInfo()
}

func (t *Target) Info() (atc.Info, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.info, t.infoError
}
