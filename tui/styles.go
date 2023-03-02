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
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

type ansi16colors []ansi16color

func (a ansi16colors) code(s string) string {
	for _, v := range a {
		if s == v.name {
			return v.colorCode
		}
	}
	return ""
}

func (a ansi16colors) codeByID(i int) string {
	for _, v := range a {
		if i == v.id {
			return v.colorCode
		}
	}
	return ""
}

func (a ansi16colors) color(s string) ansi16color {
	for _, v := range a {
		if strings.TrimSpace(s) == strings.TrimSpace(v.name) {
			return v
		}
	}
	return ansi16color{}
}

func (a ansi16colors) colorByID(i int) ansi16color {
	for _, v := range a {
		if i == v.id {
			return v
		}
	}
	return ansi16color{}
}

type ansi16color struct {
	id        int
	name      string
	colorCode string
}

var textColors = ansi16colors{
	{id: 0, name: "black", colorCode: "\033[0;30m"},
	{id: 1, name: "red", colorCode: "\033[0;31m"},
	{id: 2, name: "green", colorCode: "\033[0;32m"},
	{id: 3, name: "yellow", colorCode: "\033[0;33m"},
	{id: 4, name: "blue", colorCode: "\033[0;34m"},
	{id: 5, name: "magenta", colorCode: "\033[0;35m"},
	{id: 6, name: "cyan", colorCode: "\033[0;36m"},
	{id: 7, name: "white", colorCode: "\033[0;37m"},
	{id: 8, name: "bright black", colorCode: "\033[1;30m"},
	{id: 9, name: "bright red", colorCode: "\033[1;31m"},
	{id: 10, name: "bright green", colorCode: "\033[1;32m"},
	{id: 11, name: "bright yellow", colorCode: "\033[1;33m"},
	{id: 12, name: "bright blue", colorCode: "\033[1;34m"},
	{id: 13, name: "bright magenta", colorCode: "\033[1;35m"},
	{id: 14, name: "bright cyan", colorCode: "\033[1;36m"},
	{id: 15, name: "bright white", colorCode: "\033[1;37m"},
	{id: 8, name: "dark grey", colorCode: "\033[1;30m"},
	{id: 8, name: "dark gray", colorCode: "\033[1;30m"},
	{id: 7, name: "light grey", colorCode: "\033[0;37m"},
	{id: 7, name: "light gray", colorCode: "\033[0;37m"},
}
var backgroundColors = ansi16colors{
	{id: -1, name: "blank", colorCode: ""},
	{id: 0, name: "black", colorCode: "\033[0;40m"},
	{id: 1, name: "red", colorCode: "\033[0;41m"},
	{id: 2, name: "green", colorCode: "\033[0;42m"},
	{id: 3, name: "yellow", colorCode: "\033[0;43m"},
	{id: 4, name: "blue", colorCode: "\033[0;44m"},
	{id: 5, name: "magenta", colorCode: "\033[0;45m"},
	{id: 6, name: "cyan", colorCode: "\033[0;46m"},
	{id: 7, name: "white", colorCode: "\033[0;47m"},
	{id: 8, name: "bold on black", colorCode: "\033[0;40m"},
	{id: 9, name: "bold on red", colorCode: "\033[0;41m"},
	{id: 10, name: "bold on green", colorCode: "\033[0;42m"},
	{id: 11, name: "bold on yellow", colorCode: "\033[0;43m"},
	{id: 12, name: "bold on blue", colorCode: "\033[0;44m"},
	{id: 13, name: "bold on magenta", colorCode: "\033[0;45m"},
	{id: 14, name: "bold on cyan", colorCode: "\033[0;46m"},
	{id: 15, name: "bold on white", colorCode: "\033[0;47m"},
}

type dsAdaptiveColor struct {
	light             ansi16color
	dark              ansi16color
	blankOnCloudShell bool
}

func (a dsAdaptiveColor) code() string {
	if a.blankOnCloudShell && os.Getenv("GOOGLE_CLOUD_SHELL") != "" {
		return clear
	}

	if termenv.HasDarkBackground() {
		return a.dark.colorCode
	}
	return a.light.colorCode
}

var clear = "\033[0m"

type dsStyle struct {
	style      lipgloss.Style
	foreground dsAdaptiveColor
	background dsAdaptiveColor
	bright     bool
	underline  bool
	bold       bool
}

func (d dsStyle) Render(s string) string {

	startFg := d.foreground.code()
	if d.underline {
		// Replace the right character with the underline trigger
		sl := strings.Split(startFg, "")
		sl[2] = "4"
		startFg = strings.Join(sl, "")
	}
	startBg := d.background.code()
	content := d.style.Render(s)

	return fmt.Sprintf("%s%s%s%s", startFg, startBg, content, clear)
}

