// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package offset

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// calc is a app-wide offset calculator. To initialize it, call Initialize().
var calc *calculator

const (
	// Have to use ansi escape codes to ensure lipgloss doesn't consider ID's as
	// part of the width of the view.

	identStart    = '\x1B' // ANSI escape code.
	identStartLen = len(string(identStart))
	identEnd      = '\x9C' // ANSI termination code.
	identEndLen   = len(string(identEnd))

	areaStart = "_bound_start"
	areaEnd   = "_bound_end"
)

// Coords is a struct that holds the X and Y coordinates of an offset identifier.
type Coords struct {
	id string // raw id of the offset.
	X  int    // X coordinate, starting from 0 (width).
	Y  int    // Y coordinate, starting from 0 (height).
}

// IsZero returns true if Coords doesn't reference an offset. Useful when calling
// Get() using an ID that hasn't been registered to have an offset yet.
func (xy Coords) IsZero() bool {
	return xy.id == ""
}

// InBounds returns true if the mouse event is within the bounds of the given
// components dimensions. Returns false if Coords references an ID that hasn't
// been set yet.
func (xy Coords) InBounds(h, w int, e tea.MouseMsg) bool {
	if xy.IsZero() {
		return false
	}

	return e.X >= xy.X && e.X < xy.X+w && e.Y >= xy.Y && e.Y < xy.Y+h
}

type Area struct {
	Start Coords
	End   Coords
}

func (a Area) IsZero() bool {
	return a.Start.IsZero() || a.End.IsZero()
}

func (a Area) InBounds(e tea.MouseMsg) bool {
	if a.IsZero() {
		return false
	}

	if a.Start.X > a.End.X || a.Start.Y > a.End.Y {
		return false
	}

	if e.X < a.Start.X || e.Y < a.Start.Y {
		return false
	}

	if e.X >= a.End.X || e.Y > a.End.Y {
		return false
	}

	return true
}

func (a Area) Pos(msg tea.MouseMsg) (x, y int) {
	return msg.X - a.Start.X, msg.Y - a.Start.Y
}

// calculator holds the state of the offset calculator, including ID mappings and
// offsets of components.
type calculator struct {
	ctx     context.Context
	cancel  func()
	setChan chan Coords

	mapMu   sync.RWMutex
	mapping map[string]Coords

	idMu      sync.RWMutex
	idCounter int64
	ids       map[string]string // user ID -> generated control sequence ID.
	rids      map[string]string // generated control sequence ID -> user ID.
}

// Initialize initializes the offset calculator globally, starting a goroutine
// to update mapped component offsets. Close() can be called to stop the worker.
func Initialize() {
	if calc != nil {
		return
	}

	calc = &calculator{
		setChan:   make(chan Coords, 200),
		mapping:   make(map[string]Coords),
		ids:       make(map[string]string),
		rids:      make(map[string]string),
		idCounter: 500,
	}

	calc.ctx, calc.cancel = context.WithCancel(context.Background())
	go worker()
}

// Close stops the offset calculator worker.
func Close() {
	if calc == nil {
		return
	}
	calc.cancel()
}

// ID generates an injectable ID for the given component. The resulting ID can be
// stored and re-used as long as the input ID hasn't changed. The ID uses ANSI
// escape codes and an incrementing counter to ensure it doesn't conflict with
// lipgloss's width checks.
func ID(id string) string {
	if calc == nil {
		panic("offset: not initialized")
	}

	calc.idMu.RLock()
	if id, ok := calc.ids[id]; ok {
		calc.idMu.RUnlock()
		return id
	}
	calc.idMu.RUnlock()

	calc.idMu.Lock()
	calc.idCounter++
	counter := fmt.Sprint(calc.idCounter)
	calc.ids[id] = string(identStart) + counter + string(identEnd)
	calc.rids[counter] = id
	calc.idMu.Unlock()
	return calc.ids[id]
}

// Clear removes any stored offsets for the given component ID.
func Clear(id string) {
	calc.mapMu.Lock()
	delete(calc.mapping, id)
	calc.mapMu.Unlock()
}

// Get returns the offset of the given component ID. If the ID is not known (yet),
// Get() returns a Coords with IsZero() == true.
func Get(id string) (xy Coords) {
	if calc == nil {
		panic("offset: not initialized")
	}

	calc.mapMu.RLock()
	xy = calc.mapping[id]
	calc.mapMu.RUnlock()
	return xy
}

// GetReverse returns the component ID from a generated ID (that includes ANSI
// escape codes).
func GetReverse(id string) (resolved string) {
	calc.idMu.RLock()
	resolved = calc.rids[id]
	calc.idMu.RUnlock()
	return resolved
}

func AreaID(id, v string) string {
	return ID(id+areaStart) + v + ID(id+areaEnd)
}

func GetArea(id string) (a Area) {
	a.Start = Get(id + areaStart)
	a.End = Get(id + areaEnd)

	return a
}

func worker() {
	for {
		select {
		case <-calc.ctx.Done():
			return
		case xy := <-calc.setChan:
			// log.Warnf("offset: %#v", xy)
			calc.mapMu.Lock()
			calc.mapping[GetReverse(xy.id)] = xy
			calc.mapMu.Unlock()
		}
	}
}

// Scan will scan the view output, searching for offset identifiers, returning the
// original view output with the offset identifiers stripped. Scan() should be used
// by the outer most model/component of your application, and not inside of a
// model/component child.
//
// Scan buffers the offsets to be stored, so an immediate call to Get(id) may not
// return the correct offset. Thus it's recommended to primarily use Get(id) for
// actions like mouse events, which don't occur immediately after a view shift
// (where the previously stored offset might be different).
func Scan(v string) string {
	if calc == nil {
		panic("offset: not initialized")
	}

	vLen := len(v)
	start := -1
	var end, i, w, newlines, lastNewline, width int
	var id string
	var r rune

	for {
		if i+1 >= vLen {
			return v
		}

		r, w = utf8.DecodeRuneInString(v[i:])

		switch r {
		case utf8.RuneError:
			i += w // Skip invalid rune.
			continue
		case identStart:
			start = i
			i += w
		case identEnd:
			i += w
			if start == -1 {
				continue
			}
			end = i

			id = v[start+identStartLen : end-identEndLen]
			if !isID(id) {
				continue
			}

			// calculate the offset here.
			newlines = strings.Count(v[:start], "\n")
			lastNewline = strings.LastIndex(v[:start], "\n")
			if lastNewline == -1 {
				lastNewline = 0
			}
			width = lipgloss.Width(v[lastNewline:start])

			calc.setChan <- Coords{id: id, X: width, Y: newlines}
			v = v[:start] + v[end:]
			i = start
			vLen = len(v)

			start = -1
		default:
			i += w
		}
	}
}

func isID(rid string) bool {
	for _, r := range rid {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
