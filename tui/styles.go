// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

const (
	// TERMCYAN is the terminal code for cyan text
	TERMCYAN = "\033[0;36m"
	// TERMCYANB is the terminal code for bold cyan text
	TERMCYANB = "\033[1;36m"
	// TERMCYANREV is the terminal code for black on cyan text
	TERMCYANREV = "\u001b[46m"
	// TERMRED is the terminal code for red text
	TERMRED = "\033[0;31m"
	// TERMREDB is the terminal code for bold red text
	TERMREDB = "\033[1;31m"
	// TERMREDREV is the terminal code for black on red text
	TERMREDREV = "\033[41m"
	// TERMCLEAR is the terminal code for the clear out color text
	TERMCLEAR = "\033[0m"
	// TERMCLEARSCREEN is the terminal code for clearning the whole screen.
	TERMCLEARSCREEN = "\033[2J"
	// TERMGREY is the terminal code for grey text
	TERMGREY = "\033[1;30m"
)

var colors = ansiColors{
	"blank":   ansiColor{id: -1},
	"black":   ansiColor{id: 0},
	"red":     ansiColor{id: 1},
	"green":   ansiColor{id: 2},
	"yellow":  ansiColor{id: 3},
	"blue":    ansiColor{id: 4},
	"magenta": ansiColor{id: 5},
	"cyan":    ansiColor{id: 6},
	"white":   ansiColor{id: 7},
	"grey":    ansiColor{id: 8},
	"gray":    ansiColor{id: 8},
}

var clear = "\033[0m"

type ansiColor struct {
	id int
}

func (a ansiColor) bright() string {
	if a.id == -1 {
		return ""
	}
	return fmt.Sprintf("\033[1;3%dm", a.id)
}

func (a ansiColor) regular() string {
	if a.id == -1 {
		return ""
	}
	return fmt.Sprintf("\033[0;3%dm", a.id)
}

func (a ansiColor) bold() string {
	if a.id == -1 {
		return ""
	}
	return fmt.Sprintf("\033[1;3%dm", a.id)
}

func (a ansiColor) underline() string {
	if a.id == -1 {
		return ""
	}
	return fmt.Sprintf("\033[4;3%dm", a.id)
}

func (a ansiColor) background() string {
	if a.id == -1 {
		return ""
	}
	return fmt.Sprintf("\033[1;4%dm", a.id)
}

type ansiColors map[string]ansiColor

func (a ansiColors) get(s string) ansiColor {

	if s == "copy" {
		if termenv.HasDarkBackground() {
			return colors["white"]
		}
		return colors["black"]
	}
	return colors[s]
}

type dsStyle struct {
	style      lipgloss.Style
	foreground ansiColor
	background ansiColor
	bright     bool
	underline  bool
	bold       bool
}

func (d dsStyle) Render(s string) string {

	startFg := d.foreground.regular()
	if d.bright {
		startFg = d.foreground.bright()
	}
	if d.bold {
		startFg = d.foreground.bold()
	}
	if d.underline {
		startFg = d.foreground.underline()
	}
	startBg := d.background.background()

	content := d.style.Render(s)

	return fmt.Sprintf("%s%s%s%s", startBg, startFg, content, clear)
}

func newDsStyle() dsStyle {
	r := dsStyle{style: lipgloss.NewStyle()}
	r.foreground = colors.get("copy")
	r.background = ansiColor{id: -1}
	return r
}

func (d dsStyle) Foreground(a ansiColor) dsStyle {
	d.foreground = a
	return d
}

func (d dsStyle) Background(a ansiColor) dsStyle {
	d.background = a
	return d
}

func (d dsStyle) Bright(t bool) dsStyle {
	d.bright = t
	return d
}

func (d dsStyle) Bold(t bool) dsStyle {
	d.bold = t
	return d
}

func (d dsStyle) Underline(t bool) dsStyle {
	d.underline = t
	return d
}

func (d dsStyle) Width(i int) dsStyle {
	d.style = d.style.Width(i)
	return d
}

func (d dsStyle) Height(i int) dsStyle {
	d.style = d.style.Height(i)
	return d
}

func (d dsStyle) MarginLeft(i int) dsStyle {
	d.style = d.style.MarginLeft(i)
	return d
}
func (d dsStyle) MarginTop(i int) dsStyle {
	d.style = d.style.MarginTop(i)
	return d
}
func (d dsStyle) MarginRight(i int) dsStyle {
	d.style = d.style.MarginRight(i)
	return d
}
func (d dsStyle) MarginBottom(i int) dsStyle {
	d.style = d.style.MarginBottom(i)
	return d
}

func (d dsStyle) Margin(i ...int) dsStyle {
	d.style = d.style.Margin(i...)
	return d
}

func (d dsStyle) PaddingLeft(i int) dsStyle {
	d.style = d.style.PaddingLeft(i)
	return d
}
func (d dsStyle) PaddingTop(i int) dsStyle {
	d.style = d.style.PaddingTop(i)
	return d
}
func (d dsStyle) PaddingRight(i int) dsStyle {
	d.style = d.style.PaddingRight(i)
	return d
}
func (d dsStyle) PaddingBottom(i int) dsStyle {
	d.style = d.style.PaddingBottom(i)
	return d
}

