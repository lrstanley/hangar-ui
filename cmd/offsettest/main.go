// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/ui/offset"
)

func main() {
	offset.Initialize()

	s := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#383838"))

	test := offset.Calc.ID("1234567111111111111111111") + "testing"
	fmt.Println("width:", lipgloss.Width(offset.Calc.Scan(test)))
	fmt.Println("width:", lipgloss.Width("testing"))
	fmt.Printf("%#v\n", test)
	out := offset.Calc.Scan("test " + s.Render(offset.Calc.ID("foo")+"testing") + " test")

	// fmt.Println("\003")

	fmt.Println(out)
	fmt.Printf("%#v\n", out)

	time.Sleep(51 * time.Millisecond)
	fmt.Printf("%#v\n", offset.Calc.Get("foo"))
}
