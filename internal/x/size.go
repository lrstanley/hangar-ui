// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package x

import "github.com/charmbracelet/lipgloss"

var (
	// (W) Width returns the cell width of characters in the string. ANSI sequences
	// are ignored and characters wider than one cell (such as Chinese characters
	// and emojis) are appropriately measured.
	//
	// You should use this instead of len(string) len([]rune(string) as neither
	// will give you accurate results.
	W = lipgloss.Width

	// (H) Height returns height of a string in cells. This is done simply by
	// counting \n characters. If your strings use \r\n for newlines you should
	// convert them to \n first, or simply write a separate function for measuring
	// height.
	H = lipgloss.Height
)

func HMulti(input ...string) (h int) {
	for _, s := range input {
		h += H(s)
	}
	return
}

func WMulti(input ...string) (w int) {
	for _, s := range input {
		w += W(s)
	}
	return
}
