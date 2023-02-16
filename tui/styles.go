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
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
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

//			Normal  Bright
// Black	0		8
// Red		1		9
// Green	2		10
// Yellow	3		11
// Blue		4		12
// Purple	5		13
// Cyan		6		14
// White	7		15

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

	strong = lipgloss.NewStyle().
		Foreground(highlight)

	normal = lipgloss.NewStyle().
		Foreground(basicText)

	url = lipgloss.NewStyle().
		Foreground(highlight).
		Underline(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(basicText)

	purchaseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(alert).
			Background(grayWeak)

	subTitleStyle = lipgloss.NewStyle().
			MaxWidth(hardWidthLimit).
			Bold(false).
			Italic(true).
			Foreground(basicText)

	margins = lipgloss.NewStyle().
		MarginLeft(0).
		MarginRight(0).
		Padding(0, 3)

	headerStyle = margins.Copy().
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false).
			BorderBottom(true).
			MaxWidth(hardWidthLimit).
			BorderForeground(gray).
			Width(width)

	bodyStyle = margins.Copy().
			Foreground(basicText).
			Width(width).
			MaxWidth(hardWidthLimit)

	docStyle = lipgloss.NewStyle().
			Foreground(basicText).
			Padding(0, 2)

	promptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(simClearColor).
			Background(highlight)

	alertStyle = bodyStyle.Copy().
			Foreground(alert)

	instructionStyle = lipgloss.NewStyle().
				PaddingLeft(3)

	textStyle = lipgloss.NewStyle().
			Foreground(basicText)

	textInputDefaultStyle = lipgloss.NewStyle().
				Foreground(highlight)

	tableStyle = table.DefaultStyles()

	inputText = bodyStyle.Copy().
			Foreground(highlight)

	componentStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			MarginLeft(0)

	billingDisabledStyle = lipgloss.NewStyle().
				Foreground(gray)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(simClearColor).
				Background(highlight)

	paginationStyle = list.DefaultStyles().
			PaginationStyle.PaddingLeft(4)

	helpStyle = list.DefaultStyles().
			HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(grayWeak)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(highlight)

	textInputPrompt = helpStyle.Copy().
			PaddingLeft(3)
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
