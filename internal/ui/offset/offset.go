// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package offset

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/ansi"
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

	areaStart = "__start"
	areaEnd   = "__end"
)

// Coords is a struct that holds the X and Y coordinates of an offset identifier.
type Coords struct {
	id string // raw id of the offset.
	X  int    // X coordinate, starting from 0 (width).
	Y  int    // Y coordinate, starting from 0 (height).
}

// IsZero returns true if Coords doesn't reference an offset. Useful when calling
// Get() using an ID that hasn't been registered to have an offset yet.
func (xy *Coords) IsZero() bool {
	if xy == nil {
		return true
	}
	return xy.id == ""
}

type Area struct {
	Start *Coords
	End   *Coords
}

func (a *Area) IsZero() bool {
	if a == nil {
		return true
	}
	return a.Start.IsZero() || a.End.IsZero()
}

func (a *Area) InBounds(e tea.MouseMsg) bool {
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

func (a *Area) Pos(msg tea.MouseMsg) (x, y int) {
	return msg.X - a.Start.X, msg.Y - a.Start.Y
}

// calculator holds the state of the offset calculator, including ID mappings and
// offsets of components.
type calculator struct {
	ctx     context.Context
	cancel  func()
	setChan chan *Coords

	mapMu   sync.RWMutex
	mapping map[string]*Coords

	idMu      sync.RWMutex
	idCounter int
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
		setChan:   make(chan *Coords, 200),
		mapping:   make(map[string]*Coords),
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
func ID(id, v string) string {
	if calc == nil {
		panic("offset: not initialized")
	}

	startID := id + areaStart
	endID := id + areaEnd

	calc.idMu.RLock()
	start := calc.ids[startID]
	end := calc.ids[endID]
	calc.idMu.RUnlock()

	if start != "" && end != "" {
		return start + v + end
	}

	calc.idMu.Lock()

	calc.idCounter++
	counter := strconv.Itoa(calc.idCounter)
	start = string(identStart) + counter + string(identEnd)
	calc.ids[startID] = start
	calc.rids[counter] = startID // TODO: should this be counter, or start?

	calc.idCounter++
	counter = strconv.Itoa(calc.idCounter)
	end = string(identStart) + counter + string(identEnd)
	calc.ids[endID] = end
	calc.rids[counter] = endID

	calc.idMu.Unlock()
	return start + v + end
}

// Clear removes any stored offsets for the given ID.
func Clear(id string) {
	calc.mapMu.Lock()
	delete(calc.mapping, id+areaStart)
	delete(calc.mapping, id+areaEnd)
	calc.mapMu.Unlock()
}

// Get returns the offset info of the given ID. If the ID is not known (yet),
// Get() returns a Coords with IsZero() == true.
func Get(id string) (a *Area) {
	if calc == nil {
		panic("offset: not initialized")
	}

	calc.mapMu.RLock()
	a = &Area{
		Start: calc.mapping[id+areaStart],
		End:   calc.mapping[id+areaEnd],
	}
	calc.mapMu.RUnlock()
	return a
}

// getReverse returns the component ID from a generated ID (that includes ANSI
// escape codes).
func getReverse(id string) (resolved string) {
	calc.idMu.RLock()
	resolved = calc.rids[id]
	calc.idMu.RUnlock()
	return resolved
}

func worker() {
	for {
		select {
		case <-calc.ctx.Done():
			return
		case xy := <-calc.setChan:
			calc.mapMu.Lock()
			calc.mapping[getReverse(xy.id)] = xy
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
			if !isNumber(id) {
				continue
			}

			// calculate the offset here.
			newlines = countNewlines(v[:start])
			lastNewline = strings.LastIndex(v[:start], "\n")
			if lastNewline == -1 {
				lastNewline = 0
			}
			width = ansi.PrintableRuneWidth(v[lastNewline:start])

			calc.setChan <- &Coords{id: id, X: width, Y: newlines}
			v = v[:start] + v[end:]
			i = start
			vLen = len(v)

			start = -1
		default:
			i += w
		}
	}
}

func isNumber(rid string) bool {
	for _, r := range rid {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
