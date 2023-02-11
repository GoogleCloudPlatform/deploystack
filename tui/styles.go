package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	width          = 100
	hardWidthLimit = width
	gray           = lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	grayWeak       = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}
	simClearColor  = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}
	promptBGround  = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}
	highlight      = lipgloss.AdaptiveColor{Light: "6", Dark: "14"}
	basicText      = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}
	alert          = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}

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

	headerCopyStyle = lipgloss.NewStyle().
			MaxWidth(hardWidthLimit)

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

	cursorPromptStyle = lipgloss.NewStyle().
				Foreground(highlight).
				Background(simClearColor)

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

	alertStrongStyle = bodyStyle.Copy().
				Foreground(alert).
				PaddingLeft(3).Bold(true)

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

	quitTextStyle = lipgloss.NewStyle().
			Margin(1, 0, 2, 4)

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
