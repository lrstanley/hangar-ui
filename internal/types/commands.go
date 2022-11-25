// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func DelayCmd(delay time.Duration, cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay)

		return cmd()
	}
}

func DelayMsg(delay time.Duration, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay)

		return msg
	}
}

func MsgAsCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

type FilterMsg struct {
	Filter string
}

func (m FilterMsg) Value() string {
	return m.Filter
}

type FlyEvent int

const (
	FlyTargetsUpdated FlyEvent = iota + 1
	FlyActiveTargetUpdated
)

type LoadingMsg struct {
	Text string
}

type CancelLoadingMsg struct{}