func (d dsStyle) Padding(i ...int) dsStyle {
	d.style = d.style.Padding(i...)
	return d
}

func (d dsStyle) MaxWidth(i int) dsStyle {
	d.style = d.style.MaxWidth(i)
	return d
}

func (d dsStyle) Italic(t bool) dsStyle {
	d.style = d.style.Italic(t)
	return d
}

func (d dsStyle) Copy() dsStyle {
	r := dsStyle{}

	r.style = d.style.Copy()
	r.foreground = d.foreground
	r.background = d.background
	r.bright = d.bright
	r.underline = d.underline
	r.bold = d.bold

	return r
}

func (d dsStyle) BorderLeft(t bool) dsStyle {
	d.style = d.style.BorderLeft(t)
	return d
}
func (d dsStyle) BorderTop(t bool) dsStyle {
	d.style = d.style.BorderTop(t)
	return d
}
func (d dsStyle) BorderRight(t bool) dsStyle {
	d.style = d.style.BorderRight(t)
	return d
}
func (d dsStyle) BorderBottom(t bool) dsStyle {
	d.style = d.style.BorderBottom(t)
	return d
}

func (d dsStyle) BorderStyle(b lipgloss.Border) dsStyle {
	d.style = d.style.Border(b)
	return d
}

func (d dsStyle) BorderForeground(b lipgloss.TerminalColor) dsStyle {
	d.style = d.style.BorderForeground(b)
	return d
}

var (
	width          = 100
	hardWidthLimit = width
	gray           = lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	grayWeak       = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}
	simClearColor  = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}
	highlight      = lipgloss.AdaptiveColor{Light: "6", Dark: "14"}
	basicText      = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}
	alert          = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}
	completeColor  = lipgloss.AdaptiveColor{Light: "8", Dark: "8"}
	pendingColor   = lipgloss.AdaptiveColor{Light: "6", Dark: "6"}

	strong = newDsStyle().
		Foreground(colors.get("cyan"))

	normal = newDsStyle().
		Foreground(colors.get("copy"))

	url = newDsStyle().
		Foreground(colors.get("cyan")).
		Underline(true)

	titleStyle = newDsStyle().
			Bold(true).
			Foreground(colors.get("copy"))

	purchaseStyle = newDsStyle().
			Bold(true).
			Foreground(colors.get("red")).
			Background(colors.get("grey"))

	subTitleStyle = newDsStyle().
			MaxWidth(hardWidthLimit).
			Bold(false).
			Italic(true).
			Foreground(colors.get("copy"))

	headerCopyStyle = newDsStyle().
			MaxWidth(hardWidthLimit)

	headerStyle = newDsStyle().
			MarginLeft(0).
			MarginRight(0).
			Padding(0, 3).
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false).
			BorderBottom(true).
			MaxWidth(hardWidthLimit).
			BorderForeground(gray).
			Width(width)

	cursorPromptStyle = newDsStyle().
				Foreground(colors.get("cyan"))

	bodyStyle = newDsStyle().
			MarginLeft(0).
			MarginRight(0).
			Padding(0, 3).
			Foreground(colors.get("copy")).
			Width(width).
			MaxWidth(hardWidthLimit)

	docStyle = newDsStyle().
			Foreground(colors.get("copy")).
			Padding(0, 2)

	promptStyle = newDsStyle().
			Bold(true).
			Background(colors.get("cyan"))

	alertStyle = bodyStyle.Copy().
			Foreground(colors.get("red"))

	alertStrongStyle = bodyStyle.Copy().
				Foreground(colors.get("red")).
				PaddingLeft(3).Bold(true)

	instructionStyle = newDsStyle().
				PaddingLeft(3)

	textStyle = newDsStyle().
			Foreground(colors.get("copy"))

	textInputDefaultStyle = newDsStyle().
				Foreground(colors.get("cyan"))

	tableStyle = table.DefaultStyles()

	inputText = bodyStyle.Copy().
			Foreground(colors.get("cyan"))

	componentStyle = newDsStyle().
			PaddingLeft(1).
			MarginLeft(0)

	billingDisabledStyle = newDsStyle().
				Foreground(colors.get("grey"))

	itemStyle = newDsStyle().
			PaddingLeft(4)

	selectedItemStyle = newDsStyle().
				PaddingLeft(2).
				Background(colors.get("cyan")).
				Foreground(colors.get("white"))

	paginationStyle = list.DefaultStyles().
			PaginationStyle.PaddingLeft(4)

	helpStyle = list.DefaultStyles().
			HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(grayWeak)

	quitTextStyle = newDsStyle().
			Margin(1, 0, 2, 4)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(highlight)

	textInputPrompt = helpStyle.Copy().
			PaddingLeft(3)

	completeStyle = newDsStyle().Foreground(colors.get("cyan"))

	pendingStyle = newDsStyle().Foreground(colors.get("white"))
)

func init() {
	width, _, _ = term.GetSize(int(os.Stdout.Fd()))

	tableStyle.Header.
		BorderStyle(lipgloss.HiddenBorder()).
		BorderBottom(true).
		Bold(false)

	tableStyle.Selected.
		Foreground(basicText).
		Bold(false)
	tableStyle.Cell.Foreground(basicText).
		Padding(0)
	tableStyle.Header.Padding(0)
}
