// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/fly/rc"
)

type TargetInfo struct {
	TargetName string
	Target     rc.Target

	Info  atc.Info
	Error error
}

type TargetInfoMsg []TargetInfo

// QueryTargetInfo queries all known targets for their status, and returns
// api.TargetInfoMsg.
func (m *apiManager) QueryTargetInfo(delay time.Duration) func() tea.Msg {
	return func() tea.Msg {
		time.Sleep(delay)
		msg := make(TargetInfoMsg, 0)

		targets := m.Targets()

		for name := range targets {
			t, err := rc.LoadTarget(name, false)
			if err != nil {
				continue
			}

			info, err := t.Client().GetInfo()

			msg = append(msg, TargetInfo{
				TargetName: string(name),
				Target:     t,
				Info:       info,
				Error:      err,
			})
		}

		m.logger.Debugf("%#v", msg)

		return msg
	}
}
