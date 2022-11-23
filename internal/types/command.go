// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

type FilterMsg struct {
	Filter string
}

func (m FilterMsg) Value() string {
	return m.Filter
}

type FlyMsg int

const (
	FlyTargetsUpdated FlyMsg = iota + 1
	FlyActiveTargetUpdated
)
