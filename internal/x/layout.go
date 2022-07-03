// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package x

import "github.com/charmbracelet/lipgloss"

var (
	// X is a utility function for horizontally joining two
	// potentially multi-lined strings along a vertical axis. The first argument is
	// the position, with 0 being all the way at the top and 1 being all the way
	// at the bottom.
	//
	// If you just want to align to the left, right or center you may as well just
	// use the helper constants Top, Center, and Bottom.
	//
	// Example:
	//
	//     blockB := "...\n...\n..."
	//     blockA := "...\n...\n...\n...\n..."
	//
	//     // Join 20% from the top
	//     str := lipgloss.JoinHorizontal(0.2, blockA, blockB)
	//
	//     // Join on the top edge
	//     str := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB)
	//
	X = lipgloss.JoinHorizontal
	// Y is a utility function for vertically joining two potentially
	// multi-lined strings along a horizontal axis. The first argument is the
	// position, with 0 being all the way to the left and 1 being all the way to
	// the right.
	//
	// If you just want to align to the left, right or center you may as well just
	// use the helper constants Left, Center, and Right.
	//
	// Example:
	//
	//     blockB := "...\n...\n..."
	//     blockA := "...\n...\n...\n...\n..."
	//
	//     // Join 20% from the top
	//     str := lipgloss.JoinVertical(0.2, blockA, blockB)
	//
	//     // Join on the right edge
	//     str := lipgloss.JoinVertical(lipgloss.Right, blockA, blockB)
	//
	Y = lipgloss.JoinVertical

	// PlaceX places a string or text block horizontally in an unstyled block of
	// a given width. If the given width is shorter than the max width of the
	// string (measured by it's longest line) this will be a noöp.
	PlaceX = lipgloss.PlaceHorizontal

	// PlaceY places a string or text block vertically in an unstyled block of a
	// given height. If the given height is shorter than the height of the string
	// (measured by it's newlines) then this will be a noöp.
	PlaceY = lipgloss.PlaceVertical

	// Place places a string or text block vertically in an unstyled box of a given
	// width or height.
	Place = lipgloss.Place

	// Top is a helper constant for the top edge of a block.
	Top = lipgloss.Top
	// Bottom is a helper constant for the bottom edge of a block.
	Bottom = lipgloss.Bottom
	// Left is a helper constant for the left edge of a block.
	Left = lipgloss.Left
	// Right is a helper constant for the right edge of a block.
	Right = lipgloss.Right
)
