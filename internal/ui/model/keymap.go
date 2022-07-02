// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package model

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/lrstanley/hangar-ui/internal/types"
)

var _ help.KeyMap = (*KeyMap)(nil) // Validate interface.

type KeyMap struct {
	app types.App

	Binds map[types.Viewable][]key.Binding
}

func NewKeyMap(app types.App) *KeyMap {
	return &KeyMap{
		app: app,
		Binds: map[types.Viewable][]key.Binding{
			types.ViewRoot: {
				types.KeyHelp,
				types.KeyQuit,
			},
			types.ViewHelp: {
				types.KeyCancel,
			},
			types.ViewCommandBar: {
				types.KeyCancel,
				types.KeyCmdFilter,
				types.KeyCmdInvoke,
				types.KeyEnter,
			},
		},
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *KeyMap) ShortHelp() (kb []key.Binding) {
	if active := k.app.Active(); active != types.ViewRoot {
		kb = append(kb, k.Binds[k.app.Active()]...)
	}
	kb = append(kb, k.Binds[types.ViewRoot]...)
	return kb
}

func (k *KeyMap) FullHelp() (kb [][]key.Binding) {
	// kb = append(kb, k.binds[types.ViewRoot])
	// if active := k.app.Active(); active != types.ViewRoot {
	// 	kb = append(kb, k.binds[k.app.Active()])
	// }
	// return kb

	return nil
}
