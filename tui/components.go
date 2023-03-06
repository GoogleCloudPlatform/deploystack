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
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type component interface {
	render() string
}

type productList []struct {
	item    string
	product string
}

func (p productList) longest(k string) int {
	longest := 0

	for _, v := range p {
		if k == "item" {
			if len(v.item) > longest {
				longest = len(v.item)
			}
		} else {
			if len(v.product) > longest {
				longest = len(v.product)
			}
		}
	}
	return longest
}

type description struct {
	stack *config.Stack
}

func newDescription(stack *config.Stack) description {
	return description{stack: stack}
}

func (d *description) parse() (productList, []string) {
	p := productList{}

	if len(d.stack.Config.Products) > 0 {

		for _, v := range d.stack.Config.Products {
			tmp := struct{ item, product string }{}
			tmp.item = strings.TrimSpace(v.Info)
			tmp.product = strings.TrimSpace(v.Product)
			p = append(p, tmp)
		}

	}

	return p, []string{d.stack.Config.Description}
}

func (d description) render() string {
	doc := strings.Builder{}

	list, additionalText := d.parse()

	columns := []table.Column{
		{Title: "", Width: list.longest("item") + 10},
		{Title: "", Width: list.longest("product") + 10},
	}

	rows := []table.Row{}

	for _, v := range list {
		rows = append(rows, table.Row{
			titleStyle.Render(v.item),
			strong.Render(v.product),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(list)),
	)

	t.SetStyles(tableStyle)

	if len(list) > 0 {
		doc.WriteString(normal.Render("This process will install the following resources:"))
		doc.WriteString(t.View())
		doc.WriteString("\n\n")
	}

	for _, v := range additionalText {
		doc.WriteString(normal.Render(v))
		doc.WriteString("\n\n")
	}

	doc.WriteString(normal.Render("It's going to take around "))
	doc.WriteString(strong.Render(strconv.Itoa(d.stack.Config.Duration)))

	if d.stack.Config.Duration == 1 {
		doc.WriteString(normal.Render(" minute."))
	} else {
		doc.WriteString(normal.Render(" minutes."))
	}
	doc.WriteString("\n\n")

	if len(d.stack.Config.DocumentationLink) > 0 {
		doc.WriteString(normal.Render("If you would like more information about this stack, "))
		doc.WriteString(normal.Render("please read the documentation at: "))
		doc.WriteString(url.Render(d.stack.Config.DocumentationLink))
		doc.WriteString("\n\n")
	}

	return doc.String()
}

type errorAlert struct {
	err errMsg
}

func (e errorAlert) Render() string {
	sb := strings.Builder{}

	height := len(e.err.Error()) / width
	style := errorAlertStyle.Copy()
	style.Height(height)

	sb.WriteString("\n")
	sb.WriteString(boldAlert.Render("There was an error!"))
	sb.WriteString("\n")
	if e.err.usermsg != "" {
		sb.WriteString(e.err.usermsg)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString("Details: \n")
	sb.WriteString(e.err.Error())
	sb.WriteString("\n")
	sb.WriteString("\n")

	if !e.err.quit {
		sb.WriteString("You can exit the program by typing ")
		sb.WriteString(cmdStyle.Render("ctr+c."))
	}

	if e.err.target != "" {
		text := " Press the Enter Key to go back and change choice "

		if e.err.target == "quit" {
			text = " Press the Enter Key exit "
		}

		sb.WriteString("\n")
		sb.WriteString("\n")
		sb.WriteString(bodyStyle.Render(promptStyle.Render(text)))
		sb.WriteString("\n")
	}

	return style.Render(sb.String())
}

type header struct {
	title    string
	subtitle string
}

func newHeader(title, subtitle string) header {
	return header{title: title, subtitle: subtitle}
}

func (h header) render() string {
	doc := strings.Builder{}

	content := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("%s%s%s", textColors.code("bright cyan"), titleStyle.Render(h.title), clear),
		subTitleStyle.Render(h.subtitle),
	)

	doc.WriteString(headerStyle.Render(content))

	doc.WriteString("\n\n")
	return doc.String()
}

func drawProgress(percent int) string {

	sb := strings.Builder{}

	label := "   Progress "
	sb.WriteString(textStyle.Render(label))

	totalWidth := hardWidthLimit - len(label)
	completeLength := int(float32(totalWidth) * (float32(percent) / float32(100)))
	pendingLength := totalWidth - completeLength

	comp := strings.Builder{}
	for i := 0; i < completeLength; i++ {
		comp.WriteString("█")
	}

	pend := strings.Builder{}
	for i := 0; i < pendingLength; i++ {
		pend.WriteString("░")
	}

	sb.WriteString(completeStyle.Render(comp.String()))
	sb.WriteString(pendingStyle.Render(pend.String()))

	return sb.String()
}

type settingsTable struct {
	stack *config.Stack
}

func newSettingsTable(s *config.Stack) settingsTable {
	return settingsTable{stack: s}
}

func (s settingsTable) render() string {
	doc := strings.Builder{}
	wSetting := 0
	wValue := 0

	s.stack.Settings.Sort()

	rows := []table.Row{}

	if s := s.stack.Settings.Find("stack_name"); s != nil && len(s.Value) > 0 {
		rows = append(rows, table.Row{
			titleStyle.Render("Stack Name"),
			strong.Render(s.Value),
		})
	}

	if s := s.stack.Settings.Find("project_name"); s != nil && len(s.Value) > 0 {
		rows = append(rows, table.Row{
			titleStyle.Render("Project Name"),
			strong.Render(s.Value),
		})
	}

	if s := s.stack.Settings.Find("project_id"); s != nil && len(s.Value) > 0 {
		rows = append(rows, table.Row{
			titleStyle.Render("Project ID"),
			strong.Render(s.Value),
		})
	}

	if s := s.stack.Settings.Find("project_number"); s != nil && len(s.Value) > 0 {
		rows = append(rows, table.Row{
			titleStyle.Render("Project Number"),
			strong.Render(s.Value),
		})
	}

	for _, setting := range s.stack.Settings {

		rawValue := setting.TFvarsValue()
		rawValue = strings.Trim(rawValue, "\"")
		value := strong.Render(strings.TrimSpace(rawValue))

		if len(rawValue) > 45 {
			value = strong.Render(rawValue[:45] + "...")
		}

		if len(setting.Name) > wSetting {
			wSetting = len(setting.Name)
		}

		if len(value) > wValue {
			wValue = len(value)
		}

		if setting.Name == "project_id" ||
			setting.Name == "project_number" ||
			setting.Name == "project_name" ||
			setting.Name == "stack_name" {
			continue
		}
		if len(value) < 1 {
			continue
		}

		nameRaw := strings.TrimSpace(setting.Name)
		nameRaw = strings.ReplaceAll(nameRaw, "_", " ")
		nameRaw = strings.ReplaceAll(nameRaw, "-", " ")
		formatted := cases.Title(language.English).String(nameRaw)
		rows = append(rows, table.Row{titleStyle.Render(formatted), value})

	}

	columns := []table.Column{
		{Title: "Setting", Width: 35},
		{Title: "Value", Width: 55},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(s.stack.Settings)),
	)

	t.SetStyles(tableStyle)
	doc.WriteString("\n")
	doc.WriteString(t.View())
	doc.WriteString("\n")

	return doc.String()
}

type textBlock string

func (t textBlock) render() string    { return string(t) }
func newTextBlock(s string) textBlock { return textBlock(s) }
