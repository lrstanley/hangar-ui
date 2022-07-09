// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package offset

import (
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestMain(m *testing.M) {
	Initialize()

	testsScan = []scanTestCase{
		{"empty", "", "", nil},
		{"single", "a", "a", nil},
		{"double", "aa", "aa", nil},
		{"triple", "aaa", "aaa", nil},
		{"quad", "aaaa", "aaaa", nil},
		{"lipgloss-empty", testStyle.Render(""), testStyle.Render(""), nil},
		{"lipgloss-basic", testStyle.Render("testing"), testStyle.Render("testing"), nil},
		{"lipgloss-basic-start", "a" + testStyle.Render("testing"), "a" + testStyle.Render("testing"), nil},
		{"lipgloss-basic-end", testStyle.Render("testing") + "a", testStyle.Render("testing") + "a", nil},
		{"lipgloss-basic-start-end", "a" + testStyle.Render("testing") + "a", "a" + testStyle.Render("testing") + "a", nil},
		{"lipgloss-basic-between", testStyle.Render("testing") + "a" + testStyle.Render("testing"), testStyle.Render("testing") + "a" + testStyle.Render("testing"), nil},
		{"id-empty", ID("testing"), "", []string{"testing"}},
		{"id-single-start", "a" + ID("testing"), "a", []string{"testing"}},
		{"id-single-end", ID("testing") + "a", "a", []string{"testing"}},
		{"id-single-start-end", "a" + ID("testing") + "a", "aa", []string{"testing"}},
		{"id-single-between", ID("testing") + "a" + ID("testing"), "a", []string{"testing"}},
		{"id-with-lipgloss-start", testStyle.Render(ID("testing") + "testing"), testStyle.Render("testing"), []string{"testing"}},
		{"id-with-lipgloss-end", testStyle.Render("testing" + ID("testing")), testStyle.Render("testing"), []string{"testing"}},
		{"id-multi-empty", ID("foo") + ID("bar"), "", []string{"foo", "bar"}},
		{"id-multi-start", "a" + ID("foo") + ID("bar"), "a", []string{"foo", "bar"}},
		{"id-multi-end", ID("foo") + ID("bar") + "a", "a", []string{"foo", "bar"}},
		{"id-multi-start-end", "a" + ID("foo") + ID("bar") + "a", "aa", []string{"foo", "bar"}},
	}

	m.Run()
	Close()
}

type scanTestCase struct {
	name string
	in   string
	want string
	ids  []string
}

var (
	testStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#383838")).
			Bold(true).
			Italic(true).
			Blink(true)
	testsScan []scanTestCase
)

func BenchmarkScan(b *testing.B) {
	for _, test := range testsScan {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Scan(test.in)
			}
		})
	}
}

func TestScan(t *testing.T) {
	for _, test := range testsScan {
		t.Run(test.name, func(t *testing.T) {
			got := Scan(test.in)
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
			if len(test.ids) > 0 {
				time.Sleep(15 * time.Millisecond)
				for _, id := range test.ids {
					if xy := Get(id); xy.IsZero() {
						t.Errorf("id %q not found", id)
					}
				}
			}
		})
	}
}

func FuzzScan(f *testing.F) {
	for _, test := range testsScan {
		f.Add(test.in)
		f.Add(test.want)
	}

	f.Fuzz(func(t *testing.T, a string) {
		_ = Scan(a)
	})
}
