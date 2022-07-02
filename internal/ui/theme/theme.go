// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package theme

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/hangar-ui/internal/types"
)

var (
	CurrentScheme = types.NewAtomicComparable("Gotham")
	Schemes       = map[string]*Scheme{}

	//go:embed themes.json
	data []byte
)

type raw struct {
	Name   string   `json:"name"`
	Fg     string   `json:"foreground"`
	Bg     string   `json:"background"`
	Colors []string `json:"color"`
}

type Scheme struct {
	Name string
	Fg   lipgloss.AdaptiveColor
	Bg   lipgloss.AdaptiveColor
	C    []lipgloss.AdaptiveColor
}

func init() {
	schemes := map[string]raw{}
	err := json.Unmarshal(data, &schemes)
	if err != nil {
		panic(err)
	}

	for name, scheme := range schemes {
		Schemes[name] = &Scheme{
			Name: name,
			Fg:   lipgloss.AdaptiveColor{Dark: scheme.Fg, Light: scheme.Fg},
			Bg:   lipgloss.AdaptiveColor{Dark: scheme.Bg, Light: scheme.Bg},
			C: []lipgloss.AdaptiveColor{
				{Dark: scheme.Colors[0], Light: scheme.Colors[8]},
				{Dark: scheme.Colors[1], Light: scheme.Colors[9]},
				{Dark: scheme.Colors[2], Light: scheme.Colors[10]},
				{Dark: scheme.Colors[3], Light: scheme.Colors[11]},
				{Dark: scheme.Colors[4], Light: scheme.Colors[12]},
				{Dark: scheme.Colors[5], Light: scheme.Colors[13]},
				{Dark: scheme.Colors[6], Light: scheme.Colors[14]},
				{Dark: scheme.Colors[7], Light: scheme.Colors[15]},
			},
		}
	}
}

func Set(name string) error {
	if _, ok := Schemes[name]; !ok {
		return fmt.Errorf("unknown theme: %s", name)
	}

	CurrentScheme.Store(name)
	return nil
}

func Random() {
	keys := make([]string, 0, len(Schemes))
	for k := range Schemes {
		keys = append(keys, k)
	}

	Set(keys[rand.Intn(len(keys))])
}

func Get() *Scheme {
	return Schemes[CurrentScheme.Load()]
}