func newDsStyle() dsStyle {
	blankBG := backgroundColors.color("blank")
	black := textColors.color("black")
	white := textColors.color("light grey")

	r := dsStyle{style: lipgloss.NewStyle()}
	r.foreground = dsAdaptiveColor{light: black, dark: white, blankOnCloudShell: true}
	r.background = dsAdaptiveColor{light: blankBG, dark: blankBG}
	return r
}

func (d dsStyle) Foreground(a dsAdaptiveColor) dsStyle {
	d.foreground = a
	return d
}

func (d dsStyle) Background(a dsAdaptiveColor) dsStyle {
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
	lgbasicText    = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}
	lggray         = lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	lggrayWeak     = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}
	lgalert        = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}

	gray          = dsAdaptiveColor{light: textColors.color("white"), dark: textColors.color("dark grey")}
	grayWeak      = dsAdaptiveColor{light: textColors.color("dark grey"), dark: textColors.color("white")}
	simClearColor = dsAdaptiveColor{light: textColors.color("bright white"), dark: textColors.colorByID(0)}
	highlight     = dsAdaptiveColor{light: textColors.color("cyan"), dark: textColors.color("bright cyan")}
	basicText     = dsAdaptiveColor{light: textColors.color("black"), dark: textColors.color("light grey"), blankOnCloudShell: true}
	alert         = dsAdaptiveColor{light: textColors.color("red"), dark: textColors.color("bright red")}
	completeColor = dsAdaptiveColor{light: textColors.color("dark grey"), dark: textColors.color("dark grey")}
	pendingColor  = dsAdaptiveColor{light: textColors.color("cyan"), dark: textColors.color("bright cyan")}
	highlightBG   = dsAdaptiveColor{light: backgroundColors.color("bold on cyan"), dark: backgroundColors.color("cyan")}

	strong = newDsStyle().
		Foreground(highlight)

	normal = newDsStyle().
		Foreground(basicText)

	url = newDsStyle().
		Foreground(highlight).
		Underline(true)

	titleStyle = newDsStyle().
			Bold(true).
			Foreground(basicText)

	purchaseStyle = newDsStyle().
			Bold(true).
			Foreground(alert).
			Background(gray)

	subTitleStyle = newDsStyle().
			MaxWidth(hardWidthLimit).
			Bold(false).
			Foreground(basicText)

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
			BorderForeground(lggray).
			Width(width)

	cursorPromptStyle = newDsStyle().
				Foreground(highlight)

	bodyStyle = newDsStyle().
			MarginLeft(0).
			MarginRight(0).
			Padding(0, 3).
			Foreground(basicText).
			Width(width).
			MaxWidth(hardWidthLimit)

	docStyle = newDsStyle().
			Foreground(basicText).
			Padding(0, 2)

	promptStyle = newDsStyle().
			Bold(true).
			Background(highlightBG).
			Foreground(dsAdaptiveColor{light: textColors.color("white"), dark: textColors.color("white")})

	alertStyle = bodyStyle.Copy().
			Foreground(alert)

	alertStrongStyle = bodyStyle.Copy().
				Foreground(alert).
				PaddingLeft(3).Bold(true)

	instructionStyle = newDsStyle().
				PaddingLeft(3)

	textStyle = newDsStyle().
			Foreground(basicText)

	textInputDefaultStyle = newDsStyle().
				Foreground(highlight)

	tableStyle = table.DefaultStyles()

	inputText = bodyStyle.Copy().
			Foreground(highlight)

	componentStyle = newDsStyle().
			PaddingLeft(1).
			MarginLeft(0)

	billingDisabledStyle = newDsStyle().
				Foreground(gray)

	itemStyle = newDsStyle().
			PaddingLeft(4)

	selectedItemStyle = newDsStyle().
				PaddingLeft(2).
				Background(highlightBG).
				Foreground(basicText)

	paginationStyle = list.DefaultStyles().
			PaginationStyle.PaddingLeft(4)

	helpStyle = list.DefaultStyles().
			HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(lggrayWeak)

	quitTextStyle = newDsStyle().
			Margin(1, 0, 2, 4)

	spinnerStyle = newDsStyle().Foreground(highlight)

	textInputPrompt = helpStyle.Copy().
			PaddingLeft(3)

	completeStyle = newDsStyle().Foreground(highlight)

	pendingStyle = newDsStyle().Foreground(grayWeak)

	errorAlertStyle = lipgloss.NewStyle().
			Width(100).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lgalert).
			PaddingLeft(3).
			Foreground(lggrayWeak)

	boldAlert = lipgloss.NewStyle().Bold(true).Foreground(lgalert)
	cmdStyle  = lipgloss.NewStyle().Background(lggrayWeak).Foreground(lgalert)
)

func init() {

	width, _, _ = term.GetSize(int(os.Stdout.Fd()))

	tableStyle.Header.
		BorderStyle(lipgloss.HiddenBorder()).
		BorderBottom(true).
		Bold(false)

	tableStyle.Selected.
		Foreground(lgbasicText).
		Bold(false)
	tableStyle.Cell.Foreground(lgbasicText).
		Padding(0)
	tableStyle.Header.Padding(0)
}
