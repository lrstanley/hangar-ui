// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import "github.com/charmbracelet/lipgloss"

type ThemeConfig struct {
	Fg lipgloss.AdaptiveColor
	Bg lipgloss.AdaptiveColor

	SuccessFg lipgloss.AdaptiveColor
	FailureFg lipgloss.AdaptiveColor

	ViewBorderActiveFg   lipgloss.AdaptiveColor
	ViewBorderInactiveFg lipgloss.AdaptiveColor
	ViewBorderBg         lipgloss.AdaptiveColor

	NavActiveFg   lipgloss.AdaptiveColor
	NavActiveBg   lipgloss.AdaptiveColor
	NavInactiveFg lipgloss.AdaptiveColor
	NavInactiveBg lipgloss.AdaptiveColor

	InputFg            lipgloss.AdaptiveColor
	InputCursorFg      lipgloss.AdaptiveColor
	InputPlaceholderFg lipgloss.AdaptiveColor

	TitleFg lipgloss.AdaptiveColor
	TitleBg lipgloss.AdaptiveColor

	StatusBarTargetFg  lipgloss.AdaptiveColor
	StatusBarTargetBg  lipgloss.AdaptiveColor
	StatusBarFg        lipgloss.AdaptiveColor
	StatusBarKeyDescFg lipgloss.AdaptiveColor
	StatusBarBg        lipgloss.AdaptiveColor
	StatusBarURLFg     lipgloss.AdaptiveColor
	StatusBarURLBg     lipgloss.AdaptiveColor
	StatusBarLogoFg    lipgloss.AdaptiveColor
	StatusBarLogoBg    lipgloss.AdaptiveColor
}

// Theme is the default theme for the application.
var Theme *ThemeConfig

// SetTheme allows changing the theme of the application.
func SetTheme(style string) {
	switch style {
	case "default":
		Theme = &ThemeConfig{
			Fg: lipgloss.AdaptiveColor{Dark: "#98D1CE", Light: "#98D1CE"},
			Bg: lipgloss.AdaptiveColor{Dark: "#0A0F14", Light: "#0A0F14"},

			SuccessFg: lipgloss.AdaptiveColor{Dark: "#69ff94", Light: "#69ff94"},
			FailureFg: lipgloss.AdaptiveColor{Dark: "#ff6e6e", Light: "#ff6e6e"},

			ViewBorderActiveFg:   lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"},
			ViewBorderInactiveFg: lipgloss.AdaptiveColor{Dark: "#D9DCCF", Light: "#D9DCCF"},
			ViewBorderBg:         lipgloss.AdaptiveColor{Dark: "#0A0F14", Light: "#0A0F14"},

			NavActiveFg:   lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			NavActiveBg:   lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},
			NavInactiveFg: lipgloss.AdaptiveColor{Dark: "#D9DCCF", Light: "#D9DCCF"},
			NavInactiveBg: lipgloss.AdaptiveColor{Dark: "#353533", Light: "#353533"},

			InputFg:            lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			InputCursorFg:      lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			InputPlaceholderFg: lipgloss.AdaptiveColor{Dark: "#525252", Light: "#525252"},

			TitleFg: lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			TitleBg: lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},

			StatusBarTargetFg:  lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			StatusBarTargetBg:  lipgloss.AdaptiveColor{Dark: "#CC6699", Light: "#CC6699"},
			StatusBarFg:        lipgloss.AdaptiveColor{Light: "#C1C6B2", Dark: "#C1C6B2"},
			StatusBarKeyDescFg: lipgloss.AdaptiveColor{Light: "#c4c4c4", Dark: "#c4c4c4"},
			StatusBarBg:        lipgloss.AdaptiveColor{Light: "#353533", Dark: "#353533"},
			StatusBarURLFg:     lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			StatusBarURLBg:     lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"},
			StatusBarLogoFg:    lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#FFFFFF"},
			StatusBarLogoBg:    lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},
		}
	}
}
